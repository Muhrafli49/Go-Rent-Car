package models

import (
	"time"
)

type Booking struct {
	ID         int       `db:"id" json:"id"`
	CustomerID int       `db:"customer_id" json:"customer_id"`
	CarID      int       `db:"car_id" json:"car_id"`
	StartRent  time.Time `db:"start_rent" json:"start_rent"`
	EndRent    time.Time `db:"end_rent" json:"end_rent"`
	TotalCost  float64   `db:"total_cost" json:"total_cost"`
	Finished   bool      `db:"finished" json:"finished"`
}
