package login

import (
	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/model/utils"

	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"runtime/debug"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/services"
	"c2f_admin_webapp/session"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/encoding/base64"

	"c2f_admin_webapp/model/util/log"
	"c2f_admin_webapp/model/util/pbkdf2"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

type Login struct {
	// beego.Controller
	middleware.BaseController
}

type MapmyIndia struct {
	Access_token string `json:"access_token"`
}

type LoginData struct {
	Uname string `form:"username" valid:"Required"`
	Pass  string `form:"password" valid:"Required;MinSize(6);MaxSize(16)"`
}

func (c *Login) Get() {
	defer func() {
		if l_exception := recover(); l_exception != nil {
			stack := debug.Stack()
			log.Println(beego.AppConfig.String("loglevel"), "Exception", string(stack))
			session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
			c.Data["DisplayMessage"] = "Something went wrong.Please Contact CustomerCare."
			c.TplName = "error/error.html"
		}
		return
	}()
	//	utils.SetHTTPHeader(c.Ctx)

	session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", c.Ctx.Input.IP())
	c.TplName = "login/login.html"
	return
}

func (c *Login) Post() {

	log.Println(beego.AppConfig.String("loglevel"), "Info", "Login Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)

	var err error

	defer func() {

		if err != nil {
			if l_exception := recover(); l_exception != nil {
				stack := debug.Stack()
				log.Println(beego.AppConfig.String("loglevel"), "Exception", string(stack))
				session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
				c.TplName = "error/error.html"
			} else {
				c.Data["DisplayMessage"] = err.Error()
			}
			c.TplName = "login/login.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", " Login Page Fail")
		} else {

			c.TplName = "login/login.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Login Page Success")
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
	defer func() {
		utils.EventLogs(c.Ctx, sess, c.Ctx.Input.Method(), c.Input(), c.Data, err)
		sess.SessionRelease(c.Ctx.ResponseWriter)
	}()
	sessionId := sess.SessionID()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Session ID - ", sessionId)
	var l LoginData
	if err := c.ParseForm(&l); err != nil {
		err = errors.New("Invalid Request Received")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Form Data - ", l)
	c.Data["FormData"] = l
	valid := validation.Validation{}
	b, err := valid.Valid(&l)
	if err != nil {
		err = errors.New("Parameter validation failed")
		return
	}

	if !b {
		for _, err := range valid.Errors {
			log.Println(beego.AppConfig.String("loglevel"), "Debug", err.Key, ":", err.Message, ":", err.Field, ":", err.LimitValue, ":", err.Name, ":", err.Tmpl, ":", err.Value)
		}
		err = errors.New("Invalid Input values")
		return
	}

	if l.Uname == "" || l.Pass == "" {
		log.Println(beego.AppConfig.String("loglevel"), "Error", "Blank User Name")
		err = errors.New("UserName/Password can't be blank.")
		return
	}

	err = session.CheckUserSession(l.Uname)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	}

	err = SearchUser(l.Uname)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("User Not Found/INACTIVE")
		return
	}

	name, id, user_type, language, fullname, err := authinticate(l.Uname, l.Pass)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("User Authentication Failed")
		return
	}

	////////////////map my india auth key ///////////

	url := "https://outpost.mapmyindia.com/api/security/oauth/token"
	method := "POST"

	payload := strings.NewReader("grant_type=client_credentials&client_id=33OkryzDZsL8xDuUWLIfkx9al4FuD4Is6v58GStr2Q7kHQzIH4Jb4oHSeey_avtnmxwnGcY6kPUz4FWSxgGLaKd9FaLagPaq&client_secret=lrFxI-iSEg-RJcrIjsjyvqMssYQ4HVKsMqjP6k_5ro8Qx2SU1_Kd5qQ8kvliiWWhljYzJ2trODBWHizoADQinY_x3J26JuU3dmffgNDhZUE%3D")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
	///fmt.Println(body.access_token)

	var mapindia MapmyIndia
	json.Unmarshal([]byte(string(body)), &mapindia)
	fmt.Printf("token: %s", mapindia.Access_token)

	////////////////map my india auth key ///////////

	sess.Set("uname", l.Uname)
	sess.Set("name", name)
	sess.Set("uid", id)
	sess.Set("user_type", user_type)
	sess.Set("language", language)
	sess.Set("fullname", fullname)
	sess.Set("token", mapindia.Access_token)

	// Add these lines to prevent hijacking
	sess.Set("bind_ip", c.Ctx.Input.IP())
	sess.Set("bind_ua", c.Ctx.Request.UserAgent())

	// Retrieve and print
	bindIP := sess.Get("bind_ip")
	bindUA := sess.Get("bind_ua")

	fmt.Println("Bind IP:", bindIP)
	fmt.Println("Bind UA:", bindUA)

	c_name, err := services.GetMappedClusterforUseratLogin(id, user_type)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Cluster deatils for user")
		return
	}

	if (user_type == "CLUSTERMANAGER" || user_type == "CLUSTEREXECUTIVE") && c_name == "" {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Cluster deatils for user")
		return

	}

	sess.Set("clustername", c_name)

	// if user_type == "CLUSTERMANAGER" {
	// 	data, err := services.GetClusterAsPerMapping(id)
	// 	if err != nil {
	// 		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 		err = errors.New("Cluster user fetch Failed")
	// 		return
	// 	}
	// 	if len(data) <= 0 {
	// 		err = errors.New("Cluster user  Not Found")
	// 		return
	// 	}

	// 	sess.Set("clustername", data[0][3])

	// } else if user_type == "CLUSTEREXECUTIVE" {

	// 	data, err := services.GetClusterAsPerMappingExec(id)
	// 	if err != nil {
	// 		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 		err = errors.New("Cluster user fetch Failed")
	// 		return
	// 	}
	// 	if len(data) <= 0 {
	// 		err = errors.New("Cluster user  Not Found")
	// 		return
	// 	}

	// 	sess.Set("clustername", data[0][3])

	// } else {

	// 	sess.Set("clustername", "")

	// }

	sess.Set("Delhi", beego.AppConfig.String("DELHI_COORDINATES"))
	sess.Set("Andhra", beego.AppConfig.String("ANDHRA_PRADESH_COORDINATES"))

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "uid", id)

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "user_type", user_type)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "MapmyIndia Auth Key ", mapindia.Access_token)

	if user_type == "HOADMIN" {
		sess.Set("uname", l.Uname)
		sess.Set("name", name)
		sess.Set("uid", id)
		sess.Set("user_type", user_type)

		// Add these lines to prevent hijacking
		sess.Set("bind_ip", c.Ctx.Input.IP())
		sess.Set("bind_ua", c.Ctx.Request.UserAgent())

		err = session.SetUserSession(sess.SessionID(), l.Uname)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("System is unable to process your request.Please contact customer care")
			return
		}
		c.Redirect("/Dashboard", 302)
		return

	} else if user_type == "CLUSTERMANAGER" {

		sess.Set("uname", l.Uname)
		sess.Set("name", name)
		sess.Set("uid", id)
		sess.Set("user_type", user_type)

		// Add these lines to prevent hijacking
		sess.Set("bind_ip", c.Ctx.Input.IP())
		sess.Set("bind_ua", c.Ctx.Request.UserAgent())

		err = session.SetUserSession(sess.SessionID(), l.Uname)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("System is unable to process your request.Please contact customer care")
			return
		}
		c.Redirect("/ClusterDashboard", 302)
		return
	} else if user_type == "CLUSTEREXECUTIVE" {

		sess.Set("uname", l.Uname)
		sess.Set("name", name)
		sess.Set("uid", id)
		sess.Set("user_type", user_type)

		// Add these lines to prevent hijacking
		sess.Set("bind_ip", c.Ctx.Input.IP())
		sess.Set("bind_ua", c.Ctx.Request.UserAgent())

		err = session.SetUserSession(sess.SessionID(), l.Uname)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("System is unable to process your request.Please contact customer care")
			return
		}
		c.Redirect("/ClusterExecutiveDashboard", 302)
		return
	} else if user_type == "HOHELPDESK" {

		sess.Set("uname", l.Uname)
		sess.Set("name", name)
		sess.Set("uid", id)
		sess.Set("user_type", user_type)

		// Add these lines to prevent hijacking
		sess.Set("bind_ip", c.Ctx.Input.IP())
		sess.Set("bind_ua", c.Ctx.Request.UserAgent())

		err = session.SetUserSession(sess.SessionID(), l.Uname)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("System is unable to process your request.Please contact customer care")
			return
		}
		c.Redirect("/Dashboard", 302)
		return
	} else if user_type == "HOFINANCE" {

		sess.Set("uname", l.Uname)
		sess.Set("name", name)
		sess.Set("uid", id)
		sess.Set("user_type", user_type)

		// Add these lines to prevent hijacking
		sess.Set("bind_ip", c.Ctx.Input.IP())
		sess.Set("bind_ua", c.Ctx.Request.UserAgent())

		err = session.SetUserSession(sess.SessionID(), l.Uname)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("System is unable to process your request.Please contact customer care")
			return
		}
		c.Redirect("/Dashboard", 302)
		return
	} else if user_type == "REPORTS" {

		sess.Set("uname", l.Uname)
		sess.Set("name", name)
		sess.Set("uid", id)
		sess.Set("user_type", user_type)

		// Add these lines to prevent hijacking
		sess.Set("bind_ip", c.Ctx.Input.IP())
		sess.Set("bind_ua", c.Ctx.Request.UserAgent())

		err = session.SetUserSession(sess.SessionID(), l.Uname)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("System is unable to process your request.Please contact customer care")
			return
		}
		c.Redirect("/Dashboard", 302)
		return
	}
	return
}

