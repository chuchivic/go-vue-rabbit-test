package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"os"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"messages", // name
		"fanout",   // type
		false,      // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	for {
		time.Sleep(200 * time.Millisecond)
		body := bodyFrom(os.Args)
		err = ch.Publish(
			"messages", // exchange
			"",         // routing key
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		failOnError(err, "Failed to publish a message")
		log.Printf(" [x] Sent %s", body)
	}
}

type Message struct {
	Id       int    `json:"id"`
	NewState string `json:"newstate"`
}

func bodyFrom(args []string) []byte {
	num := rand.Intn(100)
	state := "operational"
	if num > 30 && num < 70 {
		state = "maintenance"
	} else if num >= 70 {
		state = "error"
	} else if num <= 30 {
		state = "operational"
	}
	id := rand.Intn(6)
	message := Message{Id: id, NewState: state}
	res, _ := json.Marshal(message)
	return res

}
