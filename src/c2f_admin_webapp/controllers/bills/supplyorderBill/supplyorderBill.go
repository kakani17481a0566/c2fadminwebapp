package supplyorderBill

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
	Id                string
	SupplyorderId     string
	SupplyorderDate   string
	ProducerMob       string
	Status            string
	ProductName       string
	ProductQty        string
	ProductPackageQty string
	TotalBuyPrice     string
}
type Display5 struct {
	Fields5 []Field5
}

type Field5 struct {
	Id     string
	Mobile string
	Name   string
}
type SupplyorderBill struct {
	// beego.Controller
	middleware.BaseController
}

func (c *SupplyorderBill) Get() {
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
			c.TplName = "bills/producersupplyorderBill/producersupplyorderBill.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Reports Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "bills/producersupplyorderBill/producersupplyorderBill.html"
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

	id := sess.Get("uid").(string)
	c.Data["id"] = id

	auth, err := utils.IsAuthorized(utype, "SupplyOrderBill")
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
		var Dis Display5
		for i := range data {
			var d Field5
			d.Id = data[i][0]
			d.Mobile = data[i][1]
			d.Name = data[i][2]
			Dis.Fields5 = append(Dis.Fields5, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

	} else if utype == "CLUSTEREXECUTIVE" {

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
		var Dis Display5
		for i := range data {
			var d Field5
			d.Id = data[i][0]
			d.Mobile = data[i][1]
			d.Name = data[i][2]
			Dis.Fields5 = append(Dis.Fields5, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

	}

	// row, err := db.Db.Query(`select id,create_date,supplyordertxn_id, producer_mob,status,product_name,product_qty,product_package_qty,total_buy_price from SupplyOrderTxn ORDER BY create_date DESC`)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Unable to fetch order")
	// 	return
	// }
	// defer sql.Close(row)
	// _, data, err1 := sql.Scan(row)
	// if err1 != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
	// 	err1 = errors.New("Unable to fetch order")
	// 	return
	// }

	// log.Println(beego.AppConfig.String("loglevel"), "Debug", data)
	// var result []Row

	// for i := range data {
	// 	var r Row
	// 	r.Id = data[i][0]
	// 	r.SupplyorderDate = data[i][1]
	// 	r.SupplyorderId = data[i][2]
	// 	r.ProducerMob = data[i][3]
	// 	r.Status = data[i][4]
	// 	r.ProductName = data[i][5]
	// 	r.ProductQty = data[i][6]
	// 	r.ProductPackageQty = data[i][7]
	// 	r.TotalBuyPrice = data[i][8]

	// 	result = append(result, r)
	// }
	// c.Data["CustomerData"] = result

	return
}

func (c *SupplyorderBill) Post() {
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
			c.TplName = "bills/producersupplyorderBill/producersupplyorderBill.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchOrder Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "bills/producersupplyorderBill/producersupplyorderBill.html"
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
	var Dis Display5
	for i := range data1 {
		var d Field5
		d.Id = data1[i][0]
		d.Mobile = data1[i][1]
		d.Name = data1[i][2]
		Dis.Fields5 = append(Dis.Fields5, d)
	}
	c.Data["Dis"] = Dis
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

	auth, err := utils.IsAuthorized(utype, "SupplyOrderBill")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	input_rnn_no := c.Input().Get("input_rnn_no")
	input_status := c.Input().Get("input_status")
	input_card_mask := c.Input().Get("input_card_mask")
	input_producer_mob := c.Input().Get("input_producer_mob")
	//input_status := c.Input().Get("input_status")
	input_SupplyorderId := c.Input().Get("input_SupplyorderId")

	c.Data["RRN"] = input_rnn_no
	c.Data["Mask"] = input_card_mask
	c.Data["DateRange"] = c.Input().Get("daterange")
	c.Data["Status"] = input_status
	c.Data["SupplyorderId"] = input_SupplyorderId
	tmp := strings.Split(c.Input().Get("daterange"), " - ")

	from := ""
	to := ""
	if len(tmp) == 2 {
		from = tmp[0]
		to = tmp[1]
	}

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
		var Dis Display5
		for i := range data2 {
			var d Field5
			d.Id = data2[i][0]
			d.Mobile = data2[i][1]
			d.Name = data2[i][2]
			Dis.Fields5 = append(Dis.Fields5, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

		row, err := db.Db.Query(`select id,create_date,supplyordertxn_id, producer_mob,status,product_name,product_qty,product_package_qty,total_buy_price,cluster_name from SupplyOrderTxn 
where cluster_name=? AND (status='' OR status like ?) AND (supplyordertxn_id='' OR supplyordertxn_id like ?) AND (producer_mob='' OR producer_mob like ?)  AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) AND (status="SUPPLY_COMPLETED") ORDER BY create_date DESC`, data1[0][3], input_status+"%", input_SupplyorderId+"%", input_producer_mob+"%", from, to)

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
			r.Id = data[i][0]
			r.SupplyorderDate = data[i][1]
			r.SupplyorderId = data[i][2]
			r.ProducerMob = data[i][3]
			r.Status = data[i][4]
			r.ProductName = data[i][5]
			r.ProductQty = data[i][6]
			r.ProductPackageQty = data[i][7]
			r.TotalBuyPrice = data[i][8]

			result = append(result, r)
		}
		c.Data["CustomerData"] = result

	} else if utype == "CLUSTEREXECUTIVE" {

		//NO CODE FOR CLUSTEREXECUTIVE

	} else {
		data1, err1 := services.GetProducers()
		if err1 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
			err1 = errors.New("Producer fetch Failed")
			return
		}
		if len(data1) <= 0 {
			err1 = errors.New("Producer  Not Found")
			return
		}
		var Dis Display5
		for i := range data1 {
			var d Field5
			d.Id = data1[i][0]
			d.Mobile = data1[i][1]
			d.Name = data1[i][2]
			Dis.Fields5 = append(Dis.Fields5, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

		row, err := db.Db.Query(`select id,create_date,supplyordertxn_id, producer_mob,status,product_name,product_qty,product_package_qty,total_buy_price from SupplyOrderTxn 
where (status='' OR status like ?) AND (supplyordertxn_id='' OR supplyordertxn_id like ?) AND (producer_mob='' OR producer_mob like ?)  AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) AND (status="SUPPLY_COMPLETED") ORDER BY create_date DESC`, input_status+"%", input_SupplyorderId+"%", input_producer_mob+"%", from, to)

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
			r.Id = data[i][0]
			r.SupplyorderDate = data[i][1]
			r.SupplyorderId = data[i][2]
			r.ProducerMob = data[i][3]
			r.Status = data[i][4]
			r.ProductName = data[i][5]
			r.ProductQty = data[i][6]
			r.ProductPackageQty = data[i][7]
			r.TotalBuyPrice = data[i][8]

			result = append(result, r)
		}
		c.Data["CustomerData"] = result

	}

	return
}

// row, err := db.Db.Query(`select producer_mob, product_name, cluster_name, qty, unit_price,price, start_date, end_date, create_date from SupplyForecastTxn`)

// if err != nil {
// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
// 	err = errors.New("Unable to fetch order")
// 	return
// }
// defer sql.Close(row)
// _, data, err := sql.Scan(row)
// if err != nil {
// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
// 	err = errors.New("Unable to fetch order")
// 	return
// }

// log.Println(beego.AppConfig.String("loglevel"), "Debug", data)
// var result []Row

// for i := range data {
// 	var r Row
// 	r.Producer_Mob = data[i][0]
// 	r.Product_Name = data[i][1]
// 	r.Cluster_Name = data[i][2]
// 	r.Qty = data[i][3]
// 	r.Price = data[i][4]
// 	r.Start_Date = data[i][5]
// 	r.End_Date = data[i][6]
// 	r.TimeStmp = data[i][7]
// 	result = append(result, r)
// }
// c.Data["CustomerData"] = result

// return
// }
