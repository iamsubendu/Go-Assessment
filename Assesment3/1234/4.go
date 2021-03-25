//go get github.com/google/uuid

package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	db, err := sql.Open("mysql", "root:somu7chotu8@tcp(127.0.0.1:3306)/go_employee")
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	fmt.Println("Connected to Mysql database")

	router := gin.Default()

	router.LoadHTMLGlob("cards/*")

	type Product struct {
		ID    int    `json; "id"`
		Name  string `json: "name"`
		Price string `json: "price"`
	}

	router.GET("/:id", func(c *gin.Context) {
		var (
			product Product
			result  gin.H
		)
		id := c.Param("id")
		row := db.QueryRow("select * from product where id = ?;", id)
		err = row.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			// If no results send null
			result = gin.H{
				"Result": "their is some error please check",
			}
		} else {
			result = gin.H{
				"Result": product,
				"Count":  1,
			}
		}
		c.JSON(http.StatusOK, result)
	})

	// GET all product
	router.GET("/", func(c *gin.Context) {
		var (
			product  Product
			products []Product
		)
		rows, err := db.Query("select * from product;")
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&product.ID, &product.Name, &product.Price)
			products = append(products, product)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"Result": products,
			"Count":  len(products),
		})
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":   "Home Page",
			"payload": products,
		},
		)
	})

	// POST new product details
	router.POST("/", func(c *gin.Context) {
		var buffer bytes.Buffer
		id := c.PostForm("id")
		name := c.PostForm("name")
		price := c.PostForm("price")
		stmt, err := db.Prepare("insert into product (id, name, price) values(?,?,?);")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(id, name, price)

		if err != nil {
			fmt.Print(err.Error())
		}

		// Fastest way to append strings
		buffer.WriteString(name)
		buffer.WriteString(" ")
		buffer.WriteString(price)
		defer stmt.Close()
		datanya := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"Post": fmt.Sprintf(" Added product name and price of data %s ", datanya),
		})
	})

	// PUT - update a product details
	router.PUT("/", func(c *gin.Context) {
		var buffer bytes.Buffer
		id := c.PostForm("id")
		name := c.PostForm("name")
		price := c.PostForm("price")
		stmt, err := db.Prepare("update product set name= ?, price = ?;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(name, price, id)
		if err != nil {
			fmt.Print(err.Error())
		}

		// Fastest way to append strings
		buffer.WriteString(name)
		buffer.WriteString(" ")
		buffer.WriteString(price)
		defer stmt.Close()
		datanya := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"Updated": fmt.Sprintf("Updated resource of  Id %s Product name and price %s", id, datanya),
		})
	})

	//Delete resource
	router.DELETE("/", func(c *gin.Context) {
		id := c.PostForm("id")
		stmt, err := db.Prepare("delete from product where id= ?;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(id)
		if err != nil {
			fmt.Print(err.Error())
		}
		c.JSON(http.StatusOK, gin.H{
			"Deleted": fmt.Sprintf("Deleted resource %s", id),
		})
	})

	router.Run(":8080")
}
