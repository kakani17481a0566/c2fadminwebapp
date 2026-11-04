package searchChannelPartner

import (
	"c2f_admin_webapp/session"

	"runtime/debug"

	"c2f_admin_webapp/model/utils"
	"errors"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/services"

	"c2f_admin_webapp/model/db"

	"strings"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

type Row struct {
	Id                 string
	Timestamp          string
	FirstName          string
	MiddleName         string
	LastName           string
	Mobile             string
	Email              string
	Usertype           string
	Status             string
	ClusterName        string
	ChannelPartnerType string
}

type SearchChannelPartner struct {
	// beego.Controller
	middleware.BaseController
}

type Display struct {
	Fields []Field
}
type Field struct {
	Id    string
	Name  string
	Email string
}

func (c *SearchChannelPartner) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Channel Partner Page Start")
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
			c.TplName = "channelpartner/searchChannelPartner/searchChannelPartner.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Channel Partner Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "channelpartner/searchChannelPartner/searchChannelPartner.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Channel Partner Page Success")
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

	id := sess.Get("uid").(string)
	c.Data["id"] = id

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

	auth, err := utils.IsAuthorized(utype, "ChannelPartner")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	row, err := db.Db.Query(`select ChannelPartner.id,ChannelPartner.first_name,ChannelPartner.mob_no,ChannelPartner.email,ChannelPartner.status,
		ChannelPartner.create_date,ChannelPartnerType.name from ChannelPartner
		LEFT JOIN ChannelPartnerType ON ChannelPartner.channel_partner_type_id = ChannelPartnerType.id  
		 ORDER BY ChannelPartner.id DESC`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get the error message info")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Message Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the error message info")
		return
	}

	var result []Row

	for i := range data {
		var r Row
		r.Id = data[i][0]
		r.FirstName = data[i][1]
		r.Mobile = data[i][2]
		r.Email = data[i][3]
		r.Status = data[i][4]
		r.Timestamp = data[i][5]
		r.ChannelPartnerType = data[i][6]
		result = append(result, r)
	}

	c.Data["CustomerData"] = result

	data1, err := services.GetChannelPartnerType()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Channel Partner Type fetch Failed")
		return
	}
	if len(data1) <= 0 {
		err = errors.New("Channel Partner Type Not Found")
		return
	}
	var Dis1 Display
	for i := range data1 {
		var d Field
		d.Id = data1[i][0]
		d.Name = data1[i][1]
		Dis1.Fields = append(Dis1.Fields, d)
	}
	c.Data["Dis1"] = Dis1
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis1)

	return
}

func (c *SearchChannelPartner) Post() {
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
			c.TplName = "channelpartner/searchChannelPartner/searchChannelPartner.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Channel Partner Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "channelpartner/searchChannelPartner/searchChannelPartner.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Channel Partner Page Success")
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

	id := sess.Get("uid").(string)
	c.Data["id"] = id

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

	auth, err := utils.IsAuthorized(utype, "ChannelPartner")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	input_name := c.Input().Get("input_name")
	input_email := c.Input().Get("input_email")
	input_status := c.Input().Get("input_status")
	input_mob := c.Input().Get("input_mob")
	// input_cluster_name := c.Input().Get("input_cluster_name")
	select_channel_partner_type := c.Input().Get("select_channel_partner_type")

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_name - ", input_name)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_email - ", input_email)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_status - ", input_status)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "select_channel_partner_type - ", select_channel_partner_type)

	c.Data["DateRange"] = c.Input().Get("daterange")
	tmp := strings.Split(c.Input().Get("daterange"), " - ")

	from := ""
	to := ""
	if len(tmp) == 2 {
		from = tmp[0]
		to = tmp[1]
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "From", from)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "To", to)

	data1, err := services.GetChannelPartnerType()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Channel Partner Type fetch Failed")
		return
	}
	if len(data1) <= 0 {
		err = errors.New("Channel Partner Type Not Found")
		return
	}
	var Dis1 Display
	for i := range data1 {
		var d Field
		d.Id = data1[i][0]
		d.Name = data1[i][1]
		Dis1.Fields = append(Dis1.Fields, d)
	}
	c.Data["Dis1"] = Dis1
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis1)

	row, err := db.Db.Query(`select ChannelPartner.id,ChannelPartner.first_name,ChannelPartner.mob_no,ChannelPartner.email,ChannelPartner.status,
		ChannelPartner.create_date,ChannelPartnerType.name from ChannelPartner
		LEFT JOIN ChannelPartnerType ON ChannelPartner.channel_partner_type_id = ChannelPartnerType.id 
		where (ChannelPartner.first_name='' OR ChannelPartner.first_name like ?)AND (ChannelPartner.status='' OR ChannelPartner.status like ?)AND (ChannelPartner.email='' OR ChannelPartner.email like ?) AND (ChannelPartner.mob_no='' OR ChannelPartner.mob_no like ?)AND (ChannelPartnerType.name='' OR ChannelPartnerType.name like ?) AND (DATE(ChannelPartner.create_date='') OR DATE (ChannelPartner.create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(ChannelPartner.create_date='') OR DATE (ChannelPartner.create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) ORDER BY id DESC`, input_name+"%", input_status+"%", input_email+"%", input_mob+"%", select_channel_partner_type+"%", from, to)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get the error message info")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Message Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the error message info")
		return
	}

	var result []Row

	for i := range data {
		var r Row
		r.Id = data[i][0]
		r.FirstName = data[i][1]
		r.Mobile = data[i][2]
		r.Email = data[i][3]
		r.Status = data[i][4]
		r.Timestamp = data[i][5]
		r.ChannelPartnerType = data[i][6]
		result = append(result, r)
	}

	c.Data["CustomerData"] = result

	return
}
