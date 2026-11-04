package payuresponse

import (
	//"c2f_admin_webapp/model/general"
	//	"c2f_admin_webapp/model/msg91"

	"c2f_admin_webapp/model/utils"

	"c2f_admin_webapp/model/db"
	"encoding/json"
	"errors"
	"fmt"

	"runtime/debug"

	"c2f_admin_webapp/model/util/log"

	"github.com/astaxie/beego"
	//"github.com/astaxie/beego/validation"

	//"c2f_admin_webapp/model/C2fDeliveryApi"

	"c2f_admin_webapp/model/encoding/json/O2OApi"
)

type PayuResponse struct {
	beego.Controller
}

type Payload struct {
	Txnid          string `json:"txnid"`
	Status         string `json:"status"`
	Mihpayid       string `json:"mihpayid"`
	Hash           string `json:"hash"`
	Mode           string `json:"mode"`
	ErrorMessage   string `json:"error_Message"`
	Error          string `json:"error"`
	Unmappedstatus string `json:"unmappedstatus"`
}

func (c *PayuResponse) Post() {

	var err error
	var r O2OApi.Response
	// var request []byte
	var response []byte
	// var l O2OApi.Request

	// Set CORS headers
	c.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	c.Ctx.Output.Header("Access-Control-Allow-Methods", "POST")
	c.Ctx.Output.Header("Access-Control-Allow-Headers", "Content-Type")

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "PayU Transaction Id  - ", c.GetString("txnid"))
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "PayU Transaction Status  - ", c.GetString("status"))
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "mihpayid - ", c.GetString("mihpayid"))
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "hash - ", c.GetString("hash"))

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "mode - ", c.GetString("mode"))
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "error_Message - ", c.GetString("error_Message"))
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "error - ", c.GetString("error"))
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "unmappedstatus - ", c.GetString("unmappedstatus"))

	PayUTransactionId := c.GetString("txnid")
	PayUTransactionStatus := c.GetString("status")
	PayUErrorMessage := c.GetString("error_Message")
	PayUMode := c.GetString("mode")

	PayUMihapyid := c.GetString("mihpayid")
	PayUHash := c.GetString("hash")
	PayUError := c.GetString("error")
	PayUnmappedstatus := c.GetString("unmappedstatus")
	var tpl string

	if PayUTransactionStatus == "success" {
		tpl = "paymentsuccess/paymentsuccess.html"

	} else if PayUTransactionStatus == "failure" {

		tpl = "paymentfailure/paymentfailure.html"

	}

	defer func() {
		if l_exception := recover(); l_exception != nil {
			stack := debug.Stack()
			log.Println(beego.AppConfig.String("loglevel"), "Exception", string(stack))
			r.Code = "99"
			r.Message = "System is unable to process your request"
			r.Status = "PENDING"
		}
		if err != nil {
			r.Code = "11"
			r.Message = err.Error()
			r.Status = "DECLINED"
		} else {

			c.TplName = tpl

		}
		if r.Code == "99" || r.Code == "11" {
			response, _ = json.Marshal(r)
			log.Println(beego.AppConfig.String("loglevel"), "Debug", "Response : "+string(response))

			// utils.EventLogs(c.Ctx, l, c.Ctx.Input.Method(), string(request), string(response), c.Input(), c.Data, err)
			// Prepare request string
			request := c.Ctx.Input.RequestBody
			requestStr := string(request)

			// Prepare output map
			outputMap := make(map[interface{}]interface{})
			outputMap["txnid"] = PayUTransactionId
			outputMap["status"] = PayUTransactionStatus

			utils.PayUEventLogs(
				c.Ctx,
				PayUTransactionId,    // sessionID
				c.GetString("phone"), // mobile
				c.Ctx.Input.Method(),
				requestStr,
				string(response),
				c.Input(),
				outputMap,
				err,
			)
			c.Ctx.Output.Body(response)

		}
		return
	}()

	utils.SetHTTPHeader(c.Ctx)

	if PayUTransactionStatus == "success" {

		fmt.Println("------------------------Transaction Success ---------------")

		result, err := db.Db.Exec(`UPDATE  CustomerOrderTxn SET payment_status=?,payment_remarks=? WHERE payu_txnid = ?`, "APPROVED", "Payment done by "+PayUMode, PayUTransactionId)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Customer Order  updation for online payment info failed ")
			return
		}
		if err != nil {
			err = errors.New("Customer Order  updation for online payment info failed")
			return
		}

		i, err := result.RowsAffected()
		if err != nil || i == 0 {
			err = errors.New("Customer Order  updation for online payment info failed")
			return
		}

	} else if PayUTransactionStatus == "failure" {

		fmt.Println("------------------------Transaction Failed ---------------")

		result, err := db.Db.Exec(`UPDATE  CustomerOrderTxn SET status = ?,payment_status=?,payment_remarks=? WHERE payu_txnid = ?`, "ORDER_REJECTED", "DECLINED", "CANCELLED - "+PayUErrorMessage, PayUTransactionId)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Customer Order  updation for online payment info failed ")
			return
		}
		if err != nil {
			err = errors.New("Customer Order  updation for online payment info failed")
			return
		}

		i, err := result.RowsAffected()
		if err != nil || i == 0 {
			err = errors.New("Customer Order  updation for online payment info failed")
			return
		}

	}

	// JSON FOR PAYLOAD PART
	payuresponse := Payload{
		Txnid:          PayUTransactionId,
		Status:         PayUTransactionStatus,
		Mihpayid:       PayUMihapyid,
		Hash:           PayUHash,
		Mode:           PayUMode,
		ErrorMessage:   PayUErrorMessage,
		Error:          PayUError,
		Unmappedstatus: PayUnmappedstatus,
	}

	var jsonData2 []byte
	jsonData2, err = json.Marshal(payuresponse)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Payload json Data - ", string(jsonData2))

	result2, err := db.Db.Exec(`UPDATE  Partner_Logs SET response = ? WHERE payu_txnid = ?`, string(jsonData2), PayUTransactionId)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Paylogs updation failed ")
		return
	}
	if err != nil {
		err = errors.New("Paylogs updation failed")
		return
	}

	i, err := result2.RowsAffected()
	if err != nil || i == 0 {
		err = errors.New("Paylogs updation failed")
		return
	}

	// if err := c.ParseForm(&l); err != nil {
	// 	err = errors.New("Invalid Request Received")
	// }
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Form Data - ", l)

	// request, err = json.Marshal(l)
	// if err != nil {
	// 	err = errors.New("Invalid Request Received")
	// }

	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Form Data - ", string(request))
	// valid := validation.Validation{}
	// b, err := valid.Valid(&l)
	// if err != nil {
	// 	err = errors.New("Parameter validation failed")
	// }

	// if !b {
	// 	for _, err := range valid.Errors {
	// 		log.Println(beego.AppConfig.String("loglevel"), "Debug", err.Key, ":", err.Message, ":", err.Field, ":", err.LimitValue, ":", err.Name, ":", err.Tmpl, ":", err.Value)
	// 	}
	// 	err = errors.New("Invalid Input values")
	// }

	// l.PIP = c.Ctx.Input.IP()
	// if l.Mobile == "" {
	// 	err = errors.New("Mobile Can not be Empty")
	// 	return
	// }

	// if len(l.Mobile) != 10 {
	// 	err = errors.New("Mobile length incorrect")
	// 	return
	// }

	// usrInfo, err := general.GetUserInfo(l.Mobile)
	// if err != nil {
	// 	err = errors.New("User not found")
	// 	return
	// }

	// sendotp, err := msg91.SendOTP(l.Mobile)
	// if err != nil {
	// 	err = errors.New("OTP Send failed")
	// 	return
	// }

	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "sendotp - ", sendotp)

	// r.Name = usrInfo[0][0] + usrInfo[0][1]
	// r.Mobile = usrInfo[0][2]
	// r.ClusterName = usrInfo[0][4]

	return
}
