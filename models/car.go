package models


type Car struct {
	ID        int     `json:"id" db:"id"`
	Name      string  `json:"name" db:"name"`
	Stock     int     `json:"stock" db:"stock"`
	DailyRent float64 `json:"daily_rent" db:"daily_rent"`
}
