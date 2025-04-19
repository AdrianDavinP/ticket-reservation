package handler

import (
	"context"
	"log"
	"ticket-reservation/internal/service"
	"ticket-reservation/pb"
)

type GrpcHandler struct {
	pb.UnimplementedConcertServiceServer
	Service *service.BookingService
}

func NewGrpcHandler(service *service.BookingService) *GrpcHandler {
	return &GrpcHandler{Service: service}
}

func (h *GrpcHandler) GetConcerts(ctx context.Context, _ *pb.Empty) (*pb.ConcertList, error) {
	concerts, err := h.Service.GetConcerts(ctx)
	if err != nil {
		return nil, err
	}

	var pbConcerts []*pb.Concert
	for _, c := range concerts {
		pbConcerts = append(pbConcerts, &pb.Concert{
			Id:               int32(c.ID),
			Name:             c.NameConcert,
			AvailableTickets: int32(c.AvailableTickets),
			StartTime:        c.StartTime.Format("2006-01-02T15:04:05"),
			EndTime:          c.EndTime.Format("2006-01-02T15:04:05"),
		})
	}

	return &pb.ConcertList{Concerts: pbConcerts}, nil
}

func (h *GrpcHandler) BookTicket(ctx context.Context, req *pb.BookRequest) (*pb.BookResponse, error) {
	status, err := h.Service.BookTicket(ctx, int(req.ConcertId), int(req.UserId), int(req.Quantity))
	if err != nil {
		return &pb.BookResponse{
			Status:  status,
			Message: err.Error(),
		}, nil
	}

	return &pb.BookResponse{
		Status:  status,
		Message: "Ticket booked successfully",
	}, nil
}

func (h *GrpcHandler) SearchConcerts(ctx context.Context, req *pb.SearchRequest) (*pb.ConcertList, error) {
	// Menambahkan pencarian konser berdasarkan nama
	concerts, err := h.Service.GetConcerts(ctx, req.Name) // Mengirim nama ke service
	if err != nil {
		log.Printf("Error searching concerts: %v", err)
		return nil, err
	}

	var pbConcerts []*pb.Concert
	for _, c := range concerts {
		pbConcerts = append(pbConcerts, &pb.Concert{
			Id:               int32(c.ID),
			Name:             c.NameConcert,
			AvailableTickets: int32(c.AvailableTickets),
			StartTime:        c.StartTime.Format("2006-01-02T15:04:05"),
			EndTime:          c.EndTime.Format("2006-01-02T15:04:05"),
		})
	}

	return &pb.ConcertList{Concerts: pbConcerts}, nil
}
