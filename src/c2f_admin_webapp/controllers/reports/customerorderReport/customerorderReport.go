package customerorderReport

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
	TimeStamp             string
	Oder_id               string
	Customer_mob          string
	Cluster_name          string
	Product_name          string
	Product_Qty           string
	Product_Package_Qty   string
	Cluster_Selling_price string
	Product_Units         string
	Product_Total_Price   string
	Status                string
	FromDate              string
	Todate                string
	Time                  string
	Repeat                string
	TotalSellPrice        string
	Payment_Status        string
	Payment_Type          string
	Order_Type            string
	Purchase_Type         string
	ClusterType           string
}
type Display5 struct {
	Fields5 []Field5
}

type Field5 struct {
	Id     string
	Mobile string
}
type CustomerorderReport struct {
	// beego.Controller
	middleware.BaseController
}

func (c *CustomerorderReport) Get() {

	Status := c.Ctx.Input.Param(":Status")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Status", Status)
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
			c.TplName = "reports/customerorderReport/customerorderReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Reports Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "reports/customerorderReport/customerorderReport.html"
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
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "utype", utype)

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

	orderId := c.GetString("order_id")

	data, err := services.GetConsumer()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Producer fetch Failed")
		return
	}
	if len(data) <= 0 {
		err = errors.New("Producer  Not Found")
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

	auth, err := utils.IsAuthorized(utype, "CustOrderReport")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	if utype == "CLUSTERMANAGER" {

		uid := sess.Get("uid").(string)
		c.Data["utype"] = utype
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "uid", uid)

		cid, err := services.GetCid(uid)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster ID fetch Failed")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "cid", cid)

		cname, err := services.GetClustername(cid)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster Name fetch Failed")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "cname", cname)

		if Status == "" {

			row, err := db.Db.Query(`select order_date, 
	order_id,
   customer_mob,
cluster_name,
product_name,
product_qty,
product_package_qty,
product_price,
order_units,
total_product_price,
total_selling_price,
status,
orderdeliver_fromdate,
orderdeliver_todate,
orderdeliver_time,
order_repeat,
payment_status,
payment_type,
order_type,
order_purchagetype,
cluster_type

 from CustomerOrderTxn where (cluster_name='' OR cluster_name like ?) OR order_id =? ORDER BY order_date DESC`, cname, orderId)

			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to fetch order")
				return
			}
			defer sql.Close(row)
			_, data, err1 := sql.Scan(row)
			if err1 != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
				err1 = errors.New("Unable to fetch order")
				return
			}

			log.Println(beego.AppConfig.String("loglevel"), "Debug", data)
			var result []Row

			for i := range data {
				var r Row
				r.TimeStamp = data[i][0]
				r.Oder_id = data[i][1]
				r.Customer_mob = data[i][2]
				r.Cluster_name = data[i][3]
				r.Product_name = data[i][4]
				r.Product_Qty = data[i][5]
				r.Product_Package_Qty = data[i][6]
				r.Cluster_Selling_price = data[i][7]
				r.Product_Units = data[i][8]
				r.Product_Total_Price = data[i][9]
				r.TotalSellPrice = data[i][10]
				r.Status = data[i][11]
				r.FromDate = data[i][12]
				r.Todate = data[i][13]
				r.Time = data[i][14]
				r.Repeat = data[i][15]
				r.Payment_Status = data[i][16]
				r.Payment_Type = data[i][17]
				r.Order_Type = data[i][18]
				r.Purchase_Type = data[i][19]
				r.ClusterType = data[i][20]

				result = append(result, r)
			}
			c.Data["CustomerData"] = result

		} else {

			row, err := db.Db.Query(`select order_date, 
	order_id,
   customer_mob,
cluster_name,
product_name,
product_qty,
product_package_qty,
product_price,
order_units,
total_product_price,
total_selling_price,
status,
orderdeliver_fromdate,
orderdeliver_todate,
orderdeliver_time,
order_repeat,
payment_status,
payment_type,
order_type,
order_purchagetype,
cluster_type

 from CustomerOrderTxn where (cluster_name='' OR cluster_name like ?) and (status='' OR status like ?) OR order_id =? ORDER BY order_date DESC`, cname, Status, orderId)

			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to fetch order")
				return
			}
			defer sql.Close(row)
			_, data, err1 := sql.Scan(row)
			if err1 != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
				err1 = errors.New("Unable to fetch order")
				return
			}

			log.Println(beego.AppConfig.String("loglevel"), "Debug", data)
			var result []Row

			for i := range data {
				var r Row
				r.TimeStamp = data[i][0]
				r.Oder_id = data[i][1]
				r.Customer_mob = data[i][2]
				r.Cluster_name = data[i][3]
				r.Product_name = data[i][4]
				r.Product_Qty = data[i][5]
				r.Product_Package_Qty = data[i][6]
				r.Cluster_Selling_price = data[i][7]
				r.Product_Units = data[i][8]
				r.Product_Total_Price = data[i][9]
				r.TotalSellPrice = data[i][10]
				r.Status = data[i][11]
				r.FromDate = data[i][12]
				r.Todate = data[i][13]
				r.Time = data[i][14]
				r.Repeat = data[i][15]
				r.Payment_Status = data[i][16]
				r.Payment_Type = data[i][17]
				r.Order_Type = data[i][18]
				r.Purchase_Type = data[i][19]
				r.ClusterType = data[i][20]

				result = append(result, r)
			}
			c.Data["CustomerData"] = result

		}

	} else if utype == "CLUSTEREXECUTIVE" {

		uid := sess.Get("uid").(string)
		c.Data["utype"] = utype
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "uid", uid)

		cid, err := services.GetCid1(uid)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster ID fetch Failed")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "cid", cid)

		cname, err := services.GetClustername(cid)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster Name fetch Failed")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "cname", cname)

		if Status == "" {

			row, err := db.Db.Query(`select order_date, 
	order_id,
   customer_mob,
cluster_name,
product_name,
product_qty,
product_package_qty,
product_price,
order_units,
total_product_price,
total_selling_price,
status,
orderdeliver_fromdate,
orderdeliver_todate,
orderdeliver_time,
order_repeat,
payment_status,
payment_type,
order_type,
order_purchagetype,
cluster_type

 from CustomerOrderTxn where (cluster_name='' OR cluster_name like ?) OR order_id =? ORDER BY order_date DESC`, cname, orderId)

			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to fetch order")
				return
			}
			defer sql.Close(row)
			_, data, err1 := sql.Scan(row)
			if err1 != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
				err1 = errors.New("Unable to fetch order")
				return
			}

			log.Println(beego.AppConfig.String("loglevel"), "Debug", data)
			var result []Row

			for i := range data {
				var r Row
				r.TimeStamp = data[i][0]
				r.Oder_id = data[i][1]
				r.Customer_mob = data[i][2]
				r.Cluster_name = data[i][3]
				r.Product_name = data[i][4]
				r.Product_Qty = data[i][5]
				r.Product_Package_Qty = data[i][6]
				r.Cluster_Selling_price = data[i][7]
				r.Product_Units = data[i][8]
				r.Product_Total_Price = data[i][9]
				r.TotalSellPrice = data[i][10]
				r.Status = data[i][11]
				r.FromDate = data[i][12]
				r.Todate = data[i][13]
				r.Time = data[i][14]
				r.Repeat = data[i][15]
				r.Payment_Status = data[i][16]
				r.Payment_Type = data[i][17]
				r.Order_Type = data[i][18]
				r.Purchase_Type = data[i][19]
				r.ClusterType = data[i][20]

				result = append(result, r)
			}
			c.Data["CustomerData"] = result

		} else {

			row, err := db.Db.Query(`select order_date, 
	order_id,
   customer_mob,
cluster_name,
product_name,
product_qty,
product_package_qty,
product_price,
order_units,
total_product_price,
total_selling_price,
status,
orderdeliver_fromdate,
orderdeliver_todate,
orderdeliver_time,
order_repeat,
payment_status,
payment_type,
order_type,
order_purchagetype,
cluster_type

 from CustomerOrderTxn where (cluster_name='' OR cluster_name like ?) and (status='' OR status like ?) OR order_id =? ORDER BY order_date DESC`, cname, Status, orderId)

			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to fetch order")
				return
			}
			defer sql.Close(row)
			_, data, err1 := sql.Scan(row)
			if err1 != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
				err1 = errors.New("Unable to fetch order")
				return
			}

			log.Println(beego.AppConfig.String("loglevel"), "Debug", data)
			var result []Row

			for i := range data {
				var r Row
				r.TimeStamp = data[i][0]
				r.Oder_id = data[i][1]
				r.Customer_mob = data[i][2]
				r.Cluster_name = data[i][3]
				r.Product_name = data[i][4]
				r.Product_Qty = data[i][5]
				r.Product_Package_Qty = data[i][6]
				r.Cluster_Selling_price = data[i][7]
				r.Product_Units = data[i][8]
				r.Product_Total_Price = data[i][9]
				r.TotalSellPrice = data[i][10]
				r.Status = data[i][11]
				r.FromDate = data[i][12]
				r.Todate = data[i][13]
				r.Time = data[i][14]
				r.Repeat = data[i][15]
				r.Payment_Status = data[i][16]
				r.Payment_Type = data[i][17]
				r.Order_Type = data[i][18]
				r.Purchase_Type = data[i][19]
				r.ClusterType = data[i][20]

				result = append(result, r)
			}
			c.Data["CustomerData"] = result

		}

	} else {

		// HOADMIN CODE

		row, err := db.Db.Query(`select order_date, 
	order_id,
   customer_mob,
cluster_name,
product_name,
product_qty,
product_package_qty,
product_price,
order_units,
total_product_price,
total_selling_price,
status,
orderdeliver_fromdate,
orderdeliver_todate,
orderdeliver_time,
order_repeat,
payment_status,
payment_type,
order_type,
order_purchagetype,
cluster_type

from CustomerOrderTxn where (? = '' OR order_id = ?) ORDER BY order_date DESC`, orderId, orderId)

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch order")
			return
		}
		defer sql.Close(row)
		_, data, err1 := sql.Scan(row)
		if err1 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
			err1 = errors.New("Unable to fetch order")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", data)
		var result []Row

		for i := range data {
			var r Row
			r.TimeStamp = data[i][0]
			r.Oder_id = data[i][1]
			r.Customer_mob = data[i][2]
			r.Cluster_name = data[i][3]
			r.Product_name = data[i][4]
			r.Product_Qty = data[i][5]
			r.Product_Package_Qty = data[i][6]
			r.Cluster_Selling_price = data[i][7]
			r.Product_Units = data[i][8]
			r.Product_Total_Price = data[i][9]
			r.TotalSellPrice = data[i][10]
			r.Status = data[i][11]
			r.FromDate = data[i][12]
			r.Todate = data[i][13]
			r.Time = data[i][14]
			r.Repeat = data[i][15]
			r.Payment_Status = data[i][16]
			r.Payment_Type = data[i][17]
			r.Order_Type = data[i][18]
			r.Purchase_Type = data[i][19]
			r.ClusterType = data[i][20]

			result = append(result, r)
		}
		c.Data["CustomerData"] = result
	}

	return
}

