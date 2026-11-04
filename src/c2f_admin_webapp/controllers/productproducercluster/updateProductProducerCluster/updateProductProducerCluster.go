package updateProductProducerCluster

import (
	"c2f_admin_webapp/session"

	"runtime/debug"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"
	"errors"
	"strings"

	"c2f_admin_webapp/model/db"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

type Row struct {
	Id           string
	ProductName  string
	ProducerName string
	ClusterName  string
	Status       string
	Sell         string
}

type UpdateProductProducerCluster struct {
	// beego.Controller
	middleware.BaseController
}

func (c *UpdateProductProducerCluster) Get() {
	AdminId := c.Ctx.Input.Param(":AdminID")
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
			c.TplName = "productproducerclusters/updateProductProducerCluster/updateProductProducerCluster.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update ProductProducerCluster Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "productproducerclusters/updateProductProducerCluster/updateProductProducerCluster.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update ProductProducerCluster Page Success")
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

	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1

	auth, err := utils.IsAuthorized(utype, "ProductProducerCluster")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	row, err := db.Db.Query(`
SELECT Product2Producer2Cluster.id,Product2Producer2Cluster.product_id,Product2Producer2Cluster.producer_id, Product2Producer2Cluster.cluster_id, 
       Product2Producer2Cluster.status, Product2Producer2Cluster.create_date,
       Lookup_products.product_name,
       Customer.mob_no,
       Clusters.cluster_name,
       Product2Producer2Cluster.independant_sell,
       Product2Producer2Cluster.producer_CESS_charge,
       Product2Producer2Cluster.producer_OtherTAX_charge
FROM Product2Producer2Cluster
LEFT JOIN Lookup_products 
ON Product2Producer2Cluster.product_id = Lookup_products.id
LEFT JOIN Customer 
ON Product2Producer2Cluster.producer_id = Customer.id
LEFT JOIN Clusters 
ON Product2Producer2Cluster.cluster_id = Clusters.id WHERE Product2Producer2Cluster.id =?  limit 1`, AdminId)
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
		c.Data["ProductName"] = data[i][6]
		c.Data["ProducerName"] = data[i][7]
		c.Data["ClusterName"] = data[i][8]
		c.Data["Status"] = data[i][4]
		c.Data["Sell"] = data[i][9]
		c.Data["ProducerCESSCharge"] = data[i][10]
		c.Data["ProducerOtherTAXCharge"] = data[i][11]

	}

	return

}

func (c *UpdateProductProducerCluster) Post() {

	AdminId := c.Ctx.Input.Param(":AdminID")

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
			c.TplName = "productproducerclusters/updateProductProducerCluster/updateProductProducerCluster.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update ProductProducerCluster Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Product-Producer-Cluster Map Updated Successfully"
			c.TplName = "productproducerclusters/updateProductProducerCluster/updateProductProducerCluster.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update ProductProducerCluster  Page Success")
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

	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1

	auth, err := utils.IsAuthorized(utype, "ProductProducerCluster")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	input_status := c.Input().Get("input_status")
	input_sell := c.Input().Get("input_sell")

	res, err := db.Db.Exec(`UPDATE Product2Producer2Cluster SET 
	status = ?,independant_sell=?
	WHERE id = ?`, input_status, input_sell, AdminId)
	if err != nil {
		err = errors.New("Product-Producer-Cluster updation failed")
		return
	}

	i, err := res.RowsAffected()
	if err != nil || i == 0 {
		err = errors.New("Product-Producer-Cluster updation failed")
		return
	}

	row, err := db.Db.Query(`
SELECT Product2Producer2Cluster.id,Product2Producer2Cluster.product_id,Product2Producer2Cluster.producer_id, Product2Producer2Cluster.cluster_id, 
       Product2Producer2Cluster.status, Product2Producer2Cluster.create_date,
       Lookup_products.product_name,
       Customer.mob_no,
       Clusters.cluster_name
FROM Product2Producer2Cluster
LEFT JOIN Lookup_products 
ON Product2Producer2Cluster.product_id = Lookup_products.id
LEFT JOIN Customer 
ON Product2Producer2Cluster.producer_id = Customer.id
LEFT JOIN Clusters 
ON Product2Producer2Cluster.cluster_id = Clusters.id WHERE Product2Producer2Cluster.id =?  limit 1`, AdminId)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get product data")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get product data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("product data not found")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	for i := range data {

		c.Data["Id"] = data[i][0]
		c.Data["ProductName"] = data[i][6]
		c.Data["ProducerName"] = data[i][7]
		c.Data["ClusterName"] = data[i][8]
		c.Data["Status"] = data[i][4]

	}

	return
}
