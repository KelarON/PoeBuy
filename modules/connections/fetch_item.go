package connections

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"poebuy/modules/connections/models"
)

type Fetcher struct {
	Client *http.Client
	Header http.Header
}

func NewFetcher(client *http.Client, header http.Header) *Fetcher {
	return &Fetcher{
		Client: client,
		Header: header,
	}
}

func (f *Fetcher) FetchItems(items []string, code string) (*[]models.FetchItem, error) {

	entities := []models.FetchItem{}

	for _, item := range items {
		reqBody := bytes.NewBuffer([]byte{})
		ItemReq, err := http.NewRequest("GET", fmt.Sprintf("https://www.pathofexile.com/api/trade/fetch/%v?query=%v", item, code), reqBody)
		if err != nil {
			return nil, fmt.Errorf("item request creation error: %v", err)

		}
		ItemReq.Header = f.Header

		ItemResp, err := f.Client.Do(ItemReq)
		if err != nil {
			return nil, fmt.Errorf("reqvest sending error: %v", err)

		}
		defer ItemResp.Body.Close()

		if ItemResp.StatusCode != http.StatusOK {
			continue
		}

		gz, err := gzip.NewReader(ItemResp.Body)
		if err != nil {
			return nil, fmt.Errorf("request decoding error: %v", err)

		}

		bodyBytes, err := io.ReadAll(gz)
		if err != nil {
			return nil, fmt.Errorf("request reading error: %v", err)

		}

		itemInfo := models.FetchItem{}
		err = json.Unmarshal(bodyBytes, &itemInfo)
		if err != nil {
			return nil, fmt.Errorf("item unmarshal error: %v", err)
		}
		entities = append(entities, itemInfo)
	}

	return &entities, nil
}
