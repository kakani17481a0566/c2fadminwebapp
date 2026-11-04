package forecastdashboard

import (
	"c2f_admin_webapp/model/db"
	// "strings"

	"c2f_admin_webapp/model/utils"
	// "c2f_admin_webapp/services"
	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/session"
	"errors"

	"fmt"
	//"html/template"

	"runtime/debug"

	"time"

	"strconv"

	"c2f_admin_webapp/model/database/sql"

	"github.com/astaxie/beego"
	log "github.com/sirupsen/logrus"

	"encoding/json"
	// "c2f_admin_webapp/model/util/log"
)

type Display struct {
	Fields []Field
}
type Field struct {
	ProducerForecastProduct string
	ProductQty              string
	Version                 string
	CreateDate              string
}

type ForecastDashboard struct {
	// beego.Controller
	middleware.BaseController
}

type CombinedForecastDashboard struct {
	beego.Controller
}

type ForecastData struct {
	Month1      string `json:"month1"`
	Month2      string `json:"month2"`
	Month3      string `json:"month3"`
	Month4      string `json:"month4"`
	Month5      string `json:"month5"`
	Month6      string `json:"month6"`
	Month7      string `json:"month7"`
	Month8      string `json:"month8"`
	Month9      string `json:"month9"`
	Month10     string `json:"month10"`
	Month11     string `json:"month11"`
	Month12     string `json:"month12"`
	PM1         string `json:"pm1"`
	PM2         string `json:"pm2"`
	PM3         string `json:"pm3"`
	PM4         string `json:"pm4"`
	PM5         string `json:"pm5"`
	PM6         string `json:"pm6"`
	PM7         string `json:"pm7"`
	PM8         string `json:"pm8"`
	PM9         string `json:"pm9"`
	PM10        string `json:"pm10"`
	PM11        string `json:"pm11"`
	PM12        string `json:"pm12"`
	CM1         string `json:"cm1"`
	CM2         string `json:"cm2"`
	CM3         string `json:"cm3"`
	CM4         string `json:"cm4"`
	CM5         string `json:"cm5"`
	CM6         string `json:"cm6"`
	CM7         string `json:"cm7"`
	CM8         string `json:"cm8"`
	CM9         string `json:"cm9"`
	CM10        string `json:"cm10"`
	CM11        string `json:"cm11"`
	CM12        string `json:"cm12"`
	Product1    string `json:"product1"`
	Product2    string `json:"product2"`
	Product1Qty string `json:"product1qty"`
	Product2QTY string `json:"product2qty"`
}

func (c *ForecastDashboard) Get() {
	log.Println(beego.AppConfig.String("loglevel"), "Info", "ForecastDashboard Page Start")
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
			log.Println(beego.AppConfig.String("loglevel"), "Info", "ForecastDashboard Page Fail")
		} else if err != nil {
			c.Data["DisplayMessage"] = err.Error()
			c.TplName = "error/error.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "ForecastDashboard Page Fail")
		} else {
			c.TplName = "analytics/forecastdashboard.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "ForecastDashboard Page Success")
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

	language := sess.Get("language").(string)
	c.Data["language"] = language
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "language", language)

	auth, err := utils.IsAuthorized(utype, "ForecastDashboard")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	//-----------------Start Customer Forecast Data------------------------

	custrow, err := db.Db.Query(`SELECT distinct(product_name),product_qty FROM CustomerForecast`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Customer Forecasted Products")
		return
	}
	defer sql.Close(custrow)

	_, custdata, err := sql.Scan(custrow)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Customer Forecasted Products")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", custdata, "\nData len - ", len(custdata))

	var Dis1 Display
	for i := range custdata {
		var d Field
		d.ProducerForecastProduct = custdata[i][0]
		d.ProductQty = custdata[i][1]

		Dis1.Fields = append(Dis1.Fields, d)
	}
	c.Data["Dis1"] = Dis1
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Customer Forecasted Products - ", Dis1)

	row11, err := db.Db.Query(`SELECT version, MAX(create_date) AS create_date
FROM CustomerForecast
GROUP BY version`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Producer Forecasted Products")
		return
	}
	defer sql.Close(row11)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row11)
	_, data11, err := sql.Scan(row11)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Producer Forecasted Products")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data11, "\nData len - ", len(data11))

	var Dis11 Display
	for i := range data11 {
		var d Field
		d.Version = data11[i][0]
		d.CreateDate = data11[i][1]

		Dis11.Fields = append(Dis11.Fields, d)
	}
	c.Data["Dis11"] = Dis11
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Customer Forecasted Version - ", Dis11)

	// c.Data["CVersion"] = data11[0][0]

	if len(data11) > 0 && len(data11[0]) > 0 {
		c.Data["CVersion"] = data11[0][0]
	} else {
		log.Println("Warning: data11 is empty, cannot set CVersion")
		c.Data["CVersion"] = "" // or a default value
	}

	//-----------------End Customer Forecast Data------------------------

	//-----------------Start Producer Forecast Data------------------------
	row, err := db.Db.Query(`SELECT distinct(product_name),product_qty,version FROM ProducerForecast`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Producer Forecasted Products")
		return
	}
	defer sql.Close(row)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Producer Forecasted Products")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))

	var Dis Display
	for i := range data {
		var d Field
		d.ProducerForecastProduct = data[i][0]
		d.ProductQty = data[i][1]
		d.Version = data[i][2]

		Dis.Fields = append(Dis.Fields, d)
	}
	c.Data["Dis"] = Dis
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Producer Forecasted Products - ", Dis)

	row12, err := db.Db.Query(`SELECT version, MAX(create_date) AS create_date
FROM ProducerForecast
GROUP BY version`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Producer Forecasted Products")
		return
	}
	defer sql.Close(row12)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row12)
	_, data12, err := sql.Scan(row12)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Producer Forecasted Products")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data12, "\nData len - ", len(data12))

	var Dis12 Display
	for i := range data12 {
		var d Field
		d.Version = data12[i][0]
		d.CreateDate = data12[i][1]

		Dis12.Fields = append(Dis12.Fields, d)
	}
	c.Data["Dis12"] = Dis12
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Producer Forecasted Version - ", Dis12)

	// c.Data["PVersion"] = data12[0][0]

	if len(data12) > 0 && len(data12[0]) > 0 {
		c.Data["PVersion"] = data12[0][0]
	} else {
		log.Println("Warning: data12 is empty, cannot set PVersion")
		c.Data["PVersion"] = "" // or a default value
	}
	//-----------------End Producer Forecast Data------------------------

	//-----------------Start Combined Producer and Customer Forecast Data------------------------

	row2, err := db.Db.Query(`(SELECT distinct(product_name),product_qty FROM ProducerForecast)UNION(SELECT distinct(product_name),product_qty FROM CustomerForecast)`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Producer Forecasted Products")
		return
	}
	defer sql.Close(row2)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row)
	_, data2, err := sql.Scan(row2)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Producer Forecasted Products")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))

	var Dis2 Display
	for i := range data2 {
		var d Field
		d.ProducerForecastProduct = data2[i][0]
		d.ProductQty = data2[i][1]

		Dis2.Fields = append(Dis2.Fields, d)
	}
	c.Data["Dis3"] = Dis2
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Combined Forecasted Products - ", Dis2)

	row13, err := db.Db.Query(`SELECT version, MIN(create_date) AS create_date FROM (
    SELECT DISTINCT version, create_date FROM ProducerForecast
    UNION
    SELECT DISTINCT version, create_date FROM CustomerForecast) AS combined
	GROUP BY version`)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Producer Forecasted Products")
		return
	}
	defer sql.Close(row13)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row13)
	_, data13, err := sql.Scan(row13)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Producer Forecasted Products")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data13, "\nData len - ", len(data13))

	var Dis13 Display
	for i := range data13 {
		var d Field
		d.Version = data13[i][0]
		d.CreateDate = data13[i][1]

		Dis13.Fields = append(Dis13.Fields, d)
	}
	c.Data["Dis13"] = Dis13
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Combined Forecasted Products - ", Dis13)

	// c.Data["PCVersion"] = data13[0][0]

	if len(data13) > 0 && len(data13[0]) > 0 {
		c.Data["PCVersion"] = data13[0][0]
	} else {
		log.Println("Warning: data13 is empty, cannot set PCVersion")
		c.Data["PCVersion"] = "" // or a default value
	}

	//-----------------End Combined Producer and Customer Forecast Data------------------------

	currentTime := time.Now()

	// get 12 months name and year from current time

	m1year := strconv.Itoa(currentTime.Year())
	month1 := time.Month(currentTime.Month()).String()
	fmt.Println("Month 1 and Year:", month1, m1year)
	c.Data["month1"] = month1 + " " + m1year

	m2 := currentTime.AddDate(0, 0, 30)
	m2year := strconv.Itoa(m2.Year())
	month2 := time.Month(m2.Month()).String()
	fmt.Println("Month 2 and Year:", month2, m2year)
	c.Data["month2"] = month2 + " " + m2year

	m3 := currentTime.AddDate(0, 0, 60)
	m3year := strconv.Itoa(m3.Year())
	month3 := time.Month(m3.Month()).String()
	fmt.Println("Month 3 and Year:", month3, m3year)
	c.Data["month3"] = month3 + " " + m3year

	m4 := currentTime.AddDate(0, 0, 90)
	m4year := strconv.Itoa(m4.Year())
	month4 := time.Month(m4.Month()).String()
	fmt.Println("Month 4 and Year:", month4, m4year)
	c.Data["month4"] = month4 + " " + m4year

	m5 := currentTime.AddDate(0, 0, 120)
	m5year := strconv.Itoa(m5.Year())
	month5 := time.Month(m5.Month()).String()
	fmt.Println("Month 5 and Year:", month5, m5year)
	c.Data["month5"] = month5 + " " + m5year

	m6 := currentTime.AddDate(0, 0, 150)
	m6year := strconv.Itoa(m6.Year())
	month6 := time.Month(m6.Month()).String()
	fmt.Println("Month 6 and Year:", month6, m6year)
	c.Data["month6"] = month6 + " " + m6year

	m7 := currentTime.AddDate(0, 0, 180)
	m7year := strconv.Itoa(m7.Year())
	month7 := time.Month(m7.Month()).String()
	fmt.Println("Month 7 and Year:", month7, m7year)
	c.Data["month7"] = month7 + " " + m7year

	m8 := currentTime.AddDate(0, 0, 210)
	m8year := strconv.Itoa(m8.Year())
	month8 := time.Month(m8.Month()).String()
	fmt.Println("Month 8 and Year:", month8, m8year)
	c.Data["month8"] = month8 + " " + m8year

	m9 := currentTime.AddDate(0, 0, 240)
	m9year := strconv.Itoa(m9.Year())
	month9 := time.Month(m9.Month()).String()
	fmt.Println("Month 9 and Year:", month9, m9year)
	c.Data["month9"] = month9 + " " + m9year

	m10 := currentTime.AddDate(0, 0, 270)
	m10year := strconv.Itoa(m10.Year())
	month10 := time.Month(m10.Month()).String()
	fmt.Println("Month 10 and Year:", month10, m10year)
	c.Data["month10"] = month10 + " " + m10year

	m11 := currentTime.AddDate(0, 0, 300)
	m11year := strconv.Itoa(m11.Year())
	month11 := time.Month(m11.Month()).String()
	fmt.Println("Month 11 and Year:", month11, m11year)
	c.Data["month11"] = month11 + " " + m11year

	m12 := currentTime.AddDate(0, 0, 330)
	m12year := strconv.Itoa(m12.Year())
	month12 := time.Month(m12.Month()).String()
	fmt.Println("Month 12 and Year:", month12, m12year)
	c.Data["month12"] = month12 + " " + m12year

	return
}

