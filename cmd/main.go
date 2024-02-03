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
		for {
			msg := <-messageBus
			var userEvent rmq.UserEvent
			err := json.Unmarshal(msg.Body, &userEvent)
			if err != nil {
				log.Println("error during json unmarshal in messagebus: ", err)
			}
			log.Println("message : ", userEvent)
			msg.Acknowledger.Ack(1, false)
		}
	}()

	http.ListenAndServe(":9000", r)
}
