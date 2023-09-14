package publisher

import (
	"context"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"os"
	"time"
)

// https://gist.github.com/harrisonturton/c6b62d45e6117d5d03ff44e4e8e1e7f7
// https://www.ribice.ba/golang-rabbitmq-client/

var (
	ErrDisconnected  = errors.New("disconnected from rabbitmq, trying to reconnect")
	ErrAlreadyClosed = errors.New("already closed: not connected to the queue")
	// ErrNotConfirmed  = errors.New("message not confirmed")
)

const (
	reconnectDelay = 5 * time.Second
	maxRetryCount  = 5
	resendDelay    = 2 * time.Second
)

type Client struct {
	pushQueue     string
	logger        zerolog.Logger
	connection    *amqp.Connection
	channel       *amqp.Channel
	done          chan bool
	notifyClose   chan *amqp.Error
	notifyConfirm chan amqp.Confirmation
	isConnected   bool
	// alive         bool
	// threads       int
	// wg            *sync.WaitGroup
}

func New(pushQueue, addr string, l zerolog.Logger) *Client {
	// threads := runtime.GOMAXPROCS(0)
	// if numCPU := runtime.NumCPU(); numCPU > threads {
	//	threads = numCPU
	// }
	//
	client := Client{
		logger:    l,
		pushQueue: pushQueue,
		done:      make(chan bool),
	}
	// client.wg.Add(threads)

	go client.handleReconnect(addr)

	return &client
}

func (c *Client) handleReconnect(addr string) {
	for {
		c.isConnected = false
		t := time.Now()
		// c.logger.Printf("Attempting to connect to rabbitMQ: %s", addr)
		var retryCount int
		for !c.connect(addr) {
			c.logger.Printf("Failed to connect. Retrying. %d tries left", maxRetryCount-retryCount)
			time.Sleep(reconnectDelay + time.Duration(retryCount)*time.Second)
			if retryCount >= maxRetryCount-1 {
				c.logger.Printf("Unable to connect to rabbitMQ in %d tries. Exiting...", maxRetryCount)
				os.Exit(1)
			}
			retryCount++
		}
		c.logger.Printf("Connected to rabbitMQ in: %vms", time.Since(t).Milliseconds())
		select {
		case <-c.done:
			return
		case <-c.notifyClose:
		}
	}
}

func (c *Client) connect(addr string) bool {
	conn, err := amqp.Dial(addr)
	if err != nil {
		c.logger.Printf("failed to dial rabbitMQ server: %v", err)
		return false
	}

	ch, err := conn.Channel()
	if err != nil {
		c.logger.Printf("failed connecting to channel: %v", err)
		return false
	}

	err = ch.Confirm(false)
	if err != nil {
		c.logger.Print("channel does not acknowledged Confirm")
	}

	// _, err = ch.QueueDeclare(
	//	listenQueue,
	//	false,
	//	false,
	//	false,
	//	false,
	//	nil,
	// )
	// if err != nil {
	//	c.logger.Printf("failed to declare listen queue: %v", err)
	//	return false
	// }

	_, err = ch.QueueDeclare(
		c.pushQueue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		c.logger.Printf("Failed to declare push queue: %v", err)
		return false
	}

	c.changeConnection(conn, ch)
	c.isConnected = true
	return true
}

func (c *Client) changeConnection(connection *amqp.Connection, channel *amqp.Channel) {
	c.connection = connection
	c.channel = channel
	c.notifyClose = make(chan *amqp.Error)
	c.notifyConfirm = make(chan amqp.Confirmation)
	c.channel.NotifyClose(c.notifyClose)
	c.channel.NotifyPublish(c.notifyConfirm)
}

// Push will push data onto the queue and wait for a confirmation
// If no confirms are received until within the resendTimeout
// it continuously resends messages until a confirmation is received.
// This will block until the server sends a confirmation.
func (c *Client) Push(message amqp.Publishing) error {
	// if !c.isConnected {
	// 	return errors.New("failed to push: not connected")
	// }

	for {
		err := c.UnsafePush(message)
		if err != nil {
			if errors.Is(err, ErrDisconnected) {
				continue
			}
			return err
		}
		select {
		case confirm := <-c.notifyConfirm:
			if confirm.Ack {
				c.logger.Print("Push confirmed")
				return nil
			}
		case <-time.After(resendDelay):
		}
		c.logger.Print("Push didn't confirm. Retrying...")
	}
}

func (c *Client) UnsafePush(message amqp.Publishing) error {
	if !c.isConnected {
		return ErrDisconnected
	}

	// create simple Background context
	ctx := context.Background()

	return c.channel.PublishWithContext(
		ctx,
		"",
		c.pushQueue,
		false,
		false,
		message,
	)
}

func (c *Client) Close() error {
	if !c.isConnected {
		return ErrAlreadyClosed
	}

	// c.alive = false
	// fmt.Println("Waiting for current messages to be processed...")
	// c.wg.Wait()
	// for i := 1; i <= c.threads; i++ {
	//	fmt.Println("Closing consumer: ", i)
	//	err := c.channel.Cancel(consumerName(i), false)
	//	if err != nil {
	//		return fmt.Errorf("error canceling consumer %s: %v", consumerName(i), err)
	//	}
	// }

	err := c.channel.Close()
	if err != nil {
		return err
	}

	err = c.connection.Close()
	if err != nil {
		return err
	}

	close(c.done)
	c.isConnected = false
	c.logger.Print("Gracefully stopped rabbitMQ connection")
	return nil
}

// func consumerName(i int) string {
//	return fmt.Sprintf("go-consumer-%v", i)
// }
