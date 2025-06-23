package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"runtime/debug"

	// Add this line to import the strconv package

	"github.com/thiagohmm/integracaoThothConsumer/internal/domain/entities"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("VendaRepositoryDB")

type VendaRepositoryDB struct {
	DB *sql.DB
	Tx *sql.Tx
}

// Conversão para valores nulos para campos de data e numéricos
func nullDate(value string) sql.NullString {
	if value == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: value, Valid: true}
}

func formatNumber(value string) sql.NullFloat64 {
	if value == "" {
		return sql.NullFloat64{Valid: false}
	}
	var num float64
	_, err := fmt.Sscanf(value, "%f", &num)
	if err != nil {
		return sql.NullFloat64{Valid: false}
	}
	return sql.NullFloat64{Float64: num, Valid: true}
}

func (r *VendaRepositoryDB) SalvarVendaComTx(ctx context.Context, tx *sql.Tx, venda entities.Venda, uuid string) error {

	query := `INSERT INTO STG_WS_VENDA (
        CD_CHAVE_TRANSACAO, DT_TRANSACAO, HR_INICIO_TRANSACAO, HR_FIM_TRANSACAO, CD_IBM_LOJA, RAZAO_SOCIAL_LOJA, 
        CD_DEPARTAMENTO, CPF_CNPJ_CLIENTE, CD_EAN_PRODUTO, DS_PRODUTO, CD_TP_PRODUTO, CD_EAN_EMBALAGEM, CD_PROMOCAO, 
        NM_FORMA_PAGAMENTO, NM_BANDEIRA, CD_MODELO_DOCTO, CD_CCF, CD_TRANSACAO, CD_ITEM_TRANSACAO, CD_TP_TRANSACAO, 
        QT_PRODUTO, VL_PRECO_UNITARIO, VL_IMPOSTO, VL_CUSTO_UNITARIO, VL_DESCONTO, VL_FATURADO, NM_SISTEMA, 
        DT_ARQUIVO, SRC_LOAD, DT_LOAD, FL_CARGA_HISTORICA, CD_IBM_LOJA_EAGLE, CD_EAN_PRODUTO_EAGLE, VL_CUSTO_EAGLE, 
        NR_DDD_TELEFONE, NR_TELEFONE, VL_PIS, VL_COFINS, VL_ICMS
    ) VALUES (
        :1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12, :13, :14, :15, :16, :17, :18, :19, :20,
        :21, :22, :23, :24, :25, :26, :27, :28, :29, TO_TIMESTAMP(:30, 'YYYY-MM-DD"T"HH24:MI:SS.FF3"Z"'),
        :31, :32, :33, :34, :35, :36, :37, :38, :39
    )`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {

		return fmt.Errorf("erro ao preparar a query: %w", err)

	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		nullIfEmpty(venda.CD_CHAVE_TRANSACAO),
		venda.DT_TRANSACAO,
		nullIfEmpty(venda.HR_INICIO_TRANSACAO),
		nullIfEmpty(venda.HR_FIM_TRANSACAO),
		nullIfEmpty(venda.CD_IBM_LOJA),
		nullIfEmpty(venda.RAZAO_SOCIAL_LOJA),
		nullIfEmpty(venda.CD_DEPARTAMENTO),
		nullIfEmpty(venda.CPF_CNPJ_CLIENTE),
		nullIfEmpty(venda.CD_EAN_PRODUTO),
		nullIfEmpty(venda.DS_PRODUTO),
		nullIfEmpty(venda.CD_TP_PRODUTO),
		nullIfEmpty(venda.CD_EAN_EMBALAGEM),
		nullIfEmpty(venda.CD_PROMOCAO),
		nullIfEmpty(venda.NM_FORMA_PAGAMENTO),
		nullIfEmpty(venda.NM_BANDEIRA),
		nullIfEmpty(venda.CD_MODELO_DOCTO),
		nullIfEmpty(venda.CD_CCF),
		nullIfEmpty(venda.CD_TRANSACAO),
		nullIfEmpty(venda.CD_ITEM_TRANSACAO),
		nullIfEmpty(venda.CD_TP_TRANSACAO),
		formatNumber(venda.QT_PRODUTO),
		formatNumber(venda.VL_PRECO_UNITARIO),
		formatNumber(venda.VL_IMPOSTO),
		formatNumber(venda.VL_CUSTO_UNITARIO),
		formatNumber(venda.VL_DESCONTO),
		formatNumber(venda.VL_FATURADO),
		nullIfEmpty(venda.NM_SISTEMA),
		venda.DT_ARQUIVO,
		nullIfEmpty(venda.SRC_LOAD),
		venda.DT_LOAD,
		nullIfEmpty(venda.FL_CARGA_HISTORICA),
		nullIfEmpty(venda.CD_IBM_LOJA_EAGLE),
		nullIfEmpty(venda.CD_EAN_PRODUTO_EAGLE),
		formatNumber(venda.VL_CUSTO_EAGLE),
		nullIfEmpty(venda.NR_DDD_TELEFONE),
		nullIfEmpty(venda.NR_TELEFONE),
		venda.VL_PIS,
		venda.VL_COFINS,
		venda.VL_ICMS,
	)

	if err != nil {

		return fmt.Errorf("erro ao executar a query: %w", err)
	}

	log.Printf("Venda salva com sucesso: CD_IBM_LOJA %s, DT_TRANSACAO %d, uuid %s", venda.CD_IBM_LOJA, venda.DT_TRANSACAO, uuid)
	return nil
}

