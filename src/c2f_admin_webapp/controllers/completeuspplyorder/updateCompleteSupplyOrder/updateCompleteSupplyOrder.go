package updateCompleteSupplyOrder

import (
	"c2f_admin_webapp/session"

	"runtime/debug"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"
	"errors"
	"strconv"
	"strings"

	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/services"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

type Row struct {
	Id              string
	CategoryName    string
	SubCategoryName string
	ProductName     string
	ProductDesc     string
	Productstatus   string
	CreatedBy       string
}

type UpdateCompleteSupplyOrder struct {
	// beego.Controller
	middleware.BaseController
}

func (c *UpdateCompleteSupplyOrder) Get() {
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

			c.TplName = "completesupplyorder/updateCompleteSupplyOrder/updateCompleteSupplyOrder.html"

			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "completesupplyorder/updateCompleteSupplyOrder/updateCompleteSupplyOrder.html"

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

	auth, err := utils.IsAuthorized(utype, "CompleteSupplyOrder")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	row, err := db.Db.Query(`select id,supplyordertxn_id,product_name,product_sku,product_qty,product_package_qty,cluster_buying_price,product_supply_units,product_supply_total_price,status,category_name,subcategory_name,cluster_name,product_id,producer_mob,producer_mob from SupplyOrderTxn WHERE id =?  limit 1`, AdminId)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get supply order transaction data")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get supply order transaction data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data : supply order transaction - ", data, "\nData len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("supply order transaction data not found")
		return
	}

	// log.Println(beego.AppConfig.String("loglevel"), "supply order transaction", data)

	for i := range data {

		c.Data["Id"] = data[i][0]
		c.Data["SupplyTxnID"] = data[i][1]
		c.Data["ProductName"] = data[i][2]
		c.Data["ProductSKU"] = data[i][3]
		c.Data["ProductQTY"] = data[i][4]
		c.Data["ProductPackageQTY"] = data[i][5]
		c.Data["ClusterBuyPrice"] = data[i][6]
		c.Data["SupplyUnits"] = data[i][7]
		c.Data["TotalPrice"] = data[i][8]
		c.Data["Status"] = data[i][9]
		c.Data["CategoryName"] = data[i][10]
		c.Data["SubCategoryName"] = data[i][11]
		c.Data["ClusterName"] = data[i][12]
		c.Data["PID"] = data[i][13]
		c.Data["ProducerMob"] = data[i][14]

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Product id  - ", data[i][13])
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Cluster name  - ", data[i][12])
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "ProducerMob  - ", data[i][14])

		row1, err := db.Db.Query(`select id,product_quality_mesure,product_shelf_life, product_package_type, product_origin,krishivala_selling_price from Lookup_products WHERE product_sku =?`, data[i][3])
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get product data")
			return
		}
		defer sql.Close(row1)
		_, data1, err := sql.Scan(row1)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get product data")
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data product - ", data1, "\nData len - ", len(data1))
		if len(data1) <= 0 {
			err = errors.New("product data not found")
			return
		}
		for j := range data1 {
			c.Data["ProductQualityMesure"] = data1[j][1]
			c.Data["Productshelflife"] = data1[j][2]
			c.Data["ProductPackageType"] = data1[j][3]
			c.Data["ProductOrigin"] = data1[j][4]
			c.Data["KrishivalaSellPrice"] = data1[j][5]

		}

		row2, err := db.Db.Query(`select id,producer_selling_price,product_seeling_units,producer_id from SupplyForecastTxn WHERE product_sku =? AND cluster_name=? AND producer_mob=?`, data[i][3], data[i][12], data[i][14])
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get SupplyForecastTxn data")
			return
		}
		defer sql.Close(row2)
		_, data2, err := sql.Scan(row2)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get SupplyForecastTxn data")
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data SupplyForecastTxn - ", data2, "\nData len - ", len(data2))
		if len(data2) <= 0 {
			err = errors.New("SupplyForecastTxn data not found")
			return
		}
		for j := range data2 {
			c.Data["ProducerSellPrice"] = data2[j][1]
			c.Data["ProductSelingUnits"] = data2[j][2]
			c.Data["Producerid"] = data2[j][3]

		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Producerid  - ", data2[0][3])

		row22, err22 := db.Db.Query(`select id from  Clusters WHERE cluster_name=?`, data[i][12])
		if err22 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get cluster id")
			return
		}
		defer sql.Close(row22)
		_, data22, err := sql.Scan(row22)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get cluster id")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Cluster Id  - ", data22[0][0])

		row3, err3 := db.Db.Query(`select product_id,cluster_selling_price,cluster_supply_charge,cluster_packaging_charge,cluster_bank_charge,cluster_discount from Product2Cluster where product_id = ? and cluster_id = ?`, data[i][13], data22[0][0])
		if err3 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err3)
			err3 = errors.New("Unable to get Product2Cluster data")
			return
		}
		defer sql.Close(row3)
		_, data3, err := sql.Scan(row3)
		if err3 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err3)
			err3 = errors.New("Unable to get Product2Cluster data")
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data Product2Cluster- ", data3, "Datalen - ", len(data3))
		if len(data3) <= 0 {
			err3 = errors.New("cluster product  data not found")
			return
		}
		for k := range data3 {
			c.Data["ClusterSellPrice"] = data3[k][1]
			c.Data["ClusterSupplyCharge"] = data3[k][2]
			c.Data["ClusterPackagingCharge"] = data3[k][3]
			c.Data["ClusterBankCharge"] = data3[k][4]
			c.Data["ClusterDiscount"] = data3[k][5]

		}

		row4, err := db.Db.Query(`select product_id,producer_CGST_charge,producer_SGST_charge,producer_CESS_charge,producer_OtherTAX_charge from Product2Producer2Cluster WHERE product_id =? AND cluster_id=? AND producer_id=?`, data[i][13], data22[0][0], data2[0][3])
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get Product2Producer2Cluster data")
			return
		}
		defer sql.Close(row4)
		_, data4, err := sql.Scan(row4)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get Product2Producer2Cluster data")
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data Product2Producer2Cluster - ", data4, "\nData len - ", len(data4))
		if len(data4) <= 0 {
			err = errors.New("cluster selling data not found")
			return
		}
		for l := range data4 {
			c.Data["Producer_CGST_charge"] = data4[l][1]
			c.Data["Producer_SGST_charge"] = data4[l][2]
			c.Data["Producer_CESS_charge"] = data4[l][3]
			c.Data["Producer_OtherTAX_charge"] = data4[l][4]

		}
	}
	return

}

