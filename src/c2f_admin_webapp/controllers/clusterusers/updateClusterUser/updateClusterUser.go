package updateClusterUser

import (
	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/model/utils"
	"c2f_admin_webapp/services"
	"c2f_admin_webapp/session"
	"errors"
	"runtime/debug"
	"strings"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"

	"github.com/astaxie/beego"
)

type Row struct {
	Id          string
	ClusterName string
	UserEmail   string
	Status      string
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
	Id   string
	Name string
}
type UpdateClusterUser struct {
	// beego.Controller
	middleware.BaseController
}

func (c *UpdateClusterUser) Get() {
	AdminId := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User Start")
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
			c.TplName = "clusterusers/updateClusterUser/updateClusterUser.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Cluster User Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "clusterusers/updateClusterUser/updateClusterUser.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Cluster User Page Success")
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
	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1
	data2, err := services.GetClusteruser()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Cluster user fetch Failed")
		return
	}
	if len(data2) <= 0 {
		err = errors.New("Cluster user  Not Found")
		return
	}
	var Dis Display
	for i := range data2 {
		var d Field
		d.Id = data2[i][0]
		d.Name = data2[i][1]
		Dis.Fields = append(Dis.Fields, d)
	}
	c.Data["Dis"] = Dis
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)
	data1, err := services.Getusers()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("User fetch Failed")
		return
	}
	if len(data1) <= 0 {
		err = errors.New("User  Not Found")
		return
	}
	var Dis1 Display1
	for i := range data1 {
		var d Field1
		d.Id = data1[i][0]
		d.Name = data1[i][1]
		Dis1.Fields1 = append(Dis1.Fields1, d)
	}
	c.Data["Dis1"] = Dis1
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis1)
	auth, err := utils.IsAuthorized(utype, "ClusterUser")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")
	row, err := db.Db.Query(`SELECT  User2Cluster.id,User2Cluster.cluster_id, User2Cluster.user_id, 
       User2Cluster.status, User2Cluster.create_date,Clusters.cluster_name, 
       User.email  
FROM User2Cluster
LEFT JOIN Clusters 
ON User2Cluster.cluster_id = Clusters.id
LEFT JOIN User 
ON User2Cluster.user_id = User.id WHERE User2Cluster.id =?  limit 1`, AdminId)
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
		c.Data["UserEmail"] = data[i][6]
		c.Data["Status"] = data[i][3]
	}
	return
}
func (c *UpdateClusterUser) Post() {
	AdminId := c.Ctx.Input.Param(":AdminID")
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
			c.TplName = "clusterusers/updateClusterUser/updateClusterUser.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update ClusterUser Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Cluster user Updated Successfully"
			c.TplName = "clusterusers/updateClusterUser/updateClusterUser.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update ClusterUser   Page Success")
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
	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1
	auth, err := utils.IsAuthorized(utype, "ClusterUser")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")
	input_status := c.Input().Get("input_status")
	select_cluster_name := c.Input().Get("select_cluster_name")
	res, err := db.Db.Exec(`UPDATE User2Cluster SET 
        cluster_id=?,
        status = ?
        WHERE id = ? `, select_cluster_name, input_status, AdminId)
	if err != nil {
		err = errors.New("Cluster-user updation failed")
		return
	}
	i, err := res.RowsAffected()
	if err != nil || i == 0 {
		err = errors.New("Cluster-user updation failed")
		return
	}
	row, err := db.Db.Query(`SELECT  User2Cluster.id,User2Cluster.cluster_id, User2Cluster.user_id, 
       User2Cluster.status, User2Cluster.create_date,Clusters.cluster_name, 
       User.email  
FROM User2Cluster
LEFT JOIN Clusters 
ON User2Cluster.cluster_id = Clusters.id
LEFT JOIN User 
ON User2Cluster.user_id = User.id WHERE User2Cluster.id =?  limit 1`, AdminId)
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
		result = append(result, r)
	}
	return
}
