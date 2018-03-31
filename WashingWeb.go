package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type Wash struct {
	Date  time.Time
	Usage int
	Name  string
}

type WashingData struct {
	PageTitle      string
	Washings       []Wash
	SummedWashings map[string]int
}

var data = WashingData{
	PageTitle:      "My Washing list",
	Washings:       []Wash{},
	SummedWashings: make(map[string]int),
}

func register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("layout.html")
		t.Execute(w, data)
	} else {
		r.ParseForm()

		fmt.Println("Date:", time.Now())
		fmt.Println("Usage:", r.Form["usage"])
		fmt.Println("Name:", r.Form["name"])
		usageInput, _ := strconv.Atoi(r.Form["usage"][0])
		data.Washings = append(data.Washings, Wash{Date: time.Now(), Usage: usageInput, Name: r.Form["name"][0]})
		t, _ := template.ParseFiles("layout.html")
		data.SummedWashings[r.Form["name"][0]] += usageInput
		t.Execute(w, data)
	}
}

func delete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("layout.html")

		id := r.URL.Query()["id"]
		fmt.Println("Row id to remove ", id[0])
		idAsInt, _ := strconv.Atoi(id[0])
		data.SummedWashings[data.Washings[idAsInt].Name] -= data.Washings[idAsInt].Usage
		data.Washings = data.Washings[:idAsInt+copy(data.Washings[idAsInt:], data.Washings[idAsInt+1:])]

		t.Execute(w, data)
	}
}

func main() {
	tmpl := template.Must(template.ParseFiles("layout.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data.PageTitle = "Washing calculator"
		tmpl.Execute(w, data)
	})

	http.HandleFunc("/register", register)
	http.HandleFunc("/delete", delete)

	http.ListenAndServe(":8080", nil)
}
