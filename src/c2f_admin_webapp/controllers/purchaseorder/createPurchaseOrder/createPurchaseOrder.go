package createPurchaseOrder

import (
	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/model/utils"
	"c2f_admin_webapp/services"
	"c2f_admin_webapp/session"
	"errors"
	"runtime/debug"

	"strconv"
	"strings"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"

	"c2f_admin_webapp/model/util/txnno"

	"github.com/astaxie/beego"
)

//	type Row struct {
//		Sourceclustername      string
//		Beneficiaryclustername string
//		Productname            string
//		Productsku             string
//		Requiredunits          string
//	}
type Display struct {
	Fields []Field
}
type Field struct {
	Id   string
	Name string
	Sku  string
}
type CreatePurchaseOrder struct {
	// beego.Controller
	middleware.BaseController
}

func (c *CreatePurchaseOrder) Get() {
	AdminId := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "AdminId - ", AdminId)

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
			c.TplName = "purchaseorder/createPurchaseOrder/createPurchaseOrder.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "purchaseorder/createPurchaseOrder/createPurchaseOrder.html"
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

	maptoken := sess.Get("token")
	c.Data["maptoken"] = maptoken
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "maptoken", maptoken)

	auth, err := utils.IsAuthorized(utype, "PurchaseOrder")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	row, err := db.Db.Query(`select id,cluster_name, product_name, 
        product_sku, 
        inventory_available ,
        inventory_delivered,
        cluster_selling_price,
        product_package_qty,
        product_qty
    from Inventory where id=? limit 1`, AdminId)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	defer sql.Close(row)
	_, data3, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data3, "\nData len - ", len(data3))
	if len(data3) <= 0 {
		err = errors.New("User data not found")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", data3)
	for i := range data3 {
		c.Data["ID"] = data3[i][0]
		c.Data["Clustername"] = data3[i][1]
		c.Data["Productname"] = data3[i][2]
		c.Data["Productsku"] = data3[i][3]
		c.Data["Inventoryavailable"] = data3[i][4]
		c.Data["Inventorydelivered"] = data3[i][5]
		c.Data["Clustersellingprice"] = data3[i][6]
		c.Data["ProductPackageQty"] = data3[i][7]
		c.Data["ProductQty"] = data3[i][8]
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Source Cluster Name - ", data3[0][1])

	// Get all clusters

	data, err := services.GetClusters3(data3[0][1])
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Clusters not found")
		return
	}
	if len(data) <= 0 {
		err = errors.New("cluster name not found")
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

	return
}
func (c *CreatePurchaseOrder) Post() {
	msg11 := ""
	AdminId := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "AdminId - ", AdminId)

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
			c.TplName = "purchaseorder/createPurchaseOrder/createPurchaseOrder.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = msg11
			c.TplName = "purchaseorder/createPurchaseOrder/createPurchaseOrder.html"
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

	auth, err := utils.IsAuthorized(utype, "PurchaseOrder")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "adminid - ", AdminId)

	input_source_clustername := c.Input().Get("input_source_clustername")
	input_product_name := c.Input().Get("input_product_name")

	input_required_units := c.Input().Get("input_required_units")
	input_available_units := c.Input().Get("input_available_units")
	input_id := c.Input().Get("input_id")
	input_cluster_price := c.Input().Get("input_cluster_price")
	input_status := c.Input().Get("input_status")
	//input_delivery_charge := c.Input().Get("input_delivery_charge")
	cluster_productpurchase_txnid := txnno.Generate13Digit()

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_status - ", input_status)

	// Get product sku of source cluster
	input_product_sku := c.Input().Get("input_product_sku")

	// Get the clustername of beneficary
	input_beneficiary_clustername := c.Input().Get("input_beneficiary_clustername")

	// Check in inventory if the product sku is there for beneficary cluster

	row5, err := db.Db.Query(`SELECT count(*) FROM Inventory where product_sku=? and cluster_name=?`, input_product_sku, input_beneficiary_clustername)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get beneficary inventory data")
		return
	}
	defer sql.Close(row5)

	_, data5, err := sql.Scan(row5)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get beneficary inventory data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data5, "\nData len - ", len(data5))
	if len(data5) <= 0 {
		err = errors.New("Unable to get inventory data")
		return
	}

	// checking with count if it is 0 add new row inventory , if >0 update exisisting inventory

	count, err := strconv.ParseFloat(data5[0][0], 64)

	if count <= 0 {

		// If ProductSKU NOT there in beneficary cluster

		row, err := db.Db.Query(`select id,cluster_name, 
         product_name, 
        product_sku, 
        inventory_available ,
        inventory_delivered,
        cluster_selling_price,
        product_package_qty,
        product_qty
    from Inventory where id=? limit 1`, input_id)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get user data")
			msg11 = "Unable to get user data"
			return
		}
		defer sql.Close(row)
		_, data3, err := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get user data")
			msg11 = "Unable to get user data"
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data3, "\nData len - ", len(data3))
		if len(data3) <= 0 {
			err = errors.New("User data not found")
			msg11 = "User data not found"
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", data3)
		for i := range data3 {
			c.Data["ID"] = data3[i][0]
			c.Data["Clustername"] = data3[i][1]
			c.Data["Productname"] = data3[i][2]
			c.Data["Productsku"] = data3[i][3]
			c.Data["Inventoryavailable"] = data3[i][4]
			c.Data["Inventorydelivered"] = data3[i][5]
			c.Data["Clustersellingprice"] = data3[i][6]
			c.Data["ProductPackageQty"] = data3[i][7]
			c.Data["ProductQty"] = data3[i][8]

		}

		input_delivered_units := data3[0][5]
		//input_oder_id := c.Input().Get("input_oder_id")

		//input_courier_details := c.Input().Get("input_courier_details")
		input_product_package_qty := c.Input().Get("input_product_package_qty")
		input_product_qty := c.Input().Get("input_product_qty")

		aunits, err := strconv.ParseFloat(input_available_units, 64)
		runits, err := strconv.ParseFloat(input_required_units, 64)
		dunits, err := strconv.ParseFloat(input_delivered_units, 64)
		//dcharge, err := strconv.ParseFloat(input_delivery_charge, 64)

		cprice, err := strconv.ParseFloat(input_cluster_price, 64)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "cprice - ", cprice)

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "aunits - ", aunits)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "runits - ", runits)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "dunits - ", dunits)
		//log.Println(beego.AppConfig.String("loglevel"), "Debug", "dcharge - ", dcharge)

		TA := aunits - runits
		TA1 := strconv.FormatFloat(TA, 'f', 0, 32)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "TA1 - ", TA1)

		TD := dunits + runits
		TD1 := strconv.FormatFloat(TD, 'f', 0, 32)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "TD1 - ", TD1)

		TotalPrice := (cprice * runits)
		TotalPrice1 := strconv.FormatFloat(TotalPrice, 'f', 0, 32)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "TotalPrice1 - ", TotalPrice1)

		c.Data["TotalPrice1"] = TotalPrice1

		if runits > aunits {

			//log.Println(beego.AppConfig.String("loglevel"), "Debug", "inside - ")

			err = errors.New("Insufficent Available Units")
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			msg11 = "Insufficent Available Units"

			return
		}

		if runits == 0 {

			//log.Println(beego.AppConfig.String("loglevel"), "Debug", "inside - ")

			err = errors.New("Required Units must be greater than Zeroo")
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			msg11 = "Required Units must be greater than Zero"

			return
		}

		if input_product_qty == "kg-bulk" && runits < 10 {
			err = errors.New("Required Units must be at least 10 kg for kg-bulk product")
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			msg11 = "Required Units must be at least 10 kg for kg-bulk product"
			return
		}

		result, err := db.Db.Exec(`INSERT INTO ClusterPurchaseProduct 
	   (cluster_productpurchase_txnid,
	    source_cluster_name, 
        beneficiary_cluster_name,
         product_name, 
        product_sku,  
        number_of_units,
        total_price,
        status,
        product_price,
        product_package_qty,
        sku_type,
        product_qty,
        create_date)
        VALUES (?,?,?,?,?,?,?,?,?,?,?,?,now())`,
			cluster_productpurchase_txnid,
			input_source_clustername,
			input_beneficiary_clustername,
			input_product_name,
			input_product_sku,
			input_required_units,
			TotalPrice1,
			input_status,
			input_cluster_price,
			input_product_package_qty,
			"OTHER",
			input_product_qty)
		if err != nil {
			err = errors.New("Cluster Purchase Product creation failed")
			msg11 = "Cluster Purchase Product creation failed"
			return
		}
		i, err := result.RowsAffected()
		if err != nil || i == 0 {
			err = errors.New("Cluster Purchase Product creation failed")
			msg11 = "Cluster Purchase Product creation failed"
			return
		}

	} else {

		// For exsisting product sku in beneficary cluster

		row, err := db.Db.Query(`select id,cluster_name, 
         product_name, 
        product_sku, 
        inventory_available ,
        inventory_delivered,
        cluster_selling_price,
        product_package_qty,
        product_qty
    from Inventory where id=? limit 1`, input_id)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get user data")
			msg11 = "Unable to get user data"
			return
		}
		defer sql.Close(row)
		_, data3, err := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get user data")
			msg11 = "Unable to get user data"
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data3, "\nData len - ", len(data3))
		if len(data3) <= 0 {
			err = errors.New("User data not found")
			msg11 = "User data not found"
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", data3)
		for i := range data3 {
			c.Data["ID"] = data3[i][0]
			c.Data["Clustername"] = data3[i][1]
			c.Data["Productname"] = data3[i][2]
			c.Data["Productsku"] = data3[i][3]
			c.Data["Inventoryavailable"] = data3[i][4]
			c.Data["Inventorydelivered"] = data3[i][5]
			c.Data["Clustersellingprice"] = data3[i][6]
			c.Data["ProductPackageQty"] = data3[i][7]
			c.Data["ProductQty"] = data3[i][8]

		}

		input_delivered_units := data3[0][5]
		//input_oder_id := c.Input().Get("input_oder_id")

		//input_courier_details := c.Input().Get("input_courier_details")
		input_product_package_qty := c.Input().Get("input_product_package_qty")
		input_product_qty := c.Input().Get("input_product_qty")

		aunits, err := strconv.ParseFloat(input_available_units, 64)
		runits, err := strconv.ParseFloat(input_required_units, 64)
		dunits, err := strconv.ParseFloat(input_delivered_units, 64)
		//dcharge, err := strconv.ParseFloat(input_delivery_charge, 64)

		cprice, err := strconv.ParseFloat(input_cluster_price, 64)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "cprice - ", cprice)

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "aunits - ", aunits)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "runits - ", runits)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "dunits - ", dunits)
		//log.Println(beego.AppConfig.String("loglevel"), "Debug", "dcharge - ", dcharge)

		TA := aunits - runits
		TA1 := strconv.FormatFloat(TA, 'f', 0, 32)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "TA1 - ", TA1)

		TD := dunits + runits
		TD1 := strconv.FormatFloat(TD, 'f', 0, 32)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "TD1 - ", TD1)

		TotalPrice := (cprice * runits)
		TotalPrice1 := strconv.FormatFloat(TotalPrice, 'f', 0, 32)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "TotalPrice1 - ", TotalPrice1)

		c.Data["TotalPrice1"] = TotalPrice1

		if runits > aunits {

			//log.Println(beego.AppConfig.String("loglevel"), "Debug", "inside - ")

			err = errors.New("Insufficent Available Units")
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			msg11 = "Insufficent Available Units"

			return
		}

		if runits == 0 {

			//log.Println(beego.AppConfig.String("loglevel"), "Debug", "inside - ")

			err = errors.New("Required Units must be greater than Zero")
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			msg11 = "Required Units must be greater than Zero"

			return
		}

		if input_product_qty == "kg-bulk" && runits < 10 {
			err = errors.New("Required Units must be at least 10 kg for kg-bulk product")
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			msg11 = "Required Units must be at least 10 kg for kg-bulk product"
			return
		}

		result, err := db.Db.Exec(`INSERT INTO ClusterPurchaseProduct 
	   (cluster_productpurchase_txnid,
	    source_cluster_name, 
        beneficiary_cluster_name,
         product_name, 
        product_sku,  
        number_of_units,
        total_price,
        status,
        product_price,
        product_package_qty,
        sku_type,
        product_qty,
        create_date)
        VALUES (?,?,?,?,?,?,?,?,?,?,?,?,now())`,
			cluster_productpurchase_txnid,
			input_source_clustername,
			input_beneficiary_clustername,
			input_product_name,
			input_product_sku,
			input_required_units,
			TotalPrice1,
			input_status,
			input_cluster_price,
			input_product_package_qty,
			"SOURCE",
			input_product_qty)
		if err != nil {
			err = errors.New("Cluster Purchase Product creation failed")
			msg11 = "Cluster Purchase Product creation failed"
			return
		}
		i, err := result.RowsAffected()
		if err != nil || i == 0 {
			err = errors.New("Cluster Purchase Product creation failed")
			msg11 = "Cluster Purchase Product creation failed"
			return
		}

	}
	msg11 = "Create Cluster Purchase Product Order Successfully"

	// res2, err := db.Db.Exec(`UPDATE Inventory SET inventory_delivered = ? ,inventory_available=?
	//                              WHERE id = ?`, TD1, TA1, input_id)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Inventory Get Failed")
	// 	return
	// }
	// i2, err := res2.RowsAffected()
	// if err != nil || i2 == 0 {
	// 	err = errors.New(" CustomerOrderTxn Delivery updation failed")
	// 	return
	// }
	return
}
