package controllers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"Arutalalab-technical-test/models"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	DB *sql.DB
}

// Create a constructor function for the controller
func NewController(db *sql.DB) *Controller {
	return &Controller{DB: db}
}

// CreateProduct handles the creation of a new product
func (ctrl *Controller) CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert into the database
	result, err := ctrl.DB.Exec("INSERT INTO product (name, price, stock, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		product.Name, product.Price, product.Stock, time.Now(), time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	// Get the ID of the newly created product
	productID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product ID"})
		return
	}

	// Set the ID in the product struct
	product.ID = int(productID)

	c.JSON(http.StatusCreated, product)
}

// GetProducts handles the retrieval of all products
func (ctrl *Controller) GetProducts(c *gin.Context) {
	rows, err := ctrl.DB.Query("SELECT * FROM product")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products", "details": err.Error()})
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan product", "details": err.Error()})
			return
		}
		products = append(products, product)
	}

	c.JSON(http.StatusOK, products)
}

// GetProduct handles the retrieval of a single product by ID
func (ctrl *Controller) GetProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	err = ctrl.DB.QueryRow("SELECT * FROM product WHERE id = ?", id).
		Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProduct handles the update of a product by ID
func (ctrl *Controller) UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var updatedProduct models.Product
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the product in the database
	_, err = ctrl.DB.Exec("UPDATE product SET name = ?, price = ?, stock = ?, updated_at = ? WHERE id = ?",
		updatedProduct.Name, updatedProduct.Price, updatedProduct.Stock, time.Now(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	// Return the updated product
	updatedProduct.ID = id
	c.JSON(http.StatusOK, updatedProduct)
}

// DeleteProduct handles the deletion of a product by ID
func (ctrl *Controller) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Delete the product from the database
	_, err = ctrl.DB.Exec("DELETE FROM product WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
