package inventoryReport

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

type Row struct {
	TimeStamp                string
	Cluster_name             string
	Category_name            string
	Subcategory_name         string
	Product_name             string
	Product_sku              string
	Product_qty              string
	Starting_inventory       string
	Inventory_received       string
	Product_quality_mesure   string
	Product_shelf_life       string
	Product_package_type     string
	Product_package_qty      string
	Product_origin           string
	Krishivala_selling_price string
	Cluster_buying_price     string
	Cluster_selling_price    string
	Producer_selling_price   string
	Received_units           string
	Avaliable_units          string
	Deliverd_units           string
	Supplycharge             string
	Packagecharge            string
	Bankcharge               string
	Discount                 string
}

type InventoryReport struct {
	// beego.Controller
	middleware.BaseController
}

func (c *InventoryReport) Get() {
	Cname := c.Ctx.Input.Param(":Cname")

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
			c.TplName = "reports/inventoryReport/inventoryReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Reports Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "reports/inventoryReport/inventoryReport.html"
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

	auth, err := utils.IsAuthorized(utype, "InventReport")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	if Cname == "" {
		row, err := db.Db.Query(`select create_date, 
	    cluster_name,
		category_name, 
		subcategory_name,
		product_name,
		product_sku,
		product_qty,
		starting_inventory,
		inventory_received, 
		product_quality_mesure,
		product_shelf_life,
		product_package_type,
		product_package_qty,
		product_origin,
		krishivala_selling_price,
		cluster_buying_price, 
		cluster_selling_price,
		producer_selling_price,
		inventory_received,
		inventory_available,
		inventory_delivered,
		cluster_supply_charge,
		cluster_packaging_charge,
		cluster_bank_charge,
		cluster_discount
		from Inventory`)
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
			r.Cluster_name = data[i][1]
			r.Category_name = data[i][2]
			r.Subcategory_name = data[i][3]
			r.Product_name = data[i][4]
			r.Product_sku = data[i][5]
			r.Product_qty = data[i][6]
			r.Starting_inventory = data[i][7]
			r.Inventory_received = data[i][8]
			r.Product_quality_mesure = data[i][9]
			r.Product_shelf_life = data[i][10]
			r.Product_package_type = data[i][11]
			r.Product_package_qty = data[i][12]
			r.Product_origin = data[i][13]
			r.Krishivala_selling_price = data[i][14]
			r.Cluster_buying_price = data[i][15]
			r.Cluster_selling_price = data[i][16]
			r.Producer_selling_price = data[i][17]
			r.Received_units = data[i][18]
			r.Avaliable_units = data[i][19]
			r.Deliverd_units = data[i][20]
			r.Supplycharge = data[i][21]
			r.Packagecharge = data[i][22]
			r.Bankcharge = data[i][23]
			r.Discount = data[i][24]

			result = append(result, r)
		}
		c.Data["CustomerData"] = result
	} else {

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Cname", Cname)

		row, err := db.Db.Query(`select create_date, 
	    cluster_name,
		category_name, 
		subcategory_name,
		product_name,
		product_sku,
		product_qty,
		starting_inventory,
		inventory_received, 
		product_quality_mesure,
		product_shelf_life,
		product_package_type,
		product_package_qty,
		product_origin,
		krishivala_selling_price,
		cluster_buying_price, 
		cluster_selling_price,
		producer_selling_price,
		inventory_received,
		inventory_available,
		inventory_delivered,
		cluster_supply_charge,
		cluster_packaging_charge,
		cluster_bank_charge,
		cluster_discount
		from Inventory where (cluster_name='' OR cluster_name like ?)`, Cname)
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
			r.Cluster_name = data[i][1]
			r.Category_name = data[i][2]
			r.Subcategory_name = data[i][3]
			r.Product_name = data[i][4]
			r.Product_sku = data[i][5]
			r.Product_qty = data[i][6]
			r.Starting_inventory = data[i][7]
			r.Inventory_received = data[i][8]
			r.Product_quality_mesure = data[i][9]
			r.Product_shelf_life = data[i][10]
			r.Product_package_type = data[i][11]
			r.Product_package_qty = data[i][12]
			r.Product_origin = data[i][13]
			r.Krishivala_selling_price = data[i][14]
			r.Cluster_buying_price = data[i][15]
			r.Cluster_selling_price = data[i][16]
			r.Producer_selling_price = data[i][17]
			r.Received_units = data[i][18]
			r.Avaliable_units = data[i][19]
			r.Deliverd_units = data[i][20]
			r.Supplycharge = data[i][21]
			r.Packagecharge = data[i][22]
			r.Bankcharge = data[i][23]
			r.Discount = data[i][24]

			result = append(result, r)
		}
		c.Data["CustomerData"] = result
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

		row, err := db.Db.Query(`select create_date,
	    cluster_name,
		category_name,
		subcategory_name,
		product_name,
		product_sku,
		product_qty,
		starting_inventory,
		inventory_received,
		product_quality_mesure,
		product_shelf_life,
		product_package_type,
		product_package_qty,
		product_origin,
		krishivala_selling_price,
		cluster_buying_price,
		cluster_selling_price,
		producer_selling_price,
		inventory_received,
		inventory_available,
		inventory_delivered,
		cluster_supply_charge,
		cluster_packaging_charge,
		cluster_bank_charge,
		cluster_discount
		from Inventory where (cluster_name='' OR cluster_name like ?)`, cname)
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
			r.Cluster_name = data[i][1]
			r.Category_name = data[i][2]
			r.Subcategory_name = data[i][3]
			r.Product_name = data[i][4]
			r.Product_sku = data[i][5]
			r.Product_qty = data[i][6]
			r.Starting_inventory = data[i][7]
			r.Inventory_received = data[i][8]
			r.Product_quality_mesure = data[i][9]
			r.Product_shelf_life = data[i][10]
			r.Product_package_type = data[i][11]
			r.Product_package_qty = data[i][12]
			r.Product_origin = data[i][13]
			r.Krishivala_selling_price = data[i][14]
			r.Cluster_buying_price = data[i][15]
			r.Cluster_selling_price = data[i][16]
			r.Producer_selling_price = data[i][17]
			r.Received_units = data[i][18]
			r.Avaliable_units = data[i][19]
			r.Deliverd_units = data[i][20]
			r.Supplycharge = data[i][21]
			r.Packagecharge = data[i][22]
			r.Bankcharge = data[i][23]
			r.Discount = data[i][24]

			result = append(result, r)
		}
		c.Data["CustomerData"] = result

	}
	if utype == "CLUSTEREXECUTIVE" {

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

		row, err := db.Db.Query(`select create_date,
	    cluster_name,
		category_name,
		subcategory_name,
		product_name,
		product_sku,
		product_qty,
		starting_inventory,
		inventory_received,
		product_quality_mesure,
		product_shelf_life,
		product_package_type,
		product_package_qty,
		product_origin,
		krishivala_selling_price,
		cluster_buying_price,
		cluster_selling_price,
		producer_selling_price,
		inventory_received,
		inventory_available,
		inventory_delivered,
		cluster_supply_charge,
		cluster_packaging_charge,
		cluster_bank_charge,
		cluster_discount
		from Inventory where (cluster_name='' OR cluster_name like ?)`, cname)
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
			r.Cluster_name = data[i][1]
			r.Category_name = data[i][2]
			r.Subcategory_name = data[i][3]
			r.Product_name = data[i][4]
			r.Product_sku = data[i][5]
			r.Product_qty = data[i][6]
			r.Starting_inventory = data[i][7]
			r.Inventory_received = data[i][8]
			r.Product_quality_mesure = data[i][9]
			r.Product_shelf_life = data[i][10]
			r.Product_package_type = data[i][11]
			r.Product_package_qty = data[i][12]
			r.Product_origin = data[i][13]
			r.Krishivala_selling_price = data[i][14]
			r.Cluster_buying_price = data[i][15]
			r.Cluster_selling_price = data[i][16]
			r.Producer_selling_price = data[i][17]
			r.Received_units = data[i][18]
			r.Avaliable_units = data[i][19]
			r.Deliverd_units = data[i][20]
			r.Supplycharge = data[i][21]
			r.Packagecharge = data[i][22]
			r.Bankcharge = data[i][23]
			r.Discount = data[i][24]

			result = append(result, r)
		}
		c.Data["CustomerData"] = result

	}

	return
}

