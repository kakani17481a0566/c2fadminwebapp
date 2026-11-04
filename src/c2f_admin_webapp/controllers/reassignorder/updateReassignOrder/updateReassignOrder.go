package updateReassignOrder

import (
	"c2f_admin_webapp/session"
	"runtime/debug"
	"time"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"
	"errors"
	"fmt"
	"strings"

	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/services"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

type Row struct {
	Id              string
	CategoryName    string
	SubCategoryName string
	ProductName     string
	ProductDesc     string
	Productstatus   string
	CreatedBy       string
}

type OderDeatils struct {
	Orderid         string          `json:"Orderid,omitempty"`
	DeliveryDate    string          `json:"DeliveryDate,omitempty"`
	DeliveryTime    string          `json:"DeliveryTime,omitempty"`
	Price           string          `json:"Price,omitempty"`
	Status          string          `json:"Status,omitempty"`
	CustomerDeatils CustomerDeatils `json:"CustomerDeatils"`
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

type Display struct {
	Fields []Field
}

type Field struct {
	Id    string
	Name  string
	Email string
}

type UpdateReassignOrder struct {
	// beego.Controller
	middleware.BaseController
}

func (c *UpdateReassignOrder) Get() {
	AdminId := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User Start")
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

			c.TplName = "reassignorder/updateReassignOrder/updateReassignOrder.html"

			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "reassignorder/updateReassignOrder/updateReassignOrder.html"

			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User Page Success")
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

	id := sess.Get("uid").(string)
	c.Data["id"] = id

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

	row, err := db.Db.Query(`select id, order_date,oder_id,status, deliveryuser_mob, create_date,JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.clusterName')),JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.totalSellPrice')),
	JSON_UNQUOTE(JSON_EXTRACT(customer_details, '$.name')),
	JSON_UNQUOTE(JSON_EXTRACT(customer_details, '$.mobile')),
	JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.productSKU')),
	JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.productUnits')),
	delivery_dates,JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.deliveryTo')),
	JSON_UNQUOTE(JSON_EXTRACT(customer_details, '$.address1')),
	JSON_UNQUOTE(JSON_EXTRACT(customer_details, '$.address2'))
	from Oder2Delivery  WHERE id =?  limit 1`, AdminId)
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

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	for i := range data {

		c.Data["Id"] = data[i][0]
		c.Data["Orderdate"] = data[i][1]
		c.Data["OrderID"] = data[i][2]
		c.Data["Status"] = data[i][3]
		c.Data["DeliveryUserMobile"] = data[i][4]
		c.Data["ClusterName"] = data[i][6]
		c.Data["TotalPrice"] = data[i][7]
		c.Data["CustomerName"] = data[i][8]
		c.Data["CustomerMobile"] = data[i][9]
		c.Data["ProductSKU"] = data[i][10]
		c.Data["ProductUnits"] = data[i][11]
		c.Data["Deliverydate"] = data[i][12]
		c.Data["DeliveryCompletionDate"] = data[i][13]
		c.Data["Address"] = data[i][14] + data[i][15]

	}

	ClusterName := data[0][6]
	DeliveryUserMobile := data[0][4]

	fmt.Println("ClusterName", ClusterName)
	fmt.Println("DeliveryUserMobile", DeliveryUserMobile)

	if utype == "CLUSTERMANAGER" {

		data2, err2 := services.GetClusterAsPerMapping(id)
		if err2 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err2)
			err2 = errors.New("Cluster user fetch Failed")
			return
		}
		if len(data2) <= 0 {
			err2 = errors.New("Cluster user  Not Found")
			return
		}

		data1, err1 := services.GetClustermappedDeliveryuserForRessignOrder(data2[0][1], ClusterName, DeliveryUserMobile)
		if err1 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Delivery user fetch Failed")
			return
		}
		if len(data1) <= 0 {
			err1 = errors.New("Delivery user  Not Found")
			return
		}
		var Dis Display
		for i := range data1 {
			var d Field
			d.Id = data1[i][0]
			d.Name = data1[i][1]
			d.Email = data1[i][2]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

	} else if utype == "CLUSTEREXECUTIVE" {

		data2, err2 := services.GetClusterAsPerMappingExec(id)
		if err2 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err2)
			err2 = errors.New("Cluster user fetch Failed")
			return
		}
		if len(data2) <= 0 {
			err2 = errors.New("Cluster user  Not Found")
			return
		}

		data1, err1 := services.GetClustermappedDeliveryuserForRessignOrder(data2[0][1], ClusterName, DeliveryUserMobile)
		if err1 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Delivery user fetch Failed")
			return
		}
		if len(data1) <= 0 {
			err1 = errors.New("Delivery user  Not Found")
			return
		}
		var Dis Display
		for i := range data1 {
			var d Field
			d.Id = data1[i][0]
			d.Name = data1[i][1]
			d.Email = data1[i][2]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

	} else {

		data1, err1 := services.GetDeliveryuserForRessignOrder(ClusterName, DeliveryUserMobile)
		if err1 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Delivery user fetch Failed")
			return
		}
		if len(data1) <= 0 {
			err1 = errors.New("Delivery user  Not Found")
			return
		}
		var Dis Display
		for i := range data1 {
			var d Field
			d.Id = data1[i][0]
			d.Name = data1[i][1]
			d.Email = data1[i][2]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

	}

	auth, err := utils.IsAuthorized(utype, "ReassignOrder")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	return

}

