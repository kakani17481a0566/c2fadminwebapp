package searchCompleteSupplyOrder

import (
	"c2f_admin_webapp/services"
	"c2f_admin_webapp/session"
	"runtime/debug"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"
	"errors"

	"c2f_admin_webapp/model/db"

	"strings"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

type Row struct {
	Id                   string
	SupplyTxnID          string
	ProductName          string
	ProductSKU           string
	ProductQTY           string
	ProductPackageQTY    string
	ClusterBuyPrice      string
	SupplyUnits          string
	TotalPrice           string
	Status               string
	ProducerMobileNumber string
}

type Display struct {
	Fields []Field
}

type Field struct {
	Id     string
	Mobile string
	Name   string
}
type SearchCompleteSupplyOrder struct {
	// beego.Controller
	middleware.BaseController
}

func (c *SearchCompleteSupplyOrder) Get() {
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

			c.TplName = "completesupplyorder/searchCompleteSupplyOrder/searchCompleteSupplyOrder.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "completesupplyorder/searchCompleteSupplyOrder/searchCompleteSupplyOrder.html"
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

	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1

	auth, err := utils.IsAuthorized(utype, "CompleteSupplyOrder")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

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

	} else {

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

	}

	// row, err := db.Db.Query(`select id,supplyordertxn_id,product_name,product_sku,product_qty,product_package_qty,cluster_buying_price,product_supply_units,product_supply_total_price,status from SupplyOrderTxn`)

	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Unable to get the error message info")
	// 	return
	// }
	// defer sql.Close(row)
	// _, data, err1 := sql.Scan(row)
	// if err1 != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
	// 	err1 = errors.New("Message Scan Fail")
	// 	return
	// }
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	// if len(data) <= 0 {
	// 	err1 = errors.New("Unable to get the error message info")
	// 	return
	// }

	// var result []Row

	// for i := range data {
	// 	var r Row
	// 	r.Id = data[i][0]
	// 	r.SupplyTxnID = data[i][1]
	// 	r.ProductName = data[i][2]
	// 	r.ProductSKU = data[i][3]
	// 	r.ProductQTY = data[i][4]
	// 	r.ProductPackageQTY = data[i][5]
	// 	r.ClusterBuyPrice = data[i][6]
	// 	r.SupplyUnits = data[i][7]
	// 	r.TotalPrice = data[i][8]
	// 	r.Status = data[i][9]
	// 	result = append(result, r)
	// }

	// c.Data["CustomerData"] = result

	return
}

func (c *SearchCompleteSupplyOrder) Post() {
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
			c.TplName = "completesupplyorder/searchCompleteSupplyOrder/searchCompleteSupplyOrder.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "completesupplyorder/searchCompleteSupplyOrder/searchCompleteSupplyOrder.html"
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

	auth, err := utils.IsAuthorized(utype, "CompleteSupplyOrder")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	//	select_producer_name := c.Input().Get("select_producer_name")
	select_producer_name := c.Input().Get("select_producer_name")
	input_status := c.Input().Get("input_status")
	input_supplyordertxn_id := c.Input().Get("input_supplyordertxn_id")

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

	// row, err := db.Db.Query(`SELECT id,first_name,email,status FROM sysuser.credential`)

	// row, err := db.Db.Query(`select first_name,mob_no,email,status,user_type from Sysuser_info where
	// DATE(create_date) >= STR_TO_DATE(?, '%m/%d/%Y') AND DATE(create_date) <= STR_TO_DATE(?, '%m/%d/%Y') ORDER BY create_date DESC`, from, to)

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
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err2 = errors.New("Producer fetch Failed")
			return
		}
		if len(data2) <= 0 {
			err = errors.New("Producer  Not Found")
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

		row, err := db.Db.Query(`select id,supplyordertxn_id,product_name,product_sku,product_qty,product_package_qty,cluster_buying_price,product_supply_units,product_supply_total_price,status,producer_mob from SupplyOrderTxn where cluster_name=? AND  (producer_mob like ?) AND (status='' OR status like ?) AND (supplyordertxn_id='' OR supplyordertxn_id like ?) AND status="SUPPLY_RECEIVED"`, data1[0][3], select_producer_name+"%", input_status+"%", input_supplyordertxn_id+"%")
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
			err = errors.New("Unable to get the error message info")
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
			r.ProducerMobileNumber = data[i][10]

			result = append(result, r)
		}

		c.Data["CustomerData"] = result

	} else if utype == "CLUSTEREXECUTIVE" {

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

		row, err := db.Db.Query(`select id,supplyordertxn_id,product_name,product_sku,product_qty,product_package_qty,cluster_buying_price,product_supply_units,product_supply_total_price,status,producer_mob from SupplyOrderTxn where cluster_name=? AND (producer_mob like ?) AND (status='' OR status like ?) AND (supplyordertxn_id='' OR supplyordertxn_id like ?) AND status="SUPPLY_RECEIVED"`, data1[0][3], select_producer_name+"%", input_status+"%", input_supplyordertxn_id+"%")
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
			err = errors.New("Unable to get the error message info")
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
			r.ProducerMobileNumber = data[i][10]

			result = append(result, r)
		}

		c.Data["CustomerData"] = result

	} else {
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

		row, err := db.Db.Query(`select id,supplyordertxn_id,product_name,product_sku,product_qty,product_package_qty,cluster_buying_price,product_supply_units,product_supply_total_price,status,producer_mob from SupplyOrderTxn where (producer_mob like ?) AND (status='' OR status like ?) AND (supplyordertxn_id='' OR supplyordertxn_id like ?) AND status="SUPPLY_RECEIVED"`, select_producer_name+"%", input_status+"%", input_supplyordertxn_id+"%")
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
			err = errors.New("Unable to get the error message info")
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
			r.ProducerMobileNumber = data[i][10]

			result = append(result, r)
		}

		c.Data["CustomerData"] = result

	}

	return
}
