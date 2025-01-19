package models

type DriverIncentive struct {
	ID        int     `json:"id" db:"id"`
	BookingID int     `json:"booking_id" db:"booking_id"`
	Incentive float64 `json:"incentive" db:"incentive"`
}
