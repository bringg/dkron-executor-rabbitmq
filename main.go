package main

import (
	"log"

	"github.com/hashicorp/go-plugin"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"github.com/victorcoder/dkron/dkron"
	dkplugin "github.com/victorcoder/dkron/plugin"
	"strconv"
)

type rabbitMQExecutor struct {
	conn *amqp.Connection
	ch   *amqp.Channel
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

func fetchConfig(config map[string]string, key string, defaultValue string) string {
	value, ok := config[key]
	if ok {
		return value
	} else {
		return defaultValue
	}
}

func (s *rabbitMQExecutor) Execute(args *dkron.ExecuteRequest) ([]byte, error) {
	queueName := args.Config["queue_name"]
	log.Printf("[rabbitMQExecutor] [%s] will publish to queue [%s]\n", args.JobName, queueName)

	mandatory, err := strconv.ParseBool(fetchConfig(args.Config, "mandatory", "false"))
	if err != nil {
		log.Printf("[rabbitMQExecutor] [%s] 'mandatory' have invalid value - %v", args.JobName, err)
	}

	immediate, err := strconv.ParseBool(fetchConfig(args.Config, "immediate", "false"))
	if err != nil {
		log.Printf("[rabbitMQExecutor] [%s] 'immediate' have invalid value - %v", args.JobName, err)
	}

	err = s.ch.Publish(
		fetchConfig(args.Config, "exchange", ""),
		args.Config["queue_name"],
		mandatory,
		immediate,
		amqp.Publishing{
			ContentType: fetchConfig(args.Config, "header", "application/json"),
			Body:        []byte(args.Config["payload"]),
		})

	if err != nil {
		return []byte(err.Error()), err
	}

	return []byte("OK"), nil
}

func main() {

	viper.SetDefault("rabbit_host", "amqp://guest:guest@localhost:5672/")

	viper.SetConfigName("dkron-executor-rabbitmq") // name of config file (without extension)
	viper.AddConfigPath("/etc/dkron")              // call multiple times to add many search paths
	viper.AddConfigPath("$HOME/.dkron")            // call multiple times to add many search paths
	viper.AddConfigPath("./config")                // call multiple times to add many search paths
	viper.SetEnvPrefix("dkron_executor_rabbitmq")  // will be uppercased automatically
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
