package rabbitmq

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
	"github.com/thiagohmm/integracaoThothConsumer/configuration"
)

type RabbitMQConsumer struct {
	URL                   string
	GetRabbitMQConnection func(string) (*amqp.Connection, error)
}

// HealthChecker define o contrato de saúde para o RabbitMQConsumer.
type HealthChecker interface {
	IsHealthy() bool
}

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
func (c RabbitMQConsumer) IsHealthy() bool {
	// tenta abrir e fechar uma conexão rápida
	rabbitmqURL := os.Getenv("ENV_RABBITMQ")
	if rabbitmqURL == "" {
		cfg, err := loadConfig()
		if err != nil {
			log.Printf("Failed to load config: %v", err)
			return false
		}
		rabbitmqURL = cfg.ENV_RABBITMQ
	}
	_, err := c.GetRabbitMQConnection(rabbitmqURL)
	if err != nil {
		log.Printf("RabbitMQ connection failed: %v", err, rabbitmqURL)
		return false
	}
	//conn.Close()
	return true
}

func loadConfig() (*configuration.Conf, error) {
	// Carrega as configurações do arquivo .env
	cfg, err := configuration.LoadConfig("../../.env")
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}
	return cfg, err
}
