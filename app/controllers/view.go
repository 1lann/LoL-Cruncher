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
	"cruncher/app/models/query"
	"github.com/revel/revel"
	"strings"
)

type View struct {
	*revel.Controller
}

func (c View) Index() revel.Result {
	return c.Render()
}

func (c View) About() revel.Result {
	return c.Render()
}

func (c View) Robots() revel.Result {
	return c.RenderText("Sitemap: https://lolcruncher.tk/sitemap.xml")
}

func (c View) Request(region, name string) revel.Result {
	if !(region == "na" || region == "euw" || region == "eune" ||
		region == "lan" || region == "las" || region == "oce" ||
		region == "br" || region == "ru" || region == "kr" ||
		region == "tr") {
		c.Flash.Error("Sorry, that region isn't supported!")
		return c.Redirect(View.Index)
	}

	region = strings.ToLower(region)

	resolvedName, player, new, err := query.GetStats(name, region)
	if err != nil {
		if err.Error() == "database error" {
			return c.RenderTemplate("errors/database.html")
		} else if err.Error() == "database down" {
			return c.RenderTemplate("errors/down.html")
		} else if err.Error() == "Not Found" {
			c.Flash.Error("Sorry, that summoner could not be found!")
			return c.Redirect(View.Index)
		} else {
			c.Flash.Error("Could not connect to Riot Games' servers! Try again later.")
			return c.Redirect(View.Index)
		}
	}

	if strings.Trim(resolvedName, " ") != strings.Trim(name, " ") {
		return c.Redirect("/" + region + "/" + strings.Trim(resolvedName, " "))
	}

	c.RenderArgs["new"] = new
	c.RenderArgs["player"] = player
	c.RenderArgs["name"] = resolvedName
	c.RenderArgs["titleName"] = resolvedName + " - LoL Cruncher"
	c.RenderArgs["description"] = "View " + resolvedName + "'s statistics for all queues (since " + player.RecordStart + ")"
	return c.Render()
}
