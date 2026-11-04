package supplyorderAssignmentReport

import (
	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/model/utils"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/services"
	"c2f_admin_webapp/session"

	"strings"

	"runtime/debug"

	"errors"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

// type Row struct {
// 	Producer_Mob string
// 	Product_Name string
// 	Cluster_Name string
// 	Qty          string
// 	Price        string
// 	Start_Date   string
// 	End_Date     string
// 	TimeStmp     string
// }

type Row struct {
	TimeStmp                 string
	SupplyOrderAssignmentID  string
	Mobile                   string
	Address1                 string
	Address2                 string
	Address                  string
	Product_Name             string
	Product_Qty              string
	SupplyUnits              string
	SupplyExteciutiveCluster string
	Product_Sku              string
	ProducerName             string
	City                     string
	Zipcode                  string
	Pckqty                   string
	Status                   string
	Id                       string
}

type SupplyorderAssignmentReport struct {
	// beego.Controller
	middleware.BaseController
}

func (c *SupplyorderAssignmentReport) Get() {
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
			c.TplName = "reports/supplyorderAssignmentReport/supplyorderAssignmentReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Supply Order Assignment Report Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "reports/supplyorderAssignmentReport/supplyorderAssignmentReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Supply Order Assignment Report Page Success")
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
	//	uid := sess.Get("uid").(string)
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

	auth, err := utils.IsAuthorized(utype, "SupplyOrAssigReport")
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

		row, err := db.Db.Query(`select create_date,supply_order_id,supplyexecutive_mob,JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_name')),
		JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_qty')),JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_supply_units')),
		JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_sku')),JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.address1')),
		JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.address2')),JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.name')),
		JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.city')),JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.zipcode')),
		JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_package_qty')),supplyexecutive_cluster,status,id from Oder2SupplyExecutive
		 where (supplyexecutive_cluster='' OR supplyexecutive_cluster like ?)`, cname)

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
			r.TimeStmp = data[i][0]
			r.SupplyOrderAssignmentID = data[i][1]
			r.Mobile = data[i][2]
			r.Product_Name = data[i][3]
			r.Product_Qty = data[i][4]
			r.SupplyUnits = data[i][5]
			r.Product_Sku = data[i][6]
			r.Address1 = data[i][7]
			r.Address2 = data[i][8]
			r.ProducerName = data[i][9]
			r.City = data[i][10]
			r.Zipcode = data[i][11]
			r.Pckqty = data[i][12] + data[i][4]
			r.SupplyExteciutiveCluster = data[i][13]
			r.Status = data[i][14]
			r.Id = data[i][15]

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

		row, err := db.Db.Query(`select create_date,supply_order_id,supplyexecutive_mob,JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_name')),
		JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_qty')),JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_supply_units')),
		JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_sku')),JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.address1')),
		JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.address2')),JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.name')),
		JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.city')),JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.zipcode')),
		JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_package_qty')),supplyexecutive_cluster,status,id from Oder2SupplyExecutive
		 where (supplyexecutive_cluster='' OR supplyexecutive_cluster like ?)`, cname)
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
			r.TimeStmp = data[i][0]
			r.SupplyOrderAssignmentID = data[i][1]
			r.Mobile = data[i][2]
			r.Product_Name = data[i][3]
			r.Product_Qty = data[i][4]
			r.SupplyUnits = data[i][5]
			r.Product_Sku = data[i][6]
			r.Address1 = data[i][7]
			r.Address2 = data[i][8]
			r.ProducerName = data[i][9]
			r.City = data[i][10]
			r.Zipcode = data[i][11]
			r.Pckqty = data[i][12] + data[i][4]
			r.SupplyExteciutiveCluster = data[i][13]
			r.Status = data[i][14]
			r.Id = data[i][15]

			result = append(result, r)
		}
		c.Data["CustomerData"] = result
	} else {
		row, err := db.Db.Query(`select create_date,supply_order_id,supplyexecutive_mob,JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_name')),
		JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_qty')),JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_supply_units')),
		JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_sku')),JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.address1')),
		JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.address2')),JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.name')),
		JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.city')),JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.zipcode')),
		JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_package_qty')),supplyexecutive_cluster,status,id from Oder2SupplyExecutive`)

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
			r.TimeStmp = data[i][0]
			r.SupplyOrderAssignmentID = data[i][1]
			r.Mobile = data[i][2]
			r.Product_Name = data[i][3]
			r.Product_Qty = data[i][4]
			r.SupplyUnits = data[i][5]
			r.Product_Sku = data[i][6]
			r.Address1 = data[i][7]
			r.Address2 = data[i][8]
			r.ProducerName = data[i][9]
			r.City = data[i][10]
			r.Zipcode = data[i][11]
			r.Pckqty = data[i][12] + data[i][4]
			r.SupplyExteciutiveCluster = data[i][13]
			r.Status = data[i][14]
			r.Id = data[i][15]

			result = append(result, r)

		}
		c.Data["CustomerData"] = result

	}

	return
}

