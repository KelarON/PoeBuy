package models

type LivesearchNewItem struct {
	Result string `json:"result"`
	Count  int    `json:"count"`
}

type LivesearchAuthStatus struct {
	Auth bool `json:"auth"`
}
