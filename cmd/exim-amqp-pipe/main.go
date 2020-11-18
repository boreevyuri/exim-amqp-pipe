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

	//dataChan := make(chan []byte)
	//done := make(chan bool)

	conf.GetConf(confFile)
	queue := publisher.New(conf.Amqp)

	//reader.Parse(conf.Parse)
	dataChan := reader.Parse(conf.Parse)

	for _, data := range <-dataChan {
		queue.Publish(data)
	}
	//mail := reader.ReadStdin()
	//queue.Publish(mail)
}
