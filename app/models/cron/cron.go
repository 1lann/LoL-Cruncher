
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


var ProcessingRunning bool = false

func processPlayers(players []dataFormat.BasicPlayer) {
	for _, player := range players {
		time.Sleep(time.Duration(2) * time.Second)

		revel.INFO.Printf("Processing player: %v", player.Id)

		playerGames, err := riotapi.GetRecentGames(player.Id, player.Region)
		if err != nil {
			revel.ERROR.Println("Failed to load recent games!")
			revel.ERROR.Println(err)
			continue
		}

		playerData, resp := database.GetSummonerData(player.Id,
			player.Region)

		if resp != database.Yes {
			revel.ERROR.Println("Non-ok response from " +
				"database while processing!")
			continue
		}

		newPlayer := crunch.Crunch(playerData, playerGames)
		newPlayer.NextUpdate = crunch.GetNextUpdate(playerGames)

		revel.INFO.Printf("Next update time in hours: %v",
			time.Since(newPlayer.NextUpdate).Hours())

		resp = database.StoreSummonerData(newPlayer)
		if resp != database.Yes {
			revel.ERROR.Println("Non-ok response for storing player data")
		}
	}
}

func processTiers(players []dataFormat.BasicPlayer) {
	for _, player := range players {
		time.Sleep(time.Duration(2) * time.Second)

		revel.INFO.Printf("Processing tier for: %v", player.Id)

		tier, err := riotapi.GetTier(player.Id, player.Region)
		if err != nil {
			revel.ERROR.Println("Failed to load player tier!")
			revel.ERROR.Println(err)
			continue
		}

		nextLongUpdate := time.Now().Add(time.Duration(72) * time.Hour)

		resp := database.StoreTier(player.Id, player.Region, tier,
			nextLongUpdate)
		if resp != database.Yes {
			revel.ERROR.Println("Non-ok response from " +
				"database while storing tier!")
			continue
		}
	}
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
			time.Sleep(time.Duration(60) * time.Minute)
			continue
		}

		for ProcessingRunning {
			time.Sleep(time.Duration(10) * time.Second)
		}

		ProcessingRunning = true
		processPlayers(players)
		ProcessingRunning = false

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
			time.Sleep(time.Hour)
			continue
		}

		for ProcessingRunning {
			time.Sleep(time.Duration(10) * time.Second)
		}

		ProcessingRunning = true
		processTiers(players)
		ProcessingRunning = false

		time.Sleep(time.Duration(24) * time.Hour)
	}
}

func Start() {
	go RecordMonitor()
	go LongTermMonitor()
	revel.INFO.Println("Now running monitors")
}
