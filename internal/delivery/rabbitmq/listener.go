package listener

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
	"github.com/thiagohmm/integracaoThothConsumer/internal/infraestructure/cache"
	infraestructure "github.com/thiagohmm/integracaoThothConsumer/internal/infraestructure/rabbitmq"
	"github.com/thiagohmm/integracaoThothConsumer/internal/usecases"
	"go.opentelemetry.io/otel"
)

type Listener struct {
	CompraUC  *usecases.CompraUseCase
	EstoqueUC *usecases.EstoqueUseCase
	VendaUC   *usecases.VendaUseCase
	StatusUC  *usecases.StatusUseCase
	Cache     *cache.Cache
	Workers   int // número de workers concorrentes

}

func (l *Listener) getConnectionWithWait(rabbitmqurl string) (*amqp.Connection, error) {
	for {
		conn, err := infraestructure.GetRabbitMQConnection(rabbitmqurl)
		if err == nil {
			return conn, nil
		}
		log.Printf("Erro conectando ao RabbitMQ: %v. Tentando novamente em 5 segundos...", err)
		time.Sleep(5 * time.Second)
	}
}

func (l *Listener) ListenToQueue(rabbitmqurl string) error {
	if rabbitmqurl == "" {
		return fmt.Errorf("rabbitmq URL cannot be empty")
	}

	if l.Workers <= 0 {
		l.Workers = 20 // default to 20 workers if not set
	}

	conn, err := l.getConnectionWithWait(rabbitmqurl)
	if err != nil {
		// Esse ponto nunca deve ser alcançado já que getConnectionWithWait
		// só retorna em caso de sucesso.
		return fmt.Errorf("erro conectando ao RabbitMQ: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("error creating RabbitMQ channel: %w", err)
	}
	defer ch.Close()

	// Configurar prefetch count para controlar quantas mensagens cada worker recebe
	err = ch.Qos(
		l.Workers, // prefetch count
		0,         // prefetch size
		false,     // global
	)
	if err != nil {
		return fmt.Errorf("error setting QoS: %w", err)
	}

	queue := "thothQueue"
	msgs, err := ch.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("error consuming messages: %w", err)
	}

	// Canal para controle de término dos workers (usado como sinal de finalização)
	done := make(chan bool, l.Workers)

	// Iniciar workers
	for i := 0; i < l.Workers; i++ {
		go l.worker(i, msgs, done)
	}

	// Aguardar o término dos workers caso o canal de mensagens seja fechado
	for i := 0; i < l.Workers; i++ {
		<-done
	}
	log.Printf("All workers finished processing messages")

	return nil
}

func (l *Listener) worker(id int, msgs <-chan amqp.Delivery, done chan<- bool) {
	// 1) Recover para capturar panic
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Worker %d recovered from panic: %v", id, r)
		}
	}()

	log.Printf("Worker %d started", id)
	tracer := otel.Tracer("Workers")
	ctx := context.Background()

	for msg := range msgs {
		log.Printf("Worker %d processing message", id)

		if err, uuid := l.processMessage(msg); err != nil {
			// 2) Abra span, registre erro e finalize sem defer
			ctxSpan, span := tracer.Start(ctx, uuid)
			log.Printf("Worker %d - Error processing message: %v", id, err)
			l.Cache.AtualizaStatusProcesso(ctxSpan, uuid, "erro")
			span.RecordError(err)
			span.End()
		}

		msg.Ack(false)
	}

	done <- true
}

func (l *Listener) processMessage(msg amqp.Delivery) (error, string) {
	var message map[string]interface{}
	if err := json.Unmarshal(msg.Body, &message); err != nil {
		return fmt.Errorf("error parsing message: %w", err), ""
	}

	processa, ok := message["processa"].(string)
	if !ok {
		return fmt.Errorf("invalid or missing 'processa' field"), ""
	}

	uuid, ok := message["processo"].(string)
	if !ok {
		return fmt.Errorf("invalid or missing 'processo' field"), uuid
	}

	dados, ok := message["dados"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid or missing 'dados' field"), uuid
	}

	ctx := context.WithValue(context.Background(), "uuid", uuid)

	var err error

	switch processa {
	case "compra":
		_, err = l.CompraUC.ProcessarCompra(ctx, dados)
		log.Printf("Processing purchase with UUID: %s", uuid)
	case "venda":
		err = l.VendaUC.ProcessarVenda(ctx, dados)
	case "estoque":
		err = l.EstoqueUC.ProcessarEstoque(ctx, dados)
	default:
		return fmt.Errorf("unknown process type: %s", processa), uuid
	}

	if err != nil {
		l.Cache.AtualizaStatusProcesso(context.Background(), uuid, "erro")
		l.StatusUC.UpdateStatusProcesso(context.Background(), uuid, "erro")
		log.Printf("Error processing message: %v", err)
		return err, uuid
	}

	log.Printf("Message processed successfully: %s", processa)
	if updErr := l.Cache.AtualizaStatusProcesso(context.Background(), uuid, "Sucesso"); updErr != nil {
		log.Printf("Error updating process status in Redis: %v", updErr)
	}

	if updErr := l.StatusUC.UpdateStatusProcesso(context.Background(), uuid, "Sucesso"); updErr != nil {
		log.Printf("Error updating process status in database: %v", updErr)
	}
	log.Printf("Process status updated successfully for UUID: %s", uuid)
	return nil, uuid
}
