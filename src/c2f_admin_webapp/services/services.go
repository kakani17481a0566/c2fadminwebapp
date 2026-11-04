package services

import (
	"c2f_admin_webapp/model/db"
	"errors"

	"fmt"
	"io/ioutil"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"

	"time"

	"c2f_admin_webapp/model/database/sql"

	"github.com/scorredoira/email"

	"crypto/rand"

	"c2f_admin_webapp/model/encoding/base64"

	"c2f_admin_webapp/model/util/pbkdf2"

	"github.com/astaxie/beego"

	"log"
)

type unencryptedAuth struct {
	smtp.Auth
}

func SendEmail(emilid, name, password, emailtemplate string) {
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "called - ")

	uname := beego.AppConfig.String("EMAIL_NOTIFY_USERNAME")
	pass := beego.AppConfig.String("EMAIL_NOTIFY_PASSWORD")
	url := beego.AppConfig.String("EMAIL_NOTIFY_URL")
	to := beego.AppConfig.String("EMAIL_NOTIFY_TIMEOUT")
	loginurl := beego.AppConfig.String("EMAIL_APPLICATION_LOGIN_URL")
	recipients := strings.Split(emilid, "||")

	tmpFile := emailtemplate

	buff, err := ioutil.ReadFile(tmpFile)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", "read file -", err)
		return
	}

	msg := string(buff)
	msg = strings.Replace(string(msg), "{{.Name}}", name, -1)
	msg = strings.Replace(string(msg), "{{.Email}}", emilid, -1)
	msg = strings.Replace(string(msg), "{{.Password}}", password, -1)
	msg = strings.Replace(string(msg), "{{.LoginURL}}", loginurl, -1)

	m := email.NewHTMLMessage("Email", msg)
	m.From = mail.Address{Name: "KRISHIVALS", Address: uname}
	m.To = recipients

	// send it
	//auth := smtp.PlainAuth("", uname, pass, url)

	config := beego.AppConfig.String("EMAIL_AUTH_CONFIG_MODE")
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "config", config)

	if config == "1" {
		auth := smtp.PlainAuth("", uname, pass, url)
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "auth")
		if err = email.Send(url+":"+to, auth, m); err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			return
		}

	} else if config == "2" {

		auth := unencryptedAuth{
			smtp.PlainAuth(
				"",
				uname,
				pass,
				url,
			),
		}
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "no tls auth")
		if err = email.Send(url+":"+to, auth, m); err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			return
		}
	} else {
		log.Println(beego.AppConfig.String("loglevel"), "Debug", "no auth")
		if err = email.Send(url+":"+to, nil, m); err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			return
		}
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Email sent successfully")
	return
}

func EncryptPassword(pass string) (out []byte, err error) {

	//commenting display of password in logs
	//log.Println(beego.AppConfig.String("loglevel"), "Debug", "inside encryption password rec: ", pass)
	b := make([]byte, 32)
	_, err = rand.Read(b)
	var pbkdf pbkdf2.Pbkdf2
	pbkdf.Itr = 32
	pbkdf.KeyLen = 32
	pbkdf.Plain = []byte(pass)
	pbkdf.Salt = b
	err = pbkdf.Encrypt()
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to encrypt password")
		return
	}
	var tmp []byte
	tmp = append(tmp, pbkdf.Salt...)
	tmp = append(tmp, pbkdf.Cipher...)

	out, err = base64.Encode(tmp)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to encrypt password")
		return
	}
	//	log.Printf("%s %s \n", "inside encrypt pass after encryption:", out)
	return
}

func (a unencryptedAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	s := *server
	s.TLS = true
	return a.Auth.Start(&s)
}

func GetClusteruser() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT id,cluster_name from Clusters where status="ACTIVE" `)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Cluster Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Cluster Name Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func GetClusters3(cname string) (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT id,cluster_name from Clusters where status="ACTIVE" AND cluster_name NOT LIKE ?`, cname)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Cluster Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Cluster Name Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func GetMappedClusterforUseratLogin(id string, usertype string) (cluster string, err error) {

	if usertype == "CLUSTERMANAGER" {
		row, err := db.Db.Query(`SELECT User2Cluster.user_id,User2Cluster.cluster_id,Clusters.id,Clusters.cluster_name  
	from User2Cluster
	left join
	Clusters on
	User2Cluster.cluster_id = Clusters.id
	where  User2Cluster.user_id=? and Clusters.status="ACTIVE"`, id)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster manager's cluster Get Fail")
			return cluster, err
		}
		defer sql.Close(row)
		_, data, err := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster manager's cluster Scan Fail")
			return cluster, err
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data of Cluster manager's cluster - ", data, "\nData len - ", len(data))
		if len(data) > 0 {
			cluster = data[0][3]
		}
		return cluster, err

	} else if usertype == "CLUSTEREXECUTIVE" {

		row, err := db.Db.Query(`SELECT User2Cluster.clusterexecutive_id,User2Cluster.cluster_id,Clusters.id,Clusters.cluster_name  
	from User2Cluster
	left join
	Clusters on
	User2Cluster.cluster_id = Clusters.id
	where  User2Cluster.clusterexecutive_id=? and Clusters.status="ACTIVE"`, id)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster executive's cluster Get Fail")
			return cluster, err
		}
		defer sql.Close(row)
		_, data, err := sql.Scan(row)
		if err != nil {
			log.Println(beego.AppConfig.String("loglevel"), "Error", err)
			err = errors.New("Cluster executive's cluster Scan Fail")
			return cluster, err
		}

		log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data of Cluster executive's cluster - ", data, "\nData len - ", len(data))
		if len(data) > 0 {
			cluster = data[0][3]
		}
		return cluster, err

	} else {

		cluster = ""

	}

	return

}