func (c *ForecastDashboard) Post() {

	var payload string

	log.Println(beego.AppConfig.String("loglevel"), "Info", "ForecastDashboard Page Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	var err error
	var Autherr error
	defer func() {
		if l_exception := recover(); l_exception != nil {
			c.Ctx.Output.Body([]byte(`{"Tittle":"ERROR","Message": "` + err.Error() + `","Type": "error"}`))
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
			log.Println(beego.AppConfig.String("loglevel"), "Info", "ForecastDashboard Page Fail")
		} else if err != nil {
			c.Ctx.Output.Body([]byte(`{"Tittle":"ERROR","Message": "` + err.Error() + `","Type": "error"}`))
			c.Data["DisplayMessage"] = err.Error()
			c.TplName = "error/error.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "ForecastDashboard Page Fail")
		} else {
			c.Ctx.Output.Body([]byte(`{
			"Tittle":"SUCCESS",
			"Message": "Forecast Data sent successfully",
			"Type": "success",
			"payloaddata": ` + payload + `
			
			}`))
			c.TplName = "analytics/forecastdashboard.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "ForecastDashboard Page Success")
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

	language := sess.Get("language").(string)
	c.Data["language"] = language
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "language", language)

	auth, err := utils.IsAuthorized(utype, "ForecastDashboard")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	var input_forecasted_product, input_forecasted_product2 string
	input_combined_product := c.Input().Get("input_combined_product")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_combined_product", input_combined_product)

	input_forecasted_product = c.Input().Get("input_forecasted_product")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_forecasted_product", input_forecasted_product)
	input_forecasted_product2 = c.Input().Get("input_forecasted_product2")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_forecasted_product2", input_forecasted_product2)

	input_combined_version := c.Input().Get("input_combined_version")
	input_producer_version := c.Input().Get("input_producer_version")
	input_customer_version := c.Input().Get("input_customer_version")

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_combined_version", input_combined_version)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_producer_version", input_producer_version)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_customer_version", input_customer_version)

	// if input_combined_product != "" {

	// 	input_forecasted_product = input_combined_product
	// 	input_forecasted_product2 = input_combined_product

	// } else {

	// }

	c.Data["ProducerProduct"] = input_forecasted_product
	c.Data["CustomerProduct"] = input_forecasted_product2

	var pdctqty1, pdctqty2 string

	if input_forecasted_product != "" || input_producer_version != "" {

		row99, err := db.Db.Query(`SELECT product_qty FROM ProducerForecast where product_name=? or version=?`, input_forecasted_product, input_producer_version)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get Producer Forecasted ProductQty")
			return
		}
		defer sql.Close(row99)

		_, data99, err := sql.Scan(row99)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get Producer Forecasted ProductQty")
			return
		}

		pdctqty1 = data99[0][0]

	}

	if input_forecasted_product2 != "" || input_customer_version != "" {

		row99, err := db.Db.Query(`SELECT product_qty FROM CustomerForecast where product_name=? or version=?`, input_forecasted_product2, input_customer_version)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get Customer Forecasted ProductQty")
			return
		}
		defer sql.Close(row99)

		_, data99, err := sql.Scan(row99)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get Customer Forecasted ProductQty")
			return
		}

		pdctqty2 = data99[0][0]

	}

	c.Data["pdctqty1"] = pdctqty1
	c.Data["pdctqty2"] = pdctqty2

	currentTime := time.Now()

	// get 12 months name and year from current time

	m1year := strconv.Itoa(currentTime.Year())
	month1 := time.Month(currentTime.Month()).String()
	fmt.Println("Month 1 and Year:", month1, m1year)
	c.Data["month1"] = month1 + " " + m1year

	m2 := currentTime.AddDate(0, 0, 30)
	m2year := strconv.Itoa(m2.Year())
	month2 := time.Month(m2.Month()).String()
	fmt.Println("Month 2 and Year:", month2, m2year)
	c.Data["month2"] = month2 + " " + m2year

	m3 := currentTime.AddDate(0, 0, 60)
	m3year := strconv.Itoa(m3.Year())
	month3 := time.Month(m3.Month()).String()
	fmt.Println("Month 3 and Year:", month3, m3year)
	c.Data["month3"] = month3 + " " + m3year

	m4 := currentTime.AddDate(0, 0, 90)
	m4year := strconv.Itoa(m4.Year())
	month4 := time.Month(m4.Month()).String()
	fmt.Println("Month 4 and Year:", month4, m4year)
	c.Data["month4"] = month4 + " " + m4year

	m5 := currentTime.AddDate(0, 0, 120)
	m5year := strconv.Itoa(m5.Year())
	month5 := time.Month(m5.Month()).String()
	fmt.Println("Month 5 and Year:", month5, m5year)
	c.Data["month5"] = month5 + " " + m5year

	m6 := currentTime.AddDate(0, 0, 150)
	m6year := strconv.Itoa(m6.Year())
	month6 := time.Month(m6.Month()).String()
	fmt.Println("Month 6 and Year:", month6, m6year)
	c.Data["month6"] = month6 + " " + m6year

	m7 := currentTime.AddDate(0, 0, 180)
	m7year := strconv.Itoa(m7.Year())
	month7 := time.Month(m7.Month()).String()
	fmt.Println("Month 7 and Year:", month7, m7year)
	c.Data["month7"] = month7 + " " + m7year

	m8 := currentTime.AddDate(0, 0, 210)
	m8year := strconv.Itoa(m8.Year())
	month8 := time.Month(m8.Month()).String()
	fmt.Println("Month 8 and Year:", month8, m8year)
	c.Data["month8"] = month8 + " " + m8year

	m9 := currentTime.AddDate(0, 0, 240)
	m9year := strconv.Itoa(m9.Year())
	month9 := time.Month(m9.Month()).String()
	fmt.Println("Month 9 and Year:", month9, m9year)
	c.Data["month9"] = month9 + " " + m9year

	m10 := currentTime.AddDate(0, 0, 270)
	m10year := strconv.Itoa(m10.Year())
	month10 := time.Month(m10.Month()).String()
	fmt.Println("Month 10 and Year:", month10, m10year)
	c.Data["month10"] = month10 + " " + m10year

	m11 := currentTime.AddDate(0, 0, 300)
	m11year := strconv.Itoa(m11.Year())
	month11 := time.Month(m11.Month()).String()
	fmt.Println("Month 11 and Year:", month11, m11year)
	c.Data["month11"] = month11 + " " + m11year

	m12 := currentTime.AddDate(0, 0, 330)
	m12year := strconv.Itoa(m12.Year())
	month12 := time.Month(m12.Month()).String()
	fmt.Println("Month 12 and Year:", month12, m12year)
	c.Data["month12"] = month12 + " " + m12year

	var m1end, m2end, m3end, m4end, m5end, m6end, m7end, m8end, m9end, m10end, m11end, m12end string
	var m1startdate, m2startdate, m3startdate, m4startdate, m5startdate, m6startdate, m7startdate, m8startdate, m9startdate, m10startdate, m11startdate, m12startdate string
	var m1enddate, m2enddate, m3enddate, m4enddate, m5enddate, m6enddate, m7enddate, m8enddate, m9enddate, m10enddate, m11enddate, m12enddate string

	// Month1

	if currentTime.Month() == 1 || currentTime.Month() == 3 || currentTime.Month() == 5 || currentTime.Month() == 7 || currentTime.Month() == 8 || currentTime.Month() == 10 || currentTime.Month() == 12 {

		m1end = "31"

	} else {

		m1end = "30"

	}

	if len(strconv.Itoa(int(currentTime.Month()))) == 1 {

		m1startdate = m1year + "-" + "0" + strconv.Itoa(int(currentTime.Month())) + "-" + "01"
		m1enddate = m1year + "-" + "0" + strconv.Itoa(int(currentTime.Month())) + "-" + m1end

	} else {

		m1startdate = m1year + "-" + strconv.Itoa(int(currentTime.Month())) + "-" + "01"
		m1enddate = m1year + "-" + strconv.Itoa(int(currentTime.Month())) + "-" + m1end

	}

	row1, err := db.Db.Query(`SELECT sum(units)  FROM ProducerForecast where start_date>=?  and end_date <= ? and product_name=? and version=?`, m1startdate, m1enddate, input_forecasted_product, input_producer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 1 sum")
		return
	}
	defer sql.Close(row1)

	_, data1, err := sql.Scan(row1)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 1 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "pm1 - ", data1[0][0])
	c.Data["pm1"] = data1[0][0]

	custrow1, err := db.Db.Query(`SELECT sum(units)  FROM CustomerForecast where start_date>=?  and end_date <= ? and product_name=? and version=?`, m1startdate, m1enddate, input_forecasted_product2, input_customer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 1 sum")
		return
	}
	defer sql.Close(custrow1)

	_, custdata1, err := sql.Scan(custrow1)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 1 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cm1 - ", custdata1[0][0])
	c.Data["cm1"] = custdata1[0][0]

	// Month2

	if m2.Month() == 1 || m2.Month() == 3 || m2.Month() == 5 || m2.Month() == 7 || m2.Month() == 8 || m2.Month() == 10 || m2.Month() == 12 {

		m2end = "31"

	} else {

		m2end = "30"

	}

	if len(strconv.Itoa(int(m2.Month()))) == 1 {

		m2startdate = m2year + "-" + "0" + strconv.Itoa(int(m2.Month())) + "-" + "01"
		m2enddate = m2year + "-" + "0" + strconv.Itoa(int(m2.Month())) + "-" + m2end

	} else {

		m2startdate = m2year + "-" + strconv.Itoa(int(m2.Month())) + "-" + "01"
		m2enddate = m2year + "-" + strconv.Itoa(int(m2.Month())) + "-" + m2end

	}

	row2, err := db.Db.Query(`SELECT sum(units)  FROM ProducerForecast where start_date>=?  and end_date <= ? and product_name=? and version=?`, m2startdate, m2enddate, input_forecasted_product, input_producer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 2 sum")
		return
	}
	defer sql.Close(row2)

	_, data2, err := sql.Scan(row2)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 2 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "pm2 - ", data2[0][0])
	c.Data["pm2"] = data2[0][0]

	custrow2, err := db.Db.Query(`SELECT sum(units)  FROM CustomerForecast where start_date>=?  and end_date <= ? and product_name=? and version=?`, m2startdate, m2enddate, input_forecasted_product2, input_customer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 2 sum")
		return
	}
	defer sql.Close(custrow2)

	_, custdata2, err := sql.Scan(custrow2)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 2 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cm2 - ", custdata2[0][0])
	c.Data["cm2"] = custdata2[0][0]

	// Month3

	if m3.Month() == 1 || m3.Month() == 3 || m3.Month() == 5 || m3.Month() == 7 || m3.Month() == 8 || m3.Month() == 10 || m3.Month() == 12 {

		m3end = "31"

	} else {

		m3end = "30"

	}

	if len(strconv.Itoa(int(m3.Month()))) == 1 {

		m3startdate = m3year + "-" + "0" + strconv.Itoa(int(m3.Month())) + "-" + "01"
		m3enddate = m3year + "-" + "0" + strconv.Itoa(int(m3.Month())) + "-" + m3end

	} else {

		m3startdate = m3year + "-" + strconv.Itoa(int(m3.Month())) + "-" + "01"
		m3enddate = m3year + "-" + strconv.Itoa(int(m3.Month())) + "-" + m3end

	}

	row3, err := db.Db.Query(`SELECT sum(units)  FROM ProducerForecast where start_date>=?  and end_date <= ? and product_name=? and version=?`, m3startdate, m3enddate, input_forecasted_product, input_producer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 3 sum")
		return
	}
	defer sql.Close(row3)

	_, data3, err := sql.Scan(row3)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 3 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "pm3 - ", data3[0][0])
	c.Data["pm3"] = data3[0][0]

	custrow3, err := db.Db.Query(`SELECT sum(units)  FROM CustomerForecast where start_date>=?  and end_date <= ? and product_name=? and version=?`, m3startdate, m3enddate, input_forecasted_product2, input_customer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 3 sum")
		return
	}
	defer sql.Close(custrow3)

	_, custdata3, err := sql.Scan(custrow3)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 3 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cm3 - ", custdata3[0][0])
	c.Data["cm3"] = custdata3[0][0]

	// Month4

	if m4.Month() == 1 || m4.Month() == 3 || m4.Month() == 5 || m4.Month() == 7 || m4.Month() == 8 || m4.Month() == 10 || m4.Month() == 12 {

		m4end = "31"

	} else {

		m4end = "30"

	}

	if len(strconv.Itoa(int(m4.Month()))) == 1 {

		m4startdate = m4year + "-" + "0" + strconv.Itoa(int(m4.Month())) + "-" + "01"
		m4enddate = m4year + "-" + "0" + strconv.Itoa(int(m4.Month())) + "-" + m4end

	} else {

		m4startdate = m4year + "-" + strconv.Itoa(int(m4.Month())) + "-" + "01"
		m4enddate = m4year + "-" + strconv.Itoa(int(m4.Month())) + "-" + m4end

	}

	row4, err := db.Db.Query(`SELECT sum(units)  FROM ProducerForecast where start_date>=?  and end_date <= ? and product_name=? and version=?`, m4startdate, m4enddate, input_forecasted_product, input_producer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 4 sum")
		return
	}
	defer sql.Close(row4)

	_, data4, err := sql.Scan(row4)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 4 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "pm4 - ", data4[0][0])
	c.Data["pm4"] = data4[0][0]

	custrow4, err := db.Db.Query(`SELECT sum(units)  FROM CustomerForecast where start_date>=?  and end_date <= ? and product_name=? and version=?`, m4startdate, m4enddate, input_forecasted_product2, input_customer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 4 sum")
		return
	}
	defer sql.Close(custrow4)

	_, custdata4, err := sql.Scan(custrow4)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 4 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cm4 - ", custdata4[0][0])
	c.Data["cm4"] = custdata4[0][0]

	// Month5

	if m5.Month() == 1 || m5.Month() == 3 || m5.Month() == 5 || m5.Month() == 7 || m5.Month() == 8 || m5.Month() == 10 || m5.Month() == 12 {

		m5end = "31"

	} else {

		m5end = "30"

	}

	if len(strconv.Itoa(int(m5.Month()))) == 1 {

		m5startdate = m5year + "-" + "0" + strconv.Itoa(int(m5.Month())) + "-" + "01"
		m5enddate = m5year + "-" + "0" + strconv.Itoa(int(m5.Month())) + "-" + m5end

	} else {

		m5startdate = m5year + "-" + strconv.Itoa(int(m5.Month())) + "-" + "01"
		m5enddate = m5year + "-" + strconv.Itoa(int(m5.Month())) + "-" + m5end

	}

	row5, err := db.Db.Query(`SELECT sum(units)  FROM ProducerForecast where start_date>=?  and end_date <= ? and product_name=? and version=?`, m5startdate, m5enddate, input_forecasted_product, input_producer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 5 sum")
		return
	}
	defer sql.Close(row5)

	_, data5, err := sql.Scan(row5)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 5 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "pm5 - ", data5[0][0])
	c.Data["pm5"] = data5[0][0]

	custrow5, err := db.Db.Query(`SELECT sum(units)  FROM CustomerForecast where start_date>=?  and end_date <= ? and product_name=? and version=?`, m5startdate, m5enddate, input_forecasted_product2, input_customer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 5 sum")
		return
	}
	defer sql.Close(custrow5)

	_, custdata5, err := sql.Scan(custrow5)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 5 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cm5 - ", custdata5[0][0])
	c.Data["cm5"] = custdata5[0][0]

	// Month6

	if m6.Month() == 1 || m6.Month() == 3 || m6.Month() == 5 || m6.Month() == 7 || m6.Month() == 8 || m6.Month() == 10 || m6.Month() == 12 {

		m6end = "31"

	} else {

		m6end = "30"

	}

	if len(strconv.Itoa(int(m6.Month()))) == 1 {

		m6startdate = m6year + "-" + "0" + strconv.Itoa(int(m6.Month())) + "-" + "01"
		m6enddate = m6year + "-" + "0" + strconv.Itoa(int(m6.Month())) + "-" + m6end

	} else {

		m6startdate = m6year + "-" + strconv.Itoa(int(m6.Month())) + "-" + "01"
		m6enddate = m6year + "-" + strconv.Itoa(int(m6.Month())) + "-" + m6end

	}

	row6, err := db.Db.Query(`SELECT sum(units)  FROM ProducerForecast where start_date>=?  and end_date <= ? and product_name=? and version=?`, m6startdate, m6enddate, input_forecasted_product, input_producer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 6 sum")
		return
	}
	defer sql.Close(row6)

	_, data6, err := sql.Scan(row6)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 6 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "pm6 - ", data6[0][0])
	c.Data["pm6"] = data6[0][0]

	custrow6, err := db.Db.Query(`SELECT sum(units)  FROM CustomerForecast where start_date>=?  and end_date <= ? and product_name=? and version=?`, m6startdate, m6enddate, input_forecasted_product2, input_customer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 6 sum")
		return
	}
	defer sql.Close(custrow6)

	_, custdata6, err := sql.Scan(custrow6)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 6 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cm6 - ", custdata6[0][0])
	c.Data["cm6"] = custdata6[0][0]

	// Month7

	if m7.Month() == 1 || m7.Month() == 3 || m7.Month() == 5 || m7.Month() == 7 || m7.Month() == 8 || m7.Month() == 10 || m7.Month() == 12 {

		m7end = "31"

	} else {

		m7end = "30"

	}

	if len(strconv.Itoa(int(m7.Month()))) == 1 {

		m7startdate = m7year + "-" + "0" + strconv.Itoa(int(m7.Month())) + "-" + "01"
		m7enddate = m7year + "-" + "0" + strconv.Itoa(int(m7.Month())) + "-" + m7end

	} else {

		m7startdate = m7year + "-" + strconv.Itoa(int(m7.Month())) + "-" + "01"
		m7enddate = m7year + "-" + strconv.Itoa(int(m7.Month())) + "-" + m7end

	}

	row7, err := db.Db.Query(`SELECT sum(units)  FROM ProducerForecast where start_date>=?  and end_date <= ? and product_name=? and version=?`, m7startdate, m7enddate, input_forecasted_product, input_producer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 7 sum")
		return
	}
	defer sql.Close(row7)

	_, data7, err := sql.Scan(row7)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 7 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "pm7 - ", data7[0][0])
	c.Data["pm7"] = data7[0][0]

	custrow7, err := db.Db.Query(`SELECT sum(units)  FROM CustomerForecast where start_date>=?  and end_date <= ? and product_name=? and version=?`, m7startdate, m7enddate, input_forecasted_product2, input_customer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 7 sum")
		return
	}
	defer sql.Close(custrow7)

	_, custdata7, err := sql.Scan(custrow7)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 7 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cm7 - ", custdata7[0][0])
	c.Data["cm7"] = custdata7[0][0]

	// Month8

	if m8.Month() == 1 || m8.Month() == 3 || m8.Month() == 5 || m8.Month() == 7 || m8.Month() == 8 || m8.Month() == 10 || m8.Month() == 12 {

		m8end = "31"

	} else {

		m8end = "30"

	}

	if len(strconv.Itoa(int(m8.Month()))) == 1 {

		m8startdate = m8year + "-" + "0" + strconv.Itoa(int(m8.Month())) + "-" + "01"
		m8enddate = m8year + "-" + "0" + strconv.Itoa(int(m8.Month())) + "-" + m8end

	} else {

		m8startdate = m8year + "-" + strconv.Itoa(int(m8.Month())) + "-" + "01"
		m8enddate = m8year + "-" + strconv.Itoa(int(m8.Month())) + "-" + m8end

	}

	row8, err := db.Db.Query(`SELECT sum(units)  FROM ProducerForecast where start_date>=?  and end_date <= ? and product_name=? and version=?`, m8startdate, m8enddate, input_forecasted_product, input_producer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 8 sum")
		return
	}
	defer sql.Close(row8)

	_, data8, err := sql.Scan(row8)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 8 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "pm8 - ", data8[0][0])
	c.Data["pm8"] = data8[0][0]

	custrow8, err := db.Db.Query(`SELECT sum(units)  FROM CustomerForecast where start_date>=?  and end_date <= ? and product_name=? and version=?`, m8startdate, m8enddate, input_forecasted_product2, input_customer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 8 sum")
		return
	}
	defer sql.Close(custrow8)

	_, custdata8, err := sql.Scan(custrow8)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 8 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cm8 - ", custdata8[0][0])
	c.Data["cm8"] = custdata8[0][0]

	// Month9

	if m9.Month() == 1 || m9.Month() == 3 || m9.Month() == 5 || m9.Month() == 7 || m9.Month() == 8 || m9.Month() == 10 || m9.Month() == 12 {

		m9end = "31"

	} else {

		m9end = "30"

	}

	if len(strconv.Itoa(int(m9.Month()))) == 1 {

		m9startdate = m9year + "-" + "0" + strconv.Itoa(int(m9.Month())) + "-" + "01"
		m9enddate = m9year + "-" + "0" + strconv.Itoa(int(m9.Month())) + "-" + m9end

	} else {

		m9startdate = m9year + "-" + strconv.Itoa(int(m9.Month())) + "-" + "01"
		m9enddate = m9year + "-" + strconv.Itoa(int(m9.Month())) + "-" + m9end

	}

	row9, err := db.Db.Query(`SELECT sum(units)  FROM ProducerForecast where start_date>=?  and end_date <= ? and product_name=? and version=?`, m9startdate, m9enddate, input_forecasted_product, input_producer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 9 sum")
		return
	}
	defer sql.Close(row9)

	_, data9, err := sql.Scan(row9)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 9 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "pm9 - ", data9[0][0])
	c.Data["pm9"] = data9[0][0]

	custrow9, err := db.Db.Query(`SELECT sum(units)  FROM CustomerForecast where start_date>=?  and end_date <= ? and product_name=? and version=?`, m9startdate, m9enddate, input_forecasted_product2, input_customer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 9 sum")
		return
	}
	defer sql.Close(custrow9)

	_, custdata9, err := sql.Scan(custrow9)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 9 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cm9 - ", custdata9[0][0])
	c.Data["cm9"] = custdata9[0][0]

	// Month10

	if m10.Month() == 1 || m10.Month() == 3 || m10.Month() == 5 || m10.Month() == 7 || m10.Month() == 8 || m10.Month() == 10 || m10.Month() == 12 {

		m10end = "31"

	} else {

		m10end = "30"

	}

	if len(strconv.Itoa(int(m10.Month()))) == 1 {

		m10startdate = m10year + "-" + "0" + strconv.Itoa(int(m10.Month())) + "-" + "01"
		m10enddate = m10year + "-" + "0" + strconv.Itoa(int(m10.Month())) + "-" + m10end

	} else {

		m10startdate = m10year + "-" + strconv.Itoa(int(m10.Month())) + "-" + "01"
		m10enddate = m10year + "-" + strconv.Itoa(int(m10.Month())) + "-" + m10end

	}

	row10, err := db.Db.Query(`SELECT sum(units)  FROM ProducerForecast where start_date>=?  and end_date <= ? and product_name=? and version=?`, m10startdate, m10enddate, input_forecasted_product, input_producer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 10 sum")
		return
	}
	defer sql.Close(row10)

	_, data10, err := sql.Scan(row10)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 10 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "pm10 - ", data10[0][0])
	c.Data["pm10"] = data10[0][0]

	custrow10, err := db.Db.Query(`SELECT sum(units)  FROM CustomerForecast where start_date>=?  and end_date <= ? and product_name=? and version=?`, m10startdate, m10enddate, input_forecasted_product2, input_customer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 10 sum")
		return
	}
	defer sql.Close(custrow10)

	_, custdata10, err := sql.Scan(custrow10)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 10 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cm10 - ", custdata10[0][0])
	c.Data["cm10"] = custdata10[0][0]

	// Month11

	if m11.Month() == 1 || m11.Month() == 3 || m11.Month() == 5 || m11.Month() == 7 || m11.Month() == 8 || m11.Month() == 10 || m11.Month() == 12 {

		m11end = "31"

	} else {

		m11end = "30"

	}

	if len(strconv.Itoa(int(m11.Month()))) == 1 {

		m11startdate = m11year + "-" + "0" + strconv.Itoa(int(m11.Month())) + "-" + "01"
		m11enddate = m11year + "-" + "0" + strconv.Itoa(int(m11.Month())) + "-" + m11end

	} else {

		m11startdate = m11year + "-" + strconv.Itoa(int(m11.Month())) + "-" + "01"
		m11enddate = m11year + "-" + strconv.Itoa(int(m11.Month())) + "-" + m11end

	}

	row11, err := db.Db.Query(`SELECT sum(units)  FROM ProducerForecast where start_date>=?  and end_date <= ? and product_name=? and version=?`, m11startdate, m11enddate, input_forecasted_product, input_producer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 11 sum")
		return
	}
	defer sql.Close(row11)

	_, data11, err := sql.Scan(row11)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 11 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "pm11 - ", data11[0][0])
	c.Data["pm11"] = data11[0][0]

	custrow11, err := db.Db.Query(`SELECT sum(units)  FROM CustomerForecast where start_date>=?  and end_date <= ? and product_name=? and version=?`, m11startdate, m11enddate, input_forecasted_product2, input_customer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 11 sum")
		return
	}
	defer sql.Close(custrow11)

	_, custdata11, err := sql.Scan(custrow11)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 11 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cm11 - ", custdata11[0][0])
	c.Data["cm11"] = custdata11[0][0]

	// Month12

	if m12.Month() == 1 || m12.Month() == 3 || m12.Month() == 5 || m12.Month() == 7 || m12.Month() == 8 || m12.Month() == 10 || m12.Month() == 12 {

		m12end = "31"

	} else {

		m12end = "30"

	}

	if len(strconv.Itoa(int(m12.Month()))) == 1 {

		m12startdate = m12year + "-" + "0" + strconv.Itoa(int(m12.Month())) + "-" + "01"
		m12enddate = m12year + "-" + "0" + strconv.Itoa(int(m12.Month())) + "-" + m12end

	} else {

		m12startdate = m12year + "-" + strconv.Itoa(int(m12.Month())) + "-" + "01"
		m12enddate = m12year + "-" + strconv.Itoa(int(m12.Month())) + "-" + m12end

	}

	row12, err := db.Db.Query(`SELECT sum(units)  FROM ProducerForecast where start_date>=?  and end_date <= ? and product_name=? and version=?`, m12startdate, m12enddate, input_forecasted_product, input_producer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 12 sum")
		return
	}
	defer sql.Close(row12)

	_, data12, err := sql.Scan(row12)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 12 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "pm12 - ", data12[0][0])
	c.Data["pm12"] = data12[0][0]

	custrow12, err := db.Db.Query(`SELECT sum(units)  FROM CustomerForecast where start_date>=?  and end_date <= ? and product_name=? and version=?`, m12startdate, m12enddate, input_forecasted_product2, input_customer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get  CustomerProduct Month 12 sum")
		return
	}
	defer sql.Close(custrow12)

	_, custdata12, err := sql.Scan(custrow12)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get  CustomerProduct Month 12 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cm12 - ", custdata12[0][0])
	c.Data["cm12"] = custdata12[0][0]

	forecastdatapayload := ForecastData{

		Month1:      month1 + " " + m1year,
		Month2:      month2 + " " + m2year,
		Month3:      month3 + " " + m3year,
		Month4:      month4 + " " + m4year,
		Month5:      month5 + " " + m5year,
		Month6:      month6 + " " + m6year,
		Month7:      month7 + " " + m7year,
		Month8:      month8 + " " + m8year,
		Month9:      month9 + " " + m9year,
		Month10:     month10 + " " + m10year,
		Month11:     month11 + " " + m11year,
		Month12:     month12 + " " + m12year,
		PM1:         data1[0][0],
		PM2:         data2[0][0],
		PM3:         data3[0][0],
		PM4:         data4[0][0],
		PM5:         data5[0][0],
		PM6:         data6[0][0],
		PM7:         data7[0][0],
		PM8:         data8[0][0],
		PM9:         data9[0][0],
		PM10:        data10[0][0],
		PM11:        data11[0][0],
		PM12:        data12[0][0],
		CM1:         custdata1[0][0],
		CM2:         custdata2[0][0],
		CM3:         custdata3[0][0],
		CM4:         custdata4[0][0],
		CM5:         custdata5[0][0],
		CM6:         custdata6[0][0],
		CM7:         custdata7[0][0],
		CM8:         custdata8[0][0],
		CM9:         custdata9[0][0],
		CM10:        custdata10[0][0],
		CM11:        custdata11[0][0],
		CM12:        custdata12[0][0],
		Product1:    input_forecasted_product,
		Product2:    input_forecasted_product2,
		Product1Qty: pdctqty1,
		Product2QTY: pdctqty2,
	}

	var jsonData2 []byte
	jsonData2, err = json.Marshal(forecastdatapayload)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	}

	payload = f(string(jsonData2))
	fmt.Println("forecast payload in return", payload)

	return
}

