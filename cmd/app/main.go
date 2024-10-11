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

	// Inicializa a conexão com o banco de dados
	database, err := database.ConectarBanco(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Inicializa o repositório de compra
	compraRepo := repositories.CompraRepositoryDB{
		DB: database,
	}

	compraUseCase := usecases.CompraUseCase{
		Repo: compraRepo,
	}

	// Configura a conexão com RabbitMQ
	rabbitmqURL := cfg.ENV_RABBITMQ
	if rabbitmqURL == "" {
		log.Fatalf("RabbitMQ URL is not defined.")
	}

	// Inicializa o listener da fila
	//rabbitmq.ListenQueue(rabbitmqURL, compraUseCase)

	rabbitmqListener := rabbitmq.Listener{
		CompraUC: compraUseCase,
		// Adicione mais usecases conforme necessário
	}

	if err := rabbitmqListener.ListenToQueue(rabbitmqURL); err != nil {
		log.Fatalf("Error listening to RabbitMQ queue: %v", err)
	}

	// Manter o programa em execução
	select {}
}
