package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net"
	"time"

	"ticket-reservation/internal/model"
	"ticket-reservation/internal/repository"
	"ticket-reservation/pb"

	"google.golang.org/grpc"
)

type BookingService struct {
	DB   *sql.DB
	Repo *repository.ConcertRepo
}

func (s *BookingService) GetConcerts(ctx context.Context) ([]model.Concert, error) {
	return s.Repo.GetConcerts(ctx)
}

func (s *BookingService) BookTicket(ctx context.Context, concertID, userID, quantity int) (string, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return "FAILED", err
	}
	defer tx.Rollback()

	concert, err := s.Repo.LockConcertByID(ctx, tx, concertID)
	if err != nil {
		return "FAILED", err
	}

	now := time.Now()
	if now.Before(concert.StartTime) || now.After(concert.EndTime) {
		return "FAILED", errors.New("booking not in time window")
	}

	if concert.AvailableTickets < quantity {
		return "FAILED", errors.New("not enough tickets available")
	}

	booking := &model.Booking{
		ConcertID: concert.ID,
		UserID:    userID,
		Quantity:  quantity,
		BookedAt:  now,
	}
	if err := s.Repo.InsertBooking(ctx, tx, booking); err != nil {
		return "FAILED", err
	}

	if err := s.Repo.UpdateTicketStock(ctx, tx, concert.ID, quantity); err != nil {
		return "FAILED", err
	}

	if err := tx.Commit(); err != nil {
		return "FAILED", err
	}

	return "SUCCESS", nil
}

func StartGRPCServer() {
	// Membuka port untuk server gRPC
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	// Membuat server gRPC
	server := grpc.NewServer()

	// Mendaftarkan server ke gRPC
	pb.RegisterConcertServiceServer(server, BookingService)

	// Mulai server
	log.Println("Starting gRPC server on port 50051...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
