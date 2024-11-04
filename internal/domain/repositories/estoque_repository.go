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

func (r *EstoqueRepositoryDB) SalvarEstoque(ctx context.Context, estoque entities.Estoque) error {
	query := `INSERT INTO STG_WS_ESTOQUE (
        DT_ESTOQUE, DT_ULTIMA_COMPRA, CD_IBM_LOJA, RAZAO_SOCIAL_LOJA, CD_EAN_PRODUTO, DS_PRODUTO, CD_TP_PRODUTO, 
        VL_PRECO_UNITARIO, QT_INVENTARIO_ENTRADA, QT_INVENTARIO_SAIDA, QT_INICIAL_PRODUTO, QT_FINAL_PRODUTO, 
        VL_TOTAL_ESTOQUE, VL_CUSTO_MEDIO, NM_SISTEMA, SRC_LOAD, DT_LOAD
    ) VALUES (
        :1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12, :13, :14, :15, :16, :17
    )`

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("erro ao preparar a query: %w", err)
	}
	defer stmt.Close()

	// Logar os valores antes da inserção
	log.Printf("Valores para inserção: DT_ESTOQUE: %d, DT_ULTIMA_COMPRA: %d, CD_IBM_LOJA: %s, RAZAO_SOCIAL_LOJA: %s, CD_EAN_PRODUTO: %s, DS_PRODUTO: %s, CD_TP_PRODUTO: %s, VL_PRECO_UNITARIO: %f, QT_INVENTARIO_ENTRADA: %f, QT_INVENTARIO_SAIDA: %f, QT_INICIAL_PRODUTO: %f, QT_FINAL_PRODUTO: %f, VL_TOTAL_ESTOQUE: %f, VL_CUSTO_MEDIO: %f, NM_SISTEMA: %s, SRC_LOAD: %s, DT_LOAD: %s",
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
		time.Now().Format("2006-01-02 15:04:05"), // Formatar DT_LOAD corretamente
	)

	_ = stmt.QueryRowContext(ctx,
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
	)
	if err != nil {
		return fmt.Errorf("erro ao executar a query: %w", err)
	}

	log.Printf("Estoque salvo com sucesso: ")
	return nil
}

func (r *EstoqueRepositoryDB) DeleteByIBMEstoque(ctx context.Context, ibm string, dtEstoque string) error {
	query := `DELETE FROM STG_WS_ESTOQUE WHERE CD_IBM_LOJA = :1 AND DT_ESTOQUE = :2`

	_, err := r.DB.ExecContext(ctx, query, ibm, dtEstoque)
	if err != nil {
		return fmt.Errorf("erro ao deletar estoque por IBM e Estoque: %w", err)
	}

	log.Printf("Estoque deletado com sucesso: CD_IBM_LOJA %s, DT_ESTOQUE %d", ibm, dtEstoque)
	return nil
}
