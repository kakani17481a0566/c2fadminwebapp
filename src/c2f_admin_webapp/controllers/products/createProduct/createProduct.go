package createProduct

import (
	"c2f_admin_webapp/session"

	"runtime/debug"

	"c2f_admin_webapp/middleware"
	"c2f_admin_webapp/model/utils"
	"errors"

	//"fmt"

	// "log"
	//"io/ioutil"
	//"net/http"
	"strings"

	// "io"

	"os"
	"path"

	// "path/filepath"
	// "time"
	"encoding/json"

	"c2f_admin_webapp/model/db"

	"c2f_admin_webapp/services"

	"c2f_admin_webapp/model/database/sql"

	"c2f_admin_webapp/model/util/log"

	"github.com/astaxie/beego"
	//	"c2f_admin_webapp/model/util/sku"
)

type CreateProduct struct {
	// beego.Controller
	middleware.BaseController
}

type GetGST struct {
	// beego.Controller
	middleware.BaseController
}

type UploadProductImage struct {
	// beego.Controller
	middleware.BaseController
}

type Display struct {
	Fields1 []Field1
}
type Field1 struct {
	Id           string
	CategoryName string
}

type Display2 struct {
	Fields2 []Field2
}
type Field2 struct {
	Id              string
	SubCategoryName string
}

type Display3 struct {
	Fields3 []Field3
}
type Field3 struct {
	Id             string
	Commodity_Name string
}
type ImagesPaths struct {
	ProductImagesPaths []string `json:"imgs_paths"`
}

type ProductImages struct {
	Img1   string `json:"img1"`
	Img2   string `json:"img2"`
	Img3   string `json:"img3"`
	Img4   string `json:"img4"`
	Video1 string `json:"video1"`
}

func (c *CreateProduct) Get() {
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
			c.TplName = "products/createProduct/createProduct.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "Search Customer Page Fail")
		} else {
			c.Data["DisplayMessage"] = ""
			c.TplName = "products/createProduct/createProduct.html"
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

	auth, err := utils.IsAuthorized(utype, "Product")
	if !auth {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "UnAuthorized")
		Autherr = errors.New("UnAuthorized")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "IsAuthorized - ", "Authorized")

	data, err := services.GetCategory()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Category fetch Failed")
		return
	}
	if len(data) <= 0 {
		err = errors.New("Category  Not Found")
		return
	}
	var Dis Display
	for i := range data {
		var d Field1
		d.Id = data[i][0]
		d.CategoryName = data[i][1]
		Dis.Fields1 = append(Dis.Fields1, d)
	}
	c.Data["Dis"] = Dis
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis)

	////////////////////////////

	data2, err := services.GetSubCategory()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Sub Category fetch Failed")
		return
	}
	if len(data2) <= 0 {
		err = errors.New("Sub Category  Not Found")
		return
	}
	var Dis1 Display2
	for i := range data2 {
		var d Field2
		d.Id = data2[i][0]
		d.SubCategoryName = data2[i][1]
		Dis1.Fields2 = append(Dis1.Fields2, d)
	}
	c.Data["Dis1"] = Dis1
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Object Data - ", Dis1.Fields2)

	////////////////////

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

	return
}

//////////////

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

////////////////

