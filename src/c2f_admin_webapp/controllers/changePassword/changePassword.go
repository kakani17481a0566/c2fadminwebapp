package changePassword

import (
	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"
	"c2f_admin_webapp/session"

	"crypto/rand"
	"errors"
	"runtime/debug"

	//	"unicode"
	"regexp"

	"net/smtp"
	"strings"

	"c2f_admin_webapp/model/util/pbkdf2"

	"c2f_admin_webapp/model/db"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/encoding/base64"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

type ChangePassword struct {
	// beego.Controller
	middleware.BaseController
}

type unencryptedAuth struct {
	smtp.Auth
}

func (c *ChangePassword) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "ChangePassword Start")
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
			c.TplName = "changePassword/changePassword.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "changePassword/changePassword.html"
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
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "utype", utype)

	language := sess.Get("language").(string)
	c.Data["language"] = language
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "language", language)

	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1

	auth, err := utils.IsAuthorized(utype, "ChangePassword")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	return
}

func (c *ChangePassword) Post() {

	log.Println(beego.AppConfig.String("loglevel"), "Info", "ChangePassword post page")
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
			c.TplName = "changePassword/changePassword.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Password Changed Successfully"
			c.TplName = "changePassword/changePassword.html"
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
	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "utype", utype)

	language := sess.Get("language").(string)
	c.Data["language"] = language
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "language", language)

	auth, err := utils.IsAuthorized(utype, "ChangePassword")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}

	// old_password := c.Input().Get("input_old_password")
	// //log.Println(beego.AppConfig.String("loglevel"), "Debug", "old_password - ", old_password)

	// new_password := c.Input().Get("input_new_password")
	// //log.Println(beego.AppConfig.String("loglevel"), "Debug", "new_password - ", new_password)

	// confirm_password := c.Input().Get("input_confirm_password")
	// //log.Println(beego.AppConfig.String("loglevel"), "Debug", "confirm_password - ", confirm_password)

	// uname := sess.Get("uname")

	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "User Name - ", uname)

	old_password := c.Input().Get("input_old_password")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "old_password - ", old_password)

	new_password := c.Input().Get("input_new_password")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "new_password - ", new_password)

	confirm_password := c.Input().Get("input_confirm_password")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "confirm_password - ", confirm_password)

	uname := sess.Get("uname")
	// language := sess.Get("language").(string)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "User Name - ", uname)

	if new_password == "" || confirm_password == "" {
		log.Println(beego.AppConfig.String("loglevel"), "Error", "Blank Passwords")
		err = errors.New("Passwords can't be blank")
		return
	}

	if new_password != confirm_password {
		log.Println(beego.AppConfig.String("loglevel"), "Error", "New Password Mismatch")
		err = errors.New("New passwords can't be different")
		return
	}

	if len(new_password) < 8 {
		err = errors.New("Password must be at least 8 characters long")
		return
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(new_password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(new_password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(new_password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(new_password)

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		err = errors.New("Password must contain uppercase, lowercase, number, and special character")
		return
	}

	err = Authenticate(uname.(string), old_password, "ADMIN", "")

	if err != nil {

		if sess.Get("language") == "English" {
			err = errors.New(beego.AppConfig.String("ENGLISH_ADMIN_AUTHENTICATION_FAILED"))
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			return
		}
		if sess.Get("language") == "French" {
			err = errors.New(beego.AppConfig.String("FRENCH_ADMIN_AUTHENTICATION_FAILED"))
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			return
		}
		// log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		// err = errors.New("Admin Authentication Failed")
		// return
	}

	err = UpdatePassword(uname.(string), new_password)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Admin Update Password Failed")
		return
	}

	c.Data["MenuJson"] = sess.Get("menujson")

	return
}

func Authenticate(uname, pass, userType, lang string) (err error) {

	row, err := db.Db.Query("SELECT id,password,status FROM User where email=? limit 1", uname)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System User Not Found")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System User Detail Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("System User Not Found")
		return
	}

	cp, err := base64.Decode([]byte(data[0][1]))
	if err != nil {
		err = errors.New("System User Password Decoding Fail")
		return
	}

	var pbkdf pbkdf2.Pbkdf2
	pbkdf.Itr = 32
	pbkdf.KeyLen = 32
	pbkdf.Plain = []byte(pass)
	pbkdf.Salt = cp[:32]
	pbkdf.Cipher = cp[32:]
	result, err := pbkdf.Compare()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System User Password MisMatch")
		return
	}

	if !result {
		log.Println(beego.AppConfig.String("loglevel"), "Error", result)
		err = errors.New("System User Password MisMatch")
		return
	}

	if data[0][2] != "ACTIVE" {

		err = errors.New("System User Not Active")

		return
	}

	return
}

func UpdatePassword(uname, password string) (err error) {

	b := make([]byte, 32)
	_, err = rand.Read(b)
	var pbkdf pbkdf2.Pbkdf2
	pbkdf.Itr = 32
	pbkdf.KeyLen = 32
	pbkdf.Plain = []byte(password)
	pbkdf.Salt = b
	err = pbkdf.Encrypt()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System User Password Update Fail")
		return
	}
	var tmp []byte
	tmp = append(tmp, pbkdf.Salt...)
	tmp = append(tmp, pbkdf.Cipher...)

	out, err := base64.Encode(tmp)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System User Password Update Fail")
		return
	}

	result, err := db.Db.Exec("update User set password=?, update_date=now() where email=? ", string(out), uname)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System User Password Update Fail")
		return
	}
	n, err := result.RowsAffected()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System User Password Update Fail")
		return
	}

	if n != 1 {
		err = errors.New("System User Password Update Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "pass", password)
	return
}
