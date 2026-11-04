package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"net/smtp"
	"strings"
	"text/template"

	"github.com/astaxie/beego/context"

	"c2f_admin_webapp/model/util/log"

	"c2f_admin_webapp/model/encoding/base64"

	"c2f_admin_webapp/model/util/pbkdf2"

	p "c2f_admin_webapp/model/util/password"

	"github.com/astaxie/beego"

	"net/url"

	"c2f_admin_webapp/model/db"

	"github.com/astaxie/beego/session"
)

type MenusStruct struct {
	Menus []string
}

type PayULogRequest struct {
	SessionID string
	Mobile    string
}

func GeneratePassword() (login_pass, encrypted_pwd string, err error) {
	login_pass, _ = p.Numeric(6)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Login Password", login_pass)

	b := make([]byte, 32)
	_, err = rand.Read(b)
	var pbkdf pbkdf2.Pbkdf2
	pbkdf.Itr = 32
	pbkdf.KeyLen = 32
	pbkdf.Plain = []byte(login_pass)
	pbkdf.Salt = b
	err = pbkdf.Encrypt()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to create password")
		return
	}
	var tmp []byte
	tmp = append(tmp, pbkdf.Salt...)
	tmp = append(tmp, pbkdf.Cipher...)

	out, err := base64.Encode(tmp)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to create password")
		return
	}
	encrypted_pwd = string(out)
	return
}

func IsAuthorized(role, page string) (result bool, err error) {

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "role", role)

	result = false
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cluster", beego.AppConfig.String("CLUSTERADMIN_MENU_ARRAY"))
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "report", beego.AppConfig.String("REPORTS_MENU_ARRAY"))
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "clusterexecutive", beego.AppConfig.String("CLUSTEREXECUTIVE_MENU_ARRAY"))
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "hofinanace", beego.AppConfig.String("HOFINANCE_MENU_ARRAY"))

	var menusjson string

	menusjson = ""

	if role == "HOADMIN" {
		menusjson = beego.AppConfig.String("ADMIN_MENU_ARRAY")
	} else if role == "CLUSTERMANAGER" {
		menusjson = beego.AppConfig.String("CLUSTERADMIN_MENU_ARRAY")
	} else if role == "REPORTS" {
		menusjson = beego.AppConfig.String("REPORTS_MENU_ARRAY")
	} else if role == "CLUSTEREXECUTIVE" {
		menusjson = beego.AppConfig.String("CLUSTEREXECUTIVE_MENU_ARRAY")
	} else if role == "HOHELPDESK" {
		menusjson = beego.AppConfig.String("HOHELPDESK_MENU_ARRAY")
	} else if role == "HOFINANCE" {
		menusjson = beego.AppConfig.String("HOFINANCE_MENU_ARRAY")
	}

	if role == "" {
		err = errors.New("Role invalid")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "menusjson", menusjson)

	var menus MenusStruct
	err = json.Unmarshal([]byte(menusjson), &menus)
	if err != nil {
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "menus", menus)

	for _, men := range menus.Menus {
		if strings.EqualFold(men, page) {
			result = true
		}
	}
	return
}

func SendEmail(host string, port int, userName string, password string, to []string, subject string, message string) (err error) {
	parameters := struct {
		From    string
		To      string
		Subject string
		Message string
	}{
		userName,
		strings.Join([]string(to), ","),
		subject,
		message,
	}

	buffer := new(bytes.Buffer)

	template := template.Must(template.New("emailTemplate").Parse(emailScript()))
	template.Execute(buffer, &parameters)

	auth := smtp.PlainAuth("", userName, password, host)

	err = smtp.SendMail(
		fmt.Sprintf("%s:%d", host, port),
		auth,
		userName,
		to,
		buffer.Bytes())

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	}
	return err
}
func emailScript() (script string) {
	return `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}
MIME-version: 1.0
Content-Type: text/html; charset="UTF-8"

{{.Message}}`
}

func SetHTTPHeader(Ctx *context.Context) {
	Ctx.Output.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	Ctx.Output.Header("Pragma", "no-cache")
	Ctx.Output.Header("Expires", "0")
	Ctx.Output.Header("X-Content-Type-Options", "nosniff")
	Ctx.Output.Header("Strict-Transport-Security", "max-age=31536000 ; includeSubDomains")
	Ctx.Output.Header("X-Frame-Options", "SAMEORIGIN")
	Ctx.Output.Header("X-XSS-Protection", "1; mode=block")
	Ctx.Output.Header("X-Content-Security-Policy", "default-src 'self'")
	//Ctx.Output.Header("X-WebKit-CSP", "default-src 'self'")
}

func EventLogs(c *context.Context, sess session.Store, method string, input url.Values, output map[interface{}]interface{}, err_r error) (res string, err error) {
	m2 := make(map[string]interface{})

	for key, value := range output {
		switch key := key.(type) {
		case string:
			m2[key] = value
		}
	}
	event := make(map[string]interface{})

	event["PIP"] = c.Input.IP()
	event["URL"] = c.Input.URL()
	event["SessionID"] = sess.SessionID()
	if sess.Get("uname") != nil {
		event["UserName"] = sess.Get("uname")
	}
	event["Host"] = c.Input.Host()
	event["Method"] = method
	if err_r != nil {
		event["Status"] = err_r

	} else {
		event["Status"] = "Success"
	}

	event["Input"] = input
	event["Output"] = m2

	jsonString, err := json.Marshal(event)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Host Data - ", c.Input.Host())
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "URL Data - ", c.Input.URL())
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "User Data - ", sess.Get("uname"))
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "method Data - ", method)

	event_json := string(jsonString)

	_, err = db.Db.Exec("INSERT INTO Web_event(event,created_on,user_id,url,ip,host,status) VALUES ( ?,now(),?,?,?,?,?)", event_json, sess.Get("uname"), c.Input.URL(), c.Input.IP(), c.Input.Host(), event["Status"])

	return
}

func PayUEventLogs(
	c *context.Context,
	sessionID string,
	mobile string,
	method string,
	request string,
	response string,
	input url.Values,
	output map[interface{}]interface{},
	err_r error,
) (res string, err error) {

	m2 := make(map[string]interface{})

	for key, value := range output {
		if k, ok := key.(string); ok {
			m2[k] = value
		}
	}

	event := make(map[string]interface{})

	event["PIP"] = c.Input.IP()
	event["URL"] = c.Input.URL()
	event["SessionID"] = sessionID

	if mobile != "" {
		event["UserName"] = mobile
	}

	event["Host"] = c.Input.Host()
	event["Method"] = method

	if err_r != nil {
		event["Status"] = err_r.Error()
	} else {
		event["Status"] = "Success"
	}

	event["Input"] = input
	event["Output"] = m2

	jsonString, err := json.Marshal(event)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}

	event_json := string(jsonString)

	log.Println(beego.AppConfig.String("loglevel"), "Debug", response)

	_, err = db.Db.Exec(
		"INSERT INTO Events(event, request, response, dump_time) VALUES (?,?,?,now())",
		event_json, request, response,
	)

	return
}
