package selectDeliveryOder

import (
	"c2f_admin_webapp/session"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/model/utils"
	"c2f_admin_webapp/services"
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"
	"time"

	"strconv"
	"strings"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

type DeliveryOrders struct {
	Deliveryuser string     `json:"deliveryuser"`
	Oderinfo     []Oderinfo `json:"oderinfo"`
}
type Oderinfo struct {
	Orderid           string          `json:"orderid"`
	OderDate          string          `json:"oderDate"`
	DeliveryFrom      string          `json:"deliveryFrom"`
	DeliveryTo        string          `json:"deliveryTo"`
	DeliveryTime      string          `json:"deliveryTime"`
	OrderType         string          `json:"orderType"`
	RepeatOrder       string          `json:"repeatOrder"`
	CategoryName      string          `json:"categoryName"`
	SubcategroryName  string          `json:"subcategroryName"`
	ProductName       string          `json:"productName"`
	ProductSKU        string          `json:"productSKU"`
	ProductQty        string          `json:"productQty"`
	ProductPackageQty string          `json:"productPackageQty"`
	ProductUnits      string          `json:"productUnits"`
	ProductPrice      string          `json:"productPrice"`
	TotalPrice        string          `json:"totalPrice"`
	TotalSellPrice    string          `json:"totalSellPrice"`
	ClusterName       string          `json:"clusterName"`
	Name              string          `json:"name"`
	Address1          string          `json:"address1"`
	Address2          string          `json:"address2"`
	City              string          `json:"city"`
	Zipcode           string          `json:"zipcode"`
	Mobile            string          `json:"mobile"`
	Email             string          `json:"email"`
	CustomerDeatils   CustomerDeatils `json:"CustomerDeatils"`
}
type CustomerDeatils struct {
	Name        string `json:"name"`
	Address1    string `json:"address1"`
	Address2    string `json:"address2"`
	City        string `json:"city"`
	Zipcode     string `json:"zipcode"`
	Mobile      string `json:"mobile"`
	Email       string `json:"email"`
	Latlong     string `json:"latlong"`
	LocationURL string `json:"locationURL"`
}

type Order struct {
	OrderID string
	Msg     interface{}
}

// type Sound struct {
// 	Description string
// 	Authority   string
// }

type DeliveryOderinfo struct {
	SubOrderid        string `json:"orderid"`
	DeliveryDate      string `json:"deliveryDate"`
	DeliveryTime      string `json:"deliveryTime"`
	ProductName       string `json:"productName"`
	ProductSKU        string `json:"productSKU"`
	ProductQty        string `json:"productQty"`
	ProductPackageQty string `json:"productPackageQty"`
	ProductUnits      string `json:"productUnits"`
	ProductPrice      string `json:"productPrice"`
	ClusterName       string `json:"clusterName"`
	Name              string `json:"name"`
	Address1          string `json:"address1"`
	Address2          string `json:"address2"`
	City              string `json:"city"`
	Zipcode           string `json:"zipcode"`
	Mobile            string `json:"mobile"`
	Email             string `json:"email"`
}

type Row struct {
	Id           string
	TimeStamp    string
	Oder_id      string
	Cluster_name string
	Product_name string
	Status       string
	Name         string
	Address1     string
	Address2     string
	City         string
	Zipcode      string
	Latlong      string
	LocationURL  string
	Email        string
	Mobile       string
	Custdetails  string
	CustomerID   string
	OrderType    string
	DeliveryDate string
	Units        string
}

type Display struct {
	Fields []Field
}

type Field struct {
	Id      string
	Name    string
	Email   string
	Cluster string
}

type SelectDeliveryOder struct {
	// beego.Controller
	middleware.BaseController
}
type GetDeliveryDateOrders struct {
	// beego.Controller
	middleware.BaseController
}

func (c *SelectDeliveryOder) Get() {
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
			c.TplName = "deliveryorder/selectDeliveryOder/selectDeliveryOder.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "deliveryorder/selectDeliveryOder/selectDeliveryOder.html"
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
	id := sess.Get("uid").(string)
	c.Data["id"] = id

	language := sess.Get("language").(string)
	c.Data["language"] = language
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "language", language)

	clustername := sess.Get("clustername").(string)
	clustername1 := strings.ToUpper(clustername)
	c.Data["clustername"] = clustername1
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "clustername", clustername1)

	auth, err := utils.IsAuthorized(utype, "DeliveryOder")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	return
}

