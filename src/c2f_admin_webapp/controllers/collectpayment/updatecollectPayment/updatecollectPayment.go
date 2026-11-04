package updatecollectPayment

import (
	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/model/utils"
	"c2f_admin_webapp/services"
	"c2f_admin_webapp/session"
	"errors"
	"runtime/debug"
	"strings"

	"strconv"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"

	"github.com/astaxie/beego"
)

type Row struct {
	Oder_id        string
	PaymentStatus  string
	PaymentType    string
	PaymentRemarks string
	Amount         string
	//OrderId        string
	CustomerMobile string
	//CustomerName   string
	CollectionagentMobile string
}
type Display2 struct {
	Fields2 []Field2
}
type Field2 struct {
	Id     string
	Name   string
	Mobile string
}
type UpdatecollectPayment struct {
	// beego.Controller
	middleware.BaseController
}

func (c *UpdatecollectPayment) Get() {
	AdminId := c.Ctx.Input.Param(":AdminID")
	//log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User Start")
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
			c.TplName = "collectpayment/updatecollectPayment/updatecollectPayment.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Collect Payment  Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "collectpayment/updatecollectPayment/updatecollectPayment.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Collect Payment  Page Success")
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

	data2, err := services.Getcollectionagent()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Producer fetch Failed")
		return
	}
	if len(data2) <= 0 {
		err = errors.New("Producer  Not Found")
		return
	}
	var Dis2 Display2
	for i := range data2 {
		var d Field2
		d.Id = data2[i][0]
		d.Name = data2[i][2]
		d.Mobile = data2[i][3]
		Dis2.Fields2 = append(Dis2.Fields2, d)
	}
	c.Data["Dis2"] = Dis2
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis2)
	auth, err := utils.IsAuthorized(utype, "CollectPayment")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	row5, err := db.Db.Query(`select JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.productUnits'))
	from Oder2Delivery
    where oder_id=? AND status="DELIVERY_COMPLETED" limit 1`, AdminId)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get supply transaction data")
		return
	}
	defer sql.Close(row5)
	_, data5, err := sql.Scan(row5)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data5, "\nData len - ", len(data5))
	if len(data5) <= 0 {
		err = errors.New("User data not found")
		return
	}

	if len(data5) <= 0 {
		err = errors.New("Unable to get the error message info")
		return
	}
	u5 := data5[0][0]

	row2, err := db.Db.Query(`select SUM(JSON_EXTRACT(oder_deatils, "$.productPrice")) as oder_deatils
	from Oder2Delivery
    where oder_id=? AND status="DELIVERY_COMPLETED" limit 1`, AdminId)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get supply transaction data")
		return
	}
	defer sql.Close(row2)
	_, data3, err := sql.Scan(row2)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data3, "\nData len - ", len(data3))
	if len(data3) <= 0 {
		err = errors.New("User data not found")
		return
	}

	if len(data3) <= 0 {
		err = errors.New("Unable to get the error message info")
		return
	}

	p5 := data2[0][0]

	u6, err := strconv.ParseFloat(u5, 64)
	p6, err := strconv.ParseFloat(p5, 64)
	totalprice5 := (u6 * p6)
	totalprice6 := strconv.FormatFloat(totalprice5, 'f', 0, 32)

	//c.Data["TotalAmount"] = totalprice6

	c.Data["TotalAmount"] = totalprice6

	row, err := db.Db.Query(`select CustomerOrderTxn.order_id,
        CustomerOrderTxn.payment_status,
        CustomerOrderTxn.payment_type,
        CustomerOrderTxn.total_selling_price,
        CustomerOrderTxn.customer_mob,
        Order2Collectionagent.oder_id,
        Order2Collectionagent.collectionagent_mob,
        CustomerOrderTxn.cluster_supply_charge,
	CustomerOrderTxn.cluster_packaging_charge,
	CustomerOrderTxn.cluster_bank_charge,
	CustomerOrderTxn.cluster_discount,
	CustomerOrderTxn.krishivala_CGST_charge,
	
	CustomerOrderTxn.krishivala_SGST_charge,
	CustomerOrderTxn.krishivala_CESS_charge,
	CustomerOrderTxn.krishivala_OtherTAX_charge,
	CustomerOrderTxn.total_selling_price,
	CustomerOrderTxn.total_days
         FROM CustomerOrderTxn 
LEFT JOIN Order2Collectionagent 
ON CustomerOrderTxn.order_id = Order2Collectionagent.oder_id
        WHERE CustomerOrderTxn.order_id =?  limit 1`, AdminId)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Customer data")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Customer data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Customer data not found")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	c.Data["Oder_id"] = data[0][0]
	c.Data["PaymentStatus"] = data[0][1]
	c.Data["PaymentType"] = data[0][2]
	//c.Data["PaymentRemarks"] = data[0][3]
	c.Data["Amount"] = data[0][3]
	c.Data["CustomerMobile"] = data[0][4]
	c.Data["CollectionagentMobile"] = data[0][6]

	c.Data["Clustersupplycharge"] = data[0][7]
	c.Data["Clusterpackagingcharge"] = data[0][8]
	c.Data["Clusterbankcharge"] = data[0][9]
	c.Data["Clusterdiscountcharge"] = data[0][10]

	c.Data["KrishivalaCGSTcharge"] = data[0][11]
	c.Data["KrishivalaSGSTcharge"] = data[0][12]
	c.Data["KrishivalaCESScharge"] = data[0][13]
	c.Data["KrishivalaOtherTAXcharge"] = data[0][14]
	c.Data["TotalSalePrice"] = data[0][15]
	c.Data["TotalDays"] = data[0][16]

	TTotalAmount := data3[0][0]
	PProductUnits := data5[0][0]

	cluster_supply_charge := data[0][7]
	cluster_packaging_charge := data[0][8]
	cluster_bank_charge := data[0][9]
	cluster_discount := data[0][10]

	krishivala_CGST_charge := data[0][11]
	krishivala_SGST_charge := data[0][12]
	krishivala_CESS_charge := data[0][13]
	krishivala_OtherTAX_charge := data[0][14]

	scluster_supply_charge, err := strconv.ParseFloat(cluster_supply_charge, 64)
	scluster_packaging_charge, err := strconv.ParseFloat(cluster_packaging_charge, 64)
	scluster_bank_charge, err := strconv.ParseFloat(cluster_bank_charge, 64)
	scluster_discount, err := strconv.ParseFloat(cluster_discount, 64)

	TTotalAAmount, err := strconv.ParseFloat(TTotalAmount, 64)
	PProductUUnits, err := strconv.ParseFloat(PProductUnits, 64)

	skrishivala_CGST_charge, err := strconv.ParseFloat(krishivala_CGST_charge, 64)
	skrishivala_SGST_charge, err := strconv.ParseFloat(krishivala_SGST_charge, 64)
	skrishivala_CESS_charge, err := strconv.ParseFloat(krishivala_CESS_charge, 64)
	skrishivala_OtherTAX_charge, err := strconv.ParseFloat(krishivala_OtherTAX_charge, 64)

	charges := scluster_supply_charge + scluster_packaging_charge + scluster_bank_charge
	taxs := skrishivala_CGST_charge + skrishivala_SGST_charge + skrishivala_CESS_charge + skrishivala_OtherTAX_charge

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "charges - ", charges)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "taxs - ", taxs)

	TPP := (TTotalAAmount * PProductUUnits) + (charges) - scluster_discount + (taxs)

	tpp := strconv.FormatFloat(TPP, 'f', 0, 32)

	//	log.Println(beego.AppConfig.String("loglevel"), "Debug", "select_Customer  - ", select_Customer)

	//  TotalBuying_price = cluster.buyingprice  + (cluster.supplycarge + cluster.packagingcharge + cluster.bankcharge ) - cluster.discount  + (producer.gcst + producer.scst + producer.C29 + producer.othertax)
	TotalBillAmount := tpp

	c.Data["TotalBillAmount"] = TotalBillAmount

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "TPP  - ", TPP)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "TotalBillAmount  - ", TotalBillAmount)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "TPP  - ", TPP)

	return
}
func (c *UpdatecollectPayment) Post() {
	AdminId := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Info", "add asset post page")
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
			c.TplName = "collectpayment/updatecollectPayment/updatecollectPayment.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Collect Payment Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Collect Payment Updated Successfully"
			c.TplName = "collectpayment/updatecollectPayment/updatecollectPayment.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Collect Payment  Page Success")
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

	auth, err := utils.IsAuthorized(utype, "CollectPayment")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "adminid - ", AdminId)

	// row, err := db.Db.Query(`select CustomerOrderTxn.order_id,
	//        CustomerOrderTxn.payment_status,
	//        CustomerOrderTxn.payment_type,
	//        CustomerOrderTxn.total_selling_price,
	//        CustomerOrderTxn.customer_mob,
	//        Order2Collectionagent.oder_id,
	//        Order2Collectionagent.collectionagent_mob
	//         FROM CustomerOrderTxn
	// LEFT JOIN Order2Collectionagent
	// ON CustomerOrderTxn.order_id = Order2Collectionagent.oder_id
	//        WHERE CustomerOrderTxn.order_id =?  limit 1`, AdminId)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Unable to get Customer data")
	// 	return
	// }
	// defer sql.Close(row)
	// _, data, err := sql.Scan(row)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Unable to get Customer data")
	// 	return
	// }
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	// if len(data) <= 0 {
	// 	err = errors.New("Customer data not found")
	// 	return
	// }
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", data)
	// for i := range data {
	// 	c.Data["Oder_id"] = data[i][0]
	// 	c.Data["PaymentStatus"] = data[i][1]
	// 	c.Data["PaymentType"] = data[i][2]
	// 	//c.Data["PaymentRemarks"] = data[i][3]
	// 	c.Data["Amount"] = data[i][3]
	// 	c.Data["CustomerMobile"] = data[i][4]
	// 	c.Data["CollectionagentMobile"] = data[i][6]
	// }
	input_payment_status := c.Input().Get("input_payment_status")
	input_payment_type := c.Input().Get("input_payment_type")
	input_payment_remarks := c.Input().Get("input_payment_remarks")
	input_customer_mob := c.Input().Get("input_customer_mob")
	//input_agent_mob := c.Input().Get("input_agent_mob")
	//select_agent_name := c.Input().Get("select_agent_name")
	//input_customer_name := c.Input().Get("input_customer_name")
	input_order_id := c.Input().Get("input_order_id")
	input_collected_amount := c.Input().Get("input_collected_amount")
	select_Collection_Agent := c.Input().Get("select_Collection_Agent")
	// CollectionAgent_id, _ := services.GetCollectionAgentId(select_Collection_Agent)
	// if err != nil {
	//      beego.Error(err)
	//      err = errors.New("Invalid CollectionAgent  ")
	//      return
	// }
	CollectionAgent_name, Cluster_name, agent_id, _ := services.GetCollectionAgentName(select_Collection_Agent)
	if err != nil {
		beego.Error(err)
		err = errors.New("Invalid CollectionAgent  ")
		return
	}
	Customer_name, _ := services.GetCustomerName(input_customer_mob)
	if err != nil {
		beego.Error(err)
		err = errors.New("Invalid Customer  ")
		return
	}
	if input_payment_status == "" || input_payment_remarks == "" || select_Collection_Agent == "" {
		err = errors.New("* fields cannot be empty")
		return
	}
	res, err := db.Db.Exec(`UPDATE CustomerOrderTxn SET 
        payment_status = ?, payment_type = ?, payment_remarks = ?
        WHERE order_id = ?`, input_payment_status, input_payment_type, input_payment_remarks, AdminId)
	if err != nil {
		err = errors.New("Collect Payment updation failed")
		return
	}
	i, err := res.RowsAffected()
	if err != nil || i == 0 {
		err = errors.New("Collect Payment updation failed")
		return
	}
	result1, err := db.Db.Exec(`INSERT INTO CollectionOrderTxn (agent_mobile, 
        agent_name,
        customer_mobile,
        customer_name,
        order_id, 
        collected_amount,
        payment_type,
        payment_status, 
        payment_remarks,
        agent_cluster,
        agent_id,
        create_date)
        VALUES (?,?,?,?,?,?,?,?,?,?,?,now())`,
		select_Collection_Agent,
		CollectionAgent_name,
		input_customer_mob,
		Customer_name,
		input_order_id,
		input_collected_amount,
		input_payment_type,
		input_payment_status,
		input_payment_remarks,
		Cluster_name,
		agent_id)
	if err != nil {
		err = errors.New("User creation failed")
		return
	}
	j, err := result1.RowsAffected()
	if err != nil || j == 0 {
		err = errors.New("User creation failed")
		return
	}
	// row, err := db.Db.Query(`select order_id,payment_status, payment_type,payment_remarks,total_selling_price,customer_mob from CustomerOrderTxn WHERE order_id =?  limit 1`, AdminId)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Unable to get Customer data")
	// 	return
	// }
	// defer sql.Close(row)
	// _, data, err := sql.Scan(row)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Unable to get Customer data")
	// 	return
	// }
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	// if len(data) <= 0 {
	// 	err = errors.New("Customer data not found")
	// 	return
	// }
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", data)
	// var result []Row
	// for i := range data {
	// 	var r Row
	// 	r.Oder_id = data[i][0]
	// 	r.PaymentStatus = data[i][1]
	// 	r.PaymentType = data[i][2]
	// 	r.PaymentRemarks = data[i][3]
	// 	r.Amount = data[i][4]
	// 	r.CustomerMobile = data[i][5]
	// 	result = append(result, r)
	// }
	return
}
