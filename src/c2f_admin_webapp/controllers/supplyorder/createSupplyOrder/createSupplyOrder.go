package createSupplyOrder

import (
	"c2f_admin_webapp/session"
	"runtime/debug"
	"strconv"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/util/txnno"

	"c2f_admin_webapp/model/utils"
	"errors"
	"strings"

	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/services"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"
	"c2f_admin_webapp/model/util/log"
)

type CreateSupplyOrder struct {
	// beego.Controller
	middleware.BaseController
}

type Display2 struct {
	Fields2 []Field2
}
type Field2 struct {
	Id              string
	Mobile          string
	Name            string
	ProductName     string
	ProductPackage  string
	ProductQuantity string
}

func (c *CreateSupplyOrder) Get() {
	AdminId := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Add Assets Start")
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
			c.TplName = "supplyorder/createSupplyOrder/createSupplyOrder.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "supplyorder/createSupplyOrder/createSupplyOrder.html"
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

	id := sess.Get("uid").(string)
	c.Data["id"] = id

	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "utype", utype)

	auth, err := utils.IsAuthorized(utype, "OrdertoSupplyExe")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	if utype == "CLUSTERMANAGER" {

		data1, err1 := services.GetClusterAsPerMapping(id)
		if err1 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
			err1 = errors.New("Cluster user fetch Failed")
			return
		}
		if len(data1) <= 0 {
			err1 = errors.New("Cluster user  Not Found")
			return
		}
		data2, err := services.GetClusterMappedProducers(data1[0][1])
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Producer fetch Failed")
			return
		}
		if len(data2) <= 0 {
			err = errors.New("Producer  Not Found")
			return
		}
		var Dis2 Display2
		for i := range data2 {
			var d Field2

			d.Mobile = data2[i][3]
			d.Name = data2[i][4]
			d.ProductName = data2[i][2]
			d.ProductPackage = data2[i][6]
			d.ProductQuantity = data2[i][7]

			Dis2.Fields2 = append(Dis2.Fields2, d)
		}
		c.Data["Dis2"] = Dis2
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis2)

	} else if utype == "CLUSTEREXECUTIVE" {

		data1, err1 := services.GetClusterAsPerMappingExec(id)
		if err1 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
			err1 = errors.New("Cluster user fetch Failed")
			return
		}
		if len(data1) <= 0 {
			err1 = errors.New("Cluster user  Not Found")
			return
		}
		data2, err := services.GetClusterMappedProducers(data1[0][1])
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Producer fetch Failed")
			return
		}
		if len(data2) <= 0 {
			err = errors.New("Producer  Not Found")
			return
		}
		var Dis2 Display2
		for i := range data2 {
			var d Field2

			d.Mobile = data2[i][3]
			d.Name = data2[i][4]
			d.ProductName = data2[i][2]
			d.ProductPackage = data2[i][6]
			d.ProductQuantity = data2[i][7]

			Dis2.Fields2 = append(Dis2.Fields2, d)
		}
		c.Data["Dis2"] = Dis2
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis2)
	} else {
		data2, err := services.GetProducersinfo()
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Producer fetch Failed")
			return
		}
		if len(data2) <= 0 {
			err = errors.New("Producer  Not Found")
			return
		}
		var Dis2 Display2
		for i := range data2 {
			var d Field2

			d.Mobile = data2[i][3]
			d.Name = data2[i][4]
			d.ProductName = data2[i][2]
			d.ProductPackage = data2[i][6]
			d.ProductQuantity = data2[i][7]

			Dis2.Fields2 = append(Dis2.Fields2, d)
		}
		c.Data["Dis2"] = Dis2
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis2)

	}

	// product_id, err := services.GetClusterid(AdminId)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Cluster Id fetch Failed")
	// 	return
	// }

	// cluster_name, err := services.GetClusterName(product_id)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Cluster Name fetch Failed")
	// 	return
	// }

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
    Lookup_products.product_sku
FROM
    Product2Cluster
        LEFT JOIN
    Lookup_products ON Product2Cluster.product_id = Lookup_products.id
      LEFT JOIN Clusters 
