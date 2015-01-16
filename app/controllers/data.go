
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

package controllers

import (
	"github.com/revel/revel"
	"cruncher/app/models/database"
	"cruncher/app/models/dataFormat"
	"encoding/json"
	"time"
)

type Data struct {
	*revel.Controller
}

var lastCacheUpdate time.Time
var cacheResponse string

type playerUpdate struct {
	Time int64
	Players []dataFormat.BrowserPlayer
}

func getDatabaseUpdates() string {
	if (database.LastPlayerUpdate.After(lastCacheUpdate)) {
		results, resp := database.GetBrowserPlayers()
		if resp != database.Yes {
			revel.ERROR.Println("getDatabaseUpdates error!")
			return "error"
		}

		fullResult := playerUpdate{
			Time: database.LastPlayerUpdate.Unix() + 1,
			Players: results,
		}

		result, err := json.Marshal(fullResult)
		if err != nil {
			revel.ERROR.Println("getDatabaseUpdates JSON marshal error")
			revel.ERROR.Println(err)
			return "error"
		}

		lastCacheUpdate = database.LastPlayerUpdate
		cacheResponse = string(result)
		return string(result)
	} else {
		return cacheResponse
	}
}

func (c Data) CheckDatabaseUpdates(lastupdate int) revel.Result {
	if (database.LastPlayerUpdate.After(time.Unix(int64(lastupdate), 0))) {
		databaseUpdates := getDatabaseUpdates()
		return c.RenderText(databaseUpdates)
	} else {
		return c.RenderText("false")
	}
}
