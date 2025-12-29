package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	connectionString := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		return
	}
	defer conn.Close()

	fmt.Println("Successfully connected to RabbitMQ")

	// Wait for a signal to exit
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	fmt.Println("Shutting down gracefully...")
}
