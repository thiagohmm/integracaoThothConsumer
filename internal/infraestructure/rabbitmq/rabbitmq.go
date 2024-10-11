package rabbitmq

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

func GetRabbitMQConnection() (*amqp.Connection, error) {
	rabbitmqUrl := os.Getenv("RABBITMQ_URL")
	if rabbitmqUrl == "" {
		return nil, fmt.Errorf("RABBITMQ_URL is not defined")
	}

	conn, err := amqp.Dial(rabbitmqUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	return conn, nil
}
