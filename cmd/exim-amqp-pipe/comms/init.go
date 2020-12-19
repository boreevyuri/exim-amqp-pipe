package comms

//
//import (
//	"errors"
//	"fmt"
//	"github.com/streadway/amqp"
//)
//
//type MessageBody struct {
//	Data []byte
//	Type string
//}
//
//type Message struct {
//	Queue         string
//	ReplyTo       string
//	ContentType   string
//	CorrelationID string
//	Priority      uint8
//	Body          MessageBody
//}
//
//type Connection struct {
//	name     string
//	conn     *amqp.Connection
//	channel  *amqp.Channel
//	exchange string
//	queues   []string
//	err      chan error
//}
//
//const (
//	connectError         = "failed to connect to RabbitMQ: %s"
//	connClosedError      = "connection closed"
//	channelError         = "unable to open channel: %s"
//	exchangeDeclareError = "unable to declare exchange: %s"
//	queueDeclareError    = "unable to declare the queue: %s"
//	queueBindError       = "unable to bind queue: %s"
//	publishError         = "unable to publish: %s"
//)
//
//var connectionPool = make(map[string]*Connection)
//
//func NewConnection(name, exchange string, queues []string) *Connection {
//	if c, ok := connectionPool[name]; ok {
//		return c
//	}
//
//	c := &Connection{
//		exchange: exchange,
//		queues:   queues,
//		err:      make(chan error),
//	}
//	connectionPool[name] = c
//
//	return c
//}
//
//func (c *Connection) Connect() error {
//	var err error
//
//	c.conn, err = amqp.Dial("amqp://guest:guest@127.0.0.1:5672/")
//	if err != nil {
//		return fmt.Errorf(connectError, err)
//	}
//
//	go func() {
//		// Ожидаем сообщения о закрытом соединении
//		<-c.conn.NotifyClose(make(chan *amqp.Error))
//		c.err <- errors.New(connClosedError)
//	}()
//
//	c.channel, err = c.conn.Channel()
//	if err != nil {
//		return fmt.Errorf(channelError, err)
//	}
//
//	if err := c.channel.ExchangeDeclare(
//		c.exchange,
//		"fanout",
//		false,
//		false,
//		false,
//		false,
//		nil,
//	); err != nil {
//		return fmt.Errorf(exchangeDeclareError, err)
//	}
//
//	return nil
//}
//
//func (c *Connection) BindQueue() error {
//	for _, q := range c.queues {
//		if _, err := c.channel.QueueDeclare(
//			q,
//			false,
//			false,
//			false,
//			false,
//			nil,
//		); err != nil {
//			return fmt.Errorf(queueDeclareError, err)
//		}
//
//		if err := c.channel.QueueBind(
//			q,
//			"route_key",
//			c.exchange,
//			false,
//			nil,
//		); err != nil {
//			return fmt.Errorf(queueBindError, err)
//		}
//	}
//
//	return nil
//}
//
//func (c *Connection) Reconnect() error {
//	if err := c.Connect(); err != nil {
//		return err
//	}
//
//	if err := c.BindQueue(); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (c *Connection) Publish(m Message) error {
//	// Неблокирующий канал
//	select {
//	case err := <-c.err:
//		if err != nil {
//			_ = c.Reconnect()
//		}
//	default:
//	}
//
//	p := amqp.Publishing{
//		Headers:       amqp.Table{"type": m.Body.Type},
//		ContentType:   m.ContentType,
//		CorrelationId: m.CorrelationID,
//		Body:          m.Body.Data,
//		ReplyTo:       m.ReplyTo,
//	}
//
//	if err := c.channel.Publish(
//		c.exchange,
//		m.Queue,
//		false,
//		false,
//		p,
//	); err != nil {
//		return fmt.Errorf(publishError, err)
//	}
//
//	return nil
//}
