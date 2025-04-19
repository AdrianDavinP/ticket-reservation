package server

import (
	"log"
	"net"
	"ticket-reservation/internal/service"
	"ticket-reservation/pb"

	"google.golang.org/grpc"
)

func StartGRPCServer(bookingService *service.BookingService) {
	// Open a TCP listener on port 50051
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the BookingService with the gRPC server
	pb.RegisterConcertServiceServer(grpcServer, bookingService)

	// Start the gRPC server
	log.Println("Starting gRPC server on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