func (c *InventoryReport) Post() {
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
			c.TplName = "reports/inventoryReport/inventoryReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchOrder Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "reports/inventoryReport/inventoryReport.html"
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

	auth, err := utils.IsAuthorized(utype, "InventReport")
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

		row, err := db.Db.Query(`select create_date, 
	    cluster_name,
		category_name, 
		subcategory_name,
		product_name,
		product_sku,
		product_qty,
		starting_inventory,
		inventory_received, 
		product_quality_mesure,
		product_shelf_life,
		product_package_type,
		product_package_qty,
		product_origin,
		krishivala_selling_price,
		cluster_buying_price, 
		cluster_selling_price,
		producer_selling_price,
		inventory_received,
		inventory_available,
		inventory_delivered,
		cluster_supply_charge,
		cluster_packaging_charge,
		cluster_bank_charge,
		cluster_discount
		from Inventory where (cluster_name='' OR cluster_name like ?) AND (product_name='' OR product_name like ?) AND (cluster_name='' OR cluster_name like ?)  AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) ORDER BY create_date DESC`, cname, input_product_name+"%", input_cluster_name+"%", from, to)

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
			r.Cluster_name = data[i][1]
			r.Category_name = data[i][2]
			r.Subcategory_name = data[i][3]
			r.Product_name = data[i][4]
			r.Product_sku = data[i][5]
			r.Product_qty = data[i][6]
			r.Starting_inventory = data[i][7]
			r.Inventory_received = data[i][8]
			r.Product_quality_mesure = data[i][9]
			r.Product_shelf_life = data[i][10]
			r.Product_package_type = data[i][11]
			r.Product_package_qty = data[i][12]
			r.Product_origin = data[i][13]
			r.Krishivala_selling_price = data[i][14]
			r.Cluster_buying_price = data[i][15]
			r.Cluster_selling_price = data[i][16]
			r.Producer_selling_price = data[i][17]
			r.Received_units = data[i][18]
			r.Avaliable_units = data[i][19]
			r.Deliverd_units = data[i][20]
			r.Supplycharge = data[i][21]
			r.Packagecharge = data[i][22]
			r.Bankcharge = data[i][23]
			r.Discount = data[i][24]

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

		row, err := db.Db.Query(`select create_date, 
	    cluster_name,
		category_name, 
		subcategory_name,
		product_name,
		product_sku,
		product_qty,
		starting_inventory,
		inventory_received, 
		product_quality_mesure,
		product_shelf_life,
		product_package_type,
		product_package_qty,
		product_origin,
		krishivala_selling_price,
		cluster_buying_price, 
		cluster_selling_price,
		producer_selling_price,
		inventory_received,
		inventory_available,
		inventory_delivered,
		cluster_supply_charge,
		cluster_packaging_charge,
		cluster_bank_charge,
		cluster_discount
		from Inventory where (cluster_name='' OR cluster_name like ?) AND (product_name='' OR product_name like ?) AND (cluster_name='' OR cluster_name like ?)  AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) ORDER BY create_date DESC`, cname, input_product_name+"%", input_cluster_name+"%", from, to)

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
			r.Cluster_name = data[i][1]
			r.Category_name = data[i][2]
			r.Subcategory_name = data[i][3]
			r.Product_name = data[i][4]
			r.Product_sku = data[i][5]
			r.Product_qty = data[i][6]
			r.Starting_inventory = data[i][7]
			r.Inventory_received = data[i][8]
			r.Product_quality_mesure = data[i][9]
			r.Product_shelf_life = data[i][10]
			r.Product_package_type = data[i][11]
			r.Product_package_qty = data[i][12]
			r.Product_origin = data[i][13]
			r.Krishivala_selling_price = data[i][14]
			r.Cluster_buying_price = data[i][15]
			r.Cluster_selling_price = data[i][16]
			r.Producer_selling_price = data[i][17]
			r.Received_units = data[i][18]
			r.Avaliable_units = data[i][19]
			r.Deliverd_units = data[i][20]
			r.Supplycharge = data[i][21]
			r.Packagecharge = data[i][22]
			r.Bankcharge = data[i][23]
			r.Discount = data[i][24]

			result = append(result, r)
		}
		c.Data["CustomerData"] = result
	} else {
		row, err := db.Db.Query(`select create_date, 
	    cluster_name,
		category_name, 
		subcategory_name,
		product_name,
		product_sku,
		product_qty,
		starting_inventory,
		inventory_received, 
		product_quality_mesure,
		product_shelf_life,
		product_package_type,
		product_package_qty,
		product_origin,
		krishivala_selling_price,
		cluster_buying_price, 
		cluster_selling_price,
		producer_selling_price,
		inventory_received,
		inventory_available,
		inventory_delivered,
		cluster_supply_charge,
		cluster_packaging_charge,
		cluster_bank_charge,
		cluster_discount
		from Inventory where (product_name='' OR product_name like ?) AND (cluster_name='' OR cluster_name like ?)  AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) ORDER BY create_date DESC`, input_product_name+"%", input_cluster_name+"%", from, to)

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
			r.Cluster_name = data[i][1]
			r.Category_name = data[i][2]
			r.Subcategory_name = data[i][3]
			r.Product_name = data[i][4]
			r.Product_sku = data[i][5]
			r.Product_qty = data[i][6]
			r.Starting_inventory = data[i][7]
			r.Inventory_received = data[i][8]
			r.Product_quality_mesure = data[i][9]
			r.Product_shelf_life = data[i][10]
			r.Product_package_type = data[i][11]
			r.Product_package_qty = data[i][12]
			r.Product_origin = data[i][13]
			r.Krishivala_selling_price = data[i][14]
			r.Cluster_buying_price = data[i][15]
			r.Cluster_selling_price = data[i][16]
			r.Producer_selling_price = data[i][17]
			r.Received_units = data[i][18]
			r.Avaliable_units = data[i][19]
			r.Deliverd_units = data[i][20]
			r.Supplycharge = data[i][21]
			r.Packagecharge = data[i][22]
			r.Bankcharge = data[i][23]
			r.Discount = data[i][24]

			result = append(result, r)
		}
		c.Data["CustomerData"] = result
	}

	return
}
