package customerorderBill

import (
	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/model/utils"
	"c2f_admin_webapp/services"
	"c2f_admin_webapp/session"

	"strings"

	"runtime/debug"

	"errors"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

type Row struct {
	Id            string
	OrderId       string
	OrderDate     string
	DeliveryDates string
	Status        string
	Remarks       string

	CustomerName   string
	CustomerAdd1   string
	CustomerAdd2   string
	CustomerCity   string
	CustomerMobile string

	ProductName       string
	ProductQty        string
	ProductPackageQty string
	ProductSKU        string
	UnitPrice         string
	ProductUnit       string

	DeliveryTime string
	DeliveryFrom string
	DeliveryTo   string

	ClusterName string
}
type Display5 struct {
	Fields5 []Field5
}

type Field5 struct {
	Id     string
	Mobile string
}
type CustomerOrderBill struct {
	// beego.Controller
	middleware.BaseController
}

func (c *CustomerOrderBill) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Reports Page Start")
	var err error
	var Autherr error
	sessErr := false
	defer func() {
		if l_exception := recover(); l_exception != nil {
			stack := debug.Stack()
			log.Println(beego.AppConfig.String("loglevel"), "Exception", string(stack))
			session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
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
			c.TplName = "bills/customerorderBill/customerorderBill.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Reports Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "bills/customerorderBill/customerorderBill.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Reports Page Success")
		}
		return
	}()
	utils.SetHTTPHeader(c.Ctx)
	sess, err := session.GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System is unable to process your request.Please contact customer care")
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

	data, err := services.GetConsumer()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Consumer fetch Failed")
		return
	}
	if len(data) <= 0 {
		err = errors.New("Consumer  Not Found")
		return
	}
	var Dis Display5
	for i := range data {
		var d Field5
		d.Id = data[i][0]
		d.Mobile = data[i][1]
		Dis.Fields5 = append(Dis.Fields5, d)
	}
	c.Data["Dis"] = Dis
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

	auth, err := utils.IsAuthorized(utype, "CustOrderBill")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	return
}

func (c *CustomerOrderBill) Post() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchOrder Page Start")
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
			c.TplName = "bills/customerorderBill/customerorderBill.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchOrder Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "bills/customerorderBill/customerorderBill.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchOrder  Page Success")
		}
		return
	}()
	utils.SetHTTPHeader(c.Ctx)
	sess, err := session.GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System is unable to process your request.Please contact customer care")
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

	auth, err := utils.IsAuthorized(utype, "CustOrderBill")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	input_SupplyorderId := c.Input().Get("input_SupplyorderId")

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_SupplyorderId", input_SupplyorderId)

	// c.Data["RRN"] = input_rnn_no
	// c.Data["Mask"] = input_card_mask
	// c.Data["DateRange"] = c.Input().Get("daterange")
	// c.Data["Status"] = input_status
	// c.Data["SupplyorderId"] = input_SupplyorderId
	// tmp := strings.Split(c.Input().Get("daterange"), " - ")

	// from := ""
	// to := ""
	// if len(tmp) == 2 {
	// 	from = tmp[0]
	// 	to = tmp[1]
	// }

	row, err := db.Db.Query(`select 
	oder_id,
	order_date,
	delivery_dates,
    status,
    remarks,
    
    JSON_UNQUOTE(JSON_EXTRACT(customer_details, '$.name')),
	JSON_UNQUOTE(JSON_EXTRACT(customer_details, '$.address1')),
	JSON_UNQUOTE(JSON_EXTRACT(customer_details, '$.address2')),
	JSON_UNQUOTE(JSON_EXTRACT(customer_details, '$.city')),
	JSON_UNQUOTE(JSON_EXTRACT(customer_details, '$.mobile')),
	
	JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.productName')),
    JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.productPackageQty')),
	JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.productQty')),
    JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.productSKU')),
    JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.productPrice')),
	JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.productUnits')),
	
	JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.deliveryTime')),
	JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.deliveryFrom')),
	JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.deliveryTo')),
    JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.clusterName'))
         from Oder2Delivery
    where 
    oder_id=?`, input_SupplyorderId)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch order")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch order")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)
	var result []Row

	for i := range data {
		var r Row

		r.OrderId = data[i][0]
		r.OrderDate = data[i][1]
		r.DeliveryDates = data[i][2]
		r.Status = data[i][3]
		r.Remarks = data[i][4]

		r.CustomerName = data[i][5]
		r.CustomerAdd1 = data[i][6]
		r.CustomerAdd2 = data[i][7]
		r.CustomerCity = data[i][8]
		r.CustomerMobile = data[i][9]

		r.ProductName = data[i][10]
		r.ProductPackageQty = data[i][11] + data[i][12]
		r.ProductQty = data[i][12]
		r.ProductSKU = data[i][13]
		r.UnitPrice = data[i][14]
		r.ProductUnit = data[i][15]

		r.DeliveryTime = data[i][16]
		r.DeliveryFrom = data[i][17]
		r.DeliveryTo = data[i][18]
		r.ClusterName = data[i][19]

		result = append(result, r)
	}
	c.Data["CustomerData"] = result

	return
}
