package createCategory

import (
	"c2f_admin_webapp/session"

	"runtime/debug"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"
	"errors"

	"c2f_admin_webapp/model/db"

	"os"
	"path"
	"strings"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
	//	"c2f_admin_webapp/model/util/sku"
)

type CreateCategory struct {
	// beego.Controller
	middleware.BaseController
}

func (c *CreateCategory) Get() {
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
			c.TplName = "category/createCategory/createCategory.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "category/createCategory/createCategory.html"
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

	auth, err := utils.IsAuthorized(utype, "Category")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	return
}
func (c *CreateCategory) Post() {

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
			c.TplName = "category/createCategory/createCategory.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Category Added Successfully"
			c.TplName = "category/createCategory/createCategory.html"
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

	auth, err := utils.IsAuthorized(utype, "Category")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	input_category_name := c.Input().Get("input_category_name")
	input_category_status := c.Input().Get("input_category_status")

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "name  - ", input_category_name)

	if input_category_name == "" || input_category_status == "" {
		err = errors.New("category name and category status cannot be empty")
		return
	}
	// if input_product_status == "" {
	// 	err = errors.New("status cannot be empty")
	// 	return
	// }

	err = CheckUserAlreadyExists(input_category_name)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Category already Exists")
		return
	}

	// Here It Creates a New Directory with the name as Category Name

	//Get Path to store image from conf

	catogory_image_path := beego.AppConfig.String("CATEGORY_IMAGE_PATH")

	err1 := os.Mkdir(catogory_image_path+input_category_name, 0755)

	if err != nil {
		// log.Fatal(err)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "err1", err1)
	}

	file, information, err := c.GetFile("myFile")

	if err != nil {
		// log.Println("File retrieval failure")
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "-", "File retrieval failure")
		return
	}
	defer file.Close()

	filename := information.Filename

	err = c.SaveToFile("myFile", path.Join(catogory_image_path+input_category_name, filename))

	if err != nil {
		// log.Println("File upload failed!")
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "-", "File upload failed!")
	} else {
		// log.Println("File upload succeed!")
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "-", "File upload succeed!")
		// c.Data["DisplayMessage"] = "File upload succeed!"
		//Display information after successful upload
	}

	filepath := "/static/dist/img/category_images/" + input_category_name + "/" + filename

	result, err := db.Db.Exec(`INSERT INTO Category (
	name, 
	status,
	image, 
	created_by, 
	create_date)
		VALUES (?,?,?,?,now())`,
		input_category_name,
		input_category_status,
		filepath,
		id)

	if err != nil {
		err = errors.New("Category creation failed")
		return
	}

	i, err := result.RowsAffected()
	if err != nil || i == 0 {
		err = errors.New("Category creation failed")
		return
	}

	return
}

func CheckUserAlreadyExists(email string) (err error) {

	err = nil

	row, err := db.Db.Query(`select count(*) from Category WHERE name = ?`, email)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch Category")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch Category")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	countlen := data[0][0]

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "countlen", countlen)

	if countlen != "0" {
		err = errors.New("Category already exists")
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}

	return

}
