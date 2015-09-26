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
	"cruncher/app/models/crunch"
	"cruncher/app/models/dataFormat"
	"cruncher/app/models/database"
	"cruncher/app/models/riotapi"
	"github.com/revel/revel"
	"github.com/revel/revel/cache"
	"sync"
	"time"
)

var updateRate int

var updateHealth = 5
var longUpdateHealth = 5
var updateLock = &sync.Mutex{}

func UpdatePlayer(player dataFormat.Player) {
	if updateHealth <= 0 {
		return
	}

	games, err := riotapi.GetRecentGames(player.SummonerId, player.Region)
	if err != nil {
		revel.WARN.Println("cron: failed to get games for player:", player.InternalId)
		revel.WARN.Println(err)
		updateHealth -= 1
		return
	}

	updateHealth += 1

	crunch.Crunch(player, games)

	go cache.Delete(player.Region + ":" + player.NormalizedName)
}

func LongUpdatePlayer(player dataFormat.Player) {
	if longUpdateHealth <= 0 {
		return
	}

	tier, err := riotapi.GetTier(player.SummonerId, player.Region)
	if err != nil {
		revel.WARN.Println("cron: failed to get tier for player:", player.InternalId)
		revel.WARN.Println(err)
		longUpdateHealth -= 1
		return
	}

	longUpdateHealth += 1

	tierData := struct {
		Tier           string    `gorethink:"t"`
		NextLongUpdate time.Time `gorethink:"nl"`
	}{tier, time.Now().Add(time.Hour * 168)}

	if err = database.UpdatePlayerInformation(player, tierData); err != nil {
		revel.WARN.Println("cron: failed to update player tier:",
			player.InternalId)
		revel.WARN.Println(err)
	}
}

func RecordMonitor() {
	for {
		updateLock.Lock()
		players := []dataFormat.Player{}
		for {
			var err error
			players, err = database.GetUpdatePlayers()
			if err != nil {
				revel.WARN.Println("cron: failed to get player updates:", err)
				time.Sleep(time.Minute)
			} else {
				break
			}
		}

		updateHealth = 5

		for _, player := range players {
			go UpdatePlayer(player)
			time.Sleep(time.Millisecond * time.Duration(updateRate))
		}

		if updateHealth <= 0 {
			revel.WARN.Println("cron: player update stopped due to no health")
		}
		updateLock.Unlock()

		time.Sleep(time.Hour)
	}
}

func LongMonitor() {
	for {
		updateLock.Lock()
		players := []dataFormat.Player{}
		for {
			var err error
			players, err = database.GetLongUpdatePlayers()
			if err != nil {
				revel.WARN.Println("cron: failed to get player updates:", err)
				time.Sleep(time.Minute)
			} else {
				break
			}
		}

		longUpdateHealth = 5

		for _, player := range players {
			go LongUpdatePlayer(player)
			time.Sleep(time.Millisecond * time.Duration(updateRate))
		}

		if longUpdateHealth <= 0 {
			revel.WARN.Println("cron: player update stopped due to no health")
		}
		updateLock.Unlock()

		time.Sleep(time.Hour * 24)
	}
}

func Start() {
	updateRate = revel.Config.IntDefault("updaterate", 2000)
	revel.INFO.Printf("Update rate set to: %v ms", updateRate)
	go RecordMonitor()
	go LongMonitor()
	revel.INFO.Println("Now running monitors")
}
