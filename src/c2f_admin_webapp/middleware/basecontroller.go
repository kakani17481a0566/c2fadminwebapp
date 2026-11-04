package middleware

import (
	"c2f_admin_webapp/session"
	// "log"

	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
}

func (c *BaseController) Prepare() {
	path := c.Ctx.Request.URL.Path
	if path == "/Login" || path == "/Logout" {
		return
	}

	sess, err := session.GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	if err != nil {
		// log.Println("[BaseController] Session start failed:", err)
		c.Redirect("/Login", 302)
		return
	}

	storedIP := sess.Get("bind_ip")
	storedUA := sess.Get("bind_ua")
	currentIP := c.Ctx.Input.IP()
	currentUA := c.Ctx.Request.UserAgent()

	// log.Println("[BaseController] storedIP:", storedIP, "currentIP:", currentIP)
	// log.Println("[BaseController] storedUA:", storedUA, "currentUA:", currentUA)

	if storedIP != nil && storedUA != nil {
		if storedIP != currentIP || storedUA != currentUA {
			// log.Println("[SECURITY] Session Hijack Detected! Destroying session.")
			session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
			c.Redirect("/Login", 302)
			return
		}
	}
}
