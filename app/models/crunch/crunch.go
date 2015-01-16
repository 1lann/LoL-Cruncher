
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
)

func chomp(playerData *dataFormat.Player, game dataFormat.Game) {
	normalSR := game.Type == "NORMAL"
	rankedSR := (game.Type == "RANKED_SOLO_5x5") ||
		(game.Type == "RANKED_PREMADE_5x5")
	teamSR := (game.Type == "RANKED_TEAM_5x5")
	normalTT := (game.Type == "NORMAL_3x3")
	rankedTT := (game.Type == "RANKED_PREMADE_3x3")
	teamTT := (game.Type == "RANKED_TEAM_3x3")
	if !(normalSR || rankedSR || teamSR || normalTT || rankedTT || teamTT) {
		return
	}

	var parsedType string

	if normalSR {
		parsedType = "Summoner's Rift Normals"
	} else if rankedSR {
		parsedType = "Summoner's Rift Ranked"
	} else if teamSR {
		parsedType = "Summoner's Rift Ranked Team"
	} else if normalTT {
		parsedType = "Twisted Treeline Normals"
	} else if rankedTT {
		parsedType = "Twisted Treeline Ranked"
	} else if teamTT {
		parsedType = "Twisted Treeline Ranked Team"
	}

	allAll := playerData.All.All
	allGameType := playerData.All.GameTypeStats[parsedType]
	allChampion := playerData.All.Champions[game.ChampionId]

	monthlyAll := playerData.MonthlyStats[game.YearMonth].All
	monthlyGameType := playerData.MonthlyStats[game.YearMonth].
		GameTypeStats[parsedType]
	monthlyAllChampion := playerData.MonthlyStats[game.YearMonth].
		Champions[game.ChampionId]


	if game.DidWin {
		allAll.Wins++
		allGameType.Wins++
		allChampion.Wins++

		monthlyAll.Wins++
		monthlyGameType.Wins++
		monthlyAllChampion.Wins++
		if game.IsOnBlue {
			allAll.Blue.Wins++
			allGameType.Blue.Wins++

			monthlyAll.Blue.Wins++
			monthlyGameType.Blue.Wins++
		} else {
			allAll.Red.Wins++
			allGameType.Red.Wins++

			monthlyAll.Red.Wins++
			monthlyGameType.Red.Wins++
		}
	} else {
		allAll.Losses++
		allGameType.Losses++
		allChampion.Losses++
		monthlyAll.Losses++
		monthlyGameType.Losses++
		monthlyAllChampion.Losses++
		if game.IsOnBlue {
			allAll.Blue.Losses++
			allGameType.Blue.Losses++

			monthlyAll.Blue.Losses++
			monthlyGameType.Blue.Losses++
		} else {
			allAll.Red.Losses++
			allGameType.Red.Losses++

			monthlyAll.Red.Losses++
			monthlyGameType.Red.Losses++
		}
	}

	allAll.TimePlayed += game.Duration
	allGameType.TimePlayed += game.Duration
	allChampion.TimePlayed += game.Duration

	monthlyAll.TimePlayed += game.Duration
	monthlyGameType.TimePlayed += game.Duration
	monthlyAllChampion.TimePlayed += game.Duration



	allAll.Kills += game.Kills
	allGameType.Kills += game.Kills
	allChampion.Kills += game.Kills

	monthlyAll.Kills += game.Kills
	monthlyGameType.Kills += game.Kills
	monthlyAllChampion.Kills += game.Kills



	allAll.Assists += game.Assists
	allGameType.Assists += game.Assists
	allChampion.Assists += game.Assists

	monthlyAll.Assists += game.Assists
	monthlyGameType.Assists += game.Assists
	monthlyAllChampion.Assists += game.Assists



	allAll.Deaths += game.Deaths
	allGameType.Deaths += game.Deaths
	allChampion.Deaths += game.Deaths

	monthlyAll.Deaths += game.Deaths
	monthlyGameType.Deaths += game.Deaths
	monthlyAllChampion.Deaths += game.Deaths



	allAll.MinionsKilled += game.MinionsKilled
	allGameType.MinionsKilled += game.MinionsKilled
	allChampion.MinionsKilled += game.MinionsKilled

	monthlyAll.MinionsKilled += game.MinionsKilled
	monthlyGameType.MinionsKilled += game.MinionsKilled
	monthlyAllChampion.MinionsKilled += game.MinionsKilled



	allAll.MonstersKilled += game.MonstersKilled
	allGameType.MonstersKilled += game.MonstersKilled
	allChampion.MonstersKilled += game.MonstersKilled

	monthlyAll.MonstersKilled += game.MonstersKilled
	monthlyGameType.MonstersKilled += game.MonstersKilled
	monthlyAllChampion.MonstersKilled += game.MonstersKilled



	allAll.WardsPlaced += game.WardsPlaced
	allGameType.WardsPlaced += game.WardsPlaced
	allChampion.WardsPlaced += game.WardsPlaced

	monthlyAll.WardsPlaced += game.WardsPlaced
	monthlyGameType.WardsPlaced += game.WardsPlaced
	monthlyAllChampion.WardsPlaced += game.WardsPlaced



	allAll.DoubleKills += game.DoubleKills
	allGameType.DoubleKills += game.DoubleKills

	monthlyAll.DoubleKills += game.DoubleKills
	monthlyGameType.DoubleKills += game.DoubleKills



	allAll.TripleKills += game.TripleKills
	allGameType.TripleKills += game.TripleKills

	monthlyAll.TripleKills += game.TripleKills
	monthlyGameType.TripleKills += game.TripleKills



	allAll.QuadraKills += game.QuadraKills
	allGameType.QuadraKills += game.QuadraKills

	monthlyAll.QuadraKills += game.QuadraKills
	monthlyGameType.QuadraKills += game.QuadraKills



	allAll.PentaKills += game.PentaKills
	allGameType.PentaKills += game.PentaKills

	monthlyAll.PentaKills += game.PentaKills
	monthlyGameType.PentaKills += game.PentaKills



	allAll.GoldEarned += game.GoldEarned
	allGameType.GoldEarned += game.GoldEarned

	monthlyAll.GoldEarned += game.GoldEarned
	monthlyGameType.GoldEarned += game.GoldEarned


	allAll.WardsKilled += game.WardsKilled
	allGameType.WardsKilled += game.WardsKilled

	monthlyAll.WardsKilled += game.WardsKilled
	monthlyGameType.WardsKilled += game.WardsKilled

	// allAll := playerData.All.All
	// allGameType := playerData.All.GameTypeStats[parsedType]
	// allChampion := playerData.All.Champions[game.ChampionId]

	// monthlyAll := playerData.MonthlyStats[game.YearMonth].All
	// monthlyGameType := playerData.MonthlyStats[game.YearMonth].
	// 	GameTypeStats[parsedType]
	// monthlyAllChampion := playerData.MonthlyStats[game.YearMonth].
	// 	Champions[game.ChampionId]

	playerData.All.All = allAll
	if playerData.All.GameTypeStats == nil {
		playerData.All.GameTypeStats = make(map[string]dataFormat.DetailedNumberOf)
	}
	playerData.All.GameTypeStats[parsedType] = allGameType
	if playerData.All.Champions == nil {
		playerData.All.Champions = make(map[string]dataFormat.BasicNumberOf)
	}
	playerData.All.Champions[game.ChampionId] = allChampion

	if playerData.MonthlyStats == nil {
		playerData.MonthlyStats = make(map[string]dataFormat.Stats)
	}

	monthlyCopy := playerData.MonthlyStats[game.YearMonth]
	monthlyCopy.All = monthlyAll

	if monthlyCopy.GameTypeStats == nil {
		monthlyCopy.GameTypeStats = make(map[string]dataFormat.DetailedNumberOf)
	}
	monthlyCopy.GameTypeStats[parsedType] = monthlyGameType

	if monthlyCopy.Champions == nil {
		monthlyCopy.Champions = make(map[string]dataFormat.BasicNumberOf)
	}
	monthlyCopy.Champions[game.ChampionId] = monthlyAllChampion

	playerData.MonthlyStats[game.YearMonth] = monthlyCopy

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
