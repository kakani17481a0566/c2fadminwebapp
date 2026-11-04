package updateChannelPartnerCluster

import (
	//	"c2f_admin_webapp/services"
	"c2f_admin_webapp/session"

	"runtime/debug"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"
	"errors"
	"strings"

	"c2f_admin_webapp/model/db"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

type Row struct {
	Id                  string
	ClusterName         string
	ChannelPartnerEmail string
	Status              string
}

// type Display struct {
// 	Fields []Field
// }

// type Display1 struct {
// 	Fields1 []Field1
// }
// type Field struct {
// 	Id    string
// 	Name  string
// 	Email string
// }

// type Field1 struct {
// 	Id   string
// 	Name string
// }

type UpdateChannelPartnerCluster struct {
	// beego.Controller
	middleware.BaseController
}

func (c *UpdateChannelPartnerCluster) Get() {
	AdminId := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Channel Partner Cluster Page Start")
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
			c.TplName = "channelpartnercluster/updateChannelPartnerCluster/updateChannelPartnerCluster.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Channel Partner Cluster Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "channelpartnercluster/updateChannelPartnerCluster/updateChannelPartnerCluster.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Channel Partner Cluster Page Success")
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
ON ChannelPartner2Cluster.channelpartner_id = ChannelPartner.id
 WHERE ChannelPartner2Cluster.id =?  limit 1`, AdminId)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("User data not found")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)
	for i := range data {
		c.Data["Id"] = data[i][0]
		c.Data["ClusterName"] = data[i][5]
		c.Data["ChannelPartnerEmail"] = data[i][6]
		c.Data["Status"] = data[i][3]
	}

	return
}

func (c *UpdateChannelPartnerCluster) Post() {

	AdminId := c.Ctx.Input.Param(":AdminID")

	log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Channel Partner Cluster page")
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
			c.TplName = "channelpartnercluster/updateChannelPartnerCluster/updateChannelPartnerCluster.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Channel Partner Cluster Page Fail")
		} else {
			c.Data["DisplayMessage"] = "ChannelPartner-Cluster Updated Successfully"
			c.TplName = "channelpartnercluster/updateChannelPartnerCluster/updateChannelPartnerCluster.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Channel Partner Cluster Page Success")
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

	input_status := c.Input().Get("input_status")

	res, err := db.Db.Exec(`UPDATE ChannelPartner2Cluster SET 
	status = ?
	WHERE id = ?`, input_status, AdminId)
	if err != nil {
		err = errors.New("ChannelPartner-Cluster updation failed")
		return
	}

	i, err := res.RowsAffected()
	if err != nil || i == 0 {
		err = errors.New("ChannelPartner-Cluster updation failed")
		return
	}

	row, err := db.Db.Query(`SELECT ChannelPartner2Cluster.id,ChannelPartner2Cluster.cluster_id, ChannelPartner2Cluster.channelpartner_id, 
       ChannelPartner2Cluster.status, ChannelPartner2Cluster.create_date,Clusters.cluster_name, 
       ChannelPartner.email,
    ChannelPartner.mob_no  
FROM ChannelPartner2Cluster   
LEFT JOIN Clusters 
ON ChannelPartner2Cluster.cluster_id = Clusters.id
LEFT JOIN ChannelPartner 
ON ChannelPartner2Cluster.channelpartner_id = ChannelPartner.id
 WHERE ChannelPartner2Cluster.id =?  limit 1`, AdminId)

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
		c.Data["Id"] = data[i][0]
		c.Data["ClusterName"] = data[i][5]
		c.Data["ChannelPartnerEmail"] = data[i][6]
		c.Data["Status"] = data[i][3]

		result = append(result, r)
	}

	return

}
