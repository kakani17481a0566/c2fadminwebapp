package selectOrderCollectionagent

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

	//	"strconv"
	"strings"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

type DeliveryOrders struct {
	Deliveryuser string     `json:"deliveryuser"`
	Oderinfo     []Oderinfo `json:"oderinfo"`
}
type Oderinfo struct {
	Orderid           string          `json:"orderid"`
	OderDate          string          `json:"oderDate"`
	DeliveryFrom      string          `json:"deliveryFrom"`
	DeliveryTo        string          `json:"deliveryTo"`
	DeliveryTime      string          `json:"deliveryTime"`
	OrderType         string          `json:"orderType"`
	RepeatOrder       string          `json:"repeatOrder"`
	CategoryName      string          `json:"categoryName"`
	SubcategroryName  string          `json:"subcategroryName"`
	ProductName       string          `json:"productName"`
	ProductSKU        string          `json:"productSKU"`
	ProductQty        string          `json:"productQty"`
	ProductPackageQty string          `json:"productPackageQty"`
	ProductUnits      string          `json:"productUnits"`
	ProductPrice      string          `json:"productPrice"`
	TotalPrice        string          `json:"totalPrice"`
	TotalSellPrice    string          `json:"totalSellPrice"`
	ClusterName       string          `json:"clusterName"`
	Name              string          `json:"name"`
	Address1          string          `json:"address1"`
	Address2          string          `json:"address2"`
	City              string          `json:"city"`
	Zipcode           string          `json:"zipcode"`
	Mobile            string          `json:"mobile"`
	Email             string          `json:"email"`
	CustomerDeatils   CustomerDeatils `json:"CustomerDeatils"`
}
type CustomerDeatils struct {
	Name        string `json:"name"`
	Address1    string `json:"address1"`
	Address2    string `json:"address2"`
	City        string `json:"city"`
	Zipcode     string `json:"zipcode"`
	Mobile      string `json:"mobile"`
	Email       string `json:"email"`
	Latlong     string `json:"latlong"`
	LocationURL string `json:"locationURL"`
}

type Row struct {
	Id           string
	TimeStamp    string
	Oder_id      string
	Cluster_name string
	Product_name string
	Status       string
	Name         string
	Address1     string
	Address2     string
	City         string
	Zipcode      string
	Latlong      string
	LocationURL  string
	Email        string
	Mobile       string
	Custdetails  string
	CustomerID   string
}

type Display struct {
	Fields []Field
}

type Field struct {
	Id      string
	Name    string
	Mobile  string
	Cluster string
}

type SelectOrderCollectionagent struct {
	// beego.Controller
	middleware.BaseController
}

