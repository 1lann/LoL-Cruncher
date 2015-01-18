
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
	"github.com/revel/revel"
	"gopkg.in/mgo.v2"
	"time"
)

var IsConnected bool
var isConnecting bool
var activeSession *mgo.Session
var players *mgo.Collection
var playerIds *mgo.Collection

func isDisconnected(err string) bool {
	if err == "EOF" || err == "no reachable servers" ||
		err == "Closed explicitly" {
		return true
	} else {
		return false
	}
}

func Connect() {
	if !isConnecting {
		IsConnected = false

		if activeSession != nil {
			activeSession.Close()
		}

		isConnecting = true
		revel.INFO.Println("Connecting...")

		databaseIp, found := revel.Config.String("database.ip")

		if !found {
			revel.ERROR.Println("Missing database.ip in conf/app.conf!")
			panic("Missing database.ip in conf/app.conf!")
			return
		}

		session, err := mgo.DialWithTimeout(databaseIp, time.Second*3)
		if err != nil {
			isConnecting = false
			IsConnected = false
			revel.ERROR.Println(err)
			return
		}

		session.SetMode(mgo.Monotonic, true)
		session.SetSafe(&mgo.Safe{})
		session.SetSyncTimeout(time.Second*3)
		session.SetSocketTimeout(time.Second*3)

		activeSession = session

		players = session.DB("cruncher").C("players")
		playerIds = session.DB("cruncher").C("playerids")

		IsConnected = true
		isConnecting = false
	}
}
