package searchChannelPartnerCluster

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
	Id          string
	Timestamp   string
	ClusterName string
	UserEmail   string
	Status      string
	Mobile      string
}

type SearchChannelPartnerCluster struct {
	// beego.Controller
	middleware.BaseController
}

func (c *SearchChannelPartnerCluster) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Channel Partner Cluster Page Start")
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
			c.TplName = "channelpartnercluster/searchChannelPartnerCluster/searchChannelPartnerCluster.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Channel Partner Cluster Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "channelpartnercluster/searchChannelPartnerCluster/searchChannelPartnerCluster.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Channel Partner Cluster  Page Success")
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

	auth, err := utils.IsAuthorized(utype, "ChannelPartnerCluster")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	row, err := db.Db.Query(`SELECT ChannelPartner2Cluster.id,ChannelPartner2Cluster.cluster_id, ChannelPartner2Cluster.channelpartner_id, 
       ChannelPartner2Cluster.status, ChannelPartner2Cluster.create_date,Clusters.cluster_name, 
       ChannelPartner.email,
    ChannelPartner.mob_no  
FROM ChannelPartner2Cluster   
LEFT JOIN Clusters 
ON ChannelPartner2Cluster.cluster_id = Clusters.id
LEFT JOIN ChannelPartner 
ON ChannelPartner2Cluster.channelpartner_id = ChannelPartner.id ORDER BY ChannelPartner2Cluster.create_date DESC`)

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
		r.ClusterName = data[i][5]
		r.UserEmail = data[i][6]
		r.Status = data[i][3]
		r.Timestamp = data[i][4]
		r.Mobile = data[i][7]

		result = append(result, r)
	}

	c.Data["CustomerData"] = result

	return
}

func (c *SearchChannelPartnerCluster) Post() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Channel Partner Cluster page")
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
			c.TplName = "channelpartnercluster/searchChannelPartnerCluster/searchChannelPartnerCluster.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Channel Partner Cluster Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "channelpartnercluster/searchChannelPartnerCluster/searchChannelPartnerCluster.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Channel Partner Cluster Page Success")
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

	auth, err := utils.IsAuthorized(utype, "ChannelPartnerCluster")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	//input_name := c.Input().Get("input_name")
	input_channelpartneremail := c.Input().Get("input_channelpartneremail")
	input_Clusterscluster_name := c.Input().Get("input_Clusterscluster_name")
	input_status := c.Input().Get("input_status")
	input_mob := c.Input().Get("input_mob")

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_Clusterscluster_name - ", input_Clusterscluster_name)
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

	row, err := db.Db.Query(`SELECT ChannelPartner2Cluster.id,ChannelPartner2Cluster.cluster_id, ChannelPartner2Cluster.channelpartner_id, 
       ChannelPartner2Cluster.status, ChannelPartner2Cluster.create_date,Clusters.cluster_name, 
       ChannelPartner.email,
    ChannelPartner.mob_no  
FROM ChannelPartner2Cluster   
LEFT JOIN Clusters 
ON ChannelPartner2Cluster.cluster_id = Clusters.id
LEFT JOIN ChannelPartner 
ON ChannelPartner2Cluster.channelpartner_id = ChannelPartner.id  where (Clusters.cluster_name='' OR Clusters.cluster_name like ?)AND (ChannelPartner.email ='' OR ChannelPartner.email  like ?)  AND (ChannelPartner2Cluster.status ='' OR ChannelPartner2Cluster.status  like ?) AND (ChannelPartner.mob_no ='' OR ChannelPartner.mob_no  like ?)  AND (DATE(ChannelPartner2Cluster.create_date='') OR DATE (ChannelPartner2Cluster.create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(ChannelPartner2Cluster.create_date='') OR DATE (ChannelPartner2Cluster.create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) ORDER BY ChannelPartner2Cluster.create_date DESC`, input_Clusterscluster_name+"%", input_channelpartneremail+"%", input_status+"%", input_mob+"%", from, to)

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
		r.ClusterName = data[i][5]
		r.UserEmail = data[i][6]
		r.Status = data[i][3]
		r.Timestamp = data[i][4]
		r.Mobile = data[i][7]

		result = append(result, r)
	}

	c.Data["CustomerData"] = result

	return
}
