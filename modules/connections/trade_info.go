package connections

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"poebuy/modules/connections/headers"
	"poebuy/modules/connections/models"
	"poebuy/utils"
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
		return nil, fmt.Errorf("error making tradeinfo request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, ErrorBadPoessid
	}

	bt, err := utils.ReadEncodedResponse(res)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	info := &models.TradeInfo{}

	nicknameMatches := regexp.MustCompile(`account/view-profile/.*?">(.+?)<`).FindSubmatch(bt)
	if len(nicknameMatches) == 0 {
		return nil, ErrorBadPoessid
	}

	info.Nickname = string(nicknameMatches[1])

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

	bt, err = utils.ReadEncodedResponse(res)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
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
