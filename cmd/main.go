package main

import (
	"codelabx-consumer/redis"
	"codelabx-consumer/rmq"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/mux"
)

var (
	consumer    *rmq.RmqConsumer
	redisClient *redis.RedisClient
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
				runFile(&userEvent)
				msg.Ack(false)
				log.Println("acknowledged the message...")
			}
		}
	}()

	http.ListenAndServe(":9010", r)
}

func init() {
	createFiles()
	redisClient = redis.GetRedisClient()
}

func createFiles() {
	err := os.Mkdir("res", os.ModePerm)
	if err != nil {
		log.Println("Failed to create res dir with : ", err)
		return
	}
	p, err := os.Create("res/codelabx.py")
	if err != nil {
		log.Println("error in py file creation: ", err)
	}
	p.Close()
}

func writeToFile(userEvent *rmq.UserEvent) {

	path := "res/codelabx.py"

	file, err := os.OpenFile(path, os.O_WRONLY, 0333)
	if err != nil {
		log.Println("err in file Writting: ", err)
	}
	defer file.Close()
	file.Truncate(0)
	file.WriteString(userEvent.Code)
}

func runFile(userEvent *rmq.UserEvent) {
	runPythonFile(userEvent)
}

func runPythonFile(userEvent *rmq.UserEvent) {
	out, err := exec.Command("python", "res/codelabx.py").CombinedOutput()

	if err != nil {
		log.Println("err in runPython: ", err)
	}
	log.Println("output: ", string(out))
	writeToRedis(userEvent.UserName, string(out))
}

func writeToRedis(username string, stdout string) {
	ctx := context.Background()
	err := redisClient.Rdb.Set(ctx, username, stdout, 0).Err()
	if err != nil {
		log.Println("error in inserting into redis: ", err)
	}
}
