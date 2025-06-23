package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel/trace"
)

// ServiceStatus representa o status de um serviço dependente.
type ServiceStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// HealthHandler verifica a saúde dos serviços dependentes.
type HealthHandler struct {
	redisClient *redis.Client
	oracleDB    *sql.DB
	// Use um tipo (ou interface) em vez de chamar uma função:
	rabbitConsumer interface{} // agora é interface vazia para permitir type assertion
	tracerProvider trace.TracerProvider
}

// HealthChecker define o contrato de saúde para o RabbitMQPublisher
// (adicione este trecho em infraestructure/rabbitmq ou num arquivo de interfaces).
type HealthChecker interface {
	IsHealthy() bool
}

// NewHealthHandler cria uma nova instância de HealthHandler.
func NewHealthHandler(
	redisClient *redis.Client,
	oracleDB *sql.DB,
	rabbitConsumer interface{},
	tp trace.TracerProvider,
) *HealthHandler {
	return &HealthHandler{
		redisClient:    redisClient,
		oracleDB:       oracleDB,
		rabbitConsumer: rabbitConsumer,
		tracerProvider: tp,
	}
}

// Check é o handler HTTP para o endpoint /health.
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second) // Aumentado para permitir múltiplas verificações
	defer cancel()

	var wg sync.WaitGroup
	statusChan := make(chan ServiceStatus, 4) // Para Redis, Oracle, RabbitMQ, Jaeger

	overallHealthy := true

	// 1. Verificar Redis
	wg.Add(1)
	go func() {
		defer wg.Done()
		serviceName := "Redis"
		if err := h.redisClient.Ping(ctx).Err(); err != nil {
			statusChan <- ServiceStatus{Name: serviceName, Status: "unhealthy", Error: err.Error()}
			return
		}
		statusChan <- ServiceStatus{Name: serviceName, Status: "healthy"}
	}()

	// 2. Verificar Oracle DB
	wg.Add(1)
	go func() {
		defer wg.Done()
		serviceName := "OracleDB"
		if err := h.oracleDB.PingContext(ctx); err != nil {
			statusChan <- ServiceStatus{Name: serviceName, Status: "unhealthy", Error: err.Error()}
			return
		}
		statusChan <- ServiceStatus{Name: serviceName, Status: "healthy"}
	}()

	// 3. Verificar RabbitMQ
	wg.Add(1)
	go func() {
		defer wg.Done()
		serviceName := "RabbitMQ"
		// O RabbitMQPublisher agora tem ensureChannel que tenta conectar/reconectar.
		// Uma forma de testar é tentar obter o canal ou declarar uma fila de teste.
		// A função `IsConnected` seria ideal no publisher.
		// Por simplicidade, vamos assumir que se `ensureChannel` não retornar erro, está ok.
		// A forma mais simples é checar se a conexão interna do publisher não está nil e não está fechada.
		// (Esta parte depende da implementação exata do seu RabbitMQPublisher)
		hc, ok := h.rabbitConsumer.(HealthChecker)

		if !ok || !hc.IsHealthy() {
			statusChan <- ServiceStatus{Name: serviceName, Status: "unhealthy", Error: "RabbitMQ publisher is not healthy"}
			return
		}
		statusChan <- ServiceStatus{Name: serviceName, Status: "healthy"}
	}()

	// 4. Verificar Jaeger (Tracer)
	wg.Add(1)
	go func() {
		defer wg.Done()
		serviceName := "JaegerTracer"
		// Se o tracerProvider foi injetado e não é nil, e podemos obter um tracer.
		if h.tracerProvider == nil {
			statusChan <- ServiceStatus{Name: serviceName, Status: "unhealthy", Error: "TracerProvider is nil"}
			return
		}
		tracer := h.tracerProvider.Tracer("health-check-tracer")
		if tracer == nil { // Esta verificação pode ser redundante se o provider não for nil
			statusChan <- ServiceStatus{Name: serviceName, Status: "unhealthy", Error: "Failed to get tracer"}
			return
		}
		// Para um teste mais real, você poderia tentar iniciar e finalizar um span.
		_, span := tracer.Start(ctx, "health-check-span")
		if span == nil { // Pouco provável se o tracer for obtido
			statusChan <- ServiceStatus{Name: serviceName, Status: "unhealthy", Error: "Failed to start span"}
			return
		}
		span.End()
		statusChan <- ServiceStatus{Name: serviceName, Status: "healthy"}
	}()

	go func() {
		wg.Wait()
		close(statusChan)
	}()

	var services []ServiceStatus
	for status := range statusChan {
		services = append(services, status)
		if status.Status == "unhealthy" {
			overallHealthy = false
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if !overallHealthy {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	response := map[string]interface{}{
		"overall_status": "healthy",
		"services":       services,
	}
	if !overallHealthy {
		response["overall_status"] = "unhealthy"
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Se não conseguir encodar, a resposta HTTP já pode ter sido parcialmente enviada.
		// Apenas loga.
		http.Error(w, `{"error":"failed to encode health status"}`, http.StatusInternalServerError)
	}
}
