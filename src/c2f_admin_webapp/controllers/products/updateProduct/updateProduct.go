package updateProduct

import (
	"c2f_admin_webapp/session"

	"runtime/debug"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"
	"errors"
	"strings"

	"c2f_admin_webapp/model/db"

	"github.com/astaxie/beego"

	"c2f_admin_webapp/services"
	"path"
	"strconv"

	"c2f_admin_webapp/model/database/sql"
	"c2f_admin_webapp/model/util/log"
)

type Row struct {
	Id              string
	CategoryName    string
	SubCategoryName string
	ProductName     string
	ProductDesc     string
	Productstatus   string
	ProductQuantity string
	QualityType     string
	QualityMeasure  string
	ShelLife        string
	PackageType     string
	PackageQty      string
	ProductOrigin   string

	Kselprice string
	kbuyprice string
	kmarprice string

	comname   string
	cgst      string
	sgst      string
	hsn       string
	cess      string
	othertax  string
	CreatedBy string
}

type UpdateProduct struct {
	// beego.Controller
	middleware.BaseController
}

type GetGST struct {
	// beego.Controller
	middleware.BaseController
}

type Display3 struct {
	Fields3 []Field3
}
type Field3 struct {
	Id             string
	Commodity_Name string
}

func (c *UpdateProduct) Get() {
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
			c.TplName = "products/updateProduct/updateProduct.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Update Admin User Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "products/updateProduct/updateProduct.html"
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

	data3, err := services.GetCommodityName()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Commodity_Name fetch Failed")
		return
	}
	if len(data3) <= 0 {
		err = errors.New("Commodity_Name   Not Found")
		return
	}
	var Dis3 Display3
	for i := range data3 {
		var d Field3
		d.Id = data3[i][0]
		d.Commodity_Name = data3[i][1]
		Dis3.Fields3 = append(Dis3.Fields3, d)
	}
	c.Data["Dis3"] = Dis3
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis3.Fields3)

	auth, err := utils.IsAuthorized(utype, "Product")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	row, err := db.Db.Query(`select id,category_name,subcategory_name,product_name,product_desc,product_status, 
product_qty,
product_quality_type,
product_quality_mesure,
product_origin, 
product_package_type,
product_package_qty,
product_shelf_life,
krishivala_selling_price,
krishivala_buying_price,
krishivala_market_price,
commodity_name,
krishivala_CGST_charge,
krishivala_SGST_charge,
hsn_code,
krishivala_CESS_charge,
krishivala_OtherTAX_charge ,
JSON_EXTRACT(images_path,'$.img1'),
JSON_EXTRACT(images_path,'$.img2'),
JSON_EXTRACT(images_path,'$.img3'),
JSON_EXTRACT(images_path,'$.img4'),
JSON_EXTRACT(images_path,'$.video1'),
product_shelf_life_unit
from Lookup_products WHERE id =?  limit 1`, AdminId)
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
		c.Data["CategoryName"] = data[i][1]
		c.Data["SubCategoryName"] = data[i][2]
		c.Data["ProductName"] = data[i][3]
		c.Data["ProductDesc"] = data[i][4]
		c.Data["Status"] = data[i][5]
		c.Data["ProductQuantity"] = data[i][6]
		c.Data["QualityType"] = data[i][7]
		c.Data["QualityMeasure"] = data[i][8]
		c.Data["ShelLife"] = data[i][12]
		c.Data["PackageType"] = data[i][10]
		c.Data["PackageQty"] = data[i][11]
		c.Data["ProductOrigin"] = data[i][9]

		c.Data["Kselprice"] = data[i][13]
		c.Data["kbuyprice"] = data[i][14]
		c.Data["kmarprice"] = data[i][15]

		c.Data["comname"] = data[i][16]
		c.Data["cgst"] = data[i][17]
		c.Data["sgst"] = data[i][18]
		c.Data["hsn"] = data[i][19]
		c.Data["cess"] = data[i][20]
		c.Data["othertax"] = data[i][21]
		c.Data["Img1"], err = strconv.Unquote(data[i][22])
		c.Data["Img2"], err = strconv.Unquote(data[i][23])
		c.Data["Img3"], err = strconv.Unquote(data[i][24])
		c.Data["Img4"], err = strconv.Unquote(data[i][25])
		c.Data["Video1"], err = strconv.Unquote(data[i][26])
		c.Data["ShelLifeUnit"] = data[i][27]

	}

	return

}

