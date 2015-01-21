
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
	"time"
	"math"
)

func chomp(playerData *dataFormat.Player, game dataFormat.Game) {
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

	var parsedType string

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

	allQAllM := playerData.AllQueues.AllMonths.All
	allQAllMThisC := playerData.AllQueues.AllMonths.Champions[game.ChampionId]
	allQThisM := playerData.AllQueues.MonthlyStats[game.YearMonth].All
	allQThisMThisC := playerData.AllQueues.MonthlyStats[game.YearMonth].Champions[game.ChampionId]

	thisQAllM := playerData.QueueStats[parsedType].AllMonths.All
	thisQAllMThisC := playerData.QueueStats[parsedType].AllMonths.Champions[game.ChampionId]
	thisQThisM := playerData.QueueStats[parsedType].MonthlyStats[game.YearMonth].All
	thisQThisMThisC := playerData.QueueStats[parsedType].MonthlyStats[game.YearMonth].Champions[game.ChampionId]

	if game.DidWin {
		allQAllM.Wins++
		allQAllMThisC.Wins++
		allQThisM.Wins++
		allQThisMThisC.Wins++

		thisQAllM.Wins++
		thisQAllMThisC.Wins++
		thisQThisM.Wins++
		thisQThisMThisC.Wins++

		if game.IsOnBlue {
			allQAllM.Blue.Wins++
			allQThisM.Blue.Wins++

			thisQAllM.Blue.Wins++
			thisQThisM.Blue.Wins++
		} else {
			allQAllM.Red.Wins++
			allQThisM.Red.Wins++

			thisQAllM.Red.Wins++
			thisQThisM.Red.Wins++
		}
	} else {
		allQAllM.Losses++
		allQAllMThisC.Losses++
		allQThisM.Losses++
		allQThisMThisC.Losses++

		thisQAllM.Losses++
		thisQAllMThisC.Losses++
		thisQThisM.Losses++
		thisQThisMThisC.Losses++

		if game.IsOnBlue {
			allQAllM.Blue.Losses++
			allQThisM.Blue.Losses++

			thisQAllM.Blue.Losses++
			thisQThisM.Blue.Losses++
		} else {
			allQAllM.Red.Losses++
			allQThisM.Red.Losses++

			thisQAllM.Red.Losses++
			thisQThisM.Red.Losses++
		}
	}

	allQAllM.TimePlayed += game.Duration
	allQAllMThisC.TimePlayed += game.Duration
	allQThisM.TimePlayed += game.Duration
	allQThisMThisC.TimePlayed += game.Duration

	thisQAllM.TimePlayed += game.Duration
	thisQAllMThisC.TimePlayed += game.Duration
	thisQThisM.TimePlayed += game.Duration
	thisQThisMThisC.TimePlayed += game.Duration



	allQAllM.Kills += game.Kills
	allQAllMThisC.Kills += game.Kills
	allQThisM.Kills += game.Kills
	allQThisMThisC.Kills += game.Kills

	thisQAllM.Kills += game.Kills
	thisQAllMThisC.Kills += game.Kills
	thisQThisM.Kills += game.Kills
	thisQThisMThisC.Kills += game.Kills



	allQAllM.Assists += game.Assists
	allQAllMThisC.Assists += game.Assists
	allQThisM.Assists += game.Assists
	allQThisMThisC.Assists += game.Assists

	thisQAllM.Assists += game.Assists
	thisQAllMThisC.Assists += game.Assists
	thisQThisM.Assists += game.Assists
	thisQThisMThisC.Assists += game.Assists



	allQAllM.Deaths += game.Deaths
	allQAllMThisC.Deaths += game.Deaths
	allQThisM.Deaths += game.Deaths
	allQThisMThisC.Deaths += game.Deaths

	thisQAllM.Deaths += game.Deaths
	thisQAllMThisC.Deaths += game.Deaths
	thisQThisM.Deaths += game.Deaths
	thisQThisMThisC.Deaths += game.Deaths



	allQAllM.MinionsKilled += game.MinionsKilled
	allQAllMThisC.MinionsKilled += game.MinionsKilled
	allQThisM.MinionsKilled += game.MinionsKilled
	allQThisMThisC.MinionsKilled += game.MinionsKilled

	thisQAllM.MinionsKilled += game.MinionsKilled
	thisQAllMThisC.MinionsKilled += game.MinionsKilled
	thisQThisM.MinionsKilled += game.MinionsKilled
	thisQThisMThisC.MinionsKilled += game.MinionsKilled



	allQAllM.MonstersKilled += game.MonstersKilled
	allQAllMThisC.MonstersKilled += game.MonstersKilled
	allQThisM.MonstersKilled += game.MonstersKilled
	allQThisMThisC.MonstersKilled += game.MonstersKilled

	thisQAllM.MonstersKilled += game.MonstersKilled
	thisQAllMThisC.MonstersKilled += game.MonstersKilled
	thisQThisM.MonstersKilled += game.MonstersKilled
	thisQThisMThisC.MonstersKilled += game.MonstersKilled



	allQAllM.WardsPlaced += game.WardsPlaced
	allQAllMThisC.WardsPlaced += game.WardsPlaced
	allQThisM.WardsPlaced += game.WardsPlaced
	allQThisMThisC.WardsPlaced += game.WardsPlaced

	thisQAllM.WardsPlaced += game.WardsPlaced
	thisQAllMThisC.WardsPlaced += game.WardsPlaced
	thisQThisM.WardsPlaced += game.WardsPlaced
	thisQThisMThisC.WardsPlaced += game.WardsPlaced


	allQAllM.GoldEarned += game.GoldEarned
	allQAllMThisC.GoldEarned += game.GoldEarned
	allQThisM.GoldEarned += game.GoldEarned
	allQThisMThisC.GoldEarned += game.GoldEarned

	thisQAllM.GoldEarned += game.GoldEarned
	thisQAllMThisC.GoldEarned += game.GoldEarned
	thisQThisM.GoldEarned += game.GoldEarned
	thisQThisMThisC.GoldEarned += game.GoldEarned



	allQAllM.DoubleKills += game.DoubleKills
	allQThisM.DoubleKills += game.DoubleKills

	thisQAllM.DoubleKills += game.DoubleKills
	thisQThisM.DoubleKills += game.DoubleKills



	allQAllM.TripleKills += game.TripleKills
	allQThisM.TripleKills += game.TripleKills

	thisQAllM.TripleKills += game.TripleKills
	thisQThisM.TripleKills += game.TripleKills



	allQAllM.QuadraKills += game.QuadraKills
	allQThisM.QuadraKills += game.QuadraKills

	thisQAllM.QuadraKills += game.QuadraKills
	thisQThisM.QuadraKills += game.QuadraKills



	allQAllM.PentaKills += game.PentaKills
	allQThisM.PentaKills += game.PentaKills

	thisQAllM.PentaKills += game.PentaKills
	thisQThisM.PentaKills += game.PentaKills



	allQAllM.WardsKilled += game.WardsKilled
	allQThisM.WardsKilled += game.WardsKilled

	thisQAllM.WardsKilled += game.WardsKilled
	thisQThisM.WardsKilled += game.WardsKilled


	//
	//	allQAllM 1/8
	//
	playerData.AllQueues.AllMonths.All = allQAllM

	if playerData.AllQueues.AllMonths.Champions == nil {
		playerData.AllQueues.AllMonths.Champions =
			make(map[string]dataFormat.BasicNumberOf)
	}

	//
	//	allQAllMThisC 2/8
	//
	playerData.AllQueues.AllMonths.Champions[game.ChampionId] = allQAllMThisC

	if playerData.AllQueues.MonthlyStats == nil {
		playerData.AllQueues.MonthlyStats = make(map[string]dataFormat.Stats)
	}

	allQueuesThisMonth := playerData.AllQueues.MonthlyStats[game.YearMonth]

	//
	//	allQThisM 3/8
	//
	allQueuesThisMonth.All = allQThisM

	if allQueuesThisMonth.Champions == nil {
		allQueuesThisMonth.Champions = make(map[string]dataFormat.BasicNumberOf)
	}

	//
	//	allQThisMThisC 4/8
	//
	allQueuesThisMonth.Champions[game.ChampionId] = allQThisMThisC


	if playerData.QueueStats == nil {
		playerData.QueueStats = make(map[string]dataFormat.QueueStats)
	}

	thisQueueStats := playerData.QueueStats[parsedType]

	//
	//	thisQAllM 5/8
	//
	thisQueueStats.AllMonths.All = thisQAllM

	if thisQueueStats.AllMonths.Champions == nil {
		thisQueueStats.AllMonths.Champions =
			make(map[string]dataFormat.BasicNumberOf)
	}

	//
	// thisQAllMThisC 6/8
	//
	thisQueueStats.AllMonths.Champions[game.ChampionId] = thisQAllMThisC

	if thisQueueStats.MonthlyStats == nil {
		thisQueueStats.MonthlyStats = make(map[string]dataFormat.Stats)
	}

	thisQueueThisMonthStats := thisQueueStats.MonthlyStats[game.YearMonth]

	//
	//	thisQThisM 7/8
	//
	thisQueueThisMonthStats.All = thisQThisM

	if thisQueueThisMonthStats.Champions == nil {
		thisQueueThisMonthStats.Champions =
			make(map[string]dataFormat.BasicNumberOf)
	}

	//
	//	thisQThisMThisC 8/8
	//
	thisQueueThisMonthStats.Champions[game.ChampionId] = thisQThisMThisC

	//
	//	Put it back together
	//
	playerData.AllQueues.MonthlyStats[game.YearMonth] = allQueuesThisMonth
	thisQueueStats.MonthlyStats[game.YearMonth] = thisQueueThisMonthStats
	playerData.QueueStats[parsedType] = thisQueueStats

	return
}

func hasBeenProcessed(games []string, query string) bool {
	for _, game := range games {
		if game == query {
			return true
		}
	}
	return false
}

func GetNextUpdate(games []dataFormat.Game) time.Time {
	checkIndex := int(math.Min(float64(2), float64(len(games) - 1)))
	intervalDuration := time.Since(games[checkIndex].Date)

	if intervalDuration.Hours() > 24 {
		intervalDuration = time.Duration(24) * time.Hour
	}

	return time.Now().Add(intervalDuration)
}


func Crunch(playerData dataFormat.Player,
	games []dataFormat.Game) dataFormat.Player {
	var processedList []string
	for _, game := range games {
		processedList = append(processedList, game.Id)
		if !hasBeenProcessed(playerData.ProcessedGames, game.Id) {
			chomp(&playerData, game)
		}
	}

	playerData.ProcessedGames = processedList



	return playerData
}
