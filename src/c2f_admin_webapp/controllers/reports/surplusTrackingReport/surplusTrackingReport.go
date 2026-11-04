package surplusTrackingReport

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
	Id                    string
	Timestamp             string
	ProducerName          string
	ProducerMobile        string
	ProducerEmail         string
	ProducerAdd1          string
	ProducerAdd2          string
	ProducerCity          string
	ProducerState         string
	ProducerZipcode       string
	ProducerStatus        string
	ProducerClusterName   string
	LocationInfo          string
	ChannelPartnerName    string
	ChannelPartnerMobile  string
	ChannelPartnerEmail   string
	Address1              string
	Address2              string
	City                  string
	State                 string
	Pincode               string
	ChannelPartnerStatus  string
	ChannelPartnerType    string
	ProductName           string
	ClusterName           string
	ProductSku            string
	ProducerAddress       string
	ChannelPartnerAddress string
	ProductQty            string
}

type Display5 struct {
	Fields5 []Field5
}

type Field5 struct {
	Id     string
	Mobile string
}
type Display struct {
	Fields []Field
}
type Field struct {
	Id    string
	Name  string
	Email string
}

type SurplusTrackingReport struct {
	// beego.Controller
	middleware.BaseController
}

func (c *SurplusTrackingReport) Get() {

	Status := c.Ctx.Input.Param(":Status")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Status", Status)
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
			c.TplName = "reports/surplusTrackingReport/surplusTrackingReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Reports Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "reports/surplusTrackingReport/surplusTrackingReport.html"
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
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "utype", utype)

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

	auth, err := utils.IsAuthorized(utype, "SplusTraReport")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	// HOADMIN CODE

	row, err := db.Db.Query(`SELECT 
	id,
    created_on,

    JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.Name')) AS producer_name,
    JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.Mobile')) AS producer_mobile,
    JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.Email')) AS producer_email,
    JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.Add1')) AS producer_add1,
    JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.Add2')) AS producer_add2,
    JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.City')) AS producer_city,
    JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.State')) AS producer_state,
    JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.Zipcode')) AS producer_zipcode,
    JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.Status')) AS producer_status,
    JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.ClusterName')) AS producer_cluster,
    JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.LocationInfo')) AS producer_locationinfo,

    JSON_UNQUOTE(JSON_EXTRACT(channel_partner_details, '$.Name')) AS partner_name,
    JSON_UNQUOTE(JSON_EXTRACT(channel_partner_details, '$.Mobile')) AS partner_mobile,
    JSON_UNQUOTE(JSON_EXTRACT(channel_partner_details, '$.Email')) AS partner_email,
    JSON_UNQUOTE(JSON_EXTRACT(channel_partner_details, '$.Address1')) AS partner_address1,
    JSON_UNQUOTE(JSON_EXTRACT(channel_partner_details, '$.Address2')) AS partner_address2,
    JSON_UNQUOTE(JSON_EXTRACT(channel_partner_details, '$.City')) AS partner_city,
    JSON_UNQUOTE(JSON_EXTRACT(channel_partner_details, '$.State')) AS partner_state,
    JSON_UNQUOTE(JSON_EXTRACT(channel_partner_details, '$.Pincode')) AS partner_pincode,
    JSON_UNQUOTE(JSON_EXTRACT(channel_partner_details, '$.Status')) AS partner_status,
    JSON_UNQUOTE(JSON_EXTRACT(channel_partner_details, '$.ChannelPartnerType')) AS partner_type,
    JSON_UNQUOTE(JSON_EXTRACT(product_details, '$.ProductName')) AS product_name,
    JSON_UNQUOTE(JSON_EXTRACT(product_details, '$.ProductSKU')) AS product_sku,
    JSON_UNQUOTE(JSON_EXTRACT(product_details, '$.ProductQty')) AS product_qty,
    cluster
	FROM Surplus_tracking`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch order")
		return
	}
	defer sql.Close(row)
	_, data, err1 := sql.Scan(row)
	if err1 != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
		err1 = errors.New("Unable to fetch order")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)
	var result []Row

	for i := range data {
		var r Row
		r.Id = data[i][0]
		r.Timestamp = data[i][1]
		r.ProducerName = data[i][2]
		r.ProducerMobile = data[i][3]
		r.ProducerEmail = data[i][4]
		r.ProducerAdd1 = data[i][5]
		r.ProducerAdd2 = data[i][6]
		r.ProducerCity = data[i][7]
		r.ProducerState = data[i][8]
		r.ProducerZipcode = data[i][9]
		r.ProducerStatus = data[i][10]
		r.ProducerClusterName = data[i][11]
		r.LocationInfo = data[i][12]
		r.ChannelPartnerName = data[i][13]
		r.ChannelPartnerMobile = data[i][14]
		r.ChannelPartnerEmail = data[i][15]
		r.Address1 = data[i][16]
		r.Address2 = data[i][17]
		r.City = data[i][18]
		r.State = data[i][19]
		r.Pincode = data[i][20]
		r.ChannelPartnerStatus = data[i][21]
		r.ChannelPartnerType = data[i][22]
		r.ProductName = data[i][23]
		r.ProductSku = data[i][24]
		r.ProductQty = data[i][25]
		r.ClusterName = data[i][26]

		r.ProducerAddress = data[i][5] + ", " + data[i][6] + ", " + data[i][7] + ", " + data[i][8] + ", " + data[i][9]

		r.ChannelPartnerAddress = data[i][16] + ", " + data[i][17] + ", " + data[i][18] + ", " + data[i][19] + ", " + data[i][20]

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

func (c *SurplusTrackingReport) Post() {
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
			c.TplName = "reports/surplusTrackingReport/surplusTrackingReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchOrder Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "reports/surplusTrackingReport/surplusTrackingReport.html"
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

	auth, err := utils.IsAuthorized(utype, "SplusTraReport")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}

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

	input_producer_name := c.Input().Get("input_producer_name")
	input_producer_mobile := c.Input().Get("input_producer_mobile")
	input_producer_email := c.Input().Get("input_producer_email")
	input_cluster_name := c.Input().Get("input_cluster_name")
	input_product_name := c.Input().Get("input_product_name")
	select_channel_partner_type := c.Input().Get("select_channel_partner_type")
	// input_producer_status := c.Input().Get("input_producer_status")
	input_channelpartner_name := c.Input().Get("input_channelpartner_name")
	input_cp_mobile := c.Input().Get("input_cp_mobile")
	input_cp_email := c.Input().Get("input_cp_email")
	input_cp_status := c.Input().Get("input_cp_status")

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_producer_name - ", input_producer_name)

	c.Data["DateRange"] = c.Input().Get("daterange")
	tmp := strings.Split(c.Input().Get("daterange"), " - ")

	from := ""
	to := ""
	if len(tmp) == 2 {
		from = tmp[0]
		to = tmp[1]
	}

	row, err := db.Db.Query(`SELECT 
	id,
    created_on,

    JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.Name')) AS producer_name,
    JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.Mobile')) AS producer_mobile,
    JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.Email')) AS producer_email,
    JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.Add1')) AS producer_add1,
    JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.Add2')) AS producer_add2,
    JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.City')) AS producer_city,
    JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.State')) AS producer_state,
    JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.Zipcode')) AS producer_zipcode,
    JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.Status')) AS producer_status,
    JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.ClusterName')) AS producer_cluster,
    JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.LocationInfo')) AS producer_locationinfo,

    JSON_UNQUOTE(JSON_EXTRACT(channel_partner_details, '$.Name')) AS partner_name,
    JSON_UNQUOTE(JSON_EXTRACT(channel_partner_details, '$.Mobile')) AS partner_mobile,
    JSON_UNQUOTE(JSON_EXTRACT(channel_partner_details, '$.Email')) AS partner_email,
    JSON_UNQUOTE(JSON_EXTRACT(channel_partner_details, '$.Address1')) AS partner_address1,
    JSON_UNQUOTE(JSON_EXTRACT(channel_partner_details, '$.Address2')) AS partner_address2,
    JSON_UNQUOTE(JSON_EXTRACT(channel_partner_details, '$.City')) AS partner_city,
    JSON_UNQUOTE(JSON_EXTRACT(channel_partner_details, '$.State')) AS partner_state,
    JSON_UNQUOTE(JSON_EXTRACT(channel_partner_details, '$.Pincode')) AS partner_pincode,
    JSON_UNQUOTE(JSON_EXTRACT(channel_partner_details, '$.Status')) AS partner_status,
    JSON_UNQUOTE(JSON_EXTRACT(channel_partner_details, '$.ChannelPartnerType')) AS partner_type,
    JSON_UNQUOTE(JSON_EXTRACT(product_details, '$.ProductName')) AS product_name,
    JSON_UNQUOTE(JSON_EXTRACT(product_details, '$.ProductSKU')) AS product_sku,
    JSON_UNQUOTE(JSON_EXTRACT(product_details, '$.ProductQty')) AS product_qty,
    cluster
FROM Surplus_tracking
WHERE
    (JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.Name')) = '' 
     OR JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.Name')) LIKE ?) 
AND (JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.Mobile')) = '' 
     OR JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.Mobile')) LIKE ?) 
AND (JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.Email')) = '' 
     OR JSON_UNQUOTE(JSON_EXTRACT(producer_details, '$.Email')) LIKE ?) 
AND (cluster ='' OR cluster LIKE ?) 
AND (JSON_UNQUOTE(JSON_EXTRACT(product_details, '$.ProductName')) = '' 
     OR JSON_UNQUOTE(JSON_EXTRACT(product_details, '$.ProductName')) LIKE ?)
AND (JSON_UNQUOTE(JSON_EXTRACT(channel_partner_details, '$.ChannelPartnerType')) = '' 
     OR JSON_UNQUOTE(JSON_EXTRACT(channel_partner_details, '$.ChannelPartnerType')) LIKE ?) 
AND (JSON_UNQUOTE(JSON_EXTRACT(channel_partner_details, '$.Name')) = '' 
     OR JSON_UNQUOTE(JSON_EXTRACT(channel_partner_details, '$.Name')) LIKE ?)
AND (JSON_UNQUOTE(JSON_EXTRACT(channel_partner_details, '$.Mobile')) = '' 
     OR JSON_UNQUOTE(JSON_EXTRACT(channel_partner_details, '$.Mobile')) LIKE ?)
AND (JSON_UNQUOTE(JSON_EXTRACT(channel_partner_details, '$.Email')) = '' 
     OR JSON_UNQUOTE(JSON_EXTRACT(channel_partner_details, '$.Email')) LIKE ?)
AND (JSON_UNQUOTE(JSON_EXTRACT(channel_partner_details, '$.Status')) = '' 
     OR JSON_UNQUOTE(JSON_EXTRACT(channel_partner_details, '$.Status')) LIKE ?)
AND (DATE(created_on='') OR DATE (created_on) >= STR_TO_DATE(?, '%m/%d/%Y')) AND
    (DATE(created_on='') OR DATE (created_on) <= STR_TO_DATE(?, '%m/%d/%Y'))`,
		input_producer_name+"%", input_producer_mobile+"%", input_producer_email+"%",
		input_cluster_name+"%", input_product_name+"%", select_channel_partner_type+"%",
		input_channelpartner_name+"%", input_cp_mobile+"%",
		input_cp_email+"%", input_cp_status+"%", from, to)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch order")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch order")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug post", data)
	var result []Row

	for i := range data {
		var r Row
		r.Id = data[i][0]
		r.Timestamp = data[i][1]
		r.ProducerName = data[i][2]
		r.ProducerMobile = data[i][3]
		r.ProducerEmail = data[i][4]
		r.ProducerAdd1 = data[i][5]
		r.ProducerAdd2 = data[i][6]
		r.ProducerCity = data[i][7]
		r.ProducerState = data[i][8]
		r.ProducerZipcode = data[i][9]
		r.ProducerStatus = data[i][10]
		r.ProducerClusterName = data[i][11]
		r.LocationInfo = data[i][12]
		r.ChannelPartnerName = data[i][13]
		r.ChannelPartnerMobile = data[i][14]
		r.ChannelPartnerEmail = data[i][15]
		r.Address1 = data[i][16]
		r.Address2 = data[i][17]
		r.City = data[i][18]
		r.State = data[i][19]
		r.Pincode = data[i][20]
		r.ChannelPartnerStatus = data[i][21]
		r.ChannelPartnerType = data[i][22]
		r.ProductName = data[i][23]
		r.ProductSku = data[i][24]
		r.ProductQty = data[i][25]
		r.ClusterName = data[i][26]

		r.ProducerAddress = data[i][5] + ", " + data[i][6] + ", " + data[i][7] + ", " + data[i][8] + ", " + data[i][9]

		r.ChannelPartnerAddress = data[i][16] + ", " + data[i][17] + ", " + data[i][18] + ", " + data[i][19] + ", " + data[i][20]

		result = append(result, r)
	}
	c.Data["CustomerData"] = result

	return

}
