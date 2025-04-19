package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"ticket-reservation/internal/model"
	"ticket-reservation/internal/repository"
	"ticket-reservation/pb"
)

type BookingService struct {
	DB   *sql.DB
	Repo *repository.ConcertRepo
}

// Implementasi GetConcerts sesuai dengan ConcertServiceServer interface yang dihasilkan oleh gRPC
func (s *BookingService) GetConcerts(ctx context.Context, req *pb.Empty) (*pb.ConcertList, error) {
	// Ambil data konser dari repository
	concerts, err := s.Repo.GetConcerts(ctx)
	if err != nil {
		return nil, err
	}

	var pbConcerts []*pb.Concert
	for _, concert := range concerts {
		pbConcerts = append(pbConcerts, &pb.Concert{
			Id:               int32(concert.ID),
			Name:             concert.Name,
			StartTime:        concert.StartTime.String(),
			EndTime:          concert.EndTime.String(),
			AvailableTickets: int32(concert.AvailableTickets),
		})
	}

	// Mengembalikan ConcertList
	return &pb.ConcertList{Concerts: pbConcerts}, nil
}

// Implementasi BookTicket sesuai dengan ConcertServiceServer interface yang dihasilkan oleh gRPC
func (s *BookingService) BookTicket(ctx context.Context, req *pb.BookRequest) (*pb.BookResponse, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return &pb.BookResponse{Status: "FAILED", Message: err.Error()}, err
	}
	defer tx.Rollback()

	// Lock konser berdasarkan ID untuk menghindari kondisi race
	concert, err := s.Repo.LockConcertByID(ctx, tx, int(req.ConcertId))
	if err != nil {
		return &pb.BookResponse{Status: "FAILED", Message: err.Error()}, err
	}

	// Validasi apakah booking dilakukan dalam window waktu yang benar
	now := time.Now()
	if now.Before(concert.StartTime) || now.After(concert.EndTime) {
		return &pb.BookResponse{Status: "FAILED", Message: "booking not in time window"}, errors.New("booking not in time window")
	}

	// Validasi ketersediaan tiket
	if concert.AvailableTickets < int(req.Quantity) {
		return &pb.BookResponse{Status: "FAILED", Message: "not enough tickets available"}, errors.New("not enough tickets available")
	}

	// Buat entri booking
	booking := &model.Booking{
		ConcertID: concert.ID,
		UserID:    int(req.UserId),
		Quantity:  int(req.Quantity),
		BookedAt:  now,
	}
	if err := s.Repo.InsertBooking(ctx, tx, booking); err != nil {
		return &pb.BookResponse{Status: "FAILED", Message: err.Error()}, err
	}

	// Update jumlah tiket yang tersedia setelah booking
	if err := s.Repo.UpdateTicketStock(ctx, tx, concert.ID, int(req.Quantity)); err != nil {
		return &pb.BookResponse{Status: "FAILED", Message: err.Error()}, err
	}

	// Commit transaksi
	if err := tx.Commit(); err != nil {
		return &pb.BookResponse{Status: "FAILED", Message: err.Error()}, err
	}

	return &pb.BookResponse{Status: "SUCCESS", Message: "Booking successful"}, nil
}
