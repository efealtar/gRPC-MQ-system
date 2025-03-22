# Multi-Service gRPC Showcase

This repository demonstrates a simple multi-service architecture using Go and gRPC. The project consists of three services that work together to simulate a payment processing workflow, and a shared module containing the Protocol Buffers definitions.

## Overview

- **senderService**:

  - Exposes an HTTP POST endpoint (e.g., `/send`).
  - Receives JSON with `amount` and `address`.
  - Uses a gRPC client (from the shared proto module) to forward the payment details to the messageQueueService.

- **messageQueueService**:

  - Acts as a gRPC server listening on port `50051`.
  - Receives payment messages from senderService.
  - Forwards the message to the listenerService.

- **listenerService**:

  - Acts as a gRPC server listening on port `50052`.
  - Validates the payment by checking the provided `amount` against an expected amount for a given `address`.
  - Returns a success or error message accordingly.

- **protos**:
  - Contains the shared Protocol Buffers definitions (in `sender/sender.proto`) and the generated Go code.
  - Used by all services for inter-service communication.

## Directory Structure

.
├── listenerService
│ ├── go.mod
│ ├── go.sum
│ └── main.go
├── messageQueueService
│ ├── go.mod
│ ├── go.sum
│ └── main.go
├── protos
│ ├── go.mod
│ ├── go.sum
│ └── sender
│ ├── sender.proto
│ ├── sender.pb.go
│ └── sender_grpc.pb.go
├── run_all.sh
└── senderService
├── go.mod
├── go.sum
└── main.go

## Prerequisites

- Go (version 1.18 or later recommended)
- Protocol Buffers Compiler (protoc) with the Go plugins installed
- (Optional) Git for cloning the repository

## Setup

### 1. Clone the Repository

Clone this repository to your local machine:

    git clone https://github.com/yourusername/your-repository.git
    cd your-repository

### 2. Generate gRPC Code (if needed)

If you modify the Protocol Buffers definition (`sender.proto`), regenerate the Go code as follows:

1. Navigate to the `protos/sender` directory:

   cd protos/sender

2. Run the Protocol Buffers compiler with the source-relative option to avoid deep folder nesting:

   protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. sender.proto

3. Return to the root directory:

   cd ../../

### 3. Update Go Modules for Each Service

Each service (senderService, messageQueueService, listenerService) imports the shared proto package from `github.com/efealtar/protos/sender`. Since this module is local and not hosted remotely, ensure you have a replace directive in each service's `go.mod`. For example, in **senderService/go.mod** add:

    replace github.com/efealtar/protos => ../protos

Then run the following in each service's directory:

    go mod tidy

Repeat for messageQueueService and listenerService.

## How It Works

1. **senderService**

   - Receives an HTTP POST request at `/send` with a JSON payload containing `amount` and `address`.
   - Creates a gRPC client using the shared proto definitions.
   - Forwards the payment data to messageQueueService via gRPC.

2. **messageQueueService**

   - Listens on port `50051` for incoming gRPC requests.
   - Logs the received message.
   - Forwards the message to listenerService (assumed to be running on port `50052`).

3. **listenerService**
   - Listens on port `50052` for incoming gRPC requests.
   - Validates the payment against a predefined list of addresses and expected amounts.
   - Returns a status message:
     - If the amount is greater than the expected value, it returns a success message.
     - Otherwise, it returns an error message indicating an invalid payment.

## Running the Services

### Using the Provided Shell Script

A shell script (`run_all.sh`) is provided at the repository root to start all three services concurrently. This is useful for local development and testing.

1. **Make the Script Executable**

   chmod +x run_all.sh

2. **Run the Script**

   ./run_all.sh

The script changes directories into each service folder and runs `go run main.go` in the background for each service.

### Running Services Manually

Alternatively, you can start each service in separate terminal windows:

- **senderService:**

       cd senderService
       go run main.go

- **messageQueueService:**

       cd messageQueueService
       go run main.go

- **listenerService:**

       cd listenerService
       go run main.go

## Testing the Setup

Once all services are running, you can test the end-to-end workflow:

1.  **Send a POST Request to senderService**

    Use `curl` (or a tool like Postman) to send a payment request. For example:

        curl -X POST -H "Content-Type: application/json" \
        -d '{"amount": 100.5, "address": "abc123"}' \
        http://localhost:8080/send

2.  **Expected Flow:**

    - **senderService** receives the POST request and forwards the data via gRPC.
    - **messageQueueService** logs the received message and forwards it to listenerService.
    - **listenerService** validates the payment:
      - If `amount` > expected value for `address` (e.g., expected is 50.0 for "abc123"), it returns a success message:
        `payment is valid for address: abc123`
      - Otherwise, it returns an error message:
        `no address found or amount is less than expected skipping for address: abc123`
    - The final response is relayed back to the HTTP client.

## Additional Notes

- **Logging:** Each service logs its operations to the console. Check the terminal output for debugging.
- **Error Handling:** If a service (e.g., listenerService) is not running, the calling service will log an error. Ensure all services are running for full end-to-end functionality.
- **Development Workflow:**
  - Modify the `protos/sender/sender.proto` file as needed.
  - Regenerate the Go code for protos using the provided `protoc` command.
  - Update and run `go mod tidy` in each service directory if module paths change.
- **Production Considerations:**
  - For a production environment, consider containerizing each service and using Docker Compose or Kubernetes for orchestration.

## Conclusion

This project demonstrates a basic microservice architecture with Go and gRPC, using a shared Protocol Buffers module for inter-service communication. Use this repository as a foundation for building more complex distributed systems.
