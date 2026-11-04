package searchProduct

import (
	"c2f_admin_webapp/session"

	"runtime/debug"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"
	"errors"

	"c2f_admin_webapp/services"

	"c2f_admin_webapp/model/db"

	"strings"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

type Row struct {
	Id                string
	Timestamp         string
	CategoryName      string
	SubCategoryName   string
	ProductName       string
	ProductDesc       string
	Productstatus     string
	ProductQty        string
	ProductpackageQty string
	ProductSKU        string
	CreatedBy         string
}

type SearchProduct struct {
	// beego.Controller
	middleware.BaseController
}

func (c *SearchProduct) Get() {
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
			c.TplName = "products/searchProduct/searchProduct.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "products/searchProduct/searchProduct.html"
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

	id := sess.Get("uid").(string)
	c.Data["id"] = id

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

	auth, err := utils.IsAuthorized(utype, "Product")
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

		row, err := db.Db.Query(`SELECT DISTINCT 
		Product2Cluster.cluster_product_id,
		Product2Cluster.cluster_id, 
		Product2Cluster.product_id, 
        Product2Cluster.status, 
        Product2Cluster.create_date,
        Clusters.cluster_name, 
        Lookup_products.product_name,
     	Lookup_products.product_package_qty,
     	Lookup_products.product_qty,
		Lookup_products.id,
		Lookup_products.category_name,
		Lookup_products.subcategory_name,
    	Lookup_products.product_desc,
		Lookup_products.product_status,
		Lookup_products.created_by,
		Lookup_products.product_sku, 
		Lookup_products.create_date
		FROM C2F_NEW.Product2Cluster
		LEFT JOIN Clusters 
		ON C2F_NEW.Product2Cluster.cluster_id = Clusters.id
		LEFT JOIN Lookup_products 
		ON C2F_NEW.Product2Cluster.product_id = Lookup_products.id where Product2Cluster.cluster_id=? ORDER BY Product2Cluster.create_date DESC`, data1[0][1])

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get the error message info")
			return
		}
		defer sql.Close(row)
		_, data, err := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Message Scan Fail")
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
		if len(data) <= 0 {
			err = errors.New("Unable to get the error message info")
			return
		}

		var result []Row

		for i := range data {
			var r Row
			r.Id = data[i][9]
			r.CategoryName = data[i][10]
			r.SubCategoryName = data[i][11]
			r.ProductName = data[i][6]
			r.ProductDesc = data[i][12]
			r.Productstatus = data[i][13]
			r.CreatedBy = data[i][14]
			//r.ProductQty = data[i][7]
			r.ProductpackageQty = data[i][7] + data[i][8]
			r.ProductSKU = data[i][15]
			r.Timestamp = data[i][16]
			result = append(result, r)
		}

		c.Data["CustomerData"] = result
	} else if utype == "CLUSTEREXECUTIVE" {

		//NO CODE FOR CUSTEREXECUTIVE
	} else {

		row, err := db.Db.Query(`select id,category_name,subcategory_name,product_name,product_desc,product_status,created_by,product_qty,product_package_qty,product_sku,create_date from Lookup_products ORDER BY id DESC`)

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get the error message info")
			return
		}
		defer sql.Close(row)
		_, data, err := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Message Scan Fail")
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
		if len(data) <= 0 {
			err = errors.New("Unable to get the error message info")
			return
		}

		var result []Row

		for i := range data {
			var r Row
			r.Id = data[i][0]
			r.CategoryName = data[i][1]
			r.SubCategoryName = data[i][2]
			r.ProductName = data[i][3]
			r.ProductDesc = data[i][4]
			r.Productstatus = data[i][5]
			r.CreatedBy = data[i][6]
			//r.ProductQty = data[i][7]
			r.ProductpackageQty = data[i][8] + data[i][7]
			r.ProductSKU = data[i][9]
			r.Timestamp = data[i][10]
			result = append(result, r)
		}

		c.Data["CustomerData"] = result

	}

	return
}

