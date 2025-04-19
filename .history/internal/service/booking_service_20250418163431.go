package service

import (
	"context"
	"database/sql"
	"errors"
	"ticket-reservation/internal/model"
	"ticket-reservation/internal/repository"
	"ticket-reservation/pb" // Import gRPC generated files
	"time"
)

type BookingService struct {
	pb.UnimplementedConcertServiceServer
	DB   *sql.DB
	Repo *repository.ConcertRepo
}

// Implement the GetConcerts method to fulfill the gRPC interface
func (s *BookingService) GetConcerts(ctx context.Context, req *pb.Empty) (*pb.ConcertList, error) {
	concerts, err := s.Repo.GetConcerts(ctx)
	if err != nil {
		return nil, err
	}

	var pbConcerts []*pb.Concert
	for _, concert := range concerts {
		pbConcerts = append(pbConcerts, &pb.Concert{
			Id:               int32(concert.ID),
			Name:             concert.NameConcert,
			StartTime:        concert.StartTime.String(),
			EndTime:          concert.EndTime.String(),
			AvailableTickets: int32(concert.AvailableTickets),
		})
	}

	return &pb.ConcertList{Concerts: pbConcerts}, nil
}

// Implement the BookTicket method to fulfill the gRPC interface
func (s *BookingService) BookTicket(ctx context.Context, req *pb.BookRequest) (*pb.BookResponse, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return &pb.BookResponse{Status: "FAILED", Message: err.Error()}, err
	}
	defer tx.Rollback()

	// Lock the concert by ID to prevent race condition
	concert, err := s.Repo.LockConcertByID(ctx, tx, int(req.ConcertId))
	if err != nil {
		return &pb.BookResponse{Status: "FAILED", Message: err.Error()}, err
	}

	now := time.Now()
	if now.Before(concert.StartTime) || now.After(concert.EndTime) {
		return &pb.BookResponse{Status: "FAILED", Message: "booking not in time window"}, errors.New("booking not in time window")
	}

	if concert.AvailableTickets < int(req.Quantity) {
		return &pb.BookResponse{Status: "FAILED", Message: "not enough tickets available"}, errors.New("not enough tickets available")
	}

	// Create booking
	booking := &model.Booking{
		ConcertID: concert.ID,
		UserID:    int(req.UserId),
		Quantity:  int(req.Quantity),
		BookedAt:  now,
	}
	if err := s.Repo.InsertBooking(ctx, tx, booking); err != nil {
		return &pb.BookResponse{Status: "FAILED", Message: err.Error()}, err
	}

	if err := s.Repo.UpdateTicketStock(ctx, tx, concert.ID, int(req.Quantity)); err != nil {
		return &pb.BookResponse{Status: "FAILED", Message: err.Error()}, err
	}

	if err := tx.Commit(); err != nil {
		return &pb.BookResponse{Status: "FAILED", Message: err.Error()}, err
	}

	return &pb.BookResponse{Status: "SUCCESS", Message: "Booking successful"}, nil
}
