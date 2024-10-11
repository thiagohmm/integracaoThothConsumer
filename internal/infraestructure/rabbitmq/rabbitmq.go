package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

func GetRabbitMQConnection(rabbitmqURL string) (*amqp.Connection, error) {
	rabbitmqUrl := rabbitmqURL
	if rabbitmqUrl == "" {
		return nil, fmt.Errorf("RABBITMQ_URL is not defined")
	}

	conn, err := amqp.Dial(rabbitmqUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	return conn, nil
}
