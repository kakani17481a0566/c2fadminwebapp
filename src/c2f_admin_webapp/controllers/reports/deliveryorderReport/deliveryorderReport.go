package deliveryorderReport

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

// type Row struct {
// 	Producer_Mob string
// 	Product_Name string
// 	Cluster_Name string
// 	Qty          string
// 	Price        string
// 	Start_Date   string
// 	End_Date     string
// 	TimeStmp     string
// }

type Row struct {
	Order_id           string
	Order_date         string
	DeliveryUserMob    string
	DeliveryassignDate string
	Status             string
	Remarks            string
	ClusterName        string
}
type Display5 struct {
	Fields5 []Field5
}

type Field5 struct {
	Id     string
	Mobile string
}
type DeliveryorderReport struct {
	// beego.Controller
	middleware.BaseController
}

func (c *DeliveryorderReport) Get() {
	Status := c.Ctx.Input.Param(":Status")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Status", Status)

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
			c.TplName = "reports/deliveryorderReport/deliveryorderReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Reports Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "reports/deliveryorderReport/deliveryorderReport.html"
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

	auth, err := utils.IsAuthorized(utype, "DveryOrderReport")
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

		if Status == "" {
			row, err := db.Db.Query(`select order_date,oder_id, deliveryuser_mob, create_date,status,remarks,JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.clusterName')) from Oder2Delivery where deliveryuser_cluster=? `, cname)

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
				r.Order_date = data[i][0]
				r.Order_id = data[i][1]
				r.DeliveryUserMob = data[i][2]
				r.DeliveryassignDate = data[i][3]
				r.Status = data[i][4]
				r.Remarks = data[i][5]
				r.ClusterName = data[i][6]

				result = append(result, r)
			}
			c.Data["CustomerData"] = result

		} else {
			row, err := db.Db.Query(`select order_date,oder_id, deliveryuser_mob, create_date,status,remarks,JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.clusterName')) from Oder2Delivery where (status='' OR status like ?) AND deliveryuser_cluster=?`, Status, cname)

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
				r.Order_date = data[i][0]
				r.Order_id = data[i][1]
				r.DeliveryUserMob = data[i][2]
				r.DeliveryassignDate = data[i][3]
				r.Status = data[i][4]
				r.Remarks = data[i][5]
				r.ClusterName = data[i][6]

				result = append(result, r)
			}
			c.Data["CustomerData"] = result
		}
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
		if Status == "" {
			row, err := db.Db.Query(`select order_date,oder_id, deliveryuser_mob, create_date,status,remarks,JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.clusterName')) from Oder2Delivery where deliveryuser_cluster=?`, cname)

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
				r.Order_date = data[i][0]
				r.Order_id = data[i][1]
				r.DeliveryUserMob = data[i][2]
				r.DeliveryassignDate = data[i][3]
				r.Status = data[i][4]
				r.Remarks = data[i][5]
				r.ClusterName = data[i][6]

				result = append(result, r)
			}
			c.Data["CustomerData"] = result

		} else {
			row, err := db.Db.Query(`select order_date,oder_id, deliveryuser_mob, create_date,status,remarks,JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.clusterName')) from Oder2Delivery where (status='' OR status like ?) AND deliveryuser_cluster=?`, Status, cname)

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
				r.Order_date = data[i][0]
				r.Order_id = data[i][1]
				r.DeliveryUserMob = data[i][2]
				r.DeliveryassignDate = data[i][3]
				r.Status = data[i][4]
				r.Remarks = data[i][5]
				r.ClusterName = data[i][6]

				result = append(result, r)
			}
			c.Data["CustomerData"] = result
		}
	} else {
		if Status == "" {
			row, err := db.Db.Query(`select order_date,oder_id, deliveryuser_mob, create_date,status,remarks,JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.clusterName')) from Oder2Delivery `)

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
				r.Order_date = data[i][0]
				r.Order_id = data[i][1]
				r.DeliveryUserMob = data[i][2]
				r.DeliveryassignDate = data[i][3]
				r.Status = data[i][4]
				r.Remarks = data[i][5]
				r.ClusterName = data[i][6]

				result = append(result, r)
			}
			c.Data["CustomerData"] = result

		} else {
			row, err := db.Db.Query(`select order_date,oder_id, deliveryuser_mob, create_date,status,remarks,JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.clusterName')) from Oder2Delivery where (status='' OR status like ?)`, Status)

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
				r.Order_date = data[i][0]
				r.Order_id = data[i][1]
				r.DeliveryUserMob = data[i][2]
				r.DeliveryassignDate = data[i][3]
				r.Status = data[i][4]
				r.Remarks = data[i][5]
				r.ClusterName = data[i][6]

				result = append(result, r)
			}
			c.Data["CustomerData"] = result
		}
	}

	return
}

