package adminAduitReport

import (
	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/model/utils"
	"c2f_admin_webapp/services"
	"c2f_admin_webapp/session"

	"strings"

	"runtime/debug"

	"errors"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

type Row struct {
	TimeStamp   string
	Adminid     string
	URL         string
	Status      string
	IP          string
	Host        string
	HTTPSMethod string
	SessionID   string
}
type Display5 struct {
	Fields5 []Field5
}

type Field5 struct {
	Id     string
	Mobile string
}
type AdminAduitReport struct {
	// beego.Controller
	middleware.BaseController
}

func (c *AdminAduitReport) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Reports Page Start")
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
			c.TplName = "reports/adminAduitReport/adminAduitReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Reports Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "reports/adminAduitReport/adminAduitReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Reports Page Success")
		}
		return
	}()
	utils.SetHTTPHeader(c.Ctx)
	sess, err := session.GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System is unable to process your request.Please contact customer care")
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

	data, err := services.GetConsumer()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Producer fetch Failed")
		return
	}
	if len(data) <= 0 {
		err = errors.New("Producer  Not Found")
		return
	}
	var Dis Display5
	for i := range data {
		var d Field5
		d.Id = data[i][0]
		d.Mobile = data[i][1]
		Dis.Fields5 = append(Dis.Fields5, d)
	}
	c.Data["Dis"] = Dis
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

	auth, err := utils.IsAuthorized(utype, "AduitReport")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")
	// row, err := db.Db.Query(`select created_on,
	// JSON_UNQUOTE(JSON_EXTRACT(event, '$.UserName')),
	// JSON_UNQUOTE(JSON_EXTRACT(event, '$.URL')),
	// JSON_UNQUOTE(JSON_EXTRACT(event, '$.Status')),
	// JSON_UNQUOTE(JSON_EXTRACT(event, '$.PIP')),
	// JSON_UNQUOTE(JSON_EXTRACT(event, '$.Host')),
	// JSON_UNQUOTE(JSON_EXTRACT(event, '$.Method')),
	// JSON_UNQUOTE(JSON_EXTRACT(event, '$.SessionID'))
	// FROM Web_event`)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Unable to fetch buyer")
	// 	return
	// }
	// defer sql.Close(row)
	// _, data, err1 := sql.Scan(row)
	// if err1 != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
	// 	err1 = errors.New("Unable to fetch buyer")
	// 	return
	// }

	// log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	// var result []Row

	// for i := range data {
	// 	var r Row
	// 	r.TimeStamp = data[i][0]
	// 	r.Adminid = data[i][1]
	// 	r.URL = data[i][2]
	// 	r.Status = data[i][3]
	// 	r.IP = data[i][4]
	// 	r.Host = data[i][5]
	// 	r.HTTPSMethod = data[i][6]
	// 	r.SessionID = data[i][7]
	// 	result = append(result, r)
	// }
	// c.Data["CustomerData"] = result

	return
}

func (c *AdminAduitReport) Post() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchOrder Page Start")
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
			c.TplName = "reports/adminAduitReport/adminAduitReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchOrder Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "reports/adminAduitReport/adminAduitReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchOrder  Page Success")
		}
		return
	}()
	utils.SetHTTPHeader(c.Ctx)
	sess, err := session.GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System is unable to process your request.Please contact customer care")
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

	auth, err := utils.IsAuthorized(utype, "AduitReport")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	admin := c.Input().Get("admin_id")
	url := c.Input().Get("input_url")
	ipp := c.Input().Get("input_ip")
	status := c.Input().Get("input_status")
	host := c.Input().Get("input_host")
	//httpsmethod := c.Input().Get("input_https_method")

	c.Data["DateRange"] = c.Input().Get("daterange")
	tmp := strings.Split(c.Input().Get("daterange"), " - ")

	from := ""
	to := ""
	if len(tmp) == 2 {
		from = tmp[0]
		to = tmp[1]
	}

	row, err := db.Db.Query(`select created_on, 
JSON_UNQUOTE(JSON_EXTRACT(event, '$.UserName')),
JSON_UNQUOTE(JSON_EXTRACT(event, '$.URL')),
JSON_UNQUOTE(JSON_EXTRACT(event, '$.Status')),
JSON_UNQUOTE(JSON_EXTRACT(event, '$.PIP')),
JSON_UNQUOTE(JSON_EXTRACT(event, '$.Host')),
JSON_UNQUOTE(JSON_EXTRACT(event, '$.Method')),
JSON_UNQUOTE(JSON_EXTRACT(event, '$.SessionID'))
FROM Web_event WHERE (user_id ='' OR user_id like ?)
AND(url ='' OR url like ?)AND
(ip='' OR ip like ?)AND
(host='' OR host like ?)AND
(status='' OR status like ?)AND
DATE(created_on) >= STR_TO_DATE(?, '%m/%d/%Y') AND DATE(created_on) <= STR_TO_DATE(?, '%m/%d/%Y') ORDER BY created_on DESC`, admin+"%", url+"%", ipp+"%", host+"%", status+"%", from, to)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch buyer")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch buyer")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	var result []Row

	for i := range data {
		var r Row
		r.TimeStamp = data[i][0]
		r.Adminid = data[i][1]
		r.URL = data[i][2]
		r.Status = data[i][3]
		r.IP = data[i][4]
		r.Host = data[i][5]
		r.HTTPSMethod = data[i][6]
		r.SessionID = data[i][7]
		result = append(result, r)
	}
	c.Data["CustomerData"] = result

	return
}
