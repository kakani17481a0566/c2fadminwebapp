package customerBulkUpload

import (
	"c2f_admin_webapp/session"

	"c2f_admin_webapp/model/utils"
	"errors"
	"runtime/debug"
	"strings"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/util/sku"

	"c2f_admin_webapp/model/db"

	"github.com/astaxie/beego"

	"encoding/csv"
	"fmt"
	"strconv"

	"log"

	"c2f_admin_webapp/services"

	"c2f_admin_webapp/model/database/sql"

	_ "github.com/go-sql-driver/mysql"

	"c2f_admin_webapp/model/util/password"
)

type CustomerBulkUpload struct {
	// beego.Controller
	middleware.BaseController
}

type CustomerData struct {
	FName  string
	MName  string
	LName  string
	MobNum string
	Email  string
	UType  string

	Add1     string
	Add2     string
	City     string
	Zipcode  string
	State    string
	Language string
	Status   string

	BillColl string
	BillType string
	ColMode  string
	Town     string
	LocType  string
	LocInfo  string
	CluName  string
}

func (c *CustomerBulkUpload) Get() {
	//	AdminId := c.Ctx.Input.Param(":AdminID")
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Customer Bulk Upload Start")
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
			c.TplName = "customer/customerBulkUpload/customerBulkUpload.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Customer Bulk Upload Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "customer/customerBulkUpload/customerBulkUpload.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Customer Bulk Upload Page Success")
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

	auth, err := utils.IsAuthorized(utype, "Customer")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	return
}

