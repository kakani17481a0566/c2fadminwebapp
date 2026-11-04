package updateUser

import (
	"c2f_admin_webapp/session"

	"runtime/debug"

	"c2f_admin_webapp/model/utils"
	"errors"

	"strings"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/db"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"
	"c2f_admin_webapp/model/util/log"
)

type Row struct {
	Id         string
	FirstName  string
	MiddleName string
	LastName   string
	Mobile     string
	Email      string
	Usertype   string
	Status     string
}

type UpdateUser struct {
	// beego.Controller
	middleware.BaseController
}

func (c *UpdateUser) Get() {
	AdminId := c.Ctx.Input.Param(":AdminID")
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
			c.TplName = "users/updateUsers/updateUsers.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "users/updateUsers/updateUsers.html"
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

	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1

	maptoken := sess.Get("token")
	c.Data["maptoken"] = maptoken
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "maptoken", maptoken)

	auth, err := utils.IsAuthorized(utype, "User")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	row, err := db.Db.Query(`select 
	id,
	first_name,
	middle_name,
	last_name,
	mob_no,
	email,
	status,
	user_type,
	address1,
	address2,
	town,
	city,
	zipcode,
	state
	from Customer WHERE id =?  limit 1`, AdminId)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("User data not found")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	for i := range data {

		c.Data["Id"] = data[i][0]
		c.Data["FirstName"] = data[i][1]
		c.Data["MiddleName"] = data[i][2]
		c.Data["LastName"] = data[i][3]
		c.Data["Mobile"] = data[i][4]
		c.Data["Email"] = data[i][5]
		c.Data["UStatus"] = data[i][6]
		c.Data["UserType"] = data[i][7]
		c.Data["Address1"] = data[i][8]
		c.Data["Address2"] = data[i][9]
		c.Data["Town"] = data[i][10]
		c.Data["City"] = data[i][11]
		c.Data["Zipcode"] = data[i][12]
		c.Data["State"] = data[i][13]

	}

	return

}

func (c *UpdateUser) Post() {

	AdminId := c.Ctx.Input.Param(":AdminID")

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
			c.TplName = "users/updateUsers/updateUsers.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Producer Updated Successfully"
			c.TplName = "users/updateUsers/updateUsers.html"
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

	auth, err := utils.IsAuthorized(utype, "User")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	maptoken := sess.Get("token")
	c.Data["maptoken"] = maptoken
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "maptoken", maptoken)

	row1, err := db.Db.Query(`select id,first_name,middle_name,last_name,mob_no,email,status,user_type from Customer WHERE id =?  limit 1`, AdminId)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	defer sql.Close(row1)
	_, data1, err := sql.Scan(row1)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data1, "\nData len - ", len(data1))
	if len(data1) <= 0 {
		err = errors.New("User data not found")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data1)

	for i := range data1 {

		c.Data["Id"] = data1[i][0]
		c.Data["FirstName"] = data1[i][1]
		c.Data["MiddleName"] = data1[i][2]
		c.Data["LastName"] = data1[i][3]
		c.Data["Mobile"] = data1[i][4]
		c.Data["Email"] = data1[i][5]
		c.Data["UStatus"] = data1[i][6]
		c.Data["UserType"] = data1[i][7]

	}

	input_first_name := c.Input().Get("input_first_name")
	input_middle_name := c.Input().Get("input_middle_name")
	input_last_name := c.Input().Get("input_last_name")
	input_email := c.Input().Get("input_email")
	input_mobile := c.Input().Get("input_mobile")
	input_user_type := c.Input().Get("input_user_type")
	input_user_status := c.Input().Get("input_user_status")
	input_location_info := c.Input().Get("input_location_info")
	input_location_type := c.Input().Get("input_location_type")
	input_address1 := c.Input().Get("input_address1")
	input_address2 := c.Input().Get("input_address2")
	input_town := c.Input().Get("input_town")
	input_city := c.Input().Get("input_city")
	input_zipcode := c.Input().Get("input_zipcode")
	input_state := c.Input().Get("input_state")

	if input_user_status == "" {
		input_user_status = "ACTIVE"
	}

	if input_first_name == "" || input_user_status == "" || input_email == "" || input_location_info == "" || input_location_type == "" ||
		input_mobile == "" || input_user_type == "" || input_address1 == "" || input_address2 == "" || input_town == "" || input_city == "" ||
		input_zipcode == "" || input_state == "" {
		err = errors.New("* fields cannot be empty")
		return
	}

	res, err := db.Db.Exec(`UPDATE Customer SET 
	first_name = ?, middle_name = ?, last_name = ?, email = ?, mob_no = ?,status=? , user_type=?,location_info=?,location_type=?,
	address1 = ?, address2 = ?, city = ?,zipcode = ?, state= ?, town = ?
	WHERE id = ?`, input_first_name, input_middle_name, input_last_name, input_email, input_mobile, input_user_status, input_user_type, input_location_info, input_location_type,
		input_address1, input_address2, input_city, input_zipcode, input_state, input_town, AdminId)
	if err != nil {
		err = errors.New("User updation failed")
		return
	}

	i, err := res.RowsAffected()
	if err != nil || i == 0 {
		err = errors.New("Producer updation failed due to no new changes found")
		return
	}

	row, err := db.Db.Query(`select id,first_name,middle_name,last_name,mob_no,email,status,user_type from Customer WHERE id = ? limit 1`, AdminId)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get user data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("User data not found")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	for i := range data {

		c.Data["Id"] = data[i][0]
		c.Data["FirstName"] = data[i][1]
		c.Data["MiddleName"] = data[i][2]
		c.Data["LastName"] = data[i][3]
		c.Data["Mobile"] = data[i][4]
		c.Data["Email"] = data[i][5]
		c.Data["UStatus"] = data[i][6]
		c.Data["UserType"] = data[i][7]

	}

	return
}
