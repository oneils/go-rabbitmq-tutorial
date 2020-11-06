# 4 Routing Receiving messages selectively

https://www.rabbitmq.com/tutorials/tutorial-four-go.html

## Producer

To run producer which will generate 10 messages with random severity

```bash
go run example4/logs_producer.go
```

## Consumer

To run consumer for listening following severities:

```bash
go run example4/logs_consumer.go info warning error
```

To consume specific severity:

```bash
go run example4/logs_consumer.go error
```