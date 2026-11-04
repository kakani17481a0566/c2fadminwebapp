package sql

import (
	"log"
	"testing"
	"time"
)

func TestDB(t *testing.T) {
	var obj Database
	obj.Username = "postgres"
	obj.Password = "hangover"
	obj.Ip = "localhost"
	obj.Port = "5432"
	obj.Schema = "BillingEngine"
	obj.Type = "postgres"

	err := obj.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	log.Println("Connect Done")
	time.Sleep(10 * time.Second)
	row, err := obj.Query("select * from logs.transactions_copy")
	defer row.Close()
	if err != nil {
		t.Error(err)
		return

	}
	log.Println("Row Done")
	time.Sleep(10 * time.Second)
	_, data, err := Scan(row)
	if err != nil {
		t.Error(err)
		return
	}
	log.Println("Scan Done")
	time.Sleep(20 * time.Second)
	log.Println(data[0][0])
	time.Sleep(20 * time.Second)
}