func (c *SupplyorderAssignmentReport) Post() {
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
			c.TplName = "reports/supplyorderAssignmentReport/supplyorderAssignmentReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Supply Order Assignment Report Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "reports/supplyorderAssignmentReport/supplyorderAssignmentReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Supply Order Assignment Report  Page Success")
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
	//uid := sess.Get("uid").(string)
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

	// cid, err := services.GetCid(uid)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Cluster ID fetch Failed")
	// 	return
	// }

	// cname, err := services.GetClustername(cid)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Cluster Name fetch Failed")
	// 	return
	// }

	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "cname", cname)

	auth, err := utils.IsAuthorized(utype, "SupplyOrAssigReport")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	input_rnn_no := c.Input().Get("input_rnn_no")
	input_status := c.Input().Get("input_status")
	input_card_mask := c.Input().Get("input_card_mask")
	input_product_name := c.Input().Get("input_product_name")
	input_cluster_name := c.Input().Get("input_cluster_name")
	input_supplyexe_mob := c.Input().Get("input_supplyexe_mob")
	input_order_id := c.Input().Get("input_order_id")

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
		row, err := db.Db.Query(`select create_date,supply_order_id,supplyexecutive_mob,JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_name')),
		JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_qty')),JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_supply_units')),
		JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_sku')),JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.address1')),
		JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.address2')),JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.name')),
		JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.city')),JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.zipcode')),
		JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_package_qty')),supplyexecutive_cluster,status,id from Oder2SupplyExecutive
		 where (supplyexecutive_cluster='' OR supplyexecutive_cluster like ?) AND (JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_name'))='' OR JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_name')) like ?)  
		AND  (supplyexecutive_mob='' OR supplyexecutive_mob like ?)  
		AND (supplyexecutive_cluster='' OR supplyexecutive_cluster like ?) AND (supply_order_id='' OR supply_order_id like ?) 
		AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) 
		ORDER BY create_date DESC`, cname, input_product_name+"%", input_supplyexe_mob+"%", input_cluster_name+"%", input_order_id+"%", from, to)

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
			r.TimeStmp = data[i][0]
			r.SupplyOrderAssignmentID = data[i][1]
			r.Mobile = data[i][2]
			r.Product_Name = data[i][3]
			r.Product_Qty = data[i][4]
			r.SupplyUnits = data[i][5]
			r.Product_Sku = data[i][6]
			r.Address1 = data[i][7]
			r.Address2 = data[i][8]
			r.ProducerName = data[i][9]
			r.City = data[i][10]
			r.Zipcode = data[i][11]
			r.Pckqty = data[i][12] + data[i][4]
			r.SupplyExteciutiveCluster = data[i][13]
			r.Status = data[i][14]
			r.Id = data[i][15]

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
		row, err := db.Db.Query(`select create_date,supply_order_id,supplyexecutive_mob,JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_name')),
		JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_qty')),JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_supply_units')),
		JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_sku')),JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.address1')),
		JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.address2')),JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.name')),
		JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.city')),JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.zipcode')),
		JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_package_qty')),supplyexecutive_cluster,status,id from Oder2SupplyExecutive
		 where (supplyexecutive_cluster='' OR supplyexecutive_cluster like ?) AND (JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_name'))='' OR JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_name')) like ?)  
		AND  (supplyexecutive_mob='' OR supplyexecutive_mob like ?)  
		AND (supplyexecutive_cluster='' OR supplyexecutive_cluster like ?) AND (supply_order_id='' OR supply_order_id like ?) 
		AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) 
		ORDER BY create_date DESC`, cname, input_product_name+"%", input_supplyexe_mob+"%", input_cluster_name+"%", input_order_id+"%", from, to)
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
			r.TimeStmp = data[i][0]
			r.SupplyOrderAssignmentID = data[i][1]
			r.Mobile = data[i][2]
			r.Product_Name = data[i][3]
			r.Product_Qty = data[i][4]
			r.SupplyUnits = data[i][5]
			r.Product_Sku = data[i][6]
			r.Address1 = data[i][7]
			r.Address2 = data[i][8]
			r.ProducerName = data[i][9]
			r.City = data[i][10]
			r.Zipcode = data[i][11]
			r.Pckqty = data[i][12] + data[i][4]
			r.SupplyExteciutiveCluster = data[i][13]
			r.Status = data[i][14]
			r.Id = data[i][15]
			result = append(result, r)
		}
		c.Data["CustomerData"] = result
	} else {
		row, err := db.Db.Query(`select create_date,supply_order_id,supplyexecutive_mob,JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_name')),
		JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_qty')),JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_supply_units')),
		JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_sku')),JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.address1')),
		JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.address2')),JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.name')),
		JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.city')),JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.zipcode')),
		JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_package_qty')),supplyexecutive_cluster,status,id from Oder2SupplyExecutive
		 where (JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_name'))='' OR JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_name')) like ?)  
		AND  (supplyexecutive_mob='' OR supplyexecutive_mob like ?)  
		AND (supplyexecutive_cluster='' OR supplyexecutive_cluster like ?) AND (supply_order_id='' OR supply_order_id like ?) 
		AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) 
		ORDER BY create_date DESC`, input_product_name+"%", input_supplyexe_mob+"%", input_cluster_name+"%", input_order_id+"%", from, to)
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
			r.TimeStmp = data[i][0]
			r.SupplyOrderAssignmentID = data[i][1]
			r.Mobile = data[i][2]
			r.Product_Name = data[i][3]
			r.Product_Qty = data[i][4]
			r.SupplyUnits = data[i][5]
			r.Product_Sku = data[i][6]
			r.Address1 = data[i][7]
			r.Address2 = data[i][8]
			r.ProducerName = data[i][9]
			r.City = data[i][10]
			r.Zipcode = data[i][11]
			r.Pckqty = data[i][12] + data[i][4]
			r.SupplyExteciutiveCluster = data[i][13]
			r.Status = data[i][14]
			r.Id = data[i][15]
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
