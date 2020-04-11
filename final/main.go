package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)
type name struct {
	username, email string
	psword          string
}
func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //Parse url parameters passed, then parse the response packet for the POST body (request body)
	// attention: If you do not call ParseForm method, the following data can not be obtained form
	fmt.Println(r.Form) // print information on server side.
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") // write data to response
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/login.html")
			t.Execute(w, nil)
	} else {
		r.ParseForm()
		// logic part of log in
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}
}
func register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/register.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// logic part of log in
		uName:= r.FormValue("username")
		email:= r.FormValue("email")
		pwd:= r.FormValue("password")
		pwdConfirm:= r.FormValue("confirmpassword")
		if pwd == pwdConfirm {
			hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
			if err != nil {
				// TODO: Properly handle error
				log.Fatal(err)
			}
			fmt.Println("hash",string(hash))
			data:=name{username: uName,email: email,psword:string(hash) }
			fmt.Println("data",data)
			file, _ := json.MarshalIndent(data, "", " ")

			_ = ioutil.WriteFile("test.json", file, 0644)
			_, _ = fmt.Fprintln(w, "Registration successful.")
		} else {
			fmt.Fprintln(w, "Password information must be the same.")
		}
	}
}
func main() {
	http.HandleFunc("/", sayhelloName) // setting router rule
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	err := http.ListenAndServe(":9090", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}