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

package dataFormat

import (
	"strings"
	"time"
)

type SideStats struct {
	Wins   int
	Losses int
}

type DetailedNumberOf struct {
	InternalPlayerId string `gorethink:"ip" json:"-"`
	TimePeriod       string `gorethink:"p" json:"p"`
	Queue            string `gorethink:"q" json:"q"`

	Wins           int `gorethink:"w" json:"w"`
	Losses         int `gorethink:"l" json:"l"`
	TimePlayed     int `gorethink:"t" json:"t"`
	Kills          int `gorethink:"k" json:"k"`
	Assists        int `gorethink:"a" json:"a"`
	Deaths         int `gorethink:"d" json:"d"`
	DoubleKills    int `gorethink:"dk" json:"dk"`
	TripleKills    int `gorethink:"tk" json:"tk"`
	QuadraKills    int `gorethink:"qk" json:"qk"`
	PentaKills     int `gorethink:"pk" json:"pk"`
	GoldEarned     int `gorethink:"g" json:"g"`
	MinionsKilled  int `gorethink:"m" json:"m"`
	MonstersKilled int `gorethink:"n" json:"n"`
	WardsPlaced    int `gorethink:"wp" json:"wp"`
	WardsKilled    int `gorethink:"wk" json:"wk"`
	Blue           struct {
		Wins   int `gorethink:"w" json:"w"`
		Losses int `gorethink:"l" json:"l"`
	} `gorethink:"b"`
	Red struct {
		Wins   int `gorethink:"w" json:"w"`
		Losses int `gorethink:"l" json:"l"`
	} `gorethink:"r"`
}

type BasicNumberOf struct {
	InternalPlayerId string `gorethink:"ip" json:"-"`
	Champion         string `gorethink:"c" json:"c"`
	TimePeriod       string `gorethink:"p" json:"p"`
	Queue            string `gorethink:"q" json:"q"`

	Wins           int `gorethink:"w" json:"w"`
	Losses         int `gorethink:"l" json:"l"`
	TimePlayed     int `gorethink:"t" json:"t"`
	Kills          int `gorethink:"k" json:"k"`
	Assists        int `gorethink:"a" json:"a"`
	Deaths         int `gorethink:"d" json:"d"`
	GoldEarned     int `gorethink:"g" json:"g"`
	MinionsKilled  int `gorethink:"m" json:"m"`
	MonstersKilled int `gorethink:"n" json:"n"`
	WardsPlaced    int `gorethink:"wp" json:"wp"`
}

type Player struct {
	// nr for normalized name region (internally)
	Region         string    `gorethink:"r" json:"region"`
	Tier           string    `gorethink:"t" json:"-"`
	SummonerId     string    `gorethink:"pi" json:"-"`
	InternalId     string    `gorethink:"id,omitempty" json:"-"`
	SummonerName   string    `gorethink:"sn" json:"summonerName"`
	NormalizedName string    `gorethink:"nn" json:"-"`
	RecordStart    string    `gorethink:"rs" json:"-"` // Date of first ever game recorded
	NextUpdate     time.Time `gorethink:"nu" json:"-"`
	NextLongUpdate time.Time `gorethink:"nl" json:"-"`
	ProcessedGames []string  `gorethink:"p" json:"-"`
}

type PlayerData struct {
	Detailed       []DetailedNumberOf `gorethink:"detailed" json:"detailed"`
	Basic          []BasicNumberOf    `gorethink:"basic" json:"basic"`
	SummonerName   string             `gorethink:"sn" json:"summonerName"`
	Region         string             `gorethink:"r" json:"r"`
	RecordStart    string             `gorethink:"rs" json:"rs"`
	ProcessedGames []string           `gorethink:"p" json:"-"`
}

type Champion struct {
	Name        string
	Title       string
	SquareURL   string
	SplashURL   string
	PortraitURL string
}

type Game struct {
	DidWin         bool
	IsOnBlue       bool
	IsNormals      bool // AKA Not Custom
	ChampionId     string
	Duration       int
	Id             string
	Type           string
	Kills          int
	Assists        int
	Deaths         int
	DoubleKills    int
	TripleKills    int
	QuadraKills    int
	PentaKills     int
	GoldEarned     int
	MinionsKilled  int
	MonstersKilled int
	WardsPlaced    int
	WardsKilled    int
	YearMonth      string
	Date           time.Time
}

func NormalizeName(name string) string {
	return strings.ToLower(strings.Replace(name, " ", "", -1))
}
