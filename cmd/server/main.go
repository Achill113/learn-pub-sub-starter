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

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
		return
	}
	defer channel.Close()

	err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
		return
	}

	fmt.Println("Successfully connected to RabbitMQ")

	gamelogic.PrintServerHelp()

Gameloop:
	for true {
		inputs := gamelogic.GetInput()

		if len(inputs) == 0 {
			continue
		}

		switch inputs[0] {
		case "pause":
			fmt.Println("Pausing the game...")
			err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
			if err != nil {
				log.Printf("Failed to publish pause message: %v", err)
			}
		case "resume":
			fmt.Println("Resuming the game...")
			err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
			if err != nil {
				log.Printf("Failed to publish resume message: %v", err)
			}
		case "quit":
			fmt.Println("Quitting the game...")
			break Gameloop
		default:
			fmt.Println("Unknown command. Please use 'pause', 'resume', or 'quit'.")
		}
	}

	// Wait for a signal to exit
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	fmt.Println("Shutting down gracefully...")
}
