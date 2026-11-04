package selectPurchaseOrder

import (
	"c2f_admin_webapp/session"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/model/utils"
	"c2f_admin_webapp/services"

	//	"encoding/json"
	"errors"
	//"fmt"
	"runtime/debug"
	//	"time"

	//	"strconv"
	"strings"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

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

type Order struct {
	OrderID string
	Msg     interface{}
}

// type Sound struct {
// 	Description string
// 	Authority   string
// }

type Row struct {
	Id                    string
	Cluster_name          string
	Product_name          string
	Inventory_available   string
	Product_sku           string
	Cluster_selling_price string
}

type Display struct {
	Fields []Field
}

type Field struct {
	Id   string
	Name string
}

type SelectPurchaseOrder struct {
	// beego.Controller
	middleware.BaseController
}

func (c *SelectPurchaseOrder) Get() {
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
			c.TplName = "purchaseorder/selectPurchaseOrder/selectPurchaseOrder.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "purchaseorder/selectPurchaseOrder/selectPurchaseOrder.html"
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

	id := sess.Get("uid").(string)
	c.Data["id"] = id

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

	if utype == "CLUSTERMANAGER" {

		data, err := services.GetClusterAsPerMapping(id)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster user fetch Failed")
			return
		}
		if len(data) <= 0 {
			err = errors.New("Cluster user  Not Found")
			return
		}
		var Dis Display
		for i := range data {
			var d Field
			d.Id = data[i][2]
			d.Name = data[i][3]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)
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
		row, err := db.Db.Query(` select id,
   cluster_name,
product_name,
inventory_available,
product_sku,
cluster_selling_price
 from Inventory where (cluster_name='' OR cluster_name like ?)`, cname)

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
			r.Cluster_name = data4[i][1]
			r.Product_name = data4[i][2]
			r.Inventory_available = data4[i][3]
			r.Product_sku = data4[i][4]
			r.Cluster_selling_price = data4[i][5]

			log.Println(beego.AppConfig.String("loglevel"), "Debug", "data4[i][2]", data4[i][2])

			result = append(result, r)
		}
		c.Data["CustomerData"] = result

	} else if utype == "CLUSTEREXECUTIVE" {
		data, err := services.GetClusterAsPerMappingExec(id)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster Name fetch Failed")
			return
		}
		if len(data) <= 0 {
			err = errors.New("Cluster Name  Not Found")
			return
		}
		var Dis Display
		for i := range data {
			var d Field
			d.Id = data[i][2]
			d.Name = data[i][3]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

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
		row, err := db.Db.Query(` select id,
   cluster_name,
product_name,
inventory_available,
product_sku,
cluster_selling_price
 from Inventory where (cluster_name='' OR cluster_name like ?)`, cname)

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
			r.Cluster_name = data4[i][1]
			r.Product_name = data4[i][2]
			r.Inventory_available = data4[i][3]
			r.Product_sku = data4[i][4]
			r.Cluster_selling_price = data4[i][5]

			log.Println(beego.AppConfig.String("loglevel"), "Debug", "data4[i][2]", data4[i][2])

			result = append(result, r)
		}
		c.Data["CustomerData"] = result

	} else {
		data, err := services.GetClusteruser()
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster Name fetch Failed")
			return
		}
		if len(data) <= 0 {
			err = errors.New("Cluster Name  Not Found")
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

		row, err := db.Db.Query(` select id,
   cluster_name,
product_name,
inventory_available,
product_sku,
cluster_selling_price
 from Inventory `)

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
			r.Cluster_name = data4[i][1]
			r.Product_name = data4[i][2]
			r.Inventory_available = data4[i][3]
			r.Product_sku = data4[i][4]
			r.Cluster_selling_price = data4[i][5]

			log.Println(beego.AppConfig.String("loglevel"), "Debug", "data4[i][2]", data4[i][2])

			result = append(result, r)
		}
		c.Data["CustomerData"] = result
	}

	auth, err := utils.IsAuthorized(utype, "PurchaseOrder")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	return
}

