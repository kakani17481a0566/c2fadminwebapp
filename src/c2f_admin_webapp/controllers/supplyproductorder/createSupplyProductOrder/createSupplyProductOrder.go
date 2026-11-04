package createSupplyProductOrder

import (
	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/model/utils"
	"c2f_admin_webapp/services"
	"c2f_admin_webapp/session"
	"errors"
	"runtime/debug"

	"strconv"
	"strings"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"

	"github.com/astaxie/beego"
)

//	type Row struct {
//		Sourceclustername      string
//		Beneficiaryclustername string
//		Productname            string
//		Productsku             string
//		Requiredunits          string
//	}
type Display struct {
	Fields []Field
}
type Field struct {
	Id    string
	Name  string
	Email string
}
type CreateSupplyProductOrder struct {
	// beego.Controller
	middleware.BaseController
}

func (c *CreateSupplyProductOrder) Get() {
	AdminId := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "AdminId - ", AdminId)

	log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User Start")
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
			c.TplName = "supplyproductorder/createSupplyProductOrder/createSupplyProductOrder.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "supplyproductorder/createSupplyProductOrder/createSupplyProductOrder.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User Page Success")
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

	maptoken := sess.Get("token")
	c.Data["maptoken"] = maptoken
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "maptoken", maptoken)

	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1

	data, err := services.GetClusteruser()
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

	auth, err := utils.IsAuthorized(utype, "SupplyProductOrder")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	row, err := db.Db.Query(`select id, 
	source_cluster_name,
	beneficiary_cluster_name, 
	product_name, 
	product_sku,
	number_of_units, 
	delivery_charges, 
	total_price, 
	status,
	product_package_qty,
	product_qty
 from ClusterPurchaseProduct where id=? limit 1`, AdminId)
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
		c.Data["ID"] = data3[i][0]
		c.Data["Source_cluster_name"] = data3[i][1]
		c.Data["Beneficiary_cluster_name"] = data3[i][2]
		c.Data["Product_name"] = data3[i][3]
		c.Data["Product_sku"] = data3[i][4]
		c.Data["Number_of_units"] = data3[i][5]
		c.Data["Delivery_charges"] = data3[i][6]
		c.Data["Total_price"] = data3[i][7]
		c.Data["Status"] = data3[i][8]
		c.Data["ProductPackageQty"] = data3[i][9]
		c.Data["ProductQty"] = data3[i][10]
	}

	return
}
func (c *CreateSupplyProductOrder) Post() {
	AdminId := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "AdminId - ", AdminId)

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
			c.TplName = "supplyproductorder/createSupplyProductOrder/createSupplyProductOrder.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Update Cluster Purchase Product Order Successfully"
			c.TplName = "supplyproductorder/createSupplyProductOrder/createSupplyProductOrder.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer  Page Success")
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

	auth, err := utils.IsAuthorized(utype, "SupplyProductOrder")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "adminid - ", AdminId)

	row, err := db.Db.Query(`select id,sku_type from ClusterPurchaseProduct where id=?`, AdminId)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get sku type data")
		return
	}
	defer sql.Close(row)
	_, data3, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get sku type data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data sku type - ", data3, "\nData len - ", len(data3))
	if len(data3) <= 0 {
		err = errors.New("skutype data not found")
		return
	}

	skutype := data3[0][1]

	input_source_clustername := c.Input().Get("input_source_clustername")
	//	input_product_name := c.Input().Get("input_product_name")
	input_product_sku := c.Input().Get("input_product_sku")
	input_beneficiary_clustername := c.Input().Get("input_beneficiary_clustername")
	input_required_units := c.Input().Get("input_required_units")
	input_status := c.Input().Get("input_status")

	//	input_total_price := c.Input().Get("input_total_price")

	input_product_package_qty := c.Input().Get("input_product_package_qty")
	input_product_qty := c.Input().Get("input_product_qty")
	input_id := c.Input().Get("input_id")
	//input_delivery_charge := c.Input().Get("input_delivery_charge")

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_product_sku - ", input_product_sku)

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_source_clustername - ", input_source_clustername)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_beneficiary_clustername - ", input_beneficiary_clustername)

	sourceaunits, sourceclustersprice, sourcereunits, err := services.Getinventoryavalibeunits(input_product_sku, input_source_clustername)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Inventory Get Failed4")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "sourceaunits - ", sourceaunits)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "sourceclustersprice - ", sourceclustersprice)

	sourceaunits1, err := strconv.ParseFloat(sourceaunits, 64)
	sourcereunits1, err := strconv.ParseFloat(sourcereunits, 64)

	runits, err := strconv.ParseFloat(input_required_units, 64)
	//dcharge, err := strconv.ParseFloat(input_delivery_charges, 64)

	input_courier_details := c.Input().Get("input_courier_details")
	input_delivery_channel_partner_type := c.Input().Get("input_delivery_channel_partner_type")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_status - ", input_status)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_delivery_channel_partner_type - ", input_delivery_channel_partner_type)

	input_courier_details1 := input_delivery_channel_partner_type + "::" + input_courier_details

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_courier_details1 - ", input_courier_details1)

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "sourceaunits1 - ", sourceaunits1)
	//	log.Println(beego.AppConfig.String("loglevel"), "Debug", "dcharge - ", dcharge)

	// dunits, err := strconv.ParseFloat(input_delivered_units, 64)

	TS := sourceaunits1 - runits
	TS1 := strconv.FormatFloat(TS, 'f', 0, 32)

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "TS1 - ", TS1)

	TSR := sourcereunits1 - runits
	TSR1 := strconv.FormatFloat(TSR, 'f', 0, 32)

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "TSR1 - ", TSR1)

	if runits > sourceaunits1 {

		//log.Println(beego.AppConfig.String("loglevel"), "Debug", "inside - ")

		err = errors.New("Insufficent Available Units")
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)

		return
	}

	if input_status == "SUPPLY_COMPLETED" && input_product_qty == "Lts" {

		var input_delivery_charges, tpp string
		var TPP1 float64

		input_delivery_charges = c.Input().Get("input_delivery_charges")

		if input_delivery_charges == "" || input_courier_details1 == "" {
			err = errors.New(" * fields cannot be empty ")
			return
		}

		input_delivery_charges1, err := strconv.ParseFloat(input_delivery_charges, 64)
		if err != nil {

			err = errors.New("Unable to get Delivery Charge")
			return
		}
		input_total_price := c.Input().Get("input_total_price")

		input_total_price1, err := strconv.ParseFloat(input_total_price, 64)
		TPP1 = input_total_price1 + input_delivery_charges1
		tpp = strconv.FormatFloat(TPP1, 'f', 0, 32)

		c.Data["tpp"] = tpp
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "total price = (productprice * units) + delivery charge", tpp)

		res1, err := db.Db.Exec(`UPDATE ClusterPurchaseProduct SET status = ?,delivery_charges=?,courier_details=?,total_price=?
                             WHERE id=? AND product_sku = ? AND beneficiary_cluster_name=?`, input_status, input_delivery_charges, input_courier_details1, tpp, input_id, input_product_sku, input_beneficiary_clustername)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Inventory Get Failed1")
			return
		}
		i1, err := res1.RowsAffected()
		if err != nil || i1 == 0 {
			err = errors.New(" Cluster Purchase Product updation failed1")
			return
		}

	} else if input_status == "SUPPLY_COMPLETED" && input_product_qty == "kg-bulk" {

		var input_delivery_charges, tpp string
		var TPP1 float64

		input_delivery_charges = c.Input().Get("input_delivery_charges")

		if input_delivery_charges == "" || input_courier_details1 == "" {
			err = errors.New(" * fields cannot be empty ")
			return
		}

		input_delivery_charges1, err := strconv.ParseFloat(input_delivery_charges, 64)
		if err != nil {
			err = errors.New("Unable to get Delivery Charge")
			return
		}
		input_total_price := c.Input().Get("input_total_price")

		input_total_price1, err := strconv.ParseFloat(input_total_price, 64)
		TPP1 = input_total_price1 + input_delivery_charges1
		tpp = strconv.FormatFloat(TPP1, 'f', 0, 32)

		c.Data["tpp"] = tpp
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "total price = (productprice * units) + delivery charge", tpp)

		res1, err := db.Db.Exec(`UPDATE ClusterPurchaseProduct SET status = ?,delivery_charges=?,courier_details=?,total_price=?
                             WHERE id=? AND product_sku = ? AND beneficiary_cluster_name=?`, input_status, input_delivery_charges, input_courier_details1, tpp, input_id, input_product_sku, input_beneficiary_clustername)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Inventory Get Failed1")
			return
		}
		i1, err := res1.RowsAffected()
		if err != nil || i1 == 0 {
			err = errors.New(" Cluster Purchase Product updation failed1")
			return
		}

	} else if input_status == "SUPPLY_COMPLETED" && input_product_qty == "gms" {

		if input_courier_details1 == "" {
			err = errors.New(" * fields cannot be empty ")
			return
		}

		var input_delivery_charges, tpp string
		var TPP1 float64

		input_delivery_charges, _ = services.GetDeliveryCharge(input_product_package_qty, input_required_units)
		if err != nil {
			beego.Error(err)
			err = errors.New("Unable to get Delivery Charge")
			return
		}

		input_delivery_charges1, err := strconv.ParseFloat(input_delivery_charges, 64)
		if err != nil {

			err = errors.New("Unable to get Delivery Charge")
			return
		}
		input_total_price := c.Input().Get("input_total_price")

		input_total_price1, err := strconv.ParseFloat(input_total_price, 64)
		TPP1 = input_total_price1 + input_delivery_charges1
		tpp = strconv.FormatFloat(TPP1, 'f', 0, 32)

		c.Data["tpp"] = tpp
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "tpp", tpp)

		res1, err := db.Db.Exec(`UPDATE ClusterPurchaseProduct SET status = ?,delivery_charges=?,courier_details=?,total_price=?
                             WHERE id=? AND product_sku = ? AND beneficiary_cluster_name=?`, input_status, input_delivery_charges1, input_courier_details1, tpp, input_id, input_product_sku, input_beneficiary_clustername)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Inventory Get Failed1")
			return
		}
		i1, err := res1.RowsAffected()
		if err != nil || i1 == 0 {
			err = errors.New(" Cluster Purchase Product updation failed1")
			return
		}

	} else if input_status == "DEMAND_ORDER_ACCEPTED" {
		res1, err := db.Db.Exec(`UPDATE ClusterPurchaseProduct SET status = ?
                             WHERE id=? AND product_sku = ? AND beneficiary_cluster_name=?`, input_status, input_id, input_product_sku, input_beneficiary_clustername)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Inventory Get Failed1")
			return
		}
		i1, err := res1.RowsAffected()
		if err != nil || i1 == 0 {
			err = errors.New(" Cluster Purchase Product updation failed1")
			return
		}

	} else {

	}

	if input_status == "SUPPLY_COMPLETED" && skutype == "SOURCE" {

		// If product sku is  present in beneficary cluster

		res2, err := db.Db.Exec(`UPDATE Inventory SET  inventory_available=?
                             WHERE product_sku = ? AND cluster_name=? `, TS1, input_product_sku, input_source_clustername)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Source Inventory Get Failed2")
			return
		}
		i2, err := res2.RowsAffected()
		if err != nil || i2 == 0 {
			err = errors.New(" Source Inventory updation failed1")
			return
		}

		beneficiaryaunits, beneficiarycsprice, beneficiaryreunits, err := services.Getinventoryavalibeunits(input_product_sku, input_beneficiary_clustername)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Inventory Get Failed5")
			return
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "beneficiaryaunits - ", beneficiaryaunits)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "beneficiarycsprice - ", beneficiarycsprice)

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "beneficiaryreunits - ", beneficiaryreunits)

		// beneficiarycsprice1, err := strconv.ParseFloat(beneficiarycsprice, 64)

		// TotalPrice := (beneficiarycsprice1 * runits) + dcharge
		// TotalPrice1 := strconv.FormatFloat(TotalPrice, 'f', 0, 32)

		// c.Data["TotalPrice1"] = TotalPrice1
		// log.Println(beego.AppConfig.String("loglevel"), "Debug", "TotalPrice1 - ", TotalPrice1)

		beneficiaryaunits1, err := strconv.ParseFloat(beneficiaryaunits, 64)

		beneficiaryreunits1, err := strconv.ParseFloat(beneficiaryreunits, 64)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "beneficiaryaunits1 - ", beneficiaryaunits1)
		TB := beneficiaryaunits1 + runits
		TB1 := strconv.FormatFloat(TB, 'f', 0, 32)

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "TB1 - ", TB1)

		TBR := beneficiaryreunits1 + runits
		TBR1 := strconv.FormatFloat(TBR, 'f', 0, 32)

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "TBR1 - ", TBR1)

		res3, err := db.Db.Exec(`UPDATE Inventory SET  inventory_available=?,inventory_received=?
                             WHERE product_sku = ? AND cluster_name=? `, TB1, TBR1, input_product_sku, input_beneficiary_clustername)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Beneficary Inventory Get Failed3")
			return
		}
		i3, err := res3.RowsAffected()
		if err != nil || i3 == 0 {
			err = errors.New(" Beneficary Inventory updation failed2")
			return
		}

	} else if input_status == "SUPPLY_COMPLETED" && skutype == "OTHER" {

		// If product sku is not present in beneficary cluster

		// 1. Update source inventory data

		res2, err := db.Db.Exec(`UPDATE Inventory SET  inventory_available=?
                             WHERE product_sku = ? AND cluster_name=? `, TS1, input_product_sku, input_source_clustername)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Source Inventory Get Failed2")
			return
		}
		i2, err := res2.RowsAffected()
		if err != nil || i2 == 0 {
			err = errors.New(" Source Inventory updation failed1")
			return
		}

		// 2. insert beneficary inventory data

		// row5, err5 := db.Db.Query(`select id from Clusters where cluster_name=? limit 1`, input_source_clustername)
		// if err5 != nil {
		// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		// 	err5 = errors.New("Unable to source cluster id data")
		// 	return
		// }
		// defer sql.Close(row5)
		// _, data5, err := sql.Scan(row5)
		// if err5 != nil {
		// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		// 	err5 = errors.New("Unable to get source cluster id data")
		// 	return
		// }

		row, err := db.Db.Query(`select 
	Lookup_products.id,
	Lookup_products.category_name,
	Lookup_products.subcategory_name,
	Lookup_products.product_name,
	Lookup_products.product_qty,
	Lookup_products.product_quality_mesure,
	Lookup_products.product_shelf_life,
	Lookup_products.product_package_type,
	Lookup_products.product_package_qty,
	Lookup_products.product_origin,
	Lookup_products.krishivala_selling_price,
	Inventory.cluster_buying_price,
	Inventory.cluster_selling_price,
	Inventory.cluster_supply_charge,
	Inventory.cluster_packaging_charge,
	Inventory.cluster_bank_charge,
	Inventory.cluster_discount
	from Lookup_products
	left join Inventory on  Lookup_products.id=Inventory.product_id
	where Lookup_products.product_sku=? and Inventory.cluster_name=? limit 1`, input_product_sku, input_source_clustername)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to product data")
			return
		}
		defer sql.Close(row)
		_, data3, err := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get product data")
			return
		}

		result, err := db.Db.Exec(`INSERT INTO Inventory (
		product_id,
		cluster_name,
		category_name, 
		subcategory_name,
		product_name,
		product_sku,
		product_qty,
		starting_inventory,
		inventory_received, 
		status,
		product_quality_mesure,
		product_shelf_life,
		product_package_type,
		product_package_qty,
		product_origin,
		krishivala_selling_price,
		cluster_buying_price, 
		cluster_selling_price,
		producer_selling_price,
		inventory_available,
		create_by,
		cluster_supply_charge,
		cluster_packaging_charge,
		cluster_bank_charge,
		cluster_discount,
		create_date )
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,now())`,
			data3[0][0],
			input_beneficiary_clustername,
			data3[0][1],
			data3[0][2],
			data3[0][3],
			input_product_sku,
			data3[0][4],
			input_required_units,
			input_required_units,
			input_status,
			data3[0][5],
			data3[0][6],
			data3[0][7],
			data3[0][8],
			data3[0][9],
			data3[0][10],
			data3[0][11],
			data3[0][12],
			"",
			input_required_units,
			id,
			data3[0][13],
			data3[0][14],
			data3[0][15],
			data3[0][16])
		if err != nil {
			err = errors.New("Inventory creation failed")

		}

		i, err := result.RowsAffected()
		if err != nil || i == 0 {
			err = errors.New("Inventory creation failed")

		}

	} else {

		// Do nothing

	}

	return

}
