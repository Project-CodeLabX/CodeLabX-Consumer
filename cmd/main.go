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
				createResFolder()
				createFiles(userEvent.FileName)
				writeToFile(userEvent.FileName, &userEvent)
				runFile(userEvent.FileName, &userEvent)
				deleteFiles()
				msg.Ack(false)
				log.Println("acknowledged the message...")
			}
		}
	}()

	http.ListenAndServe(":9020", r)
}

func init() {
	// createResFolder()
	redisClient = redis.GetRedisClient()
}

func createResFolder() {
	err := os.MkdirAll("res", os.ModeDir)
	if err != nil {
		log.Println("error in res folder creation : ", err)
	}
}

func deleteFiles() {
	err := os.RemoveAll("res")
	if err != nil {
		log.Println("error in res folder deletion : ", err)
	}
}

func createFiles(fileName string) {
	path := "res/" + fileName + ".java"
	log.Println("path to crate file : ", path)
	p, err := os.Create(path)
	if err != nil {
		log.Println("error in java file creation: ", err)
	}
	p.Close()

}

func writeToFile(fileName string, userEvent *rmq.UserEvent) {
	path := "res/" + fileName + ".java"
	// var path string = "res/codelabx.java"

	file, err := os.OpenFile(path, os.O_WRONLY, 0333)
	if err != nil {
		log.Println("err in file Writting: ", err)
	}
	defer file.Close()
	file.Truncate(0)
	file.WriteString("package res;\n\n" + userEvent.Code)
}

func runFile(filename string, userEvent *rmq.UserEvent) {
	runJavaFile(filename, userEvent)
}

func runJavaFile(fileName string, userEvent *rmq.UserEvent) {
	path := "res/" + fileName + ".java"
	out1, err := exec.Command("javac", path).CombinedOutput()
	if err != nil {
		log.Println("err in runJavac: ", err)
	}
	log.Println("javac output: ", string(out1))
	out, err := exec.Command("java", "res/"+fileName).CombinedOutput()
	if err != nil {
		log.Println("err in runJava: ", err)
	}
	log.Println("java output: ", string(out))
	writeToRedis(userEvent.UserName, string(out1)+"\n"+string(out))
}

func writeToRedis(username string, stdout string) {
	ctx := context.Background()
	err := redisClient.Rdb.Set(ctx, username, stdout, 0).Err()
	if err != nil {
		log.Println("error in inserting into redis: ", err)
	}
}
