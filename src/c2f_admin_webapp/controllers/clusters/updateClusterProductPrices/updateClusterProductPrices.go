package updateClusterProductPrices

import (
	"c2f_admin_webapp/session"

	"c2f_admin_webapp/model/utils"
	"errors"
	"runtime/debug"
	"strings"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/services"

	"github.com/astaxie/beego"
	"github.com/xuri/excelize"

	// "encoding/csv"
	"fmt"
	// "net/http"
	"strconv"

	// "time"

	"log"

	"c2f_admin_webapp/model/database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type UpdateClusterProductPrice struct {
	// beego.Controller
	middleware.BaseController
}

type Row struct {
	Cluster        string
	Product        string
	SKU            string
	BuyPrice       string
	SellPrice      string
	SupplyCharge   string
	PacakageCharge string
	BankCharge     string
	Discount       string
}

func (c *UpdateClusterProductPrice) Get() {
	//	AdminId := c.Ctx.Input.Param(":AdminID")
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
			c.TplName = "clusters/updateClusterProductPrices/updateClusterProductPrices.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "clusters/updateClusterProductPrices/updateClusterProductPrices.html"
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

	id := sess.Get("uid").(string)
	c.Data["id"] = id
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "id", id)

	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1

	clustername := sess.Get("clustername").(string)
	clustername1 := strings.ToUpper(clustername)
	c.Data["clustername"] = clustername1
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "clustername", clustername1)

	auth, err := utils.IsAuthorized(utype, "Cluster")
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

		row2, err := db.Db.Query(`
		select Lookup_products.product_sku,
		Lookup_products.product_name,
		Inventory.cluster_name,
		Inventory.cluster_buying_price,
		Inventory.cluster_selling_price,
		Inventory.cluster_supply_charge,
		Inventory.cluster_packaging_charge,
		Inventory.cluster_bank_charge,
		Inventory.cluster_discount
		from Inventory
		left join  Lookup_products on Inventory.product_id=Lookup_products.id where Inventory.cluster_name=?`, data1[0][3])
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get Cluster-Product data")
			return
		}
		defer sql.Close(row2)
		_, data2, err := sql.Scan(row2)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get Cluster-Product data")
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data Cluster-Product - ", data2, "\nData len - ", len(data2))
		if len(data2) <= 0 {
			err = errors.New("Cluster-Product data not found")
			return
		}

		var result []Row

		for j := range data2 {

			var r Row
			r.Cluster = data2[j][2]
			r.Product = data2[j][1]
			r.SKU = data2[j][0]
			r.BuyPrice = data2[j][3]
			r.SellPrice = data2[j][4]
			r.SupplyCharge = data2[j][5]
			r.PacakageCharge = data2[j][6]
			r.BankCharge = data2[j][7]
			r.Discount = data2[j][8]

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
			err1 = errors.New("Cluster user  Not Found")
			return
		}

		row2, err := db.Db.Query(`
		select Lookup_products.product_sku,
		Lookup_products.product_name,
		Inventory.cluster_name,
		Inventory.cluster_buying_price,
		Inventory.cluster_selling_price,
		Inventory.cluster_supply_charge,
		Inventory.cluster_packaging_charge,
		Inventory.cluster_bank_charge,
		Inventory.cluster_discount
		from Inventory
		left join  Lookup_products on Inventory.product_id=Lookup_products.id where Inventory.cluster_name=?`, data1[0][3])
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get Cluster-Product data")
			return
		}
		defer sql.Close(row2)
		_, data2, err := sql.Scan(row2)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get Cluster-Product data")
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data Cluster-Product - ", data2, "\nData len - ", len(data2))
		if len(data2) <= 0 {
			err = errors.New("Cluster-Product data not found")
			return
		}

		var result []Row

		for j := range data2 {

			var r Row
			r.Cluster = data2[j][2]
			r.Product = data2[j][1]
			r.SKU = data2[j][0]
			r.BuyPrice = data2[j][3]
			r.SellPrice = data2[j][4]
			r.SupplyCharge = data2[j][5]
			r.PacakageCharge = data2[j][6]
			r.BankCharge = data2[j][7]
			r.Discount = data2[j][8]

			result = append(result, r)

		}

		c.Data["CustomerData"] = result

	} else {

		// HOADMIN

		row2, err := db.Db.Query(`
		select Lookup_products.product_sku,
		Lookup_products.product_name,
		Inventory.cluster_name,
		Inventory.cluster_buying_price,
		Inventory.cluster_selling_price,
		Inventory.cluster_supply_charge,
		Inventory.cluster_packaging_charge,
		Inventory.cluster_bank_charge,
		Inventory.cluster_discount
		from Inventory
		left join  Lookup_products on Inventory.product_id=Lookup_products.id`)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get Cluster-Product data")
			return
		}
		defer sql.Close(row2)
		_, data2, err := sql.Scan(row2)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get Cluster-Product data")
			return
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data Cluster-Product - ", data2, "\nData len - ", len(data2))
		if len(data2) <= 0 {
			err = errors.New("Cluster-Product data not found")
			return
		}

		var result []Row

		for j := range data2 {

			var r Row
			r.Cluster = data2[j][2]
			r.Product = data2[j][1]
			r.SKU = data2[j][0]
			r.BuyPrice = data2[j][3]
			r.SellPrice = data2[j][4]
			r.SupplyCharge = data2[j][5]
			r.PacakageCharge = data2[j][6]
			r.BankCharge = data2[j][7]
			r.Discount = data2[j][8]

			result = append(result, r)

		}

		c.Data["CustomerData"] = result

	}

	return
}

