package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/thiagohmm/integracaoThothConsumer/internal/domain/entities"
)

type CompraRepositoryDB struct {
	DB *sql.DB
}

func (r *CompraRepositoryDB) Save(ctx context.Context, compra entities.Compra) error {
	// Inicia uma transação
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Insere os dados na tabela de Compras
	query := `INSERT INTO compras (dtaentrada) VALUES ($1) RETURNING id`
	var compraID int
	if err := tx.QueryRowContext(ctx, query, compra.Compras.DtaEntrada).Scan(&compraID); err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao inserir compra: %w", err)
	}

	// Insere os dados de IBM
	for _, ibm := range compra.Compras.Ibms {
		queryIBM := `INSERT INTO ibms (compra_id, nro, razao) VALUES ($1, $2, $3)`
		if _, err := tx.ExecContext(ctx, queryIBM, compraID, ibm.Nro, ibm.Razao); err != nil {
			tx.Rollback()
			return fmt.Errorf("erro ao inserir IBM: %w", err)
		}
	}

	// Se tudo estiver OK, commit a transação
	return tx.Commit()
}

func (r *CompraRepositoryDB) DeleteByIBMAndEntrada(ctx context.Context, ibm string, dtaEntrada string) error {
	query := `DELETE FROM compras WHERE dtaentrada = $1 AND id IN (SELECT compra_id FROM ibms WHERE nro = $2)`
	_, err := r.DB.Exec(query, dtaEntrada, ibm)
	return err
}
