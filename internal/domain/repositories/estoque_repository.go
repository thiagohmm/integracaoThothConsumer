package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/thiagohmm/integracaoThothConsumer/internal/domain/entities"
)

type EstoqueRepositoryDB struct {
	DB *sql.DB
}

func (r *EstoqueRepositoryDB) Save(ctx context.Context, estoque entities.Estoque) error {
	query := `INSERT INTO STG_WS_ESTOQUE (
        DT_ESTOQUE, DT_ULTIMA_COMPRA, CD_IBM_LOJA, RAZAO_SOCIAL_LOJA, CD_EAN_PRODUTO, DS_PRODUTO, CD_TP_PRODUTO, 
        VL_PRECO_UNITARIO, QT_INVENTARIO_ENTRADA, QT_INVENTARIO_SAIDA, QT_INICIAL_PRODUTO, QT_FINAL_PRODUTO, 
        VL_TOTAL_ESTOQUE, VL_CUSTO_MEDIO, NM_SISTEMA, SRC_LOAD, DT_LOAD
    ) VALUES (
        :1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12, :13, :14, :15, :16, :17
    ) RETURNING ID_STG_WS_ESTOQUE, DT_ESTOQUE, DT_ULTIMA_COMPRA, CD_IBM_LOJA, RAZAO_SOCIAL_LOJA, CD_EAN_PRODUTO, 
    DS_PRODUTO, CD_TP_PRODUTO, VL_PRECO_UNITARIO, QT_INVENTARIO_ENTRADA, QT_INVENTARIO_SAIDA, QT_INICIAL_PRODUTO, 
    QT_FINAL_PRODUTO, VL_TOTAL_ESTOQUE, VL_CUSTO_MEDIO, NM_SISTEMA, SRC_LOAD, DT_LOAD, FL_CARGA_HISTORICA, 
    CD_IBM_LOJA_EAGLE, CD_EAN_PRODUTO_EAGLE INTO :18, :19, :20, :21, :22, :23, :24, :25, :26, :27, :28, :29, :30, :31, 
    :32, :33, :34, :35, :36, :37, :38`

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("erro ao preparar a query: %w", err)
	}
	defer stmt.Close()

	var id int64
	err = stmt.QueryRowContext(ctx,
		estoque.DT_ESTOQUE,
		estoque.DT_ULTIMA_COMPRA,
		estoque.CD_IBM_LOJA,
		estoque.RAZAO_SOCIAL_LOJA,
		estoque.CD_EAN_PRODUTO,
		estoque.DS_PRODUTO,
		estoque.CD_TP_PRODUTO,
		estoque.VL_PRECO_UNITARIO,
		estoque.QT_INVENTARIO_ENTRADA,
		estoque.QT_INVENTARIO_SAIDA,
		estoque.QT_INICIAL_PRODUTO,
		estoque.QT_FINAL_PRODUTO,
		estoque.VL_TOTAL_ESTOQUE,
		estoque.VL_CUSTO_MEDIO,
		estoque.NM_SISTEMA,
		estoque.SRC_LOAD,
		time.Now(), // Passando time.Time diretamente para DT_LOAD
	).Scan(
		&id,
		&estoque.DT_ESTOQUE,
		&estoque.DT_ULTIMA_COMPRA,
		&estoque.CD_IBM_LOJA,
		&estoque.RAZAO_SOCIAL_LOJA,
		&estoque.CD_EAN_PRODUTO,
		&estoque.DS_PRODUTO,
		&estoque.CD_TP_PRODUTO,
		&estoque.VL_PRECO_UNITARIO,
		&estoque.QT_INVENTARIO_ENTRADA,
		&estoque.QT_INVENTARIO_SAIDA,
		&estoque.QT_INICIAL_PRODUTO,
		&estoque.QT_FINAL_PRODUTO,
		&estoque.VL_TOTAL_ESTOQUE,
		&estoque.VL_CUSTO_MEDIO,
		&estoque.NM_SISTEMA,
		&estoque.SRC_LOAD,
		&estoque.DT_LOAD,
		&estoque.FL_CARGA_HISTORICA,
		&estoque.CD_IBM_LOJA_EAGLE,
		&estoque.CD_EAN_PRODUTO_EAGLE,
	)
	if err != nil {
		return fmt.Errorf("erro ao executar a query: %w", err)
	}

	log.Printf("Estoque salvo com sucesso: ID %d", id)
	return nil
}

func (r *EstoqueRepositoryDB) DeleteByIBMEstoque(ctx context.Context, ibm string, dtEstoque int64) error {
	query := `DELETE FROM STG_WS_ESTOQUE WHERE CD_IBM_LOJA = :1 AND DT_ESTOQUE = :2`

	_, err := r.DB.ExecContext(ctx, query, ibm, dtEstoque)
	if err != nil {
		return fmt.Errorf("erro ao deletar estoque por IBM e Estoque: %w", err)
	}

	log.Printf("Estoque deletado com sucesso: CD_IBM_LOJA %s, DT_ESTOQUE %d", ibm, dtEstoque)
	return nil
}
