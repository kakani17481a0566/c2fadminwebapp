package createChannelPartnerType

import (
	"c2f_admin_webapp/session"

	"runtime/debug"

	"c2f_admin_webapp/model/utils"
	"errors"

	"c2f_admin_webapp/model/db"

	"strings"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
	//	"c2f_admin_webapp/model/util/sku"
)

type CreateChannelPartnerType struct {
	// beego.Controller
	middleware.BaseController
}

func (c *CreateChannelPartnerType) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Create Channel Partner Type Start")
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
			c.TplName = "channelpartnertype/createChannelPartnerType/createChannelPartnerType.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Create Channel Partner Type Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "channelpartnertype/createChannelPartnerType/createChannelPartnerType.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Create Channel Partner Type  Page Success")
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

	auth, err := utils.IsAuthorized(utype, "ChannelPartnerType")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	return
}
func (c *CreateChannelPartnerType) Post() {

	log.Println(beego.AppConfig.String("loglevel"), "Info", "Craete Channel Partner Type post page")
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
			c.TplName = "channelpartnertype/createChannelPartnerType/createChannelPartnerType.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Channel Partner Type Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Channe Partner Type Created Successfully"
			c.TplName = "channelpartnertype/createChannelPartnerType/createChannelPartnerType.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Channel Partner Type  Page Success")
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

	auth, err := utils.IsAuthorized(utype, "ChannelPartnerType")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	input_channel_partner_type_name := c.Input().Get("input_channel_partner_type_name")
	input_channel_partner_type_desc := c.Input().Get("input_channel_partner_type_desc")
	input_channel_type_status := c.Input().Get("input_channel_type_status")

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "channel partner type name  - ", input_channel_partner_type_name)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "channel partner type desc  - ", input_channel_partner_type_desc)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "channel partner type status  - ", input_channel_type_status)

	if input_channel_partner_type_name == "" || input_channel_partner_type_desc == "" || input_channel_type_status == "" {
		err = errors.New("*All Fileds cannot be empty")
		return
	}

	err = CheckUserAlreadyExists(input_channel_partner_type_name)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Category already Exists")
		return
	}

	result, err := db.Db.Exec(`INSERT INTO ChannelPartnerType (
	name, 
	description,
	status,
	created_by, 
	created_date)
		VALUES (?,?,?,?,now())`,
		input_channel_partner_type_name,
		input_channel_partner_type_desc,
		input_channel_type_status,
		id)

	if err != nil {
		err = errors.New("Channel Partner Type Creation Failed")
		return
	}

	i, err := result.RowsAffected()
	if err != nil || i == 0 {
		err = errors.New("Channel Partner Type Creation Failed")
		return
	}

	return
}

func CheckUserAlreadyExists(name string) (err error) {

	err = nil

	row, err := db.Db.Query(`select count(*) from ChannelPartnerType WHERE name = ?`, name)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch Channel Partner Type")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch Channel Partner Type")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	countlen := data[0][0]

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "countlen", countlen)

	if countlen != "0" {
		err = errors.New("Channel Partner Type already exists")
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}

	return

}
