package createCustomerForecast

import (
	"c2f_admin_webapp/session"

	"c2f_admin_webapp/model/utils"
	// "c2f_admin_webapp/services"
	"c2f_admin_webapp/middleware"
	"errors"
	"runtime/debug"
	"strings"

	"c2f_admin_webapp/model/db"

	"github.com/astaxie/beego"

	"encoding/csv"
	"fmt"

	"log"
	"strconv"

	"c2f_admin_webapp/model/database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Row struct {
	Id                string
	CustomerMob       string
	Clustername       string
	Productname       string
	ProductQty        string
	Productunits      string
	Producersellprice string
	From              string
	To                string
	CustomerName      string
	Timestamp         string
	Version           string
	Uploaded_By       string
}

type CreateCustomerForecast struct {
	// beego.Controller
	middleware.BaseController
}

type ForecastData struct {
	CustomerMob  string
	CustomerName string
	Clustername  string
	Productname  string
	ProductQty   string
	Productunits string
	From         string
	To           string
	UploadedBy   string
}

func (c *CreateCustomerForecast) Get() {

	AdminID := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "AdminID", AdminID)

	log.Println(beego.AppConfig.String("loglevel"), "Info", "Create Customer Forecast Start")
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
			c.TplName = "forecast/createCustomerForecast/createCustomerForecast.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Create Customer Forecast Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "forecast/createCustomerForecast/createCustomerForecast.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Create Customer Forecast  Page Success")
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
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "id", id)

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

	auth, err := utils.IsAuthorized(utype, "CustomerForecast")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	return
}

