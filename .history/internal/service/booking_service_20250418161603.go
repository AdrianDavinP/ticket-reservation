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
	"ticket-reservation/pb" // Import gRPC generated file

	"google.golang.org/grpc"
)

type BookingService struct {
	DB   *sql.DB
	Repo *repository.ConcertRepo
}

// Implement the GetConcerts method to fulfill the gRPC interface
func (s *BookingService) GetConcerts(ctx context.Context, req *pb.GetConcertsRequest) (*pb.GetConcertsResponse, error) {
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

	return &pb.GetConcertsResponse{Concerts: pbConcerts}, nil
}

// Implement the BookTicket method to fulfill the gRPC interface
func (s *BookingService) BookTicket(ctx context.Context, req *pb.BookTicketRequest) (*pb.BookTicketResponse, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return &pb.BookTicketResponse{Status: "FAILED"}, err
	}
	defer tx.Rollback()

	// Lock the concert by ID to prevent race condition
	concert, err := s.Repo.LockConcertByID(ctx, tx, int(req.ConcertId))
	if err != nil {
		return &pb.BookTicketResponse{Status: "FAILED"}, err
	}

	now := time.Now()
	if now.Before(concert.StartTime) || now.After(concert.EndTime) {
		return &pb.BookTicketResponse{Status: "FAILED"}, errors.New("booking not in time window")
	}

	if concert.AvailableTickets < int(req.Quantity) {
		return &pb.BookTicketResponse{Status: "FAILED"}, errors.New("not enough tickets available")
	}

	// Create booking
	booking := &model.Booking{
		ConcertID: concert.ID,
		UserID:    int(req.UserId),
		Quantity:  int(req.Quantity),
		BookedAt:  now,
	}
	if err := s.Repo.InsertBooking(ctx, tx, booking); err != nil {
		return &pb.BookTicketResponse{Status: "FAILED"}, err
	}

	if err := s.Repo.UpdateTicketStock(ctx, tx, concert.ID, int(req.Quantity)); err != nil {
		return &pb.BookTicketResponse{Status: "FAILED"}, err
	}

	if err := tx.Commit(); err != nil {
		return &pb.BookTicketResponse{Status: "FAILED"}, err
	}

	return &pb.BookTicketResponse{Status: "SUCCESS"}, nil
}

// StartGRPCServer will start the gRPC server
func StartGRPCServer(bookingService *BookingService) {
	// Membuka port untuk server gRPC
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	// Membuat server gRPC
	server := grpc.NewServer()

	// Mendaftarkan server ke gRPC
	pb.RegisterConcertServiceServer(server, bookingService)

	// Mulai server
	log.Println("Starting gRPC server on port 50051...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