func (c *CustomerBulkUpload) Post() {
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
			c.TplName = "customer/customerBulkUpload/customerBulkUpload.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Customer Bulk Upload Page Fail")
		} else {
			c.Data["DisplayMessage"] = "File Uploaded Successfully"
			c.TplName = "customer/customerBulkUpload/customerBulkUpload.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Customer Bulk Upload Page Success")
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

	clustername := sess.Get("clustername").(string)
	clustername1 := strings.ToUpper(clustername)
	c.Data["clustername"] = clustername1
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "clustername", clustername1)

	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1

	id := sess.Get("uid").(string)
	c.Data["id"] = id
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "id", id)

	auth, err := utils.IsAuthorized(utype, "Customer")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	c.TplName = "customerBulkUpload.html"

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

	csvLines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	for _, line := range csvLines[1:] {
		cstmr := CustomerData{

			FName:  line[0],
			MName:  line[1],
			LName:  line[2],
			MobNum: line[3],
			Email:  line[4],

			Add1:     line[5],
			Add2:     line[6],
			City:     line[7],
			Zipcode:  line[8],
			State:    line[9],
			Language: line[10],
			Status:   line[11],

			BillColl: line[12],
			BillType: line[13],
			ColMode:  line[14],
			Town:     line[15],
			LocType:  line[16],
			LocInfo:  line[17],
			CluName:  line[18],
		}

		fmt.Println("First Name :" + cstmr.FName)
		fmt.Println("Middle Name :" + cstmr.MName)
		fmt.Println("Last Name :" + cstmr.LName)
		fmt.Println("Mobile Number :" + cstmr.MobNum)
		fmt.Println("Email :" + cstmr.Email)
		fmt.Println("Address 1 :" + cstmr.Add1)
		fmt.Println("Address 2 :" + cstmr.Add2)
		fmt.Println("City :" + cstmr.City)
		fmt.Println("Zipcode :" + cstmr.Zipcode)
		fmt.Println("State :" + cstmr.State)
		fmt.Println("Language :" + cstmr.Language)
		fmt.Println("Status :" + cstmr.Status)
		//		fmt.Println("Status :" + cstmr.CustType)
		fmt.Println("Bill Collection :" + cstmr.BillColl)
		fmt.Println("Bill Type :" + cstmr.BillType)
		fmt.Println("Collection Mode :" + cstmr.ColMode)
		fmt.Println("Village / Town :" + cstmr.Town)
		fmt.Println("Location Type :" + cstmr.LocType)
		fmt.Println("Location Info :" + cstmr.LocInfo)
		fmt.Println("Cluster Name :" + cstmr.CluName)

		row1, err1 := db.Db.Query(`select count(*) from Customer WHERE mob_no = ?`, cstmr.MobNum)
		if err1 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
			err = errors.New("User doesnot exist" + cstmr.MobNum)
			c.Data["DisplayMessage"] = "User doesnot exist"
			return
		}
		defer sql.Close(row1)
		_, data, err1 := sql.Scan(row1)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
			err = errors.New("User doesnot exist" + cstmr.MobNum)
			c.Data["DisplayMessage"] = "User doesnot exist"
			return
		}

		countlen1 := data[0][0]
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "countlen", countlen1)
		if countlen1 == "0" {
			if cstmr.FName == "" || cstmr.MName == "" || cstmr.LName == "" || cstmr.MobNum == "" || cstmr.Email == "" ||
				cstmr.Add1 == "" || cstmr.Add2 == "" || cstmr.City == "" || cstmr.Zipcode == "" || cstmr.State == "" || cstmr.Language == "" ||
				cstmr.Status == "" || cstmr.BillColl == "" || cstmr.BillType == "" || cstmr.ColMode == "" || cstmr.Town == "" || cstmr.LocType == "" ||
				cstmr.CluName == "" {
				err = errors.New("No fields should be empty")
				return
			}

			var namelen, mobilelen, emaillen, address1len, citylen, clusternamelen int
			var firstname, middlename, lastname, mobile, email, address1, address2, city, pincode, state, language, status, billcoll, billtype, collmode, village, loctype, locinfo, clustername int

			namelen, err = strconv.Atoi(beego.AppConfig.String("FIRST_NAME_LEN"))
			mobilelen, err = strconv.Atoi(beego.AppConfig.String("MOBILE_LEN"))
			clusternamelen, err = strconv.Atoi(beego.AppConfig.String("CLUSTER_LEN"))
			emaillen, err = strconv.Atoi(beego.AppConfig.String("EMAIL_LEN"))
			address1len, err = strconv.Atoi(beego.AppConfig.String("ADD1_LEN"))
			citylen, err = strconv.Atoi(beego.AppConfig.String("CITY_LEN"))

			// total 20 entry in single row

			firstname = len(cstmr.FName)
			if firstname > namelen {
				err = errors.New("First name length not correct for : " + cstmr.FName)
				return
			}
			middlename = len(cstmr.MName)
			if middlename > namelen {
				err = errors.New("Middle name length not correct for : " + cstmr.MName)
				return
			}
			lastname = len(cstmr.LName)
			if lastname > namelen {
				err = errors.New("Last name length not correct for : " + cstmr.LName)
				return
			}
			mobile = len(cstmr.MobNum)
			if mobile != mobilelen {
				err = errors.New("Mobile length not correct for : " + cstmr.MobNum)
				return
			}
			email = len(cstmr.Email)
			if email > emaillen {
				err = errors.New("Email length not correct for : " + cstmr.Email)
				return
			}
			address1 = len(cstmr.Add1)
			if address1 > address1len {
				err = errors.New("Address1 length not correct for : " + cstmr.Add1)
				return
			}
			address2 = len(cstmr.Add2)
			if address2 > address1len {
				err = errors.New("Address2 length not correct for : " + cstmr.Add2)
				return
			}
			city = len(cstmr.City)
			if city > citylen {
				err = errors.New("City length not correct for : " + cstmr.City)
				return
			}
			pincode = len(cstmr.Zipcode)
			if pincode != 6 {
				err = errors.New("Zipcode length not correct for : " + cstmr.Zipcode)
				return
			}
			state = len(cstmr.State)
			if state > citylen {
				err = errors.New("State length not correct for : " + cstmr.State)
				return
			}
			village = len(cstmr.Town)
			if village > citylen {
				err = errors.New("Town length not correct for : " + cstmr.Town)
				return
			}
			language = len(cstmr.Language)
			if language != 2 {
				err = errors.New("Language length not correct for : " + cstmr.Language)
				return
			}
			status = len(cstmr.Status)
			if status > 8 {
				err = errors.New("Status length not correct for : " + cstmr.Status)
				return
			}
			billcoll = len(cstmr.BillColl)
			if billcoll > 10 {
				err = errors.New("Bill Collection length not correct for : " + cstmr.BillColl)
				return
			}
			billtype = len(cstmr.BillType)
			if billtype > 8 {
				err = errors.New("Bill Type length not correct for : " + cstmr.BillType)
				return
			}
			collmode = len(cstmr.ColMode)
			if collmode > 6 {
				err = errors.New("Collection Mode length not correct for : " + cstmr.ColMode)
				return
			}
			loctype = len(cstmr.LocType)
			if loctype > 30 {
				err = errors.New("Location Type length not correct for : " + cstmr.LocType)
				return
			}
			locinfo = len(cstmr.LocInfo)
			if locinfo > 30 {
				err = errors.New("Location Information length not correct for : " + cstmr.LocInfo)
				return
			}
			clustername = len(cstmr.CluName)
			if clustername > clusternamelen {
				err = errors.New("Cluster name length not correct for : " + cstmr.CluName)
				return
			}

			// 	fmt.Println(cstmr.MobNum + "User  exist in table")

			loginPass, _ := password.AlphaNumericSpecial(6)
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "loginPass - ", loginPass)

			pass, err := services.EncryptPassword(loginPass)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				return
			}

			customerid := "C" + sku.SKU(7)

			res, err := db.Db.Exec(`INSERT INTO Customer (first_name,
		password,
		middle_name,
		last_name,
		mob_no,
		email,
		user_type,
		address1,
		address2,
		city,
		zipcode,
		state,
		language,
		status,
		bill_collection,
		bill_type,
		collection_mode,
		customer_id,
		town,
		location_type,
		location_info,
		cluster_name,
		create_date)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,now())`,
				cstmr.FName,
				pass,
				cstmr.MName,
				cstmr.LName,
				cstmr.MobNum,
				cstmr.Email,
				"CONSUMER",
				cstmr.Add1,
				cstmr.Add2,
				cstmr.City,
				cstmr.Zipcode,
				cstmr.State,
				cstmr.Language,
				cstmr.Status,
				cstmr.BillColl,
				cstmr.BillType,
				cstmr.ColMode,
				customerid,
				cstmr.Town,
				cstmr.LocType,
				"",
				cstmr.CluName)

			if err != nil {
				err = errors.New("Customer creation failed for")
				c.Data["DisplayMessage"] = "Customer creation failed"

				return
			}

			i, err := res.RowsAffected()
			if err != nil || i == 0 {
				err = errors.New("Producer creation failed")
				c.Data["DisplayMessage"] = "Customer creation failed"

				return
			}

		} else {
			fmt.Println(cstmr.MobNum, "already exists")
		}
	}

}
