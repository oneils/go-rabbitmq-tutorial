package main

import (
	"github.com/streadway/amqp"
	"log"
	"os"
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
		"logs_topic", // name
		"topic",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
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

	if len(os.Args) < 2 {
		log.Printf("Usage: %s [binding_key]...", os.Args[0])
		os.Exit(0)
	}

	for _, s := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing key %s",
			q.Name, "logs_topic", s)
		err = channel.QueueBind(
			q.Name,       // queue name
			s,            // routing key
			"logs_topic", // exchange
			false,
			nil)
		failOnErrorConsumer(err, "Failed to bind a queue")
	}

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
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever

}
