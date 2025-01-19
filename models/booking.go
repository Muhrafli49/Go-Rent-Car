package models

type Booking struct {
	ID              int     `json:"id" db:"id"`
	CustomerID      int     `json:"customer_id" db:"customer_id"`
	CarID           int     `json:"car_id" db:"car_id"`
	StartRent       string  `json:"start_rent" db:"start_rent"` // Tanggal mulai sewa
	EndRent         string  `json:"end_rent" db:"end_rent"`     // Tanggal selesai sewa
	TotalCost       float64 `json:"total_cost" db:"total_cost"` // Total biaya sewa
	Finished        bool    `json:"finished" db:"finished"`     // Status selesai
	Discount        float64 `json:"discount" db:"discount"`     // Diskon yang diterapkan
	BookingTypeID   int     `json:"booking_type_id" db:"booking_type_id"` // ID jenis booking
	DriverID        int     `json:"driver_id" db:"driver_id"`         // ID supir
	TotalDriverCost float64 `json:"total_driver_cost" db:"total_driver_cost"` // Biaya supir
}
