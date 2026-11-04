package collectionAgentReport

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
	Oder_id         string
	AgentMobile     string
	AgentName       string
	CustomerMobile  string
	CollectedAmount string
	PaymentType     string
	PaymentStatus   string
	PaymentRemarks  string
	TimeStamp       string
}
type Display5 struct {
	Fields5 []Field5
}

type Field5 struct {
	Id     string
	Mobile string
}
type CollectionAgentReport struct {
	// beego.Controller
	middleware.BaseController
}

func (c *CollectionAgentReport) Get() {
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
			c.TplName = "reports/collectionAgentReport/collectionAgentReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Reports Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "reports/collectionAgentReport/collectionAgentReport.html"
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

	auth, err := utils.IsAuthorized(utype, "CollAgentReport")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	if utype == "CLUSTERMANAGER" {

		uid := sess.Get("uid").(string)
		c.Data["utype"] = utype
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "uid", uid)

		cid, err := services.GetCid(uid)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster ID fetch Failed")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "cid", cid)

		cname, err := services.GetClustername(cid)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster Name fetch Failed")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "cname", cname)
		row, err := db.Db.Query(`SELECT order_id, 
	    agent_mobile, 
		agent_name, 
        customer_mobile, 
		collected_amount, 
		payment_type, 
		payment_status, 
		payment_remarks, 
	    create_date
 FROM CollectionOrderTxn where agent_cluster=?`, cname)

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
			r.Oder_id = data[i][0]
			r.AgentMobile = data[i][1]
			r.AgentName = data[i][2]
			r.CustomerMobile = data[i][3]
			r.CollectedAmount = data[i][4]
			r.PaymentType = data[i][5]
			r.PaymentStatus = data[i][6]
			r.PaymentRemarks = data[i][7]
			r.TimeStamp = data[i][8]

			result = append(result, r)
		}
		c.Data["CustomerData"] = result
	} else if utype == "CLUSTEREXECUTIVE" {
		uid := sess.Get("uid").(string)
		c.Data["utype"] = utype
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "uid", uid)

		cid, err := services.GetCid1(uid)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster ID fetch Failed")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "cid", cid)

		cname, err := services.GetClustername(cid)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster Name fetch Failed")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "cname", cname)
		row, err := db.Db.Query(`SELECT order_id, 
	    agent_mobile, 
		agent_name, 
        customer_mobile, 
		collected_amount, 
		payment_type, 
		payment_status, 
		payment_remarks, 
	    create_date
 FROM CollectionOrderTxn where agent_cluster=?`, cname)

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
			r.Oder_id = data[i][0]
			r.AgentMobile = data[i][1]
			r.AgentName = data[i][2]
			r.CustomerMobile = data[i][3]
			r.CollectedAmount = data[i][4]
			r.PaymentType = data[i][5]
			r.PaymentStatus = data[i][6]
			r.PaymentRemarks = data[i][7]
			r.TimeStamp = data[i][8]

			result = append(result, r)
		}
		c.Data["CustomerData"] = result
	} else {
		row, err := db.Db.Query(`SELECT order_id, 
	    agent_mobile, 
		agent_name, 
        customer_mobile, 
		collected_amount, 
		payment_type, 
		payment_status, 
		payment_remarks, 
	    create_date
 FROM CollectionOrderTxn ORDER BY create_date DESC `)
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
			r.Oder_id = data[i][0]
			r.AgentMobile = data[i][1]
			r.AgentName = data[i][2]
			r.CustomerMobile = data[i][3]
			r.CollectedAmount = data[i][4]
			r.PaymentType = data[i][5]
			r.PaymentStatus = data[i][6]
			r.PaymentRemarks = data[i][7]
			r.TimeStamp = data[i][8]

			result = append(result, r)
		}
		c.Data["CustomerData"] = result
	}

	return
}

