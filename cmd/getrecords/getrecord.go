package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func main() {
	uri := "amqp://guest:guest@127.0.0.1:5672/exim"
	conn, err := amqp.Dial(uri)
	failOnError(err, "Failed to connect to RabbitMQ:")

	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open channel:")

	defer ch.Close()

	_, err = ch.QueueDeclare(
		"exim",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue:")

	files, err := ch.Consume(
		"exim",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for file := range files {
			log.Printf(" [x] %s", file.Body)
		}
	}()

	log.Printf(" [*] Waiting for file...")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s, %s", msg, err)
	}
}
