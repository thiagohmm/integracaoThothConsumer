// // listener.go

package listener

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/thiagohmm/integracaoThothConsumer/internal/infraestructure/cache"
	infraestructure "github.com/thiagohmm/integracaoThothConsumer/internal/infraestructure/rabbitmq"
	"github.com/thiagohmm/integracaoThothConsumer/internal/usecases"
)

type Listener struct {
	CompraUC  *usecases.CompraUseCase
	EstoqueUC *usecases.EstoqueUseCase
	VendaUC   *usecases.VendaUseCase
	Cache     *cache.Cache
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
		var processErr error // Initialize err for use case processing
		var message map[string]interface{}
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			log.Printf("Error parsing message JSON: %v, raw message: %s", err, string(msg.Body))
			msg.Nack(false, false) // NACK the message
			continue
		}

		processa := message["processa"].(string)
		if processa == "compra" {
			uuid := message["processo"].(string)
			ctx := context.WithValue(context.Background(), "uuid", uuid)
			_, processErr = l.CompraUC.ProcessarCompra(ctx, message["dados"].(map[string]interface{}))

			if processErr != nil {
				l.Cache.AtualizaStatusProcesso(context.Background(), uuid, "erro")
				log.Printf("Error processing '%s' for UUID %s: %v", processa, uuid, processErr)
			} else {
				l.Cache.AtualizaStatusProcesso(context.Background(), uuid, "Sucesso")
				log.Printf("Successfully processed '%s' for UUID %s", processa, uuid)
			}
		} else if processa == "venda" {
			// Chame o caso de uso para venda
			uuid := message["processo"].(string)
			ctx := context.WithValue(context.Background(), "uuid", uuid)
			processErr = l.VendaUC.ProcessarVenda(ctx, message["dados"].(map[string]interface{}))

			if processErr != nil {
				l.Cache.AtualizaStatusProcesso(context.Background(), uuid, "erro")
				log.Printf("Error processing '%s' for UUID %s: %v", processa, uuid, processErr)
			} else {
				l.Cache.AtualizaStatusProcesso(context.Background(), uuid, "Sucesso")
				log.Printf("Successfully processed '%s' for UUID %s", processa, uuid)
			}
		} else if processa == "estoque" {
			// Chame o caso de uso para estoque
			uuid := message["processo"].(string)
			ctx := context.WithValue(context.Background(), "uuid", uuid)
			processErr = l.EstoqueUC.ProcessarEstoque(ctx, message["dados"].(map[string]interface{}))

			if processErr != nil {
				l.Cache.AtualizaStatusProcesso(context.Background(), uuid, "erro")
				log.Printf("Error processing '%s' for UUID %s: %v", processa, uuid, processErr)
			} else {
				l.Cache.AtualizaStatusProcesso(context.Background(), uuid, "Sucesso")
				log.Printf("Successfully processed '%s' for UUID %s", processa, uuid)
			}
		}

		// After processing, decide whether to ACK or NACK
		if processErr != nil {
			msg.Nack(false, false)
		} else {
			msg.Ack(false)
		}
	}

	return nil
}
