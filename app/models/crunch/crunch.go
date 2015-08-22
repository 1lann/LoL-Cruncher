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

package crunch

import (
	"cruncher/app/models/dataFormat"
	"cruncher/app/models/database"
	"encoding/json"
	"github.com/revel/revel"
	"math"
	"sort"
	"time"
)

func hasBeenProcessed(games []string, query string) bool {
	for _, game := range games {
		if game == query {
			return true
		}
	}
	return false
}

type byDuration []time.Duration

func (a byDuration) Len() int           { return len(a) }
func (a byDuration) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byDuration) Less(i, j int) bool { return a[i] < a[j] }

func GetNextUpdate(games []dataFormat.Game) time.Time {
	// Get the average time diff for your last 5 games, multipled by 3
	if len(games) <= 0 {
		return time.Now().Add(time.Duration(24) * time.Hour)
	}

	sinceLastGame := time.Since(games[0].Date)
	if sinceLastGame > time.Duration(48)*time.Hour {
		return time.Now().Add(time.Duration(24) * time.Hour)
	}

	maxIndex := len(games) - 1

	lastGame := time.Since(games[maxIndex].Date)
	sortedDiffs := []time.Duration{}

	for i := maxIndex - 1; i >= 0; i-- {
		sortedDiffs = append(sortedDiffs, lastGame-time.Since(games[i].Date))
		lastGame = time.Since(games[i].Date)
	}

	sort.Sort(byDuration(sortedDiffs))

	maxCheck := int(math.Min(float64(4), float64(maxIndex)))

	total := time.Duration(0)
	for i := 0; i < maxCheck; i++ {
		total = total + sortedDiffs[i]
	}

	averageDiff := (total / time.Duration(maxCheck))
	intervalDuration := averageDiff * time.Duration(3)

	if intervalDuration.Hours() > 24 {
		intervalDuration = time.Duration(24) * time.Hour
	}

	return time.Now().Add(intervalDuration)
}

func displayError(data interface{}, dispErr error) {
	revel.ERROR.Println("crunch: data discrepancy warning!")
	revel.ERROR.Println("The following update query failed:")

	resp, err := json.Marshal(data)
	if err != nil {
		revel.ERROR.Println(data)
	} else {
		revel.ERROR.Println(string(resp))
	}

	revel.ERROR.Println(dispErr)
}

func chomp(player dataFormat.Player, game dataFormat.Game) {
	normalSR := (game.Type == "NORMAL")
	rankedSR := (game.Type == "RANKED_SOLO_5x5") ||
		(game.Type == "RANKED_PREMADE_5x5")
	teamBuilder := (game.Type == "CAP_5x5")
	teamSR := (game.Type == "RANKED_TEAM_5x5")
	normalTT := (game.Type == "NORMAL_3x3")
	rankedTT := (game.Type == "RANKED_PREMADE_3x3")
	teamTT := (game.Type == "RANKED_TEAM_3x3")
	aram := (game.Type == "ARAM_UNRANKED_5x5")
	if !(normalSR || rankedSR || teamSR || normalTT || rankedTT || teamTT ||
		teamBuilder || aram) {
		return
	}

	parsedType := ""

	if normalSR {
		parsedType = "Summoner's Rift normals"
	} else if rankedSR {
		parsedType = "Summoner's Rift ranked"
	} else if teamSR {
		parsedType = "Summoner's Rift ranked team"
	} else if normalTT {
		parsedType = "Twisted Treeline normals"
	} else if rankedTT {
		parsedType = "Twisted Treeline ranked"
	} else if teamTT {
		parsedType = "Twisted Treeline ranked team"
	} else if teamBuilder {
		parsedType = "team builder"
	} else if aram {
		parsedType = "all random all mid (ARAM)"
	}

	detailed := dataFormat.DetailedNumberOf{
		InternalPlayerId: player.InternalId,
		TimePeriod:       "all",
		Queue:            "all",
	}

	basic := dataFormat.BasicNumberOf{
		InternalPlayerId: player.InternalId,
		Champion:         game.ChampionId,
		TimePeriod:       "all",
		Queue:            "all",
	}

	if game.DidWin {
		detailed.Wins = 1
		basic.Wins = 1
	} else {
		detailed.Losses = 1
		basic.Losses = 1
	}

	detailed.TimePlayed = game.Duration
	basic.TimePlayed = game.Duration
	detailed.Kills = game.Kills
	basic.Kills = game.Kills
	detailed.Assists = game.Assists
	basic.Assists = game.Assists
	detailed.Deaths = game.Deaths
	basic.Deaths = game.Assists
	detailed.DoubleKills = game.DoubleKills
	detailed.TripleKills = game.TripleKills
	detailed.PentaKills = game.PentaKills
	detailed.GoldEarned = game.GoldEarned
	basic.GoldEarned = game.GoldEarned
	detailed.MinionsKilled = game.MinionsKilled
	basic.MinionsKilled = game.MinionsKilled
	detailed.MonstersKilled = game.MonstersKilled
	basic.MonstersKilled = game.MonstersKilled
	detailed.WardsPlaced = game.WardsPlaced
	basic.WardsPlaced = game.WardsPlaced
	detailed.WardsKilled = game.WardsKilled

	if game.IsOnBlue {
		detailed.Blue.Wins = detailed.Wins
		detailed.Blue.Losses = detailed.Losses
	} else {
		detailed.Red.Wins = detailed.Wins
		detailed.Red.Losses = detailed.Losses
	}

	if err := database.AddToDetailedPlayer(detailed); err != nil {
		displayError(detailed, err)
		return
	}

	detailed.TimePeriod = game.YearMonth
	if err := database.AddToDetailedPlayer(detailed); err != nil {
		displayError(detailed, err)
		return
	}

	detailed.Queue = parsedType
	if err := database.AddToDetailedPlayer(detailed); err != nil {
		displayError(detailed, err)
		return
	}

	detailed.TimePeriod = "all"
	if err := database.AddToDetailedPlayer(detailed); err != nil {
		displayError(detailed, err)
		return
	}

	if err := database.AddToBasicPlayer(basic); err != nil {
		displayError(basic, err)
		return
	}

	basic.TimePeriod = game.YearMonth
	if err := database.AddToBasicPlayer(basic); err != nil {
		displayError(basic, err)
		return
	}

	basic.Queue = parsedType
	if err := database.AddToBasicPlayer(basic); err != nil {
		displayError(basic, err)
		return
	}

	basic.TimePeriod = "all"
	if err := database.AddToBasicPlayer(basic); err != nil {
		displayError(basic, err)
		return
	}
}

func Crunch(player dataFormat.Player, games []dataFormat.Game) {
	var processedList []string
	for _, game := range games {
		processedList = append(processedList, game.Id)
		if !hasBeenProcessed(player.ProcessedGames, game.Id) {
			chomp(player, game)
		}
	}

	playerChanges := struct {
		ProcessedGames []string  `gorethink:"p"`
		NextUpdate     time.Time `gorethink:"nu"`
	}{processedList, GetNextUpdate(games)}

	if err := database.UpdatePlayerInformation(player,
		playerChanges); err != nil {
		revel.ERROR.Println("crunch: failed to update update information for player:",
			player.InternalId)
		revel.ERROR.Println(err)
	}
}
