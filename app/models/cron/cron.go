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
	// "github.com/revel/revel/cache"
	"sync"
	"time"
)

var regions = []string{
	"na", "euw", "eune", "lan", "las", "oce", "br", "ru", "kr", "tr",
}

var updateRate int

var updateHealth = make(map[string]int)
var longUpdateHealth = make(map[string]int)
var updateLock = &sync.Mutex{}
var updateWg = &sync.WaitGroup{}

func UpdatePlayer(player dataFormat.Player) {
	defer func() {
		if r := recover(); r != nil {
			revel.ERROR.Println("UpdatePlayer: recovered from panic")
			revel.ERROR.Println(r)
		}
	}()

	defer updateWg.Done()

	if updateHealth[player.Region] <= 0 {
		return
	}

	games, err := riotapi.GetRecentGames(player.SummonerId, player.Region)
	if err != nil {
		revel.WARN.Println("cron: failed to get games for player:", player.InternalId)
		revel.WARN.Println(err)

		if err != riotapi.ErrNotFound {
			updateHealth[player.Region] -= 1
		}
		return
	}

	updateHealth[player.Region] += 1

	crunch.Crunch(player, games)
	// go cache.Delete(player.Region + ":" + player.NormalizedName)
}

func LongUpdatePlayer(player dataFormat.Player) {
	defer func() {
		if r := recover(); r != nil {
			revel.ERROR.Println("LongUpdatePlayer: recovered from panic")
			revel.ERROR.Println(r)
		}
	}()

	defer updateWg.Done()

	if longUpdateHealth[player.Region] <= 0 {
		return
	}

	tier, err := riotapi.GetTier(player.SummonerId, player.Region)
	if err != nil {
		revel.WARN.Println("cron: failed to get tier for player:", player.InternalId)
		revel.WARN.Println(err)

		if err != riotapi.ErrNotFound {
			longUpdateHealth[player.Region] -= 1
		}
		return
	}

	longUpdateHealth[player.Region] += 1

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
		revel.INFO.Println("cron: starting player updates")

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

		if len(players) != 0 {
			for _, region := range regions {
				updateHealth[region] = 5
			}

			for _, player := range players {
				updateWg.Add(1)
				go UpdatePlayer(player)
				time.Sleep(time.Millisecond * time.Duration(updateRate))
			}

			updateWg.Wait()

			for _, region := range regions {
				if updateHealth[region] <= 0 {
					revel.WARN.Println("cron: player update for region " + region +
						" stopped due to no health")
				}
			}
		}

		revel.INFO.Println("cron: finished player updates")
		updateLock.Unlock()

		time.Sleep(time.Minute * 10)
	}
}

func LongMonitor() {
	for {
		updateLock.Lock()
		revel.INFO.Println("cron: starting long updates")

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

		if len(players) != 0 {
			for _, region := range regions {
				longUpdateHealth[region] = 5
			}

			for _, player := range players {
				updateWg.Add(1)
				go LongUpdatePlayer(player)
				time.Sleep(time.Millisecond * time.Duration(updateRate))
			}

			updateWg.Wait()

			for _, region := range regions {
				if longUpdateHealth[region] <= 0 {
					revel.WARN.Println("cron: player update for region " + region +
						" stopped due to no health")
				}
			}
		}

		revel.INFO.Println("cron: finished long updates")
		updateLock.Unlock()

		time.Sleep(time.Hour)
	}
}

func Start() {
	updateRate = revel.Config.IntDefault("updaterate", 2000)
	revel.INFO.Printf("Update rate set to: %v ms", updateRate)
	go RecordMonitor()
	go LongMonitor()
	revel.INFO.Println("Now running monitors")
}
