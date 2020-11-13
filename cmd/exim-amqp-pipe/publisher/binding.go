package publisher

import "github.com/streadway/amqp"

type ExchangeType string

type Binding struct {
	Exchange string `yaml:"exchange"`

	// аргументы точки обмена
	ExchangeArgs amqp.Table

	// имя очереди
	Queue string `yaml:"queue"`

	// аргументы очереди
	QueueArgs amqp.Table

	// тип точки обмена
	Type ExchangeType `yaml:"type"`

	//// ключ маршрутизации
	//Routing string `yaml:"routing"`
}
