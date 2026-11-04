package selectSupplyProductOrder

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
	Id                       string
	Source_cluster_name      string
	Beneficiary_cluster_name string
	Product_name             string
	Product_sku              string
	Number_of_units          string
	Delivery_charges         string
	Total_price              string
	Status                   string
}

type Display struct {
	Fields []Field
}

type Field struct {
	Id   string
	Name string
}

type SelectSupplyProductOrder struct {
	// beego.Controller
	middleware.BaseController
}

func (c *SelectSupplyProductOrder) Get() {
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
			c.TplName = "supplyproductorder/selectSupplyProductOrder/selectSupplyProductOrder.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "supplyproductorder/selectSupplyProductOrder/selectSupplyProductOrder.html"
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

	auth, err := utils.IsAuthorized(utype, "SupplyProductOrder")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1

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
	source_cluster_name,
	beneficiary_cluster_name, 
	product_name, 
	product_sku,
	number_of_units, 
	delivery_charges, 
	total_price, 
	status
 from ClusterPurchaseProduct where beneficiary_cluster_name=? AND status="DEMAND_ORDER_CREATED"`, cname)

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
			r.Source_cluster_name = data4[i][1]
			r.Beneficiary_cluster_name = data4[i][2]
			r.Product_name = data4[i][3]
			r.Product_sku = data4[i][4]
			r.Number_of_units = data4[i][5]
			r.Delivery_charges = data4[i][6]
			r.Total_price = data4[i][7]
			r.Status = data4[i][8]

			log.Println(beego.AppConfig.String("loglevel"), "Debug", "data4[i][2]", data4[i][2])

			result = append(result, r)
		}
		c.Data["CustomerData"] = result

	} else if utype == "CLUSTEREXECUTIVE" {

		// data2, err2 := services.GetClusterAsPerMappingExec(id)
		// if err2 != nil {
		// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err2)
		// 	err2 = errors.New("Cluster user fetch Failed")
		// 	return
		// }
		// if len(data2) <= 0 {
		// 	err2 = errors.New("Cluster user  Not Found")
		// 	return
		// }
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
	source_cluster_name,
	beneficiary_cluster_name, 
	product_name, 
	product_sku,
	number_of_units, 
	delivery_charges, 
	total_price, 
	status
 from ClusterPurchaseProduct where beneficiary_cluster_name=? AND status="DEMAND_ORDER_CREATED"`, cname)

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
			r.Source_cluster_name = data4[i][1]
			r.Beneficiary_cluster_name = data4[i][2]
			r.Product_name = data4[i][3]
			r.Product_sku = data4[i][4]
			r.Number_of_units = data4[i][5]
			r.Delivery_charges = data4[i][6]
			r.Total_price = data4[i][7]
			r.Status = data4[i][8]

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
	source_cluster_name,
	beneficiary_cluster_name, 
	product_name, 
	product_sku,
	number_of_units, 
	delivery_charges, 
	total_price, 
	status
 from ClusterPurchaseProduct where status="DEMAND_ORDER_ACCEPTED" OR status="DEMAND_ORDER_CREATED"`)

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
			r.Source_cluster_name = data4[i][1]
			r.Beneficiary_cluster_name = data4[i][2]
			r.Product_name = data4[i][3]
			r.Product_sku = data4[i][4]
			r.Number_of_units = data4[i][5]
			r.Delivery_charges = data4[i][6]
			r.Total_price = data4[i][7]
			r.Status = data4[i][8]

			log.Println(beego.AppConfig.String("loglevel"), "Debug", "data4[i][2]", data4[i][2])

			result = append(result, r)
		}
		c.Data["CustomerData"] = result
	}
	return
}

func (c *SelectSupplyProductOrder) Post() {

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
			c.TplName = "supplyproductorder/selectSupplyProductOrder/selectSupplyProductOrder.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "supplyproductorder/selectSupplyProductOrder/selectSupplyProductOrder.html"
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

	uid := sess.Get("uid").(string)
	c.Data["id"] = uid
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "id", uid)

	auth, err := utils.IsAuthorized(utype, "SupplyProductOrder")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	select_sourcecluster_name := c.Input().Get("select_sourcecluster_name")
	select_beneficiarycluster_name := c.Input().Get("select_beneficiarycluster_name")

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "select_sourcecluster_name - ", select_sourcecluster_name)

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

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "cname", cname)
		row, err := db.Db.Query(` select id, 
	source_cluster_name,
	beneficiary_cluster_name, 
	product_name, 
	product_sku,
	number_of_units, 
	delivery_charges, 
	total_price, 
	status
 from ClusterPurchaseProduct where source_cluster_name=? AND (source_cluster_name='' OR source_cluster_name like ?) AND  (source_cluster_name='' OR source_cluster_name like ?) AND (beneficiary_cluster_name='' OR beneficiary_cluster_name like ?) AND status="DEMAND_ORDER_ACCEPTED" OR status="DEMAND_ORDER_CREATED" ORDER BY id DESC`, data1[0][1], cname, select_sourcecluster_name+"%", select_beneficiarycluster_name+"%")

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
			r.Source_cluster_name = data4[i][1]
			r.Beneficiary_cluster_name = data4[i][2]
			r.Product_name = data4[i][3]
			r.Product_sku = data4[i][4]
			r.Number_of_units = data4[i][5]
			r.Delivery_charges = data4[i][6]
			r.Total_price = data4[i][7]
			r.Status = data4[i][8]

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
	source_cluster_name,
	beneficiary_cluster_name, 
	product_name, 
	product_sku,
	number_of_units, 
	delivery_charges, 
	total_price, 
	status
 from ClusterPurchaseProduct where (source_cluster_name='' OR source_cluster_name like ?) AND (source_cluster_name='' OR source_cluster_name like ?) AND (beneficiary_cluster_name='' OR beneficiary_cluster_name like ?) AND status="DEMAND_ORDER_ACCEPTED" OR status="DEMAND_ORDER_CREATED" ORDER BY id DESC`, cname, select_sourcecluster_name+"%", select_beneficiarycluster_name+"%")

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
			r.Source_cluster_name = data4[i][1]
			r.Beneficiary_cluster_name = data4[i][2]
			r.Product_name = data4[i][3]
			r.Product_sku = data4[i][4]
			r.Number_of_units = data4[i][5]
			r.Delivery_charges = data4[i][6]
			r.Total_price = data4[i][7]
			r.Status = data4[i][8]

			log.Println(beego.AppConfig.String("loglevel"), "Debug", "data4[i][2]", data4[i][2])

			result = append(result, r)
		}
		c.Data["CustomerData"] = result
	} else {
		row, err := db.Db.Query(` select id, 
	source_cluster_name,
	beneficiary_cluster_name, 
	product_name, 
	product_sku,
	number_of_units, 
	delivery_charges, 
	total_price, 
	status
 from ClusterPurchaseProduct where (source_cluster_name='' OR source_cluster_name like ?) AND (beneficiary_cluster_name='' OR beneficiary_cluster_name like ?) AND status="DEMAND_ORDER_ACCEPTED" OR status="DEMAND_ORDER_CREATED" ORDER BY id DESC`, select_sourcecluster_name+"%", select_beneficiarycluster_name+"%")

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
			r.Source_cluster_name = data4[i][1]
			r.Beneficiary_cluster_name = data4[i][2]
			r.Product_name = data4[i][3]
			r.Product_sku = data4[i][4]
			r.Number_of_units = data4[i][5]
			r.Delivery_charges = data4[i][6]
			r.Total_price = data4[i][7]
			r.Status = data4[i][8]

			log.Println(beego.AppConfig.String("loglevel"), "Debug", "data4[i][2]", data4[i][2])

			result = append(result, r)
		}
		c.Data["CustomerData"] = result

		// row, err := db.Db.Query(` select id,
		//   cluster_name,
		// product_name,
		// inventory_available,
		// product_sku
		// from Inventory where (cluster_name='' OR cluster_name like ?) ORDER BY id DESC`, select_cluster_name+"%")

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
