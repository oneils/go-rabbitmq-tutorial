package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatal("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil) // arguments)

	failOnError(err, "Failed to declare a queue")

	for i := 0; i < 10; i++ {
		sendMessage(channel, queue, i)
	}

	failOnError(err, "Failed to publish a message")
}

func sendMessage(channel *amqp.Channel, queue amqp.Queue, times int) {
	var dotsAmount string
	for i := 0; i < times; i++ {
		dotsAmount += "."
	}
	body := fmt.Sprintf("Hello #%d %s. %v", times, dotsAmount, time.Now().Local())

	channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		},
	)
	log.Printf(" [x] Sent %s", body)
}
