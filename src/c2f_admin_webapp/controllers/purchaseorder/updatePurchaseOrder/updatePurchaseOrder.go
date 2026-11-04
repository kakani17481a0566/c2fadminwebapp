package updatePurchaseOrder

import (
	"c2f_admin_webapp/session"

	"runtime/debug"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"
	"c2f_admin_webapp/services"
	"errors"
	"strconv"
	"strings"

	"c2f_admin_webapp/model/db"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

type Display struct {
	Fields []Field
}
type Field struct {
	Id    string
	Name  string
	Email string
}
type UpdatePurchaseOrder struct {
	// beego.Controller
	middleware.BaseController
}

func (c *UpdatePurchaseOrder) Get() {
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
			c.TplName = "purchaseorder/updatePurchaseOrder/updatePurchaseOrder.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "purchaseorder/updatePurchaseOrder/updatePurchaseOrder.html"
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

	auth, err := utils.IsAuthorized(utype, "PurchaseOrder")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	data, err := services.GetClusteruser()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Delivery user fetch Failed")
		return
	}
	if len(data) <= 0 {
		err = errors.New("Delivery user  Not Found")
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

	row, err := db.Db.Query(`select id,cluster_name, 
         product_name, 
        product_sku, 
        inventory_available 
    from Inventory WHERE id =?  limit 1`, AdminId)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	defer sql.Close(row)
	_, data3, err := sql.Scan(row)
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
	log.Println(beego.AppConfig.String("loglevel"), "Debug", data3)
	for i := range data3 {
		c.Data["Clustername"] = data3[i][1]
		c.Data["Productname"] = data3[i][2]
		c.Data["Productsku"] = data3[i][3]
		c.Data["Inventoryavailable"] = data3[i][4]
	}

	return

}

func (c *UpdatePurchaseOrder) Post() {

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
			c.TplName = "purchaseorder/updatePurchaseOrder/updatePurchaseOrder.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Product Updated Successfully"
			c.TplName = "purchaseorder/updatePurchaseOrder/updatePurchaseOrder.html"
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

	language := sess.Get("language").(string)
	c.Data["language"] = language
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "language", language)

	clustername := sess.Get("clustername").(string)
	clustername1 := strings.ToUpper(clustername)
	c.Data["clustername"] = clustername1
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "clustername", clustername1)

	auth, err := utils.IsAuthorized(utype, "PurchaseOrder")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	row, err := db.Db.Query(`select id,cluster_name, 
         product_name, 
        product_sku, 
        inventory_available ,
        inventory_delivered
    from Inventory WHERE id =?  limit 1`, AdminId)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	defer sql.Close(row)
	_, data3, err := sql.Scan(row)
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
	log.Println(beego.AppConfig.String("loglevel"), "Debug", data3)
	for i := range data3 {
		c.Data["Clustername"] = data3[i][1]
		c.Data["Productname"] = data3[i][2]
		c.Data["Productsku"] = data3[i][3]
		c.Data["Inventoryavailable"] = data3[i][4]
	}

	input_source_clustername := c.Input().Get("input_source_clustername")
	input_product_name := c.Input().Get("input_product_name")
	input_product_sku := c.Input().Get("input_product_sku")
	input_beneficiary_clustername := c.Input().Get("input_beneficiary_clustername")
	input_required_units := c.Input().Get("input_required_units")
	input_available_units := c.Input().Get("input_available_units")
	input_delivered_units := data3[0][5]

	result, err := db.Db.Exec(`INSERT INTO ClusterPurchaseProduct (source_cluster_name, 
        beneficiary_cluster_name,
         product_name, 
        product_sku,  
        number_of_units,
        create_date)
                VALUES (?,?,?,?,?,now())`,
		input_source_clustername,
		input_beneficiary_clustername,
		input_product_name,
		input_product_sku,
		input_required_units)
	if err != nil {
		err = errors.New("Cluster Purchase Product creation failed")
		return
	}
	i, err := result.RowsAffected()
	if err != nil || i == 0 {
		err = errors.New("Cluster Purchase Product creation failed")
		return
	}
	aunits, err := strconv.ParseFloat(input_available_units, 64)
	runits, err := strconv.ParseFloat(input_required_units, 64)
	dunits, err := strconv.ParseFloat(input_delivered_units, 64)
	TA := aunits - runits
	TA1 := strconv.FormatFloat(TA, 'f', 0, 32)
	TD := runits + dunits
	TD1 := strconv.FormatFloat(TD, 'f', 0, 32)
	res2, err := db.Db.Exec(`UPDATE Inventory SET inventory_delivered = ? ,inventory_available=?
                              WHERE id = ?`, TD1, TA1, AdminId)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Inventory Get Failed")
		//              msg11 = "Inventory Get Failed"
		return
	}
	i2, err := res2.RowsAffected()
	if err != nil || i2 == 0 {
		err = errors.New(" CustomerOrderTxn Delivery updation failed")
		return
	}

	return
}
