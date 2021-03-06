package main

import (
	"flag"
	"github.com/boreevyuri/exim-amqp-pipe/cmd/exim-amqp-pipe/config"
	"github.com/boreevyuri/exim-amqp-pipe/cmd/exim-amqp-pipe/publisher"
	"github.com/boreevyuri/exim-amqp-pipe/cmd/exim-amqp-pipe/reader"
)

const (
	defaultConfigFile = "/etc/exim-amqp-pipe.yaml"
)

func main() {
	var (
		confFile string
		conf     config.Conf
	)

	flag.StringVar(&confFile, "c", defaultConfigFile, "configuration file")
	flag.Parse()

	conf.GetConf(confFile)

	emlFiles := flag.Args()

	donePublish := make(chan bool)
	emailChan := make(chan reader.Email)

	go publisher.Publisher(donePublish, emailChan, conf.AMQP)
	go reader.ReadInput(emailChan, emlFiles, conf.Parse)

	<-donePublish
}
