package searchSupplyExecutive

import (
	"c2f_admin_webapp/session"

	"runtime/debug"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"
	"c2f_admin_webapp/services"
	"errors"

	"c2f_admin_webapp/model/db"

	"strings"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"
	"c2f_admin_webapp/model/util/log"
)

type Row struct {
	Id         string
	FirstName  string
	MiddleName string
	LastName   string
	Mobile     string
	Email      string
	Usertype   string
	Status     string
	CreateDate string
}
type Display5 struct {
	Fields5 []Field5
}

type Field5 struct {
	Id     string
	Mobile string
}
type SearchSupplyExecutive struct {
	// beego.Controller
	middleware.BaseController
}

func (c *SearchSupplyExecutive) Get() {
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
			c.TplName = "supplyexecutive/searchSupplyExecutive/searchSupplyExecutive.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchSupplyExecutive Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "supplyexecutive/searchSupplyExecutive/searchSupplyExecutive.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchSupplyExecutive  Page Success")
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

	data, err := services.GetDeliveryuser()
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

	auth, err := utils.IsAuthorized(utype, "SupplyExecutive")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")
	if utype == "CLUSTERMANAGER" {

		//

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

		row, err := db.Db.Query(`SELECT 
SupplyExecutive.id,
SupplyExecutive.create_date,
SupplyExecutive.first_name,
SupplyExecutive.mob_no,
SupplyExecutive.email,
SupplyExecutive.status,
SupplyExecutive.user_type FROM SupplyExecutive
LEFT JOIN Cluster2SupplyExecutive 
ON SupplyExecutive.id = Cluster2SupplyExecutive.supplyexecutive_id
LEFT JOIN Clusters 
ON Cluster2SupplyExecutive.cluster_id = Clusters.id  where Cluster2SupplyExecutive.cluster_id =? ORDER BY SupplyExecutive.create_date DESC `, data1[0][1])

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get the error message info")
			return
		}
		defer sql.Close(row)
		_, data, err1 := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
			err1 = errors.New("Message Scan Fail")
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data for hoadmin before filter - ", data, " Data len - ", len(data))
		if len(data) <= 0 {
			err1 = errors.New("Unable to get the error message info")
			return
		}

		var result []Row

		for i := range data {
			var r Row
			r.Id = data[i][0]
			r.CreateDate = data[i][1]
			r.FirstName = data[i][2]
			r.Mobile = data[i][3]
			r.Email = data[i][4]
			r.Status = data[i][5]
			r.Usertype = data[i][6]
			result = append(result, r)
		}

		c.Data["CustomerData"] = result

	} else if utype == "CLUSTEREXECUTIVE" {

		//NO CODE FOR CLUSTEREXECUTIVE

	} else {

		row, err := db.Db.Query(`select id,create_date,first_name,mob_no,email,status,user_type from SupplyExecutive ORDER BY create_date DESC `)

		//AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y'))

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get the error message info")
			return
		}
		defer sql.Close(row)
		_, data, err1 := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
			err1 = errors.New("Message Scan Fail")
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data for hoadmin before filter - ", data, " Data len - ", len(data))
		if len(data) <= 0 {
			err1 = errors.New("Unable to get the error message info")
			return
		}

		var result []Row

		for i := range data {
			var r Row
			r.Id = data[i][0]
			r.CreateDate = data[i][1]
			r.FirstName = data[i][2]
			r.Mobile = data[i][3]
			r.Email = data[i][4]
			r.Status = data[i][5]
			r.Usertype = data[i][6]
			result = append(result, r)
		}

		c.Data["CustomerData"] = result

	}

	return
}