func authinticate(uname, pass string) (name, id, user_type, language, fullname string, err error) {

	row, err := db.Db.Query("select password,email,id,role,language,first_name,middle_name,last_name from User where email=? and status='ACTIVE' limit 1", uname)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to authenticate user")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to authenticate user")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "PAVAN")

	if len(data) <= 0 {
		err = errors.New("User not registered")
		return
	}

	cp, err := base64.Decode([]byte(data[0][0]))
	if err != nil {
		err = errors.New("Unable to authenticate user")
		return
	}

	var pbkdf pbkdf2.Pbkdf2
	pbkdf.Itr = 32
	pbkdf.KeyLen = 32
	pbkdf.Plain = []byte(pass)
	pbkdf.Salt = cp[:32]
	pbkdf.Cipher = cp[32:]
	result, err := pbkdf.Compare()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("User password incorrect")
		return
	}
	if !result {
		err = errors.New("User password incorrect")
		return
	}

	name = data[0][1]
	id = data[0][2]
	language = data[0][4]
	fullname = data[0][5]
	// fullname = data[0][5] + " " + data[0][7]

	log.Println(beego.AppConfig.String("loglevel"), "Debug", " data[0][3]", data[0][3])

	row1, err := db.Db.Query("select id,role_name from Lookup_roles where id=?", data[0][3])
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to authenticate user")
		return
	}
	defer sql.Close(row1)
	_, data1, err := sql.Scan(row1)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to authenticate user")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data1))

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "PAVAN")

	if len(data1) <= 0 {
		err = errors.New("User not registered")
		return
	}

	user_type = data1[0][1]
	return
}

func SearchUser(uname string) (err error) {

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", uname)

	row, err := db.Db.Query("SELECT id,email,first_name,status FROM User where email=? limit 1", uname)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System User Not Found")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System User Detail Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("System User Not Found")
		return
	}

	if data[0][3] != "ACTIVE" {
		err = errors.New("User Not Active")
		return
	}

	return
}
