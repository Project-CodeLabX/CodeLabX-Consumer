package main

import (
	"codelabx-consumer/rmq"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"

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
				runFile(userEvent.Language)
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
	p, err := os.Create("res/codelabx.py")
	if err != nil {
		log.Println("error in py file creation: ", err)
	}
	p.Close()
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

	file, err := os.OpenFile(path, os.O_WRONLY, 0333)
	if err != nil {
		log.Println("err in file Writting: ", err)
	}
	defer file.Close()
	file.Truncate(0)
	file.WriteString(userEvent.Code)
}

func runFile(lang string) {
	if lang == "python" {
		runPythonFile()
	} else if lang == "java" {

	} else {

	}
}

func runPythonFile() {
	out, err := exec.Command("python", "res/codelabx.py").CombinedOutput()

	if err != nil {
		log.Println("err in runPython: ", err)
	}
	log.Println("output: ", string(out))
}
