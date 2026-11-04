package collectPayment

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
	Id                  string
	Order_id            string
	Order_date          string
	Collectionagent_mob string
	DeliveryassignDate  string
	Status              string
	Name                string
	Add1                string
	Add2                string
	City                string
	Mobile              string
	Latlong             string
	LocationURL         string
	ClusterName         string
	ProductName         string
	ProductPackageQty   string
	ProductQty          string
	DeliveryTime        string
	DeliveryFrom        string
	DeliveryTo          string
	ProductPrice        string
	Productunits        string
	PaymentStatus       string
}
type Display5 struct {
	Fields5 []Field5
}

type Field5 struct {
	Id     string
	Mobile string
	Name   string
}
type CollectPayment struct {
	// beego.Controller
	middleware.BaseController
}

func (c *CollectPayment) Get() {
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
			c.TplName = "collectpayment/collectPayment/collectPayment.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Reports Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "collectpayment/collectPayment/collectPayment.html"
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

	id := sess.Get("uid").(string)
	c.Data["id"] = id

	language := sess.Get("language").(string)
	c.Data["language"] = language
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "language", language)

	auth, err := utils.IsAuthorized(utype, "CollectPayment")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	clustername := sess.Get("clustername").(string)
	clustername1 := strings.ToUpper(clustername)
	c.Data["clustername"] = clustername1
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "clustername", clustername1)

	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1

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

		data, err := services.GetClustermappedagent(data1[0][1])
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
			d.Mobile = data[i][3]
			d.Name = data[i][2]
			Dis.Fields5 = append(Dis.Fields5, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

	} else if utype == "CLUSTEREXECUTIVE" {

		data, err := services.Getcollectionagent()
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
			d.Mobile = data[i][3]
			d.Name = data[i][2]
			Dis.Fields5 = append(Dis.Fields5, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

	} else {

		data, err := services.Getcollectionagent()
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
			d.Mobile = data[i][3]
			d.Name = data[i][2]
			Dis.Fields5 = append(Dis.Fields5, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

	}

	// row, err := db.Db.Query(`select order_date,oder_id, collectionagent_mob, create_date,status,id,
	// JSON_UNQUOTE(JSON_EXTRACT(customer_details, '$.name')),
	// JSON_UNQUOTE(JSON_EXTRACT(customer_details, '$.address1')),
	// JSON_UNQUOTE(JSON_EXTRACT(customer_details, '$.address2')),
	// JSON_UNQUOTE(JSON_EXTRACT(customer_details, '$.city')),
	// JSON_UNQUOTE(JSON_EXTRACT(customer_details, '$.mobile')),
	// JSON_UNQUOTE(JSON_EXTRACT(customer_details, '$.latlong')),
	// JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.clusterName')),
	// JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.productName')),
	// JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.productPackageQty')),
	// JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.productQty')),
	// JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.deliveryTime')),
	// JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.deliveryFrom')),
	// JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.deliveryTo')),
	// JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.productPrice')),
	// JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.productUnits')),
	// JSON_UNQUOTE(JSON_EXTRACT(customer_details, '$.locationURL'))
	//  from Order2Collectionagent`)

	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Unable to fetch order")
	// 	return
	// }
	// defer sql.Close(row)
	// _, data, err1 := sql.Scan(row)
	// if err1 != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
	// 	err1 = errors.New("Unable to fetch order")
	// 	return
	// }

	// log.Println(beego.AppConfig.String("loglevel"), "Debug", data)
	// var result []Row

	// for i := range data {
	// 	var r Row
	// 	r.Order_date = data[i][0]
	// 	r.Order_id = data[i][1]
	// 	r.Collectionagent_mob = data[i][2]
	// 	r.DeliveryassignDate = data[i][3]
	// 	r.Status = data[i][4]
	// 	r.Id = data[i][5]
	// 	r.Name = data[i][6]
	// 	r.Add1 = data[i][7] + data[i][8]
	// 	//r.Add2 = data[i][8]
	// 	r.City = data[i][9]
	// 	r.Mobile = data[i][10]
	// 	r.Latlong = data[i][11]
	// 	r.ClusterName = data[i][12]
	// 	r.ProductName = data[i][13]
	// 	r.ProductPackageQty = data[i][14] + data[i][15]
	// 	//r.ProductQty = data[i][15]
	// 	r.DeliveryTime = data[i][16]
	// 	r.DeliveryFrom = data[i][17] + "to" + data[i][18]
	// 	//r.DeliveryTo = data[i][18]
	// 	r.ProductPrice = data[i][19]
	// 	r.Productunits = data[i][20]
	// 	r.LocationURL = data[i][21]

	// 	result = append(result, r)
	// }
	// c.Data["CustomerData"] = result

	return
}

func (c *CollectPayment) Post() {
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
			c.TplName = "collectpayment/collectPayment/collectPayment.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchOrder Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "collectpayment/collectPayment/collectPayment.html"
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

	data1, err := services.Getcollectionagent()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Producer fetch Failed")
		return
	}
	if len(data1) <= 0 {
		err = errors.New("Producer  Not Found")
		return
	}
	var Dis Display5
	for i := range data1 {
		var d Field5
		d.Id = data1[i][0]
		d.Mobile = data1[i][3]
		d.Name = data1[i][2]
		Dis.Fields5 = append(Dis.Fields5, d)
	}
	c.Data["Dis"] = Dis
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

	auth, err := utils.IsAuthorized(utype, "CollectPayment")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	input_rnn_no := c.Input().Get("input_rnn_no")
	//input_status := c.Input().Get("input_status")
	input_card_mask := c.Input().Get("input_card_mask")
	input_collectionagent_mob := c.Input().Get("input_collectionagent_mob")
	input_order_id := c.Input().Get("input_order_id")
	//	input_cluster_name := c.Input().Get("input_cluster_name")
	//input_status := c.Input().Get("input_status")

	c.Data["RRN"] = input_rnn_no
	//c.Data["Status"] = input_status
	c.Data["Mask"] = input_card_mask
	c.Data["DateRange"] = c.Input().Get("daterange")
	//tmp := strings.Split(c.Input().Get("daterange"), " - ")

	// from := ""
	// to := ""
	// if len(tmp) == 2 {
	// 	from = tmp[0]
	// 	to = tmp[1]
	// }

	row, err := db.Db.Query(`select Order2Collectionagent.order_date,Order2Collectionagent.oder_id, Order2Collectionagent.collectionagent_mob,Order2Collectionagent.create_date,Order2Collectionagent.status,Order2Collectionagent.id,
	JSON_UNQUOTE(JSON_EXTRACT(Order2Collectionagent.customer_details, '$.name')),
	JSON_UNQUOTE(JSON_EXTRACT(Order2Collectionagent.customer_details, '$.address1')),
	JSON_UNQUOTE(JSON_EXTRACT(Order2Collectionagent.customer_details, '$.address2')),
	JSON_UNQUOTE(JSON_EXTRACT(Order2Collectionagent.customer_details, '$.city')),
	JSON_UNQUOTE(JSON_EXTRACT(Order2Collectionagent.customer_details, '$.mobile')),
	JSON_UNQUOTE(JSON_EXTRACT(Order2Collectionagent.customer_details, '$.latlong')),
	JSON_UNQUOTE(JSON_EXTRACT(Order2Collectionagent.oder_deatils, '$.clusterName')),
	JSON_UNQUOTE(JSON_EXTRACT(Order2Collectionagent.oder_deatils, '$.productName')),
	JSON_UNQUOTE(JSON_EXTRACT(Order2Collectionagent.oder_deatils, '$.productPackageQty')),
	JSON_UNQUOTE(JSON_EXTRACT(Order2Collectionagent.oder_deatils, '$.productQty')),
	JSON_UNQUOTE(JSON_EXTRACT(Order2Collectionagent.oder_deatils, '$.deliveryTime')),
	JSON_UNQUOTE(JSON_EXTRACT(Order2Collectionagent.oder_deatils, '$.deliveryFrom')),
	JSON_UNQUOTE(JSON_EXTRACT(Order2Collectionagent.oder_deatils, '$.deliveryTo')),
	JSON_UNQUOTE(JSON_EXTRACT(Order2Collectionagent.oder_deatils, '$.productPrice')),
	JSON_UNQUOTE(JSON_EXTRACT(Order2Collectionagent.oder_deatils, '$.productUnits')),
	JSON_UNQUOTE(JSON_EXTRACT(Order2Collectionagent.customer_details, '$.locationURL')),
	CustomerOrderTxn.order_id,
	CustomerOrderTxn.payment_status
FROM Order2Collectionagent
LEFT JOIN CustomerOrderTxn
ON Order2Collectionagent.oder_id = CustomerOrderTxn.order_id 
WHERE (Order2Collectionagent.oder_id='' OR Order2Collectionagent.oder_id like ?) AND (Order2Collectionagent.collectionagent_mob='' OR Order2Collectionagent.collectionagent_mob like ?) AND (CustomerOrderTxn.payment_status="PENDING") ORDER BY Order2Collectionagent.create_date DESC`, input_order_id+"%", input_collectionagent_mob+"%")

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
		r.Collectionagent_mob = data[i][2]
		r.DeliveryassignDate = data[i][3]
		r.Status = data[i][4]
		r.Id = data[i][5]
		r.Name = data[i][6]
		r.Add1 = data[i][7] + data[i][8]
		//r.Add2 = data[i][8]
		r.City = data[i][9]
		r.Mobile = data[i][10]
		r.Latlong = data[i][11]
		r.ClusterName = data[i][12]
		r.ProductName = data[i][13]
		r.ProductPackageQty = data[i][14] + data[i][15]
		//r.ProductQty = data[i][15]
		r.DeliveryTime = data[i][16]
		r.DeliveryFrom = data[i][17] + "to" + data[i][18]
		//r.DeliveryTo = data[i][18]
		r.ProductPrice = data[i][19]
		r.Productunits = data[i][20]
		r.LocationURL = data[i][21]
		r.PaymentStatus = data[i][23]

		result = append(result, r)
	}
	c.Data["CustomerData"] = result

	return
}