func GetClusterAsPerMapping(id string) (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT User2Cluster.user_id,User2Cluster.cluster_id,Clusters.id,Clusters.cluster_name  
	from User2Cluster
	left join
	Clusters on
	User2Cluster.cluster_id = Clusters.id
	where  User2Cluster.user_id=? and Clusters.status="ACTIVE"`, id)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("mapped Cluster Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("mapped Cluster Name Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data of Cluster as per mapped - ", data, "\nData len - ", len(data))
	return

}

func GetClusterAsPerMappingExec(id string) (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT User2Cluster.clusterexecutive_id,User2Cluster.cluster_id,Clusters.id,Clusters.cluster_name  
	from User2Cluster
	left join
	Clusters on
	User2Cluster.cluster_id = Clusters.id
	where  User2Cluster.clusterexecutive_id=? and Clusters.status="ACTIVE"`, id)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("mapped Cluster Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("mapped Cluster Name Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func GetClustermappedDeliveryuser(id string) (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT DeliveryUser.id,DeliveryUser.mob_no,DeliveryUser.first_name,DeliveryUser.email,Cluster2DeliveryUser.cluster_id,Cluster2DeliveryUser.deliveryuser_id,Clusters.id,Clusters.cluster_name
	from Cluster2DeliveryUser 
	left join DeliveryUser
	on DeliveryUser.id=Cluster2DeliveryUser.deliveryuser_id
	left join Clusters
	on Cluster2DeliveryUser.cluster_id=Clusters.id
	where DeliveryUser.status="ACTIVE" and Cluster2DeliveryUser.cluster_id=?`, id)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Delivery User Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Delivery User Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func GetClustermappedDeliveryuserForRessignOrder(id, clustername, deliveryusermob string) (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT DeliveryUser.id, DeliveryUser.mob_no, DeliveryUser.first_name, DeliveryUser.email, 
	Cluster2DeliveryUser.cluster_id, Cluster2DeliveryUser.deliveryuser_id, Clusters.id, Clusters.cluster_name
	FROM Cluster2DeliveryUser 
	LEFT JOIN DeliveryUser ON DeliveryUser.id = Cluster2DeliveryUser.deliveryuser_id
	LEFT JOIN Clusters ON Cluster2DeliveryUser.cluster_id = Clusters.id
	WHERE DeliveryUser.status = "ACTIVE" 
	AND Cluster2DeliveryUser.cluster_id = ? 
	AND Clusters.cluster_name = ? 
	AND DeliveryUser.mob_no != ?`, id, clustername, deliveryusermob)

	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Delivery User Name Get Fail")
		return
	}
	defer sql.Close(row)

	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Delivery User Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return
}

func GetClustermappedSupplyExecutive(id string) (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT SupplyExecutive.id,SupplyExecutive.mob_no,SupplyExecutive.first_name,SupplyExecutive.email,Cluster2SupplyExecutive.cluster_id,Cluster2SupplyExecutive.supplyexecutive_id,Clusters.id,Clusters.cluster_name
	from Cluster2SupplyExecutive 
	left join SupplyExecutive
	on SupplyExecutive.id=Cluster2SupplyExecutive.supplyexecutive_id
	left join Clusters
	on Cluster2SupplyExecutive.cluster_id=Clusters.id
	where SupplyExecutive.status="ACTIVE" and Cluster2SupplyExecutive.cluster_id=?`, id)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("SupplyExecutive Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("SupplyExecutive Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func GetDeliveryExecutive() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT id,mob_no,first_name,email from DeliveryUser where status="ACTIVE"`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Delivery User Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Delivery User Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func GetNonmappedDeliveryExecutive() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT DeliveryUser.id,DeliveryUser.mob_no,DeliveryUser.first_name,DeliveryUser.email
	 from DeliveryUser
	left join Cluster2DeliveryUser on DeliveryUser.id=Cluster2DeliveryUser.deliveryuser_id
	 where Cluster2DeliveryUser.deliveryuser_id is null and DeliveryUser.status="ACTIVE"`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Non-Mapped Delivery Users Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Non-mapped Delivery Users Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func GetDeliveryuserid2(mob string) (id, cluster string, err error) {

	row, err := db.Db.Query(`select id from DeliveryUser where mob_no = ? `, mob)
	if err != nil {
		beego.Error(err)
		err = errors.New("Unable to get the DeliveryUser Id  ")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		beego.Error(err)
		err = errors.New("Message Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the DeliveryUser ID info")
		return
	}

	row1, err1 := db.Db.Query(`SELECT DeliveryUser.id,DeliveryUser.mob_no,DeliveryUser.first_name,DeliveryUser.email,Cluster2DeliveryUser.cluster_id,Cluster2DeliveryUser.deliveryuser_id,Clusters.id,Clusters.cluster_name
	from Cluster2DeliveryUser 
	left join DeliveryUser
	on DeliveryUser.id=Cluster2DeliveryUser.deliveryuser_id
	left join Clusters
	on Cluster2DeliveryUser.cluster_id=Clusters.id
	where Cluster2DeliveryUser.deliveryuser_id=?`, data[0][0])
	if err1 != nil {
		beego.Error(err)
		err = errors.New("Unable to get the Collectionagent Name  ")
		return
	}
	defer sql.Close(row1)
	_, data1, err2 := sql.Scan(row1)
	if err2 != nil {
		beego.Error(err)
		err2 = errors.New("Message Scan Fail")
		return
	}
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data1) <= 0 {
		err = errors.New("Unable to get the Collectionagent  info")
		return
	}
	id = data[0][0]
	cluster = data1[0][7]
	return
}

func GetClustermappedagent(id string) (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT Collection_Agent.id,Collection_Agent.email,Collection_Agent.first_name,Collection_Agent.mob_no,Cluster2CollectionAgent.cluster_id,Cluster2CollectionAgent.agent_id,Clusters.id,Clusters.cluster_name
	from Cluster2CollectionAgent
	left join Collection_Agent on Collection_Agent.id=Cluster2CollectionAgent.agent_id
	left join Clusters on Clusters.id=Cluster2CollectionAgent.cluster_id
	where Collection_Agent.status="ACTIVE" and Cluster2CollectionAgent.cluster_id=?`, id)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Agent Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Agent  Name Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return
}