func (c *GetDeliveryDateOrders) Post() {
	// msg11 := ""
	// log.Println(beego.AppConfig.String("loglevel"), "Info", "node post page")
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
			c.TplName = "deliveryorder/selectDeliveryOder/selectDeliveryOder.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SelectDeloveryOrder Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "deliveryorder/selectDeliveryOder/selectDeliveryOder.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SelectDeloveryOrder  Page Success")
		}
		return
	}()
	utils.SetHTTPHeader(c.Ctx)
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
	defer sess.SessionRelease(c.Ctx.ResponseWriter)

	utype := sess.Get("user_type").(string)
	c.Data["utype"] = utype

	language := sess.Get("language").(string)
	c.Data["language"] = language
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "language", language)

	id := sess.Get("uid").(string)
	c.Data["id"] = id

	clustername := sess.Get("clustername").(string)
	clustername1 := strings.ToUpper(clustername)
	c.Data["clustername"] = clustername1
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "clustername", clustername1)

	uid := sess.Get("uid").(string)
	c.Data["id"] = uid
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "id", uid)

	auth, err := utils.IsAuthorized(utype, "DeliveryOder")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	delivery_date := c.Input().Get("delivery_date")
	if delivery_date == "" {
		err = errors.New("Delivery Date cannot be empty")
		return

	}
	fmt.Println("Selected Delivery date is :- ", delivery_date)

	c.Data["datebutton"] = "HIDE"
	c.Data["duserbutton"] = "SHOW"

	if utype == "CLUSTERMANAGER" {

		data1, err1 := services.GetClusterAsPerMapping(id)
		if err1 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
			err1 = errors.New("Cluster user fetch Failed")
			return
		}
		if len(data1) <= 0 {
			err = errors.New("Cluster user  Not Found")
			return
		}

		data, err := services.GetClustermappedDeliveryuser(data1[0][1])
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Delivery user fetch Failed")
			return
		}
		if len(data) <= 0 {
			err = errors.New("Delivery user  Not Found")
			return
		}
		var Dis Display
		for i := range data {
			var d Field
			d.Id = data[i][0]
			d.Name = data[i][1]
			d.Email = data[i][2]
			d.Cluster = data[i][7]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

		auth, err := utils.IsAuthorized(utype, "DeliveryOder")
		if !auth {
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
			Autherr = errors.New("UnAuthorized")
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

		row, err := db.Db.Query(` select id,
 order_date, 
	order_id,
   cluster_name,
product_name,
status,
customer_id,
order_type,
orderdeliver_fromdate,
orderdeliver_todate,
order_units
 from CustomerOrderTxn where status= "ORDER_RECEIVED" AND order_purchagetype="KRISHIVALA" AND cluster_type="SOURCE" AND cluster_name=?`, data1[0][3])

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch order")
			return
		}
		defer sql.Close(row)
		_, data4, err := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch order")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", data4)
		var result []Row

		for i := range data4 {

			// For Loop to get All Orders with condition ORDER_RECEIVED and KRISHIVALA and SOURCE and not bulk orders
			var r Row

			for rd := rangeDate(date(data4[i][8]), date(data4[i][9])); ; {

				// For Loop to get every date between from date and to date

				date := rd()
				if date.IsZero() {
					break
				}

				if delivery_date == date.Format("2006-01-02") {

					// Condition to check if Selected ( Delivery Date ) and  (Date) from the loop  between From and To date-range matches

					row, err := db.Db.Query(`select count(*) from Oder2Delivery WHERE oder_id = ? AND delivery_dates=?`, data4[i][2], date.Format("2006-01-02"))

					if err != nil {
						log.Println(beego.AppConfig.String("loglevel"), "Error", err)
						err = errors.New("Unable to fetch delivery date - order info")
						return
					}
					defer sql.Close(row)
					_, data, err := sql.Scan(row)
					if err != nil {
						log.Println(beego.AppConfig.String("loglevel"), "Error", err)
						err = errors.New("Unable to fetch  delivery date - order info")
						return
					}

					// log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

					countlen := data[0][0]

					// log.Println(beego.AppConfig.String("loglevel"), "Debug", "countlen", countlen)

					// Condition to check if (Date) from the loop  between From and To dates range is already assigned for delivery

					if countlen != "0" {

						fmt.Println("Already Delivery is Assigned for Order" + data4[i][2] + " on selected Delivery Date" + date.Format("2006-01-02"))

					} else {

						fmt.Println("Id :-" + data4[i][0])
						fmt.Println("OrderId :- " + data4[i][2])
						fmt.Println("Delivery Date :-" + date.Format("2006-01-02"))
						r.Id = data4[i][0]
						r.TimeStamp = data4[i][1]
						r.Oder_id = data4[i][2]
						r.Cluster_name = data4[i][3]
						r.Product_name = data4[i][4]
						r.Status = data4[i][5]
						r.CustomerID = data4[i][6]
						r.OrderType = data4[i][7]
						r.DeliveryDate = date.Format("2006-01-02")
						r.Units = data4[i][10]

						// log.Println(beego.AppConfig.String("loglevel"), "Debug", "data4[i][2]", data4[i][2])

						fname, add1, add2, city, zipcode, email, mob, latlong, err := services.GetCustomerinfo(data4[i][6])
						if err != nil {
							log.Println(beego.AppConfig.String("loglevel"), "Error", err)
							err = errors.New("Invalid Order ID ")
							return
						}

						var tmp = fname + add1 + add2

						// log.Println(beego.AppConfig.String("loglevel"), "Debug", "customer info1", tmp)

						var tmp2 = city + zipcode + email
						var tmp3 = mob
						fmt.Println(latlong)
						// log.Println(beego.AppConfig.String("loglevel"), "Debug", "customer info2", tmp2)

						// log.Println(beego.AppConfig.String("loglevel"), "Debug", "customer info3", tmp3)

						r.Custdetails = tmp + "," + tmp2 + "," + tmp3
						//	r.LocationURL = locationurl
						result = append(result, r)

					}

				}

			}

		}
		c.Data["CustomerData"] = result

	} else if utype == "CLUSTEREXECUTIVE" {

		// CLUSTER-EXECUTIVE CODE

		data1, err1 := services.GetClusterAsPerMappingExec(id)
		if err1 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
			err1 = errors.New("Cluster user fetch Failed")
			return
		}
		if len(data1) <= 0 {
			err = errors.New("Cluster user  Not Found")
			return
		}

		data, err := services.GetClustermappedDeliveryuser(data1[0][1])
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Delivery user fetch Failed")
			return
		}
		if len(data) <= 0 {
			err = errors.New("Delivery user  Not Found")
			return
		}
		var Dis Display
		for i := range data {
			var d Field
			d.Id = data[i][0]
			d.Name = data[i][1]
			d.Email = data[i][2]
			d.Cluster = data[i][7]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

		auth, err := utils.IsAuthorized(utype, "DeliveryOder")
		if !auth {
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
			Autherr = errors.New("UnAuthorized")
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

		row, err := db.Db.Query(` select id,
 order_date, 
	order_id,
   cluster_name,
product_name,
status,
customer_id,
order_type,
orderdeliver_fromdate,
orderdeliver_todate,
order_units
 from CustomerOrderTxn where status= "ORDER_RECEIVED" AND order_purchagetype="KRISHIVALA" AND cluster_type="SOURCE" AND cluster_name=?`, data1[0][3])

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch order")
			return
		}
		defer sql.Close(row)
		_, data4, err := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch order")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", data4)
		var result []Row

		for i := range data4 {

			// For Loop to get All Orders with condition ORDER_RECEIVED,KRISHIVALA,SOURCE,not bulk order
			var r Row

			for rd := rangeDate(date(data4[i][8]), date(data4[i][9])); ; {

				// For Loop to get every date between from date and to date

				date := rd()
				if date.IsZero() {
					break
				}

				if delivery_date == date.Format("2006-01-02") {

					// Condition to check if Selected ( Delivery Date ) and  (Date) from the loop  between From and To dates range matches

					row, err := db.Db.Query(`select count(*) from Oder2Delivery WHERE oder_id = ? AND delivery_dates=?`, data4[i][2], date.Format("2006-01-02"))

					if err != nil {
						log.Println(beego.AppConfig.String("loglevel"), "Error", err)
						err = errors.New("Unable to fetch delivery date - order info")
						return
					}
					defer sql.Close(row)
					_, data, err := sql.Scan(row)
					if err != nil {
						log.Println(beego.AppConfig.String("loglevel"), "Error", err)
						err = errors.New("Unable to fetch  delivery date - order info")
						return
					}

					// log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

					countlen := data[0][0]

					// log.Println(beego.AppConfig.String("loglevel"), "Debug", "countlen", countlen)

					if countlen != "0" {

						// Condition to check if (Date) from the loop  between From and To dates range is already assigned for delivery

						fmt.Println("Already Delivery is Assigned for Order" + data4[i][2] + " on selected Delivery Date" + date.Format("2006-01-02"))

					} else {

						fmt.Println("Id :-" + data4[i][0])
						fmt.Println("OrderId :- " + data4[i][2])
						fmt.Println("Delivery Date :-" + date.Format("2006-01-02"))
						r.Id = data4[i][0]
						r.TimeStamp = data4[i][1]
						r.Oder_id = data4[i][2]
						r.Cluster_name = data4[i][3]
						r.Product_name = data4[i][4]
						r.Status = data4[i][5]
						r.CustomerID = data4[i][6]
						r.OrderType = data4[i][7]
						r.DeliveryDate = date.Format("2006-01-02")
						r.Units = data4[i][10]

						// log.Println(beego.AppConfig.String("loglevel"), "Debug", "data4[i][2]", data4[i][2])

						fname, add1, add2, city, zipcode, email, mob, latlong, err := services.GetCustomerinfo(data4[i][6])
						if err != nil {
							log.Println(beego.AppConfig.String("loglevel"), "Error", err)
							err = errors.New("Invalid Order ID ")
							return
						}

						var tmp = fname + add1 + add2

						// log.Println(beego.AppConfig.String("loglevel"), "Debug", "customer info1", tmp)

						var tmp2 = city + zipcode + email
						var tmp3 = mob
						fmt.Println(latlong)
						// log.Println(beego.AppConfig.String("loglevel"), "Debug", "customer info2", tmp2)

						// log.Println(beego.AppConfig.String("loglevel"), "Debug", "customer info3", tmp3)

						r.Custdetails = tmp + "," + tmp2 + "," + tmp3
						//	r.LocationURL = locationurl
						result = append(result, r)

					}

				}

			}

		}
		c.Data["CustomerData"] = result

	} else {

		// HOADMIN CODE

		data, err := services.GetDeliveryuser()
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Delivery user fetch Failed")
			return
		}
		if len(data) <= 0 {
			err = errors.New("Delivery user  Not Found")
			return
		}
		var Dis Display
		for i := range data {
			var d Field
			d.Id = data[i][0]
			d.Name = data[i][1]
			d.Email = data[i][2]
			d.Cluster = data[i][7]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

		auth, err := utils.IsAuthorized(utype, "DeliveryOder")
		if !auth {
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
			Autherr = errors.New("UnAuthorized")
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

		row, err := db.Db.Query(` select id,
 order_date, 
	order_id,
   cluster_name,
product_name,
status,
customer_id,
order_type,
orderdeliver_fromdate,
orderdeliver_todate,
order_units
 from CustomerOrderTxn where status= "ORDER_RECEIVED" AND  order_purchagetype="KRISHIVALA" AND cluster_type="SOURCE"`)

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch order")
			return
		}
		defer sql.Close(row)
		_, data4, err := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch order")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", data4)
		var result []Row

		for i := range data4 {

			// For Loop to get All Orders with condition ORDER_RECEIVED,KRISHIVALA,SOURCE,not bulk order
			var r Row

			for rd := rangeDate(date(data4[i][8]), date(data4[i][9])); ; {

				// For Loop to get every date between from date and to date

				date := rd()
				if date.IsZero() {
					break
				}

				if delivery_date == date.Format("2006-01-02") {

					// Condition to check if Selected ( Delivery Date ) and  (Date) from the loop  between From and To dates range matches

					row, err := db.Db.Query(`select count(*) from Oder2Delivery WHERE oder_id = ? AND delivery_dates=?`, data4[i][2], date.Format("2006-01-02"))

					if err != nil {
						log.Println(beego.AppConfig.String("loglevel"), "Error", err)
						err = errors.New("Unable to fetch delivery date - order info")
						return
					}
					defer sql.Close(row)
					_, data, err := sql.Scan(row)
					if err != nil {
						log.Println(beego.AppConfig.String("loglevel"), "Error", err)
						err = errors.New("Unable to fetch  delivery date - order info")
						return
					}

					// log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

					countlen := data[0][0]

					// log.Println(beego.AppConfig.String("loglevel"), "Debug", "countlen", countlen)

					if countlen != "0" {

						// Condition to check if (Date) from the loop  between From and To dates range is already assigned for delivery

						fmt.Println("Already Delivery is Assigned for Order" + data4[i][2] + " on selected Delivery Date" + date.Format("2006-01-02"))

					} else {

						fmt.Println("Id :-" + data4[i][0])
						fmt.Println("OrderId :- " + data4[i][2])
						fmt.Println("Delivery Date :-" + date.Format("2006-01-02"))
						r.Id = data4[i][0]
						r.TimeStamp = data4[i][1]
						r.Oder_id = data4[i][2]
						r.Cluster_name = data4[i][3]
						r.Product_name = data4[i][4]
						r.Status = data4[i][5]
						r.CustomerID = data4[i][6]
						r.OrderType = data4[i][7]
						r.DeliveryDate = date.Format("2006-01-02")
						r.Units = data4[i][10]

						// log.Println(beego.AppConfig.String("loglevel"), "Debug", "data4[i][2]", data4[i][2])

						fname, add1, add2, city, zipcode, email, mob, latlong, err := services.GetCustomerinfo(data4[i][6])
						if err != nil {
							log.Println(beego.AppConfig.String("loglevel"), "Error", err)
							err = errors.New("Invalid Order ID ")
							return
						}

						var tmp = fname + add1 + add2

						// log.Println(beego.AppConfig.String("loglevel"), "Debug", "customer info1", tmp)

						var tmp2 = city + zipcode + email
						var tmp3 = mob
						fmt.Println(latlong)
						// log.Println(beego.AppConfig.String("loglevel"), "Debug", "customer info2", tmp2)

						// log.Println(beego.AppConfig.String("loglevel"), "Debug", "customer info3", tmp3)

						r.Custdetails = tmp + "," + tmp2 + "," + tmp3
						//	r.LocationURL = locationurl
						result = append(result, r)

					}

				}

			}

		}
		c.Data["CustomerData"] = result

	}

	return
}

