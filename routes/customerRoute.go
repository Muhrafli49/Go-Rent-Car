package routes

import (
	"rental-mobil/controllers"

	"github.com/labstack/echo/v4"
)

// RegisterCustomerRoutes untuk menangani rute pelanggan
func RegisterCustomerRoutes(e *echo.Echo) {
	e.GET("/customers", controllers.GetAllCustomers)
	e.POST("/customers", controllers.CreateCustomer)
	e.PUT("/customers/:id", controllers.UpdateCustomer)
	e.DELETE("/customers/:id", controllers.DeleteCustomer)
}
