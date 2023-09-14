package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type ExchangeType string

// const (
//	DirectExchangeType ExchangeType = "direct"
//	FanoutExchangeType ExchangeType = "fanout"
//	TopicExchangeType  ExchangeType = "topic"
// )

type Conf struct {
	AMQP  AMQPConfig  `yaml:"amqp"`
	Parse ParseConfig `yaml:"publish"`
}

type AMQPConfig struct {
	URI       string      `yaml:"uri"`
	QueueBind QueueConfig `yaml:"bindings"`
}

type QueueConfig struct {
	// Name string `yaml:"queue"`
	Exchange  string       `yaml:"exchange"`
	QueueName string       `yaml:"queue"`
	Type      ExchangeType `yaml:"type"`
	Routing   string       `yaml:"routing"`
}

type ParseConfig struct {
	AttachmentsOnly   bool `yaml:"attachments_only"`
	WithEmbeddedFiles bool `yaml:"embedded_files"`
}

func (c *Conf) GetConf(filename string) *Conf {
	config := readConfigFile(filename)

	err := yaml.Unmarshal(config, c)
	failOnError(err, "Unable to parse config file:")

	return c
}

func readConfigFile(filename string) []byte {
	if len(filename) == 0 {
		log.Fatalf("No config file specified")
	}

	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Config file not found")
	}

	return fileBytes
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s, %s", msg, err)
	}
}