func GetClustermappedcollectionagent() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT Collection_Agent.id,Collection_Agent.email,Collection_Agent.first_name,Collection_Agent.mob_no,Cluster2CollectionAgent.cluster_id,Cluster2CollectionAgent.agent_id,Clusters.id,Clusters.cluster_name
	from Cluster2CollectionAgent
	left join Collection_Agent on Collection_Agent.id=Cluster2CollectionAgent.agent_id
	left join Clusters on Clusters.id=Cluster2CollectionAgent.cluster_id
	where Collection_Agent.status="ACTIVE"`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Agent Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Agent  Name Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return
}

func GetProducers2(id string) (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT DISTINCT Customer.id,Customer.mob_no,Customer.first_name,Product2Producer2Cluster.producer_id,Product2Producer2Cluster.cluster_id
	from Product2Producer2Cluster
	left join
	Customer on Product2Producer2Cluster.producer_id=Customer.id
	where Customer.status="ACTIVE" AND Product2Producer2Cluster.cluster_id=? AND (Customer.user_type="PRODUCER" OR Customer.user_type="PRODUCER+CONSUMER")`, id)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Producer Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Producer  Name Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

// func GetClusterExecutive() (data [][]string, err error) {
// 	row, err := db.Db.Query(`SELECT id,clusterexecutive_id from User2Cluster where status="ACTIVE" `)
// 	if err != nil {
// 		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
// 		err = errors.New("Cluster Name Get Fail")
// 		return
// 	}
// 	defer sql.Close(row)
// 	_, data, err = sql.Scan(row)
// 	if err != nil {
// 		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
// 		err = errors.New("Cluster Name Scan Fail")
// 		return
// 	}

// 	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
// 	return

// }

func GetClustercollectionagent() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT id,cluster_name from Clusters where status="ACTIVE"`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Cluster Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Cluster Name Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func GetCategory() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT id,name from Category where status="ACTIVE"`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Category Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Category  Name Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func GetSubCategory() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT id,name from SubCategory where status="ACTIVE"`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("SubCategory Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("SubCategory  Name Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func GetCommodityName() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT id,commodity_name from GST_Rates`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("commodity_name  Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("commodity_name  Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func Getproducts() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT id,product_name,product_qty,product_package_qty,product_sku from Lookup_products where product_status="ACTIVE"`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Product Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Product  Name Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func GetProdctQtys() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT id,product_qty from Lookup_products where product_status="ACTIVE"`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Product Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Product  Name Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func Getusers() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT id,email,first_name from User where status="ACTIVE" AND role="4"`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("User Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("User  Name Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func GetNonmappedClusterManager() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT User.id,User.email,User.first_name 
	from User
	left join User2Cluster on User.id=User2Cluster.user_id
	 where User2Cluster.user_id is null and User.status="ACTIVE" AND User.role="4"`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Non-mapped ClusterManager Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Non-mapped ClusterManager Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}
func Getusers1() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT id,email,first_name from User where status="ACTIVE" AND role="5"`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("User Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("User  Name Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func GetNonmappedClusterExecutive() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT User.id,User.email,User.first_name 
	from User
	left join User2Cluster on User.id=User2Cluster.clusterexecutive_id
	 where User2Cluster.clusterexecutive_id is null and User.status="ACTIVE" AND User.role="5"`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Non-mapped ClusterExecutive Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Non-mapped ClusterExecutive Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func GetNonmappedSupplyExecutive() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT SupplyExecutive.id,SupplyExecutive.email,SupplyExecutive.first_name 
	from SupplyExecutive
	left join Cluster2SupplyExecutive on SupplyExecutive.id=Cluster2SupplyExecutive.supplyexecutive_id
	 where Cluster2SupplyExecutive.supplyexecutive_id is null and SupplyExecutive.status="ACTIVE"`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Non-mapped SupplyExecutive Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Non-mapped SupplyExecutive Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func Getcollectionagent() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT id,email ,first_name,mob_no from Collection_Agent where status="ACTIVE"`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Agent Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Agent  Name Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return
}

