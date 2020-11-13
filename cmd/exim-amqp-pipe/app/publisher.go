package app

type Publish struct {
	Abstract
}

func NewPublish() Application {
	return new(Publish)
}

func (p *Publish) Run() {
	App = p
	p.services = []interface{}{
		//Тут должен быть консумер к amqp
	}
	event := NewEvent(InitAppEvent)
	p.run(p, event)
}

func (p *Publish) FireRun(event *Event, abstractService interface{}) {
	service := abstractService.(PublishService)
	go service.OnPublish(event)
}

//type Config struct {
//	URI      string     `yaml:"uri"`
//	Bindings []*Binding `yaml:"bindings"`
//}
//
//type Service struct {
//	Configs []*Config `yaml:"amqp"`
//	amqpConnections map[string]*amqp.Connection
//	consumers map[string][]*Consumer
//}
//
//func NewService() Service {
//	app := new(Service)
//	app.amqpConnections = make(map[string]*amqp.Connection)
//	app.consumers = make(map[string][]*Consumer)
//}
//
//type Consumer struct {
//	amqpConnection *amqp.Connection
//	binding        *Binding
//	deliveries     <-chan amqp.Delivery
//}
//
//func NewConsumer(amqpConnection *amqp.Connection, binding *Binding) *Consumer {
//	app := new(Consumer)
//	app.amqpConnection = amqpConnection
//	app.binding = binding
//
//	return app
//}
