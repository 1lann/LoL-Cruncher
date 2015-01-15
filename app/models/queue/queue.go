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
