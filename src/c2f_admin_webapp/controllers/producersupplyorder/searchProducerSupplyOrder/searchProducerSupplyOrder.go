package searchProducerSupplyOrder

import (
	"c2f_admin_webapp/services"
	"c2f_admin_webapp/session"
	"runtime/debug"

	"c2f_admin_webapp/model/utils"
	"errors"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/db"

	"strings"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

type Row struct {
	Id                string
	SupplyTxnID       string
	ProducerMob       string
	ProductName       string
	ProductSKU        string
	ProductQTY        string
	ProductPackageQTY string
	ClusterBuyPrice   string
	SupplyUnits       string
	TotalPrice        string
	Status            string
}

type Display struct {
	Fields []Field
}

type Field struct {
	Id     string
	Mobile string
	Name   string
}
type Display1 struct {
	Fields1 []Field1
}

type Field1 struct {
	Supplyordertxn_id string
	ProducerMobile    string
	Status            string
	ProductName       string
	ClusterName       string
	Productqty        string
	PPackageQty       string
}
type SearchProducerSupplyOrder struct {
	// beego.Controller
	middleware.BaseController
}

func (c *SearchProducerSupplyOrder) Get() {
	Supplyordertxn_id := c.Ctx.Input.Param(":Supplyordertxn_id")
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

			c.TplName = "producersupplyorder/searchProducerSupplyOrder/searchProducerSupplyOrder.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "producersupplyorder/searchProducerSupplyOrder/searchProducerSupplyOrder.html"
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

	id := sess.Get("uid").(string)
	c.Data["id"] = id

	auth, err := utils.IsAuthorized(utype, "ProducerSupplyOrder")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")
	//select_producer_name := c.Input().Get("select_producer_name")

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

		data, err := services.GetProducers2(data1[0][1])
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Producer fetch Failed")
			return
		}
		if len(data) <= 0 {
			err = errors.New("Producer  Not Found")
			return
		}
		var Dis Display
		for i := range data {
			var d Field
			d.Id = data[i][0]
			d.Mobile = data[i][1]
			d.Name = data[i][2]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

	} else if utype == "CLUSTEREXECUTIVE" {

		// NO producer MGMT

	} else {

		// HOADMIN CODE

		data, err := services.GetProducers()
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Producer fetch Failed")
			return
		}
		if len(data) <= 0 {
			err = errors.New("Producer  Not Found")
			return
		}
		var Dis Display
		for i := range data {
			var d Field
			d.Id = data[i][0]
			d.Mobile = data[i][1]
			d.Name = data[i][2]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

		//select_producer_name := c.Input().Get("select_producer_name")

		data2, err := services.GetProducers1()
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Producer fetch Failed")
			return
		}
		if len(data2) <= 0 {
			err = errors.New("Producer  Not Found")
			return
		}
		var Dis1 Display1
		for i := range data2 {
			var d Field1
			d.ProducerMobile = data2[i][0]
			d.Status = data2[i][1]
			d.ProductName = data2[i][2]
			d.ClusterName = data2[i][3]
			d.Productqty = data2[i][4]
			d.PPackageQty = data2[i][5] + data2[i][4]
			d.Supplyordertxn_id = data2[i][6]
			Dis1.Fields1 = append(Dis1.Fields1, d)
		}
		c.Data["Dis1"] = Dis1
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis1)

		if Supplyordertxn_id == "" {

		} else {

			row1, err := db.Db.Query(`select id,supplyordertxn_id,product_name,product_sku,product_qty,product_package_qty,cluster_buying_price,product_supply_units,product_supply_total_price,status,producer_mob from SupplyOrderTxn where supplyordertxn_id=? AND (status="DEMAND_ORDER_CREATED" OR status="DEMAND_ORDER_ACCEPTED")`, Supplyordertxn_id)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to get the error message info")
				return
			}
			defer sql.Close(row1)
			_, data1, err := sql.Scan(row1)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Message Scan Fail")
				return
			}
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data1, " Data len - ", len(data1))
			if len(data1) <= 0 {
				err = errors.New("No Data Found")
				return
			}

			var result []Row

			for i := range data1 {
				var r Row
				r.Id = data1[i][0]
				r.SupplyTxnID = data1[i][1]
				r.ProductName = data1[i][2]
				r.ProductSKU = data1[i][3]
				r.ProductQTY = data1[i][4]
				r.ProductPackageQTY = data1[i][5]
				r.ClusterBuyPrice = data1[i][6]
				r.SupplyUnits = data1[i][7]
				r.TotalPrice = data1[i][8]
				r.Status = data1[i][9]
				r.ProducerMob = data1[i][10]
				result = append(result, r)
			}

			c.Data["CustomerData"] = result
			//c.Data["ProducerMob"] = data1[0][10]
		}

		row8, err := db.Db.Query(`SELECT status,count(producer_mob)as PSum FROM SupplyOrderTxn where (status="DEMAND_ORDER_CREATED" OR status="DEMAND_ORDER_ACCEPTED")`)

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get dashboard data")
			return
		}
		defer sql.Close(row8)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row8)
		_, data8, err := sql.Scan(row8)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get dashboard data")
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data8, "\nData len - ", len(data8))
		if len(data8) <= 0 {
			err = errors.New("Unable to get dashboard data")
			return
		}

		//log.Println(beego.AppConfig.String("loglevel"), "Debug", "PMob - ", data8[0][0])
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "PStatus - ", data8[0][0])
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "PCount - ", data8[0][1])

		//c.Data["producermobile"] = data8[0][0]
		c.Data["status"] = data8[0][0]
		c.Data["totalproducercount"] = data8[0][1]

	}

	return
}

