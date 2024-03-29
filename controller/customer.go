package controllers

import (
	"net/http"
	"strconv"
	"time"

	"Arutalalab-technical-test/models"

	"github.com/gin-gonic/gin"
)

// CreateCustomer handles the creation of a new customer
func (ctrl *Controller) CreateCustomer(c *gin.Context) {
	var customer models.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert into the database with the specified id
	result, err := ctrl.DB.Exec("INSERT INTO customer (id, name, email, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		customer.ID, customer.Name, customer.Email, time.Now(), time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
		return
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve rows affected"})
		return
	}

	// Check if any rows were affected, if none, return an error
	if rowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No rows were affected"})
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{"message": "Customer created successfully", "data": customer})
}

// GetCustomers handles the retrieval of all customers
func (ctrl *Controller) GetCustomers(c *gin.Context) {
	rows, err := ctrl.DB.Query("SELECT * FROM customer")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve customers"})
		return
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var customer models.Customer
		err := rows.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.CreatedAt, &customer.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan customer"})
			return
		}
		customers = append(customers, customer)
	}

	c.JSON(http.StatusOK, map[string]interface{}{"message": "Get Customers successfully", "data": customers})
}

// GetCustomer handles the retrieval of a single customer by ID
func (ctrl *Controller) GetCustomer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	var customer models.Customer
	err = ctrl.DB.QueryRow("SELECT * FROM customer WHERE id = ?", id).
		Scan(&customer.ID, &customer.Name, &customer.Email, &customer.CreatedAt, &customer.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"message": "Get Customers by ID successfully", "data": customer})
}

// UpdateCustomer handles the update of a customer by ID
func (ctrl *Controller) UpdateCustomer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	var updatedCustomer models.Customer
	if err := c.ShouldBindJSON(&updatedCustomer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the customer in the database
	_, err = ctrl.DB.Exec("UPDATE customer SET name = ?, email = ?, updated_at = ? WHERE id = ?",
		updatedCustomer.Name, updatedCustomer.Email, time.Now(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer"})
		return
	}

	// Return the updated customer
	updatedCustomer.ID = id
	c.JSON(http.StatusOK, map[string]interface{}{"message": "Update Customers successfully", "data": updatedCustomer})
}

// DeleteCustomer handles the deletion of a customer by ID
func (ctrl *Controller) DeleteCustomer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	// Delete the customer from the database
	_, err = ctrl.DB.Exec("DELETE FROM customer WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete customer"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