func (c *DeliveryorderReport) Post() {

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
			c.TplName = "reports/deliveryorderReport/deliveryorderReport.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchOrder Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "reports/deliveryorderReport/deliveryorderReport.html"
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

	auth, err := utils.IsAuthorized(utype, "DveryOrderReport")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	input_rnn_no := c.Input().Get("input_rnn_no")
	input_status := c.Input().Get("input_status")
	input_card_mask := c.Input().Get("input_card_mask")
	input_order_id := c.Input().Get("input_order_id")

	input_mob_no := c.Input().Get("input_mob_no")

	//input_customer_mob := c.Input().Get("input_customer_mob")

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

		row, err := db.Db.Query(`select order_date,oder_id, deliveryuser_mob, create_date,status,remarks,JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.clusterName')) from Oder2Delivery where deliveryuser_cluster=? AND (deliveryuser_mob='' OR deliveryuser_mob like ?) AND (oder_id='' OR oder_id like ?) AND (status='' OR status like ?)  AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) ORDER BY create_date DESC`, cname, input_mob_no+"%", input_order_id+"%", input_status+"%", from, to)

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
			r.Order_date = data[i][0]
			r.Order_id = data[i][1]
			r.DeliveryUserMob = data[i][2]
			r.DeliveryassignDate = data[i][3]
			r.Status = data[i][4]
			r.Remarks = data[i][5]
			r.ClusterName = data[i][6]

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
		row, err := db.Db.Query(`select order_date,oder_id, deliveryuser_mob, create_date,status,remarks,JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.clusterName')) from Oder2Delivery where deliveryuser_cluster=? AND (deliveryuser_mob='' OR deliveryuser_mob like ?) AND (oder_id='' OR oder_id like ?) AND (status='' OR status like ?)  AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) ORDER BY create_date DESC`, cname, input_mob_no+"%", input_order_id+"%", input_status+"%", from, to)

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
			r.Order_date = data[i][0]
			r.Order_id = data[i][1]
			r.DeliveryUserMob = data[i][2]
			r.DeliveryassignDate = data[i][3]
			r.Status = data[i][4]
			r.Remarks = data[i][5]
			r.ClusterName = data[i][6]

			result = append(result, r)
		}
		c.Data["CustomerData"] = result
	} else {
		row, err := db.Db.Query(`select order_date,oder_id, deliveryuser_mob, create_date,status,remarks,JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.clusterName')) from Oder2Delivery where (deliveryuser_mob='' OR deliveryuser_mob like ?) AND (oder_id='' OR oder_id like ?) AND (status='' OR status like ?)  AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y')) ORDER BY create_date DESC`, input_mob_no+"%", input_order_id+"%", input_status+"%", from, to)

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
			r.Order_date = data[i][0]
			r.Order_id = data[i][1]
			r.DeliveryUserMob = data[i][2]
			r.DeliveryassignDate = data[i][3]
			r.Status = data[i][4]
			r.Remarks = data[i][5]
			r.ClusterName = data[i][6]

			result = append(result, r)
		}
		c.Data["CustomerData"] = result
	}

	return
}

// row, err := db.Db.Query(`select producer_mob, product_name, cluster_name, qty, unit_price,price, start_date, end_date, create_date from SupplyForecastTxn`)

// if err != nil {
// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
// 	err = errors.New("Unable to fetch order")
// 	return
// }
// defer sql.Close(row)
// _, data, err := sql.Scan(row)
// if err != nil {
// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
// 	err = errors.New("Unable to fetch order")
// 	return
// }

// log.Println(beego.AppConfig.String("loglevel"), "Debug", data)
// var result []Row

// for i := range data {
// 	var r Row
// 	r.Producer_Mob = data[i][0]
// 	r.Product_Name = data[i][1]
// 	r.Cluster_Name = data[i][2]
// 	r.Qty = data[i][3]
// 	r.Price = data[i][4]
// 	r.Start_Date = data[i][5]
// 	r.End_Date = data[i][6]
// 	r.TimeStmp = data[i][7]
// 	result = append(result, r)
// }
// c.Data["CustomerData"] = result

// return
// }
