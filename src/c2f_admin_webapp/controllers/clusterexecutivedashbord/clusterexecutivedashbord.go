package clusterexecutivedashbord

import (
	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/model/utils"
	"c2f_admin_webapp/services"
	"c2f_admin_webapp/session"
	"errors"
	"runtime/debug"
	"strconv"
	"strings"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"

	"github.com/astaxie/beego"
)

type Row struct {
	SKU           string
	ProductName   string
	Received      string
	Delivered     string
	Avaliable     string
	Qty           string
	Cluster       string
	CustomerOrder string
}

type MenuData struct {
	Menus string
}

type Clusterexecutivedashbord struct {
	// beego.Controller
	middleware.BaseController
}

func (c *Clusterexecutivedashbord) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Clusterdashboard Page Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	var err error
	var Autherr error
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
		if err != nil && err.Error() == "Session Time Out.Please Logout and Login Again." {
			c.Data["DisplayMessage"] = err.Error()
			c.Redirect(beego.AppConfig.String("LOGIN_PATH"), 302)
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Clusterdashboard Page Fail")
		} else if err != nil {
			c.Data["DisplayMessage"] = err.Error()
			c.TplName = "error/error.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Clusterdashboard Page Fail")
		} else {
			c.TplName = "clusterexecutivedashbord/clusterexecutivedashbord.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Clusterdashboard Page Success")
		}
		return
	}()
	//	utils.SetHTTPHeader(c.Ctx)

	sess, err := session.GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System is unable to process your request.Please contact customer care")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "UserName - ", sess.Get("uname"))
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Session ID - ", sess.SessionID())
	if err = session.ValidateSession(sess); err != nil {
		sess.SessionRelease(c.Ctx.ResponseWriter)
		session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}
	defer func() {
		utils.EventLogs(c.Ctx, sess, c.Ctx.Input.Method(), c.Input(), c.Data, err)
		sess.SessionRelease(c.Ctx.ResponseWriter)
	}()

	utype := sess.Get("user_type").(string)
	uid := sess.Get("uid").(string)
	c.Data["utype"] = utype
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "uid", uid)

	language := sess.Get("language").(string)
	c.Data["language"] = language
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "language", language)

	clustername := sess.Get("clustername").(string)
	clustername1 := strings.ToUpper(clustername)
	c.Data["clustername"] = clustername1
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "clustername", clustername1)

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "utype", utype)

	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "fullname", fullname1)

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

	auth, err := utils.IsAuthorized(utype, "ClusterExeDashbord")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}

	row2, err := db.Db.Query(`SELECT count(*) FROM Oder2Delivery where status = 'DELIVERY_ASSIGNED' AND deliveryuser_cluster=?`, cname)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row2)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row2)
	_, data2, err := sql.Scan(row2)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data2, "\nData len - ", len(data2))
	// if len(data2) <= 0 {
	// 	err = errors.New("Unable to get dashboard data")
	// 	return
	// }

	number1, err := strconv.ParseUint(data2[0][0], 10, 32)
	finalIntNum1 := int(number1)

	c.Data["orderedcount"] = finalIntNum1

	row21, err := db.Db.Query(`SELECT count(*) FROM Oder2Delivery where status = 'DELIVERY_DECLINED' AND deliveryuser_cluster=?`, cname)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row21)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row21)
	_, data21, err := sql.Scan(row21)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data21, "\nData len - ", len(data21))
	// if len(data21) <= 0 {
	// 	err = errors.New("Unable to get dashboard data")
	// 	return
	// }

	number6, err := strconv.ParseUint(data21[0][0], 10, 32)
	finalIntNum6 := int(number6)

	c.Data["declinedcount"] = finalIntNum6

	row22, err := db.Db.Query(`SELECT count(*) FROM Oder2Delivery where status = 'DELIVERY_COMPLETED' AND deliveryuser_cluster=?`, cname)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row22)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row22)
	_, data22, err := sql.Scan(row22)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data22, "\nData len - ", len(data22))
	// if len(data22) <= 0 {
	// 	err = errors.New("Unable to get dashboard data")
	// 	return
	// }

	number3, err := strconv.ParseUint(data22[0][0], 10, 32)
	finalIntNum3 := int(number3)

	c.Data["deliverdcount"] = finalIntNum3

	row3, err := db.Db.Query(`SELECT count(*)as acount FROM Oder2Delivery where deliveryuser_cluster=? `, cname)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row3)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
	_, data3, err := sql.Scan(row3)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data3, "\nData len - ", len(data3))
	if len(data3) <= 0 {
		err = errors.New("Unable to get dashboard data")
		return
	}

	number5, err := strconv.ParseUint(data3[0][0], 10, 32)
	finalIntNum5 := int(number5)

	c.Data["recivedcount"] = finalIntNum5

	row13, err := db.Db.Query(`SELECT count(*)as acount FROM CustomerOrderTxn where cluster_name=? and status="ORDER_RECEIVED"`, cname)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row13)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row13)
	_, data13, err := sql.Scan(row13)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data13, "\nData len - ", len(data13))
	if len(data13) <= 0 {
		err = errors.New("Unable to get dashboard data")
		return
	}

	number15, err := strconv.ParseUint(data13[0][0], 10, 32)
	finalIntNum15 := int(number15)

	c.Data["pedningassignments"] = finalIntNum15

	// prunits, err := strconv.ParseInt(data3[0][0], 0, 64)
	// paunits, err := strconv.ParseInt(data2[0][0], 0, 64)

	// TD := prunits - paunits
	// TotalDel := strconv.FormatInt(TD, 10)

	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "pedningassignments - ", TotalDel)

	// c.Data["pedningassignments"] = TotalDel

	productcount, err := services.GetProductCount("ACTIVE")
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Product count fetch Failed")
		return
	}

	procount, err := strconv.ParseUint(productcount, 10, 32)
	produccount := int(procount)
	c.Data["productcount"] = produccount

	row6, err := db.Db.Query(`select SUM(total_selling_price) FROM CustomerOrderTxn as ASum where cluster_name =? and order_purchagetype="KRISHIVALA"`, cname)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row6)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row6)
	_, data6, err := sql.Scan(row6)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data6, "\nData len - ", len(data6))
	if len(data6) <= 0 {
		err = errors.New("Unable to get dashboard data")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "ASum - ", data6[0][0])

	row7, err := db.Db.Query(`select SUM(total_buy_price) FROM SupplyOrderTxn as DSum where cluster_name =?`, cname)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row7)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row7)
	_, data7, err := sql.Scan(row7)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data7, "\nData len - ", len(data7))
	if len(data7) <= 0 {
		err = errors.New("Unable to get dashboard data")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "DSum - ", data7[0][0])
	row8, err := db.Db.Query(`select SUM(collected_amount) FROM CollectionOrderTxn as CSum`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row8)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row7)
	_, data8, err := sql.Scan(row8)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data8, "\nData len - ", len(data8))
	if len(data8) <= 0 {
		err = errors.New("Unable to get dashboard data")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "CSum - ", data8[0][0])

	c.Data["totalcustomertransactionamount"] = data6[0][0]
	c.Data["totalsupplytransactionamount"] = data7[0][0]
	c.Data["totalcollectionagenttransactionamount"] = data8[0][0]
	// c.Data["totaltsuccessfulransactioncount"] = data1[0][2]
	// c.Data["totaltsuccessfulransactionamount"] = data1[0][3]
	// c.Data["toatldeclinedtransactioncount"] = data1[0][4]
	// c.Data["toatldeclinedtransactionamount"] = data1[0][5]

	row9, err := db.Db.Query(`select product_sku,
 product_name, 
	inventory_received,
   inventory_delivered,
inventory_available,
product_package_qty,
product_qty,
cluster_name
 from Inventory where cluster_name=?`, cname)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch order")
		return
	}
	defer sql.Close(row9)
	_, data9, err := sql.Scan(row9)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch order")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data9)
	var result []Row

	for i := range data9 {
		var r Row
		r.SKU = data9[i][0]
		r.ProductName = data9[i][1]
		r.Received = data9[i][2]
		r.Delivered = data9[i][3]
		r.Avaliable = data9[i][4]
		r.Qty = data9[i][5] + " " + data9[i][6]
		r.Cluster = data9[i][7]
		row41, err := db.Db.Query(`SELECT count(*)as acount FROM C2F_NEW.CustomerOrderTxn where product_sku =? AND cluster_name = ?`, data9[i][0], data9[i][7])

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get dashboard data")
			return
		}
		defer sql.Close(row41)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row41)
		_, data41, err := sql.Scan(row41)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get dashboard data")
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data41, "\nData len - ", len(data41))
		if len(data41) <= 0 {
			err = errors.New("Unable to get dashboard data")
			return
		}

		r.CustomerOrder = data41[0][0]

		result = append(result, r)
	}
	c.Data["CustomerData"] = result

	return
}