func (c *GetGST) Post() {
	krishivalaCGSTcharge := ""
	krishivalaSGSTcharge := ""
	hsncode := ""
	//CusType := ""
	mobile := ""
	log.Println(beego.AppConfig.String("loglevel"), "Info", "Verify Customer Page Start")
	pip := c.Ctx.Input.IP()
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Client IP - ", pip)
	var err error
	defer func() {
		if l_exception := recover(); l_exception != nil {
			c.Ctx.Output.Body([]byte(`{"Tittle":"EXCEPTION","Message": "System Exception","Type": "error"}`))
			stack := debug.Stack()
			log.Println(beego.AppConfig.String("loglevel"), "Exception", string(stack))
			session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)

		}
		if err != nil {
			c.Ctx.Output.Body([]byte(`{"Tittle":"ERROR","Message": "` + err.Error() + `","Type": "error"}`))

			log.Println(beego.AppConfig.String("loglevel"), "Info", "Create User Page Fail")
		} else {
			c.Ctx.Output.Body([]byte(`{"Tittle":"SUCCESS","Message": "Customer Account verification is successful. Please proceed and enter the Transaction Amount.","Type": "success","KrishivalaCGSTcharge":"` + krishivalaCGSTcharge + `","KrishivalaSGSTcharge":"` + krishivalaSGSTcharge + `","HSNcode":"` + hsncode + `","Mobile":"` + mobile + `"}`))
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Create User Page Success")
		}
		return
	}()
	utils.SetHTTPHeader(c.Ctx)
	sess, err := session.GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("System is unable to process your request.Please contact customer care")
		return
	}
	if err = session.ValidateSession(sess); err != nil {
		sess.SessionRelease(c.Ctx.ResponseWriter)
		session.GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}
	defer sess.SessionRelease(c.Ctx.ResponseWriter)

	cname := c.Input().Get("input_commodity_name")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "cname - ", cname)

	if cname == "" {
		log.Println(beego.AppConfig.String("loglevel"), "Error", "Feilds Blank")
		err = errors.New("cname can't be blank")
		return
	}

	row, err := db.Db.Query(`select id,commodity_name,cgst_rate,sgst_rate,hsn_code from GST_Rates WHERE commodity_name =?  limit 1`, cname)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get commodity_name& GST data")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to commodity_name& GST data")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("commodity_name& GST data not found")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)
	for i := range data {
		//c.Data["Id"] = data[i][0]
		krishivalaCGSTcharge = data[i][2]
		krishivalaSGSTcharge = data[i][3]
		hsncode = data[i][4]
	}
	// uname := sess.Get("uname").(string)

	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "User Name - ", uname)

	// krishivalaCGSTcharge = "123"
	// krishivalaSGSTcharge = "222"
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "userName - ", krishivalaCGSTcharge)

	return
}

