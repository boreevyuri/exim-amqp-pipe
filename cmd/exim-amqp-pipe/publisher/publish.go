package publisher

import (
	"fmt"
	"github.com/boreevyuri/exim-amqp-pipe/cmd/exim-amqp-pipe/config"
	"github.com/streadway/amqp"
	"log"
)

func NewPublish(done chan<- bool, incoming chan string, config config.AmqpConfig) {

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

	for data := range incoming {
		fmt.Printf("Incoming Data %d bytes\n", len(data))
		value := []byte(data)
		err := ch.Publish(
			"",
			binding.QueueName,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        value,
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
