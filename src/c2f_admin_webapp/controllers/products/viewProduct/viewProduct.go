package viewProduct

import (
	"c2f_admin_webapp/session"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"
	"runtime/debug"

	"errors"
	"strings"

	"c2f_admin_webapp/model/db"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"
	"c2f_admin_webapp/model/util/log"
)

type Row struct {
	Id              string
	CategoryName    string
	SubCategoryName string
	ProductName     string
	ProductDesc     string
	Productstatus   string
	ProductQuantity string
	QualityType     string
	QualityMeasure  string
	ShelLife        string
	PackageType     string
	PackageQty      string
	ProductOrigin   string

	Kselprice string
	kbuyprice string
	kmarprice string

	comname  string
	cgst     string
	sgst     string
	hsn      string
	cess     string
	othertax string

	CreatedBy string
}
type ViewProduct struct {
	// beego.Controller
	middleware.BaseController
}

func (c *ViewProduct) Get() {
	AdminId := c.Ctx.Input.Param(":AdminID")
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
			c.TplName = "products/viewProduct/viewProduct.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "products/viewProduct/viewProduct.html"
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

	auth, err := utils.IsAuthorized(utype, "Product")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	row, err := db.Db.Query(`select id,category_name,subcategory_name,product_name,product_desc,product_status, 
product_qty,
product_quality_type,
product_quality_mesure,
product_origin, 
product_package_type,
product_package_qty,
product_shelf_life,
krishivala_selling_price,
krishivala_buying_price,
krishivala_market_price,
commodity_name,
krishivala_CGST_charge,
krishivala_SGST_charge,
hsn_code,
krishivala_CESS_charge,
krishivala_OtherTAX_charge,
product_shelf_life_unit
from Lookup_products WHERE id =?  limit 1`, AdminId)
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
		c.Data["CategoryName"] = data[i][1]
		c.Data["SubCategoryName"] = data[i][2]
		c.Data["ProductName"] = data[i][3]
		c.Data["ProductDesc"] = data[i][4]
		c.Data["Status"] = data[i][5]
		c.Data["ProductQuantity"] = data[i][6]
		c.Data["QualityType"] = data[i][7]
		c.Data["QualityMeasure"] = data[i][8]
		c.Data["ShelLife"] = data[i][12]
		c.Data["PackageType"] = data[i][10]
		c.Data["PackageQty"] = data[i][11]
		c.Data["ProductOrigin"] = data[i][9]

		c.Data["Kselprice"] = data[i][13]
		c.Data["kbuyprice"] = data[i][14]
		c.Data["kmarprice"] = data[i][15]

		c.Data["comname"] = data[i][16]
		c.Data["cgst"] = data[i][17]
		c.Data["sgst"] = data[i][18]
		c.Data["hsn"] = data[i][19]
		c.Data["cess"] = data[i][20]
		c.Data["othertax"] = data[i][21]
		c.Data["ShelLifeUnit"] = data[i][22]
	}

	return

}