func (c *UpdateCompleteSupplyOrder) Post() {

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
			c.TplName = "completesupplyorder/updateCompleteSupplyOrder/updateCompleteSupplyOrder.html"

			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Complete Supply Order Updated Successfully"
			c.TplName = "completesupplyorder/updateCompleteSupplyOrder/updateCompleteSupplyOrder.html"

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

	id := sess.Get("uid").(string)
	c.Data["id"] = id
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

	auth, err := utils.IsAuthorized(utype, "CompleteSupplyOrder")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	input_product_id := c.Input().Get("input_product_id")

	input_supplytxn_id := c.Input().Get("input_supplytxn_id")
	input_category_name := c.Input().Get("input_category_name")
	input_subcategory_name := c.Input().Get("input_subcategory_name")
	input_product_name := c.Input().Get("input_product_name")

	input_product_sku := c.Input().Get("input_product_sku")
	input_productQTY := c.Input().Get("input_productQTY")
	input_productPackageQTY := c.Input().Get("input_productPackageQTY")

	input_product_qulaitymesure := c.Input().Get("input_product_qulaitymesure")
	input_product_shelf_life := c.Input().Get("input_product_shelf_life")
	input_product_origin := c.Input().Get("input_product_origin")
	input_product_package_type := c.Input().Get("input_product_package_type")

	input_krishivala_selling_price := c.Input().Get("input_krishivala_selling_price")
	input_cluster_selling_price := c.Input().Get("input_cluster_selling_price")
	input_cluster_buy_price := c.Input().Get("input_cluster_buy_price")
	input_producer_selling_price := c.Input().Get("input_producer_selling_price")

	input_supply_units := c.Input().Get("input_supply_units")
	//	input_supplied_total_price := c.Input().Get("input_supplied_total_price")
	input_cluster_name := c.Input().Get("input_cluster_name")
	//input_receivedunits := c.Input().Get("input_receivedunits")
	input_status := c.Input().Get("input_status")

	cluster_supply_charge := c.Input().Get("cluster_supply_charge")
	cluster_packaging_charge := c.Input().Get("cluster_packaging_charge")
	cluster_bank_charge := c.Input().Get("cluster_bank_charge")
	cluster_discount := c.Input().Get("cluster_discount")

	producer_CGST_charge := c.Input().Get("producer_CGST_charge")
	producer_SGST_charge := c.Input().Get("producer_SGST_charge")
	producer_CESS_charge := c.Input().Get("producer_CESS_charge")
	producer_OtherTAX_charge := c.Input().Get("producer_OtherTAX_charge")
	input_damagedunits := c.Input().Get("input_damagedunits")

	input_remarks := c.Input().Get("input_remarks")

	if input_supply_units == "" || input_status == "" {
		err = errors.New(" * fields cannot be empty ")
		return
	}

	input_damagedunits1, err := strconv.ParseFloat(input_damagedunits, 64)
	input_supply_units1, err := strconv.ParseFloat(input_supply_units, 64)

	if input_damagedunits1 > input_supply_units1 {

		err = errors.New("Damage Units Not More than Supply Unit")
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_damagedunits1 - ", input_damagedunits1)

	var input_receivedunits string

	//	log.Println(beego.AppConfig.String("loglevel"), "Debug", "totalunits2 - ", totalunits2)

	if input_remarks == "RECEIVED_GOOD" {

		input_receivedunits = input_supply_units
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_receivedunits - ", input_receivedunits)

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_supply_units - ", input_supply_units)

	} else if input_remarks == "SUPPLY_DAMAGED" {

		input_receivedunits1 := input_supply_units1 - input_damagedunits1

		input_receivedunits = strconv.FormatFloat(input_receivedunits1, 'f', 0, 32)

		if input_damagedunits == "" {
			err = errors.New(" Damagedunits field cannot be empty ")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_receivedunits - ", input_receivedunits)
	} else {
		input_receivedunits = "0"
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_receivedunits - ", input_receivedunits)
	}

	//c.Data["TotalUnits"] = totalunits

	bp, err := strconv.ParseFloat(input_cluster_buy_price, 64)
	u, err := strconv.ParseFloat(input_receivedunits, 64)

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "u - ", u)
	BP := bp * u
	totalprice := strconv.FormatFloat(BP, 'f', 0, 32)

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "BP - ", BP)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "bp - ", bp)

	updated_total_supply_price := totalprice

	scluster_supply_charge, err := strconv.ParseFloat(cluster_supply_charge, 64)
	scluster_packaging_charge, err := strconv.ParseFloat(cluster_packaging_charge, 64)
	scluster_bank_charge, err := strconv.ParseFloat(cluster_bank_charge, 64)
	scluster_discount, err := strconv.ParseFloat(cluster_discount, 64)

	sproducer_CGST_charge, err := strconv.ParseFloat(producer_CGST_charge, 64)
	sproducer_SGST_charge, err := strconv.ParseFloat(producer_SGST_charge, 64)
	sproducer_CESS_charge, err := strconv.ParseFloat(producer_CESS_charge, 64)
	sproducer_OtherTAX_charge, err := strconv.ParseFloat(producer_OtherTAX_charge, 64)

	charges := scluster_supply_charge + scluster_packaging_charge + scluster_bank_charge
	taxs := sproducer_CGST_charge + sproducer_SGST_charge + sproducer_CESS_charge + sproducer_OtherTAX_charge

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "charges - ", charges)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "taxs - ", taxs)

	TBP := BP + (charges) - scluster_discount + (taxs)

	row3, err3 := db.Db.Query(`select count(*) from Oder2SupplyExecutive WHERE supply_order_id = ?`, input_supplytxn_id)

	if err3 != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get count in Oder2SupplyExecutive")
		return
	}
	defer sql.Close(row3)
	_, data3, err := sql.Scan(row3)
	if err3 != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get count in Oder2SupplyExecutive")
		return
	}

	// log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	countlen := data3[0][0]

	var SE_charges string

	if countlen != "0" {
		res22, err := db.Db.Query(`select supply_executive_charges from  Oder2SupplyExecutive  WHERE supply_order_id = ?`, input_supplytxn_id)
		if err != nil {
			err = errors.New("supply_executive_charges fetch  failed")

			return
		}
		defer sql.Close(res22)
		_, data22, err := sql.Scan(res22)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("supply_executive_charges fetch  failed")
			return
		}
		SE_charges = data22[0][0]
		supply_exec_charges, err := strconv.ParseFloat(SE_charges, 64)
		TBP = TBP - supply_exec_charges

	}

	tbp := strconv.FormatFloat(TBP, 'f', 0, 32)

	//  TotalBuying_price = cluster.buyingprice  + (cluster.supplycarge + cluster.packagingcharge + cluster.bankcharge ) - cluster.discount  + (producer.gcst + producer.scst + producer.C29 + producer.othertax)
	TotalBuying_price := tbp

	// SELECT EXISTS(SELECT 1 FROM table_1 WHERE id = 1);

	row6, err := db.Db.Query(`select * from Inventory where product_sku = ? and cluster_name=? `, input_product_sku, input_cluster_name)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	defer sql.Close(row6)
	_, data6, err := sql.Scan(row6)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data6, "\nData len - ", len(data6))
	if len(data6) <= 0 {

		err = Inventoryupdate(input_product_id,
			input_cluster_name,
			input_category_name,
			input_subcategory_name,
			input_product_name,
			input_product_sku,
			input_productQTY,
			input_receivedunits,
			input_status,
			input_product_qulaitymesure,
			input_product_shelf_life,
			input_product_package_type,
			input_productPackageQTY,
			input_product_origin,
			input_krishivala_selling_price,
			input_cluster_buy_price,
			input_cluster_selling_price,
			input_producer_selling_price,
			input_receivedunits,
			input_receivedunits,
			"",
			"",
			input_damagedunits,
			"",
			id,
			cluster_supply_charge,
			cluster_packaging_charge,
			cluster_bank_charge,
			cluster_discount)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Inventory creation failed")
			return
		}

		res, err := db.Db.Exec(`UPDATE SupplyOrderTxn SET 
	 cluster_supply_charge=?,
     cluster_packaging_charge=?,
     cluster_bank_charge=?,
     cluster_discount_charge=?,
     producer_CGST_charge=?,
     producer_SGST_charge=?,
     producer_CESS_charge=?,
     producer_OtherTAX_charge=?,
     supply_received_units =?,
     updated_total_supply_price =?,
     status = ?,
    total_buy_price=?,
    remarks=?
	WHERE id = ?`, cluster_supply_charge, cluster_packaging_charge, cluster_bank_charge, cluster_discount, producer_CGST_charge, producer_SGST_charge, producer_CESS_charge, producer_OtherTAX_charge, input_receivedunits, updated_total_supply_price, input_status, TotalBuying_price, input_remarks, AdminId)
		if err != nil {
			err = errors.New("Order updation failed")
			return
		}

		i, err := res.RowsAffected()
		if err != nil || i == 0 {
			err = errors.New("Order updation failed")
			return
		}

		if countlen != "0" {
			res2, err := db.Db.Exec(`UPDATE Oder2SupplyExecutive SET status = ? WHERE supply_order_id = ?`, input_status, input_supplytxn_id)
			if err != nil {
				err = errors.New("Oder2SupplyExecutive Status updation failed")

				return
			}

			i2, err := res2.RowsAffected()
			if err != nil || i2 == 0 {
				err = errors.New(" Oder2SupplyExecutive Status updation failed")

				return
			}

		}

		// row, err := db.Db.Query(`select id,supplyordertxn_id,product_name,product_sku,product_qty,product_package_qty,cluster_buying_price,product_supply_units,product_supply_total_price,status,category_name,subcategory_name,cluster_name,product_id from SupplyOrderTxn WHERE id =?  limit 1`, AdminId)
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
		// 	c.Data["SupplyTxnID"] = data[i][1]
		// 	c.Data["ProductName"] = data[i][2]
		// 	c.Data["ProductSKU"] = data[i][3]
		// 	c.Data["ProductQTY"] = data[i][4]
		// 	c.Data["ProductPackageQTY"] = data[i][5]
		// 	c.Data["ClusterBuyPrice"] = data[i][6]
		// 	c.Data["SupplyUnits"] = data[i][7]
		// 	c.Data["TotalPrice"] = data[i][8]
		// 	c.Data["Status"] = data[i][9]
		// 	c.Data["CategoryName"] = data[i][10]
		// 	c.Data["SubCategoryName"] = data[i][11]
		// 	c.Data["ClusterName"] = data[i][12]
		// 	c.Data["PID"] = data[i][13]

		// 	log.Println(beego.AppConfig.String("loglevel"), "Debug", "data[i][13] - ", data[i][13])

		// 	row1, err := db.Db.Query(`select id,product_quality_mesure,product_shelf_life, product_package_type, product_origin,krishivala_selling_price from Lookup_products WHERE product_sku =?`, data[i][3])
		// 	if err != nil {
		// 		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		// 		err = errors.New("Unable to get user data")
		// 		return
		// 	}
		// 	defer sql.Close(row1)
		// 	_, data1, err := sql.Scan(row1)
		// 	if err != nil {
		// 		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		// 		err = errors.New("Unable to get user data")
		// 		return
		// 	}
		// 	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len2 - ", len(data1))
		// 	if len(data1) <= 0 {
		// 		err = errors.New("User data not found")
		// 		return
		// 	}
		// 	for j := range data1 {
		// 		c.Data["ProductQualityMesure"] = data1[j][1]
		// 		c.Data["Productshelflife"] = data1[j][2]
		// 		c.Data["ProductPackageType"] = data1[j][3]
		// 		c.Data["ProductOrigin"] = data1[j][4]
		// 		c.Data["KrishivalaSellPrice"] = data1[j][5]

		// 	}

		// 	row2, err := db.Db.Query(`select id,producer_selling_price from SupplyForecastTxn WHERE product_sku =?`, data[i][3])
		// 	if err != nil {
		// 		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		// 		err = errors.New("Unable to get user data")
		// 		return
		// 	}
		// 	defer sql.Close(row2)
		// 	_, data2, err := sql.Scan(row2)
		// 	if err != nil {
		// 		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		// 		err = errors.New("Unable to get user data")
		// 		return
		// 	}
		// 	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len2 - ", len(data2))
		// 	if len(data2) <= 0 {
		// 		err = errors.New("User data not found")
		// 		return
		// 	}
		// 	for j := range data2 {
		// 		c.Data["ProducerSellPrice"] = data2[j][1]

		// 	}

		// 	row3, err := db.Db.Query(`select product_id,cluster_selling_price,cluster_supply_charge,cluster_packaging_charge,cluster_bank_charge,cluster_discount from Product2Cluster WHERE product_id =?`, data[i][13])
		// 	if err != nil {
		// 		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		// 		err = errors.New("Unable to get user data")
		// 		return
		// 	}
		// 	defer sql.Close(row3)
		// 	_, data3, err := sql.Scan(row3)
		// 	if err != nil {
		// 		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		// 		err = errors.New("Unable to get user data")
		// 		return
		// 	}
		// 	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len2 - ", len(data3))
		// 	if len(data3) <= 0 {
		// 		err = errors.New("cluster selling data not found")
		// 		return
		// 	}
		// 	for k := range data3 {
		// 		c.Data["ClusterSellPrice"] = data3[k][1]
		// 		c.Data["ClusterSupplyCharge"] = data3[k][2]
		// 		c.Data["ClusterPackagingCharge"] = data3[k][3]
		// 		c.Data["ClusterBankCharge"] = data3[k][4]
		// 		c.Data["ClusterDiscount"] = data3[k][5]

		// 	}
		// }
		return
	} else {

		//	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

		dunits, aunits, runits, psku, cluster, damageunits, err := services.Getinventory(input_product_sku, input_cluster_name)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Inventory Get Failed")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "dunits - ", dunits)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "aunits - ", aunits)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "runits - ", runits)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "psku - ", psku)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "cluster - ", cluster)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "damageunits - ", damageunits)

		err = Inventoryupdate(input_product_id,
			input_cluster_name,
			input_category_name,
			input_subcategory_name,
			input_product_name,
			input_product_sku,
			input_productQTY,
			input_receivedunits,
			input_status,
			input_product_qulaitymesure,
			input_product_shelf_life,
			input_product_package_type,
			input_productPackageQTY,
			input_product_origin,
			input_krishivala_selling_price,
			input_cluster_buy_price,
			input_cluster_selling_price,
			input_producer_selling_price,
			runits,
			aunits,
			psku,
			cluster,
			input_damagedunits,
			damageunits,
			id,
			cluster_supply_charge,
			cluster_packaging_charge,
			cluster_bank_charge,
			cluster_discount)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Inventory creation failed")
			return
		}

		res, err := db.Db.Exec(`UPDATE SupplyOrderTxn SET 
	 cluster_supply_charge=?,
     cluster_packaging_charge=?,
     cluster_bank_charge=?,
     cluster_discount_charge=?,
     producer_CGST_charge=?,
     producer_SGST_charge=?,
     producer_CESS_charge=?,
     producer_OtherTAX_charge=?,
     supply_received_units =?,
     updated_total_supply_price =?,
     status = ?,
    total_buy_price=?,
    remarks=?
	WHERE id = ?`, cluster_supply_charge, cluster_packaging_charge, cluster_bank_charge, cluster_discount, producer_CGST_charge, producer_SGST_charge, producer_CESS_charge, producer_OtherTAX_charge, input_receivedunits, updated_total_supply_price, input_status, TotalBuying_price, input_remarks, AdminId)
		if err != nil {
			err = errors.New("Order updation failed")
			return
		}

		i, err := res.RowsAffected()
		if err != nil || i == 0 {
			err = errors.New("Order updation failed")
			return
		}

		if countlen != "0" {
			res2, err := db.Db.Exec(`UPDATE Oder2SupplyExecutive SET status = ? WHERE supply_order_id = ?`, input_status, input_supplytxn_id)
			if err != nil {
				err = errors.New("Oder2SupplyExecutive Status updation failed")

				return
			}

			i2, err := res2.RowsAffected()
			if err != nil || i2 == 0 {
				err = errors.New(" Oder2SupplyExecutive Status updation failed")

				return
			}

		}

		// row, err := db.Db.Query(`select id,supplyordertxn_id,product_name,product_sku,product_qty,product_package_qty,cluster_buying_price,product_supply_units,product_supply_total_price,status,category_name,subcategory_name,cluster_name,product_id from SupplyOrderTxn WHERE id =?  limit 1`, AdminId)
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
		// 	c.Data["SupplyTxnID"] = data[i][1]
		// 	c.Data["ProductName"] = data[i][2]
		// 	c.Data["ProductSKU"] = data[i][3]
		// 	c.Data["ProductQTY"] = data[i][4]
		// 	c.Data["ProductPackageQTY"] = data[i][5]
		// 	c.Data["ClusterBuyPrice"] = data[i][6]
		// 	c.Data["SupplyUnits"] = data[i][7]
		// 	c.Data["TotalPrice"] = data[i][8]
		// 	c.Data["Status"] = data[i][9]
		// 	c.Data["CategoryName"] = data[i][10]
		// 	c.Data["SubCategoryName"] = data[i][11]
		// 	c.Data["ClusterName"] = data[i][12]
		// 	c.Data["PID"] = data[i][13]

		// 	log.Println(beego.AppConfig.String("loglevel"), "Debug", "data[i][13] - ", data[i][13])

		// 	row1, err := db.Db.Query(`select id,product_quality_mesure,product_shelf_life, product_package_type, product_origin,krishivala_selling_price from Lookup_products WHERE product_sku =?`, data[i][3])
		// 	if err != nil {
		// 		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		// 		err = errors.New("Unable to get user data")
		// 		return
		// 	}
		// 	defer sql.Close(row1)
		// 	_, data1, err := sql.Scan(row1)
		// 	if err != nil {
		// 		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		// 		err = errors.New("Unable to get user data")
		// 		return
		// 	}
		// 	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len2 - ", len(data1))
		// 	if len(data1) <= 0 {
		// 		err = errors.New("User data not found")
		// 		return
		// 	}
		// 	for j := range data1 {
		// 		c.Data["ProductQualityMesure"] = data1[j][1]
		// 		c.Data["Productshelflife"] = data1[j][2]
		// 		c.Data["ProductPackageType"] = data1[j][3]
		// 		c.Data["ProductOrigin"] = data1[j][4]
		// 		c.Data["KrishivalaSellPrice"] = data1[j][5]

		// 	}

		// 	row2, err := db.Db.Query(`select id,producer_selling_price from SupplyForecastTxn WHERE product_sku =?`, data[i][3])
		// 	if err != nil {
		// 		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		// 		err = errors.New("Unable to get user data")
		// 		return
		// 	}
		// 	defer sql.Close(row2)
		// 	_, data2, err := sql.Scan(row2)
		// 	if err != nil {
		// 		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		// 		err = errors.New("Unable to get user data")
		// 		return
		// 	}
		// 	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len2 - ", len(data2))
		// 	if len(data2) <= 0 {
		// 		err = errors.New("User data not found")
		// 		return
		// 	}
		// 	for j := range data2 {
		// 		c.Data["ProducerSellPrice"] = data2[j][1]

		// 	}

		// 	row3, err := db.Db.Query(`select product_id,cluster_selling_price,cluster_supply_charge,cluster_packaging_charge,cluster_bank_charge,cluster_discount from Product2Cluster WHERE product_id =?`, data[i][13])
		// 	if err != nil {
		// 		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		// 		err = errors.New("Unable to get user data")
		// 		return
		// 	}
		// 	defer sql.Close(row3)
		// 	_, data3, err := sql.Scan(row3)
		// 	if err != nil {
		// 		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		// 		err = errors.New("Unable to get user data")
		// 		return
		// 	}
		// 	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len2 - ", len(data3))
		// 	if len(data3) <= 0 {
		// 		err = errors.New("cluster selling data not found")
		// 		return
		// 	}
		// 	for k := range data3 {
		// 		c.Data["ClusterSellPrice"] = data3[k][1]
		// 		c.Data["ClusterSupplyCharge"] = data3[k][2]
		// 		c.Data["ClusterPackagingCharge"] = data3[k][3]
		// 		c.Data["ClusterBankCharge"] = data3[k][4]
		// 		c.Data["ClusterDiscount"] = data3[k][5]

		// 	}
		// }
	}
	return
}

