// LoL Cruncher - A Historical League of Legends Statistics Tracker
// Copyright (C) 2015  Jason Chu (1lann) 1lanncontact@gmail.com

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package riotapi

import (
	"cruncher/app/models/dataFormat"
	"encoding/json"
	"errors"
	"github.com/revel/revel"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var apiKey string

var (
	ErrBadRequest   = errors.New("riotapi: bad request")
	ErrNotFound     = errors.New("riotapi: not found")
	ErrRateLimit    = errors.New("riotapi: rate limit")
	ErrServerError  = errors.New("riotapi: server error")
	ErrUnauthorized = errors.New("riotapi: unauthorized")
)

type Player struct {
	ChampionId int
	TeamId     int
}

type Stats struct {
	GoldEarned           int
	WardPlaced           int
	WardKilled           int
	Win                  bool
	DoubleKills          int
	TripleKills          int
	QuadraKills          int
	PentaKills           int
	UnrealKills          int
	NeutralMinionsKilled int
	TimePlayed           int
	ChampionsKilled      int
	Assists              int
	NumDeaths            int
	MinionsKilled        int
}

type Game struct {
	FellowPlayers []Player
	GameType      string
	GameId        int
	TeamId        int
	GameMode      string
	ChampionId    int
	CreateDate    int64
	SubType       string
	Stats         Stats
}

type gameHistory struct {
	Games []Game
}

func convertGame(game Game) dataFormat.Game {
	isMatchedGame := (game.GameType == "MATCHED_GAME")
	isClassicGame := (game.GameMode == "CLASSIC" || game.GameMode == "ARAM")

	createTime := time.Unix(int64(game.CreateDate/1000), 0)
	createYear := strconv.Itoa(createTime.Year())
	createMonth := strconv.Itoa(int(createTime.Month()))

	return dataFormat.Game{
		DidWin:         game.Stats.Win,
		IsOnBlue:       (game.TeamId == 100),
		IsNormals:      (isMatchedGame && isClassicGame),
		ChampionId:     strconv.Itoa(game.ChampionId),
		Duration:       game.Stats.TimePlayed,
		Id:             strconv.Itoa(game.GameId),
		Type:           game.SubType,
		Kills:          game.Stats.ChampionsKilled,
		Assists:        game.Stats.Assists,
		Deaths:         game.Stats.NumDeaths,
		DoubleKills:    game.Stats.DoubleKills,
		TripleKills:    game.Stats.TripleKills,
		QuadraKills:    game.Stats.QuadraKills,
		PentaKills:     game.Stats.PentaKills,
		GoldEarned:     game.Stats.GoldEarned,
		MinionsKilled:  game.Stats.MinionsKilled,
		MonstersKilled: game.Stats.NeutralMinionsKilled,
		WardsPlaced:    game.Stats.WardPlaced,
		WardsKilled:    game.Stats.WardKilled,
		YearMonth:      createYear + " " + createMonth,
		Date:           createTime,
	}
}

func constructRecentGamesURL(id string, region string) string {
	return "https://" + region + ".api.pvp.net/api/lol/" + region +
		"/v1.3/game/by-summoner/" + id + "/recent?api_key=" + apiKey
}

func constructSummonerNameURL(name string, region string) string {
	return "https://" + region + ".api.pvp.net/api/lol/" + region +
		"/v1.4/summoner/by-name/" + name + "?api_key=" + apiKey
}

func constructLeagueURL(id string, region string) string {
	return "https://" + region + ".api.pvp.net/api/lol/" + region +
		"/v2.5/league/by-summoner/" + id + "/entry?api_key=" + apiKey
}

var rateLimitingBlock time.Time

func rateLimitBlock(seconds int) {
	rateLimitingBlock = time.Now().Add(time.Second * time.Duration(seconds))
}

func isRateBlocked() bool {
	return time.Now().Before(rateLimitingBlock)
}

// Form a request to Riot's APIs
func requestRiotAPI(url string) ([]byte, error) {
	// Emtpy response used for error responses
	emptyResponse := []byte{}
	for i := 0; i < 5; i++ {
		// Check if we had been blocked before for making too many requests
		if isRateBlocked() {
			printOut := "Requested reject due to rate blocking: "
			revel.WARN.Println(printOut + url)
			return emptyResponse, ErrRateLimit
		}

		// TODO Comment out in the future?
		revel.INFO.Println("Making request to Riot Games API")
		revel.INFO.Println(url)

		resp, err := http.Get(url)
		if err != nil {
			return emptyResponse, err
		}

		defer resp.Body.Close()

		error1xx := resp.StatusCode < 200
		error3xx := resp.StatusCode >= 300 && resp.StatusCode < 400
		error5xx := resp.StatusCode >= 500

		// Do status code processing
		if error1xx || error3xx || error5xx {
			printOut := resp.Status + " Response from: "
			revel.ERROR.Println(printOut + url)
			if resp.StatusCode == 503 {
				printOut = "Is error 503, retrying following recommendations"
				revel.WARN.Println(printOut)
				time.Sleep(time.Second)
				continue
			} else {
				return emptyResponse, ErrServerError
			}
		}

		// API Errors handled here
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			if resp.StatusCode == 400 {
				return emptyResponse, ErrBadRequest
			} else if resp.StatusCode == 401 {
				return emptyResponse, ErrUnauthorized
			} else if resp.StatusCode == 404 {
				return emptyResponse, ErrNotFound
			} else if resp.StatusCode == 429 {
				printOut := "Hit maximum query limit from: "
				revel.WARN.Println(printOut + url)
				retryHeader := resp.Header.Get("Retry-After")
				retryAfter, err := strconv.Atoi(retryHeader)
				if err != nil {
					revel.ERROR.Println("Cannot determine retry after time!")
					rateLimitBlock(5)
				} else {
					rateLimitBlock(retryAfter + 1)
				}

				return emptyResponse, ErrRateLimit
			} else {
				printOut := resp.Status + " Response from: "
				revel.ERROR.Println(printOut + url)
				return emptyResponse, ErrServerError
			}
		}

		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			printOut := "Failed to read content from: "
			revel.ERROR.Println(printOut + url)
			return emptyResponse, err
		}

		return contents, nil
	}
	return emptyResponse, ErrServerError
}

