package updateProductPrice

import (
	"c2f_admin_webapp/session"

	"c2f_admin_webapp/model/utils"
	"errors"
	"runtime/debug"
	"strings"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/db"

	"github.com/astaxie/beego"

	"encoding/csv"
	"fmt"

	"log"

	"c2f_admin_webapp/model/database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type UpdateProductPrice struct {
	// beego.Controller
	middleware.BaseController
}

type ProductData struct {
	Psku         string
	KBuyPrice    string
	KSellPrice   string
	KMarPrice    string
	CatName      string
	SbcatName    string
	ProductName  string
	ProdQnty     string
	ProdQltytype string
	ProdQM       string
	ProdShelife  string
	ProdPcktype  string
	ProdPckqty   string
	ProdOrigin   string
	ProdDesc     string
	ProdStatus   string
	ProdCGST     string
	ProdSGST     string
	ProdCESS     string
	ProdTax      string
	Commodity    string
	Hsn          string
}

func (c *UpdateProductPrice) Get() {
	//	AdminId := c.Ctx.Input.Param(":AdminID")
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
			c.TplName = "products/updateProductPrice/updateProductPrice.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "products/updateProductPrice/updateProductPrice.html"
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

	return
}

func (c *UpdateProductPrice) Post() {
	//	AdminId := c.Ctx.Input.Param(":AdminID")
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
			c.TplName = "products/updateProductPrice/updateProductPrice.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Product Price Updated Successfully"
			c.TplName = "products/updateProductPrice/updateProductPrice.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer  Page Success")
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

	auth, err := utils.IsAuthorized(utype, "Product")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	c.TplName = "updateProductPrice.html"

	// Get the  uploaded Csv File

	file, information, err := c.GetFile("uploadname")

	if err != nil {

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "File - ", "File retrieval failure")
		c.Data["DisplayMessage"] = "File retrieval failure"

		return
	}

	defer file.Close()

	filename := information.Filename
	fmt.Println(filename)

	csvLines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	for _, line := range csvLines {
		pdct := ProductData{

			KBuyPrice:  line[0],
			KSellPrice: line[1],
			KMarPrice:  line[2],
			Psku:       line[3],
		}

		row1, err1 := db.Db.Query(`select count(*) from Lookup_products WHERE product_sku = ?`, pdct.Psku)

		if err1 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
			err = errors.New("category doesnot exist" + pdct.CatName)
			c.Data["DisplayMessage"] = "category doesnot exist"

			return
		}
		defer sql.Close(row1)
		_, data, err1 := sql.Scan(row1)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
			err = errors.New("category doesnot exist" + pdct.CatName)
			c.Data["DisplayMessage"] = "category doesnot exist"

			return
		}
		countlen1 := data[0][0]
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "countlen", countlen1)

		if countlen1 != "0" {

			if pdct.Psku == "" || pdct.KBuyPrice == "" || pdct.KSellPrice == "" || pdct.KMarPrice == "" {

				err = errors.New("No fields should be empty")
				return
			}

			fmt.Println(pdct.Psku + "sku  exist in table")

			res, err := db.Db.Exec(`UPDATE Lookup_products SET krishivala_buying_price=?, krishivala_selling_price=?, krishivala_market_price=? WHERE product_sku=?`, pdct.KBuyPrice, pdct.KSellPrice, pdct.KMarPrice, pdct.Psku)
			if err != nil {
				err = errors.New("Product updation failed for" + pdct.Psku)
				c.Data["DisplayMessage"] = "Product updation failed"

				return
			}

			i, err := res.RowsAffected()
			if err != nil || i == 0 {
				err = errors.New("Product updation failed")
				c.Data["DisplayMessage"] = "Product updation failed"

				return
			}

		}
	}
}
