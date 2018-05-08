# dkron-executor-rabbitmq

RabbitMQ Executor for dkron that publishes given message to given queue
example:

```json
{
  "name": "send_emails",
  "schedule": "@every 1m",
  "shell": false,
  "executor": "rabbitmq",
  "executor_config": {
    "connection_string": "amqp://guest:guest@localhost:5672/",
    "queue_name": "email-service",
    "payload": "{\"email\": \"yosi@email.com\", \"subject\": \"hello world\"}"
  },
  "disabled": false
}
```
