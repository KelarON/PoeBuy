package connections

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

type Whisper struct {
	Client *http.Client
	Header http.Header
}

func NewWhisper(client *http.Client, header http.Header) *Whisper {
	return &Whisper{
		Client: client,
		Header: header,
	}
}

func (w *Whisper) Whisper(token string) error {
	jsonBody := []byte(fmt.Sprintf("{\"token\": \"%v\"}", token))
	whisperReq, err := http.NewRequest("POST", "https://www.pathofexile.com/api/trade/whisper", bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("whisper request creation error: %v", err)

	}
	whisperReq.Header = w.Header
	whisperResp, err := w.Client.Do(whisperReq)
	if err != nil {
		return fmt.Errorf("whisper request error: %v", err)

	}
	if whisperResp.StatusCode != 200 {
		log.Println("Whisper error: ", whisperResp.Status, whisperReq.RequestURI)
	}
	whisperResp.Body.Close()

	return nil
}
