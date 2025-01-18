package routes

import (
    "rental-mobil/controllers"
    "github.com/labstack/echo/v4"
)

func BookingRoutes(e *echo.Echo) {
    e.GET("/bookings", controllers.GetAllBookings)
    e.POST("/bookings", controllers.CreateBooking)
    e.PUT("/bookings/:id", controllers.UpdateBooking)
    e.DELETE("/bookings/:id", controllers.DeleteBooking)
}
