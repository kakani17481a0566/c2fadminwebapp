package customerForecastReport

import (
	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/model/utils"

	"c2f_admin_webapp/services"
	"c2f_admin_webapp/session"
	"strings"

	"fmt"
	"runtime/debug"

	"errors"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

type Row struct {
	Id                string
	CustomerMob       string
	Clustername       string
	Productname       string
	ProductQty        string
	Productunits      string
	Producersellprice string
	From              string
	To                string
	CustomerName      string
	Timestamp         string
	Version           string
	Uploaded_By       string
}

type CustomerForecastReport struct {
	// beego.Controller
	middleware.BaseController
}

func (c *CustomerForecastReport) Get() {

	AdminID := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "AdminID", AdminID)

	log.Println(beego.AppConfig.String("loglevel"), "Info", "Customer Forecast Report Start")
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
			c.TplName = "analytics/forecast/customerForecastReport/customerForecastReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Customer Forecast Report Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "analytics/forecast/customerForecastReport/customerForecastReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Customer Forecast Report Page Success")
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
	//uid := sess.Get("uid").(string)
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

	auth, err := utils.IsAuthorized(utype, "CustForecastReport")
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
		row, err := db.Db.Query(`select id,customer_mob,cluster_name,product_name,product_qty,units,customer_name,start_date,end_date,create_date,version,uploaded_by from CustomerForecast where (cluster_name='' OR cluster_name like ?)`, cname)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch ProducerForecastData")
			return
		}
		defer sql.Close(row)
		_, data, err := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch ProducerForecastData")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "CustomerForecast Data", data)
		var result []Row

		for i := range data {
			var r Row
			r.Id = data[i][0]
			r.CustomerMob = data[i][1]
			r.Clustername = data[i][2]
			r.Productname = data[i][3]
			r.ProductQty = data[i][4]
			r.Productunits = data[i][5]
			r.CustomerName = data[i][6]
			r.From = data[i][7]
			r.To = data[i][8]
			r.Timestamp = data[i][9]
			r.Version = data[i][10]
			r.Uploaded_By = data[i][11]
			result = append(result, r)
		}
		c.Data["CustomerForecastData"] = result

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
		row, err := db.Db.Query(`select id,customer_mob,cluster_name,product_name,product_qty,units,customer_name,start_date,end_date,create_date,version,uploaded_by from CustomerForecast where (cluster_name='' OR cluster_name like ?)`, cname)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch ProducerForecast")
			return
		}
		defer sql.Close(row)
		_, data, err := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch ProducerForecast")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "CustomerForecast data", data)
		var result []Row

		for i := range data {
			var r Row
			r.Id = data[i][0]
			r.CustomerMob = data[i][1]
			r.Clustername = data[i][2]
			r.Productname = data[i][3]
			r.ProductQty = data[i][4]
			r.Productunits = data[i][5]
			r.CustomerName = data[i][6]
			r.From = data[i][7]
			r.To = data[i][8]
			r.Timestamp = data[i][9]
			r.Version = data[i][10]
			r.Uploaded_By = data[i][11]
			result = append(result, r)
		}
		c.Data["CustomerForecastData"] = result

	} else {

		if AdminID == "" {

			row, err := db.Db.Query(`select id,customer_mob,cluster_name,product_name,product_qty,units,customer_name,start_date,end_date,create_date,version,uploaded_by from CustomerForecast`)

			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to ProducerForecast")
				return
			}
			defer sql.Close(row)
			_, data, err := sql.Scan(row)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to ProducerForecast")
				return
			}

			log.Println(beego.AppConfig.String("loglevel"), "CustomerForecast data", data)
			var result []Row

			for i := range data {
				var r Row
				r.Id = data[i][0]
				r.CustomerMob = data[i][1]
				r.Clustername = data[i][2]
				r.Productname = data[i][3]
				r.ProductQty = data[i][4]
				r.Productunits = data[i][5]
				r.CustomerName = data[i][6]
				r.From = data[i][7]
				r.To = data[i][8]
				r.Timestamp = data[i][9]
				r.Version = data[i][10]
				r.Uploaded_By = data[i][11]

				result = append(result, r)
			}
			c.Data["CustomerForecastData"] = result

		} else {

			row, err := db.Db.Query(`select id,customer_mob,cluster_name,product_name,product_qty,units,customer_name,start_date,end_date,create_date,version,uploaded_by from CustomerForecast where product_name=?`, AdminID)

			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to ProducerForecast")
				return
			}
			defer sql.Close(row)
			_, data, err := sql.Scan(row)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to ProducerForecast")
				return
			}

			log.Println(beego.AppConfig.String("loglevel"), "CustomerForecast data", data)
			var result []Row

			for i := range data {
				var r Row
				r.Id = data[i][0]
				r.CustomerMob = data[i][1]
				r.Clustername = data[i][2]
				r.Productname = data[i][3]
				r.ProductQty = data[i][4]
				r.Productunits = data[i][5]
				r.CustomerName = data[i][6]
				r.From = data[i][7]
				r.To = data[i][8]
				r.Timestamp = data[i][9]
				r.Version = data[i][10]
				r.Uploaded_By = data[i][11]

				result = append(result, r)
			}
			c.Data["CustomerForecastData"] = result

		}

	}

	return
}

