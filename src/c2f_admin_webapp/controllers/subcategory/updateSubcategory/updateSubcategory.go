package updateSubcategory

import (
	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/model/utils"
	"c2f_admin_webapp/session"
	"errors"
	"path"
	"runtime/debug"
	"strings"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"

	"github.com/astaxie/beego"
)

type Row struct {
	Id                string
	SubcategoryName   string
	Subcategorystatus string
	Image             string
	//Timestamp         string
}
type UpdateSubcategory struct {
	// beego.Controller
	middleware.BaseController
}

func (c *UpdateSubcategory) Get() {
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
			c.TplName = "subcategory/updateSubcategory/updateSubcategory.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "subcategory/updateSubcategory/updateSubcategory.html"
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

	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1

	auth, err := utils.IsAuthorized(utype, "Subcategory")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")
	row, err := db.Db.Query(`select id,name,status,create_date,image from SubCategory WHERE id =?  limit 1`, AdminId)
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
		c.Data["SubcategoryName"] = data[i][1]
		c.Data["Subcategorystatus"] = data[i][2]
		c.Data["Image"] = data[i][4]
		//c.Data["Timestamp"] = data[i][3]
	}
	return
}
func (c *UpdateSubcategory) Post() {
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
			c.TplName = "subcategory/updateSubcategory/updateSubcategory.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Sub Category Updated Successfully"
			c.TplName = "subcategory/updateSubcategory/updateSubcategory.html"
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

	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1

	id := sess.Get("uid").(string)
	c.Data["id"] = id
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "id", id)

	auth, err := utils.IsAuthorized(utype, "Subcategory")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")
	input_subcategory_name := c.Input().Get("input_subcategory_name")
	input_subcategory_status := c.Input().Get("input_subcategory_status")

	res3, err4 := db.Db.Exec(`UPDATE SubCategory SET 
		name = ?, status = ?,updated_by=?,update_date=now()
		WHERE id = ?`, input_subcategory_name, input_subcategory_status, id, AdminId)
	if err4 != nil {
		err4 = errors.New("Status updation failed")
		return
	}
	i, err5 := res3.RowsAffected()
	if err5 != nil || i == 0 {
		err5 = errors.New("Status updation failed")
		return
	}

	file, information, err1 := c.GetFile("myFile")

	if err1 != nil {
		// log.Println("File retrieval failure")
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "-", "File retrieval failure")
		return
	} else {

		defer file.Close()

		filename := information.Filename

		sub_catogory_image_path := beego.AppConfig.String("SUB_CATEGORY_IMAGE_PATH")

		err1 = c.SaveToFile("myFile", path.Join(sub_catogory_image_path+input_subcategory_name, filename))

		if err1 != nil {
			// log.Println("File upload failed!")
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "-", "File upload failed!")
		} else {
			// log.Println("File upload succeed!")
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "-", "File upload succeed!")
			// c.Data["DisplayMessage"] = "File upload succeed!"
			//Display information after successful upload
		}

		filepath := "/static/dist/img/subcategory_images/" + input_subcategory_name + "/" + filename

		res, err2 := db.Db.Exec(`UPDATE SubCategory SET image=?  WHERE id=?`, filepath, AdminId)
		if err2 != nil {
			err2 = errors.New("SubCategory image path updation failed")
			return
		}
		i, err3 := res.RowsAffected()
		if err3 != nil || i == 0 {
			err3 = errors.New("SubCategory image  path updation failed")

			return
		}
	}

	row, err := db.Db.Query(`select id,name,status from SubCategory WHERE id = ? limit 1`, AdminId)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get SubCategory data")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get SubCategory data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("SubCategory data not found")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)
	for i := range data {
		c.Data["Id"] = data[i][0]
		c.Data["SubcategoryName"] = data[i][1]
		c.Data["Subcategorystatus"] = data[i][2]
		//c.Data["Timestamp"] = data[i][3]
	}
	return
}
