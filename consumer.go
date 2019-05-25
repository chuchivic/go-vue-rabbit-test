package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// rabbitch := make(chan string)

	flag.Parse()
	log.SetFlags(0)
	connmq, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	failOnError(err, "Failed to connect to RabbitMQ")
	log.Println("Connected to rabbitMQ")
	defer connmq.Close()

	ch, err := connmq.Channel()
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

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,     // queue name
		"",         // routing key
		"messages", // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Origin") != "http://localhost:8080" {
			http.Error(w, "Origin not allowed", 403)
			return
		}
		log.Println("INICIADO")
		conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
		if err != nil {
			http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		}

		go func() {
			for d := range msgs {
				log.Printf("Received a message: %s", d.Body)
				reader(conn, d.Body)
			}
		}()
	})

	//http.Handle("/socket.io/", wsServer.Server)
	log.Println("Serving at localhost:3000...")
	log.Fatal(http.ListenAndServe(":3000", nil))

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	//<-rabbitch

}

type Message struct {
	Id       int    `json:"id"`
	NewState string `json:"newstate"`
}

func reader(conn *websocket.Conn, message []byte) {

	//var m Message

	//err := conn.ReadJSON(&m)
	//if err != nil {
	//  log.Println("Error reading json.", err)
	//}

	log.Printf("Got message: %#v\n", string(message))
	err := conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println(err)
	}
}
