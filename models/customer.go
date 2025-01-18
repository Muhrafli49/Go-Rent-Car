package models

type Customer struct {
	ID    int    `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	NIK   string `db:"nik" json:"nik"`
	Phone string `db:"phone" json:"phone"`
}
