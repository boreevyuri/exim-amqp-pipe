package publisher

//func Connect() (*amqp.Channel, amqp.Queue){
//	conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1/")
//	failOnError(err, "Failed to connect to RabbitMQ")
//	defer conn.Close()
//
//	ch, err := conn.Channel()
//	failOnError(err, "Failed to open channel")
//	defer ch.Close()
//
//	q, err := ch.QueueDeclare(
//		"exim",
//		false,
//		false,
//		false,
//		false,
//		nil,
//	)
//	failOnError(err, "Failed to declare a queue")
//
//	return ch, q
//}
