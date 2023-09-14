package main

import (
	"exim-amqp-pipe/cmd/exim-amqp-pipe/config"
	"exim-amqp-pipe/cmd/exim-amqp-pipe/publisher"
	"exim-amqp-pipe/cmd/exim-amqp-pipe/reader"
	"flag"
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
