package createsysUser

import (
	"c2f_admin_webapp/session"

	"io/ioutil"
	"net/mail"
	"net/smtp"
	"runtime/debug"
	"strings"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/model/utils"
	"errors"

	"github.com/scorredoira/email"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/services"

	"c2f_admin_webapp/model/database/sql"
	"c2f_admin_webapp/model/util/log"
	"c2f_admin_webapp/model/util/password"
)

type CreatesysUser struct {
	// beego.Controller
	middleware.BaseController
}

type unencryptedAuth struct {
	smtp.Auth
}
type Display struct {
	Fields []Field
}

type Field struct {
	Id    string
	Name  string
	Email string
}

func (c *CreatesysUser) Get() {
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
			c.TplName = "sysusers/createsysUsers/createsysUsers.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "sysusers/createsysUsers/createsysUsers.html"
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

	maptoken := sess.Get("token")
	c.Data["maptoken"] = maptoken
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "maptoken", maptoken)

	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1

	data, err := services.GetRoles()
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
		Dis.Fields = append(Dis.Fields, d)
	}
	c.Data["Dis"] = Dis
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

	auth, err := utils.IsAuthorized(utype, "SysUser")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	return
}
func (c *CreatesysUser) Post() {

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
			c.TplName = "sysusers/createsysUsers/createsysUsers.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = "User Added Successfully"
			c.TplName = "sysusers/createsysUsers/createsysUsers.html"
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

	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1

	auth, err := utils.IsAuthorized(utype, "SysUser")
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
	select_role_name := c.Input().Get("select_role_name")
	input_address1 := c.Input().Get("input_address1")
	input_address2 := c.Input().Get("input_address2")
	input_city := c.Input().Get("input_city")
	input_zipcode := c.Input().Get("input_zipcode")
	input_status := c.Input().Get("input_status")
	input_state := c.Input().Get("input_state")
	input_language := c.Input().Get("input_language")
	input_town := c.Input().Get("input_town")
	input_location_type := c.Input().Get("input_location_type")
	input_location_info := c.Input().Get("input_location_info")

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "role  - ", select_role_name)

	if input_first_name == "" || input_mobile == "" || input_email == "" || select_role_name == "" ||
		input_address1 == "" || input_address2 == "" || input_city == "" || input_zipcode == "" ||
		input_status == "" || input_state == "" || input_language == "" || input_town == "" ||
		input_location_type == "" || input_location_info == "" {
		err = errors.New("* fields cannot be empty")
		return
	}
	// if input_town == "" || input_zipcode == "" {
	// 	err = errors.New("town and zipcode cannot be empty")
	// 	return
	// }

	err = CheckUserAlreadyExists(input_email)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("User already Exists")
		return
	}

	loginPass, _ := password.AlphaNumericSpecial(6)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "loginPass - ", loginPass)

	pass, err := services.EncryptPassword(loginPass)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}

	result, err := db.Db.Exec(`INSERT INTO User (first_name,
		password,
		middle_name,
		last_name,
		mob_no,
		email,
		address1,
		address2,
		city,
		zipcode,
		state,
		language,
		status,
		role,
		town,
		location_type,
		location_info,
		create_date)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,now())`,
		input_first_name,
		pass,
		input_middle_name,
		input_last_name,
		input_mobile,
		input_email,
		input_address1,
		input_address2,
		input_city,
		input_zipcode,
		input_state,
		input_language,
		input_status,
		select_role_name,
		input_town,
		input_location_type,
		input_location_info)
	if err != nil {
		err = errors.New("User creation failed")
		return
	}

	i, err := result.RowsAffected()
	if err != nil || i == 0 {
		err = errors.New("User creation failed")
		return
	}

	go SendEmail(input_email, input_first_name, loginPass, beego.AppConfig.String("EMAIL_TEMPLATE"))

	return
}

func CheckUserAlreadyExists(email string) (err error) {

	err = nil

	row, err := db.Db.Query(`select count(*) from User WHERE email = ?`, email)

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
		err = errors.New("User already exists")
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}

	return

}

func SendEmail(emilid, name, password, emailtemplate string) {
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "called - ")

	uname := beego.AppConfig.String("EMAIL_NOTIFY_USERNAME")
	pass := beego.AppConfig.String("EMAIL_NOTIFY_PASSWORD")
	url := beego.AppConfig.String("EMAIL_NOTIFY_URL")
	to := beego.AppConfig.String("EMAIL_NOTIFY_TIMEOUT")
	loginurl := beego.AppConfig.String("EMAIL_APPLICATION_LOGIN_URL")
	recipients := strings.Split(emilid, "||")

	tmpFile := emailtemplate

	buff, err := ioutil.ReadFile(tmpFile)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", "read file -", err)
		return
	}

	msg := string(buff)
	msg = strings.Replace(string(msg), "{{.Name}}", name, -1)
	msg = strings.Replace(string(msg), "{{.Email}}", emilid, -1)
	msg = strings.Replace(string(msg), "{{.Password}}", password, -1)
	msg = strings.Replace(string(msg), "{{.LoginURL}}", loginurl, -1)

	m := email.NewHTMLMessage("Welcome to Krishivala! Here are your login details", msg)
	m.From = mail.Address{Name: "Krishivala", Address: uname}
	m.To = recipients

	// send it
	//auth := smtp.PlainAuth("", uname, pass, url)

	config := beego.AppConfig.String("EMAIL_AUTH_CONFIG_MODE")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "config", config)

	if config == "1" {
		auth := smtp.PlainAuth("", uname, pass, url)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "auth")
		if err = email.Send(url+":"+to, auth, m); err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			return
		}

	} else if config == "2" {

		auth := unencryptedAuth{
			smtp.PlainAuth(
				"",
				uname,
				pass,
				url,
			),
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "no tls auth")
		if err = email.Send(url+":"+to, auth, m); err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			return
		}
	} else {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "no auth")
		if err = email.Send(url+":"+to, nil, m); err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			return
		}
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Email sent successfully")
	return
}
