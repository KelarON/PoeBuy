package bot

import (
	"context"
	"net/http"
	"poebuy/config"
	"poebuy/modules/connections"
	"poebuy/modules/connections/headers"
	"poebuy/modules/watchers"
	"poebuy/utils"
	"time"
)

const MAX_TRADES_IN_QUEUE = 5

// App is the main application struct
type Bot struct {
	config              *config.Config
	watchers            map[string]*watchers.ItemWatcher
	errChan             chan error
	logger              *utils.Logger
	UpdateCheckmarkFunc func(int)
	hideoutVisitsQueue  *utils.AsyncQueue[string]
	visitDelay          *int
	visitorCtx          context.Context
	visitorCancel       context.CancelFunc
}

// Init initializes the application
func NewBot(cfg *config.Config, logger *utils.Logger) (*Bot, error) {

	visitorCtx, visitorCancel := context.WithCancel(context.Background())

	bot := &Bot{
		errChan:            make(chan error),
		config:             cfg,
		logger:             logger,
		watchers:           make(map[string]*watchers.ItemWatcher),
		hideoutVisitsQueue: utils.NewAsyncQueue[string](MAX_TRADES_IN_QUEUE),
		visitDelay:         &cfg.Trade.VisitDelay,
		visitorCtx:         visitorCtx,
		visitorCancel:      visitorCancel,
	}

	go bot.errorWriter()
	go bot.startVisitor(bot.visitorCtx)

	cfg.DefineErrorChannel(bot.errChan)

	return bot, nil
}

func (bot *Bot) WatchItem(code string, delay int64) error {

	var index int

	for i := range bot.config.Trade.Links {
		if bot.config.Trade.Links[i].Code == code {
			index = i
			break
		}
	}

	watcher, err := watchers.NewItemWatcher(
		bot.config.General.Poesessid,
		bot.config.Trade.League,
		code,
		bot.errChan,
		delay,
		index,
		bot.UpdateCheckmarkFunc,
		bot.hideoutVisitsQueue,
	)
	if err != nil {
		return err
	}

	bot.watchers[code] = watcher

	go watcher.Watch()

	return nil
}

func (bot *Bot) StopWatcher(code string) {
	bot.watchers[code].Stop()
	delete(bot.watchers, code)
}

// Stop closes the application and cleans up
func (bot *Bot) StopAllWatchers() {

	for _, watcher := range bot.watchers {
		watcher.Stop()
	}

}

func (bot *Bot) errorWriter() {
	for {
		err := <-bot.errChan
		bot.logger.Error(err.Error())
	}
}

func (bot *Bot) startVisitor(ctx context.Context) {
	whisper := connections.NewWhisper(&http.Client{}, headers.GetWhisperHeaders(bot.config.General.Poesessid))
	for {
		token := bot.hideoutVisitsQueue.Pop(ctx)
		if token == nil {
			break
		}
		err := whisper.Whisper(*token)
		if err != nil {
			bot.errChan <- err
			continue
		}
		time.Sleep(time.Second * time.Duration(*bot.visitDelay))
	}
}

func (bot *Bot) RestartVisitor() {

	bot.visitorCancel()
	visitorCtx, visitorCancel := context.WithCancel(context.Background())
	bot.visitorCtx = visitorCtx
	bot.visitorCancel = visitorCancel
	go bot.startVisitor(bot.visitorCtx)
}
