package model

import "time"

type Booking struct {
	ID        int       `json:"id"`
	ConcertID int       `json:"concert_id"`
	UserID    int       `json:"user_id"`
	Quantity  int       `json:"quantity"`
	BookedAt  time.Time `json:"booked_at"`
}