func (c *SelectOrderCollectionagent) Get() {
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
			c.TplName = "ordercollectionagent/selectOrderCollectionagent/selectOrderCollectionagent.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "ordercollectionagent/selectOrderCollectionagent/selectOrderCollectionagent.html"
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

	id := sess.Get("uid").(string)
	c.Data["id"] = id

	auth, err := utils.IsAuthorized(utype, "OrderCollectionAgent")
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

		data, err := services.GetClustermappedagent(data1[0][1])
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Delivery user fetch Failed")
			return
		}
		if len(data) <= 0 {
			err = errors.New("Delivery user  Not Found")
			return
		}
		var Dis Display
		for i := range data {
			var d Field
			d.Id = data[i][0]
			d.Name = data[i][2]
			d.Mobile = data[i][3]
			d.Cluster = data[i][7]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

		row, err := db.Db.Query(` select id,
 order_date, 
	order_id,
   cluster_name,
product_name,
status,
customer_id,
agent_status
 from CustomerOrderTxn where status= "DELIVERY_COMPLETED" AND payment_type= "CASH" AND agent_status IS NULL AND cluster_name=?`, data1[0][3])

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
			r.Id = data4[i][0]
			r.TimeStamp = data4[i][1]
			r.Oder_id = data4[i][2]
			r.Cluster_name = data4[i][3]
			r.Product_name = data4[i][4]
			r.Status = data4[i][5]
			r.CustomerID = data4[i][6]

			log.Println(beego.AppConfig.String("loglevel"), "Debug", "data4[i][2]", data4[i][2])

			fname, add1, add2, city, zipcode, email, mob, latlong, err := services.GetCustomerinfo(data4[i][6])
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Invalid Order ID ")
				return
			}

			var tmp = fname + add1 + add2

			log.Println(beego.AppConfig.String("loglevel"), "Debug", "ttmp", tmp)

			var tmp2 = city + zipcode + email
			var tmp3 = mob + latlong
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "tmp2", tmp2)

			log.Println(beego.AppConfig.String("loglevel"), "Debug", "tmp3", tmp3)

			r.Custdetails = tmp + "," + tmp2 + "," + tmp3
			//	r.LocationURL = locationurl

			result = append(result, r)
		}
		c.Data["CustomerData"] = result

	} else if utype == "CLUSTEREXECUTIVE" {

		data1, err1 := services.GetClusterAsPerMappingExec(id)
		if err1 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
			err1 = errors.New("Cluster user fetch Failed")
			return
		}
		if len(data1) <= 0 {
			err = errors.New("Cluster user  Not Found")
			return
		}

		data, err := services.GetClustermappedagent(data1[0][1])
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Delivery user fetch Failed")
			return
		}
		if len(data) <= 0 {
			err = errors.New("Delivery user  Not Found")
			return
		}
		var Dis Display
		for i := range data {
			var d Field
			d.Id = data[i][0]
			d.Name = data[i][2]
			d.Mobile = data[i][3]
			d.Cluster = data[i][7]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

		row, err := db.Db.Query(` select id,
 order_date, 
	order_id,
   cluster_name,
product_name,
status,
customer_id,
agent_status
 from CustomerOrderTxn where status= "DELIVERY_COMPLETED" AND payment_type= "CASH" AND agent_status IS NULL AND cluster_name=?`, data1[0][3])

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
			r.Id = data4[i][0]
			r.TimeStamp = data4[i][1]
			r.Oder_id = data4[i][2]
			r.Cluster_name = data4[i][3]
			r.Product_name = data4[i][4]
			r.Status = data4[i][5]
			r.CustomerID = data4[i][6]

			log.Println(beego.AppConfig.String("loglevel"), "Debug", "data4[i][2]", data4[i][2])

			fname, add1, add2, city, zipcode, email, mob, latlong, err := services.GetCustomerinfo(data4[i][6])
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Invalid Order ID ")
				return
			}

			var tmp = fname + add1 + add2

			log.Println(beego.AppConfig.String("loglevel"), "Debug", "ttmp", tmp)

			var tmp2 = city + zipcode + email
			var tmp3 = mob + latlong
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "tmp2", tmp2)

			log.Println(beego.AppConfig.String("loglevel"), "Debug", "tmp3", tmp3)

			r.Custdetails = tmp + "," + tmp2 + "," + tmp3
			//	r.LocationURL = locationurl

			result = append(result, r)
		}
		c.Data["CustomerData"] = result

	} else {

		data, err := services.GetClustermappedcollectionagent()
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Delivery user fetch Failed")
			return
		}
		if len(data) <= 0 {
			err = errors.New("Delivery user  Not Found")
			return
		}
		var Dis Display
		for i := range data {
			var d Field
			d.Id = data[i][0]
			d.Name = data[i][2]
			d.Mobile = data[i][3]
			d.Cluster = data[i][7]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

		row, err := db.Db.Query(` select id,
 order_date, 
	order_id,
   cluster_name,
product_name,
status,
customer_id,
agent_status
 from CustomerOrderTxn where status= "DELIVERY_COMPLETED" AND payment_type= "CASH"  AND agent_status IS NULL`)

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
			r.Id = data4[i][0]
			r.TimeStamp = data4[i][1]
			r.Oder_id = data4[i][2]
			r.Cluster_name = data4[i][3]
			r.Product_name = data4[i][4]
			r.Status = data4[i][5]
			r.CustomerID = data4[i][6]

			log.Println(beego.AppConfig.String("loglevel"), "Debug", "data4[i][2]", data4[i][2])

			fname, add1, add2, city, zipcode, email, mob, latlong, err := services.GetCustomerinfo(data4[i][6])
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Invalid Order ID ")
				return
			}

			var tmp = fname + add1 + add2

			log.Println(beego.AppConfig.String("loglevel"), "Debug", "ttmp", tmp)

			var tmp2 = city + zipcode + email
			var tmp3 = mob + latlong
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "tmp2", tmp2)

			log.Println(beego.AppConfig.String("loglevel"), "Debug", "tmp3", tmp3)

			r.Custdetails = tmp + "," + tmp2 + "," + tmp3
			//	r.LocationURL = locationurl

			result = append(result, r)
		}
		c.Data["CustomerData"] = result

	}

	return
}

