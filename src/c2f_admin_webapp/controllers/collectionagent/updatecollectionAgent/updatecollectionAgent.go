package updatecollectionAgent

import (
	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/model/utils"
	"c2f_admin_webapp/services"
	"c2f_admin_webapp/session"
	"errors"
	"runtime/debug"
	"strings"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"

	"github.com/astaxie/beego"
)

type Row struct {
	Id           string
	FirstName    string
	MiddleName   string
	LastName     string
	Mobile       string
	Email        string
	Role         string
	Status       string
	Address1     string
	Address2     string
	Town         string
	City         string
	Pincode      string
	Language     string
	LocationType string
	LocationInfo string
}
type Display struct {
	Fields []Field
}
type Field struct {
	Id    string
	Name  string
	Email string
}
type UpdatecollectionAgent struct {
	// beego.Controller
	middleware.BaseController
}

func (c *UpdatecollectionAgent) Get() {
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
			c.TplName = "collectionagent/updatecollectionAgent/updatecollectionAgent.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "collectionagent/updatecollectionAgent/updatecollectionAgent.html"
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

	maptoken := sess.Get("token")
	c.Data["maptoken"] = maptoken
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "maptoken", maptoken)

	clustername := sess.Get("clustername").(string)
	clustername1 := strings.ToUpper(clustername)
	c.Data["clustername"] = clustername1
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "clustername", clustername1)

	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1

	data0, err := services.GetRoles()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Role Failed")
		return
	}
	if len(data0) <= 0 {
		err = errors.New("Roles Not Found")
		return
	}
	var Dis Display
	for i := range data0 {
		var d Field
		d.Id = data0[i][0]
		d.Name = data0[i][1]
		Dis.Fields = append(Dis.Fields, d)
	}
	c.Data["Dis"] = Dis
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)
	auth, err := utils.IsAuthorized(utype, "CollectionAgent")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")
	row, err := db.Db.Query(`select id,first_name,middle_name,last_name,mob_no,email,status,user_type,address1,address2,city,zipcode,state,language,town,location_type,location_info from Collection_Agent WHERE id =?  limit 1`, AdminId)
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
		c.Data["City"] = data[i][10]
		c.Data["Zipcode"] = data[i][11]
		c.Data["State"] = data[i][12]
		c.Data["Language"] = data[i][13]
		c.Data["Town"] = data[i][14]
		c.Data["LocationType"] = data[i][15]
		c.Data["LocationInfo"] = data[i][16]
	}
	return
}
func (c *UpdatecollectionAgent) Post() {
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
			c.TplName = "collectionagent/updatecollectionAgent/updatecollectionAgent.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Collection Agent Updated Successfully"
			c.TplName = "collectionagent/updatecollectionAgent/updatecollectionAgent.html"
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

	auth, err := utils.IsAuthorized(utype, "CollectionAgent")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "adminid - ", AdminId)

	row, err := db.Db.Query(`select id,first_name,middle_name,last_name,mob_no,email,status,user_type,address1,address1,city,zipcode,state,language,town,location_type,location_info from Collection_Agent WHERE id =?  limit 1`, AdminId)
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
		c.Data["City"] = data[i][10]
		c.Data["Zipcode"] = data[i][11]
		c.Data["State"] = data[i][12]
		c.Data["Language"] = data[i][13]
		c.Data["Town"] = data[i][14]
		c.Data["LocationType"] = data[i][15]
		c.Data["LocationInfo"] = data[i][16]
	}

	input_first_name := c.Input().Get("input_first_name")
	input_middle_name := c.Input().Get("input_middle_name")
	input_last_name := c.Input().Get("input_last_name")
	input_mobile := c.Input().Get("input_mobile")
	select_user_type := c.Input().Get("select_user_type")
	input_address1 := c.Input().Get("input_address1")
	input_address2 := c.Input().Get("input_address2")
	input_town := c.Input().Get("input_town")
	input_city := c.Input().Get("input_city")
	input_zipcode := c.Input().Get("input_zipcode")
	input_state := c.Input().Get("input_state")
	input_language := c.Input().Get("input_language")
	LocationType := c.Input().Get("LocationType")
	LocationInfo := c.Input().Get("LocationInfo")
	input_status := c.Input().Get("input_status")
	// if select_role_name == "" {
	//      log.Println(beego.AppConfig.String("loglevel"), "Error", "Blank Passwords")
	//      err = errors.New("Role can't be blank")
	//      return
	// }

	if input_first_name == "" || input_mobile == "" || input_language == "" || LocationType == "" || LocationInfo == "" || input_status == "" ||
		input_address1 == "" || input_address2 == "" || input_town == "" || input_city == "" || input_zipcode == "" || input_state == "" {
		err = errors.New("* fields cannot be empty")
		return
	}
	res, err := db.Db.Exec(`UPDATE Collection_Agent SET 
        first_name = ?, middle_name = ?, last_name = ?, mob_no = ?,status=? , user_type=?,
        address1=?,address2=?,town=?,city=?,zipcode=?,state=?,language=?,location_type=?,location_info=?,update_date=now()
        WHERE id = ?`, input_first_name, input_middle_name, input_last_name, input_mobile, input_status, select_user_type, input_address1,
		input_address2, input_town, input_city, input_zipcode, input_state, input_language, LocationType, LocationInfo, AdminId)
	if err != nil {
		err = errors.New("Agent updation failed")
		return
	}
	i, err := res.RowsAffected()
	if err != nil || i == 0 {
		err = errors.New("Agent updation failed")
		return
	}
	// row, err := db.Db.Query(`select id,first_name,middle_name,last_name,mob_no,email,status,user_type,address1,address2,city,zipcode,state,language,town,location_type,location_info from Collection_Agent WHERE id =?  limit 1`, AdminId)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Unable to get user data")
	// 	return
	// }
	// defer sql.Close(row)
	// _, data, err := sql.Scan(row)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Unable to get user data")
	// 	return
	// }
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	// if len(data) <= 0 {
	// 	err = errors.New("User data not found")
	// 	return
	// }
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", data)
	// for i := range data {
	// 	c.Data["Id"] = data[i][0]
	// 	c.Data["FirstName"] = data[i][1]
	// 	c.Data["MiddleName"] = data[i][2]
	// 	c.Data["LastName"] = data[i][3]
	// 	c.Data["Mobile"] = data[i][4]
	// 	c.Data["Email"] = data[i][5]
	// 	c.Data["UStatus"] = data[i][6]
	// 	c.Data["UserType"] = data[i][7]
	// 	c.Data["Address1"] = data[i][8]
	// 	c.Data["Address2"] = data[i][9]
	// 	c.Data["City"] = data[i][10]
	// 	c.Data["Zipcode"] = data[i][11]
	// 	c.Data["State"] = data[i][12]
	// 	c.Data["Language"] = data[i][13]
	// 	c.Data["Town"] = data[i][14]
	// 	c.Data["LocationType"] = data[i][15]
	// 	c.Data["LocationInfo"] = data[i][16]
	// }
	return
}
