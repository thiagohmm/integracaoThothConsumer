package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"runtime/debug"

	"github.com/thiagohmm/integracaoThothConsumer/internal/domain/entities"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type CompraRepositoryDB struct {
	DB *sql.DB
}

func (r *CompraRepositoryDB) SaveCompraTx(ctx context.Context, tx *sql.Tx, compra entities.Compra, uuid string) error {
	log.Printf("Valores para inserção: DT_EMISSAO_NOTA: %d, DT_ENTRADA: %d, VL_ULTIMO_CUSTO: %f, QT_PRODUTO: %f, QT_PRODUTO_CONVERTIDA: %f, VL_PRECO_COMPRA: %f, VL_PIS: %f, VL_ALIQUOTA_PIS: %f, VL_COFINS: %f, VL_ALIQUOTA_COFINS: %f, VL_ICMS: %f, VL_ALIQUOTA_ICMS: %f, VL_IPI: %f, VL_ALIQUOTA_IPI: %f, QT_PESO: %f, VL_TOTAL_ICMS: %f, VL_TOTAL_IPI: %f, VL_TOTAL_COMPRA: %f",
		compra.DT_EMISSAO_NOTA, compra.DT_ENTRADA, compra.VL_ULTIMO_CUSTO, compra.QT_PRODUTO, compra.QT_PRODUTO_CONVERTIDA, compra.VL_PRECO_COMPRA, compra.VL_PIS, compra.VL_ALIQUOTA_PIS, compra.VL_COFINS, compra.VL_ALIQUOTA_COFINS, compra.VL_ICMS, compra.VL_ALIQUOTA_ICMS, compra.VL_IPI, compra.VL_ALIQUOTA_IPI, compra.QT_PESO, compra.VL_TOTAL_ICMS, compra.VL_TOTAL_IPI, compra.VL_TOTAL_COMPRA)

	query := `INSERT INTO STG_WS_COMPRA (
        DT_EMISSAO_NOTA, DT_ENTRADA, CD_IBM_LOJA, RAZAO_SOCIAL_LOJA, CNPJ_FORNECEDOR, NM_FORNECEDOR, 
        CNPJ_TRANSPORTADORA, NM_TRANSPORTADORA, CD_PRODUTO_FORNECEDOR, CD_EAN_PRODUTO, DS_PRODUTO, CD_TP_PRODUTO, DS_UN_MEDIDA, 
        DS_UN_MEDIDA_CONVERTIDA, NR_SERIE_NOTA, NR_NOTA_FISCAL, CD_ITEM_NOTA_FISCAL, CD_CHAVE_NOTA_FISCAL, CD_NCM, TIPO_FRETE, 
        VL_ULTIMO_CUSTO, QT_PRODUTO, QT_PRODUTO_CONVERTIDA, VL_PRECO_COMPRA, VL_PIS, VL_ALIQUOTA_PIS, VL_COFINS, VL_ALIQUOTA_COFINS, 
        VL_ICMS, VL_ALIQUOTA_ICMS, VL_IPI, VL_ALIQUOTA_IPI, QT_PESO, VL_FRETE, VL_TOTAL_ICMS, VL_TOTAL_IPI, VL_TOTAL_COMPRA, 
        NM_SISTEMA, SRC_LOAD, DT_LOAD, FL_CARGA_HISTORICA, CD_IBM_LOJA_EAGLE, CD_EAN_PRODUTO_EAGLE
    ) VALUES (
        :1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12, :13, :14, :15, :16, :17, :18, :19, :20, :21, :22, :23, :24, 
        :25, :26, :27, :28, :29, :30, :31, :32, :33, :34, :35, :36, :37, :38, :39, :40, :41, :42, :43
    )`

	_, err := tx.ExecContext(ctx, query,
		compra.DT_EMISSAO_NOTA,
		compra.DT_ENTRADA,
		compra.CD_IBM_LOJA,
		compra.RAZAO_SOCIAL_LOJA,
		compra.CNPJ_FORNECEDOR,
		compra.NM_FORNECEDOR,
		nullIfEmpty(compra.CNPJ_TRANSPORTADORA),
		nullIfEmpty(compra.NM_TRANSPORTADORA),
		nullIfEmpty(compra.CD_PRODUTO_FORNECEDOR),
		compra.CD_EAN_PRODUTO,
		compra.DS_PRODUTO,
		compra.CD_TP_PRODUTO,
		nullIfEmpty(compra.DS_UN_MEDIDA),
		nullIfEmpty(compra.DS_UN_MEDIDA_CONVERTIDA),
		compra.NR_SERIE_NOTA,
		compra.NR_NOTA_FISCAL,
		nullIfEmpty(compra.CD_ITEM_NOTA_FISCAL),
		compra.CD_CHAVE_NOTA_FISCAL,
		nullIfEmpty(compra.CD_NCM),
		nullIfEmpty(compra.TIPO_FRETE),
		compra.VL_ULTIMO_CUSTO,
		compra.QT_PRODUTO,
		compra.QT_PRODUTO_CONVERTIDA,
		compra.VL_PRECO_COMPRA,
		compra.VL_PIS,
		compra.VL_ALIQUOTA_PIS,
		compra.VL_COFINS,
		compra.VL_ALIQUOTA_COFINS,
		compra.VL_ICMS,
		compra.VL_ALIQUOTA_ICMS,
		compra.VL_IPI,
		compra.VL_ALIQUOTA_IPI,
		compra.QT_PESO,
		compra.VL_FRETE,
		compra.VL_TOTAL_ICMS,
		compra.VL_TOTAL_IPI,
		compra.VL_TOTAL_COMPRA,
		compra.NM_SISTEMA,
		compra.SRC_LOAD,
		time.Now(),
		compra.FL_CARGA_HISTORICA,
		nullIfEmpty(compra.CD_IBM_LOJA_EAGLE),
		nullIfEmpty(compra.CD_EAN_PRODUTO_EAGLE),
	)
	if err != nil {
		log.Printf("Erro ao salvar a compra (TX): %v. Valores: DT_EMISSAO_NOTA: %d, DT_ENTRADA: %d, CD_IBM_LOJA: %s, RAZAO_SOCIAL_LOJA: %s, CNPJ_FORNECEDOR: %s, NM_FORNECEDOR: %s, ...", err, compra.DT_EMISSAO_NOTA, compra.DT_ENTRADA, compra.CD_IBM_LOJA, compra.RAZAO_SOCIAL_LOJA, compra.CNPJ_FORNECEDOR, compra.NM_FORNECEDOR)
		return err
	}

	log.Println("Compra (TX) salva com sucesso. UUID: %s", uuid)
	return nil
}

