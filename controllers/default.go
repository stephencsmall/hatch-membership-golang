package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

type CheckinController struct {
	beego.Controller
}

type SignupController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.TplName = "splash.html"
}

func (c *CheckinController) Get() {
	c.TplName = "checkin.html"
}

func (c *SignupController) Get() {
	c.TplName = "signup.html"
}