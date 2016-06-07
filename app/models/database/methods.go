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

package database

import (
	"cruncher/app/models/dataFormat"
	"errors"
	// "fmt"
	r "github.com/dancannon/gorethink"
	"github.com/revel/revel"
	"strings"
	"time"
)

var LastPlayerUpdate = time.Now()

var (
	ErrDisconnected      = errors.New("database: disconnected")
	ErrNoResults         = errors.New("database: no results")
	ErrInsertDiscrepancy = errors.New("database: insert discrepancy")
)

func GetBrowserPlayers() ([]dataFormat.Player, error) {
	if !IsConnected {
		go Connect()
		return []dataFormat.Player{}, ErrDisconnected
	}

	c, err := r.Table("players").OrderBy(r.OrderByOpts{"nn"}).
		Pluck("sn", "r").Run(activeSession)
	if isDisconnected(err) {
		return []dataFormat.Player{}, ErrDisconnected
	} else if err != nil {
		return []dataFormat.Player{}, err
	}

	results := []dataFormat.Player{}
	err = c.All(&results)
	c.Close()
	if isDisconnected(err) {
		return []dataFormat.Player{}, ErrDisconnected
	} else if err == r.ErrEmptyResult {
		return []dataFormat.Player{}, ErrNoResults
	} else if err != nil {
		return []dataFormat.Player{}, err
	}

	return results, nil
}

func GetSummonerData(name string, region string) (dataFormat.PlayerData,
	error) {
	if !IsConnected {
		go Connect()
		return dataFormat.PlayerData{}, ErrDisconnected
	}

	name = dataFormat.NormalizeName(name)
	c, err := r.Table("players").
		GetAllByIndex("nn", name).
		Filter(map[string]string{"r": region}).AtIndex(0).
		Merge(func(row r.Term) interface{} {
			return map[string]interface{}{
				"detailed": r.DB("cruncher").Table("detailed").
					Between([]interface{}{
						row.Field("id"),
						r.MinVal,
						r.MinVal,
					}, []interface{}{
						row.Field("id"),
						r.MaxVal,
						r.MaxVal,
					}, r.BetweenOpts{Index: "ippq"}).Without("id", "ip").
					CoerceTo("array"),
				"basic": r.DB("cruncher").Table("basic").
					Between([]interface{}{
						row.Field("id"),
						r.MinVal,
						r.MinVal,
						r.MinVal,
					}, []interface{}{
						row.Field("id"),
						r.MaxVal,
						r.MaxVal,
						r.MaxVal,
					}, r.BetweenOpts{Index: "ippqc"}).Without("id", "ip").
					CoerceTo("array"),
			}
		}).Run(activeSession)

	if isDisconnected(err) {
		return dataFormat.PlayerData{}, ErrDisconnected
	} else if err != nil {
		if strings.Contains(err.Error(), "gorethink: Index out of bounds: 0 in:") {
			return dataFormat.PlayerData{}, ErrNoResults
		}
		return dataFormat.PlayerData{}, err
	}

	playerData := dataFormat.PlayerData{}

	err = c.One(&playerData)
	c.Close()
	if isDisconnected(err) {
		return dataFormat.PlayerData{}, ErrDisconnected
	} else if err == r.ErrEmptyResult {
		return dataFormat.PlayerData{}, ErrNoResults
	} else if err != nil {
		return dataFormat.PlayerData{}, err
	}

	return playerData, nil
}

func AddToDetailedPlayer(details dataFormat.DetailedNumberOf) error {
	if !IsConnected {
		go Connect()
		return ErrDisconnected
	}

	resp, err := r.Table("detailed").GetAllByIndex("ippq",
		[]interface{}{
			details.InternalPlayerId,
			details.TimePeriod,
			details.Queue,
		}).Update(
		map[string]interface{}{
			"w":  r.Row.Field("w").Add(details.Wins),
			"l":  r.Row.Field("l").Add(details.Losses),
			"t":  r.Row.Field("t").Add(details.TimePlayed),
			"k":  r.Row.Field("k").Add(details.Kills),
			"a":  r.Row.Field("a").Add(details.Assists),
			"d":  r.Row.Field("d").Add(details.Deaths),
			"dk": r.Row.Field("dk").Add(details.DoubleKills),
			"tk": r.Row.Field("tk").Add(details.TripleKills),
			"qk": r.Row.Field("qk").Add(details.QuadraKills),
			"pk": r.Row.Field("pk").Add(details.PentaKills),
			"g":  r.Row.Field("g").Add(details.GoldEarned),
			"m":  r.Row.Field("m").Add(details.MinionsKilled),
			"n":  r.Row.Field("n").Add(details.MonstersKilled),
			"wp": r.Row.Field("wp").Add(details.WardsPlaced),
			"wk": r.Row.Field("wk").Add(details.WardsKilled),
			"b": map[string]interface{}{
				"w": r.Row.Field("b").Field("w").Add(details.Blue.Wins),
				"l": r.Row.Field("b").Field("l").Add(details.Blue.Losses),
			},
			"r": map[string]interface{}{
				"w": r.Row.Field("r").Field("w").Add(details.Red.Wins),
				"l": r.Row.Field("r").Field("l").Add(details.Red.Losses),
			},
		}).RunWrite(activeSession)

	if isDisconnected(err) {
		return ErrDisconnected
	} else if err != nil {
		return err
	}

	if resp.Replaced+resp.Unchanged == 0 {
		// Doesn't exist, insert new
		resp, err := r.Table("detailed").Insert(details).RunWrite(activeSession)
		if isDisconnected(err) {
			return ErrDisconnected
		} else if err != nil {
			return err
		} else if resp.Inserted == 0 {
			return ErrInsertDiscrepancy
		}
	}

	return nil
}

