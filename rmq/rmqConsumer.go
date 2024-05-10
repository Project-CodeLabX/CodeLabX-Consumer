package rmq

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	queue        string = "py_events"
	consumerName        = "py_consumer"
)

var (
	url string = os.Getenv("CLX_MQ")
)

type RmqConsumer struct {
	RmqConn    *amqp.Connection
	RmqChannel *amqp.Channel
}

func NewRmqConsumer() *RmqConsumer {
	conn := ConnectToRmq()
	ch := CreateRmqChannel(conn)
	return &RmqConsumer{RmqConn: conn, RmqChannel: ch}
}

func ConnectToRmq() *amqp.Connection {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Println("error in connecto Rmq: ", err)
		return nil
	}
	return conn
}

func CreateRmqChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	if err != nil {
		log.Println("error in channel creation: ", err)
		return nil
	}
	return ch
}

func (c *RmqConsumer) Consume() (<-chan amqp.Delivery, error) {
	return c.RmqChannel.Consume(queue, consumerName, false, false, false, false, nil)
}
