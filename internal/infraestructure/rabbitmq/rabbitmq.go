package rabbitmq

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func GetRabbitMQConnection(rabbitmqURL string) (*amqp.Connection, error) {
	rabbitmqUrl := rabbitmqURL
	if rabbitmqUrl == "" {
		return nil, fmt.Errorf("RABBITMQ_URL is not defined")
	}

	var conn *amqp.Connection
	var err error

	for {
		conn, err = amqp.Dial(rabbitmqUrl)
		if err == nil {
			log.Println("Successfully connected to RabbitMQ")
			break
		}

		log.Printf("Failed to connect to RabbitMQ: %v. Retrying in 5 seconds...", err)
		time.Sleep(5 * time.Second)
	}

	go func() {
		for {
			if conn.IsClosed() {
				log.Println("Connection to RabbitMQ lost. Reconnecting...")
				for {
					conn, err = amqp.Dial(rabbitmqUrl)
					if err == nil {
						log.Println("Successfully reconnected to RabbitMQ")
						break
					}

					log.Printf("Failed to reconnect to RabbitMQ: %v. Retrying in 5 seconds...", err)
					time.Sleep(5 * time.Second)
				}
			}
			time.Sleep(1 * time.Second)
		}
	}()

	return conn, nil
}
