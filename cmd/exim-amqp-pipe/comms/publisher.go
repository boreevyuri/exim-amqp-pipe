package comms

//package main
//
//func main() {
//	conn := comms.NewConnection("my-producer",
//	"my-exchange", []string{"queue-1", "queue-2"})
//	if err := conn.Connect(); err != nil {
//		panic(err)
//	}
//	if err := conn.BindQueue(); err != nil {
//		panic(err)
//	}
//	for _, q := range c.queues {
//		m := comms.Message{
//			Queue: q,
//			//set the necessary fields
//		}
//		if err := conn.Publish(m); err != nil {
//			panic(err)
//		}
//	}
//}
//