func GetNonmappedcollectionagent() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT Collection_Agent.id,Collection_Agent.email,Collection_Agent.first_name,Collection_Agent.mob_no 
	from Collection_Agent 
	left join Cluster2CollectionAgent on Collection_Agent.id=Cluster2CollectionAgent.agent_id
	where Cluster2CollectionAgent.agent_id is null and Collection_Agent.status="ACTIVE"`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Non-mapped Agent Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Non-mapped Agent Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return
}

func GetRoles() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT id,role_name from Lookup_roles where role_status="ACTIVE"`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Role Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Role  Name Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func GetProducersByClusters(id string) (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT Customer.id,Customer.mob_no,Customer.first_name,Product2Producer2Cluster.cluster_id,Product2Producer2Cluster.producer_id from Product2Producer2Cluster left join Customer on Product2Producer2Cluster.producer_id= Customer.id where Product2Producer2Cluster.cluster_id=? AND Customer.status="ACTIVE" AND Customer.user_type="PRODUCER" OR Customer.user_type="PRODUCER+CONSUMER" `, id)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Producer Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Producer  Name Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func GetProducers() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT id,mob_no,first_name from Customer where status="ACTIVE" AND user_type="PRODUCER" OR user_type="PRODUCER+CONSUMER" `)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Producer Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Producer  Name Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func GetProducers1() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT producer_mob,status,product_name,cluster_name,product_qty,product_package_qty ,supplyordertxn_id ,create_date,product_supply_units FROM SupplyOrderTxn where status="DEMAND_ORDER_CREATED" ORDER BY create_date desc`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Producer fail with status DEMAND_ORDER_CREATED")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Producer fail with status DEMAND_ORDER_CREATED")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func GetNotificationDuration(data [][]string) (data1 []string, err error) {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}

	currentTime := time.Now().In(loc)
	currentTime1 := currentTime.Format("2006-01-02 15:04:05")

	t2, err1 := time.ParseInLocation("2006-01-02 15:04:05", currentTime1, loc)
	if err1 != nil {
		fmt.Println(err1)
		return
	}

	log.Println("Debug", "Current Time :t2 - ", t2)

	for i := range data {
		t1, err := time.ParseInLocation("2006-01-02 15:04:05", data[i][7], loc)
		if err != nil {
			fmt.Println(err)
			continue
		}

		diff := t2.Sub(t1)
		diff1 := int(diff.Minutes())
		var time1 string

		if diff.Minutes() < 60 {
			time := int(diff.Minutes())
			time1 = strconv.Itoa(time) + " mins ago"
		} else if 60 <= diff1 && diff1 < 1440 {
			time := int(diff.Minutes() / 60)
			time1 = strconv.Itoa(time) + " hrs ago"
		} else if 1440 <= diff1 && diff1 < 10080 {
			time := int(diff.Minutes() / 1440)
			time1 = strconv.Itoa(time)
			if time1 == "1" {
				time1 = time1 + " day ago"
			} else {
				time1 = time1 + " days ago"
			}
		} else {
			time := int(diff.Minutes() / 10080)
			time1 = strconv.Itoa(time)
			if time1 == "1" {
				time1 = time1 + " week ago"
			} else {
				time1 = time1 + " weeks ago"
			}
		}

		data1 = append(data1, time1)
	}

	return
}

func GetProducersinfo() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT Product2Producer2Cluster.product_id,Product2Producer2Cluster.producer_id, 
       Lookup_products.product_name,
       Customer.mob_no,
    Customer.first_name,
    Product2Producer2Cluster.status,
    Lookup_products.product_package_qty,
    Lookup_products.product_qty,
     Customer.id
    FROM Product2Producer2Cluster
LEFT JOIN Lookup_products 
ON Product2Producer2Cluster.product_id = Lookup_products.id
LEFT JOIN Customer 
ON Product2Producer2Cluster.producer_id = Customer.id WHERE Product2Producer2Cluster.status="ACTIVE" AND (Customer.user_type="PRODUCER" OR  Customer.user_type="PRODUCER+CONSUMER")`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Producer Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Producer  Name Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func GetClusterMappedProducers(id string) (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT Product2Producer2Cluster.product_id,Product2Producer2Cluster.producer_id, 
       Lookup_products.product_name,
       Customer.mob_no,
    Customer.first_name,
    Product2Producer2Cluster.status,
    Lookup_products.product_package_qty,
    Lookup_products.product_qty,
     Customer.id
    FROM Product2Producer2Cluster
LEFT JOIN Lookup_products 
ON Product2Producer2Cluster.product_id = Lookup_products.id
LEFT JOIN Customer 
ON Product2Producer2Cluster.producer_id = Customer.id WHERE Product2Producer2Cluster.cluster_id=? AND Product2Producer2Cluster.status="ACTIVE" AND (Customer.user_type="PRODUCER" OR  Customer.user_type="PRODUCER+CONSUMER") `, id)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Producer Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Producer  Name Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func GetDeliveryuser() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT DeliveryUser.id,DeliveryUser.mob_no,DeliveryUser.first_name,DeliveryUser.email,Cluster2DeliveryUser.cluster_id,Cluster2DeliveryUser.deliveryuser_id,Clusters.id,Clusters.cluster_name
	from Cluster2DeliveryUser 
	left join DeliveryUser
	on DeliveryUser.id=Cluster2DeliveryUser.deliveryuser_id
	left join Clusters
	on Cluster2DeliveryUser.cluster_id=Clusters.id
	where DeliveryUser.status="ACTIVE"`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Delivery User Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Delivery User Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data for Delivery Users - ", data, "\nData len - ", len(data))
	return

}

func GetDeliveryuserForRessignOrder(clustername, deliveryusermob string) (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT DeliveryUser.id,DeliveryUser.mob_no,DeliveryUser.first_name,DeliveryUser.email,Cluster2DeliveryUser.cluster_id,Cluster2DeliveryUser.deliveryuser_id,Clusters.id,Clusters.cluster_name
	from Cluster2DeliveryUser 
	left join DeliveryUser
	on DeliveryUser.id=Cluster2DeliveryUser.deliveryuser_id
	left join Clusters
	on Cluster2DeliveryUser.cluster_id=Clusters.id
	where DeliveryUser.status="ACTIVE"
	AND Clusters.cluster_name = ? 
	AND DeliveryUser.mob_no != ?`, clustername, deliveryusermob)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Delivery User Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Delivery User Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data for Delivery Users - ", data, "\nData len - ", len(data))
	return

}

func GetConsumer() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT id,mob_no,first_name from Customer where status="ACTIVE" AND user_type="CONSUMER" OR user_type="PRODUCER+CONSUMER"`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Producer Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Producer  Name Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func GetConsumerid(mob string) (id string, err error) {

	row, err := db.Db.Query(`select id from Customer where mob_no = ? `, mob)
	if err != nil {
		beego.Error(err)
		err = errors.New("Unable to get the Customer Id  ")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		beego.Error(err)
		err = errors.New("Message Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the Consumer ID info")
		return
	}
	id = data[0][0]
	return
}

func GetProducerid(mob string) (id string, err error) {

	row, err := db.Db.Query(`select id from Customer where mob_no = ? `, mob)
	if err != nil {
		beego.Error(err)
		err = errors.New("Unable to get the Producer Id  ")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		beego.Error(err)
		err = errors.New("Message Scan Fail")
		return
	}
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the Producer ID info")
		return
	}
	id = data[0][0]
	return
}

func GetCollectionAgentId(mob string) (id string, err error) {

	row, err := db.Db.Query(`select id from Collection_Agent where mob_no = ? `, mob)
	if err != nil {
		beego.Error(err)
		err = errors.New("Unable to get the Customer Id  ")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		beego.Error(err)
		err = errors.New("Message Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the Consumer ID info")
		return
	}
	id = data[0][0]
	return
}

func GetCollectionAgentName(mob string) (firstname, cname, id string, err error) {

	row, err := db.Db.Query(`select id,first_name from Collection_Agent where mob_no = ? `, mob)
	if err != nil {
		beego.Error(err)
		err = errors.New("Unable to get the Collectionagent Name  ")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		beego.Error(err)
		err = errors.New("Message Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the Collectionagent  info")
		return
	}

	row1, err1 := db.Db.Query(`select Cluster2CollectionAgent.agent_id,Cluster2CollectionAgent.cluster_id,Clusters.id,Clusters.cluster_name
	 from Cluster2CollectionAgent
	left join 
	Clusters
	on Cluster2CollectionAgent.cluster_id=Clusters.id
	where Cluster2CollectionAgent.agent_id = ? `, data[0][0])
	if err1 != nil {
		beego.Error(err)
		err = errors.New("Unable to get the Collectionagent Name  ")
		return
	}
	defer sql.Close(row1)
	_, data1, err2 := sql.Scan(row1)
	if err2 != nil {
		beego.Error(err)
		err2 = errors.New("Message Scan Fail")
		return
	}
	// log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data1) <= 0 {
		err = errors.New("Unable to get the Collectionagent  info")
		return
	}
	firstname = data[0][1]
	id = data[0][0]
	cname = data1[0][3]

	return
}