func (c *SelectOrderCollectionagent) Post() {
	msg11 := ""
	fmt.Println("msg11", msg11)
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
			c.TplName = "ordercollectionagent/selectOrderCollectionagent/selectOrderCollectionagent.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Order - CollectionAgent Assignmemt Completed "
			c.TplName = "ordercollectionagent/selectOrderCollectionagent/selectOrderCollectionagent.html"
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

	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1

	uid := sess.Get("uid").(string)
	c.Data["id"] = uid
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "id", uid)

	auth, err := utils.IsAuthorized(utype, "OrderCollectionAgent")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	select_collection_agent := c.Input().Get("select_collection_agent")
	//order_id := c.Input().Get("programming")
	//	select_cluster_name := c.Input().Get("select_cluster_name")

	Get_Collectionagent_Values := c.Input().Get("Get_Collectionagent_Values")
	Get_cluster_name := c.Input().Get("Get_cluster_name")

	// order_id1 := c.Input().Get("txnArr")
	order_id := c.Input().Get("txnArr")
	// fmt.Println("order_id1", order_id1)
	fmt.Println("order_id", order_id)

	if select_collection_agent == "" {
		err = errors.New(" Please Select Collection Agent ")
		return
	}

	// order_id := c.Input().Get("txnArr")

	if order_id == "" {
		err = errors.New(" Please Select Check Box ")
		return
	}

	fmt.Println("Get_Collectionagent_Values", Get_Collectionagent_Values)
	fmt.Println("Get_cluster_name", Get_cluster_name)

	// Split the Get_cluster_name string into a slice of cluster names
	clusterNames := strings.Split(Get_cluster_name, ",")

	// Flag to check if there is any mismatch
	mismatchFound := false

	// Iterate through each cluster name and compare with Get_deliveruser_Values
	for _, cluster := range clusterNames {
		cluster = strings.TrimSpace(cluster) // Trim any whitespace around cluster names
		if cluster != Get_Collectionagent_Values {
			mismatchFound = true
			break
		}
	}

	// Check if any mismatch was found
	if mismatchFound {
		fmt.Println("Both Cluster are different - Failed")
		err = errors.New("Incorrect cluster mapping for order with ID : " + order_id)
		msg11 = "Incorrect cluster mapping for order with ID :" + order_id
		return
	} else {
		fmt.Println("Both Cluster are same - Passed")
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "select_collection_agent - ", select_collection_agent)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "order_id - ", order_id)
	fruits := strings.Split(order_id, ",")

	fmt.Println(fruits)
	lg := len(fruits)
	fmt.Println(lg)

	for i := range fruits {

		h := fruits[i]

		fmt.Println(h)

		// err = CheckMapExists(select_cluster_name, select_deliveruser_name)
		// if err != nil {
		// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		// 	err = errors.New("Product Mapping Not Exists")
		// 	return
		// }

		// err = CheckmapAlreadyExists(select_deliveruser_name, h)
		// if err != nil {
		// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		// 	err = errors.New("User 2 Oder mapping already Exists")
		// 	return
		// }

		// status, err := Checkorderstatus(select_deliveruser_name, h)
		// if err != nil {
		// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		// 	err = errors.New("Order already mapped to diferent user")
		// 	return
		// }
		// if status == "DELIVERY_ASSIGNED" {
		// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		// 	err = errors.New("Order already mapped to diferent user")
		// 	return
		// }
		id, cid, oid, odate, err := services.GetOderid(h)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Invalid Order ID ")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "id - ", id)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "cid - ", cid)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "oid - ", oid)

		fname, add1, add2, city, zip, email, mob, latlong, err := services.GetCustomerinfo(cid)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Invalid Order ID ")
			return
		}

		cdeatils := CustomerDeatils{
			Name:     fname,
			Address1: add1,
			Address2: add2,
			City:     city,
			Zipcode:  zip,
			Mobile:   mob,
			Email:    email,
			Latlong:  latlong,
			//LocationURL: locationurl,
		}

		var jsonData []byte
		jsonData, err = json.Marshal(cdeatils)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "jsonData - ", string(jsonData))

		oid, odate, otype, orepeat, dfrom, dto, dtime, cname, subcname, pname, pqty, productpackageqty, pprice, punits, tprice, totalsellprice, clustername, productsku, err := services.GetOderinfo(oid)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Invalid Order ID ")
			return
		}

		odeatils := Oderinfo{
			Orderid:           oid,
			OderDate:          odate,
			OrderType:         otype,
			RepeatOrder:       orepeat,
			DeliveryFrom:      dfrom,
			DeliveryTo:        dto,
			DeliveryTime:      dtime,
			CategoryName:      cname,
			SubcategroryName:  subcname,
			ProductName:       pname,
			ProductQty:        pqty,
			ProductSKU:        productsku,
			ProductPackageQty: productpackageqty,
			ProductPrice:      pprice,
			ProductUnits:      punits,
			TotalPrice:        tprice,
			TotalSellPrice:    totalsellprice,
			ClusterName:       clustername,
			Name:              fname,
			Address1:          add1,
			Address2:          add2,
			City:              city,
			Zipcode:           zip,
			Mobile:            mob,
			Email:             email,
		}

		var jsonData1 []byte
		jsonData1, err = json.Marshal(odeatils)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "jsonData1 - ", string(jsonData1))

		result, err := db.Db.Exec(`INSERT INTO Order2Collectionagent (oder_id,
		order_date,
		collectionagent_mob,
		customer_details,
		oder_deatils,
		created_by,
		status,
		create_date)
		VALUES (?,?,?,?,?,?,?,now())`,
			oid,
			odate,
			select_collection_agent,
			jsonData,
			jsonData1,
			uid,
			"ACTIVE")
		if err != nil {
			err = errors.New("Delivery Order  creation failed")
			return
		}

		i, err := result.RowsAffected()
		if err != nil || i == 0 {
			err = errors.New("Delivery Order creation failed")
			return
		}

		res, err := db.Db.Exec(`UPDATE CustomerOrderTxn SET status = "DELIVERY_COMPLETED",agent_status="ASSIGNED"
	WHERE order_id = ?`, h)
		if err != nil {
			err = errors.New("CustomerOrderTxn Delivery updation failed")
			return
		}

		i1, err := res.RowsAffected()
		if err != nil || i1 == 0 {
			err = errors.New(" CustomerOrderTxn Delivery updation failed")
			return
		}

		// log.Println(beego.AppConfig.String("loglevel"), "Debug", "punits - ", punits)
		// log.Println(beego.AppConfig.String("loglevel"), "Debug", "productsku - ", productsku)
		// log.Println(beego.AppConfig.String("loglevel"), "Debug", "clustername - ", clustername)

		// dunits, aunits, runits, sku, cluster, err := services.Getinventory(productsku, clustername)
		// if err != nil {
		// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		// 	err = errors.New("Inventory Get Failed")
		// 	return
		// }

		// log.Println(beego.AppConfig.String("loglevel"), "Debug", "dunits - ", dunits)
		// log.Println(beego.AppConfig.String("loglevel"), "Debug", "aunits - ", aunits)
		// log.Println(beego.AppConfig.String("loglevel"), "Debug", "runits - ", runits)
		// log.Println(beego.AppConfig.String("loglevel"), "Debug", "sku - ", sku)
		// log.Println(beego.AppConfig.String("loglevel"), "Debug", "cluster - ", cluster)

		// ddunits, err := strconv.ParseInt(dunits, 0, 64)
		// rrunits, err := strconv.ParseInt(runits, 0, 64)
		// ppunits, err := strconv.ParseInt(punits, 0, 64)
		// TD := ddunits + ppunits
		// TotalDel := strconv.FormatInt(TD, 10)
		// TA := rrunits - TD
		// TotalAvl := strconv.FormatInt(TA, 10)
		// log.Println(beego.AppConfig.String("loglevel"), "Debug", "TotalDel - ", TotalDel)
		// log.Println(beego.AppConfig.String("loglevel"), "Debug", "TotalAvl - ", TotalAvl)
		// log.Println(beego.AppConfig.String("loglevel"), "Debug", "productskuuuu - ", productsku)
		// log.Println(beego.AppConfig.String("loglevel"), "Debug", "clusterrrrr - ", cluster)

		// res2, err := db.Db.Exec(`UPDATE Inventory SET inventory_delivered = ? ,inventory_available=?
		//               WHERE product_sku = ? AND cluster_name=?`, TotalDel, TotalAvl, productsku, cluster)
		// if err != nil {
		// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		// 	err = errors.New("Inventory Get Failed")
		// 	return
		// }

		// i2, err := res2.RowsAffected()
		// if err != nil || i2 == 0 {
		// 	err = errors.New(" CustomerOrderTxn Delivery updation failed")
		// 	return
		// }

	}

	return
}

func CheckmapAlreadyExists(duser, oid string) (err error) {

	err = nil

	row, err := db.Db.Query(`select count(*) from Order2Collectionagent WHERE oder_id = ? AND deliveryuser_mob=? AND status= "DELIVERY_COMPLETED"`, oid, duser)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch user")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch user")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	countlen := data[0][0]

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "countlen", countlen)

	if countlen != "0" {
		err = errors.New("Mappig already exists")
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}

	return

}

// func Checkorderstatus(duser, oid string) (status string, err error) {

// 	row, err := db.Db.Query(`select status from Oder2Delivery WHERE oder_id = ? AND deliveryuser_mob=? `, oid, duser)
// 	if err != nil {
// 		beego.Error(err)
// 		err = errors.New("Unable to get the CustomerOrderTxn Id  ")
// 		return
// 	}
// 	defer sql.Close(row)
// 	_, data, err := sql.Scan(row)
// 	if err != nil {
// 		beego.Error(err)
// 		err = errors.New("Message Scan Fail")
// 		return
// 	}
// 	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
// 	if len(data) <= 0 {
// 		err = errors.New("Unable to get the Agent Account number info")
// 		return
// 	}
// 	status = data[0][0]

// 	return
// }