func Inventoryupdate(
	input_product_id,
	input_cluster_name,
	input_category_name,
	input_subcategory_name,
	input_product_name,
	input_product_sku,
	input_productQTY,
	input_receivedunits,
	input_status,
	input_product_qulaitymesure,
	input_product_shelf_life,
	input_product_package_type,
	input_productPackageQTY,
	input_product_origin,
	input_krishivala_selling_price,
	input_cluster_buy_price,
	input_cluster_selling_price,
	input_producer_selling_price,
	runits,
	aunits,
	psku,
	cluster,
	input_damagedunits,
	damageunits,
	id,
	cluster_supply_charge,
	cluster_packaging_charge,
	cluster_bank_charge,
	cluster_discount string) (err error) {

	err = nil

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_product_sku - ", input_product_sku)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "psku - ", psku)

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_cluster_name - ", input_cluster_name)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cluster - ", cluster)

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_damagedunits - ", input_damagedunits)

	if input_product_sku == psku && input_cluster_name == cluster {

		prunits, err := strconv.ParseFloat(runits, 64)
		paunits, err := strconv.ParseFloat(aunits, 64)
		rrunits, err := strconv.ParseFloat(input_receivedunits, 64)
		damageunits1, err := strconv.ParseFloat(input_damagedunits, 64)
		damageunits2, err := strconv.ParseFloat(damageunits, 64)

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "prunits - ", prunits)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "rrunits - ", rrunits)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "paunits - ", paunits)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "damageunits - ", damageunits)

		TD := prunits + rrunits
		TotalDel := strconv.FormatFloat(TD, 'f', 0, 32)

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "TotalDel - ", TotalDel)

		AD := paunits + rrunits
		TotalAvl := strconv.FormatFloat(AD, 'f', 0, 32)

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "TotalAvl - ", TotalAvl)

		DD := damageunits1 + damageunits2
		TotalDamageUnits := strconv.FormatFloat(DD, 'f', 0, 32)

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "TotalAvl - ", TotalAvl)

		res2, err := db.Db.Exec(`UPDATE Inventory SET inventory_received = ? ,inventory_available=?,producer_selling_price=?,damaged_units=?
		              WHERE product_sku = ? AND cluster_name=?`, TotalDel, TotalAvl, input_producer_selling_price, TotalDamageUnits, input_product_sku, input_cluster_name)

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Inventory Get Failed")

		}

		i2, err := res2.RowsAffected()
		if err != nil || i2 == 0 {
			err = errors.New(" Inventory updation failed")

		}

	} else {

		result, err := db.Db.Exec(`INSERT INTO Inventory (
		product_id,
		cluster_name,
		category_name, 
		subcategory_name,
		product_name,
		product_sku,
		product_qty,
		starting_inventory,
		inventory_received, 
		status,
		product_quality_mesure,
		product_shelf_life,
		product_package_type,
		product_package_qty,
		product_origin,
		krishivala_selling_price,
		cluster_buying_price, 
		cluster_selling_price,
		producer_selling_price,
		inventory_available,
		create_by,
		cluster_supply_charge,
		cluster_packaging_charge,
		cluster_bank_charge,
		cluster_discount,
		create_date )
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,now())`,
			input_product_id,
			input_cluster_name,
			input_category_name,
			input_subcategory_name,
			input_product_name,
			input_product_sku,
			input_productQTY,
			input_receivedunits,
			input_receivedunits,
			input_status,
			input_product_qulaitymesure,
			input_product_shelf_life,
			input_product_package_type,
			input_productPackageQTY,
			input_product_origin,
			input_krishivala_selling_price,
			input_cluster_buy_price,
			input_cluster_selling_price,
			input_producer_selling_price,
			input_receivedunits,
			id,
			cluster_supply_charge,
			cluster_packaging_charge,
			cluster_bank_charge,
			cluster_discount)
		if err != nil {
			err = errors.New("Inventory creation failed")

		}

		i, err := result.RowsAffected()
		if err != nil || i == 0 {
			err = errors.New("Inventory creation failed")

		}

	}

	return

}
