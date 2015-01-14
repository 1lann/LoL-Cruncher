package controllers

import (
	"github.com/revel/revel"
	//  "cruncher/app/models/riotapi"
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
