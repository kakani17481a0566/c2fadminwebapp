package surplusmanagementReport

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
	Id                 string
	Timestamp          string
	ChannelPartnerName string
	Mobile             string
	Email              string
	ClusterName        string
	ProductName        string
	ChannelPartnerType string
	Status             string
	Address            string
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
type SurplusManagementReport struct {
	// beego.Controller
	middleware.BaseController
}

func (c *SurplusManagementReport) Get() {

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
			c.TplName = "reports/surplusmanagementReport/surplusmanagementReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Reports Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "reports/surplusmanagementReport/surplusmanagementReport.html"
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

	auth, err := utils.IsAuthorized(utype, "SplusMngReport")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	// HOADMIN CODE

	row, err := db.Db.Query(`
	SELECT 
    ChannelPartner.id,
    ChannelPartner.first_name,
    ChannelPartner.mob_no,
    ChannelPartner.email,
    ChannelPartner.status,
    ChannelPartner.create_date,
    ChannelPartnerType.name AS channel_partner_type,
    GROUP_CONCAT(DISTINCT Clusters.cluster_name ORDER BY Clusters.cluster_name SEPARATOR ', ') AS cluster_names,
    GROUP_CONCAT(DISTINCT Lookup_products.product_name ORDER BY Lookup_products.product_name SEPARATOR ', ') AS product_names,
    ChannelPartner.address1,
    ChannelPartner.address2,
    ChannelPartner.city,
    ChannelPartner.state,
    ChannelPartner.zipcode
FROM ChannelPartner
LEFT JOIN ChannelPartnerType 
    ON ChannelPartner.channel_partner_type_id = ChannelPartnerType.id
LEFT JOIN ChannelPartner2Cluster 
    ON ChannelPartner.id = ChannelPartner2Cluster.channelpartner_id
LEFT JOIN Clusters 
    ON ChannelPartner2Cluster.cluster_id = Clusters.id
LEFT JOIN ChannelPartner2Product 
    ON ChannelPartner.id = ChannelPartner2Product.channelpartner_id
LEFT JOIN Lookup_products 
    ON ChannelPartner2Product.product_id = Lookup_products.id
GROUP BY 
    ChannelPartner.id,
    ChannelPartner.first_name,
    ChannelPartner.mob_no,
    ChannelPartner.email,
    ChannelPartner.status,
    ChannelPartner.create_date,
    ChannelPartnerType.name,
    ChannelPartner.address1,
    ChannelPartner.address2,
    ChannelPartner.city,
    ChannelPartner.state,
    ChannelPartner.zipcode
ORDER BY ChannelPartner.id DESC`)

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
		r.ChannelPartnerName = data[i][1]
		r.Mobile = data[i][2]
		r.Email = data[i][3]
		r.Status = data[i][4]
		r.Timestamp = data[i][5]
		r.ChannelPartnerType = data[i][6]
		r.ClusterName = data[i][7]
		r.ProductName = data[i][8]
		r.Address = data[i][9] + ", " + data[i][10] + ", " + data[i][11] + ", " + data[i][12] + ", " + data[i][13]

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

func (c *SurplusManagementReport) Post() {
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
			c.TplName = "reports/surplusmanagementReport/surplusmanagementReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchOrder Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "reports/surplusmanagementReport/surplusmanagementReport.html"
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

	auth, err := utils.IsAuthorized(utype, "SplusMngReport")
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

	input_first_name := c.Input().Get("input_first_name")
	input_mobile := c.Input().Get("input_mobile")
	input_email := c.Input().Get("input_email")
	input_cluster_name := c.Input().Get("input_cluster_name")
	input_product_name := c.Input().Get("input_product_name")
	select_channel_partner_type := c.Input().Get("select_channel_partner_type")
	input_status := c.Input().Get("input_status")

	c.Data["DateRange"] = c.Input().Get("daterange")
	tmp := strings.Split(c.Input().Get("daterange"), " - ")

	from := ""
	to := ""
	if len(tmp) == 2 {
		from = tmp[0]
		to = tmp[1]
	}

	row, err := db.Db.Query(`
	SELECT 
    ChannelPartner.id,
    ChannelPartner.first_name,
    ChannelPartner.mob_no,
    ChannelPartner.email,
    ChannelPartner.status,
    ChannelPartner.create_date,
    ChannelPartnerType.name AS channel_partner_type,
    GROUP_CONCAT(DISTINCT Clusters.cluster_name ORDER BY Clusters.cluster_name SEPARATOR ', ') AS cluster_names,
    GROUP_CONCAT(DISTINCT Lookup_products.product_name ORDER BY Lookup_products.product_name SEPARATOR ', ') AS product_names,
    ChannelPartner.address1,
    ChannelPartner.address2,
    ChannelPartner.city,
    ChannelPartner.state,
    ChannelPartner.zipcode
FROM ChannelPartner
LEFT JOIN ChannelPartnerType 
    ON ChannelPartner.channel_partner_type_id = ChannelPartnerType.id
LEFT JOIN ChannelPartner2Cluster 
    ON ChannelPartner.id = ChannelPartner2Cluster.channelpartner_id
LEFT JOIN Clusters 
    ON ChannelPartner2Cluster.cluster_id = Clusters.id
LEFT JOIN ChannelPartner2Product 
    ON ChannelPartner.id = ChannelPartner2Product.channelpartner_id
LEFT JOIN Lookup_products 
    ON ChannelPartner2Product.product_id = Lookup_products.id
WHERE
    (ChannelPartner.first_name='' OR ChannelPartner.first_name LIKE ?) AND
    (ChannelPartner.mob_no='' OR ChannelPartner.mob_no LIKE ?) AND
    (ChannelPartner.email='' OR ChannelPartner.email LIKE ?) AND
    (Clusters.cluster_name='' OR Clusters.cluster_name LIKE ?) AND
    (Lookup_products.product_name='' OR Lookup_products.product_name LIKE ?) AND
    (ChannelPartnerType.name='' OR ChannelPartnerType.name LIKE ?) AND
    (ChannelPartner.status='' OR ChannelPartner.status LIKE ?) AND
    (DATE(ChannelPartner.create_date='') OR DATE (ChannelPartner.create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND
    (DATE(ChannelPartner.create_date='') OR DATE (ChannelPartner.create_date) <= STR_TO_DATE(?, '%m/%d/%Y'))
GROUP BY 
    ChannelPartner.id,
    ChannelPartner.first_name,
    ChannelPartner.mob_no,
    ChannelPartner.email,
    ChannelPartner.status,
    ChannelPartner.create_date,
    ChannelPartnerType.name,
    ChannelPartner.address1,
    ChannelPartner.address2,
    ChannelPartner.city,
    ChannelPartner.state,
    ChannelPartner.zipcode
ORDER BY ChannelPartner.id DESC`,
		input_first_name+"%", input_mobile+"%", input_email+"%", input_cluster_name+"%", input_product_name+"%", select_channel_partner_type+"%", input_status+"%", from, to)

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

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)
	var result []Row

	for i := range data {
		var r Row
		r.Id = data[i][0]
		r.ChannelPartnerName = data[i][1]
		r.Mobile = data[i][2]
		r.Email = data[i][3]
		r.Status = data[i][4]
		r.Timestamp = data[i][5]
		r.ChannelPartnerType = data[i][6]
		r.ClusterName = data[i][7]
		r.ProductName = data[i][8]
		r.Address = data[i][9] + ", " + data[i][10] + ", " + data[i][11] + ", " + data[i][12] + ", " + data[i][13]

		result = append(result, r)
		c.Data["CustomerData"] = result
	}

	return

}
