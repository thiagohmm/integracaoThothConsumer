package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type StatusRepositoryDB struct {
	DB *sql.DB
}

func (r *StatusRepositoryDB) UpdateStatusProcesso(ctx context.Context, uuid string, novoStatus string) error {
	query := `UPDATE LOG_ENTRADA SET STATUS = :status WHERE PROCESSO = :uuid`
	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, novoStatus, uuid)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	log.Printf("Status atualizado com sucesso para UUID: %s", uuid)
	return nil
}
