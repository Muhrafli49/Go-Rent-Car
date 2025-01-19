package controllers

import (
	// "log"
	"net/http"
	"rental-mobil/config"
	"rental-mobil/models"
	"strconv"
	"time"

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
		"page":         page,
		"limit":        limit,
		"totalPages":   totalPages,
		"totalResults": totalBookings,
		"data":         bookings,
	}

	c.Logger().Info("Fetched bookings:", bookings)
	return c.JSON(http.StatusOK, response)
}

// CreateBooking membuat data booking baru
func CreateBooking(c echo.Context) error {
    booking := new(models.Booking)

    // Bind the request data into the booking struct
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
    if booking.StartRent > booking.EndRent {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "Start rent must be before end rent"})
    }

    // Ambil data membership berdasarkan customer_id
    var membership models.Membership
    query := `
        SELECT m.id, m.name, m.discount
        FROM membership m
        JOIN customers c ON c.membership_id = m.id
        WHERE c.id = $1
    `
    err := config.DB.Get(&membership, query, booking.CustomerID)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch membership data"})
    }

    // Ambil data mobil berdasarkan car_id
    var car models.Car
    carQuery := `SELECT id, daily_rent FROM cars WHERE id = $1`
    err = config.DB.Get(&car, carQuery, booking.CarID)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch car data"})
    }

    // Mengonversi tanggal mulai dan selesai sewa ke tipe time.Time
    startRentTime, err := time.Parse("2006-01-02", booking.StartRent)
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid start rent date format"})
    }
    endRentTime, err := time.Parse("2006-01-02", booking.EndRent)
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid end rent date format"})
    }

    // Hitung durasi sewa (dalam hari)
    duration := int(endRentTime.Sub(startRentTime).Hours() / 24)

    // Validasi durasi sewa
    if duration <= 0 {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "End rent date must be after start rent date"})
    }

    // Hitung total biaya sewa
    booking.TotalCost = car.DailyRent * float64(duration) * (1 - membership.Discount/100)

    // Ambil data driver berdasarkan driver_id
    var driver models.Driver
    driverQuery := `SELECT id, daily_cost FROM driver WHERE id = $1`
    err = config.DB.Get(&driver, driverQuery, booking.DriverID)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch driver data"})
    }

    // Hitung biaya supir
    booking.TotalDriverCost = driver.DailyCost * float64(duration)

    // Simpan booking ke database
    insertQuery := `INSERT INTO bookings (customer_id, car_id, start_rent, end_rent, total_cost, finished, discount, booking_type_id, driver_id, total_driver_cost) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
    _, err = config.DB.Exec(insertQuery, booking.CustomerID, booking.CarID, booking.StartRent, booking.EndRent, booking.TotalCost, booking.Finished, membership.Discount, booking.BookingTypeID, booking.DriverID, booking.TotalDriverCost)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create booking"})
    }

    return c.JSON(http.StatusCreated, map[string]string{"message": "Booking created successfully"})
}


// UpdateBooking memperbarui data booking
func UpdateBooking(c echo.Context) error {
	id := c.Param("id")
	booking := new(models.Booking)

	// Bind data dari request body
	if err := c.Bind(booking); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
	}

	// Konversi tanggal string ke time.Time
	startRent, err := time.Parse("2006-01-02", booking.StartRent)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid start date format"})
	}

	endRent, err := time.Parse("2006-01-02", booking.EndRent)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid end date format"})
	}

	// Validasi bahwa startRent tidak setelah endRent
	if startRent.After(endRent) {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Start rent must be before end rent"})
	}

	// Periksa apakah booking dengan ID tersebut ada
	var existingBooking models.Booking
	checkQuery := `SELECT id FROM bookings WHERE id = $1`
	err = config.DB.Get(&existingBooking, checkQuery, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Booking not found"})
	}

	// Ambil harga sewa harian mobil
	var dailyRent float64
	query := `SELECT daily_rent FROM cars WHERE id = $1`
	err = config.DB.Get(&dailyRent, query, booking.CarID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch car data"})
	}

	// Hitung durasi dan total biaya
	duration := endRent.Sub(startRent).Hours() / 24
	if duration < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "End date must be at least one day after start date"})
	}
	booking.TotalCost = dailyRent * duration

	// Update data booking di database
	updateQuery := `
        UPDATE bookings 
        SET customer_id=$1, car_id=$2, start_rent=$3, end_rent=$4, total_cost=$5, finished=$6 
        WHERE id=$7`
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