func GetCustomerName(mob string) (id string, err error) {

	row, err := db.Db.Query(`select first_name from Customer where mob_no = ? `, mob)
	if err != nil {
		beego.Error(err)
		err = errors.New("Unable to get the Customer Name  ")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		beego.Error(err)
		err = errors.New("Message Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the Customer  info")
		return
	}
	id = data[0][0]
	return
}

func GetOderid(oderid string) (id, cid, oid, odate string, err error) {

	row, err := db.Db.Query(`select id,customer_id,order_id ,create_date from CustomerOrderTxn where order_id = ? `, oderid)
	if err != nil {
		beego.Error(err)
		err = errors.New("Unable to get the CustomerOrderTxn Id  ")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		beego.Error(err)
		err = errors.New("Message Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the Order info")
		return
	}
	id = data[0][0]
	cid = data[0][1]
	oid = data[0][2]
	odate = data[0][3]

	return
}

func GetCustomerinfo(cid string) (fname, add1, add2, city, zip, email, mob, latlong string, err error) {

	row, err := db.Db.Query(`select first_name,address1,address2,city,zipcode,email,mob_no,location_info from Customer where id = ? `, cid)
	if err != nil {
		beego.Error(err)
		err = errors.New("Unable to get the CustomerOrderTxn Id  ")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		beego.Error(err)
		err = errors.New("Message Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the Customer info")
		return
	}
	fname = data[0][0]
	add1 = data[0][1]
	add2 = data[0][2]
	city = data[0][3]
	zip = data[0][4]
	email = data[0][5]
	mob = data[0][6]
	latlong = data[0][7]
	//	locationtype = data[0][8]
	return
}

func Getinventory(sku, clr string) (dunits, aunits, runits, psku, cluster, damageunits string, err error) {

	row, err := db.Db.Query(`select inventory_delivered,inventory_available,inventory_received,product_sku,cluster_name,damaged_units from Inventory where product_sku = ? AND cluster_name=? `, sku, clr)
	if err != nil {
		beego.Error(err)
		err = errors.New("Unable to get the inventroy info  ")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		beego.Error(err)
		err = errors.New("Message Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the inventroy info")
		return
	}
	dunits = data[0][0]
	aunits = data[0][1]
	runits = data[0][2]
	psku = data[0][3]
	cluster = data[0][4]
	damageunits = data[0][5]

	return
}

