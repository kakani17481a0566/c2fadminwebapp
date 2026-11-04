package customerorderinvoice

import (
	"c2f_admin_webapp/session"

	"c2f_admin_webapp/model/utils"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/services"
	"errors"
	"runtime/debug"
	"strconv"
	"strings"

	"c2f_admin_webapp/model/db"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

type CustomerorderInvoice struct {
	// beego.Controller
	middleware.BaseController
}

func (c *CustomerorderInvoice) Get() {
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

			c.TplName = "invoice/customerOderInvoice/customerOderInvoice.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "invoice/customerOderInvoice/customerOderInvoice.html"
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
	create_date,
	order_id,
	category_name,
	subcategory_name,
	product_name,
	
	product_qty,
	cluster_name,
	product_price, 
	total_product_price,
	order_units,
	
	orderdeliver_fromdate,
	orderdeliver_todate,
	order_repeat,
	orderdeliver_time, 
	order_purchagetype,
	
	customer_mob,
	customer_id, 
	status,
	product_package_qty,
	product_sku,
	
	cluster_supply_charge,
	cluster_packaging_charge,
	cluster_bank_charge,
	cluster_discount,
	krishivala_CGST_charge,
	
	krishivala_SGST_charge,
	krishivala_CESS_charge,
	krishivala_OtherTAX_charge,
	total_selling_price,
	total_days,
	order_type,
	producer_id,
	cluster_type,
	delivery_charges,
	delivery_information,
	delivery_discount,
	payu_txnid,
	payu_total_tax,
	payment_type,
	delivery_address_details
	from CustomerOrderTxn WHERE order_id =?  limit 1`, AdminId)
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

	row1, err := db.Db.Query(`select first_name,email,address1, address2,city,zipcode from Customer WHERE mob_no =?  limit 1`, data[0][15])
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get customer data")
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

	// if data[0][14] == "DIRECT_PURCHASE" {

	// }

	c.Data["OrderDate"] = data[0][0]
	c.Data["OrderID"] = data[0][1]
	c.Data["CategoryName"] = data[0][2]
	c.Data["SubcategoryName"] = data[0][3]
	c.Data["ProductName"] = data[0][4]
	c.Data["ProductQty"] = data[0][5]

	c.Data["ClusterName"] = data[0][6]
	c.Data["ProductPrice"] = data[0][7]
	c.Data["TotalPrice"] = data[0][8]
	c.Data["ProductUnits"] = data[0][9]
	c.Data["DF"] = data[0][10]

	c.Data["DT"] = data[0][11]
	c.Data["RepeatOrder"] = data[0][12]
	c.Data["DeliveryTime"] = data[0][13]
	c.Data["PurchaseType"] = data[0][14]
	c.Data["CustomerMobile"] = data[0][15]

	c.Data["CustomerID"] = data[0][17]
	c.Data["Status"] = data[0][17]
	c.Data["ProductPackageQty"] = data[0][18]
	c.Data["ProductSKU"] = data[0][19]

	c.Data["Clustersupplycharge"] = data[0][20]
	c.Data["Clusterpackagingcharge"] = data[0][21]
	c.Data["Clusterbankcharge"] = data[0][22]
	c.Data["Clusterdiscountcharge"] = data[0][23]

	c.Data["KrishivalaCGSTcharge"] = data[0][24]
	c.Data["KrishivalaSGSTcharge"] = data[0][25]
	c.Data["KrishivalaCESScharge"] = data[0][26]
	c.Data["KrishivalaOtherTAXcharge"] = data[0][27]
	c.Data["TotalSalePrice"] = data[0][28]
	c.Data["TotalDays"] = data[0][29]
	c.Data["OrderType"] = data[0][30]

	c.Data["CustomerName"] = data1[0][0]
	c.Data["Address1"] = data1[0][2]
	c.Data["Address2"] = data1[0][3]
	c.Data["City"] = data1[0][4]
	c.Data["Zipcode"] = data1[0][5]
	c.Data["Mobile"] = data[0][15]
	c.Data["Email"] = data1[0][1]
	c.Data["ClusterType"] = data[0][32]
	c.Data["DeliveryCharge"] = data[0][33]

	// c.Data["DeliveryInformation"] = data[0][34]
	deliveryInfoRaw := data[0][34]
	parts := strings.SplitN(deliveryInfoRaw, "::", 2)
	if len(parts) == 2 {
		c.Data["DeliveryChannelPartnerType"] = parts[0]
		c.Data["DeliveryInformation"] = parts[1]
	} else {
		c.Data["DeliveryChannelPartnerType"] = ""
		c.Data["DeliveryInformation"] = deliveryInfoRaw
	}

	c.Data["DeliveryDiscount"] = data[0][35]
	c.Data["payutxnid"] = data[0][36]
	c.Data["payutax"] = data[0][37]
	c.Data["paymenttype"] = data[0][38]
	c.Data["DeliveryAddressDetails"] = data[0][39]

	value1Str := data[0][20]
	value2Str := data[0][21]
	totalDaysStr := data[0][29]

	value1, err1 := strconv.ParseFloat(value1Str, 64)
	if err1 != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error converting data[0][20] to float64:", err1)
		return
	}

	value2, err2 := strconv.ParseFloat(value2Str, 64)
	if err2 != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error converting data[0][21] to float64:", err2)
		return
	}

	totalDays, err3 := strconv.ParseFloat(totalDaysStr, 64)
	if err3 != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error converting data[0][29] to float64:", err3)
		return
	}

	if totalDays == 0 {
		log.Println(beego.AppConfig.String("loglevel"), "Error: Division by zero (totalDays is zero)")
		return
	}

	c.Data["SupplyChargePerDay"] = value1 / totalDays
	c.Data["PackageChargePerDay"] = value2 / totalDays

	if data[0][14] == "DIRECT_PURCHASE" {

		row2, err2 := db.Db.Query(`select first_name,
last_name,
mob_no,
address1,
address2,
town,
city,
state,
zipcode from Customer WHERE id =?`, data[0][31])
		if err2 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err2 = errors.New("Unable to get producer data")
			return
		}
		defer sql.Close(row2)
		_, data2, err := sql.Scan(row2)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get producer data")
			return
		}

		c.Data["ProducerName"] = data2[0][0] + " " + data2[0][1]
		c.Data["ProducerMobile"] = data2[0][2]
		c.Data["ProducerAddress1"] = data2[0][3]
		c.Data["ProducerAddress2"] = data2[0][4]
		c.Data["ProducerTown"] = data2[0][5]
		c.Data["ProducerCity"] = data2[0][6]
		c.Data["ProducerState"] = data2[0][7]
		c.Data["Producerpincode"] = data2[0][8]

	}

	return

}

func (c *CustomerorderInvoice) Post() {

	msg11 := ""

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

			c.TplName = "invoice/customerOderInvoice/customerOderInvoice.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = msg11
			c.TplName = "invoice/customerOderInvoice/customerOderInvoice.html"
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
	create_date,
	order_id,
	category_name,
	subcategory_name,
	product_name,
	
	product_qty,
	cluster_name,
	product_price, 
	total_product_price,
	order_units,
	
	orderdeliver_fromdate,
	orderdeliver_todate,
	order_repeat,
	orderdeliver_time, 
	order_purchagetype,
	
	customer_mob,
	customer_id, 
	status,
	product_package_qty,
	product_sku,
	
	cluster_supply_charge,
	cluster_packaging_charge,
	cluster_bank_charge,
	cluster_discount,
	krishivala_CGST_charge,
	
	krishivala_SGST_charge,
	krishivala_CESS_charge,
	krishivala_OtherTAX_charge,
	total_selling_price,
	total_days,
	order_type,
	producer_id,
	cluster_type
	 
	from CustomerOrderTxn WHERE order_id =?  limit 1`, AdminId)
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

	row1, err := db.Db.Query(`select first_name,email,address1, address2,city,zipcode from Customer WHERE mob_no =?  limit 1`, data[0][15])
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get customer data")
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

	// if data[0][14] == "DIRECT_PURCHASE" {

	// }

	// c.Data["OrderDate"] = data[0][0]
	// c.Data["OrderID"] = data[0][1]
	// c.Data["CategoryName"] = data[0][2]
	// c.Data["SubcategoryName"] = data[0][3]
	// c.Data["ProductName"] = data[0][4]
	// c.Data["ProductQty"] = data[0][5]

	// c.Data["ClusterName"] = data[0][6]
	// c.Data["ProductPrice"] = data[0][7]
	// c.Data["TotalPrice"] = data[0][8]
	// c.Data["ProductUnits"] = data[0][9]
	// c.Data["DF"] = data[0][10]

	// c.Data["DT"] = data[0][11]
	// c.Data["RepeatOrder"] = data[0][12]
	// c.Data["DeliveryTime"] = data[0][13]
	// c.Data["PurchaseType"] = data[0][14]
	// c.Data["CustomerMobile"] = data[0][15]

	// c.Data["CustomerID"] = data[0][17]
	// c.Data["Status"] = data[0][17]
	// c.Data["ProductPackageQty"] = data[0][18]
	// c.Data["ProductSKU"] = data[0][19]

	// c.Data["Clustersupplycharge"] = data[0][20]
	// c.Data["Clusterpackagingcharge"] = data[0][21]
	// c.Data["Clusterbankcharge"] = data[0][22]
	// c.Data["Clusterdiscountcharge"] = data[0][23]

	// c.Data["KrishivalaCGSTcharge"] = data[0][24]
	// c.Data["KrishivalaSGSTcharge"] = data[0][25]
	// c.Data["KrishivalaCESScharge"] = data[0][26]
	// c.Data["KrishivalaOtherTAXcharge"] = data[0][27]
	// c.Data["TotalSalePrice"] = data[0][28]
	// c.Data["TotalDays"] = data[0][29]
	// c.Data["OrderType"] = data[0][30]

	// c.Data["CustomerName"] = data1[0][0]
	// c.Data["Address1"] = data1[0][2]
	// c.Data["Address2"] = data1[0][3]
	// c.Data["City"] = data1[0][4]
	// c.Data["Zipcode"] = data1[0][5]
	// c.Data["Mobile"] = data[0][15]
	// c.Data["Email"] = data1[0][1]
	// c.Data["ClusterType"] = data[0][32]

	current_status_of_order := data[0][17]

	if data[0][14] == "DIRECT_PURCHASE" {

		row2, err2 := db.Db.Query(`select first_name,last_name,mob_no,address1,address2,town,city,state,zipcode from Customer WHERE id =?`, data[0][31])
		if err2 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err2 = errors.New("Unable to get producer data")
			return
		}
		defer sql.Close(row2)
		_, data2, err := sql.Scan(row2)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get producer data")
			return
		}

		c.Data["ProducerName"] = data2[0][0] + " " + data2[0][1]
		c.Data["ProducerMobile"] = data2[0][2]
		c.Data["ProducerAddress1"] = data2[0][3]
		c.Data["ProducerAddress2"] = data2[0][4]
		c.Data["ProducerTown"] = data2[0][5]
		c.Data["ProducerCity"] = data2[0][6]
		c.Data["ProducerState"] = data2[0][7]
		c.Data["Producerpincode"] = data2[0][8]

	}

	// delivery_charges := c.Input().Get("delivery_charges")
	// deliverychannelpartnertype := c.Input().Get("deliverychannelpartnertype")
	delivery_information := c.Input().Get("deliverychannelpartnertype") + "::" + c.Input().Get("delivery_information")
	input_status1 := c.Input().Get("input_status1")
	input_remarks1 := c.Input().Get("input_remarks1")
	input_status2 := c.Input().Get("input_status2")
	input_remarks2 := c.Input().Get("input_remarks2")

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "delivery_information------------------- - ", delivery_information)

	// if data[0][17] == "ORDER_RECEIVED" && delivery_charges == "" {

	// 	err = errors.New("delivery charges cannot be  empty")
	// 	return

	// }

	if data[0][32] == "SOURCE" && data[0][17] == "ORDER_RECEIVED" {

		if input_status1 == "" || input_remarks1 == "" {

			err = errors.New("status and remarks cannot be  empty")
			return

		}

		// update the status and remarks for source orders and non bulk product orders

		res, err := db.Db.Exec(`UPDATE CustomerOrderTxn SET status=?,remarks=? WHERE order_id = ?`, input_status1, input_remarks1, data[0][1])
		if err != nil {
			err = errors.New("status and remarks  updation failed")
			return
		}

		i, err := res.RowsAffected()
		if err != nil || i == 0 {
			err = errors.New("status and remarks updation failed")
			return
		}

		msg11 = " Status and Remarks Updated Successfully"

	}

	if data[0][32] == "OTHER" && (data[0][17] == "ORDER_RECEIVED" || data[0][17] == "ORDER_DELIVERY_ASSIGNED") {

		if input_status2 == "" || input_remarks2 == "" {

			err = errors.New("status and remarks cannot be  empty")
			return

		}

		if current_status_of_order == "ORDER_RECEIVED" && input_status2 == "ORDER_DELIVERY_COMPLETED" {

			err = errors.New("Invalid Status Input : Order cannot be deliverd without any delivery assignment")
			return

		}

		if input_status2 == "ORDER_REJECTED" {

			// update only status and remarks in case of rejection : other cluster

			res, err := db.Db.Exec(`UPDATE CustomerOrderTxn SET status=?,remarks=? WHERE order_id = ?`, input_status2, input_remarks2, data[0][1])
			if err != nil {
				err = errors.New("delivery information updation failed")
				return
			}

			i, err := res.RowsAffected()
			if err != nil || i == 0 {
				err = errors.New("delivery information updation failed")
				return
			}
			msg11 = "Order status & remarks  Updated Successfully"
		}

		if (input_status2 == "ORDER_DELIVERY_ASSIGNED" || input_status2 == "ORDER_DELIVERY_COMPLETED") && delivery_information == "" {

			err = errors.New("delivery information cannot be  empty")
			return

		}

		if input_status2 == "ORDER_DELIVERY_ASSIGNED" {

			// update all 3 i.e; status and remarks and delivery info in case of assigned : other cluster

			// Get available units at the Inventory for the given product and cluster

			aunits, err := services.Getinventory2(data[0][19], data[0][6])
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Inventory Get Failed5")
				msg11 = "Unable To Get Inventory Details , Please Check Inventory for availability of Product"
				return
			}

			log.Println(beego.AppConfig.String("loglevel"), "Debug", "Available Units at inventory - ", aunits)

			punits2, err := strconv.ParseFloat(data[0][9], 64)
			aunits1, err := strconv.ParseFloat(aunits, 64)
			TDDD, err := strconv.ParseFloat(data[0][29], 64)
			PT := punits2 * TDDD

			log.Println(beego.AppConfig.String("loglevel"), "Debug", "Total Ordered Units is ( units * days ) - ", PT)

			if PT > aunits1 {

				//log.Println(beego.AppConfig.String("loglevel"), "Debug", "inside - ")

				err = errors.New("Insufficent Available Inventory")
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				msg11 = "Insufficent Available Inventory For the Order , Check The Inventory Before Updating Delivery Status"

				return
			}

			// update the delivery information

			res, err := db.Db.Exec(`UPDATE CustomerOrderTxn SET status=?,delivery_information=?,remarks=? WHERE order_id = ?`, input_status2, delivery_information, input_remarks2, data[0][1])
			if err != nil {
				err = errors.New("delivery information updation failed")
				return
			}

			i, err := res.RowsAffected()
			if err != nil || i == 0 {
				err = errors.New("delivery information updation failed")
				return
			}

			// Update Inventory

			dunits, aunits, runits, sku, cluster, damageunits, err2 := services.Getinventory(data[0][19], data[0][6])
			if err2 != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err2 = errors.New("Inventory Get Failed1")
				msg11 = "Unable To Get Inventory Details , Please Check Inventory for availability of Product"
				return
			}

			log.Println(beego.AppConfig.String("loglevel"), "Debug", "delivered units - ", dunits)
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "available units - ", aunits)
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "sku - ", sku)
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "cluster - ", cluster)
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "damageunits - ", damageunits)

			ddunits, err := strconv.ParseFloat(dunits, 64)
			rrunits, err := strconv.ParseFloat(runits, 64)
			ppunits, err := strconv.ParseFloat(data[0][9], 64)
			TDDD1, err := strconv.ParseFloat(data[0][29], 64)
			// total ordered units of product
			ppunits1 := ppunits * TDDD1

			TD := ddunits + ppunits1
			TotalDel := strconv.FormatFloat(TD, 'f', 0, 32)
			TA := rrunits - TD
			TotalAvl := strconv.FormatFloat(TA, 'f', 0, 32)

			log.Println(beego.AppConfig.String("loglevel"), "Debug", "updated delivered units - ", TotalDel)
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "updated available units - ", TotalAvl)

			res2, err := db.Db.Exec(`UPDATE Inventory SET inventory_delivered = ? ,inventory_available=? WHERE product_sku = ? AND cluster_name=?`, TotalDel, TotalAvl, data[0][19], data[0][6])
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Inventory Get Failed3")
				msg11 = "Update Inventory Failed"
				return
			}

			i2, err := res2.RowsAffected()
			if err != nil || i2 == 0 {
				err = errors.New(" CustomerOrderTxn Delivery updation failed")
				return
			}

			msg11 = "Order status & Delivery Information Updated Successfully"

		}

		if input_status2 == "ORDER_DELIVERY_COMPLETED" {

			if current_status_of_order == "ORDER_RECEIVED" && input_status2 == "ORDER_DELIVERY_COMPLETED" {

				err = errors.New("Invalid Status Input : Order cannot be deliverd without any delivery assignment")
				return

			}

			// update the delivery information

			res, err := db.Db.Exec(`UPDATE CustomerOrderTxn SET status=?,delivery_information=?,remarks=? WHERE order_id = ?`, input_status2, delivery_information, input_remarks2, data[0][1])
			if err != nil {
				err = errors.New("delivery information updation failed")
				return
			}

			i, err := res.RowsAffected()
			if err != nil || i == 0 {
				err = errors.New("delivery information updation failed")
				return
			}

		}

		msg11 = "Order status & Delivery Information Updated Successfully"

	}

	return

}
