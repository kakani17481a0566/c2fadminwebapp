package updateClusterProduct

import (
	"c2f_admin_webapp/session"

	"runtime/debug"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"
	"errors"
	"fmt"
	"strings"

	"c2f_admin_webapp/model/db"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

type Row struct {
	Id              string
	ClusterName     string
	ProductName     string
	Status          string
	SellingPrice    string
	BuyingPrice     string
	MarketPrice     string
	SupplyCharge    string
	PackagingCharge string
	BankCharge      string
	ClusterDiscount string
}

type UpdateClusterProduct struct {
	// beego.Controller
	middleware.BaseController
}

func (c *UpdateClusterProduct) Get() {
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
			c.TplName = "clusterproducts/updateClusterProduct/updateClusterProduct.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update ClusterProduct Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "clusterproducts/updateClusterProduct/updateClusterProduct.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update ClusterProduct Page Success")
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

	auth, err := utils.IsAuthorized(utype, "ClusterProduct")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	row, err := db.Db.Query(`SELECT Product2Cluster.cluster_product_id,Product2Cluster.cluster_id, Product2Cluster.product_id, 
       Product2Cluster.status, Product2Cluster.create_date,Clusters.cluster_name, 
       Lookup_products.product_name,
     Lookup_products.product_package_qty,
     Lookup_products.product_qty,
     Product2Cluster.cluster_selling_price,
	 Product2Cluster.cluster_buying_price,
	 Product2Cluster.cluster_market_price,
	 Product2Cluster.cluster_supply_charge,
	 Product2Cluster.cluster_packaging_charge,
	 Product2Cluster.cluster_bank_charge,
	 Product2Cluster.cluster_discount,
	 Product2Cluster.delivery_discount
FROM C2F_NEW.Product2Cluster
LEFT JOIN Clusters 
ON C2F_NEW.Product2Cluster.cluster_id = Clusters.id
LEFT JOIN Lookup_products 
ON C2F_NEW.Product2Cluster.product_id = Lookup_products.id WHERE Product2Cluster.cluster_product_id =?  limit 1`, AdminId)
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
		c.Data["ClusterName"] = data[i][5]
		c.Data["ProductName"] = data[i][6]
		c.Data["Status"] = data[i][3]
		c.Data["SellingPrice"] = data[i][9]
		c.Data["BuyingPrice"] = data[i][10]
		c.Data["MarketPrice"] = data[i][11]
		c.Data["SupplyCharge"] = data[i][12]
		c.Data["PackagingCharge"] = data[i][13]
		c.Data["BankCharge"] = data[i][14]
		c.Data["ClusterDiscount"] = data[i][15]
		c.Data["DeliveryDiscount"] = data[i][16]

	}

	return

}