func (c *CustomerForecastReport) Post() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Customer Forecast Report Page Start")
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
			c.TplName = "analytics/forecast/customerForecastReport/customerForecastReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "CustomerForecast Report Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "analytics/forecast/customerForecastReport/customerForecastReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "CustomerForecast Report  Page Success")
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
	//uid := sess.Get("uid").(string)
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

	// cid, err := services.GetCid(uid)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Cluster ID fetch Failed")
	// 	return
	// }

	auth, err := utils.IsAuthorized(utype, "CustForecastReport")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	input_product_name := c.Input().Get("input_producer_id")
	input_cluster_name := c.Input().Get("input_cluster_name")
	input_customer_mob := c.Input().Get("input_customer_mob")

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
		row, err := db.Db.Query(`select id,customer_mob,cluster_name,product_name,product_qty,units,customer_name,start_date,end_date,create_date,version,uploaded_by from CustomerForecast  where (cluster_name='' OR cluster_name like ?) AND (customer_mob='' OR customer_mob like ?) AND (cluster_name='' OR cluster_name like ?) AND (product_name='' OR product_name like ?)
		AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y'))`, cname, input_customer_mob+"%", input_cluster_name+"%", input_product_name+"%", from, to)
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
			r.Id = data[i][0]
			r.CustomerMob = data[i][1]
			r.Clustername = data[i][2]
			r.Productname = data[i][3]
			r.ProductQty = data[i][4]
			r.Productunits = data[i][5]
			r.CustomerName = data[i][6]
			r.From = data[i][7]
			r.To = data[i][8]
			r.Timestamp = data[i][9]
			r.Version = data[i][10]
			r.Uploaded_By = data[i][11]
			result = append(result, r)
		}
		c.Data["CustomerForecastData"] = result
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
		row, err := db.Db.Query(`select id,customer_mob,cluster_name,product_name,product_qty,units,customer_name,start_date,end_date,create_date,version,uploaded_by from CustomerForecast  where (cluster_name='' OR cluster_name like ?) AND (customer_mob='' OR customer_mob like ?) AND (cluster_name='' OR cluster_name like ?) AND (product_name='' OR product_name like ?)
		AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y'))`, cname, input_customer_mob+"%", input_cluster_name+"%", input_product_name+"%", from, to)
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
			r.Id = data[i][0]
			r.CustomerMob = data[i][1]
			r.Clustername = data[i][2]
			r.Productname = data[i][3]
			r.ProductQty = data[i][4]
			r.Productunits = data[i][5]
			r.CustomerName = data[i][6]
			r.From = data[i][7]
			r.To = data[i][8]
			r.Timestamp = data[i][9]
			r.Version = data[i][10]
			r.Uploaded_By = data[i][11]
			result = append(result, r)
		}
		c.Data["CustomerForecastData"] = result
	} else {
		fmt.Println("input_product_name", input_product_name)
		fmt.Println("input_customer_mob", input_customer_mob)

		row, err := db.Db.Query(`select id,customer_mob,cluster_name,product_name,product_qty,units,customer_name,start_date,end_date,create_date,version,uploaded_by from CustomerForecast  where (customer_mob='' OR customer_mob like ?) AND (cluster_name='' OR cluster_name like ?) AND (product_name='' OR product_name like ?)
		AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y'))`, input_customer_mob+"%", input_cluster_name+"%", input_product_name+"%", from, to)
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
			r.Id = data[i][0]
			r.CustomerMob = data[i][1]
			r.Clustername = data[i][2]
			r.Productname = data[i][3]
			r.ProductQty = data[i][4]
			r.Productunits = data[i][5]
			r.CustomerName = data[i][6]
			r.From = data[i][7]
			r.To = data[i][8]
			r.Timestamp = data[i][9]
			r.Version = data[i][10]
			r.Uploaded_By = data[i][11]
			result = append(result, r)
		}
		c.Data["CustomerForecastData"] = result
	}

	return
}
