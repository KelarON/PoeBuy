package watchers

import (
	"fmt"
	"net/http"
	"poebuy/modules/connections"
	"poebuy/modules/connections/headers"
	"poebuy/modules/connections/models"
	"poebuy/utils"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	MAX_PROCESSABLE_ITEMS = 3
)

type ItemWatcher struct {
	wsConnection        *websocket.Conn
	fetcher             *connections.Fetcher
	whisper             *connections.Whisper
	code                string
	errChan             chan error
	working             bool
	delay               time.Duration
	readReady           bool
	index               int
	updateCheckmarkFunc func(int)
	hideoutVisitsQueue  *utils.AsyncQueue[string]
}

func NewItemWatcher(
	poesseid string,
	league string,
	code string,
	errChan chan error,
	delay int64,
	index int,
	updateCheckmarkFunc func(int),
	hideoutVisitsQueue *utils.AsyncQueue[string],
) (*ItemWatcher, error) {

	client := &http.Client{}

	wsConn, err := connections.NewWSConnection(poesseid, league, code)
	if err != nil {
		return nil, err
	}

	watcher := &ItemWatcher{
		wsConnection:        wsConn,
		fetcher:             connections.NewFetcher(client, headers.GetFetchitemHeaders(poesseid)),
		whisper:             connections.NewWhisper(client, headers.GetWhisperHeaders(poesseid)),
		code:                code,
		errChan:             errChan,
		working:             false,
		delay:               time.Millisecond * time.Duration(delay),
		readReady:           true,
		index:               index,
		updateCheckmarkFunc: updateCheckmarkFunc,
		hideoutVisitsQueue:  hideoutVisitsQueue,
	}

	return watcher, nil

}

func (w *ItemWatcher) Watch() {

	w.working = true

	if w.delay > 0 {
		go w.delayer()
	}

	for {
		if !w.working {
			return
		}
		var ls models.LivesearchNewItem
		err := w.wsConnection.ReadJSON(&ls)
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				w.errChan <- err
			}
			break
		}
		if (w.delay > 0 && !w.readReady) || ls.Result == "" {
			continue
		}

		itemsInfo, err := w.fetcher.FetchItems([]string{ls.Result}, w.code)
		if err != nil {
			w.errChan <- err
			continue
		}
		if w.delay != 0 && len(itemsInfo) > MAX_PROCESSABLE_ITEMS {
			itemsInfo = itemsInfo[:MAX_PROCESSABLE_ITEMS]
		}

		for _, itemInfo := range itemsInfo {
			if itemInfo.Result[0].Listing.WhisperToken != "" {
				err := w.whisper.Whisper(itemInfo.Result[0].Listing.WhisperToken)
				if err != nil {
					w.errChan <- err
					continue
				}
			} else {
				if itemInfo.Result[0].Listing.HideoutToken != "" {
					token := itemInfo.Result[0].Listing.HideoutToken
					w.hideoutVisitsQueue.Push(&token)
				} else {
					w.errChan <- fmt.Errorf("whisper and hideout token not found for item")
				}
			}
		}

		w.readReady = false
	}

	w.Stop()
}

func (w *ItemWatcher) Stop() {

	w.wsConnection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	w.wsConnection.Close()
	w.working = false
	w.updateCheckmarkFunc(w.index)
}

func (w *ItemWatcher) delayer() {
	for {
		if !w.working {
			return
		}
		w.readReady = true
		time.Sleep(w.delay)
	}
}
