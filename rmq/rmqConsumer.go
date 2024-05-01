package rmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	url          string = "amqps://abhi:Deadshot1060@b-195dfc46-2db6-4582-b92b-ff6bc1a3b4fd.mq.ap-south-1.amazonaws.com:5671/codelabx"
	queue        string = "py_events"
	consumerName        = "py_consumer"
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
