package clusterwise

import (
	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/model/utils"

	"c2f_admin_webapp/services"
	"c2f_admin_webapp/session"
	"errors"
	"strings"

	"fmt"
	"runtime/debug"

	// "strconv"

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
	Number        int
}

type ClusterData struct {
	ClusterNumber    int
	Clustername      string
	Clustermanager   string
	Clusterexecutive string
	Deliveryuser     string
	Collectionagent  string
	Customer         string
	Product          string
	Producer         string
	Inventory        []Row
}

type MenuData struct {
	Menus string
}

type ClusterwiseData struct {
	// beego.Controller
	middleware.BaseController
}

type Display struct {
	Fields []Field
}

type Field struct {
	Id   string
	Name string

	Email string
}

func (c *ClusterwiseData) Get() {
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
			c.TplName = "analytics/clusterwisedata.html"
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

	id := sess.Get("uid").(string)
	c.Data["id"] = id

	language := sess.Get("language").(string)
	c.Data["language"] = language
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "language", language)

	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "fullname", fullname1)

	auth, err := utils.IsAuthorized(utype, "ClusterWise")
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
			err1 = errors.New("Cluster user  Not Found")
			return
		}

		var Dis Display
		for i := range data1 {
			var d Field
			d.Id = data1[i][2]
			d.Name = data1[i][3]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

		// Get all clusters id and name

		row31, err31 := db.Db.Query(`select id,cluster_name from Clusters where id=?`, data1[0][1])

		if err31 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err31)
			err31 = errors.New("Unable to fetch Clusters Data")
			return
		}
		defer sql.Close(row31)
		_, data31, err31 := sql.Scan(row31)
		if err31 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err31)
			err31 = errors.New("Unable to fetch data of Clusters")
			return
		}

		var result2 []ClusterData

		for i := range data31 {

			var d ClusterData

			d.ClusterNumber = i + 1

			d.Clustername = strings.ToUpper(data31[i][1])

			// Cluster Manager Count
			row1, err1 := db.Db.Query(`SELECT count(*)as acount FROM User2Cluster where cluster_id=? AND clusterexecutive_id IS NULL`, data31[i][0])

			if err1 != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
				err1 = errors.New("Unable to get cluster manager count")
				return
			}
			defer sql.Close(row1)
			// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
			_, data1, err := sql.Scan(row1)
			if err1 != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
				err1 = errors.New("Unable to get cluster manager count data")
				return
			}

			d.Clustermanager = data1[0][0]

			// Cluster Executive Count

			row2, err2 := db.Db.Query(`SELECT count(*)as acount FROM User2Cluster where cluster_id=? AND user_id IS NULL`, data31[i][0])

			if err2 != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err2)
				err2 = errors.New("Unable to get cluster executive count")
				return
			}
			defer sql.Close(row2)
			// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
			_, data2, err := sql.Scan(row2)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to get cluster executive count data")
				return
			}

			d.Clusterexecutive = data2[0][0]

			// Producer Count
			row3, err3 := db.Db.Query(`SELECT count(*)as acount FROM Product2Producer2Cluster where cluster_id=?`, data31[i][0])

			if err3 != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err3)
				err3 = errors.New("Unable to get producer count")
				return
			}
			defer sql.Close(row3)
			// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
			_, data3, err := sql.Scan(row3)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to get producer count data")
				return
			}

			d.Producer = data3[0][0]

			// CustomerCount

			row4, err4 := db.Db.Query(`SELECT count(*)as acount FROM Customer where cluster_name=? and (user_type="CONSUMER" OR user_type="PRODUCER+CONSUMER")`, data31[i][1])

			if err4 != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err4)
				err4 = errors.New("Unable to get customer count")
				return
			}
			defer sql.Close(row4)
			// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
			_, data4, err := sql.Scan(row4)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to get customer count data")
				return
			}

			d.Customer = data4[0][0]

			// DeliveryUser count

			row5, err5 := db.Db.Query(`SELECT count(*)as acount FROM Cluster2DeliveryUser where cluster_id=?`, data31[i][0])

			if err5 != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err5)
				err5 = errors.New("Unable to get deliveryuser count")
				return
			}
			defer sql.Close(row5)
			// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
			_, data5, err := sql.Scan(row5)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to get deliveryuser count data")
				return
			}

			d.Deliveryuser = data5[0][0]

			// Collection Agent Count

			row6, err6 := db.Db.Query(`SELECT count(*)as acount FROM Cluster2CollectionAgent where cluster_id=?`, data31[i][0])

			if err6 != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err6)
				err6 = errors.New("Unable to get agent count")
				return
			}
			defer sql.Close(row6)
			// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
			_, data6, err := sql.Scan(row6)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to get agent count data")
				return
			}

			d.Collectionagent = data6[0][0]

			// Product Count

			row7, err7 := db.Db.Query(`SELECT count(*)as acount FROM Product2Cluster where cluster_id=?`, data31[i][0])

			if err7 != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err7)
				err7 = errors.New("Unable to get product count")
				return
			}
			defer sql.Close(row7)
			// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
			_, data7, err := sql.Scan(row7)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to get product count data")
				return
			}

			d.Product = data7[0][0]

			// Product Inventory data

			row9, err := db.Db.Query(`select product_sku,product_name, inventory_received,inventory_delivered,inventory_available,product_package_qty,product_qty,cluster_name from Inventory where cluster_name=?`, data31[i][1])

			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to inventory data")
				return
			}
			defer sql.Close(row9)
			_, data9, err := sql.Scan(row9)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to product inventory data")
				return
			}
			// fmt.Println(data9)

			var result []Row

			for j := range data9 {

				var r Row
				r.Number = j + 1
				r.SKU = data9[j][0]
				r.ProductName = data9[j][1]
				r.Received = data9[j][2]
				r.Delivered = data9[j][3]
				r.Avaliable = data9[j][4]
				r.Qty = data9[j][5] + data9[j][6]
				r.Cluster = data9[j][7]

				result = append(result, r)
			}

			d.Inventory = result
			result2 = append(result2, d)

		}

		c.Data["ClusterData"] = result2

	} else if utype == "CLUSTEREXECUTIVE" {

		data2, err2 := services.GetClusterAsPerMappingExec(id)

		if err2 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err2)
			err2 = errors.New("Cluster user fetch Failed")
			return
		}
		if len(data2) <= 0 {
			err2 = errors.New("Cluster user  Not Found")
			return
		}

		var Dis Display
		for i := range data2 {
			var d Field
			d.Id = data2[i][2]
			d.Name = data2[i][3]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

		// Get all clusters id and name

		row31, err31 := db.Db.Query(`select id,cluster_name from Clusters where id=?`, data2[0][1])

		if err31 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err31)
			err31 = errors.New("Unable to fetch Clusters Data")
			return
		}
		defer sql.Close(row31)
		_, data31, err31 := sql.Scan(row31)
		if err31 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err31)
			err31 = errors.New("Unable to fetch data of Clusters")
			return
		}

		var result2 []ClusterData

		for i := range data31 {

			var d ClusterData

			d.ClusterNumber = i + 1

			d.Clustername = strings.ToUpper(data31[i][1])

			// Cluster Manager Count
			row1, err1 := db.Db.Query(`SELECT count(*)as acount FROM User2Cluster where cluster_id=? AND clusterexecutive_id IS NULL`, data31[i][0])

			if err1 != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
				err1 = errors.New("Unable to get cluster manager count")
				return
			}
			defer sql.Close(row1)
			// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
			_, data1, err := sql.Scan(row1)
			if err1 != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
				err1 = errors.New("Unable to get cluster manager count data")
				return
			}

			d.Clustermanager = data1[0][0]

			// Cluster Executive Count

			row2, err2 := db.Db.Query(`SELECT count(*)as acount FROM User2Cluster where cluster_id=? AND user_id IS NULL`, data31[i][0])

			if err2 != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err2)
				err2 = errors.New("Unable to get cluster executive count")
				return
			}
			defer sql.Close(row2)
			// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
			_, data2, err := sql.Scan(row2)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to get cluster executive count data")
				return
			}

			d.Clusterexecutive = data2[0][0]

			// Producer Count
			row3, err3 := db.Db.Query(`SELECT count(*)as acount FROM Product2Producer2Cluster where cluster_id=?`, data31[i][0])

			if err3 != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err3)
				err3 = errors.New("Unable to get producer count")
				return
			}
			defer sql.Close(row3)
			// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
			_, data3, err := sql.Scan(row3)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to get producer count data")
				return
			}

			d.Producer = data3[0][0]

			// CustomerCount

			row4, err4 := db.Db.Query(`SELECT count(*)as acount FROM Customer where cluster_name=? and (user_type="CONSUMER" OR user_type="PRODUCER+CONSUMER")`, data31[i][1])

			if err4 != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err4)
				err4 = errors.New("Unable to get customer count")
				return
			}
			defer sql.Close(row4)
			// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
			_, data4, err := sql.Scan(row4)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to get customer count data")
				return
			}

			d.Customer = data4[0][0]

			// DeliveryUser count

			row5, err5 := db.Db.Query(`SELECT count(*)as acount FROM Cluster2DeliveryUser where cluster_id=?`, data31[i][0])

			if err5 != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err5)
				err5 = errors.New("Unable to get deliveryuser count")
				return
			}
			defer sql.Close(row5)
			// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
			_, data5, err := sql.Scan(row5)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to get deliveryuser count data")
				return
			}

			d.Deliveryuser = data5[0][0]

			// Collection Agent Count

			row6, err6 := db.Db.Query(`SELECT count(*)as acount FROM Cluster2CollectionAgent where cluster_id=?`, data31[i][0])

			if err6 != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err6)
				err6 = errors.New("Unable to get agent count")
				return
			}
			defer sql.Close(row6)
			// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
			_, data6, err := sql.Scan(row6)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to get agent count data")
				return
			}

			d.Collectionagent = data6[0][0]

			// Product Count

			row7, err7 := db.Db.Query(`SELECT count(*)as acount FROM Product2Cluster where cluster_id=?`, data31[i][0])

			if err7 != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err7)
				err7 = errors.New("Unable to get product count")
				return
			}
			defer sql.Close(row7)
			// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
			_, data7, err := sql.Scan(row7)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to get product count data")
				return
			}

			d.Product = data7[0][0]

			// Product Inventory data

			row9, err := db.Db.Query(`select product_sku,product_name, inventory_received,inventory_delivered,inventory_available,product_package_qty,product_qty,cluster_name from Inventory where cluster_name=?`, data31[i][1])

			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to inventory data")
				return
			}
			defer sql.Close(row9)
			_, data9, err := sql.Scan(row9)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to product inventory data")
				return
			}
			// fmt.Println(data9)

			var result []Row

			for j := range data9 {

				var r Row
				r.Number = j + 1
				r.SKU = data9[j][0]
				r.ProductName = data9[j][1]
				r.Received = data9[j][2]
				r.Delivered = data9[j][3]
				r.Avaliable = data9[j][4]
				r.Qty = data9[j][5] + data9[j][6]
				r.Cluster = data9[j][7]

				result = append(result, r)
			}

			d.Inventory = result
			result2 = append(result2, d)

		}

		c.Data["ClusterData"] = result2

	} else {

		// HOADMIN CODE

		// Get all clusters id and name

		row31, err31 := db.Db.Query(`select id,cluster_name from Clusters`)

		if err31 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err31)
			err31 = errors.New("Unable to fetch Clusters Data")
			return
		}
		defer sql.Close(row31)
		_, data31, err31 := sql.Scan(row31)
		if err31 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err31)
			err31 = errors.New("Unable to fetch data of Clusters")
			return
		}

		var Dis Display
		for i := range data31 {
			var d Field
			d.Id = data31[i][0]
			d.Name = data31[i][1]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

		var result2 []ClusterData

		for i := range data31 {

			var d ClusterData

			d.ClusterNumber = i + 1

			d.Clustername = strings.ToUpper(data31[i][1])

			// Cluster Manager Count
			row1, err1 := db.Db.Query(`SELECT count(*)as acount FROM User2Cluster where cluster_id=? AND clusterexecutive_id IS NULL`, data31[i][0])

			if err1 != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
				err1 = errors.New("Unable to get cluster manager count")
				return
			}
			defer sql.Close(row1)
			// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
			_, data1, err := sql.Scan(row1)
			if err1 != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
				err1 = errors.New("Unable to get cluster manager count data")
				return
			}

			d.Clustermanager = data1[0][0]

			// Cluster Executive Count

			row2, err2 := db.Db.Query(`SELECT count(*)as acount FROM User2Cluster where cluster_id=? AND user_id IS NULL`, data31[i][0])

			if err2 != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err2)
				err2 = errors.New("Unable to get cluster executive count")
				return
			}
			defer sql.Close(row2)
			// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
			_, data2, err := sql.Scan(row2)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to get cluster executive count data")
				return
			}

			d.Clusterexecutive = data2[0][0]

			// Producer Count
			row3, err3 := db.Db.Query(`SELECT count(*)as acount FROM Product2Producer2Cluster where cluster_id=?`, data31[i][0])

			if err3 != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err3)
				err3 = errors.New("Unable to get producer count")
				return
			}
			defer sql.Close(row3)
			// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
			_, data3, err := sql.Scan(row3)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to get producer count data")
				return
			}

			d.Producer = data3[0][0]

			// CustomerCount

			row4, err4 := db.Db.Query(`SELECT count(*)as acount FROM Customer where cluster_name=? and (user_type="CONSUMER" OR user_type="PRODUCER+CONSUMER")`, data31[i][1])

			if err4 != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err4)
				err4 = errors.New("Unable to get customer count")
				return
			}
			defer sql.Close(row4)
			// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
			_, data4, err := sql.Scan(row4)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to get customer count data")
				return
			}

			d.Customer = data4[0][0]

			// DeliveryUser count

			row5, err5 := db.Db.Query(`SELECT count(*)as acount FROM Cluster2DeliveryUser where cluster_id=?`, data31[i][0])

			if err5 != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err5)
				err5 = errors.New("Unable to get deliveryuser count")
				return
			}
			defer sql.Close(row5)
			// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
			_, data5, err := sql.Scan(row5)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to get deliveryuser count data")
				return
			}

			d.Deliveryuser = data5[0][0]

			// Collection Agent Count

			row6, err6 := db.Db.Query(`SELECT count(*)as acount FROM Cluster2CollectionAgent where cluster_id=?`, data31[i][0])

			if err6 != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err6)
				err6 = errors.New("Unable to get agent count")
				return
			}
			defer sql.Close(row6)
			// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
			_, data6, err := sql.Scan(row6)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to get agent count data")
				return
			}

			d.Collectionagent = data6[0][0]

			// Product Count

			row7, err7 := db.Db.Query(`SELECT count(*)as acount FROM Product2Cluster where cluster_id=?`, data31[i][0])

			if err7 != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err7)
				err7 = errors.New("Unable to get product count")
				return
			}
			defer sql.Close(row7)
			// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
			_, data7, err := sql.Scan(row7)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to get product count data")
				return
			}

			d.Product = data7[0][0]

			// Product Inventory data

			row9, err := db.Db.Query(`select product_sku,product_name, inventory_received,inventory_delivered,inventory_available,product_package_qty,product_qty,cluster_name from Inventory where cluster_name=?`, data31[i][1])

			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to inventory data")
				return
			}
			defer sql.Close(row9)
			_, data9, err := sql.Scan(row9)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to product inventory data")
				return
			}
			// fmt.Println(data9)

			var result []Row

			for j := range data9 {

				var r Row
				r.Number = j + 1
				r.SKU = data9[j][0]
				r.ProductName = data9[j][1]
				r.Received = data9[j][2]
				r.Delivered = data9[j][3]
				r.Avaliable = data9[j][4]
				r.Qty = data9[j][5] + data9[j][6]
				r.Cluster = data9[j][7]

				result = append(result, r)
			}

			d.Inventory = result
			result2 = append(result2, d)

		}

		c.Data["ClusterData"] = result2

	}

	return
}

