package models

type Driver struct {
	ID          int     `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	NIK         string  `json:"nik" db:"nik"`
	PhoneNumber string  `json:"phone_number" db:"phone_number"`
	DailyCost   float64 `json:"daily_cost" db:"daily_cost"` // Biaya supir per hari
}
