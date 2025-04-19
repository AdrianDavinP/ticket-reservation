package repository

import (
	"context"
	"database/sql"
	"errors"
	"ticket-reservation/internal/model"
)

type ConcertRepo struct {
	DB *sql.DB
}

// Get list of concerts
func (r *ConcertRepo) GetConcerts(ctx context.Context) ([]model.Concert, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id, name, available_tickets, start_time, end_time FROM concerts`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var concerts []model.Concert
	for rows.Next() {
		var c model.Concert
		if err := rows.Scan(&c.ID, &c.Name, &c.AvailableTickets, &c.StartTime, &c.EndTime); err != nil {
			return nil, err
		}
		concerts = append(concerts, c)
	}
	return concerts, nil
}

// Lock a concert row for update (inside transaction)
func (r *ConcertRepo) LockConcertByID(ctx context.Context, tx *sql.Tx, id int) (*model.Concert, error) {
	var c model.Concert
	err := tx.QueryRowContext(ctx, `SELECT id, name, available_tickets, start_time, end_time FROM concerts WHERE id = $1 FOR UPDATE`, id).
		Scan(&c.ID, &c.NameConcert, &c.AvailableTickets, &c.StartTime, &c.EndTime)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Insert a booking
func (r *ConcertRepo) InsertBooking(ctx context.Context, tx *sql.Tx, b *model.Booking) error {
	_, err := tx.ExecContext(ctx, `INSERT INTO bookings (concert_id, user_id, quantity, booked_at) VALUES ($1, $2, $3, $4)`,
		b.ConcertID, b.UserID, b.Quantity, b.BookedAt)
	return err
}

// Decrease ticket quantity
func (r *ConcertRepo) UpdateTicketStock(ctx context.Context, tx *sql.Tx, concertID int, quantity int) error {
	res, err := tx.ExecContext(ctx, `UPDATE concerts SET available_tickets = available_tickets - $1 WHERE id = $2`, quantity, concertID)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("no rows updated")
	}
	return nil
}
