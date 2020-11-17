package main

import (
	"github.com/boreevyuri/exim-amqp-pipe/cmd/exim-amqp-pipe/publisher"
	"github.com/boreevyuri/exim-amqp-pipe/cmd/exim-amqp-pipe/reader"
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

	mail := reader.ReadStdin()
	publisher.PublishMail(mail)

}
