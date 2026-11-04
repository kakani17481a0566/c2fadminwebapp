package deliveryorderinvoice

import (
	"c2f_admin_webapp/session"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"
	"encoding/json"
	"errors"
	"runtime/debug"
	"strings"

	"c2f_admin_webapp/model/db"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

type Row struct {
	Id          string
	ClusterName string
	Location    string
	Zipcode     string
	Status      string
}

type OderDeatils struct {
	Orderid           string `json:"Orderid,omitempty"`
	OderDate          string `json:"OderDate,omitempty"`
	DeliveryFrom      string `json:"DeliveryFrom,omitempty"`
	DeliveryTo        string `json:"DeliveryTo,omitempty"`
	DeliveryTime      string `json:"DeliveryTime,omitempty"`
	OrderType         string `json:"OrderType,omitempty"`
	RepeatOrder       string `json:"RepeatOrder,omitempty"`
	PaymentInfo       string `json:"PaymentInfo,omitempty"`
	CategoryName      string `json:"CategoryName,omitempty"`
	SubcategroryName  string `json:"SubcategroryName,omitempty"`
	ProductName       string `json:"ProductName,omitempty"`
	ProductQty        string `json:"ProductQty,omitempty"`
	ProductPackageQty string `json:"ProductPackageQty,omitempty"`
	ProductUnits      string `json:"ProductUnits,omitempty"`
	ProductPrice      string `json:"ProductPrice,omitempty"`
	TotalPrice        string `json:"TotalPrice,omitempty"`
	TotalSellPrice    string `json:"TotalSellPrice,omitempty"`
	ClusterName       string `json:"ClusterName,omitempty"`
}

type CustomerDeatils struct {
	Name     string `json:"Name,omitempty"`
	Address1 string `json:"Address1,omitempty"`
	Address2 string `json:"Address2,omitempty"`
	City     string `json:"City,omitempty"`
	Zipcode  string `json:"Zipcode,omitempty"`
	Mobile   string `json:"Mobile,omitempty"`
	Email    string `json:"Email,omitempty"`
}

type Deliveryorderinvoice struct {
	// beego.Controller
	middleware.BaseController
}

func (c *Deliveryorderinvoice) Get() {
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

			c.TplName = "invoice/deliveryOderInvoice/deliveryOderInvoice.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "invoice/deliveryOderInvoice/deliveryOderInvoice.html"
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

	auth, err := utils.IsAuthorized(utype, "Invoice")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	row, err := db.Db.Query(`select oder_id, order_date,deliveryuser_mob,customer_details,oder_deatils,status  from Oder2Delivery WHERE oder_id =?  limit 1`, AdminId)
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
		err = errors.New("User data not found")
		return
	}

	var cd CustomerDeatils
	json.Unmarshal([]byte(data[0][3]), &cd)

	var od OderDeatils
	json.Unmarshal([]byte(data[0][4]), &od)
	c.Data["Status"] = data[0][5]
	c.Data["TotalAmount"] = od.TotalPrice
	c.Data["OrderID"] = data[0][0]
	c.Data["OrderDate"] = data[0][1]
	c.Data["DF"] = od.DeliveryFrom
	c.Data["DT"] = od.DeliveryTo
	c.Data["DeliveryTime"] = od.DeliveryTime
	c.Data["OrderType"] = od.OrderType
	c.Data["RepeatOrder"] = od.RepeatOrder
	c.Data["ClusterName"] = od.ClusterName
	c.Data["CategoryName"] = od.CategoryName
	c.Data["SubcategoryName"] = od.SubcategroryName
	c.Data["ProductName"] = od.ProductName
	c.Data["ProductQty"] = od.ProductQty
	c.Data["ProductPackageQty"] = od.ProductPackageQty
	c.Data["ProductPrice"] = od.ProductPrice
	c.Data["TotalSellPrice"] = od.TotalSellPrice
	c.Data["ProductUnits"] = od.ProductUnits
	c.Data["CustomerName"] = cd.Name
	c.Data["Address1"] = cd.Address1
	c.Data["Address2"] = cd.Address2
	c.Data["City"] = cd.City
	c.Data["Zipcode"] = cd.Zipcode
	c.Data["Mobile"] = cd.Mobile
	c.Data["Email"] = cd.Email

	return

}
