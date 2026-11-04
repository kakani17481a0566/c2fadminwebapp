// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	currentTime := time.Now()

	// get 12 months name and year from current time

	m1year := strconv.Itoa(currentTime.Year())
	month1 := time.Month(currentTime.Month()).String()
	fmt.Println("Month 1 and Year:", month1, m1year)

	m2 := currentTime.AddDate(0, 0, 30)
	m2year := strconv.Itoa(m2.Year())
	month2 := time.Month(m2.Month()).String()
	fmt.Println("Month 2 and Year:", month2, m2year)

	m3 := currentTime.AddDate(0, 0, 60)
	m3year := strconv.Itoa(m3.Year())
	month3 := time.Month(m3.Month()).String()
	fmt.Println("Month 3 and Year:", month3, m3year)

	m4 := currentTime.AddDate(0, 0, 90)
	m4year := strconv.Itoa(m4.Year())
	month4 := time.Month(m4.Month()).String()
	fmt.Println("Month 4 and Year:", month4, m4year)

	m5 := currentTime.AddDate(0, 0, 120)
	m5year := strconv.Itoa(m5.Year())
	month5 := time.Month(m5.Month()).String()
	fmt.Println("Month 5 and Year:", month5, m5year)

	m6 := currentTime.AddDate(0, 0, 150)
	m6year := strconv.Itoa(m6.Year())
	month6 := time.Month(m6.Month()).String()
	fmt.Println("Month 6 and Year:", month6, m6year)

	m7 := currentTime.AddDate(0, 0, 180)
	m7year := strconv.Itoa(m7.Year())
	month7 := time.Month(m7.Month()).String()
	fmt.Println("Month 7 and Year:", month7, m7year)

	m8 := currentTime.AddDate(0, 0, 210)
	m8year := strconv.Itoa(m8.Year())
	month8 := time.Month(m8.Month()).String()
	fmt.Println("Month 8 and Year:", month8, m8year)

	m9 := currentTime.AddDate(0, 0, 240)
	m9year := strconv.Itoa(m9.Year())
	month9 := time.Month(m9.Month()).String()
	fmt.Println("Month 9 and Year:", month9, m9year)

	m10 := currentTime.AddDate(0, 0, 270)
	m10year := strconv.Itoa(m10.Year())
	month10 := time.Month(m10.Month()).String()
	fmt.Println("Month 10 and Year:", month10, m10year)

	m11 := currentTime.AddDate(0, 0, 300)
	m11year := strconv.Itoa(m11.Year())
	month11 := time.Month(m11.Month()).String()
	fmt.Println("Month 11 and Year:", month11, m11year)

	m12 := currentTime.AddDate(0, 0, 330)
	m12year := strconv.Itoa(m12.Year())
	month12 := time.Month(m12.Month()).String()
	fmt.Println("Month 12 and Year:", month12, m12year)

}