func (c *SearchProducerSupplyOrder) Post() {
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
			c.TplName = "producersupplyorder/searchProducerSupplyOrder/searchProducerSupplyOrder.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "producersupplyorder/searchProducerSupplyOrder/searchProducerSupplyOrder.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer  Page Success")
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

	auth, err := utils.IsAuthorized(utype, "ProducerSupplyOrder")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	select_producer_name := c.Input().Get("select_producer_name")
	//input_producer_mob := c.Input().Get("input_producer_mob")
	input_email := c.Input().Get("input_email")
	input_status := c.Input().Get("input_status")
	input_supplyordertxn_id := c.Input().Get("input_supplyordertxn_id")

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_email - ", input_email)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "select_producer_name - ", select_producer_name)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_status - ", input_status)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_supplyordertxn_id - ", input_supplyordertxn_id)

	c.Data["DateRange"] = c.Input().Get("daterange")
	tmp := strings.Split(c.Input().Get("daterange"), " - ")

	from := ""
	to := ""
	if len(tmp) == 2 {
		from = tmp[0]
		to = tmp[1]
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "From", from)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "To", to)

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

		data2, err2 := services.GetProducers2(data1[0][1])
		if err2 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err2)
			err2 = errors.New("Producer fetch Failed")
			return
		}
		if len(data2) <= 0 {
			err2 = errors.New("Producer  Not Found")
			return
		}
		var Dis Display
		for i := range data2 {
			var d Field
			d.Id = data2[i][0]
			d.Mobile = data2[i][1]
			d.Name = data2[i][2]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

		row, err := db.Db.Query(`select id,supplyordertxn_id,product_name,product_sku,product_qty,product_package_qty,cluster_buying_price,product_supply_units,product_supply_total_price,status,producer_mob from SupplyOrderTxn where (producer_mob like ?) AND (supplyordertxn_id='' OR  supplyordertxn_id like ?)AND (status='' OR status like ?) AND (status="DEMAND_ORDER_CREATED" OR status="DEMAND_ORDER_ACCEPTED") AND cluster_name=?`, select_producer_name+"%", input_supplyordertxn_id+"%", input_status+"%", data1[0][3])
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get the error message info")
			return
		}
		defer sql.Close(row)
		_, data, err := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Message Scan Fail")
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
		if len(data) <= 0 {
			err = errors.New("No Data Found")
			return
		}

		var result []Row

		for i := range data {
			var r Row
			r.Id = data[i][0]
			r.SupplyTxnID = data[i][1]
			r.ProductName = data[i][2]
			r.ProductSKU = data[i][3]
			r.ProductQTY = data[i][4]
			r.ProductPackageQTY = data[i][5]
			r.ClusterBuyPrice = data[i][6]
			r.SupplyUnits = data[i][7]
			r.TotalPrice = data[i][8]
			r.Status = data[i][9]
			r.ProducerMob = data[i][10]
			result = append(result, r)
		}

		c.Data["CustomerData"] = result

	} else if utype == "CLUSTEREXECUTIVE" {

		// NO PRODUCER MGMT

	} else {

		// HOADMIN CODE

		data1, err := services.GetProducers()
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Producer fetch Failed")
			return
		}
		if len(data1) <= 0 {
			err = errors.New("Producer  Not Found")
			return
		}
		var Dis Display
		for i := range data1 {
			var d Field
			d.Id = data1[i][0]
			d.Mobile = data1[i][1]
			d.Name = data1[i][2]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

		row, err := db.Db.Query(`select id,supplyordertxn_id,product_name,product_sku,product_qty,product_package_qty,cluster_buying_price,product_supply_units,product_supply_total_price,status,producer_mob from SupplyOrderTxn where (producer_mob like ?) AND (supplyordertxn_id='' OR  supplyordertxn_id like ?)AND (status='' OR status like ?) AND (status="DEMAND_ORDER_CREATED" OR status="DEMAND_ORDER_ACCEPTED")`, select_producer_name+"%", input_supplyordertxn_id+"%", input_status+"%")
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get the error message info")
			return
		}
		defer sql.Close(row)
		_, data, err := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Message Scan Fail")
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
		if len(data) <= 0 {
			err = errors.New("No Data Found")
			return
		}

		var result []Row

		for i := range data {
			var r Row
			r.Id = data[i][0]
			r.SupplyTxnID = data[i][1]
			r.ProductName = data[i][2]
			r.ProductSKU = data[i][3]
			r.ProductQTY = data[i][4]
			r.ProductPackageQTY = data[i][5]
			r.ClusterBuyPrice = data[i][6]
			r.SupplyUnits = data[i][7]
			r.TotalPrice = data[i][8]
			r.Status = data[i][9]
			r.ProducerMob = data[i][10]
			result = append(result, r)
		}

		c.Data["CustomerData"] = result

	}

	// row, err := db.Db.Query(`SELECT id,first_name,email,status FROM sysuser.credential`)

	// row, err := db.Db.Query(`select first_name,mob_no,email,status,user_type from Sysuser_info where
	// DATE(create_date) >= STR_TO_DATE(?, '%m/%d/%Y') AND DATE(create_date) <= STR_TO_DATE(?, '%m/%d/%Y') ORDER BY create_date DESC`, from, to)

	//id, supplyordertxn_id, category_name, subcategory_name, product_name, product_sku, product_qty, product_quality_type, product_quality_mesure, product_shelf_life, product_package_type, product_package_qty, producer_mob, cluster_name, cluster_buying_price, product_supply_units, product_supply_total_price, status, cluster_id, product_id, producer_supplied_units, cluster_received_units, create_date, create_by

	// c.Data["ProducerMob"] = data[0][10]

	// row8, err := db.Db.Query(`SELECT status,count(producer_mob)as PSum FROM SupplyOrderTxn where status="DEMAND_ORDER_CREATED"`)

	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Unable to get dashboard data")
	// 	return
	// }
	// defer sql.Close(row8)
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row8)
	// _, data8, err := sql.Scan(row8)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Unable to get dashboard data")
	// 	return
	// }
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data8, "\nData len - ", len(data8))
	// if len(data8) <= 0 {
	// 	err = errors.New("Unable to get dashboard data")
	// 	return
	// }
	// //log.Println(beego.AppConfig.String("loglevel"), "Debug", "PMob - ", data8[0][0])
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "PStatus - ", data8[0][0])
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "PCount - ", data8[0][1])

	// //c.Data["producermobile"] = data8[0][0]
	// c.Data["status"] = data8[0][0]
	// c.Data["totalproducercount"] = data8[0][1]

	return
}
