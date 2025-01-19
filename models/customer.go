package models

type Customer struct {
	ID           int    `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	NIK          string `json:"nik" db:"nik"`
	Phone        string `json:"phone_number" db:"phone"`
	MembershipID *int   `json:"membership_id" db:"membership_id"`
}
