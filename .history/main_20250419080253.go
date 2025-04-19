package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"ticket-reservation/internal/handler"
	"ticket-reservation/internal/repository"
	"ticket-reservation/internal/server"
	"ticket-reservation/internal/service"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// load config dari .env
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	// Membuat connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

	// Membuka koneksi ke database PostgreSQL
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	maxCon, err := strconv.Atoi(os.Getenv("DB_MAX_CONN"))
	if err != nil {
		log.Fatalf("Error env DB_MAX_CONN : %v", err)
	}

	maxCon, err := strconv.Atoi(os.Getenv("DB_MAX_IDLE"))
	if err != nil {
		log.Fatalf("Error env DB_MAX_CONN : %v", err)
	}

	maxCon, err := strconv.Atoi(os.Getenv("DB_MAX_CONN"))
	if err != nil {
		log.Fatalf("Error env DB_MAX_CONN : %v", err)
	}

	=5
DB_CONN_LIFETIME

	db.SetMaxOpenConns(maxCon)
	db.SetMaxIdleConns(os.Getenv("DB_USER"))
	db.SetConnMaxLifetime(30 * time.Minute)

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
	fmt.Println("Successfully connected to the database")

	// Set up repository and service layer
	repo := repository.NewConcertRepo(db)
	svc := service.NewBookingService(db, repo)

	// Create gRPC handler that implements pb.ConcertServiceServer
	grpcHandler := handler.NewGrpcHandler(svc)

	// Start gRPC server
	server.StartGRPCServer(grpcHandler)
}