func (c *UpdateProduct) Post() {

	AdminId := c.Ctx.Input.Param(":AdminID")
	//krishivalaCGSTcharge := ""
	//	krishivalaSGSTcharge := ""
	//hsncode := ""
	//CusType := ""

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
			c.TplName = "products/updateProduct/updateProduct.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Product Updated Successfully"
			c.TplName = "products/updateProduct/updateProduct.html"
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

	id := sess.Get("uid").(string)
	c.Data["id"] = id
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "id", id)

	auth, err := utils.IsAuthorized(utype, "Product")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	row, err := db.Db.Query(`select product_sku,JSON_EXTRACT(images_path,'$.img2'), JSON_EXTRACT(images_path,'$.img3'), JSON_EXTRACT(images_path,'$.img4') from Lookup_products WHERE id =?`, AdminId)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get product data")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get product data")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IMAGE Query Data - ", data, "\nData len - ", len(data))

	product_sku := data[0][0]
	// img2 := data[0][1]
	// img3 := data[0][2]
	// img4 := data[0][3]

	product_sku_1 := strings.ReplaceAll(product_sku, "/", "_")

	//input_product_qty := c.Input().Get("input_product_qty")
	//input_QualityType := c.Input().Get("input_QualityType")
	//input_QualityMeasure := c.Input().Get("input_QualityMeasure")
	input_product_status := c.Input().Get("input_product_status")
	input_ShelLife := c.Input().Get("input_ShelLife")
	//	input_PackageQty := c.Input().Get("input_PackageQty")
	input_PackageType := c.Input().Get("input_PackageType")
	input_ProductOrigin := c.Input().Get("input_ProductOrigin")
	input_product_desc := c.Input().Get("input_product_desc")
	input_krishivala_selling_price := c.Input().Get("input_krishivala_selling_price")
	input_krishivala_buying_price := c.Input().Get("input_krishivala_buying_price")
	input_krishivala_market_price := c.Input().Get("input_krishivala_market_price")
	input_commodity_name := c.Input().Get("input_commodity_name")
	krishivala_CGST_charge := c.Input().Get("krishivala_CGST_charge")
	krishivala_SGST_charge := c.Input().Get("krishivala_SGST_charge")
	hsn_code := c.Input().Get("hsn_code")
	krishivala_CESS_charge := c.Input().Get("krishivala_CESS_charge")
	krishivala_OtherTAX_charge := c.Input().Get("krishivala_OtherTAX_charge")
	input_ShelLifeUnit := c.Input().Get("input_ShelLifeUnit")

	// row, err := db.Db.Query(`select id,category_name,subcategory_name,product_name,product_desc,product_status,
	// product_qty,
	// product_quality_type,
	// product_quality_mesure,
	// product_origin,
	// product_package_type,
	// product_package_qty,
	// product_shelf_life,
	// krishivala_selling_price,
	// krishivala_buying_price,
	// krishivala_market_price,
	// commodity_name,
	// krishivala_CGST_charge,
	// krishivala_SGST_charge,
	// hsn_code,
	// krishivala_CESS_charge,
	// krishivala_OtherTAX_charge
	// from Lookup_products WHERE id =?  limit 1`, AdminId)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Unable to get product data")
	// 	return
	// }
	// defer sql.Close(row)
	// _, data, err := sql.Scan(row)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// 	err = errors.New("Unable to get product data")
	// 	return
	// }
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	// if len(data) <= 0 {
	// 	err = errors.New("product data not found")
	// 	return
	// }

	// log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	// for i := range data {

	// 	c.Data["Id"] = data[i][0]
	// 	c.Data["CategoryName"] = data[i][1]
	// 	c.Data["SubCategoryName"] = data[i][2]
	// 	c.Data["ProductName"] = data[i][3]
	// 	c.Data["ProductDesc"] = data[i][4]
	// 	c.Data["Status"] = data[i][5]
	// 	c.Data["ProductQuantity"] = data[i][6]
	// 	c.Data["QualityType"] = data[i][7]
	// 	c.Data["QualityMeasure"] = data[i][8]
	// 	c.Data["ShelLife"] = data[i][12]
	// 	c.Data["PackageType"] = data[i][10]
	// 	c.Data["PackageQty"] = data[i][11]
	// 	c.Data["ProductOrigin"] = data[i][9]

	// 	c.Data["Kselprice"] = data[i][13]
	// 	c.Data["kbuyprice"] = data[i][14]
	// 	c.Data["kmarprice"] = data[i][15]

	// 	c.Data["comname"] = data[i][16]
	// 	c.Data["cgst"] = data[i][17]
	// 	c.Data["sgst"] = data[i][18]
	// 	c.Data["hsn"] = data[i][19]
	// 	c.Data["cess"] = data[i][20]
	// 	c.Data["othertax"] = data[i][21]

	// }

	if input_product_status == "" || input_ShelLife == "" || input_PackageType == "" || input_ProductOrigin == "" || input_product_desc == "" ||
		input_krishivala_selling_price == "" || input_krishivala_buying_price == "" || input_krishivala_market_price == "" || input_commodity_name == "" || krishivala_CGST_charge == "" || krishivala_SGST_charge == "" ||
		hsn_code == "" || krishivala_CESS_charge == "" || krishivala_OtherTAX_charge == "" {
		err = errors.New("* fields  cannot be empty")
		return
	}

	file1, information1, err1 := c.GetFile("myFile1")
	product_image_path := beego.AppConfig.String("PRODUCT_IMAGE_PATH")
	if err1 != nil {

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "-", "File1 retrieval failure")

	} else {

		defer file1.Close()
		filename1 := information1.Filename
		// For Local Path

		err1 = c.SaveToFile("myFile1", path.Join(product_image_path+product_sku_1, filename1))
		if err1 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "-", "File1 upload failed!")
		} else {
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "-", "File1 upload succeed!")
		}
		filepath1 := "/static/dist/img/product_images/" + product_sku_1 + "/" + filename1
		res, err1 := db.Db.Exec(`UPDATE Lookup_products SET images_path=JSON_SET(images_path,'$.img1',?)  WHERE id=?`, filepath1, AdminId)
		if err1 != nil {
			err1 = errors.New("Product image1 path updation failed")
			return
		}
		i, err1 := res.RowsAffected()
		if err1 != nil || i == 0 {
			err1 = errors.New("Product image1  path updation failed")

			return
		}
	}

	file2, information2, err2 := c.GetFile("myFile2")
	if err2 != nil {

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "-", "File2 retrieval failure")

	} else {

		defer file2.Close()
		filename2 := information2.Filename

		err2 = c.SaveToFile("myFile2", path.Join(product_image_path+product_sku_1, filename2))
		if err2 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "-", "File2 upload failed!")
		} else {
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "-", "File2 upload succeed!")
		}
		filepath2 := "/static/dist/img/product_images/" + product_sku_1 + "/" + filename2
		res, err1 := db.Db.Exec(`UPDATE Lookup_products SET images_path=JSON_SET(images_path,'$.img2',?)  WHERE id=?`, filepath2, AdminId)
		if err1 != nil {
			err1 = errors.New("Product image2 path updation failed")
			return
		}
		i, err1 := res.RowsAffected()
		if err1 != nil || i == 0 {
			err1 = errors.New("Product image2  path updation failed")

			return
		}
	}

	file3, information3, err3 := c.GetFile("myFile3")
	if err3 != nil {

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "-", "File3 retrieval failure")

	} else {

		defer file3.Close()
		filename3 := information3.Filename

		err3 = c.SaveToFile("myFile3", path.Join(product_image_path+product_sku_1, filename3))
		if err3 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "-", "File3 upload failed!")
		} else {
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "-", "File3 upload succeed!")
		}
		filepath3 := "/static/dist/img/product_images/" + product_sku_1 + "/" + filename3
		res, err1 := db.Db.Exec(`UPDATE Lookup_products SET images_path=JSON_SET(images_path,'$.img3',?)  WHERE id=?`, filepath3, AdminId)
		if err1 != nil {
			err1 = errors.New("Product image3 path updation failed")
			return
		}
		i, err1 := res.RowsAffected()
		if err1 != nil || i == 0 {
			err1 = errors.New("Product image3  path updation failed")

			return
		}
	}
	file4, information4, err4 := c.GetFile("myFile4")
	if err4 != nil {

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "-", "File4 retrieval failure")

	} else {

		defer file4.Close()
		filename4 := information4.Filename

		err4 = c.SaveToFile("myFile4", path.Join(product_image_path+product_sku_1, filename4))
		if err4 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "-", "File4 upload failed!")
		} else {
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "-", "File4 upload succeed!")
		}
		filepath4 := "/static/dist/img/product_images/" + product_sku_1 + "/" + filename4
		res, err1 := db.Db.Exec(`UPDATE Lookup_products SET images_path=JSON_SET(images_path,'$.img4',?)  WHERE id=?`, filepath4, AdminId)
		if err1 != nil {
			err1 = errors.New("Product image4 path updation failed")
			return
		}
		i, err1 := res.RowsAffected()
		if err1 != nil || i == 0 {
			err1 = errors.New("Product image4  path updation failed")

			return
		}
	}

	file5, information5, err5 := c.GetFile("myFile5")
	if err5 != nil {

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "-", "File5 retrieval failure")

	} else {

		defer file5.Close()
		filename5 := information5.Filename

		err5 = c.SaveToFile("myFile5", path.Join(product_image_path+product_sku_1, filename5))
		if err5 != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "-", "File5 upload failed!")
		} else {
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "-", "File5 upload succeed!")
		}
		filepath5 := "/static/dist/img/product_images/" + product_sku_1 + "/" + filename5
		res, err1 := db.Db.Exec(`UPDATE Lookup_products SET images_path=JSON_SET(images_path,'$.video1',?)  WHERE id=?`, filepath5, AdminId)
		if err1 != nil {
			err1 = errors.New("Product video path updation failed")
			return
		}
		i, err1 := res.RowsAffected()
		if err1 != nil || i == 0 {
			err1 = errors.New("Product video  path updation failed")

			return
		}
	}

	res, err := db.Db.Exec(`UPDATE Lookup_products SET 
product_desc=?,
product_origin=?, 

product_package_type=?,

product_shelf_life=?,

krishivala_selling_price=?,
krishivala_buying_price=?,
krishivala_market_price=?,
commodity_name=?,
krishivala_CGST_charge=?,
krishivala_SGST_charge=?,
hsn_code=?,
krishivala_CESS_charge=?,
krishivala_OtherTAX_charge=? ,
	 product_status = ?,
	product_shelf_life_unit =?,
	updated_by=?,
	last_update=now()
	WHERE id = ?`,

		input_product_desc,
		input_ProductOrigin,

		input_PackageType,

		input_ShelLife,

		input_krishivala_selling_price,
		input_krishivala_buying_price,
		input_krishivala_market_price,
		input_commodity_name,
		krishivala_CGST_charge,
		krishivala_SGST_charge,
		hsn_code,
		krishivala_CESS_charge,
		krishivala_OtherTAX_charge,
		input_product_status,
		input_ShelLifeUnit,
		id,
		AdminId)
	if err != nil {
		err = errors.New("Product updation failed")
		return
	}

	i, err := res.RowsAffected()
	if err != nil || i == 0 {
		err = errors.New("Product updation failed")
		return
	}

	//row, err := db.Db.Query(`select id,category_name,subcategory_name,product_name,product_desc,product_status from Lookup_products WHERE id = ? limit 1`, AdminId)

	return
}
