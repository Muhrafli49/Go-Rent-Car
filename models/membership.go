package models

type Membership struct {
	ID       int     `json:"id" db:"id"`
	Name     string  `json:"membership_name" db:"name"` // Nama keanggotaan
	Discount float64 `json:"discount" db:"discount"`   // Diskon yang diterapkan
}
