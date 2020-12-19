package publisher

import (
	"github.com/boreevyuri/exim-amqp-pipe/cmd/exim-amqp-pipe/config"
	"github.com/boreevyuri/exim-amqp-pipe/cmd/exim-amqp-pipe/reader"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
	"os"
)

func Publisher(done chan<- bool, emails chan reader.Email, config config.AMQPConfig) {
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
			amqpClient.Push(atData)
		}
	}

	err := amqpClient.Close()
	if err != nil {
		logger.Printf("unable to close rabbitMQ connection: %v", err)
	}

	done <- true
}
