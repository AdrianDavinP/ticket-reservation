package repository_test

import (
	"context"
	"testing"
	"ticket-reservation/internal/model"
	"ticket-reservation/internal/repository/mocks"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSearchAvailableConcerts(t *testing.T) {
	// Membuat instance dari MockConcertRepo
	mockRepo := new(mocks.ConcertRepository)

	// Setup mock behavior untuk SearchAvailableConcerts
	mockRepo.On("SearchAvailableConcerts", mock.Anything, "Test").Return([]model.Concert{
		{ID: 1, NameConcert: "Test Concer", AvailableTickets: 100},
	}, nil)

	// Test
	ctx := context.Background()
	concerts, err := mockRepo.SearchAvailableConcerts(ctx, "Test")

	// Verifikasi
	assert.NoError(t, err)
	assert.Len(t, concerts, 1)
	assert.Equal(t, concerts[0].NameConcert, "Test Concer")

	// Verify bahwa metode mock telah dipanggil dengan argumen yang sesuai
	mockRepo.AssertExpectations(t)
}

func TestLockConcertByID(t *testing.T) {
	// Membuat instance dari MockConcertRepo
	mockRepo := new(mocks.ConcertRepository)

	// Setup mock behavior untuk LockConcertByID
	mockRepo.On("LockConcertByID", mock.Anything, mock.Anything, 1).Return(&model.Concert{
		ID:               1,
		NameConcert:      "Test Concert",
		AvailableTickets: 100,
	}, nil)

	// Test
	ctx := context.Background()
	concert, err := mockRepo.LockConcertByID(ctx, nil, 1)

	// Verifikasi
	assert.NoError(t, err)
	assert.Equal(t, concert.NameConcert, "Test Concert")
	assert.Equal(t, concert.ID, 1)

	// Verify bahwa metode mock telah dipanggil dengan argumen yang sesuai
	mockRepo.AssertExpectations(t)
}

func TestInsertBookingAndUpdateStock(t *testing.T) {
	// Membuat instance dari MockConcertRepo
	mockRepo := new(mocks.ConcertRepository)

	// Setup mock behavior untuk InsertBooking dan UpdateTicketStock
	mockRepo.On("InsertBooking", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Test
	ctx := context.Background()
	booking := &model.Booking{
		ConcertID: 1,
		UserID:    123,
		Quantity:  2,
		BookedAt:  time.Now(),
	}

	// Panggil method dengan mock repo
	err := mockRepo.InsertBooking(ctx, nil, booking)
	assert.NoError(t, err)

	err = mockRepo.UpdateTicketStock(ctx, nil, 1, 2)
	assert.NoError(t, err)

	// Verify bahwa metode mock telah dipanggil
	mockRepo.AssertExpectations(t)
}

func TestUpdateStock(t *testing.T) {
	mockRepo := new(mocks.ConcertRepository)
	mockRepo.On("UpdateTicketStock", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := context.Background()
}
