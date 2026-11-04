package createOrderCollectionagent

import (
	"c2f_admin_webapp/session"
	//	"encoding/json"
	"runtime/debug"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"
	"errors"
	"strings"

	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/services"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

type CreateOrderCollectionagent struct {
	// beego.Controller
	middleware.BaseController
}

type Display struct {
	Fields []Field
}

type Display1 struct {
	Fields1 []Field1
}
type Field struct {
	Id    string
	Name  string
	Email string
}

type Field1 struct {
	Id   string
	Name string
}

type Users struct {
	Deliveryuser string        `json:"Deliveryuser,omitempty"`
	Users        []OderDeatils `json:"OderDeatils"`
}

type OderDeatils struct {
	Orderid           string          `json:"Orderid,omitempty"`
	OderDate          string          `json:"OderDate,omitempty"`
	DeliveryFrom      string          `json:"DeliveryFrom,omitempty"`
	DeliveryTo        string          `json:"DeliveryTo,omitempty"`
	DeliveryTime      string          `json:"DeliveryTime,omitempty"`
	OrderType         string          `json:"OrderType,omitempty"`
	RepeatOrder       string          `json:"RepeatOrder,omitempty"`
	PaymentInfo       string          `json:"PaymentInfo,omitempty"`
	CategoryName      string          `json:"CategoryName,omitempty"`
	SubcategroryName  string          `json:"SubcategroryName,omitempty"`
	ProductName       string          `json:"ProductName,omitempty"`
	ProductQty        string          `json:"ProductQty,omitempty"`
	ProductPackageQty string          `json:"ProductPackageQty,omitempty"`
	ProductUnits      string          `json:"ProductUnits,omitempty"`
	ProductPrice      string          `json:"ProductPrice,omitempty"`
	TotalPrice        string          `json:"TotalPrice,omitempty"`
	TotalSellPrice    string          `json:"TotalSellPrice,omitempty"`
	ClusterName       string          `json:"ClusterName,omitempty"`
	Name              string          `json:"Name,omitempty"`
	Address1          string          `json:"Address1,omitempty"`
	Address2          string          `json:"Address2,omitempty"`
	City              string          `json:"City,omitempty"`
	Zipcode           string          `json:"Zipcode,omitempty"`
	Mobile            string          `json:"Mobile,omitempty"`
	Email             string          `json:"Email,omitempty"`
	CustomerDeatils   CustomerDeatils `json:"CustomerDeatils"`
}

type CustomerDeatils struct {
	Name     string `json:"Name,omitempty"`
	Address1 string `json:"Address1,omitempty"`
	Address2 string `json:"Address2,omitempty"`
	City     string `json:"City,omitempty"`
	Zipcode  string `json:"Zipcode,omitempty"`
	Mobile   string `json:"Mobile,omitempty"`
	Email    string `json:"Email,omitempty"`
}

func (c *CreateOrderCollectionagent) Get() {
	AdminId := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Add Assets Start")
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
			c.TplName = "ordercollectioagent/createOrderCollectioagent/createOrderCollectioagent.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "ordercollectioagent/createOrderCollectioagent/createOrderCollectioagent.html"
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

	language := sess.Get("language").(string)
	c.Data["language"] = language
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "language", language)

	clustername := sess.Get("clustername").(string)
	clustername1 := strings.ToUpper(clustername)
	c.Data["clustername"] = clustername1
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "clustername", clustername1)

	id := sess.Get("uid").(string)
	c.Data["id"] = id

	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "utype", utype)

	data, err := services.GetDeliveryuser()
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
		d.Name = data[i][1]
		Dis.Fields = append(Dis.Fields, d)
	}
	c.Data["Dis"] = Dis
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

	auth, err := utils.IsAuthorized(utype, "OrderCollectionAgent")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	row, err := db.Db.Query(`select id,
 order_date, 
	order_id,
   cluster_name,
product_name,
status
 from CustomerOrderTxn  where id=?`, AdminId)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	defer sql.Close(row)
	_, data3, err := sql.Scan(row)
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

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data3)

	for i := range data3 {
		c.Data["Orderid"] = data3[i][2]

	}

	return
}
func (c *CreateOrderCollectionagent) Post() {

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
			c.TplName = "ordercollectioagent/createOrderCollectioagent/createOrderCollectioagent.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = "User Oder Mapped Successfully"
			c.TplName = "ordercollectioagent/createOrderCollectioagent/createOrderCollectioagent.html"
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

	language := sess.Get("language").(string)
	c.Data["language"] = language
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "language", language)

	clustername := sess.Get("clustername").(string)
	clustername1 := strings.ToUpper(clustername)
	c.Data["clustername"] = clustername1
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "clustername", clustername1)

	id := sess.Get("uid").(string)
	c.Data["id"] = id
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "id", id)

	auth, err := utils.IsAuthorized(utype, "OrderCollectionAgent")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	// select_deliveruser_name := c.Input().Get("select_deliveruser_name")
	// input_oder_id := c.Input().Get("input_oder_id")

	// err = CheckmapAlreadyExists(select_deliveruser_name, input_oder_id)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("User 2 Oder mapping already Exists")
	// 	return
	// }

	// id, cid, oid, odate, err := services.GetOderid(input_oder_id)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Invalid Order ID ")
	// 	return
	// }

	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "id - ", id)
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "cid - ", cid)
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "oid - ", oid)

	// fname, add1, add2, city, zipcode, email, mob, err := services.GetCustomerinfo(cid)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Invalid Order ID ")
	// 	return
	// }

	// cdeatils := CustomerDeatils{
	// 	Name:     fname,
	// 	Address1: add1,
	// 	Address2: add2,
	// 	City:     city,
	// 	Zipcode:  zipcode,
	// 	Mobile:   mob,
	// 	Email:    email,
	// }

	// var jsonData []byte
	// jsonData, err = json.Marshal(cdeatils)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// }

	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "jsonData - ", string(jsonData))

	// oid, odate, otype, orepeat, dfrom, dto, dtime, cname, subcname, pname, pqty, productpackageqty, pprice, punits, tprice, totalsellprice, clustername,productsku err := services.GetOderinfo(oid)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Invalid Order ID ")
	// 	return
	// }

	// odeatils := OderDeatils{
	// 	Orderid:           oid,
	// 	OderDate:          odate,
	// 	OrderType:         otype,
	// 	RepeatOrder:       orepeat,
	// 	DeliveryFrom:      dfrom,
	// 	DeliveryTo:        dto,
	// 	DeliveryTime:      dtime,
	// 	CategoryName:      cname,
	// 	SubcategroryName:  subcname,
	// 	ProductName:       pname,
	// 	ProductQty:        pqty,
	// 	ProductPackageQty: productpackageqty,
	// 	ProductPrice:      pprice,
	// 	ProductUnits:      punits,
	// 	TotalPrice:        tprice,
	// 	TotalSellPrice:    totalsellprice,
	// 	ClusterName:       clustername,
	// 	Name:              fname,
	// 	Address1:          add1,
	// 	Address2:          add2,
	// 	City:              city,
	// 	Zipcode:           zipcode,
	// 	Mobile:            mob,
	// 	Email:             email,
	// }

	// var jsonData1 []byte
	// jsonData1, err = json.Marshal(odeatils)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// }

	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "jsonData1 - ", string(jsonData1))

	// //	jsondat := &myJSON{Users: []string{"one", "two", "three"}}

	// result, err := db.Db.Exec(`INSERT INTO Oder2Delivery (oder_id,
	// 	order_date,
	// 	deliveryuser_mob,
	// 	customer_details,
	// 	oder_deatils,
	// 	created_by,
	// 	status,
	// 	create_date)
	// 	VALUES (?,?,?,?,?,?,?,now())`,
	// 	oid,
	// 	odate,
	// 	select_deliveruser_name,
	// 	//		jsonData,
	// 	jsonData1,
	// 	id,
	// 	"DELIVERY_ASSIGNED")
	// if err != nil {
	// 	err = errors.New("Cluster User creation failed")
	// 	return
	// }

	// i, err := result.RowsAffected()
	// if err != nil || i == 0 {
	// 	err = errors.New("ClusterProduct creation failed")
	// 	return
	// }

	return
}

func CheckmapAlreadyExists(duser, oid string) (err error) {

	err = nil

	row, err := db.Db.Query(`select count(*) from Oder2Delivery WHERE oder_id = ? AND deliveryuser_mob=? AND status= "DELIVERY_ASSIGNED"`, oid, duser)

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
