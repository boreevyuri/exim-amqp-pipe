package comms

//package main
//
//import (
//	"comms"
//	"github.com/streadway/amqp"
//	"log"
//)
//
//func main() {
//	forever := make(chan bool)
//	conn := comms.NewConnection("my-consumer-1", "my-exchange", []string{"queue-1", "queue-2"})
//	if err := conn.Connect(); err != nil {
//		panic(err)
//	}
//	if err := conn.BindQueue(); err != nil {
//		panic(err)
//	}
//	deliveries, err := conn.Consume()
//	if err != nil {
//		panic(err)
//	}
//	for q, d := range deliveries {
//		go conn.HandleConsumedDeliveries(q, d, messageHandler)
//	}
//	<-forever
//}
//
//func messageHandler(c comms.Connection, q string, deliveries <-chan amqp.Delivery) {
//	for d := range deliveries {
//		m := comms.Message{
//			Queue:         q,
//			Body:          comms.MessageBody{Data: d.Body, Type: d.Headers["type"].(string)},
//			ContentType:   d.ContentType,
//			Priority:      d.Priority,
//			CorrelationID: d.CorrelationId,
//		}
//		//handle the custom message
//		log.Println("Got message from queue ", m.Queue)
//		d.Ack(false)
//	}
//}
