package orderToSupplyExecutive

import (
	"c2f_admin_webapp/session"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/model/utils"
	"c2f_admin_webapp/services"
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"

	// "time"

	// "strconv"
	"strings"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"
	"c2f_admin_webapp/model/util/log"
)

type ProducerDeatils struct {
	Name     string `json:"name"`
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	City     string `json:"city"`
	Zipcode  string `json:"zipcode"`
	Mobile   string `json:"mobile"`
	Email    string `json:"email"`
	Latlong  string `json:"latlong"`
	Id       string `json:"id"`
}

type Order struct {
	OrderID string
	Msg     interface{}
}

type SupplyOrderdetails struct {
	Supplyorderid      string `json:"supplyorder_id"`
	Category           string `json:"category_name"`
	Subcat             string `json:"subcategory_name"`
	ProductName        string `json:"product_name"`
	ProductSKU         string `json:"product_sku"`
	ProductQty         string `json:"product_qty"`
	ProductPckgQty     string `json:"product_package_qty"`
	ClusterBuyingprice string `json:"cluster_buying_price"`
	SupplyUnits        string `json:"product_supply_units"`
	Status             string `json:"status"`
	ProductId          string `json:"product_id"`
	TotalSupplyPrice   string `json:"total_supply_price"`
	ClusterName        string `json:"cluster"`
}

type Row struct {
	SupplyOrderId     string
	Product_name      string
	Product_sku       string
	Product_qty       string
	Product_pckgqty   string
	ProducerMob       string
	Cluster_name      string
	Cluster_buy_price string
	SupplyUnits       string
	Status            string
	TimeStamp         string
	ProducerName      string
}

type Display struct {
	Fields []Field
}

type Field struct {
	Id      string
	Name    string
	Email   string
	Cluster string
	Mobile  string
}

type OrderToSupplyExecutive struct {
	// beego.Controller
	middleware.BaseController
}

func (c *OrderToSupplyExecutive) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "OrderToSupplyExecutive Start")
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
			c.TplName = "supplyorder/ordertosupplyexecutive/ordertosupplyexecutive.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "OrderToSupplyExecutive Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "supplyorder/ordertosupplyexecutive/ordertosupplyexecutive.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "OrderToSupplyExecutive  Page Success")
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

	auth, err := utils.IsAuthorized(utype, "OrdertoSupplyExe")
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
			err = errors.New("Cluster user  Not Found")
			return
		}

		data, err := services.GetClustermappedSupplyExecutive(data1[0][1])
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("SupplyExecutive fetch Failed")
			return
		}
		if len(data) <= 0 {
			err = errors.New("SupplyExecutive Not Found")
			return
		}
		var Dis Display
		for i := range data {
			var d Field
			d.Id = data[i][0]
			d.Name = data[i][2]
			d.Mobile = data[i][1]
			d.Cluster = data[i][7]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

		auth, err := utils.IsAuthorized(utype, "OrdertoSupplyExe")
		if !auth {
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
			Autherr = errors.New("UnAuthorized")
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

		row, err := db.Db.Query(` select 
supplyordertxn_id,
product_name, 
product_sku,
product_qty,
product_package_qty,
producer_mob,
cluster_name,
cluster_buying_price,
product_supply_units,
status,
create_date
 from SupplyOrderTxn  where status= "DELIVER_BY_SUPPLY_EXECUTIVE" AND cluster_name=?`, data1[0][3])

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch order")
			return
		}
		defer sql.Close(row)
		_, data4, err := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch order")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", data4)

		var result []Row

		for i := range data4 {
			var r Row
			r.SupplyOrderId = data4[i][0]
			r.Product_name = data4[i][1]
			r.Product_sku = data4[i][2]
			r.Product_qty = data4[i][3]
			r.Product_pckgqty = data4[i][4]
			r.ProducerMob = data4[i][5]
			r.Cluster_name = data4[i][6]
			r.Cluster_buy_price = data4[i][7]
			r.SupplyUnits = data4[i][8]
			r.Status = data4[i][9]
			r.TimeStamp = data4[i][10]

			result = append(result, r)
		}

		c.Data["CustomerData"] = result

	} else if utype == "CLUSTEREXECUTIVE" {

		// CLUSTER-EXECUTIVE CODE

	} else {

		// HOADMIN CODE

		data, err := services.GetSupplyExecutive()
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("SupplyExecutive fetch Failed")
			return
		}
		if len(data) <= 0 {
			err = errors.New("SupplyExecutive  Not Found")
			return
		}
		var Dis Display
		for i := range data {
			var d Field
			d.Id = data[i][0]
			d.Name = data[i][2]
			d.Mobile = data[i][1]
			d.Cluster = data[i][7]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

		auth, err := utils.IsAuthorized(utype, "OrdertoSupplyExe")
		if !auth {
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
			Autherr = errors.New("UnAuthorized")
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

		row, err := db.Db.Query(` select 
supplyordertxn_id,
product_name, 
product_sku,
product_qty,
product_package_qty,
producer_mob,
cluster_name,
cluster_buying_price,
product_supply_units,
status,
create_date
 from SupplyOrderTxn  where status= "DELIVER_BY_SUPPLY_EXECUTIVE"`)

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch supply order")
			return
		}
		defer sql.Close(row)
		_, data4, err := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch supply order")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", data4)
		var result []Row

		for i := range data4 {
			var r Row
			r.SupplyOrderId = data4[i][0]
			r.Product_name = data4[i][1]
			r.Product_sku = data4[i][2]
			r.Product_qty = data4[i][3]
			r.Product_pckgqty = data4[i][4]
			r.ProducerMob = data4[i][5]
			r.Cluster_name = data4[i][6]
			r.Cluster_buy_price = data4[i][7]
			r.SupplyUnits = data4[i][8]
			r.Status = data4[i][9]
			r.TimeStamp = data4[i][10]

			result = append(result, r)
		}

		c.Data["CustomerData"] = result

	}

	return
}