func (c *CreateProduct) Post() {

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
			c.TplName = "products/createProduct/createProduct.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer Page Fail")
		} else {
			c.Data["DisplayMessage"] = "Product Added Successfully"
			c.TplName = "products/createProduct/createProduct.html"
			log.Println(beego.AppConfig.String("loglevel"), "Info", "SearchCustomer  Page Success")
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

	input_category_name := c.Input().Get("input_category_name")
	input_subcategory_name := c.Input().Get("input_subcategory_name")
	input_product_name := c.Input().Get("input_product_name")

	input_product_qty := c.Input().Get("input_product_qty")
	input_QualityType := c.Input().Get("input_QualityType")
	input_QualityMeasure := c.Input().Get("input_QualityMeasure")

	input_ShelLife := c.Input().Get("input_ShelLife")
	input_ShelLifeUnit := c.Input().Get("input_ShelLifeUnit")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "input_ShelLifeUnit  - ", input_ShelLifeUnit)

	input_PackageType := c.Input().Get("input_PackageType")
	input_PackageQty := c.Input().Get("input_PackageQty")

	input_product_origin := c.Input().Get("input_ProductOrigin")
	input_product_desc := c.Input().Get("input_product_desc")

	input_krishivala_selling_price := c.Input().Get("input_krishivala_selling_price")
	input_krishivala_buying_price := c.Input().Get("input_krishivala_buying_price")
	input_krishivala_market_price := c.Input().Get("input_krishivala_market_price")
	input_product_status := c.Input().Get("input_product_status")
	//input_product_image := c.Input().Get("input_product_image")

	var productsku, productsku_1 string
	productskuBase := strings.ToLower(input_product_name + "-" + input_QualityMeasure)
	productskuBase = strings.ReplaceAll(productskuBase, "%", "percent") // ✅ Replace % with 'percent'

	if input_product_qty == "kg-bulk" {
		productsku = productskuBase + "-bulk"
		productsku_1 = productsku
	} else {
		safeQty := strings.ReplaceAll(input_PackageQty, "/", "-")
		productsku = productskuBase + "-" + safeQty
		productsku_1 = productsku
	}

	krishivala_CGST_charge := c.Input().Get("krishivala_CGST_charge")
	krishivala_SGST_charge := c.Input().Get("krishivala_SGST_charge")
	krishivala_CESS_charge := c.Input().Get("krishivala_CESS_charge")
	krishivala_OtherTAX_charge := c.Input().Get("krishivala_OtherTAX_charge")

	input_commodity_name := c.Input().Get("input_commodity_name")
	hsncode := c.Input().Get("hsn_code")

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "category_name  - ", input_category_name)

	if input_category_name == "" || input_product_name == "" {
		err = errors.New("cluster name and product name cannot be empty")
		return
	}
	if krishivala_CESS_charge == "" || krishivala_OtherTAX_charge == "" {
		err = errors.New("CESS charge and OtherTAX cannot be empty")
		return
	}

	err = CheckExists(input_category_name, input_subcategory_name, input_product_name, input_PackageQty)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Product already Exists")
		return
	}

	// Here It Creates a New Directory with the name as Product SKU

	product_image_path := beego.AppConfig.String("PRODUCT_IMAGE_PATH")

	err1 := os.Mkdir(product_image_path+productsku_1, 0755)

	if err != nil {
		// log.Fatal(err)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "err1", err1)
	}

	file, information, err := c.GetFile("myFile")

	if err != nil {
		// log.Println("File retrieval failure")
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "-", "File retrieval failure")
		return
	}
	defer file.Close()

	filename := information.Filename

	err = c.SaveToFile("myFile", path.Join(product_image_path+productsku_1, filename))

	if err != nil {
		// log.Println("File upload failed!")
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "-", "File upload failed!")
	} else {
		// log.Println("File upload succeed!")
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "-", "File upload succeed!")
		// c.Data["DisplayMessage"] = "File upload succeed!"
		//Display information after successful upload
	}

	filepath := "/static/dist/img/product_images/" + productsku_1 + "/" + filename

	demoimg := "/static/dist/img/pdct_demo_image.png"
	demovideo := ""

	productimagedata := ProductImages{
		Img1:   filepath,
		Img2:   demoimg,
		Img3:   demoimg,
		Img4:   demoimg,
		Video1: demovideo,
	}

	var jsonData2 []byte
	jsonData2, err = json.Marshal(productimagedata)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "product image json Data - ", string(jsonData2))

	// b := new(strings.Builder)
	// json.NewEncoder(b).Encode(ProductImages)
	// // print(b.String())
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "product images json Data - ", b.String())

	// imagesdata := ImagesPaths{
	// 	ProductImagesPaths: ProductImages,
	// }

	// var jsonData2 []byte
	// jsonData2, err = json.Marshal(imagesdata)
	// if err != nil {
	// 	log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	// }
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "product images json Data - ", string(jsonData2))

	result, err := db.Db.Exec(`INSERT INTO Lookup_products (
	category_name, 
  	subcategory_name,
 	product_name, 
	
	product_qty,
	product_quality_type,
	product_quality_mesure,
	
	product_shelf_life,
    	product_package_type,
	product_package_qty,
	
	product_origin, 
	product_desc,
	krishivala_buying_price,
	
	krishivala_selling_price,
 	krishivala_market_price,
 	product_sku,
	
	
 	
	product_status,
	krishivala_CGST_charge,
	krishivala_SGST_charge,
	krishivala_CESS_charge,
	krishivala_OtherTAX_charge,
	commodity_name,
	hsn_code,
	images_path,
	product_shelf_life_unit,
	created_by,
 	create_date)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,now())`,
		input_category_name,
		input_subcategory_name,
		input_product_name,

		input_product_qty,
		input_QualityType,
		input_QualityMeasure,

		input_ShelLife,
		input_PackageType,
		input_PackageQty,

		input_product_origin,
		input_product_desc,
		input_krishivala_selling_price,

		input_krishivala_buying_price,
		input_krishivala_market_price,
		productsku,
		input_product_status,
		krishivala_CGST_charge,
		krishivala_SGST_charge,
		krishivala_CESS_charge,
		krishivala_OtherTAX_charge,

		input_commodity_name,
		hsncode,
		jsonData2,
		input_ShelLifeUnit,
		id)
	if err != nil {
		err = errors.New("Product creation failed")
		return
	}

	i, err := result.RowsAffected()
	if err != nil || i == 0 {
		err = errors.New("Product creation failed")
		return
	}

	return
}

func CheckExists(input_category_name, input_subcategory_name, input_product_name, input_PackageQty string) (err error) {

	err = nil

	row, err := db.Db.Query(`select count(*) from Lookup_products WHERE category_name = ? AND subcategory_name =? AND product_name=? AND product_package_qty=?`, input_category_name, input_subcategory_name, input_product_name, input_PackageQty)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch product")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to fetch product")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", data)

	countlen := data[0][0]

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "countlen", countlen)

	if countlen != "0" {
		err = errors.New("product already exists")
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		return
	}

	return

}