func (c *SelectDeliveryOder) Post() {
	msg11 := ""
	log.Println(beego.AppConfig.String("loglevel"), "Info", "node post page")
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
			c.TplName = "deliveryorder/selectDeliveryOder/selectDeliveryOder.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SelectDeliveryOrder Page Fail")
		} else {
			c.Data["DisplayMessage"] = msg11
			c.TplName = "deliveryorder/selectDeliveryOder/selectDeliveryOder.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SelectDeliveryOrder  Page Success")
		}
		return
	}()
	utils.SetHTTPHeader(c.Ctx)
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
	defer sess.SessionRelease(c.Ctx.ResponseWriter)

	utype := sess.Get("user_type").(string)
	c.Data["utype"] = utype

	language := sess.Get("language").(string)
	c.Data["language"] = language
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "language", language)

	clustername := sess.Get("clustername").(string)
	clustername1 := strings.ToUpper(clustername)
	c.Data["clustername"] = clustername1
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "clustername", clustername1)

	uid := sess.Get("uid").(string)
	c.Data["id"] = uid
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "id", uid)

	auth, err := utils.IsAuthorized(utype, "DeliveryOder")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	select_deliveruser_name := c.Input().Get("select_deliveruser_name")
	order_deliverydate := c.Input().Get("order_deliverydate")
	//order_id := c.Input().Get("programming")
	//select_cluster_name := c.Input().Get("select_cluster_name")
	Get_deliveruser_Values := c.Input().Get("Get_deliveruser_Values")
	Get_cluster_name := c.Input().Get("Get_cluster_name")

	// order_id1 := c.Input().Get("txnArr")
	order_id := c.Input().Get("txnArr")
	// fmt.Println("order_id1", order_id1)
	fmt.Println("order_id", order_id)

	if select_deliveruser_name == "" {
		err = errors.New(" Please Select Delivery user ")
		return
	}

	if order_id == "" {
		err = errors.New(" Please Select Check Box ")
		return
	}

	fmt.Println("Get_deliveruser_Values", Get_deliveruser_Values)
	fmt.Println("Get_cluster_name", Get_cluster_name)

	// Split the Get_cluster_name string into a slice of cluster names
	clusterNames := strings.Split(Get_cluster_name, ",")

	// Flag to check if there is any mismatch
	mismatchFound := false

	// Iterate through each cluster name and compare with Get_deliveruser_Values
	for _, cluster := range clusterNames {
		cluster = strings.TrimSpace(cluster) // Trim any whitespace around cluster names
		if cluster != Get_deliveruser_Values {
			mismatchFound = true
			break
		}
	}

	// Check if any mismatch was found
	if mismatchFound {
		fmt.Println("Both Cluster are different - Failed")
		err = errors.New("Incorrect cluster mapping for order with ID : " + order_id)
		msg11 = "Incorrect cluster mapping for order with ID :" + order_id
		return
	} else {
		fmt.Println("Both Cluster are same - Passed")
	}

	//var productsku1
	// order_id = services.Getinventory3(aunits)

	// if PT > aunits1 {

	// 	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "inside - ")

	// 	err = errors.New("Insufficent Available Inventory")
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	return
	// }

	// func Checkinventory(a, ou float64) {

	// if ou > a {
	// 	err = errors.New(" Insufficent Available Inventory ")
	// 	return

	// }
	//}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Delivery User - ", select_deliveruser_name)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "order_deliverydate - ", order_deliverydate)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Order ID's - ", order_id)
	fruits := strings.Split(order_id, ",")

	fmt.Println(fruits)
	lg := len(fruits)
	fmt.Println(lg)

	for i := range fruits {

		h := fruits[i]

		fmt.Println(h)

		// err = CheckMapExists(select_cluster_name, select_deliveruser_name)
		// if err != nil {
		// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		// 	err = errors.New("Product Mapping Not Exists")
		// 	return
		// }

		// err = CheckmapAlreadyExists(select_deliveruser_name, h)
		// if err != nil {
		// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		// 	err = errors.New("User 2 Oder mapping already Exists")
		// 	return
		// }

		// status, err := Checkorderstatus(select_deliveruser_name, h)
		// if err != nil {
		// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		// 	err = errors.New("Order already mapped to diferent user")
		// 	return
		// }
		// if status == "DELIVERY_ASSIGNED" {
		// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		// 	err = errors.New("Order already mapped to diferent user")
		// 	return
		// }
		id, cid, oid, odate, err := services.GetOderid(h)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Invalid Order ID ")
			msg11 = "Invalid Order ID "
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Id - ", id)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Customer Id - ", cid)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Order ID - ", oid)

		fname, add1, add2, city, zip, email, mob, latlong, err := services.GetCustomerinfo(cid)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Invalid Order ID ")
			msg11 = "Invalid Order ID "
			return
		}

		cdeatils := CustomerDeatils{
			Name:     fname,
			Address1: add1,
			Address2: add2,
			City:     city,
			Zipcode:  zip,
			Mobile:   mob,
			Email:    email,
			Latlong:  latlong,
			//LocationURL: locationurl,
		}

		var jsonData []byte
		jsonData, err = json.Marshal(cdeatils)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "jsonData - ", string(jsonData))

		oid, odate, otype, orepeat, dfrom, dto, dtime, cname, subcname, pname, pqty, productpackageqty, pprice, punits, tprice, totalsellprice, clustername, productsku, err := services.GetOderinfo(oid)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Invalid Order ID ")
			msg11 = "Invalid Order ID "
			return
		}

		odeatils := Oderinfo{
			Orderid:           oid,
			OderDate:          odate,
			OrderType:         otype,
			RepeatOrder:       orepeat,
			DeliveryFrom:      dfrom,
			DeliveryTo:        dto,
			DeliveryTime:      dtime,
			CategoryName:      cname,
			SubcategroryName:  subcname,
			ProductName:       pname,
			ProductQty:        pqty,
			ProductSKU:        productsku,
			ProductPackageQty: productpackageqty,
			ProductPrice:      pprice,
			ProductUnits:      punits,
			TotalPrice:        tprice,
			TotalSellPrice:    totalsellprice,
			ClusterName:       clustername,
			Name:              fname,
			Address1:          add1,
			Address2:          add2,
			City:              city,
			Zipcode:           zip,
			Mobile:            mob,
			Email:             email,
		}

		var jsonData1 []byte
		jsonData1, err = json.Marshal(odeatils)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "jsonData1 - ", string(jsonData1))

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "dfrom - ", dfrom)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "dto - ", dto)

		// start := time.Now()
		// end := start.AddDate(0, 0, 6)
		//fmt.Println(dfrom.Format("2006-01-02"), "-", dto.Format("2006-01-02"))

		for rd := rangeDate(date(dfrom), date(dto)); ; {
			date := rd()
			if date.IsZero() {
				break
			}
			fmt.Println(date.Format("2006-01-02"))
		}
		// s := Order{
		// 	OrderID: oid,
		// 	Msg: DeliveryOderinfo{
		// 		SubOrderid:        oid,
		// 		DeliveryDate:      date.Format("2006-01-02"),
		// 		DeliveryTime:      dtime,
		// 		ProductName:       pname,
		// 		ProductQty:        pqty,
		// 		ProductSKU:        productsku,
		// 		ProductPackageQty: productpackageqty,
		// 		ProductPrice:      pprice,
		// 		ProductUnits:      punits,
		// 		ClusterName:       clustername,
		// 		Name:              fname,
		// 		Address1:          add1,
		// 		Address2:          add2,
		// 		City:              city,
		// 		Zipcode:           zip,
		// 		Mobile:            mob,
		// 		Email:             email,
		// 	},
		// }
		// // var jsonArr []interface{}

		// jsonArr = append(jsonArr, s)
		// var Oderinfo map[string]interface{}

		// err := json.Unmarshal([]byte("{}"), &obj)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
		// obj["key1"] = "val1"
		// obj["key2"] = "val2"
		// jsonArr = append(jsonArr, s)

		// deliverydetails :=
		// var jsonData3 []byte
		// jsonData3, err = json.Marshal(s)
		// if err != nil {
		// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		// }

		// log.Println(beego.AppConfig.String("loglevel"), "Debug", "jsonData3 - ", string(jsonData3))
		//}
		// }

		// totalday := (daysBetween(date(dfrom), date(dto)))
		// totaldays := strconv.Itoa(totalday)
		//daysBetween(date("2021-08-30"), date("2021-12-30"))
		//log.Println(beego.AppConfig.String("loglevel"), "Debug", "totaldays - ", totaldays)

		// td, err := strconv.ParseFloat(totaldays, 64)
		// TDDD := td + 1
		// log.Println(beego.AppConfig.String("loglevel"), "Debug", "TDDD - ", TDDD)

		// log.Println(beego.AppConfig.String("loglevel"), "Debug", "productsku - ", productsku)

		aunits, err := services.Getinventory2(productsku, clustername)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Inventory Get Failed")
			msg11 = "Unable to get Inventory Details"
			return
		}

		// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Available Units - ", aunits)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Product SKU - ", productsku)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Cluster Name - ", clustername)

		punits2, err := strconv.ParseFloat(punits, 64)
		aunits1, err := strconv.ParseFloat(aunits, 64)

		// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Ordered Units - ", punits2)
		// log.Println(beego.AppConfig.String("loglevel"), "Debug", "aunits1 - ", aunits1)

		PT := punits2

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "OrderedUnits - ", PT)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "AvailableUnits - ", aunits1)

		// if PT > aunits1 {

		// 	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "inside - ")

		// 	err = errors.New("Insufficent Available Inventory")
		// 	msg11 = "Insufficent Available Inventory" + order_id
		// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		// 	return
		// }

		//	order_id =

		if PT > aunits1 {

			//log.Println(beego.AppConfig.String("loglevel"), "Debug", "inside - ")

			err = errors.New("Insufficent Available Inventory")
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			msg11 = "Insufficent Available Inventory For the Order : " + oid + " Check The Inventory Before Assignments"

			return
		}

		//Checkinventory(aunits1, PT)

		//number, err := strconv.ParseUint(totaldays, 10, 32)
		//finalIntNum := int(number) //Convert uint64 To int
		//td, err := strconv.ParseInt(totaldays, 0, 64)
		//tpp := strconv.FormatInt(totaldays, 10)

		// for i := 0; i < finalIntNum; i++ {

		// for rd := rangeDate(date(dfrom), date(dto)); ; {
		// 	date := rd()
		// 	if date.IsZero() {
		// 		break
		// 	}
		// 	fmt.Println(date.Format("2006-01-02"))

		// 	fmt.Printf("GeeksforGeeks\n")

		// 	log.Println(beego.AppConfig.String("loglevel"), "Debug", "array -")

		deliveryuser_id, deliveryuser_cluster, _ := services.GetDeliveryuserid2(select_deliveruser_name)
		if err != nil {
			beego.Error(err)
			err = errors.New("Invalid CollectionAgent  ")
			return
		}

		result, err := db.Db.Exec(`INSERT INTO Oder2Delivery (
		oder_id,
		order_date,
		deliveryuser_mob,
		customer_details,
		oder_deatils,
		created_by,
		status,
		delivery_dates,
		deliveryuser_id,
		deliveryuser_cluster,
		create_date)
		VALUES (?,?,?,?,?,?,?,?,?,?,now())`,
			oid,
			odate,
			select_deliveruser_name,
			jsonData,
			jsonData1,
			uid,
			"DELIVERY_ASSIGNED",
			order_deliverydate,
			deliveryuser_id,
			deliveryuser_cluster)
		if err != nil {
			err = errors.New("Delivery Order  creation failed")
			msg11 = "Delivery Order  creation failed"
			return
		}

		i, err := result.RowsAffected()
		if err != nil || i == 0 {
			err = errors.New("Delivery Order creation failed")
			msg11 = "Delivery Order  creation failed"
			return
		}

		// }
		// for _ = range totaldays {

		// }

		// result, err := db.Db.Exec(`INSERT INTO Oder2Delivery (oder_id,
		// order_date,
		// deliveryuser_mob,
		// customer_details,
		// oder_deatils,
		// created_by,
		// status,
		// create_date)
		// VALUES (?,?,?,?,?,?,?,now())`,
		// 	oid,
		// 	odate,
		// 	select_deliveruser_name,
		// 	jsonData,
		// 	jsonData1,
		// 	uid,
		// 	"DELIVERY_ASSIGNED")
		// if err != nil {
		// 	err = errors.New("Delivery Order  creation failed")
		// 	return
		// }

		// i, err := result.RowsAffected()
		// if err != nil || i == 0 {
		// 	err = errors.New("Delivery Order creation failed")
		// 	return
		// }

		if dto == order_deliverydate {

			res, err := db.Db.Exec(`UPDATE CustomerOrderTxn SET status = "DELIVERY_ASSIGNED"
	WHERE order_id = ?`, h)
			if err != nil {
				err = errors.New("CustomerOrderTxn Delivery updation failed")
				msg11 = "CustomerOrderTxn Delivery updation failed"
				return
			}

			i1, err := res.RowsAffected()
			if err != nil || i1 == 0 {
				err = errors.New(" CustomerOrderTxn Delivery updation failed")
				msg11 = "CustomerOrderTxn Delivery updation failed"
				return
			}

		}

		// log.Println(beego.AppConfig.String("loglevel"), "Debug", "punits - ", punits)
		// log.Println(beego.AppConfig.String("loglevel"), "Debug", "productsku - ", productsku)
		// log.Println(beego.AppConfig.String("loglevel"), "Debug", "clustername - ", clustername)

		dunits, aunits, runits, sku, cluster, damageunits, err := services.Getinventory(productsku, clustername)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Inventory Get Failed")
			msg11 = "Unable to get Inventory Details"
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Received - ", runits)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Delivered  - ", dunits)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Available  - ", aunits)

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "sku - ", sku)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "cluster - ", cluster)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "damageunits - ", damageunits)

		ddunits, err := strconv.ParseFloat(dunits, 64)
		rrunits, err := strconv.ParseFloat(runits, 64)
		ppunits, err := strconv.ParseFloat(punits, 64)
		ppunits1 := ppunits
		TD := ddunits + ppunits1
		TotalDel := strconv.FormatFloat(TD, 'f', 0, 32)
		TA := rrunits - TD
		TotalAvl := strconv.FormatFloat(TA, 'f', 0, 32)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Updated Delivered - ", TotalDel)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Updated Available - ", TotalAvl)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "productskuuuu - ", productsku)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "clusterrrrr - ", cluster)

		res2, err := db.Db.Exec(`UPDATE Inventory SET inventory_delivered = ? ,inventory_available=?
		              WHERE product_sku = ? AND cluster_name=?`, TotalDel, TotalAvl, productsku, cluster)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Inventory Get Failed")
			msg11 = "Update Inventory Details Failed"
			return
		}

		i2, err := res2.RowsAffected()
		if err != nil || i2 == 0 {
			err = errors.New(" CustomerOrderTxn Delivery updation failed")
			return
		}
		msg11 = "Order Assigned successfully"
	}

	return
}
func CheckMapExists(cluster, deliveryusers string) (err error) {

	err = nil
	row, err := db.Db.Query(`select count(*) from Cluster2DeliveryUser WHERE cluster_id = (select id from Clusters where cluster_name=?) AND deliveryuser_id=(select id from DeliveryUser where mob_no=?)`, cluster, deliveryusers)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch user")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch user")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	countlen := data[0][0]

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "countlen", countlen)

	if countlen == "0" {
		err = errors.New("Mappig already exists")
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}

	return

}

