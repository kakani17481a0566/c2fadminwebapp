package updateCustomer

import (
	"c2f_admin_webapp/session"

	"runtime/debug"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"
	"errors"

	"strconv"
	"strings"

	"c2f_admin_webapp/model/db"
	"c2f_admin_webapp/services"
	"fmt"

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
	Usertype    string
	Status      string
	ClusterName string
}

type UpdateCustomer struct {
	// beego.Controller
	middleware.BaseController
}

type Display struct {
	Fields []Field
}
type Field struct {
	Id    string
	Name  string
	Email string
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

func (c *UpdateCustomer) Get() {
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
			c.TplName = "customer/updateCustomer/updateCustomer.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "customer/updateCustomer/updateCustomer.html"
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

		data, err := services.GetClusteruser()
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster user fetch Failed")
			return
		}
		if len(data) <= 0 {
			err = errors.New("Cluster user  Not Found")
			return
		}
		var Dis Display
		for i := range data {
			var d Field
			//d.Id = data[i][0]
			d.Name = data[i][1]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)
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

func (c *UpdateCustomer) Post() {

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
			c.TplName = "customer/updateCustomer/updateCustomer.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Customer Updated Successfully"
			c.TplName = "customer/updateCustomer/updateCustomer.html"
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

	auth, err := utils.IsAuthorized(utype, "Customer")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	row1, err := db.Db.Query(`select id,first_name,middle_name,last_name,mob_no,email,status,user_type,address1,address1,city,zipcode,state,language,town,location_type,location_info,cluster_name,is_default from Customer WHERE id =?  limit 1`, AdminId)
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
		c.Data["Address1"] = data1[i][8]
		c.Data["Address2"] = data1[i][9]
		c.Data["City"] = data1[i][10]
		c.Data["Zipcode"] = data1[i][11]
		c.Data["State"] = data1[i][12]
		c.Data["Language"] = data1[i][13]
		c.Data["Town"] = data1[i][14]
		c.Data["LocationType"] = data1[i][15]
		c.Data["LocationInfo"] = data1[i][16]
		c.Data["ClusterName"] = data1[i][17]
		c.Data["IsDefaultMain"] = data1[i][18]

		data1, err := services.GetClusteruser()
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster user fetch Failed")
			return
		}
		if len(data1) <= 0 {
			err = errors.New("Cluster user  Not Found")
			return
		}
		var Dis Display
		for i := range data1 {
			var d Field
			//d.Id = data[i][0]
			d.Name = data1[i][1]
			Dis.Fields = append(Dis.Fields, d)
		}
		c.Data["Dis"] = Dis
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)
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

	input_first_name := c.Input().Get("input_first_name")
	input_middle_name := c.Input().Get("input_middle_name")
	input_last_name := c.Input().Get("input_last_name")
	input_mobile := c.Input().Get("input_mobile")
	input_email := c.Input().Get("input_email")
	input_user_status := c.Input().Get("input_user_status")

	input_address1 := c.Input().Get("input_address1")
	input_address2 := c.Input().Get("input_address2")
	input_town := c.Input().Get("input_town")
	input_city := c.Input().Get("input_city")
	input_zipcode := c.Input().Get("input_zipcode")
	input_state := c.Input().Get("input_state")
	input_language := c.Input().Get("input_language")
	LocationType := c.Input().Get("LocationType")
	LocationInfo := c.Input().Get("LocationInfo")
	select_cluster_name := c.Input().Get("select_cluster_name")
	input_user_type := c.Input().Get("input_user_type")
	default_address := c.Input().Get("default_address")

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "default_address--------------- - ", default_address)

	if input_user_status == "" {
		input_user_status = "ACTIVE"
	}
	if input_first_name == "" || input_mobile == "" || input_email == "" || input_address1 == "" || input_address2 == "" ||
		input_town == "" || input_city == "" || input_zipcode == "" || input_state == "" {
		err = errors.New("* fields cannot be empty")
		return
	}

	if input_language == "" || LocationType == "" || LocationInfo == "" || select_cluster_name == "" || input_user_status == "" {
		err = errors.New("* fields cannot be empty")
		return
	}

	res, err := db.Db.Exec(`UPDATE Customer SET 
	first_name = ?, middle_name = ?, last_name = ?, email = ?, mob_no = ?,status=? , user_type=?,
	address1=?,address2=?,town=?,city=?,zipcode=?,state=?,language=?,location_type=?,location_info=?,cluster_name=?,user_type=?,update_date=now()
	WHERE id = ?`, input_first_name, input_middle_name, input_last_name, input_email, input_mobile, input_user_status, "CONSUMER", input_address1,
		input_address2, input_town, input_city, input_zipcode, input_state, input_language, LocationType, LocationInfo, select_cluster_name, input_user_type, AdminId)
	if err != nil {
		err = errors.New("User updation failed")
		return
	}

	i, err := res.RowsAffected()
	if err != nil || i == 0 {
		err = errors.New("User updation failed")
		return
	}

	row, err := db.Db.Query(`select id,first_name,middle_name,last_name,mob_no,email,status,user_type,cluster_name from Customer WHERE id = ? limit 1`, AdminId)
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
		c.Data["ClusterName"] = data[i][8]

	}

	// ✔ After validating main customer fields and before address insert/update
	deletedIDs := c.GetString("deleted_ids")

	if deletedIDs != "" {
		idList := strings.Split(deletedIDs, ",")
		placeholders := strings.TrimRight(strings.Repeat("?,", len(idList)), ",")

		query := `UPDATE CustomerAddress 
          SET status='INACTIVE', updated_date=NOW()
          WHERE id IN (` + placeholders + `) AND customer_id=?`

		args := make([]interface{}, len(idList)+1)
		for i, id := range idList {
			args[i] = id
		}
		args[len(idList)] = AdminId

		_, err := db.Db.Exec(query, args...)
		if err != nil {
			err = errors.New("Failed to delete removed addresses")
			return
		}
	}
	// -------------------- HANDLE DEFAULT ADDRESS SELECTION --------------------

	defaultAddress := c.Input().Get("default_address")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Selected Default Address - ", defaultAddress)

	// Case 1 → MAIN ADDRESS SELECTED AS DEFAULT
	if defaultAddress == "main" {

		// Set main address as default
		_, err := db.Db.Exec(`UPDATE Customer 
                           SET is_default = 1 
                           WHERE id = ?`, AdminId)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", "Failed to set main default:", err.Error())
			return
		}

		// Clear ALL child addresses defaults
		_, err = db.Db.Exec(`UPDATE CustomerAddress 
                         SET is_default = 0 
                         WHERE customer_id = ?`, AdminId)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", "Failed to clear additional defaults:", err.Error())
			return
		}

		// Case 2 → ADDITIONAL ADDRESS SELECTED AS DEFAULT (index ID)
	} else if defaultAddress != "" {

		// Main address becomes NOT default
		_, err := db.Db.Exec(`UPDATE Customer 
                           SET is_default = 0 
                           WHERE id = ?`, AdminId)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", "Failed to unset main:", err.Error())
			return
		}

		// Clear defaults from all additional addresses
		_, err = db.Db.Exec(`UPDATE CustomerAddress 
                         SET is_default = 0 
                         WHERE customer_id = ?`, AdminId)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", "Failed to clear child defaults:", err.Error())
			return
		}

		// Set selected additional address as default
		_, err = db.Db.Exec(`UPDATE CustomerAddress 
                         SET is_default = 1 
                         WHERE id = ? AND customer_id = ?`,
			defaultAddress, AdminId)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", "Failed to set child default:", err.Error())
			return
		}
	}

	// -------------------- Code for additional addresses ------------------------------
	var additionalAddresses []AdditionalAddress

	for i := 0; ; i++ {
		id := c.Input().Get(fmt.Sprintf("addresses[%d][id]", i))
		name := c.Input().Get(fmt.Sprintf("addresses[%d][name]", i))

		if name == "" && id == "" {
			break
		}

		addr := AdditionalAddress{
			Id:                   id,
			Name:                 name,
			Mobile:               c.Input().Get(fmt.Sprintf("addresses[%d][mob_no]", i)),
			Email:                c.Input().Get(fmt.Sprintf("addresses[%d][email]", i)),
			Address1:             c.Input().Get(fmt.Sprintf("addresses[%d][address1]", i)),
			Address2:             c.Input().Get(fmt.Sprintf("addresses[%d][address2]", i)),
			Town:                 c.Input().Get(fmt.Sprintf("addresses[%d][town]", i)),
			City:                 c.Input().Get(fmt.Sprintf("addresses[%d][city]", i)),
			State:                c.Input().Get(fmt.Sprintf("addresses[%d][state]", i)),
			Zipcode:              c.Input().Get(fmt.Sprintf("addresses[%d][zipcode]", i)),
			LocationType:         c.Input().Get(fmt.Sprintf("addresses[%d][location_type]", i)),
			MultipleLocationInfo: c.Input().Get(fmt.Sprintf("addresses[%d][multiple_location_info]", i)),
			IsDefault:            c.Input().Get("default_address") == fmt.Sprintf("%d", i),
		}

		additionalAddresses = append(additionalAddresses, addr)

		// Backend validation for additional address fields
		if addr.Name == "" || addr.Mobile == "" || addr.Email == "" ||
			addr.Address1 == "" || addr.Address2 == "" || addr.Town == "" ||
			addr.City == "" || addr.State == "" || addr.Zipcode == "" ||
			addr.LocationType == "" || addr.MultipleLocationInfo == "" {

			err = errors.New("All address fields are required. No field should be empty.")
			return
		}

		// Mobile validation
		if len(addr.Mobile) != 10 {
			err = errors.New("Mobile number must be 10 digits")
			return
		}

		// Pincode validation
		if len(addr.Zipcode) != 6 {
			err = errors.New("Pincode must be 6 digits")
			return
		}

		// Email validation basic
		if !strings.Contains(addr.Email, "@") {
			err = errors.New("Invalid email format")
			return
		}

	}

	// Count existing rows for insert limit
	var existingCount int
	countRows, err := db.Db.Query(`SELECT COUNT(*) FROM CustomerAddress WHERE customer_id = ? AND status='ACTIVE'`, AdminId)
	if err == nil {
		defer sql.Close(countRows)
		_, countData, _ := sql.Scan(countRows)
		if len(countData) > 0 && len(countData[0]) > 0 {
			existingCount, _ = strconv.Atoi(countData[0][0])
		}
	}

	// // Update default address: make previous default = 0
	// db.Db.Exec(`UPDATE CustomerAddress SET is_default=0 WHERE customer_id=?`, AdminId)

	// Insert or Update Logic
	for _, addr := range additionalAddresses {

		if addr.Id != "" {

			// UPDATE using ID
			_, execErr := db.Db.Exec(`UPDATE CustomerAddress SET 
            name=?, mob_no=?, email=?, address1=?, address2=?, town=?, city=?, state=?, zipcode=?,
            location_type=?, location_info=?, is_default=?, updated_date=NOW()
            WHERE id=? AND customer_id=?`,
				addr.Name, addr.Mobile, addr.Email, addr.Address1, addr.Address2, addr.Town,
				addr.City, addr.State, addr.Zipcode, addr.LocationType,
				addr.MultipleLocationInfo, addr.IsDefault, addr.Id, AdminId)

			if execErr != nil {
				err = errors.New("Failed to update address")
				return
			}

			fmt.Println("Updated:", addr.Name)
		} else {

			// INSERT new address
			if existingCount >= 4 {
				err = errors.New("Maximum 4 additional addresses allowed")
				return
			}

			_, execErr := db.Db.Exec(`INSERT INTO CustomerAddress 
            (customer_id, name, mob_no, email, address1, address2, town, city, state, zipcode,
            location_type, location_info, is_default, status, updated_date)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'ACTIVE', NOW())`,
				AdminId, addr.Name, addr.Mobile, addr.Email, addr.Address1, addr.Address2,
				addr.Town, addr.City, addr.State, addr.Zipcode,
				addr.LocationType, addr.MultipleLocationInfo, addr.IsDefault)

			if execErr != nil {
				err = errors.New("Failed to insert new address")
				return
			}

			existingCount++
			fmt.Println("Inserted:", addr.Name)
		}
	}

	return
}
