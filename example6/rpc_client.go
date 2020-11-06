package main

import (
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func failOnErrorConsumer(err error, msg string) {
	if err != nil {
		log.Fatal("%s: %s", msg, err)
	}
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func fibonacciRPC(n int) (res int, err error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnErrorConsumer(err, "Failed to connect to RbbitMQ")
	defer conn.Close()

	channel, err := conn.Channel()
	failOnErrorConsumer(err, "Eror while creating channel")
	defer channel.Close()

	q, err := channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	failOnErrorConsumer(err, "Failed to declare a queue")

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

	corrId := randomString(32)

	err = channel.Publish(
		"",          //exchange
		"rpc_queue", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(strconv.Itoa(n)),
		},
	)
	failOnErrorConsumer(err, "Failed to publish a message")

	for d := range msgs {
		if d.CorrelationId == corrId {
			res, err = strconv.Atoi(string(d.Body))
			failOnErrorConsumer(err, "Failed to convert body to integer")
			break
		}
	}

	return
}

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	n := bodyFrom(os.Args)

	log.Printf(" [x] Requesting fib(%d)", n)
	res, err := fibonacciRPC(n)
	failOnErrorConsumer(err, "Failed to handle RPC request")

	log.Printf(" [.] Got %d", res)

}

func bodyFrom(args []string) int {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "30"
	} else {
		s = strings.Join(args[1:], " ")
	}
	n, err := strconv.Atoi(s)
	failOnErrorConsumer(err, "Failed to convert arg to integer")
	return n
}
