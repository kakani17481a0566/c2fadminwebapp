package searchcollectionAgent

import (
	"c2f_admin_webapp/session"

	"runtime/debug"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"
	"errors"

	"c2f_admin_webapp/services"

	"c2f_admin_webapp/model/db"

	"strings"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

type Row struct {
	Id              string
	FirstName       string
	MiddleName      string
	LastName        string
	Mobile          string
	Email           string
	UserType        string
	Status          string
	CollectionAgent string
	Timestamp       string
}

type SearchcollectionAgent struct {
	// beego.Controller
	middleware.BaseController
}

func (c *SearchcollectionAgent) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Start")
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
			c.TplName = "collectionagent/searchcollectionAgent/searchcollectionAgent.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "collectionagent/searchcollectionAgent/searchcollectionAgent.html"
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

	auth, err := utils.IsAuthorized(utype, "CollectionAgent")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	if utype == "CLUSTERMANAGER" {

		data2, err2 := services.GetClusterAsPerMapping(id)
		if err2 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err2)
			err2 = errors.New("Cluster user fetch Failed")
			return
		}
		if len(data2) <= 0 {
			err2 = errors.New("Cluster user  Not Found")
			return
		}

		row, err := db.Db.Query(`SELECT Cluster2CollectionAgent.id,Cluster2CollectionAgent.cluster_id, Cluster2CollectionAgent.agent_id, 
       Cluster2CollectionAgent.status, Cluster2CollectionAgent.create_date,Clusters.cluster_name, 
       Collection_Agent.email,Collection_Agent.id,Collection_Agent.first_name,Collection_Agent.mob_no,Collection_Agent.status,Collection_Agent.create_date 
FROM Cluster2CollectionAgent   
LEFT JOIN Clusters 
ON Cluster2CollectionAgent.cluster_id = Clusters.id
LEFT JOIN Collection_Agent 
ON Cluster2CollectionAgent.agent_id = Collection_Agent.id where Cluster2CollectionAgent.cluster_id=? ORDER BY Collection_Agent.create_date DESC`, data2[0][1])

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
			r.Id = data[i][7]
			r.FirstName = data[i][8]
			r.Mobile = data[i][9]
			r.Email = data[i][6]
			r.Status = data[i][10]
			//r.CollectionAgent = data[i][5]
			r.Timestamp = data[i][11]
			result = append(result, r)

		}
		c.Data["CustomerData"] = result
	} else if utype == "CLUSTEREXECUTIVE" {

		//NO CODE FOR CLUSTEREXECUTIVE
	} else {

		row, err := db.Db.Query(`select id,first_name,mob_no,email,status,create_date from Collection_Agent ORDER BY create_date DESC`)

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
			//r.CollectionAgent = data[i][5]
			r.Timestamp = data[i][5]
			result = append(result, r)

		}
		c.Data["CustomerData"] = result

	}

	return
}

func (c *SearchcollectionAgent) Post() {
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
			c.TplName = "collectionagent/searchcollectionAgent/searchcollectionAgent.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "collectionagent/searchcollectionAgent/searchcollectionAgent.html"
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

	auth, err := utils.IsAuthorized(utype, "CollectionAgent")
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

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_name - ", input_name)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_email - ", input_email)
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

	if utype == "CLUSTERMANAGER" {

		data2, err2 := services.GetClusterAsPerMapping(id)
		if err2 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err2)
			err2 = errors.New("Cluster user fetch Failed")
			return
		}
		if len(data2) <= 0 {
			err2 = errors.New("Cluster user  Not Found")
			return
		}

		row, err := db.Db.Query(`SELECT Cluster2CollectionAgent.id,Cluster2CollectionAgent.cluster_id, Cluster2CollectionAgent.agent_id, 
       Cluster2CollectionAgent.status, Cluster2CollectionAgent.create_date,Clusters.cluster_name, 
       Collection_Agent.email,Collection_Agent.id,Collection_Agent.first_name,Collection_Agent.mob_no,Collection_Agent.status,Collection_Agent.create_date 
FROM Cluster2CollectionAgent   
LEFT JOIN Clusters 
ON Cluster2CollectionAgent.cluster_id = Clusters.id
LEFT JOIN Collection_Agent 
ON Cluster2CollectionAgent.agent_id = Collection_Agent.id where Cluster2CollectionAgent.cluster_id=? AND (Collection_Agent.first_name='' OR Collection_Agent.first_name like ?)AND (Collection_Agent.status='' OR Collection_Agent.status like ?)AND (Collection_Agent.email='' OR Collection_Agent.email like ?) AND (Collection_Agent.mob_no='' OR Collection_Agent.mob_no like ?) AND (DATE(Collection_Agent.create_date='') OR DATE (Collection_Agent.create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(Collection_Agent.create_date='') OR DATE (Collection_Agent.create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) ORDER BY Collection_Agent.create_date DESC`, data2[0][1], input_name+"%", input_status+"%", input_email+"%", input_mob+"%", from, to)

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
			r.Id = data[i][7]
			r.FirstName = data[i][8]
			r.Mobile = data[i][9]
			r.Email = data[i][6]
			r.Status = data[i][10]
			//r.CollectionAgent = data[i][5]
			r.Timestamp = data[i][11]
			result = append(result, r)

		}
		c.Data["CustomerData"] = result
	} else if utype == "CLUSTEREXECUTIVE" {

		//NO CODE FOR CLUSTEREXECUTIVE
	} else {

		row, err := db.Db.Query(`select id,first_name,mob_no,email,status,create_date from Collection_Agent where (first_name='' OR first_name like ?)AND (status='' OR status like ?)AND (email='' OR email like ?) AND (mob_no='' OR mob_no like ?) AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y'))`, input_name+"%", input_status+"%", input_email+"%", input_mob+"%", from, to)

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
			//r.CollectionAgent = data[i][5]
			r.Timestamp = data[i][5]
			result = append(result, r)

		}
		c.Data["CustomerData"] = result

	}

	// row, err := db.Db.Query(`SELECT id,first_name,email,status FROM sysuser.credential`)

	// row, err := db.Db.Query(`select first_name,mob_no,email,status,user_type from Sysuser_info where
	// DATE(create_date) >= STR_TO_DATE(?, '%m/%d/%Y') AND DATE(create_date) <= STR_TO_DATE(?, '%m/%d/%Y') ORDER BY create_date DESC`, from, to)

	//row, err := db.Db.Query(`select id,first_name,mob_no,email,status,create_date from Collection_Agent where (first_name='' OR first_name like ?)AND (email='' OR email like ?) AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y'))`, input_name+"%", input_email+"%", from, to)

	//where (Clusters.cluster_name='' OR Clusters.cluster_name like ?)  AND (User.email ='' OR User.email  like ?)  AND (DATE(User2Cluster.create_date='') OR DATE (User2Cluster.create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(User2Cluster.create_date='') OR DATE (User2Cluster.create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) ORDER BY User2Cluster.create_date DESC`, input_Clusterscluster_name+"%", input_Useremail+"%", from, to)

	return
}