func (c *CollectionAgentReport) Post() {
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
			c.TplName = "reports/collectionAgentReport/collectionAgentReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchOrder Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "reports/collectionAgentReport/collectionAgentReport.html"
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

	auth, err := utils.IsAuthorized(utype, "CollAgentReport")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	input_rnn_no := c.Input().Get("input_rnn_no")
	input_status := c.Input().Get("input_status")
	input_card_mask := c.Input().Get("input_card_mask")
	//input_customer_mob := c.Input().Get("input_customer_mob")
	//input_order_id := c.Input().Get("input_order_id")
	//input_cluster_name := c.Input().Get("input_cluster_name")
	input_payment_status := c.Input().Get("input_payment_status")
	input_agent_name := c.Input().Get("input_agent_name")
	input_agent_mob := c.Input().Get("input_agent_mob")

	c.Data["RRN"] = input_rnn_no
	c.Data["Status"] = input_status
	c.Data["Mask"] = input_card_mask
	c.Data["DateRange"] = c.Input().Get("daterange")
	tmp := strings.Split(c.Input().Get("daterange"), " - ")

	from := ""
	to := ""
	if len(tmp) == 2 {
		from = tmp[0]
		to = tmp[1]
	}
	if utype == "CLUSTERMANAGER" {

		uid := sess.Get("uid").(string)
		c.Data["utype"] = utype
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "uid", uid)

		cid, err := services.GetCid(uid)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster ID fetch Failed")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "cid", cid)

		cname, err := services.GetClustername(cid)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster Name fetch Failed")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "cname", cname)
		row, err := db.Db.Query(`SELECT order_id, 
	    agent_mobile, 
		agent_name, 
        customer_mobile, 
		collected_amount, 
		payment_type, 
		payment_status, 
		payment_remarks, 
	    create_date
from CollectionOrderTxn where agent_cluster=? AND (agent_mobile='' OR agent_mobile like ?) AND (agent_name='' OR agent_name like ?)AND (payment_status='' OR payment_status like ?) AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) ORDER BY create_date DESC`, cname, input_agent_mob+"%", input_agent_name+"%", input_payment_status+"%", from, to)
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
			r.Oder_id = data[i][0]
			r.AgentMobile = data[i][1]
			r.AgentName = data[i][2]
			r.CustomerMobile = data[i][3]
			r.CollectedAmount = data[i][4]
			r.PaymentType = data[i][5]
			r.PaymentStatus = data[i][6]
			r.PaymentRemarks = data[i][7]
			r.TimeStamp = data[i][8]

			result = append(result, r)
		}
		c.Data["CustomerData"] = result
	} else if utype == "CLUSTEREXECUTIVE" {
		uid := sess.Get("uid").(string)
		c.Data["utype"] = utype
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "uid", uid)

		cid, err := services.GetCid1(uid)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster ID fetch Failed")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "cid", cid)

		cname, err := services.GetClustername(cid)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster Name fetch Failed")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "cname", cname)
		row, err := db.Db.Query(`SELECT order_id, 
	    agent_mobile, 
		agent_name, 
        customer_mobile, 
		collected_amount, 
		payment_type, 
		payment_status, 
		payment_remarks, 
	    create_date
from CollectionOrderTxn where agent_cluster=? AND (agent_mobile='' OR agent_mobile like ?) AND (agent_name='' OR agent_name like ?)AND (payment_status='' OR payment_status like ?) AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) ORDER BY create_date DESC`, cname, input_agent_mob+"%", input_agent_name+"%", input_payment_status+"%", from, to)
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
			r.Oder_id = data[i][0]
			r.AgentMobile = data[i][1]
			r.AgentName = data[i][2]
			r.CustomerMobile = data[i][3]
			r.CollectedAmount = data[i][4]
			r.PaymentType = data[i][5]
			r.PaymentStatus = data[i][6]
			r.PaymentRemarks = data[i][7]
			r.TimeStamp = data[i][8]

			result = append(result, r)
		}
		c.Data["CustomerData"] = result
	} else {
		row, err := db.Db.Query(`SELECT order_id, 
	    agent_mobile, 
		agent_name, 
        customer_mobile, 
		collected_amount, 
		payment_type, 
		payment_status, 
		payment_remarks, 
	    create_date
from CollectionOrderTxn where (agent_mobile='' OR agent_mobile like ?) AND (agent_name='' OR agent_name like ?)AND (payment_status='' OR payment_status like ?) AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) ORDER BY create_date DESC`, input_agent_mob+"%", input_agent_name+"%", input_payment_status+"%", from, to)

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
			r.Oder_id = data[i][0]
			r.AgentMobile = data[i][1]
			r.AgentName = data[i][2]
			r.CustomerMobile = data[i][3]
			r.CollectedAmount = data[i][4]
			r.PaymentType = data[i][5]
			r.PaymentStatus = data[i][6]
			r.PaymentRemarks = data[i][7]
			r.TimeStamp = data[i][8]

			result = append(result, r)
		}
		c.Data["CustomerData"] = result
	}

	return
}
