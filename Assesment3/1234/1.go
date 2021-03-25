package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	//unique id generator
)

// Staff is the representation of staff table
type Staff struct {
	FirstName    string `json:"first_name" validate:"required,gte=10"`
	LastName     string
	EmailID      string
	MobileNumber string `json:"password" validate:"required,gte=10"`
	Password     string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "somu7chotu8"
	dbName := "go_employee"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

// Login is gonna create new template
func Login(w http.ResponseWriter, r *http.Request) {

	tmpl.ExecuteTemplate(w, "Login", nil)
}

// LoginProcess is defining the login functionality
func LoginProcess(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		emailID := r.FormValue("emailId")
		password := r.FormValue("password")
		fmt.Println(emailID, password)

		// Validate form input
		//trim out the spaces
		if strings.Trim(emailID, " ") == "" || strings.Trim(password, " ") == "" {
			fmt.Println("Parameter's can't be empty")
			http.Redirect(w, r, "/", 301)
			return
		}

		checkUser, err := db.Query("SELECT * FROM Staff ")

		if err != nil {
			panic(err.Error())
		}
		var staff []Staff
		a := 1
		for checkUser.Next() {
			var id Staff
			err = checkUser.Scan(&id.FirstName, &id.LastName, &id.EmailID, &id.MobileNumber, &id.Password)
			if err != nil {
				panic(err.Error())
			}
			staff = append(staff, id)
			if id.EmailID == emailID && id.Password == password {
				a = 0
			}
		}

		fmt.Println(staff)
		if a == 0 {
			fmt.Println("Success")
			tmpl.ExecuteTemplate(w, "Success", staff)
		} else {
			fmt.Println("err")
			http.Redirect(w, r, "/", 301)
		}

	}
	defer db.Close()

}

// Success is the representation of success page
func Success(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "Success", nil)
}

// Register is the represntation of list page
func Register(w http.ResponseWriter, r *http.Request) {

	tmpl.ExecuteTemplate(w, "Register", nil)
}

// RegisterUser is the representation of user registration
func RegisterUser(w http.ResponseWriter, r *http.Request) {

	db := dbConn()
	if r.Method == "POST" {

		firstName := r.FormValue("firstname")
		lastName := r.FormValue("lastname")
		EmailID := r.FormValue("emailId")
		mobileNumber := r.FormValue("mobilenumber")
		password := r.FormValue("password")
		fmt.Println(EmailID)

		tmpl.Execute(w, "Register")

		u := Staff{FirstName: firstName, LastName: lastName, EmailID: EmailID, MobileNumber: mobileNumber, Password: password}

		stmt, err := db.Query("INSERT INTO Staff VALUES('" + u.FirstName + "','" + u.LastName + "','" + u.EmailID + "','" + u.MobileNumber + "','" + u.Password + "')")
		fmt.Println(stmt)
		if err != nil {
			fmt.Print(err.Error())
		}

	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	log.Println("Server started on: http://localhost:8080")

	fs := http.FileServer(http.Dir("asset/"))
	http.Handle("/asset/", http.StripPrefix("/asset/", fs))

	http.HandleFunc("/", Login)
	http.HandleFunc("/loginsubmit", LoginProcess)
	http.HandleFunc("/register", Register)
	http.HandleFunc("/registeruser", RegisterUser)
	http.HandleFunc("/success", Success)

	http.ListenAndServe(":8080", nil)
}
