package placeOrder

import (
	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/model/utils"
	"c2f_admin_webapp/services"
	"c2f_admin_webapp/session"
	"errors"

	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
	"c2f_admin_webapp/model/util/txnno"

	"github.com/astaxie/beego"
)

type Display2 struct {
	Fields2 []Field2
}
type Field2 struct {
	Id       string
	Mobile   string
	Name     string
	Cluster  string
	Address1 string
	Address2 string
	City     string
	State    string
	Zipcode  string
}
type PlaceOrder struct {
	// beego.Controller
	middleware.BaseController
}

type GetCustomerAdditionalAddress struct {
	// beego.Controller
	middleware.BaseController
}

func (c *PlaceOrder) Get() {

	AdminId := c.Ctx.Request.URL.Query().Get("AdminId")
	Para1 := c.Ctx.Request.URL.Query().Get("Para1")
	Para2 := c.Ctx.Request.URL.Query().Get("Para2")

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
			c.TplName = "customeroder/placeOrder/placeOrder.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "customeroder/placeOrder/placeOrder.html"
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
	auth, err := utils.IsAuthorized(utype, "PlaceOrder")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	// product_id, err := services.GetClusterid(AdminId)
	// if err != nil {
	//      log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	//      err = errors.New("Cluster Id fetch Failed")
	//      return
	// }
	// cluster_name, err := services.GetClusterName(product_id)
	// if err != nil {
	//      log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	//      err = errors.New("Cluster Name fetch Failed")
	//      return
	// }

	fmt.Println("AdminId :- ", AdminId)
	fmt.Println("Para1:- ", Para1)
	fmt.Println("Para2 :- ", Para2)

	product_id := Para1
	fmt.Println("Product Id :- ", product_id)
	cluster_id := Para2
	fmt.Println("Cluster Id :- ", cluster_id)

	// To check if Product is mapped to Cluster or Not
	row, err := db.Db.Query(`select count(*) from Product2Cluster WHERE cluster_id=? and product_id=?`, cluster_id, product_id)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch Mapped Product")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch  Mapped Product")
		return
	}

	countlen := data[0][0]

	if countlen != "0" {

		// For Mapped Products
		fmt.Println("Product is Mapped with Cluster")

		row, err := db.Db.Query(`
        SELECT 
        Product2Cluster.cluster_product_id,
    Product2Cluster.cluster_id,
    Product2Cluster.product_id,
    Lookup_products.category_name,
        Lookup_products.subcategory_name,
    Lookup_products.product_name,
    Lookup_products.product_qty,
    Lookup_products.product_quality_type,
    Lookup_products.product_shelf_life,
     Lookup_products.product_quality_mesure,
    Lookup_products.product_package_type,
    Lookup_products.product_package_qty,
    Lookup_products.krishivala_buying_price,
    Product2Cluster.cluster_buying_price,
    Clusters.cluster_name,
    Lookup_products.product_sku,
    Product2Cluster.cluster_selling_price,
    Lookup_products.krishivala_CGST_charge,
    Lookup_products.krishivala_SGST_charge,
    Lookup_products.krishivala_CESS_charge,
    Lookup_products.krishivala_OtherTAX_charge,
    Product2Cluster.delivery_discount
    
FROM
    Product2Cluster
        LEFT JOIN
    Lookup_products ON Product2Cluster.product_id = Lookup_products.id
      LEFT JOIN Clusters 
ON  Product2Cluster.cluster_id= Clusters.id 
WHERE Product2Cluster.cluster_product_id =?  limit 1`, AdminId)
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
			c.Data["ClusterID"] = data[i][1]
			c.Data["ProductID"] = data[i][2]
			c.Data["CategoryName"] = data[i][3]
			c.Data["SubCategoryName"] = data[i][4]
			c.Data["ProductName"] = data[i][5]
			c.Data["ProductQty"] = data[i][6]
			c.Data["ProductQualityType"] = data[i][7]
			c.Data["Product_QualityMeasure"] = data[i][9]
			c.Data["ShelLife"] = data[i][8]
			c.Data["PackageType"] = data[i][10]
			c.Data["PackageQty"] = data[i][11]
			c.Data["Krishivala_BuyingPrice"] = data[i][12]
			c.Data["Cluster_BuyingPrice"] = data[i][13]
			c.Data["Cluster_Name"] = data[i][14]
			c.Data["ProductSKU"] = data[i][15]
			c.Data["Cluster_SellingPrice"] = data[i][16]
			c.Data["krishivala_CGST_charge"] = data[i][17]
			c.Data["krishivala_SGST_charge"] = data[i][18]
			c.Data["krishivala_CESS_charge"] = data[i][19]
			c.Data["krishivala_OtherTAX_charge"] = data[i][20]
			c.Data["DeliveryDiscount"] = data[i][21]
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "data[j][15] - ", data[i][15])

			data9, err := services.GetConsumer1(data[i][14])
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Customer fetch Failed")
				return
			}
			if len(data9) <= 0 {
				err = errors.New("customer  Not Found")
				return
			}
			var Dis2 Display2
			for i := range data9 {
				var d Field2
				d.Id = data9[i][0]
				d.Mobile = data9[i][1]
				d.Name = data9[i][2]
				d.Cluster = data9[i][3]
				d.Address1 = data9[i][4]
				d.Address2 = data9[i][5]
				d.City = data9[i][6]
				d.State = data9[i][7]
				d.Zipcode = data9[i][8]
				Dis2.Fields2 = append(Dis2.Fields2, d)
			}
			c.Data["Dis2"] = Dis2
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis2)

			row2, err := db.Db.Query(`select id,id from Lookup_products WHERE product_sku =?`, data[i][15])
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
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len2 - ", len(data2))
			if len(data2) <= 0 {
				err = errors.New("User data not found")
				return
			}
			for j := range data2 {
				log.Println(beego.AppConfig.String("loglevel"), "Debug", "data[j][0] - ", data2[j][1])
				c.Data["ProductID"] = data2[j][1]
				row3, err := db.Db.Query(`select product_id,cluster_selling_price,cluster_supply_charge,cluster_packaging_charge,cluster_bank_charge,cluster_discount ,cluster_id from Product2Cluster WHERE product_id =? AND cluster_id=?`, data2[j][1], data[i][1])
				if err != nil {
					log.Println(beego.AppConfig.String("loglevel"), "Error", err)
					err = errors.New("Unable to get user data")
					return
				}
				defer sql.Close(row3)
				_, data3, err := sql.Scan(row3)
				if err != nil {
					log.Println(beego.AppConfig.String("loglevel"), "Error", err)
					err = errors.New("Unable to get user data")
					return
				}
				log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len2 - ", len(data3))
				if len(data3) <= 0 {
					err = errors.New("cluster selling data not found")
					return
				}
				for k := range data3 {
					c.Data["ClusterSellPrice"] = data3[k][1]
					c.Data["ClusterSupplyCharge"] = data3[k][2]
					c.Data["ClusterPackagingCharge"] = data3[k][3]
					c.Data["ClusterBankCharge"] = data3[k][4]
					c.Data["ClusterDiscount"] = data3[k][5]
				}
			}
		}

	} else {

		// For Non Mapped Products
		fmt.Println("Product is not Mapped with Cluster")

		row, err := db.Db.Query(`
        SELECT 
        Inventory.id,
    Clusters.id,
    Inventory.product_id,
    Lookup_products.category_name,
        Lookup_products.subcategory_name,
    Lookup_products.product_name,
    Lookup_products.product_qty,
    Lookup_products.product_quality_type,
    Lookup_products.product_shelf_life,
     Lookup_products.product_quality_mesure,
    Lookup_products.product_package_type,
    Lookup_products.product_package_qty,
    Lookup_products.krishivala_buying_price,
    Inventory.cluster_buying_price,
    Clusters.cluster_name,
    Lookup_products.product_sku,
    Inventory.cluster_selling_price,
    Lookup_products.krishivala_CGST_charge,
    Lookup_products.krishivala_SGST_charge,
    Lookup_products.krishivala_CESS_charge,
    Lookup_products.krishivala_OtherTAX_charge
    
FROM
    Inventory
        LEFT JOIN
    Lookup_products ON Inventory.product_id = Lookup_products.id
      LEFT JOIN Clusters 
ON  Inventory.cluster_name= Clusters.cluster_name
WHERE Inventory.id =?  limit 1`, AdminId)
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
			c.Data["ClusterID"] = data[i][1]
			c.Data["ProductID"] = data[i][2]
			c.Data["CategoryName"] = data[i][3]
			c.Data["SubCategoryName"] = data[i][4]
			c.Data["ProductName"] = data[i][5]
			c.Data["ProductQty"] = data[i][6]
			c.Data["ProductQualityType"] = data[i][7]
			c.Data["Product_QualityMeasure"] = data[i][9]
			c.Data["ShelLife"] = data[i][8]
			c.Data["PackageType"] = data[i][10]
			c.Data["PackageQty"] = data[i][11]
			c.Data["Krishivala_BuyingPrice"] = data[i][12]
			c.Data["Cluster_BuyingPrice"] = data[i][13]
			c.Data["Cluster_Name"] = data[i][14]
			c.Data["ProductSKU"] = data[i][15]
			c.Data["Cluster_SellingPrice"] = data[i][16]
			c.Data["krishivala_CGST_charge"] = data[i][17]
			c.Data["krishivala_SGST_charge"] = data[i][18]
			c.Data["krishivala_CESS_charge"] = data[i][19]
			c.Data["krishivala_OtherTAX_charge"] = data[i][20]
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "data[j][15] - ", data[i][15])

			data9, err := services.GetConsumer1(data[i][14])
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Customer fetch Failed")
				return
			}
			if len(data9) <= 0 {
				err = errors.New("customer  Not Found")
				return
			}
			var Dis2 Display2
			for i := range data9 {
				var d Field2
				d.Id = data9[i][0]
				d.Mobile = data9[i][1]
				d.Name = data9[i][2]
				d.Cluster = data9[i][3]
				d.Address1 = data9[i][4]
				d.Address2 = data9[i][5]
				d.City = data9[i][6]
				d.State = data9[i][7]
				d.Zipcode = data9[i][8]
				Dis2.Fields2 = append(Dis2.Fields2, d)
			}
			c.Data["Dis2"] = Dis2
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis2)

			row2, err := db.Db.Query(`select id,id from Lookup_products WHERE product_sku =?`, data[i][15])
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
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len2 - ", len(data2))
			if len(data2) <= 0 {
				err = errors.New("User data not found")
				return
			}
			for j := range data2 {
				log.Println(beego.AppConfig.String("loglevel"), "Debug", "data[j][0] - ", data2[j][1])
				c.Data["ProductID"] = data2[j][1]
				row3, err := db.Db.Query(`select product_id,cluster_selling_price,cluster_supply_charge,cluster_packaging_charge,cluster_bank_charge,cluster_discount  from Inventory WHERE product_id =? AND cluster_name=?`, data2[j][1], data[i][14])
				if err != nil {
					log.Println(beego.AppConfig.String("loglevel"), "Error", err)
					err = errors.New("Unable to get user data")
					return
				}
				defer sql.Close(row3)
				_, data3, err := sql.Scan(row3)
				if err != nil {
					log.Println(beego.AppConfig.String("loglevel"), "Error", err)
					err = errors.New("Unable to get user data")
					return
				}
				log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len2 - ", len(data3))
				if len(data3) <= 0 {
					err = errors.New("cluster selling data not found")
					return
				}
				for k := range data3 {
					c.Data["ClusterSellPrice"] = data3[k][1]
					c.Data["ClusterSupplyCharge"] = data3[k][2]
					c.Data["ClusterPackagingCharge"] = data3[k][3]
					c.Data["ClusterBankCharge"] = data3[k][4]
					c.Data["ClusterDiscount"] = data3[k][5]
				}

			}
		}

	}

	return
}
func (c *PlaceOrder) Post() {
	//AdminId := c.Ctx.Input.Param(":AdminID")
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
			c.TplName = "customeroder/placeOrder/placeOrder.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Product Order Successfully"
			c.TplName = "customeroder/placeOrder/placeOrder.html"
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

	id := sess.Get("uid").(string)
	c.Data["id"] = id
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "id", id)

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

	auth, err := utils.IsAuthorized(utype, "PlaceOrder")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")
	input_category_name := c.Input().Get("input_category_name")
	input_subcategory_name := c.Input().Get("input_subcategory_name")
	input_product_name := c.Input().Get("input_product_name")
	input_product_sku := c.Input().Get("input_product_sku")
	input_product_qty := c.Input().Get("input_product_qty")
	input_QualityType := c.Input().Get("input_QualityType")
	input_QualityMeasure := c.Input().Get("input_QualityMeasure")
	input_ShelLife := c.Input().Get("input_ShelLife")
	input_PackageType := c.Input().Get("input_PackageType")
	input_PackageQty := c.Input().Get("input_PackageQty")
	input_cluster_selling_price := c.Input().Get("input_cluster_selling_price")

	var input_deliveryfrom_date, input_deliveryto_date string

	select_Customer := c.Input().Get("select_Customer")
	input_order_type := c.Input().Get("input_order_type")

	input_Delivery_Time := c.Input().Get("input_Delivery_Time")
	input_purchage_type := c.Input().Get("input_purchage_type")
	input_payment_type := c.Input().Get("input_payment_type")
	input_cluster_name := c.Input().Get("input_cluster_name")

	var input_howmany_packets string

	input_howmany_packets1 := c.Input().Get("input_howmany_packets")
	input_qty := c.Input().Get("input_qty")
	delivery_discount := c.Input().Get("delivery_discount")

	// new line: text that must appear on invoice
	deliveryInformation := c.Input().Get("delivery_to")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "deliveryInformation - ", deliveryInformation)

	if select_Customer == "" || input_order_type == "" || input_Delivery_Time == "" || input_purchage_type == "" || input_payment_type == "" {
		err = errors.New("* fields cannot be empty")
		return
	}

	if input_product_qty != "kg-bulk" {

		if input_howmany_packets1 == "" {

			err = errors.New("Number of Packets  cannot be empty")
			return

		}

		input_howmany_packets = input_howmany_packets1

	}

	if input_product_qty == "kg-bulk" {

		if input_qty == "" {

			err = errors.New("Quantity cannot be empty")
			return

		}

		qty, err := strconv.ParseFloat(input_qty, 64)

		if err != nil {
			beego.Error(err)
			err = errors.New("Invalid qty  ")
			return
		}
		fmt.Println("no of kgs :", qty)
		if qty < 10 {

			err = errors.New("Quantity cannot be less than 10 kg")
			return

		}

		input_howmany_packets = input_qty

	}

	//  BULK CHECK FOR GMS PRODUCTS (after packets validation)
	if input_product_qty == "gms" {

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "inside gms - ")

		// input_PackageQty contains packet size: 500, 1000, etc.
		gramsPerPacket, _ := strconv.ParseFloat(input_PackageQty, 64)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "gramsPerPacket - ", gramsPerPacket)

		// number of packets entered by user
		packets, _ := strconv.ParseFloat(input_howmany_packets1, 64)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "packets - ", packets)

		totalGrams := packets * gramsPerPacket
		totalKG := totalGrams / 1000.0

		if totalKG >= 10 {

			log.Println(beego.AppConfig.String("loglevel"), "Debug", "inside packets - ", totalKG)

			// check if same category + subcategory + cluster exists with kg-bulk
			row, _ := db.Db.Query(`
            SELECT COUNT(*) 
            FROM Lookup_products
            left join Product2Cluster
			on Product2Cluster.product_id=Lookup_products.id
			left join Clusters
			on Clusters.id=Product2Cluster.cluster_id
            WHERE category_name = ?
              AND subcategory_name = ?
              AND cluster_name = ?
              AND product_qty = 'kg-bulk' And product_status = 'ACTIVE'
            LIMIT 1`,
				input_category_name,
				input_subcategory_name,
				input_cluster_name,
			)

			_, d, _ := sql.Scan(row)
			hasBulk := (d[0][0] == "1")

			if hasBulk {
				err = errors.New("This is a more than 10kg. Please place order using KG-BULK product.")
				return
			}
		}
	}

	var deliverycharges string

	if input_order_type == "SUBSCRIPTION" {
		//currentTime := time.Now()
		input_deliveryfrom_date = c.Input().Get("input_deliveryfrom_date")
		input_deliveryto_date = c.Input().Get("input_deliveryto_date")
		totalday := (daysBetween(date(input_deliveryfrom_date), date(input_deliveryto_date)))
		totaldays := strconv.Itoa(totalday)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "totaldays - ", totaldays)
		input_repeat_order := c.Input().Get("input_repeat_order")

		if input_deliveryfrom_date == "" || input_deliveryto_date == "" || input_repeat_order == "" {
			err = errors.New("FromDate & ToDate & RepeatOrder fields cannot be Empty !!")
			return
		}

		oderid := txnno.Generate()
		input_cluster_name := c.Input().Get("input_cluster_name")
		customer_id, _ := services.GetConsumerid(select_Customer)
		if err != nil {
			beego.Error(err)
			err = errors.New("Invalid Customer  ")
			return
		}
		cluster_supply_charge := c.Input().Get("cluster_supply_charge")

		fmt.Println("checking value of cluster_supply_charge", cluster_supply_charge)
		cluster_packaging_charge := c.Input().Get("cluster_packaging_charge")
		cluster_bank_charge := c.Input().Get("cluster_bank_charge")
		cluster_discount := c.Input().Get("cluster_discount")
		krishivala_CGST_charge := c.Input().Get("krishivala_CGST_charge")
		krishivala_SGST_charge := c.Input().Get("krishivala_SGST_charge")
		krishivala_CESS_charge := c.Input().Get("krishivala_CESS_charge")
		krishivala_OtherTAX_charge := c.Input().Get("krishivala_OtherTAX_charge")
		sp, err := strconv.ParseFloat(input_cluster_selling_price, 64)
		u, err := strconv.ParseFloat(input_howmany_packets, 64)
		td, err := strconv.ParseFloat(totaldays, 64)
		TDDD := td + 1
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "TDDD - ", TDDD)
		SP := sp * u * TDDD
		totalprice := strconv.FormatFloat(SP, 'f', 0, 32)
		updated_total_product_price := totalprice
		scluster_supply_charge, err := strconv.ParseFloat(cluster_supply_charge, 64)
		fmt.Println("checking value of scluster_supply_charge", scluster_supply_charge)
		scluster_packaging_charge, err := strconv.ParseFloat(cluster_packaging_charge, 64)
		scluster_supply_charge1 := scluster_supply_charge * TDDD
		fmt.Println("checking value of scluster_supply_charge1", scluster_supply_charge1)
		scluster_packaging_charge1 := scluster_packaging_charge * TDDD
		scluster_bank_charge, err := strconv.ParseFloat(cluster_bank_charge, 64)
		scluster_discount, err := strconv.ParseFloat(cluster_discount, 64)
		skrishivala_CGST_charge, err := strconv.ParseFloat(krishivala_CGST_charge, 64)
		skrishivala_SGST_charge, err := strconv.ParseFloat(krishivala_SGST_charge, 64)
		skrishivala_CESS_charge, err := strconv.ParseFloat(krishivala_CESS_charge, 64)
		skrishivala_OtherTAX_charge, err := strconv.ParseFloat(krishivala_OtherTAX_charge, 64)
		charges := scluster_supply_charge1 + scluster_packaging_charge1 + scluster_bank_charge
		taxs := skrishivala_CGST_charge + skrishivala_SGST_charge + skrishivala_CESS_charge + skrishivala_OtherTAX_charge
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "charges - ", charges)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "taxs - ", taxs)
		TPP := SP + (charges) - scluster_discount + (taxs)
		tpp := strconv.FormatFloat(TPP, 'f', 0, 32)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "select_Customer  - ", select_Customer)
		// if select_Customer == "" || input_payment_type == "" {
		// 	err = errors.New("* fields canncont be enpmty")
		// 	return
		// }
		//  TotalBuying_price = cluster.buyingprice  + (cluster.supplycarge + cluster.packagingcharge + cluster.bankcharge ) - cluster.discount  + (producer.gcst + producer.scst + producer.C29 + producer.othertax)
		TotalSelling_price := tpp
		deliverycharges := ""
		//34
		result, err := db.Db.Exec(`INSERT INTO CustomerOrderTxn (
         category_name, 
         subcategory_name,
         product_name,
         product_sku,
         product_qty,
        
         product_quality_type,
         product_quality_mesure,
         product_shelf_life,
         product_package_type, 
         product_package_qty,
        
         
         product_price,
         total_product_price,
         order_units,
         orderdeliver_fromdate,
         orderdeliver_todate,
        
         order_repeat,
         orderdeliver_time,
         order_purchagetype,
         order_id,
         cluster_name,
        
         customer_mob,
         customer_id, 
         status,
         cluster_supply_charge, 
     	 cluster_packaging_charge,
    
     	 cluster_bank_charge,
     	 cluster_discount,
     	 krishivala_CGST_charge,
     	 krishivala_SGST_charge, 
     	 krishivala_CESS_charge,
		
     	 krishivala_OtherTAX_charge, 
     	 total_selling_price,
    	 payment_status,
    	 payment_type,
    	 total_days,
		
    	 order_type,
         created_by,
         cluster_type,
         delivery_charges,
		 delivery_address_details,
         create_date,
         order_date
         )
        
                VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,now(),now())`,
			input_category_name,
			input_subcategory_name,
			input_product_name,
			input_product_sku,
			input_product_qty,

			input_QualityType,
			input_QualityMeasure,
			input_ShelLife,
			input_PackageType,
			input_PackageQty,

			input_cluster_selling_price,
			updated_total_product_price,
			input_howmany_packets,
			input_deliveryfrom_date,
			input_deliveryto_date,

			input_repeat_order,
			input_Delivery_Time,
			input_purchage_type,
			oderid,
			input_cluster_name,

			select_Customer,
			customer_id,
			"ORDER_RECEIVED",
			scluster_supply_charge1,
			scluster_packaging_charge1,

			cluster_bank_charge,
			cluster_discount,
			krishivala_CGST_charge,
			krishivala_SGST_charge,
			krishivala_CESS_charge,

			krishivala_OtherTAX_charge,
			TotalSelling_price,
			"PENDING",
			input_payment_type,
			TDDD,

			input_order_type,
			id,
			"SOURCE",
			deliverycharges,
			deliveryInformation)
		if err != nil {
			err = errors.New("Order txn creation failed")
			return
		}
		i, err := result.RowsAffected()
		if err != nil || i == 0 {
			err = errors.New("Order txn creation failed")
			return
		}
	} else {
		currentTime := time.Now()
		input_deliveryfrom_date = currentTime.Format("2006-01-02")
		input_deliveryto_date = currentTime.Format("2006-01-02")
		totalday := (daysBetween(date(input_deliveryfrom_date), date(input_deliveryto_date)))
		totaldays := strconv.Itoa(totalday)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_deliveryfrom_date - ", input_deliveryfrom_date)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_deliveryto_date - ", input_deliveryto_date)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "totaldays - ", totaldays)
		//	input_repeat_order := c.Input().Get("input_repeat_order")
		input_Delivery_Time := c.Input().Get("input_Delivery_Time")
		input_purchage_type := c.Input().Get("input_purchage_type")
		select_Customer := c.Input().Get("select_Customer")
		oderid := txnno.Generate()
		input_cluster_name := c.Input().Get("input_cluster_name")
		customer_id, _ := services.GetConsumerid(select_Customer)
		if err != nil {
			beego.Error(err)
			err = errors.New("Invalid Customer  ")
			return
		}
		cluster_supply_charge := c.Input().Get("cluster_supply_charge")
		cluster_packaging_charge := c.Input().Get("cluster_packaging_charge")
		cluster_bank_charge := c.Input().Get("cluster_bank_charge")
		cluster_discount := c.Input().Get("cluster_discount")
		krishivala_CGST_charge := c.Input().Get("krishivala_CGST_charge")
		krishivala_SGST_charge := c.Input().Get("krishivala_SGST_charge")
		krishivala_CESS_charge := c.Input().Get("krishivala_CESS_charge")
		krishivala_OtherTAX_charge := c.Input().Get("krishivala_OtherTAX_charge")
		sp, err := strconv.ParseFloat(input_cluster_selling_price, 64)
		u, err := strconv.ParseFloat(input_howmany_packets, 64)
		td, err := strconv.ParseFloat(totaldays, 64)
		TDDD := td + 1
		//		TDDD2 := strconv.FormatFloat(TDDD, 'f', 0, 32)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "TDDD - ", TDDD)
		SP := sp * u * TDDD
		totalprice := strconv.FormatFloat(SP, 'f', 0, 32)
		updated_total_product_price := totalprice
		scluster_supply_charge, err := strconv.ParseFloat(cluster_supply_charge, 64)
		scluster_packaging_charge, err := strconv.ParseFloat(cluster_packaging_charge, 64)
		scluster_supply_charge1 := scluster_supply_charge * TDDD
		scluster_packaging_charge1 := scluster_packaging_charge * TDDD
		scluster_bank_charge, err := strconv.ParseFloat(cluster_bank_charge, 64)
		scluster_discount, err := strconv.ParseFloat(cluster_discount, 64)
		skrishivala_CGST_charge, err := strconv.ParseFloat(krishivala_CGST_charge, 64)
		skrishivala_SGST_charge, err := strconv.ParseFloat(krishivala_SGST_charge, 64)
		skrishivala_CESS_charge, err := strconv.ParseFloat(krishivala_CESS_charge, 64)
		skrishivala_OtherTAX_charge, err := strconv.ParseFloat(krishivala_OtherTAX_charge, 64)
		charges := scluster_supply_charge1 + scluster_packaging_charge1 + scluster_bank_charge
		taxs := skrishivala_CGST_charge + skrishivala_SGST_charge + skrishivala_CESS_charge + skrishivala_OtherTAX_charge
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "charges - ", charges)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "taxs - ", taxs)
		// TPP := SP + (charges) - scluster_discount + (taxs)

		//Delivery discount added for kg-bulk
		sdelivery_discount, err := strconv.ParseFloat(delivery_discount, 64)
		TPP := SP + (charges) - scluster_discount + (taxs) - sdelivery_discount

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "sdelivery_discount - ", sdelivery_discount)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "TPP - ", TPP)

		var tpp string
		// var TPP1, delivery_charge_discount float64

		deliverycharges = ""

		tpp = strconv.FormatFloat(TPP, 'f', 0, 32)

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "select_Customer  - ", select_Customer)
		if select_Customer == "" || input_payment_type == "" {
			err = errors.New("* fields canncont be enpmty")
			return
		}
		//  TotalBuying_price = cluster.buyingprice  + (cluster.supplycarge + cluster.packagingcharge + cluster.bankcharge ) - cluster.discount  + (producer.gcst + producer.scst + producer.C29 + producer.othertax)
		TotalSelling_price := tpp
		//34
		result, err1 := db.Db.Exec(`INSERT INTO CustomerOrderTxn (
         category_name, 
         subcategory_name,
         product_name,
         product_sku,
         product_qty,
        
         product_quality_type,
         product_quality_mesure,
         product_shelf_life,
         product_package_type, 
         product_package_qty,
        
         
         product_price,
         total_product_price,
         order_units,
     orderdeliver_fromdate,
         orderdeliver_todate,
        
         order_repeat,
         orderdeliver_time,
         order_purchagetype,
         order_id,
         cluster_name,
        
         customer_mob,
         customer_id, 
         status,
         cluster_supply_charge, 
     cluster_packaging_charge,
    
     cluster_bank_charge,
     cluster_discount,
     krishivala_CGST_charge,
     krishivala_SGST_charge, 
     krishivala_CESS_charge,
     krishivala_OtherTAX_charge, 
     total_selling_price,
    payment_status,
    payment_type,
    total_days,
    order_type,
         created_by,
        cluster_type,
            delivery_charges,
			delivery_discount,
			delivery_address_details,
         create_date,
         order_date
         )
        
                VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,now(),now())`,
			input_category_name,
			input_subcategory_name,
			input_product_name,
			input_product_sku,
			input_product_qty,

			input_QualityType,
			input_QualityMeasure,
			input_ShelLife,
			input_PackageType,
			input_PackageQty,

			input_cluster_selling_price,
			updated_total_product_price,
			input_howmany_packets,
			input_deliveryfrom_date,
			input_deliveryto_date,

			"NO",
			input_Delivery_Time,
			input_purchage_type,
			oderid,
			input_cluster_name,

			select_Customer,
			customer_id,
			"ORDER_RECEIVED",
			scluster_supply_charge1,
			scluster_packaging_charge1,

			cluster_bank_charge,
			cluster_discount,
			krishivala_CGST_charge,
			krishivala_SGST_charge,
			krishivala_CESS_charge,

			krishivala_OtherTAX_charge,
			TotalSelling_price,
			"PENDING",
			input_payment_type,
			TDDD,

			input_order_type,
			id,
			"SOURCE",
			deliverycharges,
			delivery_discount,
			deliveryInformation)
		if err1 != nil {
			err1 = errors.New("Order txn creation failed")
			return
		}
		i, err := result.RowsAffected()
		if err != nil || i == 0 {
			err = errors.New("Order txn creation failed")
			return
		}
	}
	return
}

