package supplyorderinvoice

import (
	"c2f_admin_webapp/session"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"

	//	"encoding/json"
	"errors"
	"runtime/debug"
	"strings"

	"c2f_admin_webapp/model/db"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

type SupplyorderInvoice struct {
	// beego.Controller
	middleware.BaseController
}

func (c *SupplyorderInvoice) Get() {
	AdminId := c.Ctx.Input.Param(":AdminID")
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

			c.TplName = "invoice/supplyOderInvoice/supplyOderInvoice.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "invoice/supplyOderInvoice/supplyOderInvoice.html"
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

	auth, err := utils.IsAuthorized(utype, "Invoice")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	row, err := db.Db.Query(`select 
	supplyordertxn_id,
	create_date, 
	producer_mob,
	category_name, 
	subcategory_name, 
	
	product_name,
	product_qty,
	product_package_qty,
	product_sku,
	cluster_name, 
	
	cluster_buying_price,
	product_supply_units,
	supply_received_units, 
	updated_total_supply_price,
	total_buy_price,  
	
	status,
	cluster_supply_charge,
	cluster_packaging_charge,
	cluster_discount_charge,
	producer_CGST_charge,
	
	producer_SGST_charge,
	producer_CESS_charge,
	producer_OtherTAX_charge,
	cluster_bank_charge,
	remarks
	
	from SupplyOrderTxn WHERE supplyordertxn_id =?  limit 1`, AdminId)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get supply transaction data")
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

	row1, err := db.Db.Query(`select first_name,email,address1, address2,city,zipcode from Customer WHERE mob_no =?  limit 1`, data[0][2])
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
	c.Data["ProducerMobile"] = data[0][2]
	c.Data["CategoryName"] = data[0][3]
	c.Data["SubcategoryName"] = data[0][4]

	c.Data["ProductName"] = data[0][5]
	c.Data["ProductQty"] = data[0][6]
	c.Data["ProductPackageQty"] = data[0][7]
	c.Data["ProductSKU"] = data[0][8]
	c.Data["ClusterName"] = data[0][9]

	c.Data["ClusterBuyingPrice"] = data[0][10]
	c.Data["SuppliedUnits"] = data[0][11]
	c.Data["ReceivedUnits"] = data[0][12]
	c.Data["Updatedtotalsupplyprice"] = data[0][13]
	c.Data["TotalbuyPrice"] = data[0][14]

	c.Data["Status"] = data[0][15]
	c.Data["Clustersupplycharge"] = data[0][16]
	c.Data["Clusterpackagingcharge"] = data[0][17]
	c.Data["Clusterdiscountcharge"] = data[0][18]
	c.Data["ProducerCGSTcharge"] = data[0][19]

	c.Data["ProducerSGSTcharge"] = data[0][20]
	c.Data["ProducerCESScharge"] = data[0][21]
	c.Data["ProducerOtherTAXcharge"] = data[0][22]
	c.Data["Clusterbankcharge"] = data[0][23]

	c.Data["ProducerName"] = data1[0][0]
	c.Data["Address1"] = data1[0][2]
	c.Data["Address2"] = data1[0][3]
	c.Data["City"] = data1[0][4]
	c.Data["Zipcode"] = data1[0][5]
	c.Data["Email"] = data1[0][1]
	c.Data["Remarks"] = data[0][24]

	row3, err3 := db.Db.Query(`select count(*) from Oder2SupplyExecutive WHERE supply_order_id = ?`, data[0][0])

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
	if countlen != "0" {

		res22, err := db.Db.Query(`select supply_executive_charges from  Oder2SupplyExecutive  WHERE supply_order_id = ?`, data[0][0])
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
		SE_charges := data22[0][0]

		c.Data["SE_charges"] = SE_charges

	} else {

		c.Data["SE_charges"] = "NA"

	}

	return

}
