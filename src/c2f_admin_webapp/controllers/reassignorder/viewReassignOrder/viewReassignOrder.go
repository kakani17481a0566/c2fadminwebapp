package viewReassignOrder

import (
	"c2f_admin_webapp/session"

	"c2f_admin_webapp/model/utils"
	"runtime/debug"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/db"
	"errors"
	"strings"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

type Row struct {
	Id                 string
	SupplyTxnID        string
	ProductName        string
	ProductSKU         string
	ProductQTY         string
	ProductPackageQTY  string
	ClusterBuyPrice    string
	SupplyUnits        string
	TotalPrice         string
	Status             string
	DeliveryUserMobile string
	CustomerName       string
	CustomerMobile     string
	Address            string
}
type ViewReassignOrder struct {
	// beego.Controller
	middleware.BaseController
}

func (c *ViewReassignOrder) Get() {
	AdminId := c.Ctx.Input.Param(":AdminID")
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

			c.TplName = "reassignorder/viewReassignOrder/viewReassignOrder.html"

			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "reassignorder/viewReassignOrder/viewReassignOrder.html"
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

	clustername := sess.Get("clustername").(string)
	clustername1 := strings.ToUpper(clustername)
	c.Data["clustername"] = clustername1
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "clustername", clustername1)

	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1

	auth, err := utils.IsAuthorized(utype, "ReassignOrder")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	row, err := db.Db.Query(`select id, order_date,oder_id,status, deliveryuser_mob, create_date,JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.clusterName')),delivery_dates,JSON_UNQUOTE(JSON_EXTRACT(customer_details, '$.name')) ,JSON_UNQUOTE(JSON_EXTRACT(customer_details, '$.mobile')),
	JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.productName')) ,
	JSON_UNQUOTE(JSON_EXTRACT(customer_details, '$.address1')),
	JSON_UNQUOTE(JSON_EXTRACT(customer_details, '$.address2'))
	from Oder2Delivery WHERE id =? AND status ="DELIVERY_ASSIGNED"  limit 1`, AdminId)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Order Status Invalid")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	for i := range data {

		c.Data["Id"] = data[i][0]
		c.Data["Orderdate"] = data[i][1]
		c.Data["OrderID"] = data[i][2]
		c.Data["Status"] = data[i][3]
		c.Data["ClusterName"] = data[i][6]
		c.Data["Deliverydate"] = data[i][7]
		c.Data["DeliveryUserMobile"] = data[i][4]
		c.Data["CustomerName"] = data[i][8]
		c.Data["CustomerMobile"] = data[i][9]
		c.Data["ProductName"] = data[i][10]
		c.Data["Address"] = data[i][11] + data[i][12]

	}

	return

}