func (c *CombinedForecastDashboard) Get() {

	log.Println(beego.AppConfig.String("loglevel"), "Info", "ForecastDashboard Page Start")
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
			log.Println(beego.AppConfig.String("loglevel"), "Info", "ForecastDashboard Page Fail")
		} else if err != nil {
			c.Data["DisplayMessage"] = err.Error()
			c.TplName = "error/error.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "ForecastDashboard Page Fail")
		} else {

			c.TplName = "analytics/forecastdashboard.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "ForecastDashboard Page Success")
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

	language := sess.Get("language").(string)
	c.Data["language"] = language
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "language", language)

	auth, err := utils.IsAuthorized(utype, "Dashboard")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")
}

func (c *CombinedForecastDashboard) Post() {

	var payload string

	log.Println(beego.AppConfig.String("loglevel"), "Info", "ForecastDashboard Page Start")
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
			log.Println(beego.AppConfig.String("loglevel"), "Info", "ForecastDashboard Page Fail")
		} else if err != nil {
			c.Data["DisplayMessage"] = err.Error()
			c.TplName = "error/error.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "ForecastDashboard Page Fail")
		} else {

			c.TplName = "analytics/forecastdashboard.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "ForecastDashboard Page Success")
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

	language := sess.Get("language").(string)
	c.Data["language"] = language
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "language", language)

	auth, err := utils.IsAuthorized(utype, "Dashboard")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	//-----------------Start Customer Forecast Data------------------------

	custrow, err := db.Db.Query(`SELECT distinct(product_name),product_qty FROM CustomerForecast`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Customer Forecasted Products")
		return
	}
	defer sql.Close(custrow)

	_, custdata, err := sql.Scan(custrow)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Customer Forecasted Products")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", custdata, "\nData len - ", len(custdata))

	var Dis1 Display
	for i := range custdata {
		var d Field
		d.ProducerForecastProduct = custdata[i][0]
		d.ProductQty = custdata[i][1]

		Dis1.Fields = append(Dis1.Fields, d)
	}
	c.Data["Dis1"] = Dis1
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Customer Forecasted Products - ", Dis1)

	row21, err := db.Db.Query(`SELECT version, MAX(create_date) AS create_date
FROM CustomerForecast
GROUP BY version`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Producer Forecasted Products")
		return
	}
	defer sql.Close(row21)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row21)
	_, data21, err := sql.Scan(row21)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Producer Forecasted Products")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data21, "\nData len - ", len(data21))

	var Dis11 Display
	for i := range data21 {
		var d Field
		d.Version = data21[i][0]
		d.CreateDate = data21[i][1]

		Dis11.Fields = append(Dis11.Fields, d)
	}
	c.Data["Dis11"] = Dis11
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Customer Forecasted Version - ", Dis11)

	// c.Data["CVersion"] = data21[0][0]

	if len(data21) > 0 && len(data21[0]) > 0 {
		c.Data["CVersion"] = data21[0][0]
	} else {
		log.Println("Warning: data21 is empty, cannot set CVersion")
		c.Data["CVersion"] = "" // or a default value
	}

	//-----------------End Customer Forecast Data------------------------

	//-----------------Start Producer Forecast Data------------------------

	row, err := db.Db.Query(`SELECT distinct(product_name),product_qty FROM ProducerForecast`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Producer Forecasted Products")
		return
	}
	defer sql.Close(row)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Producer Forecasted Products")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))

	var Dis Display
	for i := range data {
		var d Field
		d.ProducerForecastProduct = data[i][0]
		d.ProductQty = data[i][1]

		Dis.Fields = append(Dis.Fields, d)
	}
	c.Data["Dis"] = Dis
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Producer Forecasted Products - ", Dis)

	row22, err := db.Db.Query(`SELECT version, MAX(create_date) AS create_date
FROM ProducerForecast
GROUP BY version`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Producer Forecasted Products")
		return
	}
	defer sql.Close(row22)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row22)
	_, data22, err := sql.Scan(row22)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Producer Forecasted Products")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data22, "\nData len - ", len(data22))

	var Dis12 Display
	for i := range data22 {
		var d Field
		d.Version = data22[i][0]
		d.CreateDate = data22[i][1]

		Dis12.Fields = append(Dis12.Fields, d)
	}
	c.Data["Dis12"] = Dis12
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Producer Forecasted Version - ", Dis12)

	// c.Data["PVersion"] = data22[0][0]

	if len(data22) > 0 && len(data22[0]) > 0 {
		c.Data["PVersion"] = data22[0][0]
	} else {
		log.Println("Warning: data22 is empty, cannot set PVersion")
		c.Data["PVersion"] = "" // or a default value
	}
	//-----------------End Customer Forecast Data------------------------

	//-----------------Start Combined Producer and Customer Forecast Data------------------------

	row222, err := db.Db.Query(`(SELECT distinct(product_name),product_qty FROM ProducerForecast)UNION(SELECT distinct(product_name),product_qty FROM CustomerForecast)`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Producer Forecasted Products")
		return
	}
	defer sql.Close(row222)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row)
	_, data222, err := sql.Scan(row222)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Producer Forecasted Products")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))

	var Dis2 Display
	for i := range data222 {
		var d Field
		d.ProducerForecastProduct = data222[i][0]
		d.ProductQty = data222[i][1]

		Dis2.Fields = append(Dis2.Fields, d)
	}
	c.Data["Dis3"] = Dis2
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Combined Forecasted Products - ", Dis2)

	row23, err := db.Db.Query(`SELECT version, MIN(create_date) AS create_date
FROM (
    SELECT DISTINCT version, create_date FROM ProducerForecast
    UNION
    SELECT DISTINCT version, create_date FROM CustomerForecast
) AS combined
GROUP BY version
`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Producer Forecasted Products")
		return
	}
	defer sql.Close(row23)
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Row Data - ", row23)
	_, data23, err := sql.Scan(row23)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Producer Forecasted Products")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data23, "\nData len - ", len(data23))

	var Dis13 Display
	for i := range data23 {
		var d Field
		d.Version = data23[i][0]
		d.CreateDate = data23[i][1]

		Dis13.Fields = append(Dis13.Fields, d)
	}
	c.Data["Dis13"] = Dis13
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Combined Forecasted Products - ", Dis13)

	// c.Data["PCVersion"] = data23[0][0]

	if len(data23) > 0 && len(data23[0]) > 0 {
		c.Data["PCVersion"] = data23[0][0]
	} else {
		log.Println("Warning: data23 is empty, cannot set PCVersion")
		c.Data["PCVersion"] = "" // or a default value
	}

	//-----------------End Combined Producer and Customer Forecast Data------------------------

	var input_forecasted_product, input_forecasted_product2, input_producer_version, input_customer_version string
	input_combined_product := c.Input().Get("input_combined_product")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_combined_product", input_combined_product)

	input_forecasted_product = input_combined_product
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_forecasted_product", input_forecasted_product)
	input_forecasted_product2 = input_combined_product
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_forecasted_product2", input_forecasted_product2)

	input_combined_version := c.Input().Get("input_combined_version")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_combined_version", input_combined_version)

	input_producer_version = input_combined_version
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_producer_version", input_producer_version)
	input_customer_version = input_combined_version
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_customer_version", input_customer_version)

	row99, err := db.Db.Query(`SELECT count(*) FROM ProducerForecast where product_name=? and version =?`, input_forecasted_product, input_producer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Producer Forecasted ProductQty Count")
		return
	}
	defer sql.Close(row99)

	_, data99, err := sql.Scan(row99)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Producer Forecasted ProductQty Count")
		return
	}

	row98, err := db.Db.Query(`SELECT count(*) FROM CustomerForecast where product_name=? and version =?`, input_forecasted_product2, input_customer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Customer Forecasted ProductQty Count")
		return
	}
	defer sql.Close(row98)

	_, data98, err := sql.Scan(row98)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get Customer Forecasted ProductQty Count")
		return
	}

	var pdctqty1 string

	if data99[0][0] != "0" {

		row999, err := db.Db.Query(`SELECT product_qty FROM ProducerForecast where product_name=? and version =?`, input_forecasted_product, input_producer_version)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get Producer Forecasted ProductQty ")
			return
		}
		defer sql.Close(row999)

		_, data999, err := sql.Scan(row999)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get Producer Forecasted ProductQty ")
			return
		}

		pdctqty1 = data999[0][0]

	}

	if data98[0][0] != "0" {

		row988, err := db.Db.Query(`SELECT product_qty FROM CustomerForecast where product_name=? and version =?`, input_forecasted_product2, input_customer_version)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get Customer Forecasted ProductQty")
			return
		}
		defer sql.Close(row988)

		_, data988, err := sql.Scan(row988)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Unable to get Customer Forecasted ProductQty")
			return
		}

		pdctqty1 = data988[0][0]

	}

	c.Data["product"] = input_combined_product
	c.Data["productqty"] = pdctqty1

	currentTime := time.Now()

	// get 12 months name and year from current time

	m1year := strconv.Itoa(currentTime.Year())
	month1 := time.Month(currentTime.Month()).String()
	fmt.Println("Month 1 and Year:", month1, m1year)
	c.Data["month1"] = month1 + " " + m1year

	m2 := currentTime.AddDate(0, 0, 30)
	m2year := strconv.Itoa(m2.Year())
	month2 := time.Month(m2.Month()).String()
	fmt.Println("Month 2 and Year:", month2, m2year)
	c.Data["month2"] = month2 + " " + m2year

	m3 := currentTime.AddDate(0, 0, 60)
	m3year := strconv.Itoa(m3.Year())
	month3 := time.Month(m3.Month()).String()
	fmt.Println("Month 3 and Year:", month3, m3year)
	c.Data["month3"] = month3 + " " + m3year

	m4 := currentTime.AddDate(0, 0, 90)
	m4year := strconv.Itoa(m4.Year())
	month4 := time.Month(m4.Month()).String()
	fmt.Println("Month 4 and Year:", month4, m4year)
	c.Data["month4"] = month4 + " " + m4year

	m5 := currentTime.AddDate(0, 0, 120)
	m5year := strconv.Itoa(m5.Year())
	month5 := time.Month(m5.Month()).String()
	fmt.Println("Month 5 and Year:", month5, m5year)
	c.Data["month5"] = month5 + " " + m5year

	m6 := currentTime.AddDate(0, 0, 150)
	m6year := strconv.Itoa(m6.Year())
	month6 := time.Month(m6.Month()).String()
	fmt.Println("Month 6 and Year:", month6, m6year)
	c.Data["month6"] = month6 + " " + m6year

	m7 := currentTime.AddDate(0, 0, 180)
	m7year := strconv.Itoa(m7.Year())
	month7 := time.Month(m7.Month()).String()
	fmt.Println("Month 7 and Year:", month7, m7year)
	c.Data["month7"] = month7 + " " + m7year

	m8 := currentTime.AddDate(0, 0, 210)
	m8year := strconv.Itoa(m8.Year())
	month8 := time.Month(m8.Month()).String()
	fmt.Println("Month 8 and Year:", month8, m8year)
	c.Data["month8"] = month8 + " " + m8year

	m9 := currentTime.AddDate(0, 0, 240)
	m9year := strconv.Itoa(m9.Year())
	month9 := time.Month(m9.Month()).String()
	fmt.Println("Month 9 and Year:", month9, m9year)
	c.Data["month9"] = month9 + " " + m9year

	m10 := currentTime.AddDate(0, 0, 270)
	m10year := strconv.Itoa(m10.Year())
	month10 := time.Month(m10.Month()).String()
	fmt.Println("Month 10 and Year:", month10, m10year)
	c.Data["month10"] = month10 + " " + m10year

	m11 := currentTime.AddDate(0, 0, 300)
	m11year := strconv.Itoa(m11.Year())
	month11 := time.Month(m11.Month()).String()
	fmt.Println("Month 11 and Year:", month11, m11year)
	c.Data["month11"] = month11 + " " + m11year

	m12 := currentTime.AddDate(0, 0, 330)
	m12year := strconv.Itoa(m12.Year())
	month12 := time.Month(m12.Month()).String()
	fmt.Println("Month 12 and Year:", month12, m12year)
	c.Data["month12"] = month12 + " " + m12year

	var m1end, m2end, m3end, m4end, m5end, m6end, m7end, m8end, m9end, m10end, m11end, m12end string
	var m1startdate, m2startdate, m3startdate, m4startdate, m5startdate, m6startdate, m7startdate, m8startdate, m9startdate, m10startdate, m11startdate, m12startdate string
	var m1enddate, m2enddate, m3enddate, m4enddate, m5enddate, m6enddate, m7enddate, m8enddate, m9enddate, m10enddate, m11enddate, m12enddate string

	// Month1

	if currentTime.Month() == 1 || currentTime.Month() == 3 || currentTime.Month() == 5 || currentTime.Month() == 7 || currentTime.Month() == 8 || currentTime.Month() == 10 || currentTime.Month() == 12 {

		m1end = "31"

	} else {

		m1end = "30"

	}

	if len(strconv.Itoa(int(currentTime.Month()))) == 1 {

		m1startdate = m1year + "-" + "0" + strconv.Itoa(int(currentTime.Month())) + "-" + "01"
		m1enddate = m1year + "-" + "0" + strconv.Itoa(int(currentTime.Month())) + "-" + m1end

	} else {

		m1startdate = m1year + "-" + strconv.Itoa(int(currentTime.Month())) + "-" + "01"
		m1enddate = m1year + "-" + strconv.Itoa(int(currentTime.Month())) + "-" + m1end

	}

	row1, err := db.Db.Query(`SELECT sum(units)  FROM ProducerForecast where start_date>=?  and end_date <= ? and product_name=? and version =?`, m1startdate, m1enddate, input_forecasted_product, input_producer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 1 sum")
		return
	}
	defer sql.Close(row1)

	_, data1, err := sql.Scan(row1)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 1 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "pm1 - ", data1[0][0])
	c.Data["pm1"] = data1[0][0]

	custrow1, err := db.Db.Query(`SELECT sum(units)  FROM CustomerForecast where start_date>=?  and end_date <= ? and product_name=? and version =?`, m1startdate, m1enddate, input_forecasted_product2, input_customer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 1 sum")
		return
	}
	defer sql.Close(custrow1)

	_, custdata1, err := sql.Scan(custrow1)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 1 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cm1 - ", custdata1[0][0])
	c.Data["cm1"] = custdata1[0][0]

	// Month2

	if m2.Month() == 1 || m2.Month() == 3 || m2.Month() == 5 || m2.Month() == 7 || m2.Month() == 8 || m2.Month() == 10 || m2.Month() == 12 {

		m2end = "31"

	} else {

		m2end = "30"

	}

	if len(strconv.Itoa(int(m2.Month()))) == 1 {

		m2startdate = m2year + "-" + "0" + strconv.Itoa(int(m2.Month())) + "-" + "01"
		m2enddate = m2year + "-" + "0" + strconv.Itoa(int(m2.Month())) + "-" + m2end

	} else {

		m2startdate = m2year + "-" + strconv.Itoa(int(m2.Month())) + "-" + "01"
		m2enddate = m2year + "-" + strconv.Itoa(int(m2.Month())) + "-" + m2end

	}

	row2, err := db.Db.Query(`SELECT sum(units)  FROM ProducerForecast where start_date>=?  and end_date <= ? and product_name=? and version =?`, m2startdate, m2enddate, input_forecasted_product, input_producer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 2 sum")
		return
	}
	defer sql.Close(row2)

	_, data2, err := sql.Scan(row2)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 2 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "pm2 - ", data2[0][0])
	c.Data["pm2"] = data2[0][0]

	custrow2, err := db.Db.Query(`SELECT sum(units)  FROM CustomerForecast where start_date>=?  and end_date <= ? and product_name=? and version =?`, m2startdate, m2enddate, input_forecasted_product2, input_customer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 2 sum")
		return
	}
	defer sql.Close(custrow2)

	_, custdata2, err := sql.Scan(custrow2)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 2 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cm2 - ", custdata2[0][0])
	c.Data["cm2"] = custdata2[0][0]

	// Month3

	if m3.Month() == 1 || m3.Month() == 3 || m3.Month() == 5 || m3.Month() == 7 || m3.Month() == 8 || m3.Month() == 10 || m3.Month() == 12 {

		m3end = "31"

	} else {

		m3end = "30"

	}

	if len(strconv.Itoa(int(m3.Month()))) == 1 {

		m3startdate = m3year + "-" + "0" + strconv.Itoa(int(m3.Month())) + "-" + "01"
		m3enddate = m3year + "-" + "0" + strconv.Itoa(int(m3.Month())) + "-" + m3end

	} else {

		m3startdate = m3year + "-" + strconv.Itoa(int(m3.Month())) + "-" + "01"
		m3enddate = m3year + "-" + strconv.Itoa(int(m3.Month())) + "-" + m3end

	}

	row3, err := db.Db.Query(`SELECT sum(units)  FROM ProducerForecast where start_date>=?  and end_date <= ? and product_name=? and version =?`, m3startdate, m3enddate, input_forecasted_product, input_producer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 3 sum")
		return
	}
	defer sql.Close(row3)

	_, data3, err := sql.Scan(row3)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 3 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "pm3 - ", data3[0][0])
	c.Data["pm3"] = data3[0][0]

	custrow3, err := db.Db.Query(`SELECT sum(units)  FROM CustomerForecast where start_date>=?  and end_date <= ? and product_name=? and version =?`, m3startdate, m3enddate, input_forecasted_product2, input_customer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 3 sum")
		return
	}
	defer sql.Close(custrow3)

	_, custdata3, err := sql.Scan(custrow3)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 3 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cm3 - ", custdata3[0][0])
	c.Data["cm3"] = custdata3[0][0]

	// Month4

	if m4.Month() == 1 || m4.Month() == 3 || m4.Month() == 5 || m4.Month() == 7 || m4.Month() == 8 || m4.Month() == 10 || m4.Month() == 12 {

		m4end = "31"

	} else {

		m4end = "30"

	}

	if len(strconv.Itoa(int(m4.Month()))) == 1 {

		m4startdate = m4year + "-" + "0" + strconv.Itoa(int(m4.Month())) + "-" + "01"
		m4enddate = m4year + "-" + "0" + strconv.Itoa(int(m4.Month())) + "-" + m4end

	} else {

		m4startdate = m4year + "-" + strconv.Itoa(int(m4.Month())) + "-" + "01"
		m4enddate = m4year + "-" + strconv.Itoa(int(m4.Month())) + "-" + m4end

	}

	row4, err := db.Db.Query(`SELECT sum(units)  FROM ProducerForecast where start_date>=?  and end_date <= ? and product_name=? and version =?`, m4startdate, m4enddate, input_forecasted_product, input_producer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 4 sum")
		return
	}
	defer sql.Close(row4)

	_, data4, err := sql.Scan(row4)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 4 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "pm4 - ", data4[0][0])
	c.Data["pm4"] = data4[0][0]

	custrow4, err := db.Db.Query(`SELECT sum(units)  FROM CustomerForecast where start_date>=?  and end_date <= ? and product_name=? and version =?`, m4startdate, m4enddate, input_forecasted_product2, input_customer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 4 sum")
		return
	}
	defer sql.Close(custrow4)

	_, custdata4, err := sql.Scan(custrow4)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 4 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cm4 - ", custdata4[0][0])
	c.Data["cm4"] = custdata4[0][0]

	// Month5

	if m5.Month() == 1 || m5.Month() == 3 || m5.Month() == 5 || m5.Month() == 7 || m5.Month() == 8 || m5.Month() == 10 || m5.Month() == 12 {

		m5end = "31"

	} else {

		m5end = "30"

	}

	if len(strconv.Itoa(int(m5.Month()))) == 1 {

		m5startdate = m5year + "-" + "0" + strconv.Itoa(int(m5.Month())) + "-" + "01"
		m5enddate = m5year + "-" + "0" + strconv.Itoa(int(m5.Month())) + "-" + m5end

	} else {

		m5startdate = m5year + "-" + strconv.Itoa(int(m5.Month())) + "-" + "01"
		m5enddate = m5year + "-" + strconv.Itoa(int(m5.Month())) + "-" + m5end

	}

	row5, err := db.Db.Query(`SELECT sum(units)  FROM ProducerForecast where start_date>=?  and end_date <= ? and product_name=? and version =?`, m5startdate, m5enddate, input_forecasted_product, input_producer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 5 sum")
		return
	}
	defer sql.Close(row5)

	_, data5, err := sql.Scan(row5)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 5 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "pm5 - ", data5[0][0])
	c.Data["pm5"] = data5[0][0]

	custrow5, err := db.Db.Query(`SELECT sum(units)  FROM CustomerForecast where start_date>=?  and end_date <= ? and product_name=? and version =?`, m5startdate, m5enddate, input_forecasted_product2, input_customer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 5 sum")
		return
	}
	defer sql.Close(custrow5)

	_, custdata5, err := sql.Scan(custrow5)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 5 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cm5 - ", custdata5[0][0])
	c.Data["cm5"] = custdata5[0][0]

	// Month6

	if m6.Month() == 1 || m6.Month() == 3 || m6.Month() == 5 || m6.Month() == 7 || m6.Month() == 8 || m6.Month() == 10 || m6.Month() == 12 {

		m6end = "31"

	} else {

		m6end = "30"

	}

	if len(strconv.Itoa(int(m6.Month()))) == 1 {

		m6startdate = m6year + "-" + "0" + strconv.Itoa(int(m6.Month())) + "-" + "01"
		m6enddate = m6year + "-" + "0" + strconv.Itoa(int(m6.Month())) + "-" + m6end

	} else {

		m6startdate = m6year + "-" + strconv.Itoa(int(m6.Month())) + "-" + "01"
		m6enddate = m6year + "-" + strconv.Itoa(int(m6.Month())) + "-" + m6end

	}

	row6, err := db.Db.Query(`SELECT sum(units)  FROM ProducerForecast where start_date>=?  and end_date <= ? and product_name=? and version =?`, m6startdate, m6enddate, input_forecasted_product, input_producer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 6 sum")
		return
	}
	defer sql.Close(row6)

	_, data6, err := sql.Scan(row6)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 6 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "pm6 - ", data6[0][0])
	c.Data["pm6"] = data6[0][0]

	custrow6, err := db.Db.Query(`SELECT sum(units)  FROM CustomerForecast where start_date>=?  and end_date <= ? and product_name=? and version =?`, m6startdate, m6enddate, input_forecasted_product2, input_customer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 6 sum")
		return
	}
	defer sql.Close(custrow6)

	_, custdata6, err := sql.Scan(custrow6)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 6 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cm6 - ", custdata6[0][0])
	c.Data["cm6"] = custdata6[0][0]

	// Month7

	if m7.Month() == 1 || m7.Month() == 3 || m7.Month() == 5 || m7.Month() == 7 || m7.Month() == 8 || m7.Month() == 10 || m7.Month() == 12 {

		m7end = "31"

	} else {

		m7end = "30"

	}

	if len(strconv.Itoa(int(m7.Month()))) == 1 {

		m7startdate = m7year + "-" + "0" + strconv.Itoa(int(m7.Month())) + "-" + "01"
		m7enddate = m7year + "-" + "0" + strconv.Itoa(int(m7.Month())) + "-" + m7end

	} else {

		m7startdate = m7year + "-" + strconv.Itoa(int(m7.Month())) + "-" + "01"
		m7enddate = m7year + "-" + strconv.Itoa(int(m7.Month())) + "-" + m7end

	}

	row7, err := db.Db.Query(`SELECT sum(units)  FROM ProducerForecast where start_date>=?  and end_date <= ? and product_name=? and version =?`, m7startdate, m7enddate, input_forecasted_product, input_producer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 7 sum")
		return
	}
	defer sql.Close(row7)

	_, data7, err := sql.Scan(row7)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 7 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "pm7 - ", data7[0][0])
	c.Data["pm7"] = data7[0][0]

	custrow7, err := db.Db.Query(`SELECT sum(units)  FROM CustomerForecast where start_date>=?  and end_date <= ? and product_name=? and version =?`, m7startdate, m7enddate, input_forecasted_product2, input_customer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 7 sum")
		return
	}
	defer sql.Close(custrow7)

	_, custdata7, err := sql.Scan(custrow7)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 7 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cm7 - ", custdata7[0][0])
	c.Data["cm7"] = custdata7[0][0]

	// Month8

	if m8.Month() == 1 || m8.Month() == 3 || m8.Month() == 5 || m8.Month() == 7 || m8.Month() == 8 || m8.Month() == 10 || m8.Month() == 12 {

		m8end = "31"

	} else {

		m8end = "30"

	}

	if len(strconv.Itoa(int(m8.Month()))) == 1 {

		m8startdate = m8year + "-" + "0" + strconv.Itoa(int(m8.Month())) + "-" + "01"
		m8enddate = m8year + "-" + "0" + strconv.Itoa(int(m8.Month())) + "-" + m8end

	} else {

		m8startdate = m8year + "-" + strconv.Itoa(int(m8.Month())) + "-" + "01"
		m8enddate = m8year + "-" + strconv.Itoa(int(m8.Month())) + "-" + m8end

	}

	row8, err := db.Db.Query(`SELECT sum(units)  FROM ProducerForecast where start_date>=?  and end_date <= ? and product_name=? and version =?`, m8startdate, m8enddate, input_forecasted_product, input_producer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 8 sum")
		return
	}
	defer sql.Close(row8)

	_, data8, err := sql.Scan(row8)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 8 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "pm8 - ", data8[0][0])
	c.Data["pm8"] = data8[0][0]

	custrow8, err := db.Db.Query(`SELECT sum(units)  FROM CustomerForecast where start_date>=?  and end_date <= ? and product_name=? and version =?`, m8startdate, m8enddate, input_forecasted_product2, input_customer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 8 sum")
		return
	}
	defer sql.Close(custrow8)

	_, custdata8, err := sql.Scan(custrow8)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 8 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cm8 - ", custdata8[0][0])
	c.Data["cm8"] = custdata8[0][0]

	// Month9

	if m9.Month() == 1 || m9.Month() == 3 || m9.Month() == 5 || m9.Month() == 7 || m9.Month() == 8 || m9.Month() == 10 || m9.Month() == 12 {

		m9end = "31"

	} else {

		m9end = "30"

	}

	if len(strconv.Itoa(int(m9.Month()))) == 1 {

		m9startdate = m9year + "-" + "0" + strconv.Itoa(int(m9.Month())) + "-" + "01"
		m9enddate = m9year + "-" + "0" + strconv.Itoa(int(m9.Month())) + "-" + m9end

	} else {

		m9startdate = m9year + "-" + strconv.Itoa(int(m9.Month())) + "-" + "01"
		m9enddate = m9year + "-" + strconv.Itoa(int(m9.Month())) + "-" + m9end

	}

	row9, err := db.Db.Query(`SELECT sum(units)  FROM ProducerForecast where start_date>=?  and end_date <= ? and product_name=? and version =?`, m9startdate, m9enddate, input_forecasted_product, input_producer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 9 sum")
		return
	}
	defer sql.Close(row9)

	_, data9, err := sql.Scan(row9)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 9 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "pm9 - ", data9[0][0])
	c.Data["pm9"] = data9[0][0]

	custrow9, err := db.Db.Query(`SELECT sum(units)  FROM CustomerForecast where start_date>=?  and end_date <= ? and product_name=? and version =?`, m9startdate, m9enddate, input_forecasted_product2, input_customer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 9 sum")
		return
	}
	defer sql.Close(custrow9)

	_, custdata9, err := sql.Scan(custrow9)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 9 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cm9 - ", custdata9[0][0])
	c.Data["cm9"] = custdata9[0][0]

	// Month10

	if m10.Month() == 1 || m10.Month() == 3 || m10.Month() == 5 || m10.Month() == 7 || m10.Month() == 8 || m10.Month() == 10 || m10.Month() == 12 {

		m10end = "31"

	} else {

		m10end = "30"

	}

	if len(strconv.Itoa(int(m10.Month()))) == 1 {

		m10startdate = m10year + "-" + "0" + strconv.Itoa(int(m10.Month())) + "-" + "01"
		m10enddate = m10year + "-" + "0" + strconv.Itoa(int(m10.Month())) + "-" + m10end

	} else {

		m10startdate = m10year + "-" + strconv.Itoa(int(m10.Month())) + "-" + "01"
		m10enddate = m10year + "-" + strconv.Itoa(int(m10.Month())) + "-" + m10end

	}

	row10, err := db.Db.Query(`SELECT sum(units)  FROM ProducerForecast where start_date>=?  and end_date <= ? and product_name=? and version =?`, m10startdate, m10enddate, input_forecasted_product, input_producer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 10 sum")
		return
	}
	defer sql.Close(row10)

	_, data10, err := sql.Scan(row10)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 10 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "pm10 - ", data10[0][0])
	c.Data["pm10"] = data10[0][0]

	custrow10, err := db.Db.Query(`SELECT sum(units)  FROM CustomerForecast where start_date>=?  and end_date <= ? and product_name=? and version =?`, m10startdate, m10enddate, input_forecasted_product2, input_customer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 10 sum")
		return
	}
	defer sql.Close(custrow10)

	_, custdata10, err := sql.Scan(custrow10)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 10 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cm10 - ", custdata10[0][0])
	c.Data["cm10"] = custdata10[0][0]

	// Month11

	if m11.Month() == 1 || m11.Month() == 3 || m11.Month() == 5 || m11.Month() == 7 || m11.Month() == 8 || m11.Month() == 10 || m11.Month() == 12 {

		m11end = "31"

	} else {

		m11end = "30"

	}

	if len(strconv.Itoa(int(m11.Month()))) == 1 {

		m11startdate = m11year + "-" + "0" + strconv.Itoa(int(m11.Month())) + "-" + "01"
		m11enddate = m11year + "-" + "0" + strconv.Itoa(int(m11.Month())) + "-" + m11end

	} else {

		m11startdate = m11year + "-" + strconv.Itoa(int(m11.Month())) + "-" + "01"
		m11enddate = m11year + "-" + strconv.Itoa(int(m11.Month())) + "-" + m11end

	}

	row11, err := db.Db.Query(`SELECT sum(units)  FROM ProducerForecast where start_date>=?  and end_date <= ? and product_name=? and version =?`, m11startdate, m11enddate, input_forecasted_product, input_producer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 11 sum")
		return
	}
	defer sql.Close(row11)

	_, data11, err := sql.Scan(row11)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 11 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "pm11 - ", data11[0][0])
	c.Data["pm11"] = data11[0][0]

	custrow11, err := db.Db.Query(`SELECT sum(units)  FROM CustomerForecast where start_date>=?  and end_date <= ? and product_name=? and version =?`, m11startdate, m11enddate, input_forecasted_product2, input_customer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 11 sum")
		return
	}
	defer sql.Close(custrow11)

	_, custdata11, err := sql.Scan(custrow11)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get CustomerProduct Month 11 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cm11 - ", custdata11[0][0])
	c.Data["cm11"] = custdata11[0][0]

	// Month12

	if m12.Month() == 1 || m12.Month() == 3 || m12.Month() == 5 || m12.Month() == 7 || m12.Month() == 8 || m12.Month() == 10 || m12.Month() == 12 {

		m12end = "31"

	} else {

		m12end = "30"

	}

	if len(strconv.Itoa(int(m12.Month()))) == 1 {

		m12startdate = m12year + "-" + "0" + strconv.Itoa(int(m12.Month())) + "-" + "01"
		m12enddate = m12year + "-" + "0" + strconv.Itoa(int(m12.Month())) + "-" + m12end

	} else {

		m12startdate = m12year + "-" + strconv.Itoa(int(m12.Month())) + "-" + "01"
		m12enddate = m12year + "-" + strconv.Itoa(int(m12.Month())) + "-" + m12end

	}

	row12, err := db.Db.Query(`SELECT sum(units)  FROM ProducerForecast where start_date>=?  and end_date <= ? and product_name=? and version =?`, m12startdate, m12enddate, input_forecasted_product, input_producer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 12 sum")
		return
	}
	defer sql.Close(row12)

	_, data12, err := sql.Scan(row12)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get ProducerProduct Month 12 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "pm12 - ", data12[0][0])
	c.Data["pm12"] = data12[0][0]

	custrow12, err := db.Db.Query(`SELECT sum(units)  FROM CustomerForecast where start_date>=?  and end_date <= ? and product_name=? and version =?`, m12startdate, m12enddate, input_forecasted_product2, input_customer_version)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get  CustomerProduct Month 12 sum")
		return
	}
	defer sql.Close(custrow12)

	_, custdata12, err := sql.Scan(custrow12)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get  CustomerProduct Month 12 sum")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cm12 - ", custdata12[0][0])
	c.Data["cm12"] = custdata12[0][0]

	forecastdatapayload := ForecastData{

		Month1:      month1 + " " + m1year,
		Month2:      month2 + " " + m2year,
		Month3:      month3 + " " + m3year,
		Month4:      month4 + " " + m4year,
		Month5:      month5 + " " + m5year,
		Month6:      month6 + " " + m6year,
		Month7:      month7 + " " + m7year,
		Month8:      month8 + " " + m8year,
		Month9:      month9 + " " + m9year,
		Month10:     month10 + " " + m10year,
		Month11:     month11 + " " + m11year,
		Month12:     month12 + " " + m12year,
		PM1:         data1[0][0],
		PM2:         data2[0][0],
		PM3:         data3[0][0],
		PM4:         data4[0][0],
		PM5:         data5[0][0],
		PM6:         data6[0][0],
		PM7:         data7[0][0],
		PM8:         data8[0][0],
		PM9:         data9[0][0],
		PM10:        data10[0][0],
		PM11:        data11[0][0],
		PM12:        data12[0][0],
		CM1:         custdata1[0][0],
		CM2:         custdata2[0][0],
		CM3:         custdata3[0][0],
		CM4:         custdata4[0][0],
		CM5:         custdata5[0][0],
		CM6:         custdata6[0][0],
		CM7:         custdata7[0][0],
		CM8:         custdata8[0][0],
		CM9:         custdata9[0][0],
		CM10:        custdata10[0][0],
		CM11:        custdata11[0][0],
		CM12:        custdata12[0][0],
		Product1:    input_forecasted_product,
		Product2:    input_forecasted_product2,
		Product1Qty: pdctqty1,
	}

	var jsonData2 []byte
	jsonData2, err = json.Marshal(forecastdatapayload)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	}

	payload = f(string(jsonData2))
	fmt.Println("forecast payload in return", payload)

	return
}

func f(key interface{}) (jsonpayload string) {
	json, err := json.Marshal(key)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(string(json))

	st, err := strconv.Unquote(string(json))
	jsonpayload = st
	// fmt.Println(st)
	return

}
