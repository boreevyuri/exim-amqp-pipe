package publisher

import (
	"github.com/boreevyuri/exim-amqp-pipe/cmd/exim-amqp-pipe/config"
	"github.com/boreevyuri/exim-amqp-pipe/cmd/exim-amqp-pipe/reader"
	"github.com/streadway/amqp"
	"log"
)

func PublishFiles(done chan<- bool, emails chan reader.Email, config config.AMQPConfig) {
	uri, binding := config.URI, config.QueueBind
	conn, err := amqp.Dial(uri)
	failOnError(err, "Failed to connect to RabbitMQ:")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open channel:")

	_, err = ch.QueueDeclare(
		binding.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue:")

	for email := range emails {
		h := map[string]interface{}{
			"Rcpt-to": email.Rcpt,
			"From":    email.Sender,
		}
		for _, at := range email.Attachments {
			h["Content-Disposition"] = at.ContentDisposition
			err := ch.Publish(
				"",
				binding.QueueName,
				false,
				false,
				amqp.Publishing{
					Headers:         h,
					ContentType:     at.ContentType,
					ContentEncoding: at.ContentEncoding,
					Body:            at.Data,
				},
			)
			failOnError(err, "Failed to publish message:")
		}
	}

	err = ch.Close()
	failOnError(err, "amqp channel already closed")

	err = conn.Close()
	failOnError(err, "amqp connection already closed")

	done <- true
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s, %s", msg, err)
	}
}