func (c *CreateCustomerForecast) Post() {
	//	AdminId := c.Ctx.Input.Param(":AdminID")
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
			c.TplName = "forecast/createCustomerForecast/createCustomerForecast.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Create Customer Forecast Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Forecast file data uploaded Successfully"
			c.TplName = "forecast/createCustomerForecast/createCustomerForecast.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Create Customer Forecast  Page Success")
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

	auth, err := utils.IsAuthorized(utype, "CustomerForecast")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	// Get the  uploaded Csv File

	file, information, err := c.GetFile("uploadname")

	if err != nil {

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "File - ", "File retrieval failure")
		c.Data["DisplayMessage"] = "File retrieval failure"

		return
	}

	defer file.Close()

	filename := information.Filename
	fmt.Println("Filename: ", filename)

	csvLines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	for _, line := range csvLines[1:] {
		// Skip blank or incomplete rows
		if len(line) < 9 {
			continue
		}

		// Skip rows where all fields are empty (extra blank line)
		allEmpty := true
		for _, field := range line {
			if strings.TrimSpace(field) != "" {
				allEmpty = false
				break
			}
		}
		if allEmpty {
			continue
		}

		pdct := ForecastData{

			CustomerMob:  line[0],
			CustomerName: line[1],
			Clustername:  line[2],
			Productname:  line[3],
			ProductQty:   line[4],
			Productunits: line[5],
			From:         line[6],
			To:           line[7],
			UploadedBy:   line[8],
		}

		fmt.Println("Customer Mobile -" + pdct.CustomerMob)
		fmt.Println("Customer Name -" + pdct.CustomerName)
		fmt.Println("Cluster Name-" + pdct.Clustername)
		fmt.Println("Product Name -" + pdct.Productname)
		fmt.Println("Product Qty -" + pdct.ProductQty)
		fmt.Println("Product Units -" + pdct.Productunits)
		fmt.Println("From -" + pdct.From)
		fmt.Println("To -" + pdct.To)
		fmt.Println("UploadedBy -" + pdct.UploadedBy)

		if pdct.CustomerMob == "" || pdct.Clustername == "" || pdct.Productname == "" || pdct.ProductQty == "" || pdct.Productunits == "" || pdct.From == "" || pdct.To == "" || pdct.UploadedBy == "" {

			err = errors.New("None of the  fields should be empty")
			return
		}

		var mobilelen, namelen, cnamelen, pdctnamelen, pqtylen, punitslen, fromlen, tolen, uploadedbylen int
		var mobile, name, cname, pdctname, pqty, punits, from, to, uploaded int

		mobilelen, err = strconv.Atoi(beego.AppConfig.String("MOBILE_LEN"))
		namelen, err = strconv.Atoi(beego.AppConfig.String("NAME_LEN"))
		cnamelen, err = strconv.Atoi(beego.AppConfig.String("CLUSTER_LEN"))
		pdctnamelen, err = strconv.Atoi(beego.AppConfig.String("PRODUCT_LEN"))
		pqtylen, err = strconv.Atoi(beego.AppConfig.String("PRODUCT_QTY_LEN"))
		punitslen, err = strconv.Atoi(beego.AppConfig.String("PRODUCT_UNITS_LEN"))
		fromlen, err = strconv.Atoi(beego.AppConfig.String("FROM_LEN"))
		tolen, err = strconv.Atoi(beego.AppConfig.String("TO_LEN"))
		uploadedbylen, err = strconv.Atoi(beego.AppConfig.String("UPLOADEDBY_LEN"))

		mobile = len(pdct.CustomerMob)
		if mobile != mobilelen {
			err = errors.New("Mobile length not correct for : " + pdct.CustomerMob)
			return
		}

		name = len(pdct.Clustername)
		if name > namelen {
			err = errors.New("Customer name length not correct for : " + pdct.Clustername)
			return
		}

		cname = len(pdct.Clustername)
		if cname > cnamelen {
			err = errors.New("Clustername length not correct for : " + pdct.Clustername)
			return
		}

		pdctname = len(pdct.Productname)
		if pdctname > pdctnamelen {
			err = errors.New("Productname length not correct for : " + pdct.Productname)
			return
		}

		pqty = len(pdct.ProductQty)
		if pqty > pqtylen {
			err = errors.New("Product Qty length not correct for : " + pdct.ProductQty)
			return
		}

		punits = len(pdct.Productunits)
		if punits > punitslen {
			err = errors.New("Product units length not correct for : " + pdct.Productunits)
			return
		}

		uploaded = len(pdct.UploadedBy)
		if uploaded > uploadedbylen {
			err = errors.New("Product Uploaded By length not correct for : " + pdct.UploadedBy)
			return
		}

		from = len(pdct.From)
		if from > fromlen {
			err = errors.New("From length not correct for : " + pdct.From)
			return
		}

		to = len(pdct.To)
		if to > tolen {
			err = errors.New("To length not correct for : " + pdct.To)
			return
		}

		substr1 := pdct.To[5:7]
		substr2 := pdct.From[5:7]

		if pdct.From[0:4] != pdct.To[0:4] {

			err = errors.New("From Year: " + pdct.From + " and To Year:  " + pdct.To + "should be same - Individual Month Forecast are considered")
			return

		}

		if substr1 != substr2 {

			err = errors.New("From Month: " + pdct.From + " and To Month:  " + pdct.To + "should be same - Individual Month Forecast are considered")
			return

		}

		i, err32 := strconv.Atoi(pdct.From[8:10])
		if err32 != nil {
			fmt.Println(err32)
		}
		j, err32 := strconv.Atoi(pdct.To[8:10])
		if err32 != nil {
			fmt.Println(err32)
		}

		if i > j {

			err = errors.New("From Date: " + pdct.From + " is greater than To Date : " + pdct.To)
			return

		}

		// check if forecast already exist for given product producer cluster from and to dates
		// Step 1: Check if an exact forecast exists for the given criteria.
		row, err := db.Db.Query(`select id FROM CustomerForecast 
		where customer_mob=? and cluster_name=? and product_name=? and start_date=? and end_date=? and product_qty=? and units=?`,
			pdct.CustomerMob, pdct.Clustername, pdct.Productname, pdct.From, pdct.To, pdct.ProductQty, pdct.Productunits)

		if err != nil {
			log.Println("Error querying forecast:", err)
			err = errors.New("Unable to check existing forecast")
			return
		}
		defer sql.Close(row)

		_, existingData, err := sql.Scan(row)
		if err != nil {
			log.Println("Error scanning forecast row:", err)
			err = errors.New("Unable to scan forecast")
			return
		}

		// If an exact forecast already exists, skip insertion.
		if len(existingData) > 0 {
			fmt.Println("Exact forecast already exists for", pdct.Productname, ". Skipping...")
			continue // ✅ go to next row
		}

		// Step 2: Get current max version to decide the new version
		row1, err := db.Db.Query(`SELECT MAX(version) FROM CustomerForecast  
    WHERE customer_mob=? AND cluster_name=? AND product_name=?`,
			pdct.CustomerMob, pdct.Clustername, pdct.Productname)
		if err != nil {
			log.Println("Error fetching version:", err)
			err = errors.New("Unable to fetch forecast version")
			return
		}
		defer sql.Close(row1)

		_, data1, err := sql.Scan(row1)
		if err != nil {
			log.Println("Error scanning version row:", err)
			err = errors.New("Unable to scan forecast version")
			return
		}

		version := 1
		if data1[0][0] != "" {
			versionParsed, err := strconv.Atoi(data1[0][0])
			if err == nil {
				version = versionParsed + 1
			} else {
				log.Println("Version parse error:", err)
			}
		}
		fmt.Println("New version to insert:", version)

		// Step 3: Insert the new forecast with the calculated version

		result, err := db.Db.Exec(`INSERT INTO CustomerForecast ( 
         customer_mob,
         customer_name, 
         cluster_name, 
         product_name, 
         product_qty,
         units, 
         start_date,
         end_date,
         version,
         uploaded_by,
         create_by,
    	 create_date)
         VALUES (?,?,?,?,?,?,?,?,?,?,?,now())`,
			pdct.CustomerMob,
			pdct.CustomerName,
			pdct.Clustername,
			pdct.Productname,
			pdct.ProductQty,
			pdct.Productunits,
			pdct.From,
			pdct.To,
			version,
			pdct.UploadedBy,
			id)
		if err != nil {
			err = errors.New("Customer Forecast creation failed")
			return
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil || rowsAffected == 0 {
			err = errors.New("Customer Forecast creation failed")
			return
		}
		fmt.Println("Inserted new forecast, version:", version)

	}
}