func (c *SearchSupplyExecutive) Post() {
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
			c.TplName = "supplyexecutive/searchSupplyExecutive/searchSupplyExecutive.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchSupplyExecutive Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "supplyexecutive/searchSupplyExecutive/searchSupplyExecutive.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchSupplyExecutive  Page Success")
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

	auth, err := utils.IsAuthorized(utype, "SupplyExecutive")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	input_name := c.Input().Get("input_name")
	input_email := c.Input().Get("input_email")
	input_first_name := c.Input().Get("input_first_name")
	input_mob_no := c.Input().Get("input_mob_no")
	input_status := c.Input().Get("input_status")

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_name - ", input_name)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_email - ", input_email)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_first_name - ", input_first_name)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_mob_no - ", input_mob_no)
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

		//

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

		row, err := db.Db.Query(`SELECT 
SupplyExecutive.id,
SupplyExecutive.create_date,
SupplyExecutive.first_name,
SupplyExecutive.mob_no,
SupplyExecutive.email,
SupplyExecutive.status,
SupplyExecutive.user_type FROM SupplyExecutive
LEFT JOIN Cluster2SupplyExecutive 
ON SupplyExecutive.id = Cluster2SupplyExecutive.supplyexecutive_id
LEFT JOIN Clusters 
ON Cluster2SupplyExecutive.cluster_id = Clusters.id  where Cluster2SupplyExecutive.cluster_id =? and (SupplyExecutive.first_name='' OR SupplyExecutive.first_name like ?) AND (SupplyExecutive.status='' OR SupplyExecutive.status like ?) AND (SupplyExecutive.mob_no='' OR SupplyExecutive.mob_no like ?)AND (DATE(SupplyExecutive.create_date='') OR DATE (SupplyExecutive.create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(SupplyExecutive.create_date='') OR DATE (SupplyExecutive.create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) ORDER BY SupplyExecutive.create_date DESC`, data1[0][1], input_first_name+"%", input_status+"%", input_mob_no+"%", from, to)

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
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data for hoadmin after filter - ", data, " Data len - ", len(data))
		if len(data) <= 0 {
			err = errors.New("Unable to get the error message info")
			return
		}

		var result []Row

		for i := range data {
			var r Row
			r.Id = data[i][0]
			r.CreateDate = data[i][1]
			r.FirstName = data[i][2]
			r.Mobile = data[i][3]
			r.Email = data[i][4]
			r.Status = data[i][5]
			r.Usertype = data[i][6]
			result = append(result, r)
		}

		c.Data["CustomerData"] = result

	} else if utype == "CLUSTEREXECUTIVE" {

		//NO CODE FOR CLUSTEREXECUTIVE
	} else {

		row, err := db.Db.Query(`select id,create_date,first_name,mob_no,email,status,user_type from SupplyExecutive where (first_name='' OR first_name like ?) AND (status='' OR status like ?) AND (mob_no='' OR mob_no like ?)AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) ORDER BY create_date DESC`, input_first_name+"%", input_status+"%", input_mob_no+"%", from, to)

		//AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y'))

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
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data for hoadmin after filter - ", data, " Data len - ", len(data))
		if len(data) <= 0 {
			err = errors.New("Unable to get the error message info")
			return
		}

		var result []Row

		for i := range data {
			var r Row
			r.Id = data[i][0]
			r.CreateDate = data[i][1]
			r.FirstName = data[i][2]
			r.Mobile = data[i][3]
			r.Email = data[i][4]
			r.Status = data[i][5]
			r.Usertype = data[i][6]
			result = append(result, r)
		}

		c.Data["CustomerData"] = result

	}

	// row, err := db.Db.Query(`SELECT id,first_name,email,status FROM sysuser.credential`)

	// row, err := db.Db.Query(`select first_name,mob_no,email,status,user_type from Sysuser_info where
	// DATE(create_date) >= STR_TO_DATE(?, '%m/%d/%Y') AND DATE(create_date) <= STR_TO_DATE(?, '%m/%d/%Y') ORDER BY create_date DESC`, from, to)

	//row, err := db.Db.Query(`select id,first_name,mob_no,create_date, email,status,user_type from DeliveryUser where (first_name='' OR first_name like ?)AND (email='' OR email like ?) `, input_name+"%", input_email+"%")

	return
}
