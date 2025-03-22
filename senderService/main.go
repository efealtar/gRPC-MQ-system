package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	pb "github.com/efealtar/protos/sender"
	"google.golang.org/grpc"
)

// RequestPayload represents the JSON structure for the POST request.
type RequestPayload struct {
    Amount  float64 `json:"amount"`
    Address string  `json:"address"`
}

// sendHandler handles the POST request, parses the payload, and sends it via gRPC.
func sendHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
        return
    }

    // Decode the JSON payload.
    var payload RequestPayload
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Establish a gRPC connection to the second service.
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        http.Error(w, "Failed to connect to gRPC service", http.StatusInternalServerError)
        log.Printf("gRPC dial error: %v", err)
        return
    }
    defer conn.Close()

    // Create a new gRPC client.
    client := pb.NewMessageServiceClient(conn)

    // Set a timeout for the gRPC call.
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Prepare the gRPC request.
    req := &pb.MessageRequest{
        Amount:  payload.Amount,
        Address: payload.Address,
    }

    // Call the SendMessage RPC.
    res, err := client.SendMessage(ctx, req)
    if err != nil {
        http.Error(w, "Error calling gRPC service: "+err.Error(), http.StatusInternalServerError)
        log.Printf("gRPC call error: %v", err)
        return
    }

    // Return the gRPC response as JSON.
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(res); err != nil {
        log.Printf("Failed to encode response: %v", err)
    }
}

func main() {
    // Register the POST endpoint.
    http.HandleFunc("/send", sendHandler)

    log.Println("Sender Service HTTP server running on port 8080...")
    // Start the HTTP server.
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
