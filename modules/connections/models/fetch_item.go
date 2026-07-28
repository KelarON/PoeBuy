package models

type FetchItem struct {
	Result []struct {
		ID      string `json:"id"`
		Listing struct {
			WhisperToken string `json:"whisper_token"`
			HideoutToken string `json:"hideout_token"`
		} `json:"listing"`
	} `json:"result"`
}
