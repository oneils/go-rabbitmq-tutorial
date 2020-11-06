package main

import (
	"bytes"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func failOnErrorConsumer(err error, msg string) {
	if err != nil {
		log.Fatal("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnErrorConsumer(err, "Failed to connect to RbbitMQ")
	defer conn.Close()

	channel, err := conn.Channel()
	failOnErrorConsumer(err, "Eror while creating channel")
	defer channel.Close()

	err = channel.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnErrorConsumer(err, "Failed to declare a queue")

	q, err := channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnErrorConsumer(err, "Failed to declare a queue")

	err = channel.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"logs", // exchange
		false,
		nil,
	)
	failOnErrorConsumer(err, "Failed to bind a queue")

	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnErrorConsumer(err, "Failed to set QoS")

	msgs, err := channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnErrorConsumer(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			dot_count := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dot_count)
			time.Sleep(t * time.Second)
			log.Printf("Done")
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