func (r *VendaRepositoryDB) DeleteByIBMVenda(ctx context.Context, ibm string, dtTransacao string, uuid string) error {
	ctx, span := tracer.Start(ctx, uuid)

	query := `DELETE FROM STG_WS_VENDA WHERE CD_IBM_LOJA = :1 AND DT_TRANSACAO = :2`

	_, err := r.DB.ExecContext(ctx, query, ibm, dtTransacao)
	if err != nil {
		// Registra erro no Jaeger (caso span esteja no contexto)
		if span := trace.SpanFromContext(ctx); span != nil {
			span.RecordError(err)
			if ctx.Err() == context.DeadlineExceeded {
				span.SetStatus(codes.Error, "timeout no delete")
				return err
			} else {
				span.SetStatus(codes.Error, "erro ao deletar venda")
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
	span.SetStatus(codes.Ok, "venda deletada com sucesso")
	span.AddEvent("Venda deletada com sucesso", trace.WithAttributes(
		attribute.String("ibm", ibm),
		attribute.String("dtTransacao", dtTransacao),
		attribute.String("uuid", uuid),
	))
	log.Printf("Venda deletada com sucesso: CD_IBM_LOJA %s, DT_TRANSACAO %s, uuid", ibm, dtTransacao, uuid)
	return nil
}

func (r *VendaRepositoryDB) DeleteByIBMVendaZerado(ctx context.Context, ibm string, dtTransacao string, uuid string) error {
	ctx, span := tracer.Start(ctx, uuid)

	query := `DELETE FROM STG_WS_VENDA_ZERADO WHERE CD_IBM_LOJA = :1 AND DT_TRANSACAO = :2`

	_, err := r.DB.ExecContext(ctx, query, ibm, dtTransacao)
	if err != nil {
		// Registra erro no Jaeger (caso span esteja no contexto)
		if span := trace.SpanFromContext(ctx); span != nil {
			span.RecordError(err)
			if ctx.Err() == context.DeadlineExceeded {
				span.SetStatus(codes.Error, "timeout no delete")
				return err
			} else {
				span.SetStatus(codes.Error, "erro ao deletar venda")
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
	span.SetStatus(codes.Ok, "venda deletada com sucesso")
	span.AddEvent("Venda deletada com sucesso", trace.WithAttributes(
		attribute.String("ibm", ibm),
		attribute.String("dtTransacao", dtTransacao),
		attribute.String("uuid", uuid),
	))
	log.Printf("Venda deletada com sucesso: CD_IBM_LOJA %s, DT_TRANSACAO %s, uuid", ibm, dtTransacao, uuid)
	return nil
}

func (r *VendaRepositoryDB) SalvarVendas(ctx context.Context, vendas []entities.Venda, uuid string) error {
	tx, err := r.iniciarTransacao(ctx)
	if err != nil {
		return err
	}
	defer r.garantirRollback(tx)

	var erros []error
	for _, venda := range vendas {
		if err := r.salvarVendaComTx(ctx, tx, venda, uuid); err != nil {
			r.logErroSalvarVenda(err)
			erros = append(erros, fmt.Errorf("erro ao salvar venda %s: %w", venda.CD_IBM_LOJA, err))
			continue // Continua processando as próximas vendas
		}
	}

	// Se houver erros, faz rollback
	if len(erros) > 0 {
		_ = tx.Rollback()
		return fmt.Errorf("erros ao salvar vendas: %v", erros)
	}

	return r.finalizarTransacao(tx)
}

func (r *VendaRepositoryDB) iniciarTransacao(ctx context.Context) (*sql.Tx, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	return tx, nil
}

func (r *VendaRepositoryDB) garantirRollback(tx *sql.Tx) {
	if p := recover(); p != nil {
		_ = tx.Rollback()
		log.Printf("Transação foi revertida devido a um erro: %v\nStacktrace:\n%s", p, debug.Stack())
		// Não propaga o panic, permitindo que a aplicação continue rodando
	}
}

func (r *VendaRepositoryDB) salvarVendaComTx(ctx context.Context, tx *sql.Tx, venda entities.Venda, uuid string) error {
	if err := r.SalvarVendaComTx(ctx, tx, venda, uuid); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("erro ao salvar venda: %w", err)
	}
	return nil
}

func (r *VendaRepositoryDB) logErroSalvarVenda(err error) {
	log.Printf("Erro ao salvar venda: %v", err)
}

func (r *VendaRepositoryDB) finalizarTransacao(tx *sql.Tx) error {
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("erro ao fazer commit: %w", err)
	}
	log.Printf("Vendas salvas com sucesso")
	return nil
}

func (r *VendaRepositoryDB) CheckIfExists(ctx context.Context, ibm string, dtTransacao string) (bool, error) {
	query := `SELECT COUNT(*) FROM STG_WS_VENDA WHERE CD_IBM_LOJA = :1 AND DT_TRANSACAO = :2`
	var count int
	err := r.DB.QueryRowContext(ctx, query, ibm, dtTransacao).Scan(&count)
	if err != nil {
		log.Printf("Erro ao verificar se existe venda: %v", err)
		return false, err
	}
	if count > 0 {
		log.Printf("Venda já existe: CD_IBM_LOJA %s, DT_TRANSACAO %s", ibm, dtTransacao)
		return true, nil
	}
	log.Printf("Venda não existe: CD_IBM_LOJA %s, DT_TRANSACAO %s", ibm, dtTransacao)
	return false, nil
}

func (r *VendaRepositoryDB) CheckIfExistsZerado(ctx context.Context, ibm string, dtTransacao string) (bool, error) {
	query := `SELECT COUNT(*) FROM STG_WS_VENDA_ZERADO WHERE CD_IBM_LOJA = :1 AND DT_TRANSACAO = :2`
	var count int
	err := r.DB.QueryRowContext(ctx, query, ibm, dtTransacao).Scan(&count)
	if err != nil {
		log.Printf("Erro ao verificar se existe venda zerada: %v", err)
		return false, err
	}
	if count > 0 {
		log.Printf("Venda zerada já existe: CD_IBM_LOJA %s, DT_TRANSACAO %s", ibm, dtTransacao)
		return true, nil
	}
	log.Printf("Venda não existe Zerada: CD_IBM_LOJA %s, DT_TRANSACAO %s", ibm, dtTransacao)
	return false, nil
}
