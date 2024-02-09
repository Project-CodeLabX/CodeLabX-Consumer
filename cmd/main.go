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
	py       *os.File
	java     *os.File
	cpp      *os.File
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

func init() {
	createFiles()
	py.Close()
	java.Close()
	cpp.Close()
	writeToFiles()
}

func createFiles() {
	p, err := os.Create("res/codelabx.py")
	if err != nil {
		log.Println("error in py file creation: ", err)
	}
	j, err := os.Create("res/codelabx.java")
	if err != nil {
		log.Println("error in java file creation: ", err)
	}
	c, err := os.Create("res/codelabx.cpp")
	if err != nil {
		log.Println("error in cpp file creation: ", err)
	}

	py = p
	java = j
	cpp = c
}

func writeToFiles() {
	py, _ := os.OpenFile("res/codelabx.py", os.O_WRONLY, 0666)
	ans, err := py.WriteString("print(\"Hello from CodeLab\")")

	log.Println("ans: ", ans)
	log.Println("err: ", err)
}
