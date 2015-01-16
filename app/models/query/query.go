
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

package query

import (
	"github.com/revel/revel"
	"cruncher/app/models/dataFormat"
	"cruncher/app/models/riotapi"
	"cruncher/app/models/database"
	"cruncher/app/models/crunch"
	"time"
	"errors"
)

// Returns id, resolved summoner name, error message
func ResolveSummonerId(name string, region string) (string, string, error) {
	playerId, resolvedName, resp := database.GetSummonerID(name, region)
	if resp == database.Yes {
		return playerId, resolvedName, nil
	} else if resp == database.Down {
		return "", "", errors.New("database down")
	} else if resp == database.No {
		// Get id from Riot Games API
		playerId, resolvedName, err := riotapi.ResolveSummonerId(name, region)
		if err != nil {
			return "", "", err
		}

		// Store result on database
		database.StoreSummonerID(resolvedName, playerId, region)

		return playerId, resolvedName, nil
	} else {
		return "", "", errors.New("database error")
	}
}

// Returns id, name, player data, errors
func RegisterSummoner(name string, region string) (string, string,
		dataFormat.Player, error) {
	revel.INFO.Println("Registering new summoner")
	id, resolvedName, err := ResolveSummonerId(name, region)
	if err != nil {
		return "", "", dataFormat.Player{}, err
	}

	tier, err := riotapi.GetTier(id, region)
	if err != nil {
		return "", "", dataFormat.Player{}, err
	}

	var newPlayer dataFormat.Player

	newPlayer.Tier = tier
	newPlayer.Id = id
	newPlayer.Region = region

	newPlayer.NextUpdate = time.Now().Add(time.Minute)
	newPlayer.NextLongUpdate = time.Now().Add(time.Hour)

	games, err := riotapi.GetRecentGames(id, region)
	if err != nil {
		return "", "", dataFormat.Player{}, err
	}

	// Get the first game and set RecordStart
	earliestDate := time.Now()
	for _, game := range games {
		if game.Date.Before(earliestDate) {
			earliestDate = game.Date
		}
	}

	newPlayer.RecordStart = earliestDate.Format("2 January 2006")

	newPlayer = crunch.Crunch(newPlayer, games)

	resp := database.StoreSummonerData(newPlayer)
	if resp == database.Yes {
		return id, resolvedName, newPlayer, nil
	} else if resp == database.Down {
		return "", "", dataFormat.Player{}, errors.New("database down")
	} else {
		return "", "", dataFormat.Player{}, errors.New("database error")
	}
}

// Returns resolved summoner name, player data, error message
func GetStats(name string, region string) (string, dataFormat.Player, error) {
	playerId, resolvedName, err := ResolveSummonerId(name, region)
	if err != nil {
		return "", dataFormat.Player{}, err
	}

	playerData, resp := database.GetSummonerData(playerId, region)
	if resp == database.Yes {
		return resolvedName, playerData, nil
	} else if resp == database.Down {
		return "", dataFormat.Player{}, errors.New("database down")
	} else if resp == database.No {
		// Go kill urself plz
		_, resolvedName, playerData, err := RegisterSummoner(name, region)
		return resolvedName, playerData, err
	} else {
		return "", dataFormat.Player{}, errors.New("database error")
	}
}
