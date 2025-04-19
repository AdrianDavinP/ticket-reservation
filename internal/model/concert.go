package model

import "time"

type Concert struct {
	ID               int       `json:"id"`
	NameConcert      string    `json:"name_concert"`
	AvailableTickets int       `json:"available_tickets"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
}
