package main

import (
	"database/sql"
	"log"
	"ticket-reservation/internal/handler"
	"ticket-reservation/internal/service"
	"ticket-reservation/pb"

	_ "github.com/lib/pq" // or your database driver (e.g., "github.com/go-sql-driver/mysql")
)

func main() {
	// Set up the database connection
	db, err := sql.Open("postgres", "your-postgres-connection-string")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	bookingService := &service.BookingService{
		DB:   db,
		Repo: repo,
	}

	grpcHandler := &handler.GrpcHandler{
		Service: bookingService,
	}

	pb.RegisterConcertServiceServer(grpcServer, grpcHandler)
}