func (c *UpdateReassignOrder) Post() {

	AdminId := c.Ctx.Input.Param(":AdminID")

	log.Println(beego.AppConfig.String("loglevel"), "Info", "add asset post page")
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
			c.TplName = "reassignorder/updateReassignOrder/updateReassignOrder.html"

			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Delivery User Updated Successfully"
			c.TplName = "reassignorder/updateReassignOrder/updateReassignOrder.html"

			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer  Page Success")
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

	//input_status := c.Input().Get("input_status")
	//input_remarks := c.Input().Get("input_remarks")
	input_order_id := c.Input().Get("input_order_id")
	input_delivery_date := c.Input().Get("input_delivery_date")
	//	input_product_sku := c.Input().Get("input_product_sku")
	input_delivery_completiondate := c.Input().Get("input_delivery_completiondate")
	select_deliveruser_name := c.Input().Get("select_deliveruser_name")

	res, err := db.Db.Exec(`UPDATE Oder2Delivery SET deliveryuser_mob = ?
	WHERE id = ? AND delivery_dates=?`, select_deliveruser_name, AdminId, input_delivery_date)
	if err != nil {
		err = errors.New("Delivery updation failed")
		return
	}

	i, err := res.RowsAffected()
	if err != nil || i == 0 {
		err = errors.New("Delivery updation failed")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", " input_delivery_completiondate - ", input_delivery_completiondate)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", " input_delivery_date - ", input_delivery_date)

	if input_delivery_completiondate == input_delivery_date {

		log.Println(beego.AppConfig.String("loglevel"), "Debug", " input_delivery_date - ", input_delivery_date)

		res1, err := db.Db.Exec(`UPDATE CustomerOrderTxn SET status = "DELIVERY_COMPLETED"
	WHERE order_id = ?`, input_order_id)
		if err != nil {
			err = errors.New("CustomerOrderTxn Delivery updation failed")
			return
		}

		i1, err := res1.RowsAffected()
		if err != nil || i1 == 0 {
			err = errors.New(" CustomerOrderTxn Delivery updation failed")
			return
		}
	}

	// row, err := db.Db.Query(`select id, order_date,oder_id,status, deliveryuser_mob, create_date,JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.ClusterName')),JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.TotalSellPrice')),JSON_UNQUOTE(JSON_EXTRACT(customer_details, '$.Name')),JSON_UNQUOTE(JSON_EXTRACT(customer_details, '$.Mobile')) from Oder2Delivery  WHERE id =?  limit 1`, AdminId)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Unable to get user data")
	// 	return
	// }
	// defer sql.Close(row)
	// _, data, err := sql.Scan(row)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Unable to get user data")
	// 	return
	// }
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	// if len(data) <= 0 {
	// 	err = errors.New("User data not found")
	// 	return
	// }

	// log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	// for i := range data {

	// 	c.Data["Id"] = data[i][0]
	// 	c.Data["Orderdate"] = data[i][1]
	// 	c.Data["OrderID"] = data[i][2]
	// 	c.Data["Status"] = data[i][3]
	// 	c.Data["ClusterName"] = data[i][6]
	// 	c.Data["TotalPrice"] = data[i][7]
	// 	c.Data["CustomerName"] = data[i][8]
	// 	c.Data["CustomerMobile"] = data[i][9]

	// }

	return
}
func daysBetween(a, b time.Time) int {
	if a.After(b) {
		a, b = b, a
	}

	days := -a.YearDay()
	for year := a.Year(); year < b.Year(); year++ {
		days += time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC).YearDay()
	}
	days += b.YearDay()

	return days
}

func date(s string) time.Time {
	d, _ := time.Parse("2006-01-02", s)
	return d
}
