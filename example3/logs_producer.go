package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
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
	failOnError(err, "Failed to declare exchange")

	for i := 0; i < 10; i++ {
		sendMessage(channel, i)
	}

	failOnError(err, "Failed to publish a message")
}

func sendMessage(channel *amqp.Channel, times int) {
	var dotsAmount string
	for i := 0; i < times; i++ {
		dotsAmount += "."
	}
	body := fmt.Sprintf("Hello #%d %s. %v", times, dotsAmount, time.Now().Local())

	channel.Publish(
		"logs", // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		},
	)
	log.Printf(" [x] Sent %s", body)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatal("%s: %s", msg, err)
	}
}
