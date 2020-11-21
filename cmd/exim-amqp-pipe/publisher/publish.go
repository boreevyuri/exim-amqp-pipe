package publisher

import (
	"github.com/boreevyuri/exim-amqp-pipe/cmd/exim-amqp-pipe/config"
	"github.com/boreevyuri/exim-amqp-pipe/cmd/exim-amqp-pipe/reader"
	"github.com/streadway/amqp"
	"log"
)

func PublishFiles(done chan<- bool, files chan reader.File, config config.AmqpConfig) {

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

	//fmt.Printf("Connection successful. Publisher waits for data...\n")

	for file := range files {
		//fmt.Printf("Incoming File %d bytes, name: %s\n", len(file.Data), file.Filename)
		err := ch.Publish(
			"",
			binding.QueueName,
			false,
			false,
			amqp.Publishing{
				ContentType:     file.ContentType,
				ContentEncoding: file.ContentEncoding,
				Body:            file.Data,
			},
		)
		if err != nil {
			failOnError(err, "Failed to publish message:")
		}
		//fmt.Printf("Data published\n")
	}

	//fmt.Printf("Incoming channel closed. Publisher exited\n")
	done <- true

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s, %s", msg, err)
	}
}
