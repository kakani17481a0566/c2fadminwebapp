package analytics

import (
	"c2f_admin_webapp/model/db"
	"strings"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"
	"c2f_admin_webapp/services"
	"c2f_admin_webapp/session"
	"errors"

	//"fmt"
	//"html/template"

	"runtime/debug"

	//"time"

	"strconv"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"

	"github.com/astaxie/beego"
)

//	type ChartController struct {
//		beego.Controller
//	}
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

type Analytics struct {
	// beego.Controller
	middleware.BaseController
}

func (c *Analytics) Get() {
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
			c.TplName = "analytics/analytics.html"
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

	auth, err := utils.IsAuthorized(utype, "Analytics")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	//row, err := db.Db.Query(`SELECT count(*)as acount,(SELECT count(*) FROM User where role = '1') as mcount FROM User where role = '2'`)

	// c.Data["DateRange"] = c.Input().Get("daterange")
	// tmp := strings.Split(c.Input().Get("daterange"), " - ")

	// from := ""
	// to := ""
	// if len(tmp) == 2 {
	// 	from = tmp[0]
	// 	to = tmp[1]
	// }
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "From", from)
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "To", to)

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

	row1, err := db.Db.Query(`SELECT count(*)as acount,(SELECT count(*) FROM Customer where (user_type = 'PRODUCER' OR user_type='PRODUCER+CONSUMER')) as mcount FROM Customer where (user_type = 'CONSUMER' OR user_type='PRODUCER+CONSUMER')`)
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

	c.Data["declinedcount"] = finalIntNum1
	c.Data["pendingcount"] = finalIntNum3
	c.Data["orderedcount"] = finalIntNum6

	row3, err := db.Db.Query(`SELECT count(*)as acount FROM CustomerOrderTxn WHERE status="ORDER_RECEIVED"`)

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

	prunits, err := strconv.ParseInt(data3[0][0], 0, 64)
	paunits, err := strconv.ParseInt(data2[0][1], 0, 64)

	TD := prunits - paunits
	TotalDel := strconv.FormatInt(TD, 10)

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "pedningassignments - ", TotalDel)

	c.Data["pedningassignments"] = TotalDel

	row12, err := db.Db.Query(`SELECT count(*)as acount FROM CustomerOrderTxn`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row12)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row12)
	_, data12, err := sql.Scan(row12)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data12, "\nData len - ", len(data12))
	if len(data12) <= 0 {
		err = errors.New("Unable to get dashboard data")
		return
	}

	number12, err := strconv.ParseUint(data12[0][0], 10, 32)
	finalIntNum12 := int(number12)

	c.Data["customerorder"] = finalIntNum12

	row20, err := db.Db.Query(`SELECT count(*)as acount FROM SupplyOrderTxn`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row20)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row20)
	_, data10, err := sql.Scan(row20)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data10, "\nData len - ", len(data10))
	if len(data10) <= 0 {
		err = errors.New("Unable to get dashboard data")
		return
	}

	number10, err := strconv.ParseUint(data10[0][0], 10, 32)
	finalIntNum10 := int(number10)

	c.Data["supplyordercount"] = finalIntNum10

	row21, err := db.Db.Query(`SELECT count(*)as acount FROM Clusters `)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row21)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row21)
	_, data11, err := sql.Scan(row21)
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

	number11, err := strconv.ParseUint(data11[0][0], 10, 32)
	finalIntNum11 := int(number11)

	c.Data["clustercount"] = finalIntNum11

	row22, err := db.Db.Query(`SELECT count(*)as acount FROM Lookup_products `)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row22)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row22)
	_, data16, err := sql.Scan(row22)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data16, "\nData len - ", len(data16))
	if len(data11) <= 0 {
		err = errors.New("Unable to get dashboard data")
		return
	}

	number16, err := strconv.ParseUint(data16[0][0], 10, 32)
	finalIntNum16 := int(number16)

	c.Data["productcount"] = finalIntNum16

	// productcount, err := services.GetProductCount("ACTIVE")
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Product count fetch Failed")
	// 	return
	// }

	// procount, err := strconv.ParseUint(productcount, 10, 32)
	// produccount := int(procount)
	// c.Data["productcount"] = produccount

	// clustercount, err := services.GetClusterCount("ACTIVE")
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Cluster count fetch Failed")
	// 	return
	// }

	// clrcount, err := strconv.ParseUint(clustercount, 10, 32)
	// clustrcount := int(clrcount)
	// c.Data["clustercount"] = clustrcount

	// customerordercount, err := services.GetCustomerOderCount("ACTIVE")
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Cluster count fetch Failed")
	// 	return
	// }

	// cordercount, err := strconv.ParseUint(customerordercount, 10, 32)
	// cordertrcount := int(cordercount)
	// c.Data["customerorder"] = cordertrcount

	// supplyordercount, err := services.GetSupplyOderCount("ACTIVE")
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("SupplyOderCount count fetch Failed")
	// 	return
	// }

	// sordercount, err := strconv.ParseUint(supplyordercount, 10, 32)
	// sordertrcount := int(sordercount)
	// c.Data["supplyordercount"] = sordertrcount

	row6, err := db.Db.Query(`select SUM(total_selling_price) FROM CustomerOrderTxn as ASum`)

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

	row7, err := db.Db.Query(`select SUM(updated_total_supply_price) FROM SupplyOrderTxn as DSum`)

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

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "totalcustomertransactionamount - ", data6[0][0])

	row9, err := db.Db.Query(`select product_sku,
 product_name, 
	inventory_received,
   inventory_delivered,
inventory_available,
product_package_qty,
product_qty,
cluster_name
 from Inventory`)

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

