package forgotpassword

import (
	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/model/utils"
	"c2f_admin_webapp/session"
	"crypto/rand"
	"errors"
	"net/smtp"
	"runtime/debug"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/encoding/base64"

	"c2f_admin_webapp/model/util/log"

	"c2f_admin_webapp/model/util/password"

	"c2f_admin_webapp/model/util/pbkdf2"

	"github.com/astaxie/beego"
	"github.com/scorredoira/email"

	"io/ioutil"
	"net/mail"

	"strings"
)

type unencryptedAuth struct {
	smtp.Auth
}
type Forgotpassword struct {
	// beego.Controller
	middleware.BaseController
}

func (c *Forgotpassword) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	var err error
	sessErr := false
	defer func() {

		if l_exception := recover(); l_exception != nil {
			stack := debug.Stack()
			log.Println(beego.AppConfig.String("loglevel"), "Exception", string(stack))
			session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
			c.Data["DisplayMessage"] = "Something went wrong.Please Contact CustomerCare."
			c.TplName = "error/error.html"
		}
		if err != nil {
			if sessErr == true {
				log.Println(beego.AppConfig.String("loglevel"), "Info", "Redirecting to login")
				c.Redirect(beego.AppConfig.String("LOGIN_PATH"), 302)

			} else {
				c.Data["DisplayMessage"] = err.Error()
			}
			c.TplName = "forgotpassword/forgotpassword.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Forgot Password  Page Fail")
		} else {

			c.TplName = "forgotpassword/forgotpassword.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Forgot Password  Page Success")
		}
		return
	}()

	return
}

func (c *Forgotpassword) Post() {

	log.Println(beego.AppConfig.String("loglevel"), "Info", "Forgot Password post page")
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
			c.TplName = "forgotpassword/forgotpassword.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Forgot Password Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Password has been reset and sent on your registered email address successfully."
			c.TplName = "forgotpassword/forgotpassword.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Forgot Password  Page Success")
		}
		return
	}()

	utils.SetHTTPHeader(c.Ctx)
	uname := c.Input().Get("username")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "User Email - ", uname)
	if uname == "" {
		log.Println(beego.AppConfig.String("loglevel"), "Error", "Blank User Name")
		err = errors.New("UserName can't be blank.")
		return
	}

	row, err := db.Db.Query("SELECT id,email,first_name,status FROM User where email=? limit 1", uname)
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

	fullname := data[0][2]
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "fullname - ", fullname)

	newPass, _ := password.AlphaNumericSpecial(6)

	err = SearchUser(uname)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("User Not Found")
		return
	}

	err = UpdatePassword(uname, newPass)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Admin Update Password Failed")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "pass", newPass)

	err = SendEmail(uname, fullname, newPass)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("USER_SENDMAIL__NOT_FOUND")
		return
	}

	return
}

func SearchUser(uname string) (err error) {

	row, err := db.Db.Query("SELECT id,email,first_name,status FROM User where email=? limit 1", uname)
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

	if data[0][3] != "ACTIVE" {
		err = errors.New("User Not Active")
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
		err = errors.New(" User Password Update Fail")
		return
	}
	var tmp []byte
	tmp = append(tmp, pbkdf.Salt...)
	tmp = append(tmp, pbkdf.Cipher...)

	out, err := base64.Encode(tmp)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New(" User2 Password Update Fail")
		return
	}

	result, err := db.Db.Exec("update User set password=?, update_date=now(),password_set=? where email=? ", string(out), "NO", uname)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New(" User3 Password Update Fail")
		return
	}
	n, err := result.RowsAffected()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New(" User4 Password Update Fail")
		return
	}

	if n != 1 {
		err = errors.New("System User Password Update Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "pass", password)
	return
}

func SendEmail(emilid string, name string, password string) (err error) {
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "called - ")

	uname := beego.AppConfig.String("EMAIL_NOTIFY_USERNAME")
	pass := beego.AppConfig.String("EMAIL_NOTIFY_PASSWORD")
	url := beego.AppConfig.String("EMAIL_NOTIFY_URL")
	to := beego.AppConfig.String("EMAIL_NOTIFY_TIMEOUT")
	recipients := strings.Split(emilid, "||")
	loginurl := beego.AppConfig.String("EMAIL_APPLICATION_LOGIN_URL")

	tmpFile := beego.AppConfig.String("PASSWORD_RESET_EMAIL_TEMPLATE")

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

	m := email.NewHTMLMessage("Krishivala Password Reset: Here are your new login details", msg)
	m.From = mail.Address{Name: "Krishivala", Address: uname}
	m.To = recipients

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
