package controllers

import (
	"net/http"
	"strconv"
	"time"

	"Arutalalab-technical-test/models"

	"github.com/gin-gonic/gin"
)

// CreateOrder handles the creation of a new order
func (ctrl *Controller) CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if product exists
	var product models.Product
	err := ctrl.DB.QueryRow("SELECT * FROM product WHERE id = ?", order.ProductID).
		Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found"})
		return
	}

	// Check if customer exists
	var customer models.Customer
	err = ctrl.DB.QueryRow("SELECT * FROM customer WHERE id = ?", order.CustomerID).
		Scan(&customer.ID, &customer.Name, &customer.Email, &customer.CreatedAt, &customer.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer not found"})
		return
	}

	// Calculate total price (assuming quantity * product price)
	order.Total = float64(order.Quantity) * product.Price

	// Insert into the database
	result, err := ctrl.DB.Exec("INSERT INTO order (product_id, customer_id, quantity, total, created_at) VALUES (?, ?, ?, ?, ?)",
		order.ProductID, order.CustomerID, order.Quantity, order.Total, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Get the ID of the newly created order
	orderID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order ID"})
		return
	}

	// Set the ID in the order struct
	order.ID = int(orderID)

	c.JSON(http.StatusCreated, order)
}

// GetOrder handles the retrieval of a single order by ID
func (ctrl *Controller) GetOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var order models.Order
	err = ctrl.DB.QueryRow("SELECT * FROM orders WHERE id = ?", id).
		Scan(&order.ID, &order.ProductID, &order.CustomerID, &order.Quantity, &order.Total, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}
