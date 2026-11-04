package collectpaymentinvoice

import (
	"c2f_admin_webapp/session"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"

	"errors"
	"runtime/debug"
	"strings"

	"fmt"
	"strconv"

	"c2f_admin_webapp/model/db"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

type Row struct {
	OrderID string
}
type CollectPaymentInvoice struct {
	// beego.Controller
	middleware.BaseController
}

func (c *CollectPaymentInvoice) Get() {
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

			c.TplName = "invoice/collectPaymentInvoice/collectPaymentInvoice.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "invoice/collectPaymentInvoice/collectPaymentInvoice.html"
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
	Oder2Delivery.oder_id,
	Oder2Delivery.order_date,
	Oder2Delivery.delivery_dates,
    Oder2Delivery.status,
    Oder2Delivery.remarks,
    
    JSON_UNQUOTE(JSON_EXTRACT(Oder2Delivery.customer_details, '$.name')),
	JSON_UNQUOTE(JSON_EXTRACT(Oder2Delivery.customer_details, '$.address1')),
	JSON_UNQUOTE(JSON_EXTRACT(Oder2Delivery.customer_details, '$.address2')),
	JSON_UNQUOTE(JSON_EXTRACT(Oder2Delivery.customer_details, '$.city')),
	JSON_UNQUOTE(JSON_EXTRACT(Oder2Delivery.customer_details, '$.mobile')),
	
	JSON_UNQUOTE(JSON_EXTRACT(Oder2Delivery.oder_deatils, '$.productName')),
    JSON_UNQUOTE(JSON_EXTRACT(Oder2Delivery.oder_deatils, '$.productPackageQty')),
	JSON_UNQUOTE(JSON_EXTRACT(Oder2Delivery.oder_deatils, '$.productQty')),
    JSON_UNQUOTE(JSON_EXTRACT(Oder2Delivery.oder_deatils, '$.productSKU')),
    JSON_UNQUOTE(JSON_EXTRACT(Oder2Delivery.oder_deatils, '$.productPrice')),
	JSON_UNQUOTE(JSON_EXTRACT(Oder2Delivery.oder_deatils, '$.productUnits')),
	
	JSON_UNQUOTE(JSON_EXTRACT(Oder2Delivery.oder_deatils, '$.deliveryTime')),
	JSON_UNQUOTE(JSON_EXTRACT(Oder2Delivery.oder_deatils, '$.deliveryFrom')),
	JSON_UNQUOTE(JSON_EXTRACT(Oder2Delivery.oder_deatils, '$.deliveryTo')),
    JSON_UNQUOTE(JSON_EXTRACT(Oder2Delivery.oder_deatils, '$.clusterName')),
    JSON_UNQUOTE(JSON_EXTRACT(Oder2Delivery.oder_deatils, '$.zipcode')),
    JSON_UNQUOTE(JSON_EXTRACT(Oder2Delivery.oder_deatils, '$.orderType')),
     JSON_UNQUOTE(JSON_EXTRACT(Oder2Delivery.oder_deatils, '$.repeatOrder')),
    JSON_UNQUOTE(JSON_EXTRACT(Oder2Delivery.oder_deatils, '$.categoryName')),
    JSON_UNQUOTE(JSON_EXTRACT(Oder2Delivery.oder_deatils, '$.subcategroryName')),
    CustomerOrderTxn.payment_status,
    CustomerOrderTxn.payment_type,
    CustomerOrderTxn.order_type,
    CustomerOrderTxn.total_days
         from Oder2Delivery
        LEFT JOIN  CustomerOrderTxn
        ON Oder2Delivery.oder_id = CustomerOrderTxn.order_id
    where Oder2Delivery.oder_id=?`, AdminId)
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

	row5, err := db.Db.Query(`select JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.productUnits'))
	from Oder2Delivery
    where oder_id=? AND status="DELIVERY_COMPLETED" limit 1`, AdminId)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get supply transaction data")
		return
	}
	defer sql.Close(row5)
	_, data5, err := sql.Scan(row5)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data5, "\nData len - ", len(data5))
	if len(data5) <= 0 {
		err = errors.New("User data not found")
		return
	}

	if len(data5) <= 0 {
		err = errors.New("Unable to get the error message info")
		return
	}
	u5 := data5[0][0]

	row2, err := db.Db.Query(`select SUM(JSON_EXTRACT(oder_deatils, "$.productPrice")) as oder_deatils
	from Oder2Delivery
    where oder_id=? AND status="DELIVERY_COMPLETED" limit 1`, AdminId)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get supply transaction data")
		return
	}
	defer sql.Close(row2)
	_, data2, err := sql.Scan(row2)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data2, "\nData len - ", len(data2))
	if len(data2) <= 0 {
		err = errors.New("User data not found")
		return
	}

	if len(data2) <= 0 {
		err = errors.New("Unable to get the error message info")
		return
	}

	p5 := data2[0][0]

	u6, err := strconv.ParseFloat(u5, 64)
	p6, err := strconv.ParseFloat(p5, 64)
	totalprice5 := (u6 * p6)
	totalprice6 := strconv.FormatFloat(totalprice5, 'f', 0, 32)

	c.Data["TotalAmount"] = totalprice6

	row3, err := db.Db.Query(`select 
	order_id,
	product_price, 
	total_product_price,
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
	status,
	delivery_charges,
	product_qty
	from CustomerOrderTxn WHERE order_id =? AND product_sku=? AND status="DELIVERY_COMPLETED"  limit 1`, AdminId, data[0][13])
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get supply transaction data")
		return
	}
	defer sql.Close(row3)
	_, data3, err := sql.Scan(row3)
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

	c.Data["OrderId"] = data[0][0]
	c.Data["OrderDate"] = data[0][1]
	c.Data["DeliveryDate"] = data[0][2]
	c.Data["Status"] = data[0][3]
	c.Data["Remarks"] = data[0][4]
	c.Data["CustomerName"] = data[0][5]
	c.Data["CustomerAdd1"] = data[0][6]

	c.Data["CustomerAdd2"] = data[0][7]
	c.Data["CustomerCity"] = data[0][8]
	c.Data["CustomerMobile"] = data[0][9]
	c.Data["ProductName"] = data[0][10]
	c.Data["ProductPackageQty"] = data[0][11]

	c.Data["ProductQty"] = data[0][12]
	c.Data["ProductSKU"] = data[0][13]
	c.Data["UnitPrice"] = data[0][14]
	c.Data["ProductUnit"] = data[0][15]
	c.Data["DeliveryTime"] = data[0][16]

	c.Data["DeliveryFrom"] = data[0][17]
	c.Data["DeliveryTo"] = data[0][18]
	c.Data["ClusterName"] = data[0][19]
	c.Data["ZipCode"] = data[0][20]
	c.Data["PurchaseType"] = data[0][21]

	c.Data["RepeatOrder"] = data[0][22]
	c.Data["CategoryName"] = data[0][23]
	c.Data["SubcategroryName"] = data[0][24]
	c.Data["PaymentStatus"] = data[0][25]

	c.Data["PaymentType"] = data[0][26]
	c.Data["OrderType"] = data[0][27]
	c.Data["TotalDays"] = data[0][28]
	//c.Data["KrishivalaOtherTAXcharge"] = data[0][27]
	// c.Data["TotalSalePrice"] = data[0][28]
	// c.Data["TotalDays"] = data[0][29]
	// c.Data["PaymentStatus"] = data[0][30]
	// c.Data["PaymentType"] = data[0][31]
	// c.Data["OrderType"] = data[0][32]

	// c.Data["CustomerName"] = data1[0][0]
	// c.Data["Address1"] = data1[0][2]
	// c.Data["Address2"] = data1[0][3]
	// c.Data["City"] = data1[0][4]
	// c.Data["Zipcode"] = data1[0][5]
	// c.Data["Mobile"] = data[0][15]
	// c.Data["Email"] = data1[0][1]

	c.Data["OrderID"] = data3[0][0]
	c.Data["ProductPrice"] = data3[0][1]
	c.Data["TotalPrice"] = data3[0][2]

	c.Data["Clustersupplycharge"] = data3[0][3]
	c.Data["Clusterpackagingcharge"] = data3[0][4]
	c.Data["Clusterbankcharge"] = data3[0][5]
	c.Data["Clusterdiscountcharge"] = data3[0][6]

	c.Data["KrishivalaCGSTcharge"] = data3[0][7]
	c.Data["KrishivalaSGSTcharge"] = data3[0][8]
	c.Data["KrishivalaCESScharge"] = data3[0][9]
	c.Data["KrishivalaOtherTAXcharge"] = data3[0][10]
	c.Data["TotalSalePrice"] = data3[0][11]
	c.Data["TotalDays"] = data3[0][12]
	c.Data["OrderType"] = data3[0][13]

	c.Data["DeliveryCharge"] = data3[0][15]
	c.Data["ProductQty"] = data3[0][16]
	TTotalAmount := totalprice6

	cluster_supply_charge := data3[0][3]
	cluster_packaging_charge := data3[0][4]
	cluster_bank_charge := data3[0][5]
	cluster_discount := data3[0][6]

	krishivala_CGST_charge := data3[0][7]
	krishivala_SGST_charge := data3[0][8]
	krishivala_CESS_charge := data3[0][9]
	krishivala_OtherTAX_charge := data3[0][10]

	scluster_supply_charge, err := strconv.ParseFloat(cluster_supply_charge, 64)
	scluster_packaging_charge, err := strconv.ParseFloat(cluster_packaging_charge, 64)
	scluster_bank_charge, err := strconv.ParseFloat(cluster_bank_charge, 64)
	scluster_discount, err := strconv.ParseFloat(cluster_discount, 64)

	TTotalAAmount, err := strconv.ParseFloat(TTotalAmount, 64)

	skrishivala_CGST_charge, err := strconv.ParseFloat(krishivala_CGST_charge, 64)
	skrishivala_SGST_charge, err := strconv.ParseFloat(krishivala_SGST_charge, 64)
	skrishivala_CESS_charge, err := strconv.ParseFloat(krishivala_CESS_charge, 64)
	skrishivala_OtherTAX_charge, err := strconv.ParseFloat(krishivala_OtherTAX_charge, 64)

	charges := scluster_supply_charge + scluster_packaging_charge + scluster_bank_charge
	taxs := skrishivala_CGST_charge + skrishivala_SGST_charge + skrishivala_CESS_charge + skrishivala_OtherTAX_charge

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "charges - ", charges)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "taxs - ", taxs)

	TPP := TTotalAAmount + (charges) - scluster_discount + (taxs)

	if data3[0][16] == "gms" {
		deliverycharges, err := strconv.ParseFloat(data3[0][15], 64)
		fmt.Println(err)

		TPP = TPP + deliverycharges

	}

	tpp := strconv.FormatFloat(TPP, 'f', 0, 32)

	//	log.Println(beego.AppConfig.String("loglevel"), "Debug", "select_Customer  - ", select_Customer)

	//  TotalBuying_price = cluster.buyingprice  + (cluster.supplycarge + cluster.packagingcharge + cluster.bankcharge ) - cluster.discount  + (producer.gcst + producer.scst + producer.C29 + producer.othertax)
	TotalBillAmount := tpp

	c.Data["TotalBillAmount"] = TotalBillAmount

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "TPP  - ", TPP)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "TotalBillAmount  - ", TotalBillAmount)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "TPP  - ", TPP)

	return

}
