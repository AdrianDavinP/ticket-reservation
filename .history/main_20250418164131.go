package main

import (
	"database/sql"
	"log"

	"ticket-reservation/internal/handler"
	"ticket-reservation/internal/repository"
	"ticket-reservation/internal/server"
	"ticket-reservation/internal/service"

	_ "github.com/lib/pq"
)

func main() {
	// Set up the database connection
	db, err := sql.Open("postgres", "your-postgres-connection-string")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Set up repository and service layer
	repo := repository.NewConcertRepo(db)
	svc := service.NewBookingService(db, repo)

	// Create gRPC handler that implements pb.ConcertServiceServer
	grpcHandler := handler.NewGrpcHandler(svc)

	// Start gRPC server
	server.StartGRPCServer(grpcHandler)
}
