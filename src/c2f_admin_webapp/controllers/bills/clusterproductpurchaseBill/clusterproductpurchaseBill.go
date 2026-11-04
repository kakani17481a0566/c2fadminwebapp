package clusterproductpurchaseBill

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
	TimeStamp                     string
	Id                            string
	Source_cluster_name           string
	Beneficiary_cluster_name      string
	Product_name                  string
	Product_sku                   string
	Number_of_units               string
	Delivery_charges              string
	Total_price                   string
	Status                        string
	Cluster_productpurchase_txnid string
	Product_price                 string
}

type Display struct {
	Fields []Field
}

type Field struct {
	Id   string
	Name string
}
type ClusterproductpurchaseBill struct {
	// beego.Controller
	middleware.BaseController
}

func (c *ClusterproductpurchaseBill) Get() {
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
			c.TplName = "bills/clusterproductpurchaseBill/clusterproductpurchaseBill.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Reports Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "bills/clusterproductpurchaseBill/clusterproductpurchaseBill.html"
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

	data, err := services.GetClusteruser()
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
		d.Name = data[i][1]
		Dis.Fields = append(Dis.Fields, d)
	}
	c.Data["Dis"] = Dis
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

	auth, err := utils.IsAuthorized(utype, "CluPductPurBill")
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
		// row, err := db.Db.Query(`select id,
		// source_cluster_name,
		// beneficiary_cluster_name,
		// product_name,
		// product_sku,
		// number_of_units,
		// delivery_charges,
		// total_price,
		// status,
		// create_date,
		// cluster_productpurchase_txnid,
		// product_price
		// from ClusterPurchaseProduct where (source_cluster_name='' OR source_cluster_name like ?)`, cname)

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
		// 	r.Source_cluster_name = data[i][1]
		// 	r.Beneficiary_cluster_name = data[i][2]
		// 	r.Product_name = data[i][3]
		// 	r.Product_sku = data[i][4]
		// 	r.Number_of_units = data[i][5]
		// 	r.Delivery_charges = data[i][6]
		// 	r.Total_price = data[i][7]
		// 	r.Status = data[i][8]
		// 	r.TimeStamp = data[i][9]
		// 	r.Cluster_productpurchase_txnid = data[i][10]
		// 	r.Product_price = data[i][11]

		// 	result = append(result, r)
		// }
		// c.Data["CustomerData"] = result
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
		// row, err := db.Db.Query(`select id,
		// source_cluster_name,
		// beneficiary_cluster_name,
		// product_name,
		// product_sku,
		// number_of_units,
		// delivery_charges,
		// total_price,
		// status,
		// create_date,
		// cluster_productpurchase_txnid,
		// product_price
		// from ClusterPurchaseProduct where (source_cluster_name='' OR source_cluster_name like ?)`, cname)

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
		// 	r.Source_cluster_name = data[i][1]
		// 	r.Beneficiary_cluster_name = data[i][2]
		// 	r.Product_name = data[i][3]
		// 	r.Product_sku = data[i][4]
		// 	r.Number_of_units = data[i][5]
		// 	r.Delivery_charges = data[i][6]
		// 	r.Total_price = data[i][7]
		// 	r.Status = data[i][8]
		// 	r.TimeStamp = data[i][9]
		// 	r.Cluster_productpurchase_txnid = data[i][10]
		// 	r.Product_price = data[i][11]

		// 	result = append(result, r)
		// }
		// c.Data["CustomerData"] = result
	} else {
		// row, err := db.Db.Query(`select id,
		// source_cluster_name,
		// beneficiary_cluster_name,
		// product_name,
		// product_sku,
		// number_of_units,
		// delivery_charges,
		// total_price,
		// status,
		// create_date,
		// cluster_productpurchase_txnid,
		// product_price
		// from ClusterPurchaseProduct ORDER BY id DESC`)

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
		// 	r.Source_cluster_name = data[i][1]
		// 	r.Beneficiary_cluster_name = data[i][2]
		// 	r.Product_name = data[i][3]
		// 	r.Product_sku = data[i][4]
		// 	r.Number_of_units = data[i][5]
		// 	r.Delivery_charges = data[i][6]
		// 	r.Total_price = data[i][7]
		// 	r.Status = data[i][8]
		// 	r.TimeStamp = data[i][9]
		// 	r.Cluster_productpurchase_txnid = data[i][10]
		// 	r.Product_price = data[i][11]

		// 	result = append(result, r)
		// }
		// c.Data["CustomerData"] = result
	}

	return
}

