package createSupplyExecutive

import (
	"c2f_admin_webapp/session"

	"runtime/debug"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"
	"errors"

	"c2f_admin_webapp/model/db"
	"strings"

	"github.com/astaxie/beego"

	// "c2f_admin_webapp/services"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
	"c2f_admin_webapp/model/util/password"
)

type User struct {
	Id string `json:"Id,omitempty"`
}
type CreateSupplyExecutive struct {
	// beego.Controller
	middleware.BaseController
}

func (c *CreateSupplyExecutive) Get() {
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
			c.TplName = "supplyexecutive/createSupplyExecutive/createSupplyExecutive.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Create SupplyExecutive Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "supplyexecutive/createSupplyExecutive/createSupplyExecutive.html"

			log.Println(beego.AppConfig.String("loglevel"), "Info", "Create SupplyExecutive  Page Success")
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

	auth, err := utils.IsAuthorized(utype, "SupplyExecutive")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	return
}
func (c *CreateSupplyExecutive) Post() {

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
			c.TplName = "supplyexecutive/createSupplyExecutive/createSupplyExecutive.html"

			log.Println(beego.AppConfig.String("loglevel"), "Info", "Create SupplyExecutive Page Fail")
		} else {
			c.Data["DisplayMessage"] = "SupplyExecutive created Successfully"
			c.TplName = "supplyexecutive/createSupplyExecutive/createSupplyExecutive.html"

			log.Println(beego.AppConfig.String("loglevel"), "Info", "Create SupplyExecutive  Page Success")
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

	auth, err := utils.IsAuthorized(utype, "SupplyExecutive")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	input_first_name := c.Input().Get("input_first_name")
	input_middle_name := c.Input().Get("input_middle_name")
	input_last_name := c.Input().Get("input_last_name")
	input_mobile := c.Input().Get("input_mobile")
	input_email := c.Input().Get("input_email")
	input_user_type := c.Input().Get("input_user_type")
	input_address1 := c.Input().Get("input_address1")
	input_address2 := c.Input().Get("input_address2")
	input_city := c.Input().Get("input_city")
	input_zipcode := c.Input().Get("input_zipcode")
	input_status := c.Input().Get("input_status")
	input_state := c.Input().Get("input_state")
	input_language := c.Input().Get("input_language")

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "user_type  - ", input_user_type)

	if input_first_name == "" || input_mobile == "" || input_email == "" || input_user_type == "" || input_address1 == "" || input_address2 == "" || input_city == "" || input_zipcode == "" ||
		input_status == "" || input_state == "" || input_language == "" {
		err = errors.New("* fields cannot be empty")
		return
	}
	if input_user_type == "" || input_language == "" {
		err = errors.New("user type and lanuguage cannot be empty")
		return
	}

	err = CheckUserAlreadyExists(input_mobile)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("SupplyExecutive already Exists")
		return
	}

	loginPass, _ := password.AlphaNumericSpecial(6)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "loginPass - ", loginPass)

	// pass, err := services.EncryptPassword(loginPass)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	return
	// }

	result, err := db.Db.Exec(`INSERT INTO SupplyExecutive (first_name,
		
		middle_name,
		last_name,
		mob_no,
		email,
		user_type,
		address1,
		address2,
		city,
		zipcode,
		state,
		language,
		status,
		create_date)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,now())`,
		input_first_name,
		input_middle_name,
		input_last_name,
		input_mobile,
		input_email,
		input_user_type,
		input_address1,
		input_address2,
		input_city,
		input_zipcode,
		input_state,
		input_language,
		input_status)
	if err != nil {
		err = errors.New("SupplyExecutive creation failed")
		return
	}

	i, err := result.RowsAffected()
	if err != nil || i == 0 {
		err = errors.New("SupplyExecutive creation failed")
		return
	}

	return
}

func CheckUserAlreadyExists(mob string) (err error) {

	err = nil

	row, err := db.Db.Query(`select count(*) from SupplyExecutive WHERE mob_no = ?`, mob)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch SupplyExecutive")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch SupplyExecutive")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	countlen := data[0][0]

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "countlen", countlen)

	if countlen != "0" {
		err = errors.New("SupplyExecutive already exists")
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}

	return

}