func AddToBasicPlayer(details dataFormat.BasicNumberOf) error {
	if !IsConnected {
		go Connect()
		return ErrDisconnected
	}

	resp, err := r.Table("basic").GetAllByIndex("ippqc",
		[]interface{}{
			details.InternalPlayerId,
			details.TimePeriod,
			details.Queue,
			details.Champion,
		}).Update(
		map[string]interface{}{
			"w":  r.Row.Field("w").Add(details.Wins),
			"l":  r.Row.Field("l").Add(details.Losses),
			"t":  r.Row.Field("t").Add(details.TimePlayed),
			"k":  r.Row.Field("k").Add(details.Kills),
			"a":  r.Row.Field("a").Add(details.Assists),
			"d":  r.Row.Field("d").Add(details.Deaths),
			"g":  r.Row.Field("g").Add(details.GoldEarned),
			"m":  r.Row.Field("m").Add(details.MinionsKilled),
			"n":  r.Row.Field("n").Add(details.MonstersKilled),
			"wp": r.Row.Field("wp").Add(details.WardsPlaced),
		}).RunWrite(activeSession)

	if isDisconnected(err) {
		return ErrDisconnected
	} else if err != nil {
		return err
	}

	if resp.Replaced+resp.Unchanged == 0 {
		// Doesn't exist, insert new
		resp, err := r.Table("basic").Insert(details).RunWrite(activeSession)
		if isDisconnected(err) {
			return ErrDisconnected
		} else if err != nil {
			return err
		} else if resp.Inserted == 0 {
			return ErrInsertDiscrepancy
		}
	}

	return nil
}

func DeletePlayer(player dataFormat.Player) error {
	if !IsConnected {
		go Connect()
		return ErrDisconnected
	}

	_, err := r.Table("players").GetAllByIndex("pi", player.SummonerId).
		Filter(map[string]string{"r": player.Region}).Delete().
		RunWrite(activeSession)
	if err != nil {
		return err
	}

	LastPlayerUpdate = time.Now()

	return nil
}

func CreatePlayer(player dataFormat.Player) (string, error) {
	if !IsConnected {
		go Connect()
		return "", ErrDisconnected
	}

	// Check if player exists already
	c, err := r.Table("players").GetAllByIndex("pi", player.SummonerId).
		Filter(map[string]string{"r": player.Region}).Field("id").
		Run(activeSession)

	if isDisconnected(err) {
		return "", ErrDisconnected
	} else if err != nil {
		return "", err
	}

	internalId := ""
	err = c.One(&internalId)
	c.Close()

	LastPlayerUpdate = time.Now()

	if err == nil {
		// Update new summoner name
		revel.WARN.Println("database: updating player summoner name for " +
			internalId)
		_, err := r.Table("players").Get(internalId).Update(
			map[string]string{
				"sn": player.SummonerName,
				"nn": player.NormalizedName,
			}).RunWrite(activeSession)
		if isDisconnected(err) {
			return "", ErrDisconnected
		} else if err != nil {
			return "", err
		}

		return internalId, nil
	} else if isDisconnected(err) {
		return "", ErrDisconnected
	} else if err != r.ErrEmptyResult {
		return "", err
	}

	changes, err := r.Table("players").Insert(player).RunWrite(activeSession)
	if isDisconnected(err) {
		return "", ErrDisconnected
	} else if err != nil {
		return "", err
	}

	if len(changes.GeneratedKeys) == 0 {
		return "", errors.New("database: missing generated keys")
	}

	return changes.GeneratedKeys[0], nil
}

func GetUpdatePlayers() ([]dataFormat.Player, error) {
	if !IsConnected {
		go Connect()
		return []dataFormat.Player{}, ErrDisconnected
	}

	c, err := r.Table("players").Between(r.MinVal, time.Now(),
		r.BetweenOpts{Index: "nu"}).
		OrderBy(r.OrderByOpts{Index: "nu"}).Pluck("p", "pi", "r", "id").
		Limit(3000).Run(activeSession)

	if isDisconnected(err) {
		return []dataFormat.Player{}, ErrDisconnected
	} else if err != nil {
		return []dataFormat.Player{}, err
	}

	players := []dataFormat.Player{}
	err = c.All(&players)

	if isDisconnected(err) {
		return []dataFormat.Player{}, ErrDisconnected
	} else if err != nil {
		return []dataFormat.Player{}, err
	}

	return players, nil
}

func GetLongUpdatePlayers() ([]dataFormat.Player, error) {
	if !IsConnected {
		go Connect()
		return []dataFormat.Player{}, ErrDisconnected
	}

	c, err := r.Table("players").Between(r.MinVal, time.Now(),
		r.BetweenOpts{Index: "nl"}).
		OrderBy(r.OrderByOpts{Index: "nl"}).Pluck("p", "pi", "r", "id").
		Limit(3000).Run(activeSession)

	if isDisconnected(err) {
		return []dataFormat.Player{}, ErrDisconnected
	} else if err != nil {
		return []dataFormat.Player{}, err
	}

	players := []dataFormat.Player{}
	err = c.All(&players)

	if isDisconnected(err) {
		return []dataFormat.Player{}, ErrDisconnected
	} else if err != nil {
		return []dataFormat.Player{}, err
	}

	return players, nil
}

func UpdatePlayerInformation(player dataFormat.Player, data interface{}) error {
	if !IsConnected {
		go Connect()
		return ErrDisconnected
	}

	_, err := r.Table("players").GetAllByIndex("id", player.InternalId).
		Update(data).RunWrite(activeSession)
	if isDisconnected(err) {
		return ErrDisconnected
	} else if err != nil {
		return err
	}

	LastPlayerUpdate = time.Now()

	return nil
}
