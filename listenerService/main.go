package main

import (
	"context"
	"log"
	"net"

	pb "github.com/efealtar/protos/sender"
	"google.golang.org/grpc"
)

// A sample list of valid addresses and their expected amounts.
var validPayments = map[string]float64{
    "abc123": 50.0,
    "xyz789": 75.0,
}

// server implements the MessageServiceServer interface.
type server struct {
    pb.UnimplementedMessageServiceServer
}

// SendMessage checks the received amount against the expected amount for the address.
func (s *server) SendMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
    expectedAmount, ok := validPayments[req.Address]
    if !ok {
        log.Printf("Address %s not found", req.Address)
        return &pb.MessageResponse{Status: "no address found or amount is less than expected skipping for address: " + req.Address}, nil
    }
    if req.Amount > expectedAmount {
        log.Printf("Payment valid for address %s: received amount %.2f is greater than expected %.2f", req.Address, req.Amount, expectedAmount)
        return &pb.MessageResponse{Status: "payment is valid for address: " + req.Address}, nil
    }
    log.Printf("Amount %.2f is less than expected %.2f for address %s", req.Amount, expectedAmount, req.Address)
    return &pb.MessageResponse{Status: "no address found or amount is less than expected skipping for address: " + req.Address}, nil
}

func main() {
    // Listen for gRPC connections on port 50052.
    lis, err := net.Listen("tcp", ":50052")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterMessageServiceServer(grpcServer, &server{})

    log.Println("ListenerService gRPC server is running on port 50052...")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
