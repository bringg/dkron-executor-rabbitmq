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
    "queue_name": "email-service",
    "payload": "{\"email\": \"yosi@email.com\", \"subject\": \"hello world\"}"
  },
  "disabled": false
}
```

## configuration

In order for `dkron-executor-rabbitmq` to know which rabbitmq it should connect, you need to pass configuration, there are multiple options:

- Next to `config/dkron.json`, create `config/dkron-executor-rabbitmq.json` with content `{ "rabbit_host": "..." }`
- Specify environment variable - `DKRON_EXECUTOR_RABBITMQ_RABBIT_HOST` with the rabbit host
