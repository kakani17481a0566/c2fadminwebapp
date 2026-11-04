package clusterproductpurchaseinvoice

import (
	"c2f_admin_webapp/session"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"

	"errors"
	"strings"

	//	"fmt"
	"runtime/debug"

	//	"strconv"

	"c2f_admin_webapp/model/db"

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
	Product_qty                   string
	Product_package_qty           string
}

type ClusterproductPurchaseInvoice struct {
	// beego.Controller
	middleware.BaseController
}

func (c *ClusterproductPurchaseInvoice) Get() {
	AdminId := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "AdminId - ", AdminId)
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

			c.TplName = "invoice/clusterproductPurchaseInvoice/clusterproductPurchaseInvoice.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "invoice/clusterproductPurchaseInvoice/clusterproductPurchaseInvoice.html"
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

	row, err := db.Db.Query(`select ClusterPurchaseProduct.id, 
	ClusterPurchaseProduct.source_cluster_name,
	ClusterPurchaseProduct.beneficiary_cluster_name, 
	ClusterPurchaseProduct.product_name, 
	ClusterPurchaseProduct.product_sku,
	ClusterPurchaseProduct.number_of_units, 
	ClusterPurchaseProduct.delivery_charges, 
	ClusterPurchaseProduct.total_price, 
	ClusterPurchaseProduct.status,
	ClusterPurchaseProduct.create_date,
	ClusterPurchaseProduct.cluster_productpurchase_txnid,
	ClusterPurchaseProduct.product_price,
	Inventory.product_qty,
	Inventory.product_package_qty,
	ClusterPurchaseProduct.courier_details
 from ClusterPurchaseProduct
LEFT JOIN Inventory ON 
ClusterPurchaseProduct.product_sku = Inventory.product_sku
    where ClusterPurchaseProduct.id=?  limit 1`, AdminId)
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
		c.Data["Source_cluster_name"] = data[i][1]
		c.Data["Beneficiary_cluster_name"] = data[i][2]
		c.Data["Product_name"] = data[i][3]
		c.Data["Product_sku"] = data[i][4]
		c.Data["Number_of_units"] = data[i][5]
		c.Data["Delivery_charges"] = data[i][6]
		c.Data["Total_price"] = data[i][7]
		c.Data["Status"] = data[i][8]
		c.Data["TimeStamp"] = data[i][9]
		c.Data["Cluster_productpurchase_txnid"] = data[i][10]
		c.Data["Product_price"] = data[i][11]
		c.Data["Product_qty"] = data[i][12]
		c.Data["Product_package_qty"] = data[i][13] + data[i][12]
		courierParts := strings.SplitN(data[i][14], "::", 2)
		if len(courierParts) == 2 {
			c.Data["Delivery_channel_partner_type"] = courierParts[0]
			c.Data["Courier_details"] = courierParts[1]
		} else {
			c.Data["Delivery_channel_partner_type"] = "-"
			c.Data["Courier_details"] = data[i][14]
		}

	}

	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "productqty", data[0][12])
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "product package qty", data[0][13])

	// ProductPrice := data[0][11]
	// NumberOfUnits := data[0][5]
	// DeliveryCharges := data[0][6]

	// ProductPrice1, err := strconv.ParseFloat(ProductPrice, 64)
	// NumberOfUnits1, err := strconv.ParseFloat(NumberOfUnits, 64)
	// DeliveryCharges1, err := strconv.ParseFloat(DeliveryCharges, 64)

	// TotalPrice1 := (ProductPrice1 * NumberOfUnits1) + DeliveryCharges1

	// TotalPrice := strconv.FormatFloat(TotalPrice1, 'f', 0, 32)

	// c.Data["TotalPrice"] = TotalPrice

	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "TotalPrice - ", TotalPrice)

	return
}
