package error

import (
	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"
	"c2f_admin_webapp/session"
	"log"

	"github.com/astaxie/beego"
)

type Error struct {
	// beego.Controller
	middleware.BaseController
}

func (c *Error) Error404() {
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "404 Path", c.Ctx.Request.URL.Path)

	c.TplName = "error/error.html"
}

func (c *Error) Error501() {
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	utils.SetHTTPHeader(c.Ctx)
	sess, _ := session.GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)

	log.Println(beego.AppConfig.String("loglevel"), "Info", "UserName Nil Found")
	uname := sess.Get("uname")
	if uname != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Info", "UserName Nil Found")
		session.SeTLogoutSession(uname.(string))
	}
	sess.SessionRelease(c.Ctx.ResponseWriter)
	session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
	c.Data["DisplayMessage"] = "501, server error"
	c.TplName = "error/error.html"
}