func (c *UpdateClusterProductPrice) Post() {
	//	AdminId := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Info", "UpdateClusterProductPrice post start")
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
			c.TplName = "clusters/updateClusterProductPrices/updateClusterProductPrices.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "UpdateClusterProductPrice Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Cluster Product Prices Updated Successfully"
			c.TplName = "clusters/updateClusterProductPrices/updateClusterProductPrices.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "UpdateClusterProductPrice  Page Success")
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

	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1

	clustername := sess.Get("clustername").(string)
	clustername1 := strings.ToUpper(clustername)
	c.Data["clustername"] = clustername1
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "clustername", clustername1)

	id := sess.Get("uid").(string)
	c.Data["id"] = id
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "id", id)

	auth, err := utils.IsAuthorized(utype, "Cluster")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	c.TplName = "updateProductPrice.html"

	// Get the  uploaded Csv File

	file, information, err := c.GetFile("uploadname")

	if err != nil {

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "File - ", "File retrieval failure")
		c.Data["DisplayMessage"] = "File retrieval failure"

		return
	}

	defer file.Close()

	filename := information.Filename
	fmt.Println(filename)

	f, err := excelize.OpenReader(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	rows := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "rows  - ", rows)

	var excelsheetrow int
	var excelsheetrow1, col1, col2, col3, col4, col5, col6, col7, col8, cell1, cell2, cell3, cell4, cell5, cell6, cell7, cell8 string

	for j := range rows {

		excelsheetrow = j + 3
		excelsheetrow1 = strconv.Itoa(excelsheetrow)

		// Get value from the cell by given worksheet name and axis.

		// ClusterName
		col1 = "A" + excelsheetrow1
		cell1 = f.GetCellValue("Sheet1", col1)

		// ProductSKU
		col2 = "C" + excelsheetrow1
		cell2 = f.GetCellValue("Sheet1", col2)

		// ClusterBuyPrice
		col3 = "D" + excelsheetrow1
		cell3 = f.GetCellValue("Sheet1", col3)
		// ClusterSellPrice
		col4 = "E" + excelsheetrow1
		cell4 = f.GetCellValue("Sheet1", col4)
		// ClusterSupplyCharge
		col5 = "F" + excelsheetrow1
		cell5 = f.GetCellValue("Sheet1", col5)
		// ClusterPackageCharge
		col6 = "G" + excelsheetrow1
		cell6 = f.GetCellValue("Sheet1", col6)
		// ClusterBankCharge
		col7 = "H" + excelsheetrow1
		cell7 = f.GetCellValue("Sheet1", col7)
		// ClusterDiscount
		col8 = "I" + excelsheetrow1
		cell8 = f.GetCellValue("Sheet1", col8)

		row33, err := db.Db.Query(`select id from Clusters where cluster_name=?`, cell1)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get cluster id")
			return
		}
		defer sql.Close(row33)
		_, data33, err := sql.Scan(row33)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get cluster id")
			return
		}

		row23, err23 := db.Db.Query(`select id from  Lookup_products WHERE product_sku=?`, cell2)
		if err23 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get product id")
			return
		}
		defer sql.Close(row23)
		_, data23, err := sql.Scan(row23)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get product id")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Cluster id  - ", data33[0][0])
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Product id  - ", data23[0][0])

		row, err25 := db.Db.Query(`select count(*) FROM Product2Cluster where product_id=? and cluster_id=? `, data23[0][0], data33[0][0])

		if err25 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err25)
			err25 = errors.New("Unable to fetch count")
			return
		}
		defer sql.Close(row)
		_, data, err := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to fetch count")
			return
		}

		countlen := data[0][0]

		if countlen != "0" {

			// Cluster to Product Mapping Exsists Update in Product2Cluster

			res, err := db.Db.Exec(`UPDATE Product2Cluster SET cluster_buying_price=?,cluster_selling_price=?,cluster_supply_charge=?,cluster_packaging_charge=?,cluster_bank_charge=?,cluster_discount=?,updatedby=?,update_date=now() WHERE product_id = ? and cluster_id=?`, cell3, cell4, cell5, cell6, cell7, cell8, id, data23[0][0], data33[0][0])
			if err != nil {
				err = errors.New("Cluster-Product updation failed")
				return
			}

			i, err := res.RowsAffected()
			if err != nil || i == 0 {
				err = errors.New("Cluster-Product updation failed no row affected")
				return
			}

			row2, err := db.Db.Query(`select count(*) from Inventory WHERE product_id = ? AND cluster_name=? `, data23[0][0], cell1)

			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to fetch count from inventory")
				return
			}
			defer sql.Close(row2)
			_, data2, err := sql.Scan(row2)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				err = errors.New("Unable to fetch count inventory")
				return
			}

			countlen2 := data2[0][0]

			if countlen2 != "0" {

				// Product mapped with Cluster is also present in Inventory , update Inventory	Table

				res2, err := db.Db.Exec(`UPDATE Inventory SET cluster_buying_price=?,cluster_selling_price=?,cluster_supply_charge=?,cluster_packaging_charge=?,cluster_bank_charge=?,cluster_discount=?,updated_by=?,update_date=now() WHERE product_id = ? and cluster_name=?`, cell3, cell4, cell5, cell6, cell7, cell8, id, data23[0][0], cell1)
				if err != nil {
					err = errors.New("Cluster-Product in inventory updation failed")
					return
				}

				i, err := res2.RowsAffected()
				if err != nil || i == 0 {
					err = errors.New("Cluster-Product  in inventory  updation failed no row affected")
					return
				}

			}

		} else {

			// Update Non Mapped Cluster Products data in Inventory Table

			res3, err := db.Db.Exec(`UPDATE Inventory SET cluster_buying_price=?,cluster_selling_price=?,cluster_supply_charge=?,cluster_packaging_charge=?,cluster_bank_charge=?,cluster_discount=?,updated_by=?,update_date=now() WHERE product_id = ? and cluster_name=?`, cell3, cell4, cell5, cell6, cell7, cell8, id, data23[0][0], cell1)
			if err != nil {
				err = errors.New("Cluster-Product in inventory updation failed")
				return
			}

			i, err := res3.RowsAffected()
			if err != nil || i == 0 {
				err = errors.New("Cluster-Product  in inventory  updation failed no row affected")
				return
			}

		}

	}

	return
}
