#!/bin/bash

# Start senderService
(cd senderService && go run main.go) &

# Start messageQueueService
(cd messageQueueService && go run main.go) &

# Start listenerService
(cd listenerService && go run main.go) &

# Wait for all background processes to finish (if needed)
wait
