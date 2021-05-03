package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"text/template"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/joho/godotenv"
)

type Order struct {
	Coupon   string
	CcNumber string
}

type Result struct {
	Status string
}

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	fmt.Println("Checkouts MS running on port :9090")

	http.HandleFunc("/", home)
	http.HandleFunc("/process", process)
	http.ListenAndServe(":9090", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	// How to use templates using go utils
	t := template.Must(template.ParseFiles("templates/home.html"))

	t.Execute(w, Result{})
}

func process(w http.ResponseWriter, r *http.Request) {
	result := makeHttpCall("http://localhost:9091", r.FormValue("coupon"), r.FormValue("cc-number"))

	t := template.Must(template.ParseFiles("templates/home.html"))

	t.Execute(w, result)
}

func makeHttpCall(urlMicroservice string, coupon string, ccNumber string) Result {
	values := url.Values{}

	values.Add("coupon", coupon)
	values.Add("ccNumber", ccNumber)

	retryClient := retryablehttp.NewClient()

	// Retries amount
	retryClient.RetryMax = 5

	res, err := retryClient.PostForm(urlMicroservice, values)

	if err != nil {
		result := Result{Status: "Sorry, service unavailable temporarily"}

		return result
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal("Error processing the result")
	}

	result := Result{}

	json.Unmarshal(data, &result)

	return result
}
