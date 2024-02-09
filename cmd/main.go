package main

import (
	"codelabx-consumer/rmq"
	"encoding/json"
	"log"
	"net/http"
	"os"

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
				writeToFile(&userEvent)
				msg.Ack(false)
				log.Println("acknowledged the message...")
			}
		}
	}()

	http.ListenAndServe(":9010", r)
}

func init() {
	createFiles()
}

func createFiles() {
	_, err := os.Create("res/codelabx.py")
	if err != nil {
		log.Println("error in py file creation: ", err)
	}
	_, err1 := os.Create("res/codelabx.java")
	if err != nil {
		log.Println("error in java file creation: ", err1)
	}
	_, err2 := os.Create("res/codelabx.cpp")
	if err != nil {
		log.Println("error in cpp file creation: ", err2)
	}
}

func writeToFile(userEvent *rmq.UserEvent) {
	var path string
	if userEvent.Language == "python" {
		path = "res/codelabx.py"
	} else if userEvent.Language == "java" {
		path = "res/codelabx.java"
	} else {
		path = "res/codelabx.cpp"
	}

	file, err := os.OpenFile(path, os.O_WRONLY, 0666)
	if err != nil {
		log.Println("err in file Writting: ", err)
	}
	defer file.Close()

	file.WriteString(userEvent.Code)
}
