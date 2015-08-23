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
	"cruncher/app/models/crunch"
	"cruncher/app/models/dataFormat"
	"cruncher/app/models/database"
	"cruncher/app/models/riotapi"
	"errors"
	"github.com/revel/revel"
	// "strings"
	"time"
)

var (
	ErrDatabaseError        = errors.New("query: database error")
	ErrDatabaseDisconnected = database.ErrDisconnected
	ErrAPIError             = errors.New("query: riot api error")
	ErrNotFound             = riotapi.ErrNotFound
)

func GetStats(name string, region string, isNew bool) (dataFormat.PlayerData,
	bool, error) {
	player, err := database.GetSummonerData(name, region)
	if err == nil {
		return player, isNew, nil
	}

	if err == ErrDatabaseDisconnected {
		return dataFormat.PlayerData{}, isNew, ErrDatabaseDisconnected
	}

	if err != database.ErrNoResults {
		revel.ERROR.Println("query: error getting stats for "+name+
			" on "+region+":", err)
		return dataFormat.PlayerData{}, isNew, ErrDatabaseError
	}

	if isNew {
		errorMessage := "query: missing player in database for " + name +
			" on " + region
		revel.ERROR.Println(errorMessage)
		return dataFormat.PlayerData{}, isNew, errors.New(errorMessage)
	}

	id, resolvedName, err := riotapi.ResolveSummonerId(name, region)
	if err == riotapi.ErrNotFound {
		return dataFormat.PlayerData{}, isNew, ErrNotFound
	} else if err != nil {
		return dataFormat.PlayerData{}, isNew, ErrAPIError
	}

	newPlayer := dataFormat.Player{
		Region:         region,
		Tier:           "UNKNOWN",
		SummonerId:     id,
		SummonerName:   resolvedName,
		NormalizedName: dataFormat.NormalizeName(resolvedName),
		NextUpdate:     time.Now(),
		NextLongUpdate: time.Now(),
	}

	newPlayer.InternalId, err = database.CreatePlayer(newPlayer)
	if err == ErrDatabaseDisconnected {
		return dataFormat.PlayerData{}, isNew, err
	} else if err != nil {
		revel.ERROR.Println("query: error creating player in database:", err)
		return dataFormat.PlayerData{}, isNew, ErrDatabaseError
	}

	err = UpdatePlayer(newPlayer)
	if err != nil {
		revel.ERROR.Println("query: error updating player in database:", err)
		_ = database.DeletePlayer(newPlayer)
		return dataFormat.PlayerData{}, isNew, err
	}

	return GetStats(name, region, true)
}

func UpdatePlayer(player dataFormat.Player) error {
	games, err := riotapi.GetRecentGames(player.SummonerId, player.Region)
	if err != nil {
		return err
	}

	earliestDate := time.Now()
	for _, game := range games {
		if game.Date.Before(earliestDate) {
			earliestDate = game.Date
		}
	}

	tier, err := riotapi.GetTier(player.SummonerId, player.Region)
	if err != nil {
		return err
	}

	newData := struct {
		Tier           string    `gorethink:"t"`
		NextLongUpdate time.Time `gorethink:"nl"`
		RecordStart    time.Time `gorethink:"rs"`
	}{tier, time.Now().Add(time.Hour * 168), earliestDate}

	if err = database.UpdatePlayerInformation(player, newData); err != nil {
		return err
	}

	crunch.Crunch(player, games)

	return nil
}
