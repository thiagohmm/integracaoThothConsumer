package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"runtime/debug"
	"time"

	"github.com/thiagohmm/integracaoThothConsumer/internal/domain/entities"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type EstoqueRepositoryDB struct {
	DB *sql.DB
}

func (r *EstoqueRepositoryDB) DeleteByIBMEstoque(ctx context.Context, ibm string, dtEstoque string, uuid string) error {
	ctx, span := tracer.Start(ctx, uuid)

	query := `DELETE FROM STG_WS_ESTOQUE WHERE CD_IBM_LOJA = :1 AND DT_ESTOQUE = :2`

	_, err := r.DB.ExecContext(ctx, query, ibm, dtEstoque)
	if err != nil {
		// Registra erro no Jaeger (caso span esteja no contexto)
		if span := trace.SpanFromContext(ctx); span != nil {
			span.RecordError(err)
			if ctx.Err() == context.DeadlineExceeded {
				span.SetStatus(codes.Error, "timeout no delete")
				return err
			} else {
				span.SetStatus(codes.Error, "erro ao deletar estoque")
				return err
			}
		}

		if ctx.Err() == context.DeadlineExceeded {
			span.RecordError(err)
			return err

		}
		span.RecordError(err)
		return err
	}
	span.SetStatus(codes.Ok, "Estoque deletada com sucesso")
	span.AddEvent("Estoque deletada com sucesso", trace.WithAttributes(
		attribute.String("ibm", ibm),
		attribute.String("dtTransacao", dtEstoque),
		attribute.String("uuid", uuid),
	))

	log.Printf("Estoque deletado com sucesso: CD_IBM_LOJA %s, DT_ESTOQUE %s", ibm, dtEstoque)
	return nil
}

func (r *EstoqueRepositoryDB) SalvarEstoqueComTx(ctx context.Context, tx *sql.Tx, estoque entities.Estoque, uuid string) error {
	query := `INSERT INTO STG_WS_ESTOQUE (
        DT_ESTOQUE, DT_ULTIMA_COMPRA, CD_IBM_LOJA, RAZAO_SOCIAL_LOJA, CD_EAN_PRODUTO, DS_PRODUTO, CD_TP_PRODUTO, 
        VL_PRECO_UNITARIO, QT_INVENTARIO_ENTRADA, QT_INVENTARIO_SAIDA, QT_INICIAL_PRODUTO, QT_FINAL_PRODUTO, 
        VL_TOTAL_ESTOQUE, VL_CUSTO_MEDIO, NM_SISTEMA, SRC_LOAD, DT_LOAD
    ) VALUES (
        :1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12, :13, :14, :15, :16, :17
    )`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("erro ao preparar a query: %w", err)
	}
	defer stmt.Close()

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
		time.Now().Format("2006-01-02 15:04:05"),
	)

	_, err = stmt.ExecContext(ctx,
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
		time.Now(), // DT_LOAD
	)

	if err != nil {
		return fmt.Errorf("erro ao executar a query: %w", err)
	}

	log.Printf("Estoque salvo com sucesso: CD_IBM_LOJA %s, DT_ESTOQUE %d UUID %s", estoque.CD_IBM_LOJA, estoque.DT_ESTOQUE, uuid)
	return nil
}

func (r *EstoqueRepositoryDB) SalvarEstoques(ctx context.Context, estoques []entities.Estoque, uuid string) error {
	tx, err := r.iniciarTransacao(ctx)
	if err != nil {
		return err
	}
	defer r.garantirRollback(tx)

	var erros []error
	for _, estoque := range estoques {
		if err := r.salvarEstoqueComTx(ctx, tx, estoque, uuid); err != nil {
			r.logErroSalvarEstoque(err)
			erros = append(erros, fmt.Errorf("erro ao salvar estoque %s: %w", estoque.CD_IBM_LOJA, err))
			continue // Continua processando os próximos estoques
		}
	}

	// Se houver erros, faz rollback
	if len(erros) > 0 {
		_ = tx.Rollback()
		return fmt.Errorf("erros ao salvar estoques: %v", erros)
	}

	return r.finalizarTransacao(tx)
}

func (r *EstoqueRepositoryDB) iniciarTransacao(ctx context.Context) (*sql.Tx, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	return tx, nil
}

func (r *EstoqueRepositoryDB) garantirRollback(tx *sql.Tx) {
	if p := recover(); p != nil {
		_ = tx.Rollback()
		log.Printf("Transação foi revertida devido a um erro: %v\nStacktrace:\n%s", p, debug.Stack())
		// Não propaga o panic, permitindo que a aplicação continue rodando
	}
}

func (r *EstoqueRepositoryDB) salvarEstoqueComTx(ctx context.Context, tx *sql.Tx, estoque entities.Estoque, uuid string) error {
	if err := r.SalvarEstoqueComTx(ctx, tx, estoque, uuid); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("erro ao salvar estoque: %w", err)
	}
	return nil
}

func (r *EstoqueRepositoryDB) logErroSalvarEstoque(err error) {
	log.Printf("Erro ao salvar estoque: %v", err)
}

func (r *EstoqueRepositoryDB) finalizarTransacao(tx *sql.Tx) error {
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("erro ao fazer commit: %w", err)
	}
	log.Println("Estoques salvos com sucesso para")
	return nil
}

func (r *EstoqueRepositoryDB) CheckIfExists(ctx context.Context, ibm string, dtEstoque string) (bool, error) {
	query := `SELECT COUNT(*) FROM STG_WS_ESTOQUE WHERE CD_IBM_LOJA = :1 AND DT_ESTOQUE = :2`
	var count int
	err := r.DB.QueryRowContext(ctx, query, ibm, dtEstoque).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar se existe estoque: %w", err)
	}
	return count > 0, nil
}
