package pickupinvoice

import (
	"c2f_admin_webapp/session"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"

	//	"encoding/json"
	"errors"
	"runtime/debug"
	"strconv"
	"strings"

	"c2f_admin_webapp/model/db"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

type PickupInvoice struct {
	// beego.Controller
	middleware.BaseController
}

func (c *PickupInvoice) Get() {
	AdminId := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Info", "PickupInvoice Start")
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

			c.TplName = "invoice/pickupInvoice/pickupInvoice.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "PickupInvoice Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "invoice/pickupInvoice/pickupInvoice.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "PickupInvoice  Page Success")
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

	auth, err := utils.IsAuthorized(utype, "Invoice")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	row, err := db.Db.Query(`select 
supply_order_id,
supply_order_date,
create_date,
JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.category_name')),
JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.subcategory_name')),
JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_name')),
JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_sku')),
JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_package_qty')),
JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_qty')),
JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.cluster_buying_price')),
JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_supply_units')),
JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.name')),
JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.address1')),
JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.address2')),
JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.city')),
JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.zipcode')),
JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.mobile')),
supplyexecutive_mob,
status,
supplyexecutive_cluster,
id,
JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.product_id')),
JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.total_supply_price')),
JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.id')),
JSON_UNQUOTE(JSON_EXTRACT(supply_order_deatils, '$.cluster')),
supply_executive_charges from Oder2SupplyExecutive WHERE id =?  limit 1`, AdminId)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Order to SupplyExecutive data")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Order to SupplyExecutive data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Order to SupplyExecutive data not found")
		return
	}

	row1, err := db.Db.Query(`select first_name,middle_name,last_name,email from SupplyExecutive WHERE mob_no =?  limit 1`, data[0][17])
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get supply transaction data")
		return
	}
	defer sql.Close(row1)
	_, data1, err := sql.Scan(row1)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	if len(data1) <= 0 {
		err = errors.New("User data not found")
		return
	}

	c.Data["SupplyOrderId"] = data[0][0]
	c.Data["SupplyOrderDate"] = data[0][1]
	c.Data["DeliveryAssignedDate"] = data[0][2]
	c.Data["CategoryName"] = data[0][3]
	c.Data["SubcategoryName"] = data[0][4]
	c.Data["ProductName"] = data[0][5]
	c.Data["ProductSKU"] = data[0][6]
	c.Data["ProductPackageQty"] = data[0][7]
	c.Data["ProductQty"] = data[0][8]
	c.Data["ClusterBuyingPrice"] = data[0][9]
	c.Data["SuppliedUnits"] = data[0][10]
	c.Data["ProducerName"] = data[0][11]
	c.Data["Address1"] = data[0][12]
	c.Data["Address2"] = data[0][13]
	c.Data["City"] = data[0][14]
	c.Data["Zipcode"] = data[0][15]
	c.Data["ProducerMobile"] = data[0][16]
	c.Data["Status"] = data[0][18]

	c.Data["SEfirstname"] = data1[0][0]
	c.Data["SEmiddlename"] = data1[0][1]
	c.Data["SElastname"] = data1[0][2]
	c.Data["SEemail"] = data1[0][3]
	c.Data["SEmobile"] = data[0][17]
	c.Data["SECluster"] = data[0][19]

	c.Data["Updatedtotalsupplyprice"] = data[0][22]

	c.Data["SupplyExecCharges"] = data[0][25]

	row22, err22 := db.Db.Query(`select id from  Clusters WHERE cluster_name=?`, data[0][24])
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

	row3, err3 := db.Db.Query(`select product_id,cluster_selling_price,cluster_supply_charge,cluster_packaging_charge,cluster_bank_charge,cluster_discount from Product2Cluster where product_id = ? and cluster_id = ?`, data[0][21], data22[0][0])
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

	cluster_SC := data3[0][2]
	cluster_PC := data3[0][3]
	cluster_BC := data3[0][4]
	cluster_discount := data3[0][5]

	row4, err := db.Db.Query(`select product_id,producer_CGST_charge,producer_SGST_charge,producer_CESS_charge,producer_OtherTAX_charge from Product2Producer2Cluster WHERE product_id =? AND cluster_id=? AND producer_id=?`, data[0][21], data22[0][0], data[0][23])
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

	cgst := data4[0][1]
	sgst := data4[0][2]
	cess := data4[0][3]
	othertax := data4[0][4]

	BP, err := strconv.ParseFloat(data[0][22], 64)

	scluster_supply_charge, err := strconv.ParseFloat(cluster_SC, 64)
	scluster_packaging_charge, err := strconv.ParseFloat(cluster_PC, 64)
	scluster_bank_charge, err := strconv.ParseFloat(cluster_BC, 64)
	scluster_discount, err := strconv.ParseFloat(cluster_discount, 64)

	sproducer_CGST_charge, err := strconv.ParseFloat(cgst, 64)
	sproducer_SGST_charge, err := strconv.ParseFloat(sgst, 64)
	sproducer_CESS_charge, err := strconv.ParseFloat(cess, 64)
	sproducer_OtherTAX_charge, err := strconv.ParseFloat(othertax, 64)

	supply_exec_charges, err := strconv.ParseFloat(data[0][25], 64)

	charges := scluster_supply_charge + scluster_packaging_charge + scluster_bank_charge
	taxs := sproducer_CGST_charge + sproducer_SGST_charge + sproducer_CESS_charge + sproducer_OtherTAX_charge

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "charges - ", charges)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "taxs - ", taxs)

	TBP := BP + (charges) - scluster_discount + (taxs) - supply_exec_charges

	tbp := strconv.FormatFloat(TBP, 'f', 0, 32)

	//  TotalBuying_price = cluster.buyingprice  + (cluster.supplycarge + cluster.packagingcharge + cluster.bankcharge ) - cluster.discount  + (producer.gcst + producer.scst + producer.C29 + producer.othertax)
	TotalBuying_price := tbp

	c.Data["TotalbuyPrice"] = TotalBuying_price

	c.Data["Clustersupplycharge"] = cluster_SC
	c.Data["Clusterpackagingcharge"] = cluster_PC
	c.Data["Clusterdiscountcharge"] = cluster_discount
	c.Data["ProducerCGSTcharge"] = cgst

	c.Data["ProducerSGSTcharge"] = sgst
	c.Data["ProducerCESScharge"] = cess
	c.Data["ProducerOtherTAXcharge"] = othertax
	c.Data["Clusterbankcharge"] = cluster_BC

	return

}
