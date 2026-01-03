package connections

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"poebuy/modules/connections/models"
)

type Whisper struct {
	client *http.Client
	header http.Header
}

func NewWhisper(client *http.Client, header http.Header) *Whisper {
	return &Whisper{
		client: client,
		header: header,
	}
}

func (w *Whisper) Whisper(token string) error {
	jsonBody := []byte(fmt.Sprintf("{\"token\": \"%v\"}", token))
	whisperReq, err := http.NewRequest("POST", "https://www.pathofexile.com/api/trade/whisper", bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("whisper request creation error: %v", err)

	}
	whisperReq.Header = w.header
	whisperResp, err := w.client.Do(whisperReq)
	if err != nil {
		return fmt.Errorf("whisper request error: %v", err)

	}
	defer whisperResp.Body.Close()

	if whisperResp.StatusCode != 200 {
		errorMsg := &models.WhisperErrorResponse{}
		r, _ := io.ReadAll(whisperResp.Body)
		json.Unmarshal(r, errorMsg)
		return fmt.Errorf("Whisper error: %v %v", whisperResp.Status, errorMsg.Error.Message)
	}

	return nil
}
