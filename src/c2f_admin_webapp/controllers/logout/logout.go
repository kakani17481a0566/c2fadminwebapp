package logout

import (
	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"
	"c2f_admin_webapp/session"
	"log"
	"runtime/debug"

	"github.com/astaxie/beego"
)

type Logout struct {
	// beego.Controller
	middleware.BaseController
}

func (c *Logout) Get() {
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	defer func() {
		if l_exception := recover(); l_exception != nil {
			stack := debug.Stack()
			log.Println(beego.AppConfig.String("loglevel"), "Exception", string(stack))
			session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
			c.Data["DisplayMessage"] = "Something went wrong.Please Contact CustomerCare."
			c.TplName = "error/error.html"
		}
		return
	}()
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

	c.Redirect("/", 302)
	return
}
