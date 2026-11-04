package createClusterClusterexecutive

import (
	"c2f_admin_webapp/session"

	"runtime/debug"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/model/utils"
	"c2f_admin_webapp/services"
	"errors"
	"strings"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

type CreateClusterClusterexecutive struct {
	// beego.Controller
	middleware.BaseController
}

type Display struct {
	Fields []Field
}

type Display1 struct {
	Fields1 []Field1
}
type Field struct {
	Id    string
	Name  string
	Email string
}

type Field1 struct {
	Id        string
	UserEmail string
	Name      string
}

func (c *CreateClusterClusterexecutive) Get() {
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
			c.TplName = "clusterclusterexecutive/createClusterClusterexecutive/createClusterClusterexecutive.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "clusterclusterexecutive/createClusterClusterexecutive/createClusterClusterexecutive.html"
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

	id := sess.Get("uid").(string)
	c.Data["id"] = id

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

	if utype == "CLUSTERMANAGER" {

		data, err := services.GetClusterAsPerMapping(id)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster user fetch Failed")
			return
		}
		if len(data) <= 0 {
			err = errors.New("Cluster user  Not Found")
			return
		}
		var Dis Display
		for i := range data {
			var d Field
			d.Id = data[i][2]
			d.Name = data[i][3]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

	} else if utype == "CLUSTEREXECUTIVE" {

		data, err := services.GetClusterAsPerMappingExec(id)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster user fetch Failed")
			return
		}
		if len(data) <= 0 {
			err = errors.New("Cluster user  Not Found")
			return
		}
		var Dis Display
		for i := range data {
			var d Field
			d.Id = data[i][2]
			d.Name = data[i][3]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

	} else {

		data, err := services.GetClusteruser()
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster user fetch Failed")
			return
		}
		if len(data) <= 0 {
			err = errors.New("Cluster user  Not Found")
			return
		}
		var Dis Display
		for i := range data {
			var d Field
			d.Id = data[i][0]
			d.Name = data[i][1]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

	}

	data1, err := services.GetNonmappedClusterExecutive()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Non-mapped ClusterExecutive fetch Failed")
		return
	}
	if len(data1) <= 0 {
		err = errors.New("Non-mapped ClusterExecutive  Not Found")
		return
	}
	var Dis1 Display1
	for i := range data1 {
		var d Field1
		d.Id = data1[i][0]
		d.UserEmail = data1[i][1]
		d.Name = data1[i][2]
		Dis1.Fields1 = append(Dis1.Fields1, d)
	}
	c.Data["Dis1"] = Dis1
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis1)

	auth, err := utils.IsAuthorized(utype, "ClusterClusterexecutive")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	return
}
func (c *CreateClusterClusterexecutive) Post() {

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
			c.TplName = "clusterclusterexecutive/createClusterClusterexecutive/createClusterClusterexecutive.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Cluster ClusterExecutive Mapped Successfully"
			c.TplName = "clusterclusterexecutive/createClusterClusterexecutive/createClusterClusterexecutive.html"
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

	id := sess.Get("uid").(string)
	c.Data["id"] = id
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "id", id)

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

	auth, err := utils.IsAuthorized(utype, "ClusterClusterexecutive")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	select_cluster_name := c.Input().Get("select_cluster_name")
	select_clusterexecutive_name := c.Input().Get("select_clusterexecutive_name")

	input_status := c.Input().Get("input_status")

	err = CheckExists(select_clusterexecutive_name)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Cluster Executive already mapped")
		return
	}

	//err = CheckmapAlreadyExists(select_cluster_name, select_clusterexecutive_name)

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cluster_id  - ", select_cluster_name)

	if select_cluster_name == "" || select_clusterexecutive_name == "" {
		err = errors.New("cluster name and user name cannot be empty")
		return
	}
	if input_status == "" {
		err = errors.New("status cannot be empty")
		return
	}
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Cluster 2 User mapping already Exists")
	// 	return
	// }

	result, err := db.Db.Exec(`INSERT INTO User2Cluster (cluster_id,
		clusterexecutive_id,
		status,
		created_by,
		create_date)
		VALUES (?,?,?,?,now())`,
		select_cluster_name,
		select_clusterexecutive_name,
		input_status,
		id)
	if err != nil {
		err = errors.New("Cluster User creation failed")
		return
	}

	i, err := result.RowsAffected()
	if err != nil || i == 0 {
		err = errors.New("ClusterProduct creation failed")
		return
	}

	return
}

func CheckExists(clusterexecutive string) (err error) {

	err = nil

	row, err := db.Db.Query(`select count(*) from User2Cluster WHERE clusterexecutive_id = ?`, clusterexecutive)

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
		err = errors.New("user already mapped")
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}

	return

}

func CheckmapAlreadyExists(cluster, clusterexecutive string) (err error) {

	err = nil

	row, err := db.Db.Query(`select count(*) from User2Cluster WHERE cluster_id = ? AND clusterexecutive_id=?`, cluster, clusterexecutive)

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
		err = errors.New("Mappig already exists")
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}

	return

}
