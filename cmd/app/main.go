package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"

	"github.com/thiagohmm/integracaoThothConsumer/configuration"
	config "github.com/thiagohmm/integracaoThothConsumer/configuration"
	listener "github.com/thiagohmm/integracaoThothConsumer/internal/delivery/rabbitmq"
	"github.com/thiagohmm/integracaoThothConsumer/internal/domain/repositories"
	"github.com/thiagohmm/integracaoThothConsumer/internal/infraestructure/cache"
	"github.com/thiagohmm/integracaoThothConsumer/internal/infraestructure/database"
	"github.com/thiagohmm/integracaoThothConsumer/internal/infraestructure/rabbitmq"

	httphandler "github.com/thiagohmm/integracaoThothConsumer/internal/domain/handler" // Ajuste o path

	httprouter "github.com/thiagohmm/integracaoThothConsumer/internal/infraestructure/router"
	"github.com/thiagohmm/integracaoThothConsumer/internal/usecases"

	go_ora "github.com/sijms/go-ora/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func LoadConfig() (*config.Conf, error) {
	// Carrega as configurações do arquivo .env
	cfg, err := configuration.LoadConfig("../../.env")
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}
	return cfg, err
}

// initTracer inicia o provedor de rastreamento do jaeger.
func initTracer() (*trace.TracerProvider, error) {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(

		jaeger.WithEndpoint(cfg.JAEGER_ENDPOINT),
	))
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("IntegracaoThothConsumer"),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

// main é a função principal do aplicativo.
func main() {
	// Inicializa o provedor de rastreamento
	// Carrega as configurações do arquivo .env
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}

	// Inicializa o provedor de rastreamento no jaeger
	tp, err := initTracer()
	if err != nil {
		// Se houver um erro ao inicializar o provedor de rastreamento,
		// o programa termina com uma mensagem de erro.
		// Isso garante que o programa não continue a executar sem rastreamento.
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// Inicia um span para a operação principal.
	ctx, span := otel.Tracer("main").Start(context.Background(), "main-operation")
	defer span.End()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.ENV_REDIS_ADDR,
		Password: cfg.ENV_REDIS_PASSWORD,
	})
	oracleURLOptions := map[string]string{"CONNECTION TIMEOUT": "10", "ssl": "true", "ssl verify": "false"}
	oracleConnStr := go_ora.BuildUrl(cfg.Host, cfg.Port, cfg.ServiceName, cfg.DBUser, cfg.DBPassword, oracleURLOptions)
	oracleDB, err := sql.Open(cfg.DBDriver, oracleConnStr)

	rabbitConsumer := rabbitmq.RabbitMQConsumer{
		GetRabbitMQConnection: rabbitmq.GetRabbitMQConnection,
	}
	if err != nil {
		log.Fatalf("Erro ao conectar ao RabbitMQ: %v", err)
	}

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		log.Fatalf("Falha ao conectar com Redis: %v", err)
	}

	// Conecta ao banco de dados
	db, err := database.ConectarBanco(cfg)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Inicia um span para a conexão com o banco de dados.
	tracer := otel.Tracer("integracaoThothConsumer")
	ctx, span = tracer.Start(context.Background(), "ConectarBanco")
	defer span.End()

	// Verifica a conexão com o banco de dados
	if err := db.Ping(); err != nil {
		log.Fatalf("Erro ao verificar a conexão: %v", err)
	}

	// Inicializa o repositório de compra
	compraRepo := &repositories.CompraRepositoryDB{
		DB: db,
	}

	estoqueRepo := &repositories.EstoqueRepositoryDB{
		DB: db,
	}

	vendaRepo := &repositories.VendaRepositoryDB{
		DB: db,
	}

	statusRepo := &repositories.StatusRepositoryDB{
		DB: db,
	}

	// Inicializa o caso de uso de compra
	compraUseCase := usecases.NewCompraUseCase(compraRepo)
	estoqueUseCase := usecases.NewEstoqueUseCase(estoqueRepo)
	vendaUseCase := usecases.NewVendaUseCase(vendaRepo)
	statusUseCase := usecases.NewStatusUseCase(statusRepo)
	// Verifica a configuração da URL do RabbitMQ
	rabbitmqURL := cfg.ENV_RABBITMQ
	if rabbitmqURL == "" {
		log.Fatalf("RabbitMQ URL não está definida.")
	}

	cache := cache.NewCache(cfg.ENV_REDIS_ADDR, cfg.ENV_REDIS_PASSWORD)
	// Inicializa o listener da fila RabbitMQ com o caso de uso de compra
	rabbitmqListener := listener.Listener{

		CompraUC:  compraUseCase,
		EstoqueUC: estoqueUseCase,
		VendaUC:   vendaUseCase,
		Cache:     cache,
		StatusUC:  statusUseCase,
		Workers:   20, // Número de workers concorrentes
		// Adicione mais usecases conforme necessário,
	}
	healthHdlr := httphandler.NewHealthHandler(redisClient, oracleDB, rabbitConsumer, tp)
	go func() { // Roteador
		chiRouter := chi.NewRouter()
		// Passe o filterHdlr para SetupRoutes
		httprouter.SetupRoutes(chiRouter, healthHdlr)

		log.Println("Servidor iniciado com sucesso na porta 3010")
		if err := http.ListenAndServe(":3010", chiRouter); err != nil {
			log.Fatalf("Falha ao iniciar servidor HTTP: %v", err)
		}
	}()
	// Inicia um span para a escuta da fila RabbitMQ.
	// Escuta a fila RabbitMQ
	_, listenSpan := tracer.Start(ctx, "ListenToQueue")
	if err := rabbitmqListener.ListenToQueue(rabbitmqURL); err != nil {
		log.Fatalf("Erro ao escutar a fila RabbitMQ: %v", err)
	}

	listenSpan.End()

	select {}
}