// returns id, name, error
func ResolveSummonerId(name string, region string) (string, string, error) {
	summonerNameURL := constructSummonerNameURL(name, region)
	contents, err := requestRiotAPI(summonerNameURL)

	if err != nil {
		return "", "", err
	}

	// Parse JSON data from server
	var data map[string]interface{}

	err = json.Unmarshal(contents, &data)
	if err != nil {
		printOut := "Failed to unmarshal content from: "
		revel.ERROR.Println(printOut + summonerNameURL)
		return "", "", err
	}

	rawUsername := ""
	for catchUsername, _ := range data {
		rawUsername = catchUsername
		break
	}

	if rawUsername == "" {
		revel.ERROR.Println("Failed to parse raw username")
		revel.ERROR.Println(data)
		return "", "", errors.New("Data parse error")
	}

	userData, ok := data[rawUsername].(map[string]interface{})
	if !ok {
		revel.ERROR.Println("Failed to parse data past raw username")
		revel.ERROR.Println(data)
		return "", "", errors.New("Data parse error")
	}

	rawUserId, ok := userData["id"].(float64)
	if !ok {
		revel.ERROR.Println("Failed to parse id data")
		revel.ERROR.Println(data)
		return "", "", errors.New("Data parse error")
	}

	userId := strconv.FormatFloat(rawUserId, 'f', 0, 64)

	username, ok := userData["name"].(string)
	if !ok {
		revel.ERROR.Println("Failed to parse resolved username")
		revel.ERROR.Println(data)
		return "", "", errors.New("Data parse error")
	}

	// FUCK YES! THIS ONE SHITTY REQUEST HAS PASSED THROUGH 15 CHECKS!
	return userId, username, nil
}

// Returns tier as BRONZE, SILVER, etc. and error
func GetTier(id string, region string) (string, error) {
	leagueURL := constructLeagueURL(id, region)
	contents, err := requestRiotAPI(leagueURL)
	if err != nil {
		if err == ErrNotFound {
			return "UNRANKED", nil
		} else {
			return "", err
		}
	}

	var data map[string]interface{}
	err = json.Unmarshal(contents, &data)
	if err != nil {
		printOut := "Failed to unmarshal content from: "
		revel.ERROR.Println(printOut + leagueURL)
		return "", err
	}

	indexId := ""
	for catchId, _ := range data {
		indexId = catchId
		break
	}

	if indexId == "" {
		revel.ERROR.Println("Failed to parse index id")
		revel.ERROR.Println(data)
		return "", errors.New("Data parse error")
	}

	leagueDataArr, ok := data[indexId].([]interface{})
	if !ok {
		revel.ERROR.Println("Failed to parse league data array")
		revel.ERROR.Println(data)
		return "", errors.New("Data parse error")
	}

	leagueData, ok := leagueDataArr[0].(map[string]interface{})
	if !ok {
		revel.ERROR.Println("Failed to parse league data")
		revel.ERROR.Println(data)
		return "", errors.New("Data parse error")
	}

	tier, ok := leagueData["tier"].(string)
	if !ok {
		revel.ERROR.Println("Failed to parse tier data")
		revel.ERROR.Println(data)
		return "", errors.New("Data parse error")
	}

	return tier, nil
}

func GetRecentGames(id string, region string) ([]dataFormat.Game, error) {
	recentGamesURL := constructRecentGamesURL(id, region)
	contents, err := requestRiotAPI(recentGamesURL)

	if err != nil {
		return []dataFormat.Game{}, err
	}

	var gameData gameHistory
	json.Unmarshal(contents, &gameData)

	var results []dataFormat.Game

	for _, game := range gameData.Games {
		results = append(results, convertGame(game))
	}

	return results, nil
}

func LoadAPIKey() {
	var found bool
	apiKey, found = revel.Config.String("riotapikey")
	if !found {
		revel.ERROR.Println("No riotapikey field found in conf/app.conf")
		panic(errors.New("No riotapikey field found in conf/app.conf"))
	}
}
