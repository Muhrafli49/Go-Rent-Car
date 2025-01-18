package controllers

import (
	"net/http"
	"rental-mobil/config"
	"rental-mobil/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetAllCustomers(c echo.Context) error {
	// Ambil parameter page dan limit dari query params
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	// Hitung offset
	offset := (page - 1) * limit

	var customers []models.Customer
	query := `SELECT * FROM customers LIMIT $1 OFFSET $2`
	err = config.DB.Select(&customers, query, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch customers"})
	}

	// Kembalikan data dengan informasi pagination
	return c.JSON(http.StatusOK, map[string]interface{}{
		"page":  page,
		"limit": limit,
		"data":  customers,
	})
}

func CreateCustomer(c echo.Context) error {
	customer := new(models.Customer)
	if err := c.Bind(customer); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
	}

	existsQuery := `SELECT COUNT(*) FROM customers WHERE nik = $1`
	var count int
	err := config.DB.QueryRow(existsQuery, customer.NIK).Scan(&count)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to check NIK"})
	}

	if count > 0 {
		return c.JSON(http.StatusConflict, map[string]string{"message": "NIK already registered"})
	}

	insertQuery := `INSERT INTO customers (name, nik, phone) VALUES ($1, $2, $3)`
	_, err = config.DB.Exec(insertQuery, customer.Name, customer.NIK, customer.Phone)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create customer"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Customer created successfully"})
}

func UpdateCustomer(c echo.Context) error {
    id := c.Param("id")

    // Cek apakah customer dengan id tersebut ada
    var count int
    checkQuery := `SELECT COUNT(*) FROM customers WHERE id = $1`
    err := config.DB.QueryRow(checkQuery, id).Scan(&count)
    if err != nil || count == 0 {
        return c.JSON(http.StatusNotFound, map[string]string{"message": "Customer not found"})
    }

    customer := new(models.Customer)
    // Bind input ke struktur customer
    if err := c.Bind(customer); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
    }

    // Validasi jika hanya field tertentu yang diubah
    if customer.Name == "" && customer.NIK == "" && customer.Phone == "" {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "No valid fields to update"})
    }

    // Cek apakah NIK sudah ada (kecuali untuk id yang sama)
    if customer.NIK != "" {
        existsQuery := `SELECT COUNT(*) FROM customers WHERE nik = $1 AND id != $2`
        var count int
        err := config.DB.QueryRow(existsQuery, customer.NIK, id).Scan(&count)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to check NIK"})
        }
        if count > 0 {
            return c.JSON(http.StatusConflict, map[string]string{"message": "NIK already registered"})
        }
    }

    // Update hanya field yang ada
    query := `UPDATE customers SET 
        name = COALESCE($1, name), 
        nik = COALESCE($2, nik), 
        phone = COALESCE($3, phone) 
        WHERE id = $4`
    _, err = config.DB.Exec(query, customer.Name, customer.NIK, customer.Phone, id)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update customer"})
    }

    return c.JSON(http.StatusOK, map[string]string{"message": "Customer updated successfully"})
}


func DeleteCustomer(c echo.Context) error {
	id := c.Param("id")

	// Cek apakah customer dengan id tersebut ada
	var count int
	checkQuery := `SELECT COUNT(*) FROM customers WHERE id = $1`
	err := config.DB.QueryRow(checkQuery, id).Scan(&count)
	if err != nil || count == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Customer not found"})
	}

	query := `DELETE FROM customers WHERE id=$1`
	result, err := config.DB.Exec(query, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete customer"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Customer not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Customer deleted successfully"})
}
