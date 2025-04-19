package service_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"ticket-reservation/internal/repository"
	"ticket-reservation/internal/service"
)

var db *sql.DB

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	// Membuat connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	code := m.Run()
	db.Close()
	os.Exit(code)
}

func setupService(t *testing.T) *service.BookingService {
	t.Helper()

	// Reset concert ID 1
	_, err := db.Exec(`UPDATE concerts SET available_tickets = 10 WHERE id = 1`)
	require.NoError(t, err)

	// Hapus semua bookings
	_, err = db.Exec(`DELETE FROM bookings`)
	require.NoError(t, err)

	repo := &repository.ConcertRepo{DB: db}
	return &service.BookingService{
		DB:   db,
		Repo: repo,
	}
}

func TestBookTicket_Success(t *testing.T) {
	svc := setupService(t)

	status, err := svc.BookTicket(context.Background(), 1, 2, 1)
	require.NoError(t, err)
	require.Equal(t, "SUCCESS", status)
}

func TestBookTicket_Fail_NotEnoughTickets(t *testing.T) {
	svc := setupService(t)

	status, err := svc.BookTicket(context.Background(), 1, 3, 20)
	require.Error(t, err)
	require.Equal(t, "FAILED", status)
}

func TestBookTicket_Concurrency(t *testing.T) {
	svc := setupService(t)

	var wg sync.WaitGroup
	var mu sync.Mutex
	successCount := 0
	failCount := 0
	totalUsers := 20

	maxCon, err := strconv.Atoi(os.Getenv("TOTAL_USERS"))
	if err != nil {
		log.Fatalf("Error env DB_MAX_CONN : %v", err)
	}

	for i := 0; i < totalUsers; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()

			status, err := svc.BookTicket(context.Background(), 1, 4, 1)

			mu.Lock()
			defer mu.Unlock()
			if err == nil && status == "SUCCESS" {
				successCount++
			} else {
				failCount++
			}
		}(i)
	}

	wg.Wait()
	require.Equal(t, 10, successCount)
	require.Equal(t, totalUsers-10, failCount)
}
