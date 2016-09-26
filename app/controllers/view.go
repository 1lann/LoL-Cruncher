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
	"cruncher/app/models/dataFormat"
	"cruncher/app/models/database"
	"cruncher/app/models/query"
	"strings"
	"time"

	"github.com/revel/revel"
	"github.com/revel/revel/cache"
	// "time"
)

type View struct {
	*revel.Controller
}

func (c View) Index() revel.Result {
	if !database.IsConnected {
		go database.Connect()
		return c.RenderTemplate("errors/down.html")
	}

	return c.Render()
}

func (c View) About() revel.Result {
	return c.Render()
}

func (c View) Robots() revel.Result {
	return c.RenderText("Sitemap: https://lolcruncher.tk/sitemap.xml")
}

func (c View) Request(region, name string) revel.Result {
	region = strings.ToLower(region)

	if !(region == "na" || region == "euw" || region == "eune" ||
		region == "lan" || region == "las" || region == "oce" ||
		region == "br" || region == "ru" || region == "kr" ||
		region == "tr") {
		c.Flash.Error("Sorry, that region isn't supported!")
		return c.Redirect(View.Index)
	}

	var err error
	player := dataFormat.PlayerData{}
	new := false

	// player, new, err = query.GetStats(name, region, false)

	if err = cache.Get(region+":"+dataFormat.NormalizeName(name),
		&player); err != nil {
		player, new, err = query.GetStats(name, region, false)
		revel.WARN.Println("Storing", name, region, "in cache.")
		go cache.Set(region+":"+dataFormat.NormalizeName(player.SummonerName),
			player, time.Hour*2)
	}

	if err != nil {
		if err == query.ErrDatabaseError {
			return c.RenderTemplate("errors/database.html")
		} else if err == query.ErrDatabaseDisconnected {
			return c.RenderTemplate("errors/down.html")
		} else if err == query.ErrNotFound {
			c.Flash.Error("Sorry, that summoner could not be found!")
			return c.Redirect(View.Index)
		} else if err == query.ErrAPIError {
			c.Flash.Error("Could not connect to Riot Games' servers! Try again later.")
			return c.Redirect(View.Index)
		} else {
			c.Flash.Error("An unknown error has occured, please try again in a few seconds.")
			return c.Redirect(View.Index)
		}
	}

	resolvedName := player.SummonerName

	if strings.Trim(resolvedName, " ") != strings.Trim(name, " ") {
		return c.Redirect("/" + region + "/" + strings.Trim(resolvedName, " "))
	}

	player.RecordStartString = player.RecordStart.Format("2 January 2006")

	c.RenderArgs["new"] = new
	c.RenderArgs["player"] = player
	c.RenderArgs["name"] = resolvedName
	c.RenderArgs["titleName"] = resolvedName + " - LoL Cruncher"
	c.RenderArgs["description"] = "View " + resolvedName +
		"'s League of Legends statistics and champion breakdowns for all " +
		"queues (since " + player.RecordStart.Format("2 January 2006") + ")"
	return c.Render()
}
