package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	redisClientInstance *redis.Client
	once                sync.Once
)

func GetRedisClient(redisAddr, redisPassword string) *redis.Client {
	once.Do(func() {
		redisClientInstance = redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: redisPassword,
			DB:       0, // usa o DB padrão
		})
	})
	return redisClientInstance
}

type Cache struct {
	RedisClient *redis.Client
}

func NewCache(redisAddr, redisPassword string) *Cache {
	client := GetRedisClient(redisAddr, redisPassword)
	return &Cache{
		RedisClient: client,
	}
}

func (c *Cache) AtualizaStatusProcesso(ctx context.Context, uuid string, novoStatus string) error {
	fmt.Println("Atualizando status do processo no Redis...", uuid, novoStatus)

	// Obter o valor atual do processo no Redis
	val, err := c.RedisClient.Get(ctx, uuid).Result()
	if err == redis.Nil {
		return fmt.Errorf("processo com UUID %s não encontrado", uuid)
	} else if err != nil {
		return fmt.Errorf("erro ao obter processo do Redis: %w", err)
	}

	// Deserializar o valor JSON
	var processo map[string]interface{}
	if err := json.Unmarshal([]byte(val), &processo); err != nil {
		return fmt.Errorf("erro ao deserializar processo: %w", err)
	}

	// Atualizar o status do processo
	processo["statusProcesso"] = novoStatus

	// Serializar o valor atualizado de volta para JSON
	updatedVal, err := json.Marshal(processo)
	if err != nil {
		fmt.Errorf("erro ao serializar processo atualizado: %w", err)
		return err
	}

	expireDaysStr := os.Getenv("ENV_REDIS_EXPIRE")
	expireDays, err := strconv.Atoi(expireDaysStr)
	if err != nil {
		expireDays = 2 // Valor padrão se a variável de ambiente não estiver definida ou for inválida
	}

	expiration := time.Duration(expireDays) * 24 * time.Hour

	// Atualizar o valor no Redis
	if err := c.RedisClient.Set(ctx, uuid, updatedVal, expiration).Err(); err != nil {
		fmt.Errorf("erro ao atualizar processo no Redis: %w", err)
		return err
	}

	fmt.Println("Status atualizado com sucesso no Redis.")
	return nil
}