func (c *SearchProduct) Post() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "node post page")
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
			c.TplName = "products/searchProduct/searchProduct.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "products/searchProduct/searchProduct.html"
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
	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "utype", utype)

	id := sess.Get("uid").(string)
	c.Data["id"] = id

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

	auth, err := utils.IsAuthorized(utype, "Product")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	input_name := c.Input().Get("input_name")
	input_email := c.Input().Get("input_email")
	input_category_name := c.Input().Get("input_category_name")
	input_product_name := c.Input().Get("input_product_name")
	input_status := c.Input().Get("input_status")

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_name - ", input_name)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_email - ", input_email)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_name - ", input_category_name)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_name - ", input_product_name)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_status - ", input_status)

	c.Data["DateRange"] = c.Input().Get("daterange")
	tmp := strings.Split(c.Input().Get("daterange"), " - ")

	from := ""
	to := ""
	if len(tmp) == 2 {
		from = tmp[0]
		to = tmp[1]
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "From", from)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "To", to)

	// row, err := db.Db.Query(`SELECT id,first_name,email,status FROM sysuser.credential`)

	// row, err := db.Db.Query(`select first_name,mob_no,email,status,user_type from Sysuser_info where
	// DATE(create_date) >= STR_TO_DATE(?, '%m/%d/%Y') AND DATE(create_date) <= STR_TO_DATE(?, '%m/%d/%Y') ORDER BY create_date DESC`, from, to)

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

		row, err := db.Db.Query(`SELECT 
		Product2Cluster.cluster_product_id,
		Product2Cluster.cluster_id, 
		Product2Cluster.product_id, 
        Product2Cluster.status, 
        Product2Cluster.create_date,
        Clusters.cluster_name, 
       	Lookup_products.product_name,
     	Lookup_products.product_package_qty,
     	Lookup_products.product_qty,
		Lookup_products.id,
		Lookup_products.
		category_name,
		Lookup_products.subcategory_name,
    	Lookup_products.product_desc,
		Lookup_products.product_status,
		Lookup_products.created_by,
		Lookup_products.product_sku,
		Lookup_products.create_date  
		FROM C2F_NEW.Product2Cluster
		LEFT JOIN Clusters 
		ON C2F_NEW.Product2Cluster.cluster_id = Clusters.id
		LEFT JOIN Lookup_products 
		ON C2F_NEW.Product2Cluster.product_id = Lookup_products.id where Product2Cluster.cluster_id=? AND (Lookup_products.category_name='' OR Lookup_products.category_name like ?)AND (Lookup_products.product_status='' OR Lookup_products.product_status like ?)AND (Lookup_products.product_name='' OR Lookup_products.product_name like ?)AND (DATE(Lookup_products.create_date='') OR DATE (Lookup_products.create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(Lookup_products.create_date='') OR DATE (Lookup_products.create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) ORDER BY Lookup_products.create_date DESC`, data1[0][1], input_category_name+"%", input_status+"%", input_product_name+"%", from, to)

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get the error message info")
			return
		}
		defer sql.Close(row)
		_, data, err := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Message Scan Fail")
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
		if len(data) <= 0 {
			err = errors.New("Unable to get the error message info")
			return
		}

		var result []Row

		for i := range data {
			var r Row
			r.Id = data[i][9]
			r.CategoryName = data[i][10]
			r.SubCategoryName = data[i][11]
			r.ProductName = data[i][6]
			r.ProductDesc = data[i][12]
			r.Productstatus = data[i][13]
			r.CreatedBy = data[i][14]
			//r.ProductQty = data[i][7]
			r.ProductpackageQty = data[i][7] + data[i][8]
			r.ProductSKU = data[i][15]
			r.Timestamp = data[i][16]
			result = append(result, r)
		}

		c.Data["CustomerData"] = result
	} else if utype == "CLUSTEREXECUTIVE" {

		//NO CODE FOR CUSTEREXECUTIVE
	} else {

		row, err := db.Db.Query(`select id,category_name,subcategory_name,product_name,product_desc,product_status,created_by,product_qty,product_package_qty,product_sku,create_date from Lookup_products where (category_name='' OR category_name like ?)AND (product_status='' OR product_status like ?)AND (product_name='' OR product_name like ?)AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) ORDER BY id DESC`, input_category_name+"%", input_status+"%", input_product_name+"%", from, to)

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get the error message info")
			return
		}
		defer sql.Close(row)
		_, data, err := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Message Scan Fail")
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
		if len(data) <= 0 {
			err = errors.New("Unable to get the error message info")
			return
		}

		var result []Row

		for i := range data {
			var r Row
			r.Id = data[i][0]
			r.CategoryName = data[i][1]
			r.SubCategoryName = data[i][2]
			r.ProductName = data[i][3]
			r.ProductDesc = data[i][4]
			r.Productstatus = data[i][5]
			r.CreatedBy = data[i][6]
			//r.ProductQty = data[i][7]
			r.ProductpackageQty = data[i][8] + data[i][7]
			r.ProductSKU = data[i][9]
			r.Timestamp = data[i][10]
			result = append(result, r)
		}

		c.Data["CustomerData"] = result
	}

	// row, err := db.Db.Query(`select id,category_name,subcategory_name,product_name,product_desc,product_status,created_by,product_qty,product_package_qty,product_sku from Lookup_products where (category_name='' OR category_name like ?)AND (product_status='' OR product_status like ?)AND (product_name='' OR product_name like ?)AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) ORDER BY id DESC`, input_category_name+"%", input_status+"%", input_product_name+"%", from, to)

	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Unable to get the error message info")
	// 	return
	// }
	// defer sql.Close(row)
	// _, data, err := sql.Scan(row)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Message Scan Fail")
	// 	return
	// }
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	// if len(data) <= 0 {
	// 	err = errors.New("Unable to get the error message info")
	// 	return
	// }

	// var result []Row

	// for i := range data {
	// 	var r Row
	// 	r.Id = data[i][0]
	// 	r.CategoryName = data[i][1]
	// 	r.SubCategoryName = data[i][2]
	// 	r.ProductName = data[i][3]
	// 	r.ProductDesc = data[i][4]
	// 	r.Productstatus = data[i][5]
	// 	r.CreatedBy = data[i][6]
	// 	//r.ProductQty = data[i][7]
	// 	r.ProductpackageQty = data[i][8] + data[i][7]
	// 	r.ProductSKU = data[i][9]
	// 	result = append(result, r)
	// }

	// c.Data["CustomerData"] = result

	return
}