func (c *ClusterproductpurchaseBill) Post() {
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
			c.TplName = "bills/clusterproductpurchaseBill/clusterproductpurchaseBill.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchOrder Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "bills/clusterproductpurchaseBill/clusterproductpurchaseBill.html"
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

	auth, err := utils.IsAuthorized(utype, "CluPductPurBill")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	data, err := services.GetClusteruser()
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
		d.Name = data[i][1]
		Dis.Fields = append(Dis.Fields, d)
	}
	c.Data["Dis"] = Dis
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

	select_sourcecluster_name := c.Input().Get("select_sourcecluster_name")
	select_beneficiarycluster_name := c.Input().Get("select_beneficiarycluster_name")
	//input_product_name := c.Input().Get("input_product_name")
	//input_status := c.Input().Get("input_status")

	// c.Data["RRN"] = input_rnn_no
	// c.Data["Status"] = input_status
	// c.Data["Mask"] = input_card_mask
	// c.Data["DateRange"] = c.Input().Get("daterange")
	// tmp := strings.Split(c.Input().Get("daterange"), " - ")

	// from := ""
	// to := ""
	// if len(tmp) == 2 {
	// 	from = tmp[0]
	// 	to = tmp[1]
	// }
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
		row, err := db.Db.Query(`select id, 
	source_cluster_name,
	beneficiary_cluster_name, 
	product_name, 
	product_sku,
	number_of_units, 
	delivery_charges, 
	total_price, 
	status,
	create_date,
	cluster_productpurchase_txnid,
	product_price
 from ClusterPurchaseProduct where source_cluster_name=?  AND (source_cluster_name='' OR source_cluster_name like ?) AND (beneficiary_cluster_name='' OR beneficiary_cluster_name like ?) ORDER BY create_date DESC`, cname, select_sourcecluster_name+"%", select_beneficiarycluster_name+"%")

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
			r.Source_cluster_name = data[i][1]
			r.Beneficiary_cluster_name = data[i][2]
			r.Product_name = data[i][3]
			r.Product_sku = data[i][4]
			r.Number_of_units = data[i][5]
			r.Delivery_charges = data[i][6]
			r.Total_price = data[i][7]
			r.Status = data[i][8]
			r.TimeStamp = data[i][9]
			r.Cluster_productpurchase_txnid = data[i][10]
			r.Product_price = data[i][11]

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
		row, err := db.Db.Query(`select id, 
	source_cluster_name,
	beneficiary_cluster_name, 
	product_name, 
	product_sku,
	number_of_units, 
	delivery_charges, 
	total_price, 
	status,
	create_date,
	cluster_productpurchase_txnid,
	product_price
 from ClusterPurchaseProduct  where source_cluster_name=? AND (source_cluster_name='' OR source_cluster_name like ?) AND (beneficiary_cluster_name='' OR beneficiary_cluster_name like ?) ORDER BY create_date DESC`, cname, select_sourcecluster_name+"%", select_beneficiarycluster_name+"%")

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
			r.Source_cluster_name = data[i][1]
			r.Beneficiary_cluster_name = data[i][2]
			r.Product_name = data[i][3]
			r.Product_sku = data[i][4]
			r.Number_of_units = data[i][5]
			r.Delivery_charges = data[i][6]
			r.Total_price = data[i][7]
			r.Status = data[i][8]
			r.TimeStamp = data[i][9]
			r.Cluster_productpurchase_txnid = data[i][10]
			r.Product_price = data[i][11]

			result = append(result, r)
		}
		c.Data["CustomerData"] = result
	} else {
		row, err := db.Db.Query(`select id, 
	source_cluster_name,
	beneficiary_cluster_name, 
	product_name, 
	product_sku,
	number_of_units, 
	delivery_charges, 
	total_price, 
	status,
	create_date,
	cluster_productpurchase_txnid,
	product_price
 from ClusterPurchaseProduct  where (source_cluster_name='' OR source_cluster_name like ?) AND (beneficiary_cluster_name='' OR beneficiary_cluster_name like ?) ORDER BY create_date DESC`, select_sourcecluster_name+"%", select_beneficiarycluster_name+"%")

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
			r.Source_cluster_name = data[i][1]
			r.Beneficiary_cluster_name = data[i][2]
			r.Product_name = data[i][3]
			r.Product_sku = data[i][4]
			r.Number_of_units = data[i][5]
			r.Delivery_charges = data[i][6]
			r.Total_price = data[i][7]
			r.Status = data[i][8]
			r.TimeStamp = data[i][9]
			r.Cluster_productpurchase_txnid = data[i][10]
			r.Product_price = data[i][11]

			result = append(result, r)
		}
		c.Data["CustomerData"] = result
	}

	return
}
