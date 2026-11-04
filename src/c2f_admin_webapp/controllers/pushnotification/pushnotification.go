package pushNotification

import (
	"c2f_admin_webapp/session"

	"runtime/debug"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"
	"errors"

	"fmt"
	"strings"

	// "c2f_admin_webapp/model/db"
	"c2f_admin_webapp/model/fcm"
	"c2f_admin_webapp/services"

	"github.com/astaxie/beego"

	// "c2f_admin_webapp/model/database/sql"
	"c2f_admin_webapp/model/util/log"
	//	"c2f_admin_webapp/model/util/sku"
)

type PushNotification struct {
	// beego.Controller
	middleware.BaseController
}

func (c *PushNotification) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Add Assets Start")
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
			c.TplName = "pushnotification/pushnotification.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Push Notification Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "pushnotification/pushnotification.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Push Notification Page Success")
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

	language := sess.Get("language").(string)
	c.Data["language"] = language
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "language", language)

	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1

	id := sess.Get("uid").(string)
	c.Data["id"] = id

	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "utype", utype)

	auth, err := utils.IsAuthorized(utype, "PushNotification")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	return
}
func (c *PushNotification) Post() {

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
			c.TplName = "pushnotification/pushnotification.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Push Notification Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Notification sent Successfully"
			c.TplName = "pushnotification/pushnotification.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Push Notification  Page Success")
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

	language := sess.Get("language").(string)
	c.Data["language"] = language
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "language", language)

	fullname := sess.Get("fullname").(string)
	fullname1 := strings.ToUpper(fullname)
	c.Data["fullname1"] = fullname1

	id := sess.Get("uid").(string)
	c.Data["id"] = id
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "id", id)

	auth, err := utils.IsAuthorized(utype, "PushNotification")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	// if input_category_name == "" || input_category_status == "" {
	// 	err = errors.New("category name and category status cannot be empty")
	// 	return
	// }

	var AllCustomerTokenArray []string

	input_notification_title := c.Input().Get("input_notification_title")
	input_notification_body := c.Input().Get("input_notification_body")
	input_notification_image := c.Input().Get("input_notification_image")

	data1, err := services.GetAllCustomerTokens()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Token data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Token Data - ", data1, "\nData len - ", len(data1))
	if len(data1) <= 0 {
		err = errors.New("Token data not found")
		return
	}
	c.Data["Total Tokens"] = len(data1)

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data1)

	for i := range data1 {
		token := data1[i][0]
		log.Println(beego.AppConfig.String("loglevel"), "Debug", fmt.Sprintf("Token [%d] = %s", i, token))
		AllCustomerTokenArray = append(AllCustomerTokenArray, token)
	}

	// var AllCustomerTokenArray []string

	// Hardcoded token for testing
	// AllCustomerTokenArray = []string{
	// 	"fMXzvx1aTsqxxm7Zo67g3M:APA91bH3C5OlUa5KvgxZU-hs2uv_V65GP--dyD68G2gkUchpwrX7FCHb93yR-r-cuax71NtQF95QH5tJsaoD_468vkJXsRNmbeXBi7tjnc-z7t02sQccYuQ",
	// }

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Title - ", input_notification_title)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Body - ", input_notification_body)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Image URL - ", input_notification_image)

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Array of Customer Tokens - ", AllCustomerTokenArray)

	// err = fcm.SendNotificationToAll(AllCustomerTokenArray, input_notification_title, input_notification_body, input_notification_image)
	// if err != nil {
	// 	err = errors.New("Notification Send Failed")
	// 	return
	// }

	err = fcm.SendNotificationToAll(AllCustomerTokenArray, input_notification_title, input_notification_body, input_notification_image)
	if err != nil {
		log.Println("Error occurred:", err.Error())
		return
	}

	return
}
