package selectSupplyOrder

import (
	"c2f_admin_webapp/session"

	"runtime/debug"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"
	"errors"
	"strings"

	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/services"
	"strconv"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"
	"c2f_admin_webapp/model/util/log"
)

type SelectSupplyOrder struct {
	// beego.Controller
	middleware.BaseController
}

type Display struct {
	Fields []Field
}

type Display1 struct {
	Fields1 []Field1
}

type Display2 struct {
	Fields2 []Field2
}

type Display3 struct {
	Fields3 []Field3
}
type Field struct {
	Id    string
	Name  string
	Email string
}

type Field1 struct {
	Id   string
	Name string
}

type Field2 struct {
	Id     string
	Mobile string
}

type Field3 struct {
	Id    string
	Units string
}

type Row struct {
	Id                   string
	CategoryName         string
	SubCategoryName      string
	Product_name         string
	Cluster_SellingPrice string
	Cluster_BuyingPrice  string
	ClusterName          string
	Product_image        string
	ProductQty           string
	ProductpackageQty    string
	Image                string
}

func (c *SelectSupplyOrder) Get() {
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
			c.TplName = "supplyorder/selectSupplyOrder/selectSupplyOrder.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "supplyorder/selectSupplyOrder/selectSupplyOrder.html"
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
	//id, category_name, subcategory_name, product_name, poduct_qty, product_desc, product_status, created_by, product_price, product_image, product_sku, create_date, last_update

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

		row, err := db.Db.Query(` SELECT 
	Product2Cluster.cluster_product_id,
    Product2Cluster.cluster_id,
    Product2Cluster.product_id,
    Product2Cluster.status,
    Lookup_products.category_name,
	Lookup_products.subcategory_name,
    Lookup_products.product_name,
    Product2Cluster.cluster_selling_price,
    Product2Cluster.cluster_buying_price,
    Product2Cluster.cluster_market_price,
    Clusters.cluster_name,
    Lookup_products.product_qty,
    Lookup_products.product_package_qty,
     JSON_EXTRACT(images_path,'$.img1')
FROM
    Product2Cluster
        LEFT JOIN
    Lookup_products ON Product2Cluster.product_id = Lookup_products.id
      LEFT JOIN Clusters 
ON  Product2Cluster.cluster_id= Clusters.id
      WHERE Product2Cluster.cluster_id=? AND Lookup_products.product_status="ACTIVE"`, data1[0][1])

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch order")
			return
		}
		defer sql.Close(row)
		_, data4, err := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch order")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", data4)
		var result []Row

		for i := range data4 {
			var r Row
			r.Id = data4[i][0]
			r.CategoryName = data4[i][4]
			r.SubCategoryName = data4[i][5]
			r.Product_name = data4[i][6]
			r.ClusterName = data4[i][10]
			r.Cluster_BuyingPrice = data4[i][8]
			r.ProductpackageQty = data4[i][12] + data4[i][11]
			r.Image, err = strconv.Unquote(data4[i][13])

			result = append(result, r)
		}
		c.Data["CustomerData"] = result

	} else if utype == "CLUSTEREXECUTIVE" {

		data2, err2 := services.GetClusterAsPerMappingExec(id)
		if err2 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err2)
			err2 = errors.New("Cluster user fetch Failed")
			return
		}
		if len(data2) <= 0 {
			err2 = errors.New("Cluster user  Not Found")
			return
		}

		row, err := db.Db.Query(` SELECT 
	Product2Cluster.cluster_product_id,
    Product2Cluster.cluster_id,
    Product2Cluster.product_id,
    Product2Cluster.status,
    Lookup_products.category_name,
	Lookup_products.subcategory_name,
    Lookup_products.product_name,
    Product2Cluster.cluster_selling_price,
    Product2Cluster.cluster_buying_price,
    Product2Cluster.cluster_market_price,
    Clusters.cluster_name,
    Lookup_products.product_qty,
    Lookup_products.product_package_qty,
     JSON_EXTRACT(images_path,'$.img1')
FROM
    Product2Cluster
        LEFT JOIN
    Lookup_products ON Product2Cluster.product_id = Lookup_products.id
      LEFT JOIN Clusters 
ON  Product2Cluster.cluster_id= Clusters.id
      WHERE Clusters.cluster_name=? AND Lookup_products.product_status="ACTIVE"`, data2[0][3])

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch order")
			return
		}
		defer sql.Close(row)
		_, data4, err := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch order")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", data4)
		var result []Row

		for i := range data4 {
			var r Row
			r.Id = data4[i][0]
			r.CategoryName = data4[i][4]
			r.SubCategoryName = data4[i][5]
			r.Product_name = data4[i][6]
			r.ClusterName = data4[i][10]
			r.Cluster_BuyingPrice = data4[i][8]
			r.ProductpackageQty = data4[i][12] + data4[i][11]
			r.Image, err = strconv.Unquote(data4[i][13])

			result = append(result, r)
		}
		c.Data["CustomerData"] = result
	} else {

		row, err := db.Db.Query(` SELECT 
	Product2Cluster.cluster_product_id,
    Product2Cluster.cluster_id,
    Product2Cluster.product_id,
    Product2Cluster.status,
    Lookup_products.category_name,
	Lookup_products.subcategory_name,
    Lookup_products.product_name,
    Product2Cluster.cluster_selling_price,
    Product2Cluster.cluster_buying_price,
    Product2Cluster.cluster_market_price,
    Clusters.cluster_name,
    Lookup_products.product_qty,
    Lookup_products.product_package_qty,
     JSON_EXTRACT(images_path,'$.img1')
FROM
    Product2Cluster
        LEFT JOIN
    Lookup_products ON Product2Cluster.product_id = Lookup_products.id
      LEFT JOIN Clusters 
ON  Product2Cluster.cluster_id= Clusters.id
      WHERE Lookup_products.product_status="ACTIVE"`)

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch order")
			return
		}
		defer sql.Close(row)
		_, data4, err := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch order")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", data4)
		var result []Row

		for i := range data4 {
			var r Row
			r.Id = data4[i][0]
			r.CategoryName = data4[i][4]
			r.SubCategoryName = data4[i][5]
			r.Product_name = data4[i][6]
			r.ClusterName = data4[i][10]
			r.Cluster_BuyingPrice = data4[i][8]
			r.ProductpackageQty = data4[i][12] + data4[i][11]
			r.Image, err = strconv.Unquote(data4[i][13])

			result = append(result, r)
		}
		c.Data["CustomerData"] = result

	}

	return
}
func (c *SelectSupplyOrder) Post() {

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
			c.TplName = "supplyorder/selectSupplyOrder/selectSupplyOrder.html"

			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Supply Forcast Done Successfully"
			c.TplName = "supplyorder/selectSupplyOrder/selectSupplyOrder.html"

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
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "id", id)

	auth, err := utils.IsAuthorized(utype, "OrdertoSupplyExe")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	select_product_name := c.Input().Get("select_product_name")
	//select_producer_name := c.Input().Get("select_producer_name")
	select_cluster_name := c.Input().Get("select_cluster_name")
	//input_quantity := c.Input().Get("input_quantity")
	//	input_units := c.Input().Get("input_units")
	//input_price := c.Input().Get("input_price")
	//input_start_date := c.Input().Get("input_start_date")

	//input_end_date := c.Input().Get("input_end_date")

	// err = CheckMapExists(select_product_name, select_producer_name, select_cluster_name)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Product Mapping Not Exists")
	// 	return
	// }

	// err = CheckExists(select_product_name, select_producer_name, select_cluster_name, input_quantity, input_price, input_start_date, input_end_date)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Order already Exists")
	// 	return
	// }

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
		row, err := db.Db.Query(` SELECT 
	Product2Cluster.cluster_product_id,
    Product2Cluster.cluster_id,
    Product2Cluster.product_id,
    Product2Cluster.status,
    Lookup_products.category_name,
	Lookup_products.subcategory_name,
    Lookup_products.product_name,
    Product2Cluster.cluster_selling_price,
    Product2Cluster.cluster_buying_price,
    Product2Cluster.cluster_market_price,
    Clusters.cluster_name,
    Lookup_products.product_qty,
    Lookup_products.product_package_qty,
     JSON_EXTRACT(images_path,'$.img1')
FROM
    Product2Cluster
        LEFT JOIN
    Lookup_products ON Product2Cluster.product_id = Lookup_products.id
      LEFT JOIN Clusters 
ON  Product2Cluster.cluster_id= Clusters.id
       where Product2Cluster.cluster_id=? AND (Lookup_products.product_name='' OR Lookup_products.product_name like ?) AND (Clusters.cluster_name='' OR Clusters.cluster_name like ?) AND Lookup_products.product_status="ACTIVE"`, data1[0][1], select_product_name+"%", select_cluster_name+"%")

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch order")
			return
		}
		defer sql.Close(row)
		_, data4, err := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch order")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", data4)
		var result []Row

		for i := range data4 {
			var r Row
			r.Id = data4[i][0]
			r.CategoryName = data4[i][4]
			r.SubCategoryName = data4[i][5]
			r.Product_name = data4[i][6]
			r.ClusterName = data4[i][10]
			r.Cluster_BuyingPrice = data4[i][8]
			r.ProductpackageQty = data4[i][12] + data4[i][11]
			r.Image, err = strconv.Unquote(data4[i][13])

			result = append(result, r)
		}
		c.Data["CustomerData"] = result

	} else if utype == "CLUSTEREXECUTIVE" {

		data2, err2 := services.GetClusterAsPerMappingExec(id)
		if err2 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err2)
			err2 = errors.New("Cluster user fetch Failed")
			return
		}
		if len(data2) <= 0 {
			err2 = errors.New("Cluster user  Not Found")
			return
		}

		row, err := db.Db.Query(` SELECT 
	Product2Cluster.cluster_product_id,
    Product2Cluster.cluster_id,
    Product2Cluster.product_id,
    Product2Cluster.status,
    Lookup_products.category_name,
	Lookup_products.subcategory_name,
    Lookup_products.product_name,
    Product2Cluster.cluster_selling_price,
    Product2Cluster.cluster_buying_price,
    Product2Cluster.cluster_market_price,
    Clusters.cluster_name,
    Lookup_products.product_qty,
    Lookup_products.product_package_qty,
     JSON_EXTRACT(images_path,'$.img1')
FROM
    Product2Cluster
        LEFT JOIN
    Lookup_products ON Product2Cluster.product_id = Lookup_products.id
      LEFT JOIN Clusters 
ON  Product2Cluster.cluster_id= Clusters.id
       where Clusters.cluster_name=? AND (Lookup_products.product_name='' OR Lookup_products.product_name like ?) AND (Clusters.cluster_name='' OR Clusters.cluster_name like ?) AND Lookup_products.product_status="ACTIVE"`, data2[0][3], select_product_name+"%", select_cluster_name+"%")

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch order")
			return
		}
		defer sql.Close(row)
		_, data4, err := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch order")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", data4)
		var result []Row

		for i := range data4 {
			var r Row
			r.Id = data4[i][0]
			r.CategoryName = data4[i][4]
			r.SubCategoryName = data4[i][5]
			r.Product_name = data4[i][6]
			r.ClusterName = data4[i][10]
			r.Cluster_BuyingPrice = data4[i][8]
			r.ProductpackageQty = data4[i][12] + data4[i][11]
			r.Image, err = strconv.Unquote(data4[i][13])

			result = append(result, r)
		}
		c.Data["CustomerData"] = result
	} else {

		row, err := db.Db.Query(` SELECT 
	Product2Cluster.cluster_product_id,
    Product2Cluster.cluster_id,
    Product2Cluster.product_id,
    Product2Cluster.status,
    Lookup_products.category_name,
	Lookup_products.subcategory_name,
    Lookup_products.product_name,
    Product2Cluster.cluster_selling_price,
    Product2Cluster.cluster_buying_price,
    Product2Cluster.cluster_market_price,
    Clusters.cluster_name,
    Lookup_products.product_qty,
    Lookup_products.product_package_qty,
     JSON_EXTRACT(images_path,'$.img1')
FROM
    Product2Cluster
        LEFT JOIN
    Lookup_products ON Product2Cluster.product_id = Lookup_products.id
      LEFT JOIN Clusters 
ON  Product2Cluster.cluster_id= Clusters.id
       where (Lookup_products.product_name='' OR Lookup_products.product_name like ?) AND (Clusters.cluster_name='' OR Clusters.cluster_name like ?) AND Lookup_products.product_status="ACTIVE"`, select_product_name+"%", select_cluster_name+"%")

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch order")
			return
		}
		defer sql.Close(row)
		_, data4, err := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch order")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", data4)
		var result []Row

		for i := range data4 {
			var r Row
			r.Id = data4[i][0]
			r.CategoryName = data4[i][4]
			r.SubCategoryName = data4[i][5]
			r.Product_name = data4[i][6]
			r.ClusterName = data4[i][10]
			r.Cluster_BuyingPrice = data4[i][8]
			r.ProductpackageQty = data4[i][12] + data4[i][11]
			r.Image, err = strconv.Unquote(data4[i][13])

			result = append(result, r)
		}
		c.Data["CustomerData"] = result

	}

	return
}
func CheckMapExists(product, producer, cluster string) (err error) {

	err = nil
	row, err := db.Db.Query(`select count(*) from Product2Producer2Cluster WHERE product_id = (select id from Lookup_products where product_name=?) AND producer_id=(select id from Customer where mob_no=?) AND cluster_id=(select id from Clusters where cluster_name=?)`, product, producer, cluster)

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
		err = errors.New("Mappig already exists")
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}

	return

}

func CheckExists(product_name, producer_name, cluster_name, quantity, price, sd, ed string) (err error) {

	err = nil

	row, err := db.Db.Query(`select count(*) from SupplyForecastTxn WHERE producer_mob=? AND  product_name=? AND cluster_name=? AND qty=? AND price=? AND start_date=? AND end_date=?`, producer_name, product_name, cluster_name, quantity, price, sd, ed)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch supply order")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch supply order")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	countlen := data[0][0]

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "countlen", countlen)

	if countlen != "0" {
		err = errors.New("Supply order already exists")
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}

	return

}