func (c *GetCustomerAdditionalAddress) Post() {

	// AdminId := c.Ctx.Input.Param(":AdminID")
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
			c.TplName = "customer/updateCustomer/updateCustomer.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "customer/updateCustomer/updateCustomer.html"
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

	auth, err := utils.IsAuthorized(utype, "PlaceOrder")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}

	customerIDStr := c.GetString("customer_id")
	mobile := c.GetString("mobile")

	customerID, err := strconv.Atoi(customerIDStr)
	if err != nil || mobile == "" {
		c.CustomAbort(400, "Invalid Inputs")
		return
	}

	// Fetch Main Address from Customer table
	mainRow, cerr := db.Db.Query(`
		SELECT first_name, mob_no, address1, town, city, state, zipcode, is_default 
		FROM Customer WHERE id = ? LIMIT 1`, customerID)
	if cerr != nil {
		c.CustomAbort(500, "Query Error")
		return
	}
	defer mainRow.Close()

	_, mainData, _ := sql.Scan(mainRow)

	var mainAddr map[string]string
	if len(mainData) > 0 {
		mainAddr = map[string]string{
			"Name":      mainData[0][0],
			"Mobile":    mainData[0][1],
			"Address1":  mainData[0][2],
			"Address2":  mainData[0][3],
			"City":      mainData[0][4],
			"State":     mainData[0][5],
			"Pincode":   mainData[0][6],
			"IsDefault": mainData[0][7],
		}
	}

	// Fetch Additional ACTIVE addresses
	addrs, derr := services.GetCustomerAddressesByIDAndMobile(customerID, mobile)
	if derr != nil {
		addrs = []map[string]string{} // ✔ FIXED
	}

	// Return JSON
	c.Data["json"] = map[string]interface{}{
		"success":         true,
		"main_address":    mainAddr,
		"additional_list": addrs,
	}
	c.ServeJSON()
}

func daysBetween(a, b time.Time) int {
	if a.After(b) {
		a, b = b, a
	}
	days := -a.YearDay()
	for year := a.Year(); year < b.Year(); year++ {
		days += time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC).YearDay()
	}
	days += b.YearDay()
	return days
}
func date(s string) time.Time {
	d, _ := time.Parse("2006-01-02", s)
	return d
}
