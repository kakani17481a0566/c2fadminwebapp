package dashboard

import (
	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/model/utils"
	"c2f_admin_webapp/services"
	"c2f_admin_webapp/session"
	"errors"
	"strings"

	// "fmt"
	"runtime/debug"
	"strconv"

	// "time"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"

	"github.com/astaxie/beego"
)

type Display1 struct {
	Fields1 []Field1
}

type Field1 struct {
	Supplyordertxn_id  string
	ProducerMobile     string
	Status             string
	ProductName        string
	ClusterName        string
	Productqty         string
	PPackageQty        string
	Createdate         string
	ProductSupplyUnits string
	NotificationTime   string
}

type Display2 struct {
	Fields2 []Field2
}

type Field2 struct {
	CustomerMobile   string
	ClusterName      string
	ProductName      string
	PPackageQty      string
	Productqty       string
	Status           string
	OrderUnits       string
	OrderId          string
	Createdate       string
	NotificationTime string
	ClusterType      string
}

// type Display2 struct {
// 	Fields2 []Field2
// }

// type Field2 struct {
// 	NotificationTime string
// }

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

type Dashboard struct {
	// beego.Controller
	middleware.BaseController
}

func (c *Dashboard) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Dashboard Page Start")
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
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Dashboard Page Fail")
		} else if err != nil {
			c.Data["DisplayMessage"] = err.Error()
			c.TplName = "error/error.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Dashboard Page Fail")
		} else {
			c.TplName = "dashboard/dashboard.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Dashboard Page Success")
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
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "fullname", fullname1)

	auth, err := utils.IsAuthorized(utype, "Dashboard")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	//row, err := db.Db.Query(`SELECT count(*)as acount,(SELECT count(*) FROM User where role = '1') as mcount FROM User where role = '2'`)
	row, err := db.Db.Query(`SELECT count(*)as acount,(SELECT count(*) FROM User where role = '2'),(SELECT count(*) FROM User where role = '3'),(SELECT count(*) FROM User where role = '4'),(SELECT count(*) FROM User where role = '5') as mcount FROM User where role = '1'`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get dashboard data")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data1 - ", data[0][0])
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data2 - ", data[0][1])
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data3 - ", data[0][2])
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data4 - ", data[0][3])
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data5 - ", data[0][4])

	num, err := strconv.ParseUint(data[0][0], 10, 32)
	fina1 := int(num)

	num1, err := strconv.ParseUint(data[0][1], 10, 32)
	fina2 := int(num1)

	num2, err := strconv.ParseUint(data[0][2], 10, 32)
	fina3 := int(num2)

	num3, err := strconv.ParseUint(data[0][3], 10, 32)
	fina4 := int(num3)

	num4, err := strconv.ParseUint(data[0][4], 10, 32)
	fina5 := int(num4)

	c.Data["admincount"] = fina1
	c.Data["financecount"] = fina2
	c.Data["helpcount"] = fina3
	c.Data["managercount"] = fina4
	c.Data["executivecount"] = fina5

	row1, err := db.Db.Query(`SELECT count(*)as acount,(SELECT count(*) FROM Customer where (user_type = 'PRODUCER' or user_type='PRODUCER+CONSUMER') ) as mcount FROM Customer where (user_type="CONSUMER" or user_type="PRODUCER+CONSUMER")`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row1)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row1)
	_, data1, err := sql.Scan(row1)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data1, "\nData len - ", len(data1))
	if len(data1) <= 0 {
		err = errors.New("Unable to get dashboard data")
		return
	}

	number, err := strconv.ParseUint(data1[0][0], 10, 32)
	finalIntNum := int(number)

	number2, err := strconv.ParseUint(data1[0][1], 10, 32)
	finalIntNum2 := int(number2)

	c.Data["producercount"] = finalIntNum
	c.Data["consumercount"] = finalIntNum2

	row4, err := db.Db.Query(`SELECT count(*)as acount FROM DeliveryUser where status = 'ACTIVE'`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row1)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row4)
	_, data4, err := sql.Scan(row4)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data4, "\nData len - ", len(data4))
	if len(data4) <= 0 {
		err = errors.New("Unable to get dashboard data")
		return
	}

	number4, err := strconv.ParseUint(data4[0][0], 10, 32)
	finalIntNum4 := int(number4)

	c.Data["deliveryusercount"] = finalIntNum4

	row5, err := db.Db.Query(`SELECT count(*)as acount FROM Collection_Agent where status = 'ACTIVE'`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row5)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row4)
	_, data5, err := sql.Scan(row5)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data5, "\nData len - ", len(data5))
	if len(data5) <= 0 {
		err = errors.New("Unable to get dashboard data")
		return
	}

	number8, err := strconv.ParseUint(data5[0][0], 10, 32)
	finalIntNum8 := int(number8)

	c.Data["collectionagentcount"] = finalIntNum8

	row2, err := db.Db.Query(`SELECT count(*)as acount,(SELECT count(*) FROM Oder2Delivery where status = 'DELIVERY_ASSIGNED'),(SELECT count(*) FROM Oder2Delivery where status = 'DELIVERY_DECLINED') as mcount FROM Oder2Delivery where status = 'DELIVERY_COMPLETED'`)

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
	if len(data2) <= 0 {
		err = errors.New("Unable to get dashboard data")
		return
	}

	number1, err := strconv.ParseUint(data2[0][0], 10, 32)
	finalIntNum1 := int(number1)

	number3, err := strconv.ParseUint(data2[0][1], 10, 32)
	finalIntNum3 := int(number3)

	number6, err := strconv.ParseUint(data2[0][2], 10, 32)
	finalIntNum6 := int(number6)

	c.Data["orderedcount"] = finalIntNum1
	c.Data["pendingcount"] = finalIntNum3
	c.Data["declinedcount"] = finalIntNum6

	row3, err := db.Db.Query(`SELECT count(*)as acount FROM Oder2Delivery`)

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

	//c = finalIntNum5 - finalIntNum6

	prunits, err := strconv.ParseInt(data3[0][0], 0, 64)
	paunits, err := strconv.ParseInt(data2[0][1], 0, 64)

	TD := prunits - paunits
	TotalDel := strconv.FormatInt(TD, 10)

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "pedningassignments - ", TotalDel)

	c.Data["pedningassignments"] = TotalDel

	productcount, err := services.GetProductCount("ACTIVE")
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Product count fetch Failed")
		return
	}

	procount, err := strconv.ParseUint(productcount, 10, 32)
	produccount := int(procount)
	c.Data["productcount"] = produccount

	clustercount, err := services.GetClusterCount("ACTIVE")
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Cluster count fetch Failed")
		return
	}

	clrcount, err := strconv.ParseUint(clustercount, 10, 32)
	clustrcount := int(clrcount)
	c.Data["clustercount"] = clustrcount

	customerordercount, err := services.GetCustomerOderCount("ACTIVE")
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Cluster count fetch Failed")
		return
	}

	cordercount, err := strconv.ParseUint(customerordercount, 10, 32)
	cordertrcount := int(cordercount)
	c.Data["customerorder"] = cordertrcount

	supplyordercount, err := services.GetSupplyOderCount("ACTIVE")
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("SupplyOderCount count fetch Failed")
		return
	}

	sordercount, err := strconv.ParseUint(supplyordercount, 10, 32)
	sordertrcount := int(sordercount)
	c.Data["supplyordercount"] = sordertrcount

	row6, err := db.Db.Query(`select SUM(total_selling_price) FROM CustomerOrderTxn as ASum where order_purchagetype="KRISHIVALA"`)

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

	row7, err := db.Db.Query(`select SUM(total_buy_price) FROM SupplyOrderTxn as DSum`)

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

	totalcustomertransactionamount, err := strconv.ParseFloat(data6[0][0], 64)
	totalsupplytransactionamount, err := strconv.ParseFloat(data7[0][0], 64)
	totalcollectionagenttransactionamount, err := strconv.ParseFloat(data8[0][0], 64)

	c.Data["totalcustomertransactionamount"] = strconv.FormatFloat(totalcustomertransactionamount, 'f', 3, 64)
	c.Data["totalsupplytransactionamount"] = strconv.FormatFloat(totalsupplytransactionamount, 'f', 3, 64)
	c.Data["totalcollectionagenttransactionamount"] = strconv.FormatFloat(totalcollectionagenttransactionamount, 'f', 3, 64)
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
 from Inventory ORDER BY product_sku`)

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

	//---------------Producer Nodification--------------

	data10, err := services.GetProducers1()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Producer fetch Failed with Status DEMAND_ORDER_CREATED")
		return
	}

	data51, err := services.GetNotificationDuration(data10)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Time Duration for Notification fetch Failed")
		return
	}

	var Dis1 Display1
	for i := range data10 {
		var d Field1
		d.ProducerMobile = data10[i][0]
		d.Status = data10[i][1]
		d.ProductName = data10[i][2]
		d.ClusterName = data10[i][3]
		d.Productqty = data10[i][4]
		d.PPackageQty = data10[i][5] + data10[i][4]
		d.Supplyordertxn_id = data10[i][6]
		d.Createdate = data10[i][7]
		d.ProductSupplyUnits = data10[i][8]
		d.NotificationTime = data51[i]

		Dis1.Fields1 = append(Dis1.Fields1, d)
	}
	c.Data["Dis1"] = Dis1
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Notification Dashboard Data - ", Dis1)

	// if len(data10) <= 0 {
	// 	err = errors.New("Time Duration for Notification fetch Failed")
	// 	return
	// }
	// var Dis2 Display2
	// for i := range data51 {
	// 	var d Field2
	// 	d.NotificationTime = data51[i]

	// 	Dis2.Fields2 = append(Dis2.Fields2, d)
	// }
	// c.Data["Dis2"] = Dis2
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis1)
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", " Notifications - ", data51, "Data len - ", len(data51))

	row11, err := db.Db.Query(`SELECT status,count(producer_mob)as PSum FROM SupplyOrderTxn where status="DEMAND_ORDER_CREATED"`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row11)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row11)
	_, data11, err := sql.Scan(row11)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data11, "\nData len - ", len(data11))
	if len(data11) <= 0 {
		err = errors.New("Unable to get dashboard data")
		return
	}

	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "PMob - ", data8[0][0])
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "PStatus - ", data11[0][0])
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "PCount - ", data11[0][1])

	//c.Data["producermobile"] = data8[0][0]
	c.Data["status"] = data11[0][0]
	c.Data["totalproducercount"] = data11[0][1]

	//---------------Customer Nodification--------------

	data12, err := services.GetCustomerOrderDetailsForNotification()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Customer Order fetch Failed with Status ORDER_RECEIVED")
		return
	}

	data52, err := services.GetCustomerOrderNotificationDuration(data12)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Time Duration for Notification fetch Failed")
		return
	}

	var Dis2 Display2
	for i := range data12 {
		var d Field2
		d.CustomerMobile = data12[i][0]
		d.ClusterName = data12[i][1]
		d.ProductName = data12[i][2]
		d.PPackageQty = data12[i][3] + data12[i][4]
		d.Productqty = data12[i][4]
		d.OrderUnits = data12[i][5]
		d.OrderId = data12[i][6]
		d.Status = data12[i][7]
		d.Createdate = data12[i][8]
		d.ClusterType = data12[i][9]
		d.NotificationTime = data52[i]

		Dis2.Fields2 = append(Dis2.Fields2, d)
	}
	c.Data["Dis2"] = Dis2
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Customer Notification Dashboard Data - ", Dis2)

	row13, err := db.Db.Query(`SELECT status,count(customer_mob)as CSum FROM CustomerOrderTxn where status="ORDER_RECEIVED"`)

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

	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "PMob - ", data8[0][0])
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "CStatus - ", data13[0][0])
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "CCount - ", data13[0][1])

	//c.Data["producermobile"] = data8[0][0]
	c.Data["status"] = data13[0][0]
	c.Data["totalcustomercount"] = data13[0][1]

	return
}
