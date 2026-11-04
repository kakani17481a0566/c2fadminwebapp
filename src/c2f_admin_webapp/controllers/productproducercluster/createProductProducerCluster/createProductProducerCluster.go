package createProductProducerCluster

import (
	"c2f_admin_webapp/session"

	"runtime/debug"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"
	"errors"
	"strings"

	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/services"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

type CreateProductProducerCluster struct {
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
type Field struct {
	Id    string
	Name  string
	Email string
}

type Field1 struct {
	Id                string
	Name              string
	ProductQty        string
	ProductpackageQty string
}

type Field2 struct {
	Id     string
	Mobile string
	Name   string
}

func (c *CreateProductProducerCluster) Get() {
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
			c.TplName = "productproducerclusters/createProductProducerCluster/createProductProducerCluster.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "productproducerclusters/createProductProducerCluster/createProductProducerCluster.html"
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

	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1

	id := sess.Get("uid").(string)
	c.Data["id"] = id

	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "utype", utype)

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

	data2, err := services.GetProducers()
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
		d.Id = data2[i][0]
		d.Mobile = data2[i][1]
		d.Name = data2[i][2]
		Dis2.Fields2 = append(Dis2.Fields2, d)
	}
	c.Data["Dis2"] = Dis2
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis1)

	auth, err := utils.IsAuthorized(utype, "ProductProducerCluster")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	return
}
func (c *CreateProductProducerCluster) Post() {

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
			c.TplName = "productproducerclusters/createProductProducerCluster/createProductProducerCluster.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Product-Producer-Cluster Added Successfully"
			c.TplName = "productproducerclusters/createProductProducerCluster/createProductProducerCluster.html"
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

	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1

	id := sess.Get("uid").(string)
	c.Data["id"] = id
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "id", id)

	auth, err := utils.IsAuthorized(utype, "ProductProducerCluster")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	data3, err3 := services.GetClusteruser()
	if err3 != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err3)
		err3 = errors.New("Cluster user fetch Failed")
		return
	}
	if len(data3) <= 0 {
		err3 = errors.New("Cluster user  Not Found")
		return
	}
	var Dis Display
	for i := range data3 {
		var d Field
		d.Id = data3[i][0]
		d.Name = data3[i][1]
		Dis.Fields = append(Dis.Fields, d)
	}
	c.Data["Dis"] = Dis
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

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

	data2, err := services.GetProducers()
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
		d.Id = data2[i][0]
		d.Mobile = data2[i][1]
		d.Name = data2[i][2]
		Dis2.Fields2 = append(Dis2.Fields2, d)
	}
	c.Data["Dis2"] = Dis2
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis1)

	select_product_name := c.Input().Get("select_product_name")
	select_producer_name := c.Input().Get("select_producer_name")
	select_cluster_name := c.Input().Get("select_cluster_name")

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "select_product_name", select_product_name)

	row, err := db.Db.Query(`SELECT id,krishivala_CGST_charge,krishivala_SGST_charge,product_name from Lookup_products WHERE id = ?`, select_product_name)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch GST")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch GST")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	producer_CGST_charge := data[0][1]
	producer_SGST_charge := data[0][2]

	log.Println(beego.AppConfig.String("loglevel"), "Debug", producer_CGST_charge)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", producer_SGST_charge)

	//producer_CGST_charge := c.Input().Get("producer_CGST_charge")
	//producer_SGST_charge := c.Input().Get("producer_SGST_charge")
	producer_CESS_charge := c.Input().Get("producer_CESS_charge")
	producer_OtherTAX_charge := c.Input().Get("producer_OtherTAX_charge")

	input_status := c.Input().Get("input_status")
	input_sell := c.Input().Get("input_sell")

	err = CheckAlreadyExists(select_product_name, select_producer_name, select_cluster_name)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Mapping already Exists")
		return
	}

	result, err := db.Db.Exec(`INSERT INTO Product2Producer2Cluster (product_id,
		producer_id,
		cluster_id,
		status,
		producer_CGST_charge,
		producer_SGST_charge,
		producer_CESS_charge,
		producer_OtherTAX_charge,
		create_by,
		independant_sell,
		create_date)
		VALUES (?,?,?,?,?,?,?,?,?,?,now())`,
		select_product_name,
		select_producer_name,
		select_cluster_name,
		input_status,
		producer_CGST_charge,
		producer_SGST_charge,
		producer_CESS_charge,
		producer_OtherTAX_charge,
		id,
		input_sell)
	if err != nil {
		err = errors.New("Product-Producer-Cluster creation failed")
		return
	}

	i, err := result.RowsAffected()
	if err != nil || i == 0 {
		err = errors.New("Product-Producer-Cluster creation failed")
		return
	}

	return
}

func CheckAlreadyExists(product, producer, cluster string) (err error) {

	err = nil

	row, err := db.Db.Query(`select count(*) from Product2Producer2Cluster WHERE product_id = ? AND producer_id=? AND cluster_id=?`, product, producer, cluster)

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