func GetConsumer1(cname string) (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT id,mob_no,first_name,cluster_name,address1,address2, city,state,zipcode from Customer where status="ACTIVE" AND (user_type="CONSUMER" OR user_type="PRODUCER+CONSUMER") AND cluster_name=?`, cname)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Customer Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Customer  Name Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func GetDashboardClusterDetails() (data01 [][]string, err error) {

	row, err := db.Db.Query(`select cluster_name, product_name, product_sku, product_qty, inventory_received, inventory_delivered, inventory_available, product_package_type from Inventory  where status="SUPPLY_COMPLETED"`)
	if err != nil {
		beego.Error(err)
		err = errors.New("Unable to get the Cluster Id  ")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		beego.Error(err)
		err = errors.New("Message Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the cluster ID info")
		return
	}
	//id = data[0][0]
	return
}

func Getinventory2(sku, clr string) (aunits string, err error) {

	row, err := db.Db.Query(`select inventory_available,product_sku,cluster_name from Inventory where product_sku = ? AND cluster_name=? `, sku, clr)
	if err != nil {
		beego.Error(err)
		err = errors.New("Unable to get the inventroy info  ")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		beego.Error(err)
		err = errors.New("Message Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the inventroy info")
		return
	}
	//dunits = data[0][0]
	aunits = data[0][0]
	//runits = data[0][2]
	//psku = data[0][1]
	//cluster = data[0][2]

	return
}

func GetInventory3(cname, psku string) (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT inventory_available,product_sku,cluster_name from Inventory where cluster_name NOT LIKE ? AND product_sku LIKE ? `, cname, psku)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Cluster Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Cluster Name Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func GetInventory5() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT inventory_available,product_sku,cluster_name from Inventory`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Cluster Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Cluster Name Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func GetOderinfo(oderid string) (oid, odate, otype, orepeat, dfrom, dto, dtime, cname, subcname, pname, pqty, productpackageqty, pprice, punits, tprice, totalsellprice, clustername, productsku string, err error) {
	row, err := db.Db.Query(`select order_id,
	create_date,
	order_purchagetype,
	order_repeat,
	orderdeliver_fromdate,
	orderdeliver_todate,
	orderdeliver_time,
	category_name,
	subcategory_name,
	product_name,
	product_qty,
	product_package_qty,
	product_price,
	order_units,
	total_product_price,
	total_selling_price,
	cluster_name,
	product_sku
	 from CustomerOrderTxn where order_id = ? `, oderid)
	if err != nil {
		beego.Error(err)
		err = errors.New("Unable to get the OrderTxn Id  ")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		beego.Error(err)
		err = errors.New("Message Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the order number info")
		return
	}
	oid = data[0][0]
	odate = data[0][1]
	otype = data[0][2]
	orepeat = data[0][3]
	dfrom = data[0][4]
	dto = data[0][5]
	dtime = data[0][6]
	cname = data[0][7]
	subcname = data[0][8]
	pname = data[0][9]
	pqty = data[0][10]
	productpackageqty = data[0][11]
	pprice = data[0][12]
	punits = data[0][13]
	tprice = data[0][14]
	totalsellprice = data[0][15]
	clustername = data[0][16]
	productsku = data[0][17]

	return
}

func GetClusterid(pid string) (id string, err error) {

	row, err := db.Db.Query(`select cluster_id from Product2Cluster where product_id = ? `, pid)
	if err != nil {
		beego.Error(err)
		err = errors.New("Unable to get the Cluster Id  ")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		beego.Error(err)
		err = errors.New("Message Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the cluster ID info")
		return
	}
	id = data[0][0]
	return
}

func GetClusterName(cid string) (name string, err error) {

	row, err := db.Db.Query(`select cluster_name from Clusters where id = ? `, cid)
	if err != nil {
		beego.Error(err)
		err = errors.New("Unable to get the Cluster Name  ")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		beego.Error(err)
		err = errors.New("Message Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the Cluster Name info")
		return
	}
	name = data[0][0]
	return
}

func GetProductCount(status string) (id string, err error) {

	row, err := db.Db.Query(`SELECT count(*)as acount from Lookup_products where product_status = ? `, status)
	if err != nil {
		beego.Error(err)
		err = errors.New("Unable to get the product count  ")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		beego.Error(err)
		err = errors.New("Message Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the product count number info")
		return
	}
	id = data[0][0]
	return
}

func GetClusterCount(status string) (id string, err error) {

	row, err := db.Db.Query(`SELECT count(*)as acount from Clusters where status = ? `, status)
	if err != nil {
		beego.Error(err)
		err = errors.New("Unable to get the clusters count  ")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		beego.Error(err)
		err = errors.New("Message Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the product count number info")
		return
	}
	id = data[0][0]
	return
}

func GetCustomerOderCount(status string) (id string, err error) {

	row, err := db.Db.Query(`SELECT count(*)as acount from CustomerOrderTxn`)
	if err != nil {
		beego.Error(err)
		err = errors.New("Unable to get the Customers Order count  ")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		beego.Error(err)
		err = errors.New("Message Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the Customers Order number info")
		return
	}
	id = data[0][0]
	return
}

func GetSupplyOderCount(status string) (id string, err error) {

	row, err := db.Db.Query(`SELECT count(*)as acount from SupplyOrderTxn`)
	if err != nil {
		beego.Error(err)
		err = errors.New("Unable to get the SupplyOrderTxn count  ")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		beego.Error(err)
		err = errors.New("Message Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the SupplyOrderTxn count number info")
		return
	}
	id = data[0][0]
	return
}

func GetClustername(id string) (name string, err error) {
	row, err := db.Db.Query(`SELECT id,cluster_name from Clusters where id=?`, id)
	if err != nil {
		beego.Error(err)
		err = errors.New("Unable to get the cluster_name   ")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		beego.Error(err)
		err = errors.New("Message Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the cluster_name  info")
		return
	}
	name = data[0][1]
	return

}

func GetCluster(cname string) (id string, err error) {
	row, err := db.Db.Query(`SELECT id from Clusters where cluster_name=?`, cname)
	if err != nil {
		beego.Error(err)
		err = errors.New("Unable to get the cluster_name   ")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		beego.Error(err)
		err = errors.New("Message Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the cluster_name  info")
		return
	}
	id = data[0][0]
	return

}

func GetProductId(pname string) (id, sku string, err error) {
	row, err := db.Db.Query(`SELECT id,product_sku from Lookup_products where product_name=?`, pname)
	if err != nil {
		beego.Error(err)
		err = errors.New("Unable to get the cluster_name   ")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		beego.Error(err)
		err = errors.New("Message Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the cluster_name  info")
		return
	}
	id = data[0][0]
	sku = data[0][1]
	return

}

func GetCid(id string) (name string, err error) {
	row, err := db.Db.Query(`SELECT id,user_id,cluster_id from User2Cluster where user_id=?`, id)
	if err != nil {
		beego.Error(err)
		err = errors.New("Unable to get the cluster_id count  ")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		beego.Error(err)
		err = errors.New("Message Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the cluster_id count number info")
		return
	}
	name = data[0][2]
	return

}

func GetCid1(id string) (name string, err error) {
	row, err := db.Db.Query(`SELECT id,clusterexecutive_id,cluster_id  from User2Cluster where clusterexecutive_id=?`, id)
	if err != nil {
		beego.Error(err)
		err = errors.New("Unable to get the cluster_id count  ")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		beego.Error(err)
		err = errors.New("Message Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the cluster_id count number info")
		return
	}
	name = data[0][2]
	return

}

func GetClustername1(sku, clu string) (cluster, psku string, err error) {
	row, err := db.Db.Query(`SELECT id,cluster_name,product_sku from Clusters where cluster_name=? AND product_sku=?`, clu, sku)
	if err != nil {
		beego.Error(err)
		err = errors.New("Unable to get the cluster_name   ")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		beego.Error(err)
		err = errors.New("Message Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the cluster_name  info")
		return
	}
	cluster = data[0][1]
	psku = data[0][2]
	return

}

func Getinventory1() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT id ,product_sku from Inventory`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Category Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Category  Name Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func GetDeliveryuserid(mob string) (id string, err error) {

	row, err := db.Db.Query(`select id from DeliveryUser where mob_no = ? `, mob)
	if err != nil {
		beego.Error(err)
		err = errors.New("Unable to get the DeliveryUser Id  ")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		beego.Error(err)
		err = errors.New("Message Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the DeliveryUser ID info")
		return
	}
	id = data[0][0]
	return
}

func GetClusterid1(duid string) (id string, err error) {

	row, err := db.Db.Query(`select cluster_id from Cluster2DeliveryUser where deliveryuser_id = ? `, duid)
	if err != nil {
		beego.Error(err)
		err = errors.New("Unable to get the Cluster Id  ")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		beego.Error(err)
		err = errors.New("Message Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the Cluster ID info")
		return
	}
	id = data[0][0]
	return
}

func GetClusterLocation(id string) (data [][]string, err error) {
	row, err := db.Db.Query(`select location,cluster_name from Clusters where id=?`, id)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get the Cluster location")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Cluster Location Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the Cluster location")
		return
	}
	return
}

func GetOrderLocation(mobile string) (data [][]string, err error) {
	row, err := db.Db.Query(`select 
	JSON_UNQUOTE(JSON_EXTRACT(customer_details, '$.latlong')),JSON_UNQUOTE(JSON_EXTRACT(customer_details, '$.name')),JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.productPackageQty')),JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.productQty')),JSON_UNQUOTE(JSON_EXTRACT(oder_deatils, '$.productName'))
	from Oder2Delivery where deliveryuser_mob=? and status="DELIVERY_ASSIGNED" `, mobile)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get the order locations")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Order Location Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the order location")
		return
	}
	return
}

