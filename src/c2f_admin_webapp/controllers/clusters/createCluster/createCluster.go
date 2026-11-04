package createCluster

import (
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

type User struct {
	Id string `json:"Id,omitempty"`
}
type CreateCluster struct {
	// beego.Controller
	middleware.BaseController
}

func (c *CreateCluster) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Add Assets Start")
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

			c.TplName = "clusters/createClusters/createClusters.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "clusters/createClusters/createClusters.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer  Page Success")
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

	maptoken := sess.Get("token")
	c.Data["maptoken"] = maptoken
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "maptoken", maptoken)

	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1

	auth, err := utils.IsAuthorized(utype, "Cluster")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	return
}
func (c *CreateCluster) Post() {

	log.Println(beego.AppConfig.String("loglevel"), "Info", "add asset post page")
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
			c.TplName = "clusters/createClusters/createClusters.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Cluster Added Successfully"
			c.TplName = "clusters/createClusters/createClusters.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer  Page Success")
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
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "utype", utype)

	language := sess.Get("language").(string)
	c.Data["language"] = language
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "language", language)

	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1

	id := sess.Get("uid").(string)
	c.Data["id"] = id

	auth, err := utils.IsAuthorized(utype, "Cluster")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	input_cluster_name := c.Input().Get("input_cluster_name")
	input_location := c.Input().Get("input_location")
	input_zipcode := c.Input().Get("input_zipcode")
	input_status := c.Input().Get("input_status")
	input_address := c.Input().Get("input_CommunicationAddress")
	input_cluster_type := c.Input().Get("input_cluster_type")
	input_clustersub_type := c.Input().Get("input_cluster_subtype")
	input_Identity := c.Input().Get("input_Identity")
	input_location_type := c.Input().Get("input_location_type")

	input_clusterpartner_type := c.Input().Get("input_clusterpartner_type")
	input_gstin := c.Input().Get("input_gstin")

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cluster_name  - ", input_cluster_name)

	if input_cluster_name == "" || input_location == "" || input_zipcode == "" || input_status == "" || input_address == "" || input_cluster_type == "" ||
		input_clustersub_type == "" || input_Identity == "" || input_clusterpartner_type == "" || input_clusterpartner_type == "" || input_gstin == "" {
		err = errors.New("* fields cannot be empty")
		return
	}
	// if input_clustersub_type == "" || input_Identity == "" {
	// 	err = errors.New("cluster subtype and cannot be empty")
	// 	return
	// }

	err = CheckExists(input_cluster_name)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Cluster already Exists")
		return
	}

	cluster_id := input_cluster_name + input_zipcode
	result, err := db.Db.Exec(`INSERT INTO Clusters (cluster_name,
		location,
		zipcode,
		status,
		address,
		cluster_type,
		cluster_subtype,
		identity,
		cluster_id,
		location_type,
		cluster_partner_type,
		GSTIN,
		created_by,
		create_date)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,now())`,
		input_cluster_name,
		input_location,
		input_zipcode,
		input_status,
		input_address,
		input_cluster_type,
		input_clustersub_type,
		input_Identity,
		cluster_id,
		input_location_type,
		input_clusterpartner_type,
		input_gstin,
		id)
	if err != nil {
		err = errors.New("Cluster creation failed")
		return
	}

	i, err := result.RowsAffected()
	if err != nil || i == 0 {
		err = errors.New("Cluster creation failed")
		return
	}

	return
}

func CheckExists(cluster_name string) (err error) {

	err = nil

	row, err := db.Db.Query(`select count(*) from Clusters WHERE cluster_name = ?`, cluster_name)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch user")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch user")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	countlen := data[0][0]

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "countlen", countlen)

	if countlen != "0" {
		err = errors.New("User already exists")
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}

	return

}
