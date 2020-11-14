package main

import (
	"github.com/boreevyuri/exim-amqp-pipe/cmd/exim-amqp-pipe/publisher"
	"github.com/boreevyuri/exim-amqp-pipe/cmd/exim-amqp-pipe/reader"
	"log"
)

const (
	defaultConfigFile = "/etc/exim-amqp-pipe.yaml"
)

func main() {
	//var (
	//	confFile string
	//)
	//
	//flag.StringVar(&confFile, "c", defaultConfigFile, "configuration file")
	//flag.Parse()

	mail, err := reader.ReadStdin()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	publisher.PublishMail(mail)

}
