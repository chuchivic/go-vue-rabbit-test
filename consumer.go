package main

import (
  "flag"
  "log"
  "net/http"
  "github.com/streadway/amqp"
  "github.com/gorilla/websocket"
  "encoding/json"

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
    "messages",   // name
    "fanout", // type
    false,     // durable
    false,    // auto-deleted
    false,    // internal
    false,    // no-wait
    nil,      // arguments
  )
  failOnError(err, "Failed to declare an exchange")

  q, err := ch.QueueDeclare(
    "", // name
    false,   // durable
    false,   // delete when unused
    false,   // exclusive
    false,   // no-wait
    nil,     // arguments
  )
  failOnError(err, "Failed to declare a queue")

  err = ch.QueueBind(
    q.Name, // queue name
    "",     // routing key
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

   rabbitch := make(chan string)

   go func() {
     for d := range msgs {
       log.Printf("Received a message: %s", d.Body)
//     rabbitch <- string(d.Body)
     }
   }()


   flag.Parse()
   log.SetFlags(0)
   http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
     if r.Header.Get("Origin") != "http://localhost:8080" {
     http.Error(w, "Origin not allowed", 403)
     return
   }
   conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
   if err != nil {
     http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
   }
   go reader(conn, rabbitch)
 })
 //http.Handle("/socket.io/", wsServer.Server)
 log.Println("Serving at localhost:3000...")
 log.Fatal(http.ListenAndServe(":3000", nil))

 log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
 <-rabbitch


 }


type Message struct {
  Id int`json:"id"`
  NewState string`json:"newstate"`
}

func reader(conn *websocket.Conn, rabbitch <- chan string) {
  for {

    //var m Message

    //err := conn.ReadJSON(&m)
    //if err != nil {
    //  log.Println("Error reading json.", err)
    //}

    log.Printf("Got message rabbith: %#v\n", rabbitch)
    jsone,_ := json.Marshal(rabbitch)
    log.Printf("Got message: %#v\n", jsone)
    err := conn.WriteJSON(jsone); 
    if err != nil {
      log.Println(err)
    }
  }
}
