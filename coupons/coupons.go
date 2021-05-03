package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Coupon struct {
	Code string
}

type Coupons struct {
	Coupon []Coupon
}

func (c Coupons) Validate(code string) string {
	for _, item := range c.Coupon {
		if code == item.Code {
			return "valid"
		}
	}

	return "invalid"
}

type Result struct {
	Status string
}

var coupons Coupons

func main() {
	fmt.Println("Coupons MS running on port :9092")

	coupon := Coupon{
		Code: "abc",
	}

	coupons.Coupon = append(coupons.Coupon, coupon)

	http.HandleFunc("/", home)
	http.ListenAndServe(":9092", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	coupon := r.FormValue("coupon")

	valid := coupons.Validate(coupon)

	result := Result{Status: valid}

	jsonResult, err := json.Marshal(result)

	if err != nil {
		log.Fatal("Error converting data to JSON")
	}

	fmt.Fprintf(w, string(jsonResult))
}
