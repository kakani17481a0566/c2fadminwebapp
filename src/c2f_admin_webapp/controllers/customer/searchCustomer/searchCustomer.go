package searchCustomer

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
	Id          string
	Timestamp   string
	FirstName   string
	MiddleName  string
	LastName    string
	Mobile      string
	Email       string
	Usertype    string
	Status      string
	ClusterName string
}

type SearchCustomer struct {
	// beego.Controller
	middleware.BaseController
}

func (c *SearchCustomer) Get() {
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
			c.TplName = "customer/searchCustomer/searchCustomer.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "customer/searchCustomer/searchCustomer.html"
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

	auth, err := utils.IsAuthorized(utype, "Customer")
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
		row, err := db.Db.Query(`select id,first_name,mob_no,email,status,cluster_name,user_type,create_date from Customer where cluster_name=? AND (user_type="CONSUMER" OR user_type="PRODUCER+CONSUMER") ORDER BY id DESC`, data1[0][3])

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
			r.ClusterName = data[i][5]
			r.Usertype = data[i][6]
			r.Timestamp = data[i][7]
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
		row, err := db.Db.Query(`select id,first_name,mob_no,email,status,cluster_name,user_type,create_date from Customer where cluster_name=? AND (user_type="CONSUMER" OR user_type="PRODUCER+CONSUMER") ORDER BY id DESC`, data2[0][3])

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
			r.ClusterName = data[i][5]
			r.Usertype = data[i][6]
			r.Timestamp = data[i][7]
			result = append(result, r)
		}

		c.Data["CustomerData"] = result
	} else {
		row, err := db.Db.Query(`select id,first_name,mob_no,email,status,cluster_name,user_type,create_date from Customer where (user_type="CONSUMER" OR user_type="PRODUCER+CONSUMER") ORDER BY id DESC`)

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
			r.ClusterName = data[i][5]
			r.Usertype = data[i][6]
			r.Timestamp = data[i][7]
			result = append(result, r)
		}

		c.Data["CustomerData"] = result

	}

	return
}

func (c *SearchCustomer) Post() {
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
			c.TplName = "customer/searchCustomer/searchCustomer.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "customer/searchCustomer/searchCustomer.html"
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

	auth, err := utils.IsAuthorized(utype, "Customer")
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
	input_cluster_name := c.Input().Get("input_cluster_name")

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
		row, err := db.Db.Query(`select id,first_name,mob_no,email,status,cluster_name,user_type,create_date from Customer where cluster_name=? AND (user_type="CONSUMER" OR user_type="PRODUCER+CONSUMER") AND (first_name='' OR first_name like ?)AND (status='' OR status like ?)AND (email='' OR email like ?) AND (mob_no='' OR mob_no like ?) AND (cluster_name='' OR cluster_name like ?) AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) AND (user_type ="CONSUMER" OR user_type ="PRODUCER+CONSUMER") ORDER BY id DESC`, data1[0][3], input_name+"%", input_status+"%", input_email+"%", input_mob+"%", input_cluster_name+"%", from, to)

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
			r.ClusterName = data[i][5]
			r.Usertype = data[i][6]
			r.Timestamp = data[i][7]
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
		row, err := db.Db.Query(`select id,first_name,mob_no,email,status,cluster_name,user_type,create_date from Customer where cluster_name=? AND (user_type="CONSUMER" OR user_type="PRODUCER+CONSUMER") AND (first_name='' OR first_name like ?)AND (status='' OR status like ?)AND (email='' OR email like ?) AND (mob_no='' OR mob_no like ?) AND (cluster_name='' OR cluster_name like ?) AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) AND (user_type ="CONSUMER" OR user_type ="PRODUCER+CONSUMER") ORDER BY id DESC`, data2[0][3], input_name+"%", input_status+"%", input_email+"%", input_mob+"%", input_cluster_name+"%", from, to)

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
			r.ClusterName = data[i][5]
			r.Usertype = data[i][6]
			r.Timestamp = data[i][7]
			result = append(result, r)
		}

		c.Data["CustomerData"] = result
	} else {
		row, err := db.Db.Query(`select id,first_name,mob_no,email,status,cluster_name,user_type,create_date from Customer where (first_name='' OR first_name like ?)AND (status='' OR status like ?)AND (email='' OR email like ?) AND (mob_no='' OR mob_no like ?) AND (cluster_name='' OR cluster_name like ?) AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) AND (user_type ="CONSUMER" OR user_type ="PRODUCER+CONSUMER") ORDER BY id DESC`, input_name+"%", input_status+"%", input_email+"%", input_mob+"%", input_cluster_name+"%", from, to)

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
			r.ClusterName = data[i][5]
			r.Usertype = data[i][6]
			r.Timestamp = data[i][7]
			result = append(result, r)
		}

		c.Data["CustomerData"] = result

	}

	// row, err := db.Db.Query(`SELECT id,first_name,email,status FROM sysuser.credential`)

	// row, err := db.Db.Query(`select first_name,mob_no,email,status,user_type from Sysuser_info where
	// DATE(create_date) >= STR_TO_DATE(?, '%m/%d/%Y') AND DATE(create_date) <= STR_TO_DATE(?, '%m/%d/%Y') ORDER BY create_date DESC`, from, to)

	//row, err := db.Db.Query(`select id,first_name,mob_no,email,status,user_type from Customer where user_type ="CONSUMER" AND (first_name='' OR first_name like ?)AND (email='' OR email like ?) AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) ORDER BY create_date DESC`, input_name+"%", input_email+"%", from, to)
	row, err := db.Db.Query(`select id,first_name,mob_no,email,status,cluster_name,user_type from Customer where (first_name='' OR first_name like ?)AND (status='' OR status like ?)AND (email='' OR email like ?) AND (mob_no='' OR mob_no like ?) AND (cluster_name='' OR cluster_name like ?) AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) AND (user_type ="CONSUMER" OR user_type ="PRODUCER+CONSUMER") ORDER BY id DESC`, input_name+"%", input_status+"%", input_email+"%", input_mob+"%", input_cluster_name+"%", from, to)

	//`select id,first_name,mob_no,email,status,user_type, create_date from Customer where (first_name='' OR first_name like ?)AND (email='' OR email like ?) AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) ORDER BY create_date DESC`, input_name+"%", input_email+"%", from, to)
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
		r.ClusterName = data[i][5]
		r.Usertype = data[i][6]
		result = append(result, r)
	}

	c.Data["CustomerData"] = result

	return
}
