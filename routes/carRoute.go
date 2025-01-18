package routes

import (
	"rental-mobil/controllers"
	"github.com/labstack/echo/v4"
)

// RegisterCarRoutes untuk menangani rute mobil
func RegisterCarRoutes(e *echo.Echo) {
	e.GET("/cars", controllers.GetAllCars)
	e.POST("/cars", controllers.CreateCar)
	e.PUT("/cars/:id", controllers.UpdateCar)
	e.DELETE("/cars/:id", controllers.DeleteCar)
}
