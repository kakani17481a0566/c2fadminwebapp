package searchChannelPartnerProduct

import (
	"c2f_admin_webapp/session"

	"runtime/debug"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"
	"errors"

	"c2f_admin_webapp/model/db"

	// "c2f_admin_webapp/services"
	"strings"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

type Row struct {
	Id                  string
	Timestamp           string
	ProductName         string
	ChannelPartnerEmail string
	Status              string
	Mobile              string
}

type SearchChannelPartnerProduct struct {
	// beego.Controller
	middleware.BaseController
}

func (c *SearchChannelPartnerProduct) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Channel Partner Product Page Start")
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
			c.TplName = "channelpartnerproduct/searchChannelPartnerProduct/searchChannelPartnerProduct.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Channel Partner Product Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "channelpartnerproduct/searchChannelPartnerProduct/searchChannelPartnerProduct.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Channel Partner Product  Page Success")
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

	auth, err := utils.IsAuthorized(utype, "ChannelPartnerProduct")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	row, err := db.Db.Query(`SELECT ChannelPartner2Product.id,ChannelPartner2Product.product_id, ChannelPartner2Product.channelpartner_id, 
        ChannelPartner2Product.status, ChannelPartner2Product.create_date,Lookup_products.product_name, 
        ChannelPartner.email,
        ChannelPartner.mob_no  
	    FROM ChannelPartner2Product   
	    LEFT JOIN Lookup_products ON ChannelPartner2Product.product_id = Lookup_products.id
 		LEFT JOIN ChannelPartner ON ChannelPartner2Product.channelpartner_id = ChannelPartner.id ORDER BY ChannelPartner2Product.create_date DESC`)

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
		r.ProductName = data[i][5]
		r.ChannelPartnerEmail = data[i][6]
		r.Status = data[i][3]
		r.Timestamp = data[i][4]
		r.Mobile = data[i][7]

		result = append(result, r)
	}

	c.Data["CustomerData"] = result

	return
}

func (c *SearchChannelPartnerProduct) Post() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Channel Partner Product page")
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
			c.TplName = "channelpartnerproduct/searchChannelPartnerProduct/searchChannelPartnerProduct.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Channel Partner Product Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "channelpartnerproduct/searchChannelPartnerProduct/searchChannelPartnerProduct.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Channel Partner Product Page Success")
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

	auth, err := utils.IsAuthorized(utype, "ChannelPartnerProduct")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	//input_name := c.Input().Get("input_name")
	input_channelpartneremail := c.Input().Get("input_channelpartneremail")
	input_product_name := c.Input().Get("input_product_name")
	input_status := c.Input().Get("input_status")
	input_mob := c.Input().Get("input_mob")

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_product_name - ", input_product_name)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_channelpartneremail - ", input_channelpartneremail)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_status - ", input_status)

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

	row, err := db.Db.Query(`SELECT ChannelPartner2Product.id,ChannelPartner2Product.product_id, ChannelPartner2Product.channelpartner_id, 
        ChannelPartner2Product.status, ChannelPartner2Product.create_date,Lookup_products.product_name, 
        ChannelPartner.email,
        ChannelPartner.mob_no  
	    FROM ChannelPartner2Product   
	    LEFT JOIN Lookup_products ON ChannelPartner2Product.product_id = Lookup_products.id
 		LEFT JOIN ChannelPartner ON ChannelPartner2Product.channelpartner_id = ChannelPartner.id 
		where (Lookup_products.product_name='' OR Lookup_products.product_name like ?)AND (ChannelPartner.email ='' OR ChannelPartner.email  like ?)  AND (ChannelPartner2Product.status ='' OR ChannelPartner2Product.status  like ?) AND (ChannelPartner.mob_no ='' OR ChannelPartner.mob_no  like ?)  AND (DATE(ChannelPartner2Product.create_date='') OR DATE (ChannelPartner2Product.create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(ChannelPartner2Product.create_date='') OR DATE (ChannelPartner2Product.create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) ORDER BY ChannelPartner2Product.create_date DESC`, input_product_name+"%", input_channelpartneremail+"%", input_status+"%", input_mob+"%", from, to)

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
		r.ProductName = data[i][5]
		r.ChannelPartnerEmail = data[i][6]
		r.Status = data[i][3]
		r.Timestamp = data[i][4]
		r.Mobile = data[i][7]

		result = append(result, r)
	}

	c.Data["CustomerData"] = result

	return
}
