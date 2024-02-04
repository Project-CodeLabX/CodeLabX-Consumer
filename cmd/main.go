package main

import (
	"codelabx-consumer/rmq"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	consumer *rmq.RmqConsumer
)

func main() {

	r := mux.NewRouter()

	consumer = rmq.NewRmqConsumer()

	messageBus, err := consumer.Consume()
	if err != nil {
		log.Println("error in Consuming: ", err)
	}

	go func() {
		defer consumer.RmqChannel.Close()
		defer consumer.RmqConn.Close()
		for msg := range messageBus {
			if msg.Body != nil {
				var userEvent rmq.UserEvent
				err := json.Unmarshal(msg.Body, &userEvent)
				if err != nil {
					log.Println("error happened in json unmarshal in msg: ", err)
					continue
				}
				log.Println("Consumed user event: ", userEvent)
				msg.Ack(false)
				log.Println("acknowledged the message...")
			}
		}
	}()

	http.ListenAndServe(":9010", r)
}