func CheckmapAlreadyExists(duser, oid string) (err error) {

	err = nil

	row, err := db.Db.Query(`select count(*) from Oder2Delivery WHERE oder_id = ? AND deliveryuser_mob=? AND status= "DELIVERY_ASSIGNED"`, oid, duser)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch user")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch user")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	countlen := data[0][0]

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "countlen", countlen)

	if countlen != "0" {
		err = errors.New("Mappig already exists")
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}

	return

}

func (c *GetDeliveryDateOrders) Get() {
	// msg11 := ""
	// log.Println(beego.AppConfig.String("loglevel"), "Info", "node post page")
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
			c.TplName = "deliveryorder/selectDeliveryOder/selectDeliveryOder.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SelectDeloveryOrder Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "deliveryorder/selectDeliveryOder/selectDeliveryOder.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SelectDeloveryOrder  Page Success")
		}
		return
	}()
	utils.SetHTTPHeader(c.Ctx)
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
	defer sess.SessionRelease(c.Ctx.ResponseWriter)

	utype := sess.Get("user_type").(string)
	c.Data["utype"] = utype

	language := sess.Get("language").(string)
	c.Data["language"] = language
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "language", language)

	id := sess.Get("uid").(string)
	c.Data["id"] = id

	clustername := sess.Get("clustername").(string)
	clustername1 := strings.ToUpper(clustername)
	c.Data["clustername"] = clustername1
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "clustername", clustername1)

	uid := sess.Get("uid").(string)
	c.Data["id"] = uid
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "id", uid)

	auth, err := utils.IsAuthorized(utype, "DeliveryOder")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	//------------this delivery date getting from notifications--------------

	delivery_date := c.GetString("delivery_date")
	order_id := c.GetString("order_id")

	if strings.Contains(delivery_date, " ") {
		delivery_date = strings.Split(delivery_date, " ")[0]
	}

	c.Data["delivery_date"] = delivery_date

	fmt.Println("Selected delivery_date is :- ", delivery_date)
	fmt.Println("Selected order_id is :- ", order_id)

	c.Data["datebutton"] = "HIDE"
	c.Data["duserbutton"] = "SHOW"

	if utype == "CLUSTERMANAGER" {

		data1, err1 := services.GetClusterAsPerMapping(id)
		if err1 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
			err1 = errors.New("Cluster user fetch Failed")
			return
		}
		if len(data1) <= 0 {
			err = errors.New("Cluster user  Not Found")
			return
		}

		data, err := services.GetClustermappedDeliveryuser(data1[0][1])
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Delivery user fetch Failed")
			return
		}
		if len(data) <= 0 {
			err = errors.New("Delivery user  Not Found")
			return
		}
		var Dis Display
		for i := range data {
			var d Field
			d.Id = data[i][0]
			d.Name = data[i][1]
			d.Email = data[i][2]
			d.Cluster = data[i][7]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

		auth, err := utils.IsAuthorized(utype, "DeliveryOder")
		if !auth {
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
			Autherr = errors.New("UnAuthorized")
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

		row, err := db.Db.Query(` select id,
 order_date, 
	order_id,
   cluster_name,
product_name,
status,
customer_id,
order_type,
orderdeliver_fromdate,
orderdeliver_todate,
order_units
 from CustomerOrderTxn where status= "ORDER_RECEIVED" AND order_purchagetype="KRISHIVALA" AND cluster_type="SOURCE" AND cluster_name=? AND order_id=?`, data1[0][3], order_id)

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch order")
			return
		}
		defer sql.Close(row)
		_, data4, err := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch order")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", data4)
		var result []Row

		for i := range data4 {

			// For Loop to get All Orders with condition ORDER_RECEIVED and KRISHIVALA and SOURCE and not bulk orders
			var r Row

			for rd := rangeDate(date(data4[i][8]), date(data4[i][9])); ; {

				// For Loop to get every date between from date and to date

				date := rd()
				if date.IsZero() {
					break
				}

				if delivery_date == date.Format("2006-01-02") {

					// Condition to check if Selected ( Delivery Date ) and  (Date) from the loop  between From and To date-range matches

					row, err := db.Db.Query(`select count(*) from Oder2Delivery WHERE oder_id = ? AND delivery_dates=?`, data4[i][2], date.Format("2006-01-02"))

					if err != nil {
						log.Println(beego.AppConfig.String("loglevel"), "Error", err)
						err = errors.New("Unable to fetch delivery date - order info")
						return
					}
					defer sql.Close(row)
					_, data, err := sql.Scan(row)
					if err != nil {
						log.Println(beego.AppConfig.String("loglevel"), "Error", err)
						err = errors.New("Unable to fetch  delivery date - order info")
						return
					}

					// log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

					countlen := data[0][0]

					// log.Println(beego.AppConfig.String("loglevel"), "Debug", "countlen", countlen)

					// Condition to check if (Date) from the loop  between From and To dates range is already assigned for delivery

					if countlen != "0" {

						fmt.Println("Already Delivery is Assigned for Order" + data4[i][2] + " on selected Delivery Date" + date.Format("2006-01-02"))

					} else {

						fmt.Println("Id :-" + data4[i][0])
						fmt.Println("OrderId :- " + data4[i][2])
						fmt.Println("Delivery Date :-" + date.Format("2006-01-02"))
						r.Id = data4[i][0]
						r.TimeStamp = data4[i][1]
						r.Oder_id = data4[i][2]
						r.Cluster_name = data4[i][3]
						r.Product_name = data4[i][4]
						r.Status = data4[i][5]
						r.CustomerID = data4[i][6]
						r.OrderType = data4[i][7]
						r.DeliveryDate = date.Format("2006-01-02")
						r.Units = data4[i][10]

						// log.Println(beego.AppConfig.String("loglevel"), "Debug", "data4[i][2]", data4[i][2])

						fname, add1, add2, city, zipcode, email, mob, latlong, err := services.GetCustomerinfo(data4[i][6])
						if err != nil {
							log.Println(beego.AppConfig.String("loglevel"), "Error", err)
							err = errors.New("Invalid Order ID ")
							return
						}

						var tmp = fname + add1 + add2

						// log.Println(beego.AppConfig.String("loglevel"), "Debug", "customer info1", tmp)

						var tmp2 = city + zipcode + email
						var tmp3 = mob
						fmt.Println(latlong)
						// log.Println(beego.AppConfig.String("loglevel"), "Debug", "customer info2", tmp2)

						// log.Println(beego.AppConfig.String("loglevel"), "Debug", "customer info3", tmp3)

						r.Custdetails = tmp + "," + tmp2 + "," + tmp3
						//	r.LocationURL = locationurl
						result = append(result, r)

					}

				}

			}

		}
		c.Data["CustomerData"] = result

	} else if utype == "CLUSTEREXECUTIVE" {

		// CLUSTER-EXECUTIVE CODE

		data1, err1 := services.GetClusterAsPerMappingExec(id)
		if err1 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
			err1 = errors.New("Cluster user fetch Failed")
			return
		}
		if len(data1) <= 0 {
			err = errors.New("Cluster user  Not Found")
			return
		}

		data, err := services.GetClustermappedDeliveryuser(data1[0][1])
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Delivery user fetch Failed")
			return
		}
		if len(data) <= 0 {
			err = errors.New("Delivery user  Not Found")
			return
		}
		var Dis Display
		for i := range data {
			var d Field
			d.Id = data[i][0]
			d.Name = data[i][1]
			d.Email = data[i][2]
			d.Cluster = data[i][7]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

		auth, err := utils.IsAuthorized(utype, "DeliveryOder")
		if !auth {
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
			Autherr = errors.New("UnAuthorized")
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

		row, err := db.Db.Query(` select id,
 order_date, 
	order_id,
   cluster_name,
product_name,
status,
customer_id,
order_type,
orderdeliver_fromdate,
orderdeliver_todate,
order_units
 from CustomerOrderTxn where status= "ORDER_RECEIVED" AND order_purchagetype="KRISHIVALA" AND cluster_type="SOURCE" AND cluster_name=? AND order_id=?`, data1[0][3], order_id)

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch order")
			return
		}
		defer sql.Close(row)
		_, data4, err := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch order")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", data4)
		var result []Row

		for i := range data4 {

			// For Loop to get All Orders with condition ORDER_RECEIVED,KRISHIVALA,SOURCE,not bulk order
			var r Row

			for rd := rangeDate(date(data4[i][8]), date(data4[i][9])); ; {

				// For Loop to get every date between from date and to date

				date := rd()
				if date.IsZero() {
					break
				}

				if delivery_date == date.Format("2006-01-02") {

					// Condition to check if Selected ( Delivery Date ) and  (Date) from the loop  between From and To dates range matches

					row, err := db.Db.Query(`select count(*) from Oder2Delivery WHERE oder_id = ? AND delivery_dates=?`, data4[i][2], date.Format("2006-01-02"))

					if err != nil {
						log.Println(beego.AppConfig.String("loglevel"), "Error", err)
						err = errors.New("Unable to fetch delivery date - order info")
						return
					}
					defer sql.Close(row)
					_, data, err := sql.Scan(row)
					if err != nil {
						log.Println(beego.AppConfig.String("loglevel"), "Error", err)
						err = errors.New("Unable to fetch  delivery date - order info")
						return
					}

					// log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

					countlen := data[0][0]

					// log.Println(beego.AppConfig.String("loglevel"), "Debug", "countlen", countlen)

					if countlen != "0" {

						// Condition to check if (Date) from the loop  between From and To dates range is already assigned for delivery

						fmt.Println("Already Delivery is Assigned for Order" + data4[i][2] + " on selected Delivery Date" + date.Format("2006-01-02"))

					} else {

						fmt.Println("Id :-" + data4[i][0])
						fmt.Println("OrderId :- " + data4[i][2])
						fmt.Println("Delivery Date :-" + date.Format("2006-01-02"))
						r.Id = data4[i][0]
						r.TimeStamp = data4[i][1]
						r.Oder_id = data4[i][2]
						r.Cluster_name = data4[i][3]
						r.Product_name = data4[i][4]
						r.Status = data4[i][5]
						r.CustomerID = data4[i][6]
						r.OrderType = data4[i][7]
						r.DeliveryDate = date.Format("2006-01-02")
						r.Units = data4[i][10]

						// log.Println(beego.AppConfig.String("loglevel"), "Debug", "data4[i][2]", data4[i][2])

						fname, add1, add2, city, zipcode, email, mob, latlong, err := services.GetCustomerinfo(data4[i][6])
						if err != nil {
							log.Println(beego.AppConfig.String("loglevel"), "Error", err)
							err = errors.New("Invalid Order ID ")
							return
						}

						var tmp = fname + add1 + add2

						// log.Println(beego.AppConfig.String("loglevel"), "Debug", "customer info1", tmp)

						var tmp2 = city + zipcode + email
						var tmp3 = mob
						fmt.Println(latlong)
						// log.Println(beego.AppConfig.String("loglevel"), "Debug", "customer info2", tmp2)

						// log.Println(beego.AppConfig.String("loglevel"), "Debug", "customer info3", tmp3)

						r.Custdetails = tmp + "," + tmp2 + "," + tmp3
						//	r.LocationURL = locationurl
						result = append(result, r)

					}

				}

			}

		}
		c.Data["CustomerData"] = result

	} else {

		// HOADMIN CODE

		data, err := services.GetDeliveryuser()
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Delivery user fetch Failed")
			return
		}
		if len(data) <= 0 {
			err = errors.New("Delivery user  Not Found")
			return
		}
		var Dis Display
		for i := range data {
			var d Field
			d.Id = data[i][0]
			d.Name = data[i][1]
			d.Email = data[i][2]
			d.Cluster = data[i][7]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

		auth, err := utils.IsAuthorized(utype, "DeliveryOder")
		if !auth {
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
			Autherr = errors.New("UnAuthorized")
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

		row, err := db.Db.Query(` select id,
 order_date, 
	order_id,
   cluster_name,
product_name,
status,
customer_id,
order_type,
orderdeliver_fromdate,
orderdeliver_todate,
order_units
 from CustomerOrderTxn where status= "ORDER_RECEIVED" AND  order_purchagetype="KRISHIVALA" AND cluster_type="SOURCE" AND order_id=?`, order_id)

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch order")
			return
		}
		defer sql.Close(row)
		_, data4, err := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch order")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", data4)
		var result []Row

		for i := range data4 {

			// For Loop to get All Orders with condition ORDER_RECEIVED,KRISHIVALA,SOURCE,not bulk order
			var r Row

			for rd := rangeDate(date(data4[i][8]), date(data4[i][9])); ; {

				// For Loop to get every date between from date and to date

				date := rd()
				if date.IsZero() {
					break
				}

				if delivery_date == date.Format("2006-01-02") {

					// Condition to check if Selected ( Delivery Date ) and  (Date) from the loop  between From and To dates range matches

					row, err := db.Db.Query(`select count(*) from Oder2Delivery WHERE oder_id = ? AND delivery_dates=?`, data4[i][2], date.Format("2006-01-02"))

					if err != nil {
						log.Println(beego.AppConfig.String("loglevel"), "Error", err)
						err = errors.New("Unable to fetch delivery date - order info")
						return
					}
					defer sql.Close(row)
					_, data, err := sql.Scan(row)
					if err != nil {
						log.Println(beego.AppConfig.String("loglevel"), "Error", err)
						err = errors.New("Unable to fetch  delivery date - order info")
						return
					}

					// log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

					countlen := data[0][0]

					// log.Println(beego.AppConfig.String("loglevel"), "Debug", "countlen", countlen)

					if countlen != "0" {

						// Condition to check if (Date) from the loop  between From and To dates range is already assigned for delivery

						fmt.Println("Already Delivery is Assigned for Order" + data4[i][2] + " on selected Delivery Date" + date.Format("2006-01-02"))

					} else {

						fmt.Println("Id :-" + data4[i][0])
						fmt.Println("OrderId :- " + data4[i][2])
						fmt.Println("Delivery Date :-" + date.Format("2006-01-02"))
						r.Id = data4[i][0]
						r.TimeStamp = data4[i][1]
						r.Oder_id = data4[i][2]
						r.Cluster_name = data4[i][3]
						r.Product_name = data4[i][4]
						r.Status = data4[i][5]
						r.CustomerID = data4[i][6]
						r.OrderType = data4[i][7]
						r.DeliveryDate = date.Format("2006-01-02")
						r.Units = data4[i][10]

						// log.Println(beego.AppConfig.String("loglevel"), "Debug", "data4[i][2]", data4[i][2])

						fname, add1, add2, city, zipcode, email, mob, latlong, err := services.GetCustomerinfo(data4[i][6])
						if err != nil {
							log.Println(beego.AppConfig.String("loglevel"), "Error", err)
							err = errors.New("Invalid Order ID ")
							return
						}

						var tmp = fname + add1 + add2

						// log.Println(beego.AppConfig.String("loglevel"), "Debug", "customer info1", tmp)

						var tmp2 = city + zipcode + email
						var tmp3 = mob
						fmt.Println(latlong)
						// log.Println(beego.AppConfig.String("loglevel"), "Debug", "customer info2", tmp2)

						// log.Println(beego.AppConfig.String("loglevel"), "Debug", "customer info3", tmp3)

						r.Custdetails = tmp + "," + tmp2 + "," + tmp3
						//	r.LocationURL = locationurl
						result = append(result, r)

					}

				}

			}

		}
		c.Data["CustomerData"] = result

	}

	return
}

// func Checkorderstatus(duser, oid string) (status string, err error) {

// 	row, err := db.Db.Query(`select status from Oder2Delivery WHERE oder_id = ? AND deliveryuser_mob=? `, oid, duser)
// 	if err != nil {
// 		beego.Error(err)
// 		err = errors.New("Unable to get the CustomerOrderTxn Id  ")
// 		return
// 	}
// 	defer sql.Close(row)
// 	_, data, err := sql.Scan(row)
// 	if err != nil {
// 		beego.Error(err)
// 		err = errors.New("Message Scan Fail")
// 		return
// 	}
// 	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
// 	if len(data) <= 0 {
// 		err = errors.New("Unable to get the Agent Account number info")
// 		return
// 	}
// 	status = data[0][0]

//		return
//	}
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

func rangeDate(start, end time.Time) func() time.Time {
	y, m, d := start.Date()
	start = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	y, m, d = end.Date()
	end = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	return func() time.Time {
		if start.After(end) {
			return time.Time{}
		}
		date := start
		start = start.AddDate(0, 0, 1)
		return date
	}
}
