package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	RedisClient *redis.Client
}

func NewCache(redisAddr, redisPassword string) *Cache {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0, // use default DB
	})

	return &Cache{
		RedisClient: redisClient,
	}
}

func (c *Cache) AtualizaStatusProcesso(ctx context.Context, uuid string, novoStatus string) error {
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
		return fmt.Errorf("erro ao serializar processo atualizado: %w", err)
	}

	expiration := 5 * 24 * time.Hour
	// Atualizar o valor no Redis
	if err := c.RedisClient.Set(ctx, uuid, updatedVal, expiration).Err(); err != nil {
		return fmt.Errorf("erro ao atualizar processo no Redis: %w", err)
	}

	return nil
}
