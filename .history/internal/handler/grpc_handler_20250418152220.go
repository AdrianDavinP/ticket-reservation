package handler

import (
	"context"
	"ticket-reservation/internal/service"
	"time"
)

type ConcertHandler struct {
	pb.UnimplementedConcertServiceServer
	Service *service.BookingService
}

func (h *ConcertHandler) GetConcerts(ctx context.Context, _ *pb.Empty) (*pb.ConcertList, error) {
	concerts, err := h.Service.GetConcerts(ctx)
	if err != nil {
		return nil, err
	}
	var result []*pb.Concert
	for _, c := range concerts {
		result = append(result, &pb.Concert{
			Id:               int32(c.ID),
			Name:             c.Name,
			AvailableTickets: int32(c.AvailableTickets),
			StartTime:        c.StartTime.Format(time.RFC3339),
			EndTime:          c.EndTime.Format(time.RFC3339),
		})
	}
	return &pb.ConcertList{Concerts: result}, nil
}

func (h *ConcertHandler) BookTicket(ctx context.Context, req *pb.BookRequest) (*pb.BookResponse, error) {
	status, err := h.Service.BookTicket(ctx, int(req.ConcertId), int(req.UserId), int(req.Quantity))
	if err != nil {
		return &pb.BookResponse{Status: status, Message: err.Error()}, nil
	}
	return &pb.BookResponse{Status: status, Message: "Ticket booked successfully"}, nil
}
