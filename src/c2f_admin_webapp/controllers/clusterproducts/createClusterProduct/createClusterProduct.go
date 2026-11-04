package createClusterProduct

import (
	"c2f_admin_webapp/session"

	"runtime/debug"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/model/utils"
	"c2f_admin_webapp/services"
	"errors"
	"strings"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

type CreateClusterProduct struct {
	// beego.Controller
	middleware.BaseController
}

type Display struct {
	Fields []Field
}

type Display1 struct {
	Fields1 []Field1
}
type Field struct {
	Id   string
	Name string

	Email string
}

type Field1 struct {
	Id                string
	Name              string
	ProductQty        string
	ProductpackageQty string
}

func (c *CreateClusterProduct) Get() {
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
			c.TplName = "clusterproducts/createClusterProduct/createClusterProduct.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "clusterproducts/createClusterProduct/createClusterProduct.html"
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

	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "utype", utype)

	if utype == "CLUSTERMANAGER" {

		data, err := services.GetClusterAsPerMapping(id)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster user fetch Failed")
			return
		}
		if len(data) <= 0 {
			err = errors.New("Cluster user  Not Found")
			return
		}
		var Dis Display
		for i := range data {
			var d Field
			d.Id = data[i][2]
			d.Name = data[i][3]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)
	} else if utype == "CLUSTEREXECUTIVE" {
		data, err := services.GetClusterAsPerMappingExec(id)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster user fetch Failed")
			return
		}
		if len(data) <= 0 {
			err = errors.New("Cluster user  Not Found")
			return
		}
		var Dis Display
		for i := range data {
			var d Field
			d.Id = data[i][2]
			d.Name = data[i][3]

			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

	} else {
		data, err := services.GetClusteruser()
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster user fetch Failed")
			return
		}
		if len(data) <= 0 {
			err = errors.New("Cluster user  Not Found")
			return
		}
		var Dis Display
		for i := range data {
			var d Field
			d.Id = data[i][0]
			d.Name = data[i][1]

			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)
	}

	data1, err := services.Getproducts()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Product fetch Failed")
		return
	}
	if len(data1) <= 0 {
		err = errors.New("Product  Not Found")
		return
	}
	var Dis1 Display1
	for i := range data1 {
		var d Field1
		d.Id = data1[i][0]
		d.Name = data1[i][1]
		d.ProductpackageQty = data1[i][3] + data1[i][2]

		Dis1.Fields1 = append(Dis1.Fields1, d)
	}
	c.Data["Dis1"] = Dis1
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis1)

	auth, err := utils.IsAuthorized(utype, "ClusterProduct")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	return
}
func (c *CreateClusterProduct) Post() {

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
			c.TplName = "clusterproducts/createClusterProduct/createClusterProduct.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Cluster Product Mapped Successfully"
			c.TplName = "clusterproducts/createClusterProduct/createClusterProduct.html"
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

	auth, err := utils.IsAuthorized(utype, "ClusterProduct")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	select_cluster_name := c.Input().Get("select_cluster_name")
	select_product_name := c.Input().Get("select_product_name")
	input_product_status := c.Input().Get("input_product_status")
	input_cluster_selling_price := c.Input().Get("input_cluster_selling_price")
	input_cluster_buying_price := c.Input().Get("input_cluster_buying_price")
	input_cluster_market_price := c.Input().Get("input_cluster_market_price")

	cluster_supply_charge := c.Input().Get("cluster_supply_charge")
	cluster_packaging_charge := c.Input().Get("cluster_packaging_charge")
	cluster_bank_charge := c.Input().Get("cluster_bank_charge")
	cluster_discount := c.Input().Get("cluster_discount")

	delivery_discount := c.Input().Get("delivery_discount")

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cluster_id  - ", select_cluster_name)

	if select_cluster_name == "" || select_product_name == "" || input_product_status == "" || input_cluster_selling_price == "" || input_cluster_buying_price == "" || input_cluster_market_price == "" ||
		cluster_supply_charge == "" || cluster_packaging_charge == "" || cluster_bank_charge == "" || cluster_discount == "" {
		err = errors.New("* fields cannot be empty")
		return
	}
	if input_product_status == "" {
		err = errors.New("status cannot be empty")
		return
	}
	err = CheckUserAlreadyExists(select_product_name, select_cluster_name)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Product mapping already Exists")
		return
	}

	result, err := db.Db.Exec(`INSERT INTO Product2Cluster (
	cluster_id, 
	product_id,
	 status,
	 cluster_selling_price,
	 cluster_buying_price,
	 cluster_market_price,
	cluster_supply_charge,
	cluster_packaging_charge,
	cluster_bank_charge,
	cluster_discount,
	 createdby, 
	delivery_discount,
	create_date
)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,now())`,
		select_cluster_name,
		select_product_name,
		input_product_status,
		input_cluster_selling_price,
		input_cluster_buying_price,
		input_cluster_market_price,
		cluster_supply_charge,
		cluster_packaging_charge,
		cluster_bank_charge,
		cluster_discount,
		id,
		delivery_discount)
	if err != nil {
		err = errors.New("ClusterProduct creation failed")
		return
	}

	i, err := result.RowsAffected()
	if err != nil || i == 0 {
		err = errors.New("ClusterProduct creation failed")
		return
	}

	// result1, err := db.Db.Exec(`UPDATE Lookup_products SET
	// supply_charge = ?, service_charge = ?, delivery_charge = ?, tax_type = ?,tax=?,discount=?
	// WHERE id = ?`, input_supply_charge, input_service_charge, input_delivery_charge, input_tax_type, input_tax_charge, input_discount, select_product_name)
	// if err != nil {
	// 	err = errors.New("Cluster Product updation failed")
	// 	return
	// }

	// i1, err := result1.RowsAffected()
	// if err != nil || i1 == 0 {
	// 	err = errors.New("ClusterProduct updation failed")
	// 	return
	// }

	return
}

func CheckUserAlreadyExists(product, cluster string) (err error) {

	err = nil

	row, err := db.Db.Query(`select count(*) from Product2Cluster WHERE product_id = ? AND cluster_id=? `, product, cluster)

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

	if countlen != "0" {
		err = errors.New("Mappig already exists")
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}

	return

}