ON  Product2Cluster.cluster_id= Clusters.id
 WHERE cluster_product_id =?  limit 1`, AdminId)

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
		c.Data["Product_QualityMeasure"] = data[i][8]
		c.Data["ShelLife"] = data[i][9]
		c.Data["PackageType"] = data[i][10]
		c.Data["PackageQty"] = data[i][11]
		c.Data["Krishivala_BuyingPrice"] = data[i][12]
		c.Data["Cluster_BuyingPrice"] = data[i][13]
		c.Data["Cluster_Name"] = data[i][14]
		c.Data["ProductSKU"] = data[i][15]

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "data[j][15] - ", data[i][15])

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
		}

	}

	return

}
func (c *CreateSupplyOrder) Post() {

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
			c.TplName = "supplyorder/createSupplyOrder/createSupplyOrder.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Supply Order Created Successfully"
			c.TplName = "supplyorder/createSupplyOrder/createSupplyOrder.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer  Page Success")
		}
		return
	}()
	utils.SetHTTPHeader(c.Ctx)
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

	auth, err := utils.IsAuthorized(utype, "OrdertoSupplyExe")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")
	input_category_name := c.Input().Get("input_category_name")
	input_subcategory_name := c.Input().Get("input_subcategory_name")
	select_product_name := c.Input().Get("select_product_name")

	input_product_sku := c.Input().Get("input_product_sku")
	input_product_qty := c.Input().Get("input_product_qty")
	input_QualityType := c.Input().Get("input_QualityType")
	input_QualityMeasure := c.Input().Get("input_QualityMeasure")

	input_ShelLife := c.Input().Get("input_ShelLife")
	input_PackageType := c.Input().Get("input_PackageType")
	input_PackageQty := c.Input().Get("input_PackageQty")

	select_producer_name := c.Input().Get("select_producer_name")
	input_cluster_name := c.Input().Get("input_cluster_name")
	input_cluster_buying_price := c.Input().Get("input_cluster_buying_price")

	product_supply_units := c.Input().Get("product_supply_units")

	//select_producer_name := c.Input().Get("select_producer_name")
	// /log.Println(beego.AppConfig.String("loglevel"), "Debug", "select_producer_name", select_producer_name)

	input_cluster_id := c.Input().Get("input_cluster_id")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_cluster_id", input_cluster_id)

	bp, err := strconv.ParseInt(input_cluster_buying_price, 0, 64)
	u, err := strconv.ParseInt(product_supply_units, 0, 64)
	BP := bp * u
	totalprice := strconv.FormatInt(BP, 10)

	product_supply_total_price := totalprice

	supplyorder_txn_id := txnno.Generate13Digit()

	input_product_id := c.Input().Get("input_product_id")
	//input_cluster_id := c.Input().Get("input_cluster_id")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_product_id", input_product_id)
	input_status := c.Input().Get("input_status")

	if select_producer_name == "" || product_supply_units == "" || input_status == "" {
		err = errors.New("* fields cannot be empty")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "select_product_name", select_product_name)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "select_producer_name", select_producer_name)
	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "select_cluster_name", select_cluster_name)
	err = CheckMapExists(select_product_name, select_producer_name, input_cluster_name)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Product-Producer-Cluster Mapping Not Exists")
		return
	}

	err = CheckExists1(input_product_id, select_producer_name, input_cluster_id)
	//	currentdate := time.Now()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Product Forecast is not available")
		return
	}

	// err = CheckExists(input_order_id)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Order already Exists")
	// 	return
	// }

	result, err := db.Db.Exec(`INSERT INTO SupplyOrderTxn (
	supplyordertxn_id,
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
	producer_mob,
	
	cluster_name, 
	cluster_buying_price,
    product_supply_units,
	product_supply_total_price,
	
	status,
	product_id, 
    create_by,
	create_date
	)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,now())`,
		supplyorder_txn_id,
		input_category_name,
		input_subcategory_name,
		select_product_name,

		input_product_sku,
		input_product_qty,
		input_QualityType,
		input_QualityMeasure,

		input_ShelLife,
		input_PackageType,
		input_PackageQty,
		select_producer_name,

		input_cluster_name,
		input_cluster_buying_price,
		product_supply_units,
		product_supply_total_price,

		input_status,
		input_product_id,

		id)
	if err != nil {
		err = errors.New("Supply Order creation failed")
		return
	}

	i, err := result.RowsAffected()
	if err != nil || i == 0 {
		err = errors.New("Supply Order creation failed")
		return
	}

	return
}

func CheckMapExists(product, producer, cluster string) (err error) {

	err = nil

	row, err := db.Db.Query(`select count(*) from Product2Producer2Cluster WHERE product_id in (select id from Lookup_products where product_name=?) AND producer_id in (select id from Customer where mob_no=?) AND cluster_id in (select id from Clusters where cluster_name=?) AND status="ACTIVE"`, product, producer, cluster)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch user")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch user")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	countlen := data[0][0]

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "countlen", countlen)

	if countlen == "0" {
		err = errors.New("Mappig Not exists")
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}

	return

}

func CheckExists(order string) (err error) {

	err = nil

	row, err := db.Db.Query(`select count(*) from SupplyOrderTxn WHERE order_id = ?`, order)

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

	countlen := data[0][0]

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "countlen", countlen)

	if countlen != "0" {
		err = errors.New("Order already exists")
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}

	return

}

func CheckExists1(product1, producer1, cluster1 string) (err error) {

	err = nil

	row, err := db.Db.Query(`SELECT count(*)
 FROM SupplyForecastTxn where product_id=? AND producer_mob=? AND cluster_id=? AND end_date >= current_date()`, product1, producer1, cluster1)

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

	countlen := data[0][0]

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "countlen", countlen)

	if countlen == "0" {
		err = errors.New("Product Forecast is not available")
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}

	return

}
