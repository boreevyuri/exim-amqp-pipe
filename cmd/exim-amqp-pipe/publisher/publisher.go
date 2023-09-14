package publisher

import (
	"exim-amqp-pipe/cmd/exim-amqp-pipe/config"
	"exim-amqp-pipe/cmd/exim-amqp-pipe/reader"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"os"
)

func Publisher(done chan<- bool, emails <-chan reader.Email, config config.AMQPConfig) {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	uri, queue := config.URI, config.QueueBind
	amqpClient := New(queue.QueueName, uri, logger)

	for email := range emails {
		headers := map[string]interface{}{
			"Rcpt-to": email.Rcpt,
			"From":    email.Sender,
		}
		for _, attachment := range email.Attachments {
			headers["Content-Disposition"] = attachment.ContentDisposition
			atData := amqp.Publishing{
				Headers:         headers,
				ContentType:     attachment.ContentType,
				ContentEncoding: attachment.ContentEncoding,
				Body:            attachment.Data,
			}
			err := amqpClient.Push(atData)
			if err != nil {
				amqpClient.logger.Printf("Unable to push: %v", err)
			}
		}
	}

	err := amqpClient.Close()
	if err != nil {
		logger.Printf("Unable to close rabbitMQ connection: %v", err)
	}

	done <- true
}
