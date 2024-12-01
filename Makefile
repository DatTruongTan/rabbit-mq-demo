# Makefile

# Declare phony targets
.PHONY: producer consumer clean

# Run the producer
producer:
	go run producer/main.go

# Run the consumer
consumer:
	go run consumer/main.go

# Clean up (optional)
clean:
	echo "No cleanup necessary for now"