func Getinventoryavalibeunits(sku, clr string) (aunits, csprice, receivedunits string, err error) {

	row, err := db.Db.Query(`select inventory_available,cluster_selling_price,inventory_received from Inventory where product_sku = ? AND cluster_name=? `, sku, clr)
	if err != nil {
		beego.Error(err)
		err = errors.New("Unable to get the inventroy info  ")
		return
	}
	defer sql.Close(row)
	_, data, err := sql.Scan(row)
	if err != nil {
		beego.Error(err)
		err = errors.New("Message Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the inventroy info")
		return
	}

	aunits = data[0][0]
	csprice = data[0][1]
	receivedunits = data[0][2]

	return
}

func GetAllCustomerTokens() (data [][]string, err error) {
	row, err := db.Db.Query(`select fcm_token
	from Customer where user_type="CONSUMER" or user_type="PRODUCER+CONSUMER" `)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Unable to get the Tokens")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Tokens Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, " Data len - ", len(data))
	if len(data) <= 0 {
		err = errors.New("Unable to get the Tokens")
		return
	}
	return
}

func GetDeliveryChargePlaceOrder(productpackageqty, units, days string) (deliverycharge string, err error) {

	std_delivery_charge := beego.AppConfig.String("STD_DELIVERY_CHARGE")
	std_delivery_charge1, err := strconv.ParseFloat(std_delivery_charge, 64)

	additional_delivery_charge := beego.AppConfig.String("ADDITIONAL_DELIVERY_CHARGE")
	additional_delivery_charge1, err := strconv.ParseFloat(additional_delivery_charge, 64)

	fmt.Println(" Product Package Quantity is gms : ", productpackageqty)
	fmt.Println(" No of Units : ", units)
	fmt.Println(" No of Days : ", days)

	productpackageqty1, err := strconv.ParseFloat(productpackageqty, 64)
	units1, err := strconv.ParseFloat(units, 64)
	days1, err := strconv.ParseFloat(days, 64)

	totalgms := units1 * productpackageqty1 * days1

	fmt.Println(" Total ordered gms of this order : ", units1, "*", productpackageqty1, "*", days1, "=", totalgms)
	additional := totalgms - 500
	additionalgms := additional / 500
	fmt.Println("Addtional 500gms repeated is", additionalgms, "times")

	if additionalgms > 0 {

		TotalAdditionalCharge := additionalgms * additional_delivery_charge1
		fmt.Println("Total Additional Charge : 16 *", additionalgms, "=", TotalAdditionalCharge)

		TotalDeliveryCharge := std_delivery_charge1 + TotalAdditionalCharge
		fmt.Println("Total Delivery Charge : 36 + ", "16 * ", additionalgms, "=", TotalDeliveryCharge)

		TotalDeliveryCharge1 := strconv.FormatFloat(TotalDeliveryCharge, 'f', 0, 32)

		deliverycharge = TotalDeliveryCharge1

	} else {

		TotalDeliveryCharge := std_delivery_charge1
		fmt.Println("Total Delivery Charge :", TotalDeliveryCharge)

		TotalDeliveryCharge1 := strconv.FormatFloat(TotalDeliveryCharge, 'f', 0, 32)

		deliverycharge = TotalDeliveryCharge1

	}

	return
}

func GetDeliveryCharge(productpackageqty, units string) (deliverycharge string, err error) {

	std_delivery_charge := beego.AppConfig.String("STD_DELIVERY_CHARGE")
	std_delivery_charge1, err := strconv.ParseFloat(std_delivery_charge, 64)

	additional_delivery_charge := beego.AppConfig.String("ADDITIONAL_DELIVERY_CHARGE")
	additional_delivery_charge1, err := strconv.ParseFloat(additional_delivery_charge, 64)

	fmt.Println(" Product Package Quantity is gms : ", productpackageqty)
	fmt.Println(" No of Units : ", units)
	//fmt.Println(" No of Days : ", days)

	productpackageqty1, err := strconv.ParseFloat(productpackageqty, 64)
	units1, err := strconv.ParseFloat(units, 64)
	//days1, err := strconv.ParseFloat(days, 64)

	totalgms := units1 * productpackageqty1

	fmt.Println(" Total ordered gms of this order : ", units1, "*", productpackageqty1, "=", totalgms)
	additional := totalgms - 500
	additionalgms := additional / 500
	fmt.Println("Addtional 500gms repeated is", additionalgms, "times")

	if additionalgms > 0 {

		TotalAdditionalCharge := additionalgms * additional_delivery_charge1
		fmt.Println("Total Additional Charge : 16 *", additionalgms, "=", TotalAdditionalCharge)

		TotalDeliveryCharge := std_delivery_charge1 + TotalAdditionalCharge
		fmt.Println("Total Delivery Charge : 36 + ", "16 * ", additionalgms, "=", TotalDeliveryCharge)

		TotalDeliveryCharge1 := strconv.FormatFloat(TotalDeliveryCharge, 'f', 0, 32)

		deliverycharge = TotalDeliveryCharge1

	} else {

		TotalDeliveryCharge := std_delivery_charge1
		fmt.Println("Total Delivery Charge :", TotalDeliveryCharge)

		TotalDeliveryCharge1 := strconv.FormatFloat(TotalDeliveryCharge, 'f', 0, 32)

		deliverycharge = TotalDeliveryCharge1

	}

	return
}

