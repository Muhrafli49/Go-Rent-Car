package main

import (
	"rental-mobil/config"
	"rental-mobil/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	// Inisialisasi database
	config.InitDB()

	// Inisialisasi Echo
	e := echo.New()

	// Daftarkan rute mobil dan pelanggan
	routes.RegisterCustomerRoutes(e)
	routes.RegisterCarRoutes(e)
	routes.BookingRoutes(e)

	// Jalankan server
	e.Logger.Fatal(e.Start(":5000"))
}
