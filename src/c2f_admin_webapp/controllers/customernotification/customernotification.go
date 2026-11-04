package customernotification

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

type Display2 struct {
	Fields2 []Field2
}

type Field2 struct {
	CustomerMobile   string
	ClusterName      string
	ProductName      string
	PPackageQty      string
	Productqty       string
	Status           string
	OrderUnits       string
	OrderId          string
	Createdate       string
	NotificationTime string
	ClusterType      string
}
type CustomerNotification struct {
	// beego.Controller
	middleware.BaseController
}

func (c *CustomerNotification) Get() {
	//AdminId := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Notification Start")
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

			c.TplName = "customernotification/customernotification.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "customernotification/customernotification.html"
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

	//---------------Customer Nodification--------------

	data12, err := services.GetCustomerOrderDetailsForNotification()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Customer Order fetch Failed with Status ORDER_RECEIVED")
		return
	}

	data52, err := services.GetCustomerOrderNotificationDuration(data12)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Time Duration for Notification fetch Failed")
		return
	}

	var Dis2 Display2
	for i := range data12 {
		var d Field2
		d.CustomerMobile = data12[i][0]
		d.ClusterName = data12[i][1]
		d.ProductName = data12[i][2]
		d.PPackageQty = data12[i][3] + data12[i][4]
		d.Productqty = data12[i][4]
		d.OrderUnits = data12[i][5]
		d.OrderId = data12[i][6]
		d.Status = data12[i][7]
		d.Createdate = data12[i][8]
		d.ClusterType = data12[i][9]
		d.NotificationTime = data52[i]

		Dis2.Fields2 = append(Dis2.Fields2, d)
	}
	c.Data["Dis2"] = Dis2
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Customer Notification Dashboard Data - ", Dis2)

	row13, err := db.Db.Query(`SELECT status,count(customer_mob)as CSum FROM CustomerOrderTxn where status="ORDER_RECEIVED"`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row13)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row13)
	_, data13, err := sql.Scan(row13)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data13, "\nData len - ", len(data13))
	if len(data13) <= 0 {
		err = errors.New("Unable to get dashboard data")
		return
	}

	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "PMob - ", data8[0][0])
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "CStatus - ", data13[0][0])
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "CCount - ", data13[0][1])

	//c.Data["producermobile"] = data8[0][0]
	c.Data["status"] = data13[0][0]
	c.Data["totalcustomercount"] = data13[0][1]

	return
}
