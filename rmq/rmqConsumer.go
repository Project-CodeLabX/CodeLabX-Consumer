package rmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	url          string = "amqp://Abhi1060:Abhi1060@localhost:5672/codelabx"
	queue        string = "user_events"
	consumerName        = "consumer_1"
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
