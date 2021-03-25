//go get -u github.com/gin-gonic/gin
//go get -u github.com/go-sql-driver/mysql

package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-gorp/gorp"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id        int64  `db:"ID" json:"id"`
	Username  string `db:"Username" json:"username"`
	Password  string `db:"Password" json:"password"`
	Firstname string `db:"Firstname" json:"firstname"`
	Lastname  string `db:"Lastname" json:"lastname"`
}

func main() {

	CreateUrlMappings()
	// Listen and server on 0.0.0.0:8080
	Router.Run(":8080")

}

var Router *gin.Engine

func CreateUrlMappings() {
	Router = gin.Default()

	Router.LoadHTMLGlob("templates/*")

	Router.Use(Cors())
	// v1 of the API
	v1 := Router.Group("/v1")
	{
		v1.GET("/users/:id", GetUserDetail)
		v1.GET("/users/", GetUser)
		v1.POST("/login/", Login)
		v1.PUT("/users/:id", UpdateUser)
		v1.POST("/users", PostUser)
	}
}

var dbmap = initDb()

func initDb() *gorp.DbMap {
	db, err := sql.Open("mysql", "root:somu7chotu8@tcp(localhost:3306)/test")
	checkErr(err, "sql.Open failed")
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func GetUser(c *gin.Context) {
	var user []User

	Router.LoadHTMLGlob("templates/*")
	_, err := dbmap.Select(&user, `select * from user `)

	if err == nil {
		c.JSON(200, user)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":   "Home Page",
			"payload": user,
		},
		)
	} else {
		c.JSON(404, gin.H{"error": "user not found"})
	}

}

func GetUserDetail(c *gin.Context) {
	id := c.Params.ByName("id")
	var user User
	Router.LoadHTMLGlob("templates/*")
	err := dbmap.SelectOne(&user, "SELECT * FROM user WHERE id=? LIMIT 1", id)

	if err == nil {
		user_id, _ := strconv.ParseInt(id, 0, 64)

		content := &User{
			Id:        user_id,
			Username:  user.Username,
			Password:  user.Password,
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
		}
		c.JSON(200, content)
		c.HTML(http.StatusOK, "User.html", gin.H{
			"title":   "Home Page",
			"payload": content,
		},
		)
	} else {
		c.JSON(404, gin.H{"error": "user not found"})
	}
}

func Login(c *gin.Context) {
	var user User
	c.Bind(&user)
	Router.LoadHTMLGlob("templates/*")
	err := dbmap.SelectOne(&user, "select * from user where Username=? LIMIT 1", user.Username)

	if err == nil {
		user_id := user.Id

		content := &User{
			Id:        user_id,
			Username:  user.Username,
			Password:  user.Password,
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
		}
		c.JSON(200, content)
		c.HTML(http.StatusOK, "Login.html", gin.H{
			"title":   "Home Page",
			"payload": content,
		},
		)
	} else {
		c.JSON(404, gin.H{"error": "user not found"})
	}

}

func PostUser(c *gin.Context) {
	var user User
	c.Bind(&user)
	Router.LoadHTMLGlob("templates/*")
	log.Println(user)

	if user.Username != "" && user.Password != "" && user.Firstname != "" && user.Lastname != "" {

		if insert, _ := dbmap.Exec(`INSERT INTO user (Username, Password, Firstname, Lastname) VALUES (?, ?, ?, ?)`, user.Username, user.Password, user.Firstname, user.Lastname); insert != nil {
			log.Println(user.Username)
			user_id, err := insert.LastInsertId()
			//to use this ...you have to set id to auto increment in db
			log.Println(user_id)
			if err == nil {
				content := &User{
					Id:        user_id,
					Username:  user.Username,
					Password:  user.Password,
					Firstname: user.Firstname,
					Lastname:  user.Lastname,
				}
				c.JSON(201, content)
				c.HTML(http.StatusOK, "Register.html", gin.H{
					"title":   "Home Page",
					"payload": content,
				},
				)
			} else {
				checkErr(err, "Insert failed")
			}
		}

	} else {
		c.JSON(400, gin.H{"error": "Fields are empty"})
	}

}

func UpdateUser(c *gin.Context) {
	id := c.Params.ByName("id")
	var user User
	Router.LoadHTMLGlob("templates/*")
	err := dbmap.SelectOne(&user, "SELECT * FROM user WHERE id=?", id)

	if err == nil {
		var json User
		c.Bind(&json)

		user_id, _ := strconv.ParseInt(id, 0, 64)

		user := User{
			Id:        user_id,
			Username:  user.Username,
			Password:  user.Password,
			Firstname: json.Firstname,
			Lastname:  json.Lastname,
		}

		if user.Firstname != "" && user.Lastname != "" {
			_, err = dbmap.Update(&user)

			if err == nil {
				c.JSON(200, user)
				c.HTML(http.StatusOK, "updateUser.html", gin.H{
					"title":   "Home Page",
					"payload": user,
				},
				)
			} else {
				checkErr(err, "Updated failed")
			}

		} else {
			c.JSON(400, gin.H{"error": "fields are empty"})
		}

	} else {
		c.JSON(404, gin.H{"error": "user not found"})
	}
}
