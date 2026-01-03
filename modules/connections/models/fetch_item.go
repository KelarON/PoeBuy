package models

import "time"

type FetchItem struct {
	Result []struct {
		ID      string `json:"id"`
		Listing struct {
			Method  string    `json:"method"`
			Indexed time.Time `json:"indexed"`
			Stash   struct {
				Name string `json:"name"`
				X    int    `json:"x"`
				Y    int    `json:"y"`
			} `json:"stash"`
			Whisper      string `json:"whisper"`
			WhisperToken string `json:"whisper_token"`
			HideoutToken string `json:"hideout_token"`
			Account      struct {
				Name   string `json:"name"`
				Online struct {
					League string `json:"league"`
				} `json:"online"`
				LastCharacterName string `json:"lastCharacterName"`
				Language          string `json:"language"`
				Realm             string `json:"realm"`
			} `json:"account"`
			Price struct {
				Type     string  `json:"type"`
				Amount   float64 `json:"amount"`
				Currency string  `json:"currency"`
			} `json:"price"`
		} `json:"listing"`
		Item struct {
			Verified   bool   `json:"verified"`
			W          int    `json:"w"`
			H          int    `json:"h"`
			Icon       string `json:"icon"`
			Support    bool   `json:"support"`
			League     string `json:"league"`
			ID         string `json:"id"`
			Name       string `json:"name"`
			TypeLine   string `json:"typeLine"`
			BaseType   string `json:"baseType"`
			Identified bool   `json:"identified"`
			Ilvl       int    `json:"ilvl"`
			Note       string `json:"note"`
			Corrupted  bool   `json:"corrupted"`
			Properties []struct {
				Name        string        `json:"name"`
				Values      []interface{} `json:"values"`
				DisplayMode int           `json:"displayMode"`
				Type        int           `json:"type,omitempty"`
			} `json:"properties"`
			Requirements []struct {
				Name        string          `json:"name"`
				Values      [][]interface{} `json:"values"`
				DisplayMode int             `json:"displayMode"`
				Type        int             `json:"type"`
			} `json:"requirements"`
			SecDescrText string   `json:"secDescrText"`
			ExplicitMods []string `json:"explicitMods"`
			DescrText    string   `json:"descrText"`
			FrameType    int      `json:"frameType"`
			Extended     struct {
				Text string `json:"text"`
			} `json:"extended"`
		} `json:"item"`
	} `json:"result"`
}
