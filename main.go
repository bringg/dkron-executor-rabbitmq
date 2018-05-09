package main

import (
	"log"

	"github.com/hashicorp/go-plugin"
	"github.com/victorcoder/dkron/dkron"
	dkplugin "github.com/victorcoder/dkron/plugin"
	"github.com/streadway/amqp"
	"github.com/spf13/viper"
)

type rabbitMQExecutor struct {
	conn *amqp.Connection
	ch *amqp.Channel
}

func createRabbitMQExecutor() (*rabbitMQExecutor, error) {
	connectionString := viper.GetString("rabbit_host")
	log.Printf("[rabbitMQExecutor] connecting to rabbit [%s]", connectionString)

	conn, err := amqp.Dial(connectionString)
	if err != nil {
		log.Printf("[rabbitMQExecutor] failed connecting to rabbitmq [%s]", connectionString)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("[rabbitMQExecutor] failed to open channel")
		return nil, err
	}

	return &rabbitMQExecutor{conn, ch}, nil

}

func (s *rabbitMQExecutor) Execute(args *dkron.ExecuteRequest) ([]byte, error) {
	// TODO: validate parameters - not nil

	queueName := args.Config["queue_name"]

	log.Printf("[rabbitMQExecutor] running job [%s] - queue [%s]\n", args.JobName, queueName)


	err := s.ch.Publish(
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

	viper.SetDefault("rabbit_host", "amqp://guest:guest@localhost:5672/")

	viper.SetConfigName("dkron-executor-rabbitmq")        // name of config file (without extension)
	viper.AddConfigPath("/etc/dkron")   // call multiple times to add many search paths
	viper.AddConfigPath("$HOME/.dkron") // call multiple times to add many search paths
	viper.AddConfigPath("./config")     // call multiple times to add many search paths
	viper.SetEnvPrefix("dkron_executor_rabbitmq")         // will be uppercased automatically
	viper.AutomaticEnv()

	viper.ReadInConfig()


	executor, err := createRabbitMQExecutor()
	if err != nil {
		log.Printf("Failed to create rabbit mq executor, error: %v", err)
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: dkplugin.Handshake,
		Plugins: map[string]plugin.Plugin{
			"executor": &dkplugin.ExecutorPlugin{Executor: executor},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