func (c *UpdateClusterProduct) Post() {

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
			c.TplName = "clusterproducts/updateClusterProduct/updateClusterProduct.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update ClusterProduct Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Cluster-Product Updated Successfully"
			c.TplName = "clusterproducts/updateClusterProduct/updateClusterProduct.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update ClusterProduct  Page Success")
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

	auth, err := utils.IsAuthorized(utype, "ClusterProduct")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	// row, err := db.Db.Query(`SELECT Product2Cluster.cluster_product_id,Product2Cluster.cluster_id, Product2Cluster.product_id,
	//       Product2Cluster.status, Product2Cluster.create_date,Clusters.cluster_name,
	//       Lookup_products.product_name,
	//     Lookup_products.product_package_qty,
	//     Lookup_products.product_qty,
	//     Product2Cluster.cluster_selling_price,
	//  Product2Cluster.cluster_buying_price,
	//  Product2Cluster.cluster_market_price,
	//  Product2Cluster.cluster_supply_charge,
	//  Product2Cluster.cluster_packaging_charge,
	//  Product2Cluster.cluster_bank_charge,
	//  Product2Cluster.cluster_discount
	// FROM C2F_NEW.Product2Cluster
	// LEFT JOIN Clusters
	// ON C2F_NEW.Product2Cluster.cluster_id = Clusters.id
	// LEFT JOIN Lookup_products
	// ON C2F_NEW.Product2Cluster.product_id = Lookup_products.id WHERE Product2Cluster.cluster_product_id =?  limit 1`, AdminId)
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
	// 	c.Data["ClusterName"] = data[i][5]
	// 	c.Data["ProductName"] = data[i][6]
	// 	c.Data["Status"] = data[i][3]
	// 	c.Data["SellingPrice"] = data[i][9]
	// 	c.Data["BuyingPrice"] = data[i][10]
	// 	c.Data["MarketPrice"] = data[i][11]
	// 	c.Data["SupplyCharge"] = data[i][12]
	// 	c.Data["PackagingCharge"] = data[i][13]
	// 	c.Data["BankCharge"] = data[i][14]
	// 	c.Data["ClusterDiscount"] = data[i][15]

	// }

	//input_product_status := c.Input().Get("input_product_status")
	input_status := c.Input().Get("input_status")
	input_cluster_selling_price := c.Input().Get("input_cluster_selling_price")
	input_cluster_buying_price := c.Input().Get("input_cluster_buying_price")
	input_cluster_market_price := c.Input().Get("input_cluster_market_price")
	input_cluster_supply_charge := c.Input().Get("input_cluster_supply_charge")
	input_cluster_packaging_charge := c.Input().Get("input_cluster_packaging_charge")
	input_cluster_bank_charge := c.Input().Get("input_cluster_bank_charge")
	input_cluster_discount := c.Input().Get("input_cluster_discount")
	select_cluster_name := c.Input().Get("select_cluster_name")
	select_product_name := c.Input().Get("select_product_name")
	delivery_discount := c.Input().Get("delivery_discount")

	if input_status == "" || input_cluster_selling_price == "" || input_cluster_buying_price == "" || input_cluster_market_price == "" || input_cluster_supply_charge == "" ||
		input_cluster_packaging_charge == "" || input_cluster_bank_charge == "" || input_cluster_discount == "" {
		err = errors.New("* fields  cannot be empty")
		return
	}

	row2, err := db.Db.Query(`select product_qty from Lookup_products WHERE product_name =?`, select_product_name)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	defer sql.Close(row2)
	_, data2, err := sql.Scan(row2)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data2, "\nData len2 - ", len(data2))

	product_qty := data2[0][0]
	fmt.Println("product_qty", product_qty)

	if delivery_discount == "" && product_qty == "kg-bulk" {
		err = errors.New("Delivery discount cannot be empty for kg-bulk product")
		return
	}

	res, err := db.Db.Exec(`UPDATE Product2Cluster SET 
	status = ?,
	cluster_selling_price=?,
	cluster_buying_price=?,
	cluster_market_price=?,
	cluster_supply_charge=?,
	cluster_packaging_charge=?,
	cluster_bank_charge=?,
	cluster_discount=?,
	delivery_discount=?,
	update_date=now()
	WHERE cluster_product_id = ?`,
		input_status,
		input_cluster_selling_price,
		input_cluster_buying_price,
		input_cluster_market_price,
		input_cluster_supply_charge,
		input_cluster_packaging_charge,
		input_cluster_bank_charge,
		input_cluster_discount,
		delivery_discount,
		AdminId)
	if err != nil {
		err = errors.New("Cluster-Product updation failed")
		return
	}

	i, err := res.RowsAffected()
	if err != nil || i == 0 {
		err = errors.New("Cluster-Product updation failed")
		return
	}

	res1, err1 := db.Db.Exec(`UPDATE Inventory SET 
	cluster_selling_price=?,
	cluster_buying_price=?,
	cluster_supply_charge=?,
	cluster_packaging_charge=?,
	cluster_bank_charge=?,
	cluster_discount=?,
	update_date=now()
	WHERE cluster_name=? AND product_name=?`,
		input_cluster_selling_price,
		input_cluster_buying_price,
		input_cluster_supply_charge,
		input_cluster_packaging_charge,
		input_cluster_bank_charge,
		input_cluster_discount,
		select_cluster_name, select_product_name)
	if err1 != nil {
		err1 = errors.New("Cluster-Product updation failed")
		return
	}

	i1, err1 := res1.RowsAffected()
	if err1 != nil || i1 == 0 {
		err1 = errors.New("Cluster-Product updation failed")
		return
	}

	return
}
