package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
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

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Failed during client welcome: %v", err)
		return
	}

	ch, _, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		fmt.Sprintf("pause.%s", username),
		routing.PauseKey,
		pubsub.TransientQueue,
	)

	if err != nil {
		log.Fatalf("Failed to declare and bind queue: %v", err)
		return
	}
	defer ch.Close()

	// Wait for a signal to exit
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	fmt.Println("Shutting down gracefully...")
}
