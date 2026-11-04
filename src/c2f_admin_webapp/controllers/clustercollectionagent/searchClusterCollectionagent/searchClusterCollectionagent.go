package searchClusterCollectionagent

import (
	"c2f_admin_webapp/session"

	"runtime/debug"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"
	"errors"

	"c2f_admin_webapp/model/db"

	"c2f_admin_webapp/services"
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

type SearchClusterCollectionagent struct {
	// beego.Controller
	middleware.BaseController
}

func (c *SearchClusterCollectionagent) Get() {
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
			c.TplName = "clustercollectionagent/searchClusterCollectionagent/searchClusterCollectionagent.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "clustercollectionagent/searchClusterCollectionagent/searchClusterCollectionagent.html"
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

	auth, err := utils.IsAuthorized(utype, "ClusterCollectionagent")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	if utype == "CLUSTERMANAGER" {

		data1, err1 := services.GetClusterAsPerMapping(id)
		if err1 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
			err1 = errors.New("Cluster user fetch Failed")
			return
		}
		if len(data1) <= 0 {
			err1 = errors.New("Cluster user  Not Found")
			return
		}

		row, err := db.Db.Query(`SELECT Cluster2CollectionAgent.id,Cluster2CollectionAgent.cluster_id, Cluster2CollectionAgent.agent_id, 
       Cluster2CollectionAgent.status, Cluster2CollectionAgent.create_date,Clusters.cluster_name, 
       Collection_Agent.email,
    Collection_Agent.mob_no  
FROM Cluster2CollectionAgent   
LEFT JOIN Clusters 
ON Cluster2CollectionAgent.cluster_id = Clusters.id
LEFT JOIN Collection_Agent 
ON Cluster2CollectionAgent.agent_id = Collection_Agent.id where Cluster2CollectionAgent.cluster_id=? ORDER BY Cluster2CollectionAgent.create_date DESC`, data1[0][1])

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
	} else if utype == "CLUSTEREXECUTIVE" {

		data2, err2 := services.GetClusterAsPerMappingExec(id)
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
       Collection_Agent.email,
    Collection_Agent.mob_no  
FROM Cluster2CollectionAgent   
LEFT JOIN Clusters 
ON Cluster2CollectionAgent.cluster_id = Clusters.id
LEFT JOIN Collection_Agent 
ON Cluster2CollectionAgent.agent_id = Collection_Agent.id where Clusters.cluster_name=? ORDER BY Cluster2CollectionAgent.create_date DESC`, data2[0][3])

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
	} else {
		row, err := db.Db.Query(`SELECT Cluster2CollectionAgent.id,Cluster2CollectionAgent.cluster_id, Cluster2CollectionAgent.agent_id, 
       Cluster2CollectionAgent.status, Cluster2CollectionAgent.create_date,Clusters.cluster_name, 
       Collection_Agent.email ,
    Collection_Agent.mob_no 
FROM Cluster2CollectionAgent   
LEFT JOIN Clusters 
ON Cluster2CollectionAgent.cluster_id = Clusters.id
LEFT JOIN Collection_Agent 
ON Cluster2CollectionAgent.agent_id = Collection_Agent.id ORDER BY Cluster2CollectionAgent.create_date DESC`)

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
	}

	return
}

func (c *SearchClusterCollectionagent) Post() {
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
			c.TplName = "clustercollectionagent/searchClusterCollectionagent/searchClusterCollectionagent.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "clustercollectionagent/searchClusterCollectionagent/searchClusterCollectionagent.html"
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

	auth, err := utils.IsAuthorized(utype, "ClusterCollectionagent")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	//input_name := c.Input().Get("input_name")
	input_Useremail := c.Input().Get("input_Useremail")
	input_Clusterscluster_name := c.Input().Get("input_Clusterscluster_name")
	input_status := c.Input().Get("input_status")
	input_mob := c.Input().Get("input_mob")

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_Clusterscluster_name - ", input_Clusterscluster_name)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_Useremail - ", input_Useremail)
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

	// row, err := db.Db.Query(`SELECT id,first_name,email,status FROM sysuser.credential`)

	// row, err := db.Db.Query(`select first_name,mob_no,email,status,user_type from Sysuser_info where
	// DATE(create_date) >= STR_TO_DATE(?, '%m/%d/%Y') AND DATE(create_date) <= STR_TO_DATE(?, '%m/%d/%Y') ORDER BY create_date DESC`, from, to)

	if utype == "CLUSTERMANAGER" {

		data1, err1 := services.GetClusterAsPerMapping(id)
		if err1 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
			err1 = errors.New("Cluster user fetch Failed")
			return
		}
		if len(data1) <= 0 {
			err1 = errors.New("Cluster user  Not Found")
			return
		}

		row, err := db.Db.Query(`SELECT Cluster2CollectionAgent.id,Cluster2CollectionAgent.cluster_id, Cluster2CollectionAgent.agent_id, 
       Cluster2CollectionAgent.status, Cluster2CollectionAgent.create_date,Clusters.cluster_name, 
       Collection_Agent.email
    Collection_Agent.mob_no  
FROM Cluster2CollectionAgent   
LEFT JOIN Clusters 
ON Cluster2CollectionAgent.cluster_id = Clusters.id
LEFT JOIN Collection_Agent 
ON Cluster2CollectionAgent.agent_id = Collection_Agent.id  where Cluster2CollectionAgent.cluster_id=? AND (Clusters.cluster_name='' OR Clusters.cluster_name like ?)AND (Collection_Agent.email ='' OR Collection_Agent.email  like ?)  AND (Cluster2CollectionAgent.status ='' OR Cluster2CollectionAgent.status  like ?) AND (Collection_Agent.mob_no ='' OR Collection_Agent.mob_no  like ?)  AND (DATE(Cluster2CollectionAgent.create_date='') OR DATE (Cluster2CollectionAgent.create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(Cluster2CollectionAgent.create_date='') OR DATE (Cluster2CollectionAgent.create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) ORDER BY Cluster2CollectionAgent.create_date DESC`, data1[0][1], input_Clusterscluster_name+"%", input_Useremail+"%", input_status+"%", input_mob+"%", from, to)

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
	} else if utype == "CLUSTEREXECUTIVE" {

		data2, err2 := services.GetClusterAsPerMappingExec(id)
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
       Collection_Agent.email,
    Collection_Agent.mob_no  
FROM Cluster2CollectionAgent   
LEFT JOIN Clusters 
ON Cluster2CollectionAgent.cluster_id = Clusters.id
LEFT JOIN Collection_Agent 
ON Cluster2CollectionAgent.agent_id = Collection_Agent.id  where Clusters.cluster_name=? And (Clusters.cluster_name='' OR Clusters.cluster_name like ?)AND (Collection_Agent.email ='' OR Collection_Agent.email  like ?)  AND (Cluster2CollectionAgent.status ='' OR Cluster2CollectionAgent.status  like ?) AND (Collection_Agent.mob_no ='' OR Collection_Agent.mob_no  like ?)  AND (DATE(Cluster2CollectionAgent.create_date='') OR DATE (Cluster2CollectionAgent.create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(Cluster2CollectionAgent.create_date='') OR DATE (Cluster2CollectionAgent.create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) ORDER BY Cluster2CollectionAgent.create_date DESC`, data2[0][3], input_Clusterscluster_name+"%", input_Useremail+"%", input_status+"%", input_mob+"%", from, to)

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
	} else {
		row, err := db.Db.Query(`SELECT Cluster2CollectionAgent.id,Cluster2CollectionAgent.cluster_id, Cluster2CollectionAgent.agent_id, 
       Cluster2CollectionAgent.status, Cluster2CollectionAgent.create_date,Clusters.cluster_name, 
       Collection_Agent.email,
    Collection_Agent.mob_no  
FROM Cluster2CollectionAgent   
LEFT JOIN Clusters 
ON Cluster2CollectionAgent.cluster_id = Clusters.id
LEFT JOIN Collection_Agent 
ON Cluster2CollectionAgent.agent_id = Collection_Agent.id  where (Clusters.cluster_name='' OR Clusters.cluster_name like ?)AND (Collection_Agent.email ='' OR Collection_Agent.email  like ?)  AND (Cluster2CollectionAgent.status ='' OR Cluster2CollectionAgent.status  like ?) AND (Collection_Agent.mob_no ='' OR Collection_Agent.mob_no  like ?)  AND (DATE(Cluster2CollectionAgent.create_date='') OR DATE (Cluster2CollectionAgent.create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(Cluster2CollectionAgent.create_date='') OR DATE (Cluster2CollectionAgent.create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) ORDER BY Cluster2CollectionAgent.create_date DESC`, input_Clusterscluster_name+"%", input_Useremail+"%", input_status+"%", input_mob+"%", from, to)

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
	}

	return
}
