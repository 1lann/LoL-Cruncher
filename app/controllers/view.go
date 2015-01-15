package controllers

import (
	"github.com/revel/revel"
	"cruncher/app/models/query"
)

type View struct {
	*revel.Controller
}

func (c View) Index() revel.Result {
	return c.Render()
}

func (c View) Test() revel.Result {
	c.RenderArgs["errorCode"] = "ErrNotFound"
	return c.RenderTemplate("errors/database.html");
}

func (c View) Request() revel.Result {
	// func GetStats(name string, region string) (string, dataFormat.Player, error) {
	name, player, err := query.GetStats("1lann", "na")
	if err != nil {
		revel.ERROR.Println(err)
		return c.RenderTemplate("errors/database.html");
	}
	c.RenderArgs["player"] = player
	c.RenderArgs["name"] = name
	return c.Render()
}
