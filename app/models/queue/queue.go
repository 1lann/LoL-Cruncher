
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

package queue

// Handles rate limiting, queueing requests, and everything inbetween

import (
	"time"
// "cruncher/app/models/database"
)

var RateLimitingBlock time.Time

// var monitorRunning bool = false

// var longMonitorRunning bool = false

func RateLimitBlock(seconds int) {
	RateLimitingBlock = time.Now().Add(time.Second * time.Duration(seconds))
}

func IsRateBlocked() bool {
	return time.Now().Before(RateLimitingBlock)
}

// func RecordMonitor() {
// 	if monitorRunning {
// 		revel.ERROR.Println("Monitor already running!")
// 		return
// 	}
// 	monitorRunning = true
// 	for true {
// 		// Scan through database in next update timestamp order
// 		// func GetRecentGames(id string, region string) ([]dataFormat.Game, error) {

// 		time.Sleep(time.Duration(60) * time.Minute)
// 	}
// }

// func LongTermMonitor() {
// 	if longMonitorRunning {
// 		revel.ERROR.Println("Long term mointor already running!")
// 		return
// 	}
// 	longMonitorRunning = true
// 	for true {

// 		time.Sleep(time.Duration(24) * time.Hour)
// 	}
// }