func (c *OrderToSupplyExecutive) Post() {
	msg11 := ""
	log.Println(beego.AppConfig.String("loglevel"), "Info", "OrderToSupplyExecutive page")
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
			c.TplName = "supplyorder/ordertosupplyexecutive/ordertosupplyexecutive.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "OrderToSupplyExecutive Page Fail")
		} else {
			c.Data["DisplayMessage"] = msg11
			c.TplName = "supplyorder/ordertosupplyexecutive/ordertosupplyexecutive.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "OrderToSupplyExecutive  Page Success")
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
	defer sess.SessionRelease(c.Ctx.ResponseWriter)

	utype := sess.Get("user_type").(string)
	c.Data["utype"] = utype

	language := sess.Get("language").(string)
	c.Data["language"] = language
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "language", language)

	clustername := sess.Get("clustername").(string)
	clustername1 := strings.ToUpper(clustername)
	c.Data["clustername"] = clustername1
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "clustername", clustername1)

	uid := sess.Get("uid").(string)
	c.Data["id"] = uid
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "id", uid)

	auth, err := utils.IsAuthorized(utype, "OrdertoSupplyExe")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	input_supplyexecutive := c.Input().Get("input_supplyexecutive")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_supplyexecutive - ", input_supplyexecutive)

	input_supplyExecCharges := c.Input().Get("input_supplyExecCharges")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_supplyExecCharges - ", input_supplyExecCharges)

	if input_supplyexecutive == "" || input_supplyExecCharges == "" {
		err = errors.New("  Supply Executive and Charges cannot be empty")
		return
	}

	row3, err := db.Db.Query(`SELECT 
SupplyExecutive.id,
SupplyExecutive.mob_no,
SupplyExecutive.first_name,
SupplyExecutive.email,
Cluster2SupplyExecutive.cluster_id,
Cluster2SupplyExecutive.supplyexecutive_id,
Clusters.id,
Clusters.cluster_name
from Cluster2SupplyExecutive 
left join SupplyExecutive
on SupplyExecutive.id=Cluster2SupplyExecutive.supplyexecutive_id
left join Clusters
on Cluster2SupplyExecutive.cluster_id=Clusters.id
where SupplyExecutive.mob_no=?`, input_supplyexecutive)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch supply order")
		return
	}
	defer sql.Close(row3)
	_, data6, err := sql.Scan(row3)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch supply order")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data6)

	supply_executive_cluster := data6[0][7]

	order_id := c.Input().Get("txnArr")

	if order_id == "" {
		err = errors.New(" Please Select Check Box ")
		return
	}

	fruits := strings.Split(order_id, ",")

	fmt.Println(fruits)
	lg := len(fruits)
	fmt.Println(lg)

	for i := range fruits {

		h := fruits[i]

		fmt.Println("Supply Order id :", h)

		row, err := db.Db.Query(`select 
		supplyordertxn_id,
		category_name,
		subcategory_name,
		product_name, 
		product_sku,
		product_qty,
		product_package_qty,
		cluster_buying_price,
		product_supply_units,
		status,
		producer_mob,
		create_date,
		product_id,
		product_supply_total_price,
		cluster_name
		from SupplyOrderTxn  where supplyordertxn_id=?`, h)

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch supply order")
			return
		}
		defer sql.Close(row)
		_, data4, err := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch supply order")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", data4)

		cluster_of_order := data4[0][14]

		if supply_executive_cluster == cluster_of_order {
			fmt.Println("Both Cluster are same - Passed")
		} else {
			fmt.Println("Both Cluster are different - Failed")
			err = errors.New("Incorrect cluster mapping for order with ID : " + h)
			msg11 = "Incorrect cluster mapping for order with ID :" + h
			return
		}

		supplyorderdetails := SupplyOrderdetails{

			Supplyorderid:      data4[0][0],
			Category:           data4[0][1],
			Subcat:             data4[0][2],
			ProductName:        data4[0][3],
			ProductSKU:         data4[0][4],
			ProductQty:         data4[0][5],
			ProductPckgQty:     data4[0][6],
			ClusterBuyingprice: data4[0][7],
			SupplyUnits:        data4[0][8],
			ProductId:          data4[0][12],
			TotalSupplyPrice:   data4[0][13],
			ClusterName:        data4[0][14],
		}

		var jsonData1 []byte
		jsonData1, err = json.Marshal(supplyorderdetails)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "jsonData1 supplyorderdetails- ", string(jsonData1))

		row2, err := db.Db.Query(`select first_name,address1,address2,city,zipcode,email,mob_no,location_info,id from Customer where mob_no = ?`, data4[0][10])

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch Producer order")
			return
		}
		defer sql.Close(row2)
		_, data5, err := sql.Scan(row2)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch Producer order")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", data5)

		producerdetails := ProducerDeatils{
			Name:     data5[0][0],
			Address1: data5[0][1],
			Address2: data5[0][2],
			City:     data5[0][3],
			Zipcode:  data5[0][4],
			Mobile:   data5[0][6],
			Email:    data5[0][5],
			Latlong:  data5[0][7],
			Id:       data5[0][8],
		}

		var jsonData2 []byte
		jsonData2, err = json.Marshal(producerdetails)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "jsonData1 supplyorderdetails- ", string(jsonData2))

		result, err := db.Db.Exec(`INSERT INTO Oder2SupplyExecutive (
		supply_order_id,
		supply_order_date,
		supplyexecutive_mob,
		producer_details,
		supply_order_deatils,
		status,
		remarks,
		supplyexecutive_id,
		supplyexecutive_cluster,
		created_by,
		supply_executive_charges,
		create_date)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,now())`,
			h,
			data4[0][11],
			input_supplyexecutive,
			jsonData2,
			jsonData1,
			"DEMAND_ORDER_ISSUED_TO_SUPPLIER",
			"",
			data6[0][0],
			data6[0][7],
			uid,
			input_supplyExecCharges)
		if err != nil {
			err = errors.New("Supply Order to Executive Assignment failed")
			msg11 = "Supply Order to Executive Assignment failed"
			return
		}

		i, err := result.RowsAffected()
		if err != nil || i == 0 {
			err = errors.New("Supply Order to Executive Assignment failed")
			msg11 = "Supply Order to Executive Assignment failed"
			return
		}

		res, err := db.Db.Exec(`UPDATE SupplyOrderTxn SET status = "DEMAND_ORDER_ISSUED_TO_SUPPLIER" WHERE supplyordertxn_id = ?`, h)
		if err != nil {
			err = errors.New("SupplyOrderTxn Status updation failed")
			msg11 = "SupplyOrderTxn Status updation failed"
			return
		}

		i1, err := res.RowsAffected()
		if err != nil || i1 == 0 {
			err = errors.New(" SupplyOrderTxn Status updation failed")
			msg11 = "SupplyOrderTxn Status updation failed"
			return
		}

		msg11 = "Supply Order  Assigned successfully"
	}

	return
}
