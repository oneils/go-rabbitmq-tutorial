# 5 Topics: Receiving messages based on a pattern (topics)

https://www.rabbitmq.com/tutorials/tutorial-five-go.html


## Topic exchange

The limitation for `routing key` is **255 bytes**.

## binding keys:

- `* (star)` can substitute for exactly one word.
- `# (hash)` can substitute for zero or more words.

`Topic exchange` is powerful and can behave like other exchanges.

When a queue is bound with `"#" (hash)` binding key - it will receive **all the messages**, regardless of the routing key - like in fanout exchange.

When special characters `"*" (star)` and `"#" (hash)` aren't used in bindings, the topic exchange will behave just like a `direct` one.

## Producer

To run producer%

```bash
go run example5/logs_producer.go "any.critical" "Any critical error from any channel"
```

```bash
go run example5/logs_producer.go "kern.warning" "Some warning from kernel"
```

```bash
go run example5/logs_producer.go "some.warning" "Some warning message"
```

```bash
go run example5/logs_producer.go "kern.critical" "A critical kernel error"
```
## Consumer

To receive all the logs:

```bash
go run example5/logs_consumer.go "#"
```

To receive all logs from the facility `"kern"`:

```bash
go run example5/logs_consumer.go "kern.*"
```

Or if you want to hear only about "critical" logs:

```bash
go run example5/logs_consumer.go "*.critical"
```