func (c *SelectPurchaseOrder) Post() {

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
			c.TplName = "purchaseorder/selectPurchaseOrder/selectPurchaseOrder.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "purchaseorder/selectPurchaseOrder/selectPurchaseOrder.html"
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

	auth, err := utils.IsAuthorized(utype, "PurchaseOrder")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	select_cluster_name := c.Input().Get("select_cluster_name")

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "select_cluster_name - ", select_cluster_name)

	// row, err := db.Db.Query(` select id,
	//   cluster_name,
	// product_name,
	// inventory_available,
	// product_sku
	// from Inventory where (cluster_name='' OR cluster_name like ?) ORDER BY id DESC`, select_cluster_name+"%")

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
		row, err := db.Db.Query(` select id,
   cluster_name,
product_name,
inventory_available,
product_sku,
cluster_selling_price
 from Inventory where (cluster_name='' OR cluster_name like ?) AND (cluster_name='' OR cluster_name like ?) ORDER BY id DESC`, cname, select_cluster_name+"%")

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
			r.Cluster_name = data4[i][1]
			r.Product_name = data4[i][2]
			r.Inventory_available = data4[i][3]
			r.Product_sku = data4[i][4]
			r.Cluster_selling_price = data4[i][5]

			log.Println(beego.AppConfig.String("loglevel"), "Debug", "data4[i][2]", data4[i][2])

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
		row, err := db.Db.Query(` select id,
   cluster_name,
product_name,
inventory_available,
product_sku,
cluster_selling_price
 from Inventory where (cluster_name='' OR cluster_name like ?) AND (cluster_name='' OR cluster_name like ?) ORDER BY id DESC`, cname, select_cluster_name+"%")

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
			r.Cluster_name = data4[i][1]
			r.Product_name = data4[i][2]
			r.Inventory_available = data4[i][3]
			r.Product_sku = data4[i][4]
			r.Cluster_selling_price = data4[i][5]

			log.Println(beego.AppConfig.String("loglevel"), "Debug", "data4[i][2]", data4[i][2])

			result = append(result, r)
		}
		c.Data["CustomerData"] = result
	} else {
		row, err := db.Db.Query(` select id,
   cluster_name,
product_name,
inventory_available,
product_sku,
cluster_selling_price
 from Inventory where (cluster_name='' OR cluster_name like ?) ORDER BY id DESC`, select_cluster_name+"%")

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
			r.Cluster_name = data4[i][1]
			r.Product_name = data4[i][2]
			r.Inventory_available = data4[i][3]
			r.Product_sku = data4[i][4]
			r.Cluster_selling_price = data4[i][5]

			log.Println(beego.AppConfig.String("loglevel"), "Debug", "data4[i][2]", data4[i][2])

			result = append(result, r)
		}
		c.Data["CustomerData"] = result
	}

	// AvailableUnits := data4[0][3]
	// ClusterPrice := data4[0][5]

	// AvailableUnits1, err := strconv.ParseFloat(AvailableUnits, 64)
	// ClusterPrice1, err := strconv.ParseFloat(ClusterPrice, 64)

	// TotalAmount := AvailableUnits1 * ClusterPrice1
	// TotalAmount1 := strconv.FormatFloat(TotalAmount, 'f', 0, 32)

	// c.Data["TotalAmount1"] = TotalAmount1
	return

}

// func daysBetween(a, b time.Time) int {
// 	if a.After(b) {
// 		a, b = b, a
// 	}

// 	days := -a.YearDay()
// 	for year := a.Year(); year < b.Year(); year++ {
// 		days += time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC).YearDay()
// 	}
// 	days += b.YearDay()

// 	return days
// }

// func date(s string) time.Time {
// 	d, _ := time.Parse("2006-01-02", s)
// 	return d
// }

// func rangeDate(start, end time.Time) func() time.Time {
// 	y, m, d := start.Date()
// 	start = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
// 	y, m, d = end.Date()
// 	end = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

// 	return func() time.Time {
// 		if start.After(end) {
// 			return time.Time{}
// 		}
// 		date := start
// 		start = start.AddDate(0, 0, 1)
// 		return date
// 	}
// }
