package database

import (
	"cruncher/app/models/dataFormat"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"strings"
)

const (
	Yes   = 1
	No    = 2
	Down  = 3
	Error = 4
)

type playerId struct {
	Id string
	Region string
	Name string
	Normalized string
}

var NewPlayerUpdate bool = true

// returns id, resolved name, and response
func GetSummonerID(name string, region string) (string, string, int) {
	if !IsConnected {
		go Connect()
		return "", "", Down
	}

	query := bson.M{"normalized": strings.ToLower(name), "region": region}
	var result playerId
	err := playerIds.Find(query).One(&result)

	if err != nil {
		if err == mgo.ErrNotFound {
			return "", "", No
		} else if isDisconnected(err.Error()) {
			go Connect()
			return "", "", Down
		} else {
			printOut := "GetSummonerID Database Error: "
			revel.ERROR.Println(printOut + err.Error())
			return "", "", Error
		}
	}

	return result.Id, result.Name, Yes
}

func StoreSummonerID(name string, id string, region string) int {
	if !IsConnected {
		go Connect()
		return Down
	}

	// Check for identical id, if so, delete the old one
	query := bson.M{"id": id, "region": region}
	var result playerId
	err := playerIds.Find(query).One(&result)
	if err == nil {
		// Clear it
		printOut := "Summoner name change detected. From: " + result.Name +
			" to " + name + " with ID: " + id
		revel.WARN.Println(printOut)
		err = playerIds.Remove(query)
	} else if err != mgo.ErrNotFound {
		if isDisconnected(err.Error()) {
			go Connect()
			return Down
		} else {
			printOut := "StoreSummonerID Database Error: "
			revel.ERROR.Println(printOut + err.Error())
			return Error
		}
	}

	NewPlayerUpdate = true

	// Continue here if everything is clean
	newPlayer := playerId{
		Id: id,
		Region: region,
		Name: name,
		Normalized: strings.ToLower(name),
	}
	err = playerIds.Insert(newPlayer)

	if err != nil {
		if isDisconnected(err.Error()) {
			go Connect()
			return Down
		} else {
			printOut := "StoreSummonerID Database Error: "
			revel.ERROR.Println(printOut + err.Error())
			return Error
		}
	}
	return Yes
}

func GetSummonerData(id string, region string) (dataFormat.Player, int) {
	if !IsConnected {
		go Connect()
		return dataFormat.Player{}, Down
	}

	query := bson.M{"id": id, "region": region}
	var result dataFormat.Player
	err := players.Find(query).One(&result)

	if err != nil {
		if err == mgo.ErrNotFound {
			return dataFormat.Player{}, No
		} else if isDisconnected(err.Error()) {
			go Connect()
			return dataFormat.Player{}, Down
		} else {
			printOut := "GetSummonerData Database Error: "
			revel.ERROR.Println(printOut + err.Error())
			return dataFormat.Player{}, Error
		}
	}

	return result, Yes
}

func GetBrowserPlayers() ([]dataFormat.BrowserPlayer, int) {
	if !IsConnected {
		go Connect()
		return []dataFormat.BrowserPlayer{}, Down
	}

	var result []dataFormat.BrowserPlayer
	query := bson.M{"name": 1, "region": 1}
	err := playerIds.Find(nil).Select(query).All(&result)

	if err != nil {
		if err == mgo.ErrNotFound {
			return []dataFormat.BrowserPlayer{}, Yes
		} else if isDisconnected(err.Error()) {
			go Connect()
			return []dataFormat.BrowserPlayer{}, Down
		} else {
			printOut := "GetBrowserPlayers Database Error: "
			revel.ERROR.Println(printOut + err.Error())
			return []dataFormat.BrowserPlayer{}, Error
		}
	}

	return result, Yes
}

func StoreSummonerData(player dataFormat.Player) int {
	if !IsConnected {
		go Connect()
		return Down
	}

	query := bson.M{"id": player.Id, "region": player.Region}
	err := players.Update(query, bson.M{"$set": player})
	if err != nil {
		if err == mgo.ErrNotFound {
			revel.INFO.Println("Player does not exist, adding")
			err = players.Insert(player)
			if err != nil {
				if isDisconnected(err.Error()) {
					go Connect()
					return Down
				} else {
					printOut := "StoreSummonerData Database Error: "
					revel.ERROR.Println(printOut + err.Error())
					return Error
				}
			}
			return Yes
		} else if isDisconnected(err.Error()) {
			go Connect()
			return Down
		}
	}
	return Yes
}

func GetUpdatePlayers() ([]dataFormat.BasicPlayer, int) {
	if !IsConnected {
		go Connect()
		return []dataFormat.BasicPlayer{}, Down
	}

	var results []dataFormat.BasicPlayer

	query := bson.M{
		"region": 1,
		"id": 1,
		"recordstart": 1,
		"nextupdate": 1,
		"nextlongupdate": 1,
	}

	it := players.Find(nil).Select(query).Sort("-nextupdate").Iter()

	var player dataFormat.BasicPlayer
	for it.Next(&player) {
		if player.NextUpdate.IsZero() {
			revel.ERROR.Println(`Zero time for next update from GetUpdates.
				This error is not self resolving and should be manually fixed.
				However in most cases, the error will be visible to the end
				user`)
			revel.ERROR.Println(player)
		} else {
			if player.NextUpdate.Before(time.Now()) {
				results = append(results, player)
			} else {
				break
			}
		}
	}

	err := it.Err()
	if err != nil {
		if isDisconnected(err.Error()) {
			go Connect()
			return []dataFormat.BasicPlayer{}, Down
		} else {
			printOut := "GetUpdates Database Error: "
			revel.ERROR.Println(printOut + err.Error())
			return []dataFormat.BasicPlayer{}, Error
		}
	}

	return results, Yes
}

func GetLongUpdatePlayers() ([]dataFormat.BasicPlayer, int) {
	if !IsConnected {
		go Connect()
		return []dataFormat.BasicPlayer{}, Down
	}

	var results []dataFormat.BasicPlayer

	query := bson.M{
		"region": 1,
		"id": 1,
		"recordstart": 1,
		"nextupdate": 1,
		"nextlongupdate": 1,
	}

	it := players.Find(nil).Select(query).Sort("-nextlongupdate").Iter()

	var player dataFormat.BasicPlayer
	for it.Next(&player) {
		if player.NextLongUpdate.IsZero() {
			revel.ERROR.Println(`Zero time for next update from GetLongUpdates.
				This error is not self resolving and should be manually fixed.
				However in most cases, the error will be visible to the end
				user`)
			revel.ERROR.Println(player)
		} else {
			if player.NextLongUpdate.Before(time.Now()) {
				results = append(results, player)
			} else {
				break
			}
		}
	}

	err := it.Err()
	if err != nil {
		if isDisconnected(err.Error()) {
			go Connect()
			return []dataFormat.BasicPlayer{}, Down
		} else {
			printOut := "GetLongUpdates Database Error: "
			revel.ERROR.Println(printOut + err.Error())
			return []dataFormat.BasicPlayer{}, Error
		}
	}

	return results, Yes
}
