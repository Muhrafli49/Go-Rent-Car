package controllers

import (
	"net/http"
	"rental-mobil/config"
	"rental-mobil/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GetAllCars mengambil semua data mobil dengan pagination
func GetAllCars(c echo.Context) error {
	// Ambil parameter page dan limit dari query params
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit < 1 || limit > 5 {
		limit = 5 // Batas maksimal limit adalah 5
	}

	// Hitung offset
	offset := (page - 1) * limit

	var cars []models.Car
	query := `SELECT id, name, stock, daily_rent FROM cars LIMIT $1 OFFSET $2`
	err = config.DB.Select(&cars, query, limit, offset)

	if err != nil {
		c.Logger().Error("Error fetching cars:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch cars"})
	}

	// Log data yang berhasil diambil
	c.Logger().Info("Fetched cars:", cars)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"page":  page,
		"limit": limit,
		"data":  cars,
	})
}

// CreateCar membuat data mobil baru
func CreateCar(c echo.Context) error {
	car := new(models.Car)
	if err := c.Bind(car); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
	}

	// Validasi name, stock, dan daily rent
	if car.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Car name is required"})
	}
	if car.Stock <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Stock must be greater than zero"})
	}
	if car.DailyRent <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Daily Rent must be greater than zero"})
	}

	// Insert data mobil baru
	insertQuery := `INSERT INTO cars (name, stock, daily_rent) VALUES ($1, $2, $3)`
	_, err := config.DB.Exec(insertQuery, car.Name, car.Stock, car.DailyRent)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create car"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Car created successfully"})
}

// UpdateCar memperbarui data mobil
func UpdateCar(c echo.Context) error {
	id := c.Param("id")
	car := new(models.Car)
	if err := c.Bind(car); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
	}

	// Validasi name, stock, dan daily rent
	if car.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Car name is required"})
	}
	if car.Stock <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Stock must be greater than zero"})
	}
	if car.DailyRent <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Daily Rent must be greater than zero"})
	}

	// Update data mobil
	query := `UPDATE cars SET name=$1, stock=$2, daily_rent=$3 WHERE id=$4`
	_, err := config.DB.Exec(query, car.Name, car.Stock, car.DailyRent, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update car"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Car updated successfully"})
}

// DeleteCar menghapus data mobil
func DeleteCar(c echo.Context) error {
	id := c.Param("id")
	query := `DELETE FROM cars WHERE id=$1`
	result, err := config.DB.Exec(query, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete car"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Car not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Car deleted successfully"})
}
