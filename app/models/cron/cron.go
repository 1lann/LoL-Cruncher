
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

package cron

import (
	"time"
	"cruncher/app/models/database"
	"cruncher/app/models/riotapi"
	"cruncher/app/models/crunch"
	"cruncher/app/models/dataFormat"
	"github.com/revel/revel"
)

var monitorRunning bool = false
var longMonitorRunning bool = false
var updateRate = 2000 // 2000 by default, also replicated at the bottom of
					  // the file in Start()

var ProcessingRunning bool = false

func processPlayers(players []dataFormat.BasicPlayer) {
	ProcessingRunning = true
	displayErrors := true

	for _, itPlayer := range players {
		time.Sleep(time.Duration(updateRate) * time.Millisecond)
		go func(player dataFormat.BasicPlayer) {
			revel.INFO.Printf("Processing player: %v", player.Id)

			playerGames, err := riotapi.GetRecentGames(player.Id, player.Region)
			if err != nil {
				if displayErrors {
					revel.ERROR.Println("Failed to load recent games!")
					revel.ERROR.Println(err)
					displayErrors = false
				}
			}

			playerData, resp := database.GetSummonerData(player.Id,
				player.Region)

			if resp != database.Yes {
				if displayErrors {
					revel.ERROR.Println("Non-ok response from " +
						"database while processing!")
					displayErrors = false
				}
			}

			newPlayer := crunch.Crunch(playerData, playerGames)
			newPlayer.NextUpdate = crunch.GetNextUpdate(playerGames)

			revel.INFO.Printf("Next update time in hours: %v",
				time.Since(newPlayer.NextUpdate).Hours())

			resp = database.StoreSummonerData(newPlayer)
			if (resp != database.Yes) && displayErrors {
				revel.ERROR.Println("Non-ok response for storing player data")
				displayErrors = false
			}
		}(itPlayer)
	}

	ProcessingRunning = false
}

func processTiers(players []dataFormat.BasicPlayer) {
	ProcessingRunning = true
	displayErrors := true

	for _, itPlayer := range players {
		time.Sleep(time.Duration(updateRate) * time.Millisecond)
		go func(player dataFormat.BasicPlayer) {
			revel.INFO.Printf("Processing tier for: %v", player.Id)

			tier, err := riotapi.GetTier(player.Id, player.Region)
			if err != nil {
				if displayErrors {
					revel.ERROR.Println("Failed to load player tier!")
					revel.ERROR.Println(err)
					displayErrors = false
				}
			}

			nextLongUpdate := time.Now().Add(time.Duration(72) * time.Hour)

			resp := database.StoreTier(player.Id, player.Region, tier,
				nextLongUpdate)
			if (resp != database.Yes) && displayErrors {
				revel.ERROR.Println("Non-ok response from " +
					"database while storing tier!")
				displayErrors = false
			}
		}(itPlayer)
	}

	ProcessingRunning = false
}

func RecordMonitor() {
	if monitorRunning {
		revel.ERROR.Println("Monitor already running!")
		return
	}
	monitorRunning = true
	for true {
		players, resp := database.GetUpdatePlayers()

		if resp != database.Yes {
			revel.ERROR.Println("RecordMonitor non-ok response from database!")
			time.Sleep(time.Duration(10) * time.Minute)
			continue
		}

		for ProcessingRunning {
			time.Sleep(time.Duration(10) * time.Second)
		}

		go processPlayers(players)

		time.Sleep(time.Duration(60) * time.Minute)
	}
}

func LongTermMonitor() {
	if longMonitorRunning {
		revel.ERROR.Println("Long term mointor already running!")
		return
	}
	longMonitorRunning = true
	for true {
		players, resp := database.GetLongUpdatePlayers()

		if resp != database.Yes {
			revel.ERROR.Println("LongTermMonitor non-ok " +
				"response from database!")
			time.Sleep(time.Duration(10) * time.Minute)
			continue
		}

		for ProcessingRunning {
			time.Sleep(time.Duration(10) * time.Second)
		}

		go processTiers(players)

		time.Sleep(time.Duration(24) * time.Hour)
	}
}

func Start() {
	updateRate = revel.Config.IntDefault("updaterate", 2000)
	revel.INFO.Printf("Update rate set to: %v ms", updateRate)
	go RecordMonitor()
	go LongTermMonitor()
	revel.INFO.Println("Now running monitors")
}
