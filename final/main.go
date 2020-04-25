package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

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
	fmt.Fprintf(w, "Hello ") // write data to response
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/login.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// logic part of log in
		pswd := r.FormValue("password")
		flag := verify(pswd)
		if flag == true {
			fmt.Fprintf(w, "success")
		}

		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}
}
func verify(pswd string) bool {
	file, err := os.Open("test.txt")
	var flag bool
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	file.Close()
	for _, l := range lines {
		strarray := strings.Fields(l)
		bytehash := []byte(strarray[2])
		err := bcrypt.CompareHashAndPassword([]byte(pswd),bytehash)
		if err != nil {
			log.Println(err)
			flag=true
			break
		}
		flag=false
	}
	return flag
}
func register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/register.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// logic part of log in
		uName := r.FormValue("username")
		email := r.FormValue("email")
		pwd := r.FormValue("password")
		pwdConfirm := r.FormValue("confirmpassword")
		if pwd == pwdConfirm {
			hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("hash", string(hash))
			data := []string{uName, email, string(hash)}
			data1 := strings.Join(data, " ")
			f, err := os.OpenFile("test.txt", os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = fmt.Fprintln(f, data1)
			http.Redirect(w, r, "/login", 301)
		} else {
			fmt.Fprintln(w, "Password information must be the same.")
		}
	}
}
func main() {
	http.HandleFunc("/", sayhelloName) // setting router rule
	http.HandleFunc("/login", login)
	//http.HandleFunc("/index", index)
	http.HandleFunc("/register", register)
	err := http.ListenAndServe(":9090", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