func GetSupplyExecutive() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT SupplyExecutive.id,SupplyExecutive.mob_no,SupplyExecutive.first_name,SupplyExecutive.email,Cluster2SupplyExecutive.cluster_id,Cluster2SupplyExecutive.supplyexecutive_id,Clusters.id,Clusters.cluster_name
	from Cluster2SupplyExecutive 
	left join SupplyExecutive
	on SupplyExecutive.id=Cluster2SupplyExecutive.supplyexecutive_id
	left join Clusters
	on Cluster2SupplyExecutive.cluster_id=Clusters.id
	where SupplyExecutive.status="ACTIVE"`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("SupplyExecutive Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("SupplyExecutive Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data for SupplyExecutive  - ", data, "\nData len - ", len(data))
	return

}

func GetChannelPartnerType() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT id,name from ChannelPartnerType where status="ACTIVE" `)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Cluster Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Cluster Name Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func GetNonmappedClusterChannelPartner() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT MIN(id) AS id, first_name, MIN(email) AS email, MIN(mob_no) AS mob_no
	FROM ChannelPartner WHERE status = 'ACTIVE' GROUP BY first_name`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Non-mapped Cluster ChannelPartner Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Non-mapped Cluster ChannelPartner Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return
}
func GetNonmappedProductChannelPartner() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT MIN(id) AS id, first_name, MIN(email) AS email, MIN(mob_no) AS mob_no
	FROM ChannelPartner WHERE status = 'ACTIVE' GROUP BY first_name`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Non-mapped Product ChannelPartner Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Non-mapped Product ChannelPartner Scan Fail")
		return
	}
	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return
}

func GetProduct() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT id,product_name from Lookup_products where product_status="ACTIVE" `)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Cluster Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Cluster Name Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func GetProductAsPerMapping(id string) (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT User2Cluster.user_id,User2Cluster.cluster_id,Clusters.id,Clusters.cluster_name,Lookup_products.product_name
	from User2Cluster
	left join
	Clusters on
	User2Cluster.cluster_id = Clusters.id
	LEFT JOIN Product2Cluster ON User2Cluster.cluster_id = Product2Cluster.cluster_id
	LEFT JOIN Lookup_products ON Product2Cluster.product_id = Lookup_products.id
	where  User2Cluster.user_id=? and Product2Cluster.status="ACTIVE"`, id)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("mapped Cluster Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("mapped Cluster Name Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data of Cluster as per mapped - ", data, "\nData len - ", len(data))
	return

}

func GetProductAsPerMappingExec(id string) (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT User2Cluster.clusterexecutive_id,User2Cluster.cluster_id,Clusters.id,Clusters.cluster_name,Lookup_products.product_name
	from User2Cluster
	left join
	Clusters on
	User2Cluster.cluster_id = Clusters.id
	LEFT JOIN Product2Cluster ON User2Cluster.cluster_id = Product2Cluster.cluster_id
	LEFT JOIN Lookup_products ON Product2Cluster.product_id = Lookup_products.id
	where  User2Cluster.clusterexecutive_id=? and Product2Cluster.status="ACTIVE"`, id)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("mapped Cluster Name Get Fail")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("mapped Cluster Name Scan Fail")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func GetCustomerAddressesByIDAndMobile(customerID int, mobile string) ([]map[string]string, error) {
	query := `
		SELECT 
			id AS AddressID,
			customer_id AS CustomerID,
			name AS Name,
			mob_no AS Mobile,
			email AS Email,
			address1 AS Line1,
			address2 AS Line2,
			city AS City,
			state AS State,
			zipcode AS Pincode,
			town AS Town,
			location_type AS LocationType,
			location_info AS LocationInfo,
			is_default AS IsDefault,
			status AS Status,
			updated_date AS UpdatedDate
		FROM CustomerAddress
		WHERE customer_id = ? AND status = 'ACTIVE'
		ORDER BY is_default DESC, id ASC
	`

	rows, err := db.Db.Query(query, customerID)
	if err != nil {
		log.Println("[Error] : DB query failed in GetCustomerAddressesByIDAndMobile:", err)
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	var results []map[string]string

	for rows.Next() {
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			log.Println("[Error] : Scan failed:", err)
			return nil, err
		}

		entry := make(map[string]string)
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				entry[col] = string(b)
			} else if val != nil {
				entry[col] = fmt.Sprint(val)
			} else {
				entry[col] = ""
			}
		}
		results = append(results, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	log.Printf("[Debug] : Fetched %d addresses for customerID=%d (mobile=%s) => %+v\n", len(results), customerID, mobile, results)
	return results, nil
}

func GetCustomerOrderDetailsForNotification() (data [][]string, err error) {
	row, err := db.Db.Query(`SELECT customer_mob,cluster_name,product_name,product_package_qty,product_qty,order_units ,order_id,status,create_date,cluster_type FROM CustomerOrderTxn where status="ORDER_RECEIVED" ORDER BY create_date desc`)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Producer fail with status DEMAND_ORDER_CREATED")
		return
	}
	defer sql.Close(row)
	_, data, err = sql.Scan(row)
	if err != nil {
		log.Println(beego.AppConfig.String("loglevel"), "Error", err)
		err = errors.New("Producer fail with status DEMAND_ORDER_CREATED")
		return
	}

	log.Println(beego.AppConfig.String("loglevel"), "Debug", "Query Data - ", data, "\nData len - ", len(data))
	return

}

func GetCustomerOrderNotificationDuration(data [][]string) (data1 []string, err error) {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}

	currentTime := time.Now().In(loc)
	currentTime1 := currentTime.Format("2006-01-02 15:04:05")

	t2, err1 := time.ParseInLocation("2006-01-02 15:04:05", currentTime1, loc)
	if err1 != nil {
		fmt.Println(err1)
		return
	}

	log.Println("Debug", "Current Time :t2 - ", t2)

	for i := range data {
		t1, err := time.ParseInLocation("2006-01-02 15:04:05", data[i][8], loc)
		if err != nil {
			fmt.Println(err)
			continue
		}

		diff := t2.Sub(t1)
		diff1 := int(diff.Minutes())
		var time1 string

		if diff.Minutes() < 60 {
			time := int(diff.Minutes())
			time1 = strconv.Itoa(time) + " mins ago"
		} else if 60 <= diff1 && diff1 < 1440 {
			time := int(diff.Minutes() / 60)
			time1 = strconv.Itoa(time) + " hrs ago"
		} else if 1440 <= diff1 && diff1 < 10080 {
			time := int(diff.Minutes() / 1440)
			time1 = strconv.Itoa(time)
			if time1 == "1" {
				time1 = time1 + " day ago"
			} else {
				time1 = time1 + " days ago"
			}
		} else {
			time := int(diff.Minutes() / 10080)
			time1 = strconv.Itoa(time)
			if time1 == "1" {
				time1 = time1 + " week ago"
			} else {
				time1 = time1 + " weeks ago"
			}
		}

		data1 = append(data1, time1)
	}

	return
}
