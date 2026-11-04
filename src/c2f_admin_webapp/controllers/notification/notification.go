package notification

import (
	"c2f_admin_webapp/session"

	"c2f_admin_webapp/model/utils"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/services"

	"errors"
	//	"fmt"
	"runtime/debug"

	//"strconv"
	"strings"

	"c2f_admin_webapp/model/db"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"
	"c2f_admin_webapp/model/util/log"
)

type Display1 struct {
	Fields1 []Field1
}

type Field1 struct {
	Supplyordertxn_id  string
	ProducerMobile     string
	Status             string
	ProductName        string
	ClusterName        string
	Productqty         string
	PPackageQty        string
	Createdate         string
	ProductSupplyUnits string
	NotificationTime   string
}

type Notification struct {
	// beego.Controller
	middleware.BaseController
}

func (c *Notification) Get() {
	//AdminId := c.Ctx.Input.Param(":AdminID")
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

			c.TplName = "notification/notification.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "notification/notification.html"
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

	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1

	auth, err := utils.IsAuthorized(utype, "Notification")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	data10, err := services.GetProducers1()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Producer fetch Failed")
		return
	}
	data51, err := services.GetNotificationDuration(data10)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Time Duration for Notification fetch Failed")
		return
	}

	var Dis1 Display1
	for i := range data10 {
		var d Field1
		d.ProducerMobile = data10[i][0]
		d.Status = data10[i][1]
		d.ProductName = data10[i][2]
		d.ClusterName = data10[i][3]
		d.Productqty = data10[i][4]
		d.PPackageQty = data10[i][5] + data10[i][4]
		d.Supplyordertxn_id = data10[i][6]
		d.Createdate = data10[i][7]
		d.ProductSupplyUnits = data10[i][8]
		d.NotificationTime = data51[i]

		Dis1.Fields1 = append(Dis1.Fields1, d)
	}
	c.Data["Dis1"] = Dis1
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis1)

	row11, err := db.Db.Query(`SELECT status,count(producer_mob)as PSum FROM SupplyOrderTxn where status="DEMAND_ORDER_CREATED"`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row11)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row11)
	_, data11, err := sql.Scan(row11)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data11, "\nData len - ", len(data11))
	if len(data11) <= 0 {
		err = errors.New("Unable to get dashboard data")
		return
	}

	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "PMob - ", data8[0][0])
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "PStatus - ", data11[0][0])
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "PCount - ", data11[0][1])

	//c.Data["producermobile"] = data8[0][0]
	c.Data["status"] = data11[0][0]
	c.Data["totalproducercount"] = data11[0][1]

	return
}
