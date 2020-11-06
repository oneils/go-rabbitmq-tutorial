package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
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
		"logs_direct", // name
		"direct",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare exchange")

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s) // initialize local pseudorandom generator

	severities := [3]string{"info", "warning", "error"}

	for i := 0; i < 10; i++ {
		indx := r.Intn(len(severities))
		severity := severities[indx]
		sendMessage(channel, i, severity)
	}

	failOnError(err, "Failed to publish a message")
}

func sendMessage(channel *amqp.Channel, times int, severity string) {
	var dotsAmount string
	// 'severity' can be one of 'info', 'warning', 'error'.
	for i := 0; i < times; i++ {
		dotsAmount += "."
	}
	body := fmt.Sprintf("%s: Hello #%d %s. %v", severity, times, dotsAmount, time.Now().Local())

	channel.Publish(
		"logs_direct", // exchange
		severity,      // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		},
	)
	log.Printf("%s: [x] Sent %s", severity, body)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatal("%s: %s", msg, err)
	}
}
