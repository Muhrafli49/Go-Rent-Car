package controllers

import (
	"net/http"
	"rental-mobil/config"
	"rental-mobil/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GetAllBookings mengambil semua data booking
func GetAllBookings(c echo.Context) error {
	// Get page and limit from query parameters, defaulting to page 1 and limit 5 if not provided
	page := 1
	limit := 5
	if p := c.QueryParam("page"); p != "" {
		var err error
		page, err = strconv.Atoi(p)
		if err != nil || page <= 0 {
			page = 1
		}
	}
	if l := c.QueryParam("limit"); l != "" {
		var err error
		limit, err = strconv.Atoi(l)
		if err != nil || limit <= 0 {
			limit = 5
		}
	}

	// Calculate offset for pagination
	offset := (page - 1) * limit

	var bookings []models.Booking
	query := `SELECT id, customer_id, car_id, start_rent, end_rent, total_cost, finished FROM bookings 
              LIMIT $1 OFFSET $2`
	err := config.DB.Select(&bookings, query, limit, offset)
	if err != nil {
		c.Logger().Error("Error fetching bookings:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch bookings"})
	}

	// Get the total number of bookings to calculate total pages
	var totalBookings int
	countQuery := `SELECT COUNT(*) FROM bookings`
	err = config.DB.Get(&totalBookings, countQuery)
	if err != nil {
		c.Logger().Error("Error counting bookings:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to count bookings"})
	}

	// Calculate the total number of pages
	totalPages := (totalBookings + limit - 1) / limit

	// Prepare the response data including pagination info
	response := map[string]interface{}{
		"page":        page,
		"limit":       limit,
		"totalPages":  totalPages,
		"totalResults": totalBookings,
		"data":        bookings,
	}

	c.Logger().Info("Fetched bookings:", bookings)
	return c.JSON(http.StatusOK, response)
}

// CreateBooking membuat data booking baru
func CreateBooking(c echo.Context) error {
	booking := new(models.Booking)
	if err := c.Bind(booking); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
	}

	// Validasi customer_id dan car_id
	if booking.CustomerID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid customer ID"})
	}
	if booking.CarID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid car ID"})
	}

	// Validasi start_rent dan end_rent
	if booking.StartRent.After(booking.EndRent) {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Start rent must be before end rent"})
	}

	// Hitung total_cost
	var dailyRent float64
	query := `SELECT daily_rent FROM cars WHERE id = $1`
	err := config.DB.Get(&dailyRent, query, booking.CarID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch car data"})
	}
	duration := booking.EndRent.Sub(booking.StartRent).Hours() / 24
	booking.TotalCost = dailyRent * duration

	// Insert data booking baru
	insertQuery := `INSERT INTO bookings (customer_id, car_id, start_rent, end_rent, total_cost, finished) 
                    VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = config.DB.Exec(insertQuery, booking.CustomerID, booking.CarID, booking.StartRent, booking.EndRent, booking.TotalCost, booking.Finished)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create booking"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Booking created successfully"})
}

// UpdateBooking memperbarui data booking
func UpdateBooking(c echo.Context) error {
    id := c.Param("id")
    booking := new(models.Booking)

    if err := c.Bind(booking); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
    }

    // Validasi customer_id dan car_id
    if booking.CustomerID <= 0 {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid customer ID"})
    }
    if booking.CarID <= 0 {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid car ID"})
    }

    // Validasi start_rent dan end_rent
    if booking.StartRent.After(booking.EndRent) {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "Start rent must be before end rent"})
    }

    // Mengecek apakah booking dengan id tersebut ada
    var existingBooking models.Booking
    checkQuery := `SELECT id FROM bookings WHERE id = $1`
    err := config.DB.Get(&existingBooking, checkQuery, id)
    if err != nil {
        return c.JSON(http.StatusNotFound, map[string]string{"message": "Booking not found"})
    }

    // Hitung total_cost
    var dailyRent float64
    query := `SELECT daily_rent FROM cars WHERE id = $1`
    err = config.DB.Get(&dailyRent, query, booking.CarID)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch car data"})
    }

    duration := booking.EndRent.Sub(booking.StartRent).Hours() / 24
    booking.TotalCost = dailyRent * duration

    // Update data booking
    updateQuery := `UPDATE bookings SET customer_id=$1, car_id=$2, start_rent=$3, end_rent=$4, total_cost=$5, finished=$6 WHERE id=$7`
    _, err = config.DB.Exec(updateQuery, booking.CustomerID, booking.CarID, booking.StartRent, booking.EndRent, booking.TotalCost, booking.Finished, id)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update booking"})
    }

    return c.JSON(http.StatusOK, map[string]string{"message": "Booking updated successfully"})
}

// DeleteBooking menghapus data booking
func DeleteBooking(c echo.Context) error {
	id := c.Param("id")
	query := `DELETE FROM bookings WHERE id=$1`
	result, err := config.DB.Exec(query, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete booking"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Booking not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Booking deleted successfully"})
}
