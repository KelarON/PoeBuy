package models

type TradeInfo struct {
	Nickname string
	Leagues  []League
}

type League struct {
	ID    string `json:"id"`
	Realm string `json:"realm"`
	Text  string `json:"text"`
}

func (t *TradeInfo) GetLeagues() []string {

	leagues := make([]string, 0, 10)

	for _, leg := range t.Leagues {
		leagues = append(leagues, leg.ID)
	}

	return leagues
}