func (c *ClusterwiseData) Post() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Dashboard Page Start")
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
			c.TplName = "analytics/clusterwisedata.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "ClusterwiseData Page Fail")
		} else {

			c.TplName = "analytics/clusterwisedata.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "ClusterwiseData  Page Success")
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

	id := sess.Get("uid").(string)
	c.Data["id"] = id

	language := sess.Get("language").(string)
	c.Data["language"] = language
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "language", language)

	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "fullname", fullname1)

	auth, err := utils.IsAuthorized(utype, "ClusterWise")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	row41, err41 := db.Db.Query(`select id,cluster_name from Clusters`)

	if err41 != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err41)
		err41 = errors.New("Unable to fetch Clusters Data")
		return
	}
	defer sql.Close(row41)
	_, data41, err41 := sql.Scan(row41)
	if err41 != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err41)
		err41 = errors.New("Unable to fetch data of Clusters")
		return
	}

	var Dis Display
	for i := range data41 {
		var d Field
		d.Id = data41[i][0]
		d.Name = data41[i][1]
		Dis.Fields = append(Dis.Fields, d)
	}
	c.Data["Dis"] = Dis
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

	select_cluster_name := c.Input().Get("select_cluster_name")

	if select_cluster_name == "" {
		err = errors.New("Please Select Cluster")
		return
	}
	fmt.Println("Selected Cluster :- ", select_cluster_name)

	// Get all clusters id and name

	row31, err31 := db.Db.Query(`select id,cluster_name from Clusters where id=?`, select_cluster_name)

	if err31 != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err31)
		err31 = errors.New("Unable to fetch Clusters Data")
		return
	}
	defer sql.Close(row31)
	_, data31, err31 := sql.Scan(row31)
	if err31 != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err31)
		err31 = errors.New("Unable to fetch data of Clusters")
		return
	}

	var result2 []ClusterData

	for i := range data31 {

		var d ClusterData

		d.ClusterNumber = i + 1

		d.Clustername = strings.ToUpper(data31[i][1])

		// Cluster Manager Count
		row1, err1 := db.Db.Query(`SELECT count(*)as acount FROM User2Cluster where cluster_id=? AND clusterexecutive_id IS NULL`, data31[i][0])

		if err1 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
			err1 = errors.New("Unable to get cluster manager count")
			return
		}
		defer sql.Close(row1)
		// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
		_, data1, err := sql.Scan(row1)
		if err1 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
			err1 = errors.New("Unable to get cluster manager count data")
			return
		}

		d.Clustermanager = data1[0][0]

		// Cluster Executive Count

		row2, err2 := db.Db.Query(`SELECT count(*)as acount FROM User2Cluster where cluster_id=? AND user_id IS NULL`, data31[i][0])

		if err2 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err2)
			err2 = errors.New("Unable to get cluster executive count")
			return
		}
		defer sql.Close(row2)
		// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
		_, data2, err := sql.Scan(row2)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get cluster executive count data")
			return
		}

		d.Clusterexecutive = data2[0][0]

		// Producer Count
		row3, err3 := db.Db.Query(`SELECT count(*)as acount FROM Product2Producer2Cluster where cluster_id=?`, data31[i][0])

		if err3 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err3)
			err3 = errors.New("Unable to get producer count")
			return
		}
		defer sql.Close(row3)
		// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
		_, data3, err := sql.Scan(row3)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get producer count data")
			return
		}

		d.Producer = data3[0][0]

		// CustomerCount

		row4, err4 := db.Db.Query(`SELECT count(*)as acount FROM Customer where cluster_name=? and (user_type="CONSUMER" OR user_type="PRODUCER+CONSUMER")`, data31[i][1])

		if err4 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err4)
			err4 = errors.New("Unable to get customer count")
			return
		}
		defer sql.Close(row4)
		// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
		_, data4, err := sql.Scan(row4)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get customer count data")
			return
		}

		d.Customer = data4[0][0]

		// DeliveryUser count

		row5, err5 := db.Db.Query(`SELECT count(*)as acount FROM Cluster2DeliveryUser where cluster_id=?`, data31[i][0])

		if err5 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err5)
			err5 = errors.New("Unable to get deliveryuser count")
			return
		}
		defer sql.Close(row5)
		// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
		_, data5, err := sql.Scan(row5)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get deliveryuser count data")
			return
		}

		d.Deliveryuser = data5[0][0]

		// Collection Agent Count

		row6, err6 := db.Db.Query(`SELECT count(*)as acount FROM Cluster2CollectionAgent where cluster_id=?`, data31[i][0])

		if err6 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err6)
			err6 = errors.New("Unable to get agent count")
			return
		}
		defer sql.Close(row6)
		// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
		_, data6, err := sql.Scan(row6)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get agent count data")
			return
		}

		d.Collectionagent = data6[0][0]

		// Product Count

		row7, err7 := db.Db.Query(`SELECT count(*)as acount FROM Product2Cluster where cluster_id=?`, data31[i][0])

		if err7 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err7)
			err7 = errors.New("Unable to get product count")
			return
		}
		defer sql.Close(row7)
		// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row3)
		_, data7, err := sql.Scan(row7)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get product count data")
			return
		}

		d.Product = data7[0][0]

		// Product Inventory data

		row9, err := db.Db.Query(`select product_sku,product_name, inventory_received,inventory_delivered,inventory_available,product_package_qty,product_qty,cluster_name from Inventory where cluster_name=?`, data31[i][1])

		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to inventory data")
			return
		}
		defer sql.Close(row9)
		_, data9, err := sql.Scan(row9)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to product inventory data")
			return
		}
		// fmt.Println(data9)

		var result []Row

		for j := range data9 {

			var r Row
			r.Number = j + 1
			r.SKU = data9[j][0]
			r.ProductName = data9[j][1]
			r.Received = data9[j][2]
			r.Delivered = data9[j][3]
			r.Avaliable = data9[j][4]
			r.Qty = data9[j][5] + data9[j][6]
			r.Cluster = data9[j][7]

			result = append(result, r)
		}

		d.Inventory = result
		result2 = append(result2, d)

	}

	c.Data["ClusterData"] = result2

	return
}