func (c *Analytics) Post() {
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
			c.TplName = "analytics/analytics.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = " "
			c.TplName = "analytics/analytics.html"
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

	language := sess.Get("language").(string)
	c.Data["language"] = language
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "language", language)

	auth, err := utils.IsAuthorized(utype, "Analytics")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	// input_name := c.Input().Get("input_name")
	// input_email := c.Input().Get("input_email")
	// input_status := c.Input().Get("input_status")

	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_name - ", input_name)
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_email - ", input_email)
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_status - ", input_status)

	c.Data["DateRange"] = c.Input().Get("daterange1")
	tmp1 := strings.Split(c.Input().Get("daterange1"), " - ")

	from1 := ""
	to1 := ""
	if len(tmp1) == 2 {
		from1 = tmp1[0]
		to1 = tmp1[1]
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "System User Details From date", from1)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "System User Details To date", to1)

	//System User Details

	row10, err := db.Db.Query(`SELECT COUNT(IF(role="1",id,NULL)),COUNT(IF(role="2",id,NULL)),COUNT(IF(role="3",id,NULL)),COUNT(IF(role="4",id,NULL)),COUNT(IF(role="5",id,NULL))FROM User WHERE (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y'))`, from1, to1)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row10)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row10)
	_, data, err := sql.Scan(row10)
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

	c.Data["DateRange"] = c.Input().Get("daterange5")
	tmp2 := strings.Split(c.Input().Get("daterange5"), " - ")

	from2 := ""
	to2 := ""
	if len(tmp2) == 2 {
		from2 = tmp2[0]
		to2 = tmp2[1]
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "System User Details From date", from2)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "System User Details To date", to2)

	//User Details

	row11, err := db.Db.Query(`SELECT COUNT(IF(user_type="PRODUCER" OR user_type="PRODUCER+CONSUMER",id,NULL)), COUNT(IF(user_type="CONSUMER" OR user_type="PRODUCER+CONSUMER",id,NULL)) FROM Customer WHERE (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y'))`, from2, to2)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row11)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row11)
	_, data1, err := sql.Scan(row11)
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

	row4, err := db.Db.Query(`SELECT count(*)as acount FROM DeliveryUser where status = 'ACTIVE' AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y'))`, from2, to2)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row11)
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

	row5, err := db.Db.Query(`SELECT COUNT(IF(user_type="COLLECTIONAGENT",id,NULL)) FROM Collection_Agent where status = 'ACTIVE' AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y'))`, from2, to2)

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

	c.Data["DateRange"] = c.Input().Get("daterange3")
	tmp4 := strings.Split(c.Input().Get("daterange3"), " - ")

	from4 := ""
	to4 := ""
	if len(tmp4) == 2 {
		from4 = tmp4[0]
		to4 = tmp4[1]
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "System User Details From date", from4)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "System User Details To date", to4)

	//Delivery Orser Details

	row2, err := db.Db.Query(`SELECT COUNT(IF(status="DELIVERY_ASSIGNED",id,NULL)),COUNT(IF(status="DELIVERY_DECLINED",id,NULL)),COUNT(IF(status="DELIVERY_COMPLETED",id,NULL)) FROM Oder2Delivery WHERE (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y'))`, from4, to4)

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

	c.Data["DateRange"] = c.Input().Get("daterange4")
	tmp15 := strings.Split(c.Input().Get("daterange4"), " - ")

	from15 := ""
	to15 := ""
	if len(tmp15) == 2 {
		from15 = tmp15[0]
		to15 = tmp15[1]
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "System User Details From date", from15)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "System User Details To date", to15)

	row26, err := db.Db.Query(`SELECT count(*)as acount FROM CustomerOrderTxn WHERE status="ORDER_RECEIVED" AND (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y'))`, from4, to4)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row26)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row26)
	_, data26, err := sql.Scan(row26)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data26, "\nData len - ", len(data26))
	if len(data26) <= 0 {
		err = errors.New("Unable to get dashboard data")
		return
	}

	number5, err := strconv.ParseUint(data26[0][0], 10, 32)
	finalIntNum5 := int(number5)

	c.Data["recivedcount"] = finalIntNum5

	supplyordercount, err := services.GetSupplyOderCount("ACTIVE")
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("SupplyOderCount count fetch Failed")
		return
	}

	sordercount, err := strconv.ParseUint(supplyordercount, 10, 32)
	sordertrcount := int(sordercount)
	c.Data["supplyordercount"] = sordertrcount

	c.Data["DateRange"] = c.Input().Get("daterange2")
	tmp10 := strings.Split(c.Input().Get("daterange2"), " - ")

	from10 := ""
	to10 := ""
	if len(tmp10) == 2 {
		from10 = tmp10[0]
		to10 = tmp10[1]
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "System User Details From date", from10)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "System User Details To date", to10)

	//Transactions Details

	row61, err := db.Db.Query(`select SUM(total_selling_price) FROM CustomerOrderTxn as ASum WHERE (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y'))`, from10, to10)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row61)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row61)
	_, data6, err := sql.Scan(row61)
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

	row7, err := db.Db.Query(`select SUM(updated_total_supply_price) FROM SupplyOrderTxn as DSum WHERE (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y'))`, from10, to10)

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

	row8, err := db.Db.Query(`select SUM(collected_amount) FROM CollectionOrderTxn as CSum WHERE (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y'))`, from10, to10)

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

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "totalcustomertransactionamount - ", data6[0][0])

	//System Details

	c.Data["DateRange"] = c.Input().Get("daterange4")
	tmp20 := strings.Split(c.Input().Get("daterange4"), " - ")

	from20 := ""
	to20 := ""
	if len(tmp20) == 2 {
		from20 = tmp20[0]
		to20 = tmp20[1]
	}

	row3, err := db.Db.Query(`SELECT count(*)as acount FROM CustomerOrderTxn WHERE (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y'))`, from20, to20)

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

	number12, err := strconv.ParseUint(data3[0][0], 10, 32)
	finalIntNum12 := int(number12)

	c.Data["customerorder"] = finalIntNum12

	row20, err := db.Db.Query(`SELECT count(*)as acount FROM SupplyOrderTxn WHERE (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y'))`, from20, to20)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row20)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row20)
	_, data10, err := sql.Scan(row20)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data10, "\nData len - ", len(data10))
	if len(data10) <= 0 {
		err = errors.New("Unable to get dashboard data")
		return
	}

	number10, err := strconv.ParseUint(data10[0][0], 10, 32)
	finalIntNum10 := int(number10)

	c.Data["supplyordercount"] = finalIntNum10

	row21, err := db.Db.Query(`SELECT count(*)as acount FROM Clusters WHERE (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y'))`, from20, to20)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row21)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row21)
	_, data11, err := sql.Scan(row21)
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

	number11, err := strconv.ParseUint(data11[0][0], 10, 32)
	finalIntNum11 := int(number11)

	c.Data["clustercount"] = finalIntNum11

	row22, err := db.Db.Query(`SELECT count(*)as acount FROM Lookup_products WHERE (DATE(create_date='') OR DATE (create_date) >= STR_TO_DATE(?, '%m/%d/%Y')) AND (DATE(create_date='') OR DATE (create_date) <= STR_TO_DATE(?, '%m/%d/%Y'))`, from20, to20)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	defer sql.Close(row22)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row22)
	_, data16, err := sql.Scan(row22)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get dashboard data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data16, "\nData len - ", len(data16))
	if len(data11) <= 0 {
		err = errors.New("Unable to get dashboard data")
		return
	}

	number16, err := strconv.ParseUint(data16[0][0], 10, 32)
	finalIntNum16 := int(number16)

	c.Data["productcount"] = finalIntNum16

	row9, err := db.Db.Query(`select product_sku,
 product_name, 
	inventory_received,
   inventory_delivered,
inventory_available,
product_package_qty,
product_qty,
cluster_name
 from Inventory`)

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

		row41, err := db.Db.Query(`SELECT count(*)as acount FROM C2F_NEW.CustomerOrderTxn `)

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
