package connections

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"poebuy/modules/connections/headers"
	"poebuy/modules/connections/models"
	"regexp"
	"strings"
)

var ErrorBadPoessid = errors.New("can't get trade info, check POESSID")

const (
	_PcLeagueId = "pc"
)

func GetTradeInfo(poesessid string) (*models.TradeInfo, error) {

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, "https://www.pathofexile.com", nil)
	if err != nil {
		return nil, err
	}
	req.Header = headers.GetFetchitemHeaders(poesessid)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, ErrorBadPoessid
	}

	gr, err := gzip.NewReader(res.Body)
	if err != nil {
		return nil, err
	}
	defer gr.Close()

	bt, err := io.ReadAll(gr)
	if err != nil {
		return nil, err
	}

	info := &models.TradeInfo{}

	info.Nickname = string(regexp.MustCompile(`account/view-profile/.*?">(.+?)<`).FindSubmatch(bt)[1])

	req, err = http.NewRequest(http.MethodGet, "https://api.pathofexile.com/league", nil)
	if err != nil {
		return nil, err
	}
	req.Header = headers.GetFetchitemHeaders(poesessid)

	res, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, ErrorBadPoessid
	}

	gr2, err := gzip.NewReader(res.Body)
	if err != nil {
		return nil, err
	}
	defer gr2.Close()

	bt, err = io.ReadAll(gr2)
	if err != nil {
		return nil, err
	}

	leagues := &models.PoeApiLeagueResponse{}

	err = json.Unmarshal(bt, leagues)
	if err != nil {
		return nil, err
	}

	for _, l := range leagues.Leagues {
		if l.Realm == _PcLeagueId && !strings.Contains(l.Description, "SSF") {
			info.Leagues = append(info.Leagues, models.League{ID: l.ID, Realm: l.Realm, Text: l.Name})
		}
	}

	return info, nil
}
