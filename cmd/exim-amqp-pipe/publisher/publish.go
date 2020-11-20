package publisher

import (
	"fmt"
	"github.com/boreevyuri/exim-amqp-pipe/cmd/exim-amqp-pipe/config"
	"github.com/boreevyuri/exim-amqp-pipe/cmd/exim-amqp-pipe/reader"
	"github.com/streadway/amqp"
	"log"
)

func NewPublish(done chan<- bool, incoming chan reader.File, config config.AmqpConfig) {

	uri, binding := config.URI, config.QueueBind
	conn, err := amqp.Dial(uri)
	failOnError(err, "Failed to connect to RabbitMQ:")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open channel:")
	defer ch.Close()

	_, err = ch.QueueDeclare(
		binding.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue:")

	fmt.Printf("Connection successful. Publisher waits for data...\n")

	for file := range incoming {
		//fmt.Printf("Incoming Data %d bytes\n", len(data))
		fmt.Printf("Incoming File %d bytes, name: %s\n", len(file.Data), file.Filename)
		//value := []byte(data)
		err := ch.Publish(
			"",
			binding.QueueName,
			false,
			false,
			amqp.Publishing{
				ContentType:     file.ContentType,
				ContentEncoding: file.ContentEncoding,
				//ContentType: "application/octet-stream",
				Body: file.Data,
			},
		)
		if err != nil {
			failOnError(err, "Failed to publish message:")
		}
		fmt.Printf("Data published\n")
	}

	fmt.Printf("Incoming channel closed. Publisher exited\n")
	done <- true

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s, %s", msg, err)
	}
}
