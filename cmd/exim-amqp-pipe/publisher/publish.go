package publisher

import (
	"github.com/boreevyuri/exim-amqp-pipe/cmd/exim-amqp-pipe/config"
	"github.com/streadway/amqp"
	"log"
)

type commandAction int

const (
	publish commandAction = iota
	finish
)

type commandData struct {
	action commandAction
	//publishing  []byte
	publishing amqp.Publishing
	result     chan<- bool
}

type processAmqp chan commandData

func (pa processAmqp) run(amqpConf config.AmqpConfig) {
	uri, binding := amqpConf.URI, amqpConf.QueueBind

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

	//exchange := binding.Exchange

	//for command := range pa {
	command := <-pa
	switch command.action {
	case publish:
		log.Printf("Got Publish signal")
		err = publishData(ch, binding.QueueName, command.publishing)
		failOnError(err, "Failed to publish a message:")
		command.result <- true
	case finish:
		log.Printf("Got Finish signal")
		close(pa)
	}
	//}
}

func (pa processAmqp) Publish(value []byte) bool {
	reply := make(chan bool)
	pa <- commandData{
		action: publish,
		publishing: amqp.Publishing{
			ContentType: "text/plain",
			Body:        value,
		},
		result: reply,
	}
	return <-reply
}

func (pa processAmqp) Finish() bool {
	reply := make(chan bool)
	pa <- commandData{
		action: finish,
		//publishing:  nil,
		result: reply,
	}
	return <-reply
}

func publishData(ch *amqp.Channel, name string, value amqp.Publishing) error {
	err := ch.Publish(
		"",
		name,
		false,
		false,
		value,
	)
	return err
}

type ProcessAmqp interface {
	Publish(value []byte) bool
	Finish() bool
}

func New(confFile string) ProcessAmqp {
	var conf config.Conf
	conf.GetConf(confFile)
	amqpPipe := make(processAmqp)
	go amqpPipe.run(conf.Amqp)
	return amqpPipe
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s, %s", msg, err)
	}
}
