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
	if err == "EOF" || err == "no reachable servers" {
		return true
	} else {
		return false
	}
}

func Connect() bool {
	if !isConnecting {
		IsConnected = false

		if activeSession != nil {
			activeSession.Close()
		}

		isConnecting = true
		revel.INFO.Println("Connecting...")
		session, err := mgo.DialWithTimeout("127.0.0.1", time.Second*3)
		if err != nil {
			isConnecting = false
			IsConnected = false
			revel.ERROR.Println(err)
			return false
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
		return true
	} else {
		return false
	}
}

func init() {
	Connect();
}
