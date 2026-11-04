package producerBulkUpload

import (
	"c2f_admin_webapp/session"

	"c2f_admin_webapp/model/utils"
	"errors"
	"runtime/debug"
	"strconv"
	"strings"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/util/sku"

	"c2f_admin_webapp/model/db"

	"github.com/astaxie/beego"

	"encoding/csv"
	"fmt"

	"log"

	"c2f_admin_webapp/services"

	"c2f_admin_webapp/model/database/sql"

	_ "github.com/go-sql-driver/mysql"

	"c2f_admin_webapp/model/util/password"
)

type ProducerBulkUpload struct {
	// beego.Controller
	middleware.BaseController
}

type ProducerData struct {
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

	CustType string

	BillColl string
	BillType string
	ColMode  string
	Town     string
	LocType  string
	LocInfo  string
	CluName  string
}

func (c *ProducerBulkUpload) Get() {
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
			c.TplName = "users/producerBulkUpload/producerBulkUpload.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Producer Bulk Upload Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "users/producerBulkUpload/producerBulkUpload.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Producer Bulk Upload Page Success")
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

	auth, err := utils.IsAuthorized(utype, "User")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	return
}

func (c *ProducerBulkUpload) Post() {
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
			c.TplName = "users/producerBulkUpload/producerBulkUpload.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Producer Bulk Upload Page Fail")
		} else {
			c.Data["DisplayMessage"] = "File Uploaded Successfully"
			c.TplName = "users/producerBulkUpload/producerBulkUpload.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Producer Bulk Upload Page Success")
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

	auth, err := utils.IsAuthorized(utype, "User")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	c.TplName = "producerBulkUpload.html"

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
		prdcr := ProducerData{

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
			CustType: line[12],

			BillColl: line[13],
			BillType: line[14],
			ColMode:  line[15],
			Town:     line[16],
			LocType:  line[17],
			LocInfo:  line[18],
			CluName:  line[19],
		}

		fmt.Println("First Name :" + prdcr.FName)
		fmt.Println("Middle Name :" + prdcr.MName)
		fmt.Println("Last Name :" + prdcr.LName)
		fmt.Println("Mobile Number :" + prdcr.MobNum)
		fmt.Println("Email :" + prdcr.Email)
		fmt.Println("Address 1 :" + prdcr.Add1)
		fmt.Println("Address 2 :" + prdcr.Add2)
		fmt.Println("City :" + prdcr.City)
		fmt.Println("Zipcode :" + prdcr.Zipcode)
		fmt.Println("State :" + prdcr.State)
		fmt.Println("Language :" + prdcr.Language)
		fmt.Println("Status :" + prdcr.Status)
		fmt.Println("CustType :" + prdcr.CustType)
		fmt.Println("Bill Collection :" + prdcr.BillColl)
		fmt.Println("Bill Type :" + prdcr.BillType)
		fmt.Println("Collection Mode :" + prdcr.ColMode)
		fmt.Println("Village / Town :" + prdcr.Town)
		fmt.Println("Location Type :" + prdcr.LocType)
		fmt.Println("Location Info :" + prdcr.LocInfo)
		fmt.Println("Cluster Name :" + prdcr.CluName)

		row1, err1 := db.Db.Query(`select count(*) from Customer WHERE mob_no = ?`, prdcr.MobNum)
		if err1 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
			err = errors.New("User doesnot exist" + prdcr.MobNum)
			c.Data["DisplayMessage"] = "User doesnot exist"
			return
		}
		defer sql.Close(row1)
		_, data, err1 := sql.Scan(row1)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err1)
			err = errors.New("User doesnot exist" + prdcr.MobNum)
			c.Data["DisplayMessage"] = "User doesnot exist"
			return
		}

		countlen1 := data[0][0]
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "countlen", countlen1)
		if countlen1 == "0" {
			if prdcr.FName == "" || prdcr.MName == "" || prdcr.LName == "" || prdcr.MobNum == "" || prdcr.Email == "" || prdcr.Add1 == "" || prdcr.Add2 == "" ||
				prdcr.City == "" || prdcr.Zipcode == "" || prdcr.State == "" || prdcr.Language == "" || prdcr.Status == "" || prdcr.CustType == "" || prdcr.BillColl == "" || prdcr.BillType == "" ||
				prdcr.ColMode == "" || prdcr.Town == "" || prdcr.LocType == "" || prdcr.CluName == "" {
				err = errors.New("No fields should be empty")
				return
			}

			var namelen, mobilelen, emaillen, address1len, citylen, clusternamelen int
			var firstname, middlename, lastname, mobile, email, address1, address2, city, pincode, state, language, status, producertype, billcoll, billtype, collmode, village, loctype, locinfo, clustername int

			namelen, err = strconv.Atoi(beego.AppConfig.String("FIRST_NAME_LEN"))
			mobilelen, err = strconv.Atoi(beego.AppConfig.String("MOBILE_LEN"))
			clusternamelen, err = strconv.Atoi(beego.AppConfig.String("CLUSTER_LEN"))
			emaillen, err = strconv.Atoi(beego.AppConfig.String("EMAIL_LEN"))
			address1len, err = strconv.Atoi(beego.AppConfig.String("ADD1_LEN"))
			citylen, err = strconv.Atoi(beego.AppConfig.String("CITY_LEN"))

			// total 20 entry in single row

			firstname = len(prdcr.FName)
			if firstname > namelen {
				err = errors.New("First name length not correct for : " + prdcr.FName)
				return
			}
			middlename = len(prdcr.MName)
			if middlename > namelen {
				err = errors.New("Middle name length not correct for : " + prdcr.MName)
				return
			}
			lastname = len(prdcr.LName)
			if lastname > namelen {
				err = errors.New("Last name length not correct for : " + prdcr.LName)
				return
			}
			mobile = len(prdcr.MobNum)
			if mobile != mobilelen {
				err = errors.New("Mobile length not correct for : " + prdcr.MobNum)
				return
			}
			email = len(prdcr.Email)
			if email > emaillen {
				err = errors.New("Email length not correct for : " + prdcr.Email)
				return
			}
			address1 = len(prdcr.Add1)
			if address1 > address1len {
				err = errors.New("Address1 length not correct for : " + prdcr.Add1)
				return
			}
			address2 = len(prdcr.Add2)
			if address2 > address1len {
				err = errors.New("Address2 length not correct for : " + prdcr.Add2)
				return
			}
			city = len(prdcr.City)
			if city > citylen {
				err = errors.New("City length not correct for : " + prdcr.City)
				return
			}
			pincode = len(prdcr.Zipcode)
			if pincode != 6 {
				err = errors.New("Zipcode length not correct for : " + prdcr.Zipcode)
				return
			}
			state = len(prdcr.State)
			if state > citylen {
				err = errors.New("State length not correct for : " + prdcr.State)
				return
			}
			village = len(prdcr.Town)
			if village > citylen {
				err = errors.New("Town length not correct for : " + prdcr.Town)
				return
			}
			language = len(prdcr.Language)
			if language != 2 {
				err = errors.New("Language length not correct for : " + prdcr.Language)
				return
			}
			status = len(prdcr.Status)
			if status > 8 {
				err = errors.New("Status length not correct for : " + prdcr.Status)
				return
			}
			producertype = len(prdcr.CustType)
			if producertype > 20 {
				err = errors.New("Producer Type not correct for : " + prdcr.CustType)
				return
			}
			billcoll = len(prdcr.BillColl)
			if billcoll > 10 {
				err = errors.New("Bill Collection length not correct for : " + prdcr.BillColl)
				return
			}
			billtype = len(prdcr.BillType)
			if billtype > 8 {
				err = errors.New("Bill Type length not correct for : " + prdcr.BillType)
				return
			}
			collmode = len(prdcr.ColMode)
			if collmode > 6 {
				err = errors.New("Collection Mode length not correct for : " + prdcr.ColMode)
				return
			}
			loctype = len(prdcr.LocType)
			if loctype > 30 {
				err = errors.New("Location Type length not correct for : " + prdcr.LocType)
				return
			}
			locinfo = len(prdcr.LocInfo)
			if locinfo > 30 {
				err = errors.New("Location Information length not correct for : " + prdcr.LocInfo)
				return
			}
			clustername = len(prdcr.CluName)
			if clustername > clusternamelen {
				err = errors.New("Cluster name length not correct for : " + prdcr.CluName)
				return
			}

			loginPass, _ := password.AlphaNumericSpecial(6)
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "loginPass - ", loginPass)

			pass, err := services.EncryptPassword(loginPass)
			if err != nil {
				log.Println(beego.AppConfig.String("loglevel"), "Error", err)
				return
			}

			customerid := "P" + sku.SKU(7)

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
		customer_type,
		bill_collection,
		bill_type,
		collection_mode,
		customer_id,
		town,
		location_type,
		location_info,
		cluster_name,
		create_date)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,now())`,
				prdcr.FName,
				pass,
				prdcr.MName,
				prdcr.LName,
				prdcr.MobNum,
				prdcr.Email,
				"PRODUCER",
				prdcr.Add1,
				prdcr.Add2,
				prdcr.City,
				prdcr.Zipcode,
				prdcr.State,
				prdcr.Language,
				prdcr.Status,
				prdcr.CustType,
				prdcr.BillColl,
				prdcr.BillType,
				prdcr.ColMode,
				customerid,
				prdcr.Town,
				prdcr.LocType,
				"",
				prdcr.CluName)

			if err != nil {
				err = errors.New("Producer creation failed for")
				c.Data["DisplayMessage"] = "Producer creation failed"

				return
			}

			i, err := res.RowsAffected()
			if err != nil || i == 0 {
				err = errors.New("Producer creation failed")
				c.Data["DisplayMessage"] = "Producer creation failed"

				return
			}

		} else {
			fmt.Println(prdcr.MobNum, "already exists")
		}
	}

}
