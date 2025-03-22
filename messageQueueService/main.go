package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "github.com/efealtar/protos/sender"
	"google.golang.org/grpc"
)

// server implements the MessageServiceServer interface.
type server struct {
    pb.UnimplementedMessageServiceServer
}

// SendMessage receives a message from senderService, then forwards it to listenerService.
func (s *server) SendMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
    log.Printf("Received message: amount=%.2f, address=%s", req.Amount, req.Address)

    // (Optional) Here you could add the message to a queue for asynchronous processing.
    // For this showcase, we forward the message directly to the listenerService.

    // Connect to the listenerService gRPC server (assumed to be on localhost:50052).
    conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
    if err != nil {
        log.Printf("Failed to connect to listenerService: %v", err)
        return &pb.MessageResponse{Status: "Error connecting to listenerService"}, err
    }
    defer conn.Close()

    listenerClient := pb.NewMessageServiceClient(conn)
    // Set a timeout for the gRPC call.
    ctxForward, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Forward the message to listenerService.
    res, err := listenerClient.SendMessage(ctxForward, req)
    if err != nil {
        log.Printf("Error forwarding message to listenerService: %v", err)
        return &pb.MessageResponse{Status: "Error forwarding message"}, err
    }

    // Return the response from listenerService.
    return res, nil
}

func main() {
    // Listen for gRPC connections on port 50051.
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterMessageServiceServer(grpcServer, &server{})

    log.Println("MessageQueueService gRPC server is running on port 50051...")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
