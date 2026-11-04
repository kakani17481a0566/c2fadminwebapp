package main

import (
	_ "c2f_admin_webapp/routers"
	"os"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	_ "github.com/astaxie/beego/session/redis"
)

var RedirectHttp = func(ctx *context.Context) {
	if !ctx.Input.IsSecure() {
		url := "https://" + ctx.Input.Domain() + ":" + beego.AppConfig.String("HttpsPort") + ctx.Input.URI()
		ctx.Redirect(302, url)
	}
}

func main() {
	if beego.AppConfig.String("EnableHTTPS") == "true" {
		beego.InsertFilter("/", beego.BeforeRouter, RedirectHttp) // for http://mysite
		beego.InsertFilter("*", beego.BeforeRouter, RedirectHttp) // for http://mysite/*
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("WEBSITES_PORT")
	}
	if port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			beego.BConfig.Listen.HTTPPort = p
		}
	}

	beego.Run()
}
