package main

import (
	"log"

	"github.com/thiagohmm/integracaoThothConsumer/configuration"
	"github.com/thiagohmm/integracaoThothConsumer/internal/delivery/rabbitmq"
	"github.com/thiagohmm/integracaoThothConsumer/internal/domain/entities/repositories"
	"github.com/thiagohmm/integracaoThothConsumer/internal/infraestructure/database"
	"github.com/thiagohmm/integracaoThothConsumer/internal/usecases"
)

func main() {

	// Carrega as configurações do arquivo .env
	cfg, err := configuration.LoadConfig("../../.env")
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}

	// Conecta ao banco de dados
	db, err := database.ConectarBanco(cfg)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Inicializa o repositório de compra
	compraRepo := &repositories.CompraRepositoryDB{
		DB: db,
	}

	// Inicializa o caso de uso de compra
	compraUseCase := usecases.NewCompraUseCase(compraRepo)

	// Verifica a configuração da URL do RabbitMQ
	rabbitmqURL := cfg.ENV_RABBITMQ
	if rabbitmqURL == "" {
		log.Fatalf("RabbitMQ URL não está definida.")
	}

	// Inicializa o listener da fila RabbitMQ com o caso de uso de compra
	rabbitmqListener := rabbitmq.Listener{
		CompraUC: compraUseCase,
	}

	// Escuta a fila RabbitMQ
	if err := rabbitmqListener.ListenToQueue(rabbitmqURL); err != nil {
		log.Fatalf("Erro ao escutar a fila RabbitMQ: %v", err)
	}

	// Mantém o programa em execução indefinidamente
	select {}
}
