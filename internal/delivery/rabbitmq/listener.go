// listener.go
package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	infraestructure "github.com/thiagohmm/integracaoThothConsumer/internal/infraestructure/rabbitmq"
	"github.com/thiagohmm/integracaoThothConsumer/internal/usecases"
)

type Listener struct {
	CompraUC  *usecases.CompraUseCase
	EstoqueUC *usecases.EstoqueUseCase
	//VendaUC *usecases.VendaUseCase
	// Adicione mais usecases conforme necessário
}

func (l *Listener) ListenToQueue(rabbitmqurl string) error {
	conn, err := infraestructure.GetRabbitMQConnection(rabbitmqurl)
	if err != nil {
		return fmt.Errorf("error connecting to RabbitMQ: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("error creating RabbitMQ channel: %w", err)
	}
	defer ch.Close()

	queue := "thothQueue"
	msgs, err := ch.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("error consuming messages: %w", err)
	}

	for msg := range msgs {
		var message map[string]interface{}
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			fmt.Println("Error parsing message:", err)
			continue
		}

		processa := message["processa"].(string)
		if processa == "compra" {
			err = l.CompraUC.ProcessarCompra(context.Background(), message["dados"].(map[string]interface{}))
		} else if processa == "venda" {
			// Chame o caso de uso para venda
		} else if processa == "estoque" {
			// Chame o caso de uso para estoque
			err = l.EstoqueUC.ProcessarEstoque(context.Background(), message["dados"].(map[string]interface{}))
		}

		if err != nil {
			fmt.Println("Error processing message:", err)
		} else {
			msg.Ack(false)
		}
	}

	return nil
}
