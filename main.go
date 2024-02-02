// main.go
package main

import (
	"log"

	"Arutalalab-technical-test/config"
	controllers "Arutalalab-technical-test/controller"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := config.ConnectDb()
	defer db.Close()

	controller := controllers.NewController(db)

	router := gin.Default()

	// Product Endpoints
	router.POST("/products", controller.CreateProduct)
	router.GET("/products", controller.GetProducts)
	router.GET("/products/:id", controller.GetProduct)
	router.PUT("/products/:id", controller.UpdateProduct)
	router.DELETE("/products/:id", controller.DeleteProduct)

	// Customer Endpoints
	router.POST("/customers", controller.CreateCustomer)
	router.GET("/customers", controller.GetCustomers)
	router.GET("/customers/:id", controller.GetCustomer)
	router.PUT("/customers/:id", controller.UpdateCustomer)
	router.DELETE("/customers/:id", controller.DeleteCustomer)

	// Order Endpoints
	router.POST("/orders", controller.CreateOrder)
	router.GET("/orders/:id", controller.GetOrder)

	// Run the server
	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