func (c *CustomerorderReport) Post() {
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
			c.TplName = "reports/customerorderReport/customerorderReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchOrder Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "reports/customerorderReport/customerorderReport.html"
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

	auth, err := utils.IsAuthorized(utype, "CustOrderReport")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	input_rnn_no := c.Input().Get("input_rnn_no")
	input_status := c.Input().Get("input_status")
	input_card_mask := c.Input().Get("input_card_mask")
	input_customer_mob := c.Input().Get("input_customer_mob")
	input_order_id := c.Input().Get("input_order_id")
	input_cluster_name := c.Input().Get("input_cluster_name")
	input_payment_status := c.Input().Get("input_payment_status")
	input_payment_type := c.Input().Get("input_payment_type")
	input_order_type := c.Input().Get("input_order_type")
	input_purchase_type := c.Input().Get("input_purchase_type")
	input_product_name := c.Input().Get("input_product_name")

	c.Data["RRN"] = input_rnn_no
	c.Data["Status"] = input_status
	c.Data["Mask"] = input_card_mask
	c.Data["DateRange"] = c.Input().Get("daterange")
	tmp := strings.Split(c.Input().Get("daterange"), " - ")

	from := ""
	to := ""
	if len(tmp) == 2 {
		from = tmp[0]
		to = tmp[1]
	}
	if utype == "CLUSTERMANAGER" {

		uid := sess.Get("uid").(string)
		c.Data["utype"] = utype
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "uid", uid)

		cid, err := services.GetCid(uid)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster ID fetch Failed")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "cid", cid)

		cname, err := services.GetClustername(cid)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster Name fetch Failed")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "cname", cname)
		row, err := db.Db.Query(`select order_date, 
	order_id,
   customer_mob,
cluster_name,
product_name,
product_qty,
product_package_qty,
product_price,
order_units,
total_product_price,
total_selling_price,
status,
orderdeliver_fromdate,
orderdeliver_todate,
orderdeliver_time,
order_repeat,
payment_status,
payment_type,
order_type,
order_purchagetype,
cluster_type

 from CustomerOrderTxn where (cluster_name='' OR cluster_name like ?) AND (customer_mob='' OR customer_mob like ?) AND (order_id='' OR order_id like ?) AND (cluster_name='' OR cluster_name like ?) AND (status='' OR status like ?) AND (payment_status='' OR payment_status like ?) AND (payment_type='' OR payment_type like ?) AND (product_name='' OR product_name like ?) AND (DATE(order_date='') OR DATE (order_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(order_date='') OR DATE (order_date) <= STR_TO_DATE(?, '%m/%d/%Y')) AND (order_purchagetype='' OR order_purchagetype like ?) ORDER BY order_date DESC`, cname, input_customer_mob+"%", input_order_id+"%", input_cluster_name+"%", input_status+"%", input_payment_status+"%", input_payment_type+"%", input_product_name+"%", from, to, input_purchase_type+"%")

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
			r.TimeStamp = data[i][0]
			r.Oder_id = data[i][1]
			r.Customer_mob = data[i][2]
			r.Cluster_name = data[i][3]
			r.Product_name = data[i][4]
			r.Product_Qty = data[i][5]
			r.Product_Package_Qty = data[i][6]
			r.Cluster_Selling_price = data[i][7]
			r.Product_Units = data[i][8]
			r.Product_Total_Price = data[i][9]
			r.TotalSellPrice = data[i][10]
			r.Status = data[i][11]
			r.FromDate = data[i][12]
			r.Todate = data[i][13]
			r.Time = data[i][14]
			r.Repeat = data[i][15]
			r.Payment_Status = data[i][16]
			r.Payment_Type = data[i][17]
			r.Order_Type = data[i][18]
			r.Purchase_Type = data[i][19]
			r.ClusterType = data[i][20]

			result = append(result, r)
		}
		c.Data["CustomerData"] = result
	} else if utype == "CLUSTEREXECUTIVE" {
		uid := sess.Get("uid").(string)
		c.Data["utype"] = utype
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "uid", uid)

		cid, err := services.GetCid1(uid)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster ID fetch Failed")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "cid", cid)

		cname, err := services.GetClustername(cid)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster Name fetch Failed")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "cname", cname)
		row, err := db.Db.Query(`select order_date, 
	order_id,
   customer_mob,
cluster_name,
product_name,
product_qty,
product_package_qty,
product_price,
order_units,
total_product_price,
total_selling_price,
status,
orderdeliver_fromdate,
orderdeliver_todate,
orderdeliver_time,
order_repeat,
payment_status,
payment_type,
order_type,
order_purchagetype,
cluster_type

 from CustomerOrderTxn where (cluster_name='' OR cluster_name like ?) AND (customer_mob='' OR customer_mob like ?) AND (order_id='' OR order_id like ?) AND (cluster_name='' OR cluster_name like ?) AND (status='' OR status like ?) AND (payment_status='' OR payment_status like ?) AND (payment_type='' OR payment_type like ?) AND (product_name='' OR product_name like ?) AND (DATE(order_date='') OR DATE (order_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(order_date='') OR DATE (order_date) <= STR_TO_DATE(?, '%m/%d/%Y')) AND (order_purchagetype='' OR order_purchagetype like ?) ORDER BY order_date DESC`, cname, input_customer_mob+"%", input_order_id+"%", input_cluster_name+"%", input_status+"%", input_payment_status+"%", input_payment_type+"%", input_product_name+"%", from, to, input_purchase_type+"%")

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
			r.TimeStamp = data[i][0]
			r.Oder_id = data[i][1]
			r.Customer_mob = data[i][2]
			r.Cluster_name = data[i][3]
			r.Product_name = data[i][4]
			r.Product_Qty = data[i][5]
			r.Product_Package_Qty = data[i][6]
			r.Cluster_Selling_price = data[i][7]
			r.Product_Units = data[i][8]
			r.Product_Total_Price = data[i][9]
			r.TotalSellPrice = data[i][10]
			r.Status = data[i][11]
			r.FromDate = data[i][12]
			r.Todate = data[i][13]
			r.Time = data[i][14]
			r.Repeat = data[i][15]
			r.Payment_Status = data[i][16]
			r.Payment_Type = data[i][17]
			r.Order_Type = data[i][18]
			r.Purchase_Type = data[i][19]
			r.ClusterType = data[i][20]

			result = append(result, r)
		}
		c.Data["CustomerData"] = result
	} else {
		row, err := db.Db.Query(`select order_date, 
	order_id,
   customer_mob,
cluster_name,
product_name,
product_qty,
product_package_qty,
product_price,
order_units,
total_product_price,
total_selling_price,
status,
orderdeliver_fromdate,
orderdeliver_todate,
orderdeliver_time,
order_repeat,
payment_status,
payment_type,
order_type,
order_purchagetype,
cluster_type

 from CustomerOrderTxn where (customer_mob='' OR customer_mob like ?) AND (order_id='' OR order_id like ?) AND (cluster_name='' OR cluster_name like ?) AND (status='' OR status like ?) AND (payment_status='' OR payment_status like ?) AND (payment_type='' OR payment_type like ?) AND (product_name='' OR product_name like ?) AND (order_type='' OR order_type like ?) AND (DATE(order_date='') OR DATE (order_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(order_date='') OR DATE (order_date) <= STR_TO_DATE(?, '%m/%d/%Y')) AND (order_purchagetype='' OR order_purchagetype like ?) ORDER BY order_date DESC`, input_customer_mob+"%", input_order_id+"%", input_cluster_name+"%", input_status+"%", input_payment_status+"%", input_payment_type+"%", input_product_name+"%", input_order_type+"%", from, to, input_purchase_type+"%")

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
			r.TimeStamp = data[i][0]
			r.Oder_id = data[i][1]
			r.Customer_mob = data[i][2]
			r.Cluster_name = data[i][3]
			r.Product_name = data[i][4]
			r.Product_Qty = data[i][5]
			r.Product_Package_Qty = data[i][6]
			r.Cluster_Selling_price = data[i][7]
			r.Product_Units = data[i][8]
			r.Product_Total_Price = data[i][9]
			r.TotalSellPrice = data[i][10]
			r.Status = data[i][11]
			r.FromDate = data[i][12]
			r.Todate = data[i][13]
			r.Time = data[i][14]
			r.Repeat = data[i][15]
			r.Payment_Status = data[i][16]
			r.Payment_Type = data[i][17]
			r.Order_Type = data[i][18]
			r.Purchase_Type = data[i][19]
			r.ClusterType = data[i][20]
			result = append(result, r)
		}
		c.Data["CustomerData"] = result
	}

	return
}
