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
	r "github.com/dancannon/gorethink"
	"github.com/revel/revel"
	"io"
	"strings"
)

var IsConnected bool
var isConnecting bool
var activeSession *r.Session

func isDisconnected(err error) bool {
	if err == r.ErrBadConn || err == r.ErrConnectionClosed ||
		err == r.ErrNoConnections || err == r.ErrNoConnectionsStarted ||
		err == io.EOF ||
		(err != nil && strings.Contains(err.Error(), "broken pipe")) {
		go Connect()
		return true
	} else {
		return false
	}
}

func databaseRecover() {
	if r := recover(); r != nil {
		revel.ERROR.Println("database: recovered from database driver panic")
		revel.ERROR.Println(r)
	}
}

func Connect() {
	if !isConnecting {
		isConnecting = true

		defer func() {
			isConnecting = false
		}()

		defer databaseRecover()

		IsConnected = false

		if activeSession != nil {
			activeSession.Close()
		}

		revel.ERROR.Println("Attempting to reconnect...")

		databaseIp, found := revel.Config.String("database.ip")

		if !found {
			revel.ERROR.Println("Missing database.ip in conf/app.conf!")
			panic("Missing database.ip in conf/app.conf!")
			return
		}

		databasePassword, hasPassword := revel.Config.String("database.password")

		if !hasPassword {
			revel.WARN.Println("No database.password in conf/app.conf, " +
				"assuming development mode with no login.")
		}

		session, err := r.Connect(r.ConnectOpts{
			Address:  databaseIp,
			Database: "cruncher",
			AuthKey:  databasePassword,
			MaxIdle:  100,
			MaxOpen:  100,
		})
		if err != nil {
			IsConnected = false
			revel.ERROR.Println(err)
			return
		}

		activeSession = session

		IsConnected = true
	}
}
