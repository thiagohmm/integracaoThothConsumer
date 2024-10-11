package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	config "github.com/thiagohmm/integracaoThothConsumer/configuration"

	go_ora "github.com/sijms/go-ora/v2"
)

func ConectarBanco(cfg *config.Conf) (*sql.DB, error) {
	// Configurar as opções de URL
	urlOptions := map[string]string{
		"CONNECTION TIMEOUT": "10",
		"ssl":                "true",  // or enable
		"ssl verify":         "false", // stop ssl certificate verification
		//"wallet":             "./wallet",
	}

	// Construir a string de conexão
	connStr := go_ora.BuildUrl(cfg.Host, cfg.Port, cfg.ServiceName, cfg.DBUser, cfg.DBPassword, urlOptions)

	// Abrir a conexão com o banco de dados
	db, err := sql.Open(cfg.DBDriver, connStr)
	if err != nil {
		log.Fatalf("Erro ao abrir a conexão: %v", err)
	}

	// Verificar a conexão
	ctx := context.Background()
	if err = db.PingContext(ctx); err != nil {
		log.Fatalf("Erro ao verificar a conexão: %v", err)
		return nil, err
	}

	fmt.Println("Conexão estabelecida com sucesso!")

	return db, nil
}
