package main

import (
	"log"

	"github.com/hashicorp/go-plugin"
	"github.com/streadway/amqp"
	"github.com/victorcoder/dkron/dkron"
	dkplugin "github.com/victorcoder/dkron/plugin"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type rabbitMQExecutor struct {
}

func (s *rabbitMQExecutor) Execute(args *dkron.ExecuteRequest) ([]byte, error) {
	queueName := args.Config["queue_name"]
	connectionString := args.Config["connection_string"]

	log.Printf("[rabbitMQExecutor] running job [%s] - queue [%s]\n", args.JobName, queueName)

	conn, err := amqp.Dial(connectionString)
	defer conn.Close()
	if err != nil {
		log.Printf("[rabbitMQExecutor] failed connecting to rabbitmq [%s]", connectionString)
		return []byte("Failed to connect to RabbitMQ"), err
	}

	ch, err := conn.Channel()
	defer ch.Close()
	if err != nil {
		log.Printf("[rabbitMQExecutor] failed to open channel")
		return []byte("Failed to open a channel"), err
	}

	err = ch.Publish(
		"",
		args.Config["queue_name"],
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(args.Config["payload"]),
		})

	if err != nil {
		return []byte(err.Error()), err
	}

	return []byte("OK"), nil
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: dkplugin.Handshake,
		Plugins: map[string]plugin.Plugin{
			"executor": &dkplugin.ExecutorPlugin{Executor: &rabbitMQExecutor{}},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