func nullIfEmpty(value string) interface{} {
	if value == "" {
		return nil
	}
	return value
}

func (r *CompraRepositoryDB) DeleteByIBMAndEntrada(ctx context.Context, ibm string, entrada string, uuid string) error {
	ctx, span := tracer.Start(ctx, uuid)

	// Garantir que os parâmetros estejam formatados corretamente
	ibm = strings.TrimSpace(ibm)
	entrada = strings.TrimSpace(entrada)

	query := `DELETE FROM STG_WS_COMPRA WHERE CD_IBM_LOJA = :1 AND DT_ENTRADA = :2`

	_, err := r.DB.ExecContext(ctx, query, ibm, entrada)
	if err != nil {
		// Registra erro no Jaeger (caso span esteja no contexto)
		if span := trace.SpanFromContext(ctx); span != nil {
			span.RecordError(err)
			if ctx.Err() == context.DeadlineExceeded {
				span.SetStatus(codes.Error, "timeout no delete")
				return err
			} else {
				span.SetStatus(codes.Error, "erro ao deletar compra")
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
	span.SetStatus(codes.Ok, "Compra deletada com sucesso")
	span.AddEvent("Compra deletada com sucesso", trace.WithAttributes(
		attribute.String("ibm", ibm),
		attribute.String("dtTransacao", entrada),
		attribute.String("uuid", uuid),
	))
	log.Printf("Compra deletada com sucesso: CD_IBM_LOJA %s, DT_ENTRADA %s, uuid", ibm, entrada, uuid)
	return nil
}

func (r *CompraRepositoryDB) SalvarCompras(ctx context.Context, compras []entities.Compra, uuid string) error {
	tx, err := r.iniciarTransacao(ctx)
	if err != nil {
		return err
	}
	defer r.garantirRollback(tx)

	var erros []error
	for _, compra := range compras {
		if err := r.salvarCompraComTx(ctx, tx, compra, uuid); err != nil {
			r.logErroSalvarCompra(err)
			erros = append(erros, fmt.Errorf("erro ao salvar compra %s: %w", compra.CD_IBM_LOJA, err))
			continue // Continua processando as próximas compras
		}
	}

	// Se houver erros, faz rollback
	if len(erros) > 0 {
		_ = tx.Rollback()
		return fmt.Errorf("erros ao salvar compras: %v", erros)
	}

	return r.finalizarTransacao(tx)
}

func (r *CompraRepositoryDB) iniciarTransacao(ctx context.Context) (*sql.Tx, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	return tx, nil
}

func (r *CompraRepositoryDB) garantirRollback(tx *sql.Tx) {
	if p := recover(); p != nil {
		_ = tx.Rollback()
		log.Printf("Transação foi revertida devido a um erro: %v\nStacktrace:\n%s", p, debug.Stack())
		// Não propaga o panic, permitindo que a aplicação continue rodando
	}
}

func (r *CompraRepositoryDB) salvarCompraComTx(ctx context.Context, tx *sql.Tx, compra entities.Compra, uuid string) error {
	if err := r.SaveCompraTx(ctx, tx, compra, uuid); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("erro ao salvar compra: %w", err)
	}
	return nil
}

func (r *CompraRepositoryDB) logErroSalvarCompra(err error) {
	log.Printf("Erro ao salvar compra: %v", err)
}

func (r *CompraRepositoryDB) finalizarTransacao(tx *sql.Tx) error {
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("erro ao fazer commit: %w", err)
	}
	log.Println("Compras salvas com sucesso.")
	return nil
}

func (r *CompraRepositoryDB) CheckIfExists(ctx context.Context, ibm string, dtEntrada string) (bool, error) {
	query := `SELECT COUNT(*) FROM STG_WS_COMPRA WHERE CD_IBM_LOJA = :1 AND DT_ENTRADA = :2`
	var count int
	err := r.DB.QueryRowContext(ctx, query, ibm, dtEntrada).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar se existe compra: %w", err)
	}
	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}

}
