package viewLoadMap

import (
	"c2f_admin_webapp/session"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"
	"runtime/debug"

	"errors"
	//	"strings"

	// "c2f_admin_webapp/model/db"
	"c2f_admin_webapp/services"

	"github.com/astaxie/beego"

	// "c2f_admin_webapp/model/database/sql"
	"c2f_admin_webapp/model/util/log"
)

type Row struct {
	ClusterLocation []string
}
type ViewLoadMap struct {
	// beego.Controller
	middleware.BaseController
}

func (c *ViewLoadMap) Get() {
	Name := c.Ctx.Input.Param(":Name")
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	var err error
	var Autherr error
	sessErr := false
	defer func() {

		if l_exception := recover(); l_exception != nil {
			stack := debug.Stack()
			log.Println(beego.AppConfig.String("loglevel"), "Exception", string(stack))
			session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
			c.Data["DisplayMessage"] = "Something went wrong.Please Contact CustomerCare."
			c.TplName = "error/error.html"
		}
		if Autherr != nil {
			c.Data["DisplayMessage"] = Autherr.Error()
			c.TplName = "error/error.html"
			return
		}
		if err != nil {
			if sessErr == true {
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Redirecting to login")
				c.Redirect(beego.AppConfig.String("LOGIN_PATH"), 302)

			} else {
				c.Data["DisplayMessage"] = err.Error()
			}
			c.TplName = "loadmap/viewLoadMap/viewLoadMap.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "loadmap/viewLoadMap/viewLoadMap.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer  Page Success")
		}
		return
	}()

	sess, err := session.GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System is unable to process your request.Please contact customer care")
		sessErr = true
		return
	}

	if err = session.ValidateSession(sess); err != nil {
		sess.SessionRelease(c.Ctx.ResponseWriter)
		session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		sessErr = true
		return
	}
	defer func() {
		utils.EventLogs(c.Ctx, sess, c.Ctx.Input.Method(), c.Input(), c.Data, err)
		sess.SessionRelease(c.Ctx.ResponseWriter)
	}()

	utype := sess.Get("user_type").(string)
	c.Data["utype"] = utype
	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "utype", utype)

	language := sess.Get("language").(string)
	c.Data["language"] = language
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "language", language)

	maptoken := sess.Get("token")
	c.Data["maptoken"] = maptoken
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "maptoken", maptoken)

	// var delhi string
	delhi := sess.Get("Delhi")
	c.Data["delhi"] = delhi
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "delhi", delhi)

	// delhi1 := strings.Trim(delhi, "\"")
	// c.Data["delhi1"] = delhi1

	auth, err := utils.IsAuthorized(utype, "LoadMap")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	var OrderLocationsArray []string
	var CustomerName []string

	deliveryuser_id, _ := services.GetDeliveryuserid(Name)
	if err != nil {
		beego.Error(err)
		err = errors.New("Invalid Delivery User  ")
		return
	}

	cluster_id, _ := services.GetClusterid1(deliveryuser_id)
	if err != nil {
		beego.Error(err)
		err = errors.New("Invalid Delivery User  ")
		return
	}

	data5, err := services.GetClusterLocation(cluster_id)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get cluster data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Cluster  Data - ", data5, "\nData len - ", len(data5))
	if len(data5) <= 0 {
		err = errors.New("cluster data not found")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data5)

	for i := range data5 {

		c.Data["ClusterLocation"] = data5[i][0]
		c.Data["ClusterName"] = data5[i][1]

	}

	data1, err := services.GetOrderLocation(Name)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get order data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Customer Order  Data - ", data1, "\nData len - ", len(data1))
	if len(data1) <= 0 {
		err = errors.New("order data not found")
		return
	}
	c.Data["length"] = len(data1)

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data1)

	for i := range data1 {

		CustomerName = append(CustomerName, data1[i][1]+","+data1[i][4]+"-"+data1[i][2]+data1[i][3])
		OrderLocationsArray = append(OrderLocationsArray, data1[i][0])

		c.Data["customername"] = CustomerName
		c.Data["orderlocations"] = OrderLocationsArray

	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "OrderLocationsArray - ", OrderLocationsArray)

	return

}

func trimQuote(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}
