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

	donePublish := make(chan bool)
	//publishChan := make(chan string)
	//parseChan := make(chan string)

	parseChan := make(chan reader.File)
	publishChan := make(chan reader.File)

	go publisher.NewPublish(donePublish, publishChan, conf.Amqp)
	go reader.Parse(parseChan, conf.Parse)

	for data := range parseChan {
		publishChan <- data
	}
	close(publishChan)

	<-donePublish
}
