package password

import (
	"math/rand"
	"time"
)

var uAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
var lAplgabet = "abcdefghijklmnopqrstuvwxyz"
var number = "0123456789"
var special = "!@#$%^&*()_+-"

func AlphaNumericSpecial(n int) (pass string, err error) {
	flag := 0
	for i := 0; i < n; i++ {
		rand.Seed(time.Now().UnixNano())
		if flag == 0 {
			randnum := rand.Intn(len(uAlphabet))
			pass = pass + string(uAlphabet[randnum])
			flag++
		} else if flag == 1 {
			randnum := rand.Intn(len(lAplgabet))
			pass = pass + string(lAplgabet[randnum])
			flag++
		} else if flag == 2 {
			randnum := rand.Intn(len(number))
			pass = pass + string(number[randnum])
			flag++
		} else {
			randnum := rand.Intn(len(special))
			pass = pass + string(special[randnum])
			flag = 0
		}
	}
	return
}

func Numeric(n int) (pass string, err error) {

	for i := 0; i < n; i++ {
		rand.Seed(time.Now().UnixNano())

		randnum := rand.Intn(len(number))
		pass = pass + string(number[randnum])

	}
	return
}
