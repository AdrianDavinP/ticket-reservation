package server

import (
	"log"
	"net"

	"ticket-reservation/pb"

	"google.golang.org/grpc"
)

func StartGRPCServer(grpcHandler pb.ConcertServiceServer) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterConcertServiceServer(grpcServer, grpcHandler)

	log.Println("🚀 gRPC server is running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
