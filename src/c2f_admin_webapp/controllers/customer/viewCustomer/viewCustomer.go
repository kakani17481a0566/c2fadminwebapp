package viewCustomer

import (
	"c2f_admin_webapp/session"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"
	"runtime/debug"

	"errors"
	"strings"

	"c2f_admin_webapp/model/db"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"
)

type Row struct {
	Id          string
	FirstName   string
	MiddleName  string
	LastName    string
	Mobile      string
	Email       string
	CompanyName string
	Usertype    string
	Status      string
	ClusterName string
}
type ViewCustomer struct {
	// beego.Controller
	middleware.BaseController
}

type AdditionalAddress struct {
	Id                   string
	Name                 string
	Mobile               string
	Email                string
	Address1             string
	Address2             string
	Town                 string
	City                 string
	State                string
	Zipcode              string
	LocationType         string
	MultipleLocationInfo string
	IsDefault            bool
}


func (c *ViewCustomer) Get() {
	AdminId := c.Ctx.Input.Param(":AdminID")
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
			c.TplName = "customer/viewCustomer/viewCustomer.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "customer/viewCustomer/viewCustomer.html"
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

	row, err := db.Db.Query(`select id,first_name,middle_name,last_name,mob_no,email,status,user_type,address1,address2,city,zipcode,state,language,town,location_type,location_info,cluster_name,is_default from Customer WHERE id =?  limit 1`, AdminId)
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
		c.Data["ClusterName"] = data[i][17]
		c.Data["IsDefaultMain"] = data[i][18]

	}

	// --- Fetch additional addresses for this customer ---
	addressRows, err := db.Db.Query(`SELECT 
    name, mob_no, email, address1, address2, town, city, state, zipcode, location_type, location_info, is_default ,id
    FROM CustomerAddress WHERE customer_id = ? AND status='ACTIVE'`, AdminId)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", "Failed to fetch additional addresses:", err)
	} else {
		defer sql.Close(addressRows)
		_, addrData, _ := sql.Scan(addressRows)

		var additionalAddresses []AdditionalAddress
		for _, row := range addrData {
			addr := AdditionalAddress{
				Name:                 row[0],
				Mobile:               row[1],
				Email:                row[2],
				Address1:             row[3],
				Address2:             row[4],
				Town:                 row[5],
				City:                 row[6],
				State:                row[7],
				Zipcode:              row[8],
				LocationType:         row[9],
				MultipleLocationInfo: row[10],
				IsDefault:            row[11] == "1" || strings.ToLower(row[11]) == "true",
				Id:                   row[12],
			}
			additionalAddresses = append(additionalAddresses, addr)
		}

		// Pass the addresses to template
		if len(additionalAddresses) == 0 {
			c.Data["Addresses"] = []AdditionalAddress{}
		} else {
			c.Data["Addresses"] = additionalAddresses
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Fetched additional addresses:", additionalAddresses)
	}

	return

}
