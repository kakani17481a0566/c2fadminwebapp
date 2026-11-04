package createReverseInventory

import (
	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/model/utils"
	"c2f_admin_webapp/services"
	"c2f_admin_webapp/session"

	"errors"
	"strings"

	"runtime/debug"
	"strconv"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/util/log"
)

type CreateReverseInventory struct {
	// beego.Controller
	middleware.BaseController
}

type Display struct {
	Fields []Field
}
type Field struct {
	Id   string
	Name string
	//ProductSku  string
}
type Display1 struct {
	Fields1 []Field1
}
type Field1 struct {
	Id  string
	Sku string
}

func (c *CreateReverseInventory) Get() {
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
			c.TplName = "reverseinventory/createReverseInventory/createReverseInventory.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "reverseinventory/createReverseInventory/createReverseInventory.html"
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
	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "utype", utype)

	language := sess.Get("language").(string)
	c.Data["language"] = language
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "language", language)

	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1

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
			err = errors.New("Role Failed")
			return
		}
		if len(data) <= 0 {
			err = errors.New("Roles Not Found")
			return
		}
		var Dis Display
		for i := range data {
			var d Field
			d.Id = data[i][2]
			d.Name = data[i][3]
			//      d.ProductSku = data[i][2]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)
	} else {
		data, err := services.GetClusteruser()
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Role Failed")
			return
		}
		if len(data) <= 0 {
			err = errors.New("Roles Not Found")
			return
		}
		var Dis Display
		for i := range data {
			var d Field
			d.Id = data[i][0]
			d.Name = data[i][1]
			//      d.ProductSku = data[i][2]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)
	}
	data1, err := services.Getinventory1()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("User fetch Failed")
		return
	}
	if len(data1) <= 0 {
		err = errors.New("User  Not Found")
		return
	}
	var Dis1 Display1
	for i := range data1 {
		var d Field1
		d.Id = data1[i][0]
		d.Sku = data1[i][1]
		Dis1.Fields1 = append(Dis1.Fields1, d)
	}
	c.Data["Dis1"] = Dis1
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis1)
	auth, err := utils.IsAuthorized(utype, "ReverseInventory")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	return
}
func (c *CreateReverseInventory) Post() {
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
			c.TplName = "reverseinventory/createReverseInventory/createReverseInventory.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Reverse Inventory Created Successfully"
			c.TplName = "reverseinventory/createReverseInventory/createReverseInventory.html"
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

	uid := sess.Get("uid").(string)
	c.Data["id"] = uid
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "id", uid)

	auth, err := utils.IsAuthorized(utype, "ReverseInventory")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	select_cluster_name := c.Input().Get("select_cluster_name")
	select_product_sku := c.Input().Get("select_product_sku")

	select_reverse_remarks := c.Input().Get("select_reverse_remarks")
	input_numberof_units := c.Input().Get("input_numberof_units")

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "select_cluster_name - ", select_cluster_name)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "select_product_sku - ", select_product_sku)

	dunits, aunits, runits, sku, cluster, damageunits, err := services.Getinventory(select_product_sku, select_cluster_name)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Inventory Get Failed")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "dunits - ", dunits)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "aunits - ", aunits)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "runits - ", runits)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "sku - ", sku)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cluster - ", cluster)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "damageunits - ", damageunits)

	ddunits, err := strconv.ParseInt(dunits, 0, 64)
	aaunits, err := strconv.ParseInt(aunits, 0, 64)
	inunits, err := strconv.ParseInt(input_numberof_units, 0, 64)

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "inunits - ", inunits)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "ddunits - ", ddunits)

	if inunits > ddunits {

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "inside - ")
		err = errors.New("Invalid Reverse Inventory Unit Value")
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return

	}

	UTD := ddunits - inunits
	TotalDel := strconv.FormatInt(UTD, 10)

	UTA := aaunits + inunits
	TotalAvl := strconv.FormatInt(UTA, 10)

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "TotalDel - ", TotalDel)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "TotalAvl - ", TotalAvl)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "select_reverse_remarks - ", select_reverse_remarks)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_numberof_units - ", input_numberof_units)

	res2, err := db.Db.Exec(`UPDATE Inventory SET inventory_delivered = ? ,inventory_available=?,reverse_remarks=?,reverse_units=?,updated_by=?,update_date=now()
	                            WHERE product_sku = ? AND cluster_name=?`, TotalDel, TotalAvl, select_reverse_remarks, input_numberof_units, uid, select_product_sku, select_cluster_name)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Inventory Get Failed")
		return
	}
	i2, err := res2.RowsAffected()
	if err != nil || i2 == 0 {
		err = errors.New(" Inventory updation failed")
		return
	}

	return
}
