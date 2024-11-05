package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	// Add this line to import the strconv package
	"github.com/thiagohmm/integracaoThothConsumer/internal/domain/entities"
)

type VendaRepositoryDB struct {
	DB *sql.DB
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

func (r *VendaRepositoryDB) SalvarVenda(ctx context.Context, venda entities.Venda) error {
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

	// Logar os valores antes da inserção
	log.Printf("Valores para inserção: CD_CHAVE_TRANSACAO: %s, DT_TRANSACAO: %d, HR_INICIO_TRANSACAO: %s, HR_FIM_TRANSACAO: %s, CD_IBM_LOJA: %s, RAZAO_SOCIAL_LOJA: %s, CD_DEPARTAMENTO: %s, CPF_CNPJ_CLIENTE: %s, CD_EAN_PRODUTO: %s, DS_PRODUTO: %s, CD_TP_PRODUTO: %s, CD_EAN_EMBALAGEM: %s, CD_PROMOCAO: %s, NM_FORMA_PAGAMENTO: %s, NM_BANDEIRA: %s, CD_MODELO_DOCTO: %s, CD_CCF: %s, CD_TRANSACAO: %s, CD_ITEM_TRANSACAO: %s, CD_TP_TRANSACAO: %s, QT_PRODUTO: %s, VL_PRECO_UNITARIO: %s, VL_IMPOSTO: %s, VL_CUSTO_UNITARIO: %s, VL_DESCONTO: %s, VL_FATURADO: %s, NM_SISTEMA: %s, DT_ARQUIVO: %s, SRC_LOAD: %s, DT_LOAD: %s, FL_CARGA_HISTORICA: %s, CD_IBM_LOJA_EAGLE: %s, CD_EAN_PRODUTO_EAGLE: %s, VL_CUSTO_EAGLE: %s, NR_DDD_TELEFONE: %s, NR_TELEFONE: %s, VL_PIS: %f, VL_COFINS: %f, VL_ICMS: %f",
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
		nullIfEmpty(venda.QT_PRODUTO),
		nullIfEmpty(venda.VL_PRECO_UNITARIO),
		nullIfEmpty(venda.VL_IMPOSTO),
		nullIfEmpty(venda.VL_CUSTO_UNITARIO),
		nullIfEmpty(venda.VL_DESCONTO),
		nullIfEmpty(venda.VL_FATURADO),
		nullIfEmpty(venda.NM_SISTEMA),
		nullIfEmpty(venda.DT_ARQUIVO),
		nullIfEmpty(venda.SRC_LOAD),
		venda.DT_LOAD, // Formatar DT_LOAD corretamente
		nullIfEmpty(venda.FL_CARGA_HISTORICA),
		nullIfEmpty(venda.CD_IBM_LOJA_EAGLE),
		nullIfEmpty(venda.CD_EAN_PRODUTO_EAGLE),
		nullIfEmpty(venda.VL_CUSTO_EAGLE),
		nullIfEmpty(venda.NR_DDD_TELEFONE),
		nullIfEmpty(venda.NR_TELEFONE),
		venda.VL_PIS,
		venda.VL_COFINS,
		venda.VL_ICMS,
	)

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("erro ao preparar a query: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		nullIfEmpty(venda.CD_CHAVE_TRANSACAO),
		venda.DT_TRANSACAO, // Agora em formato 'YYYYMMDD'
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
		venda.DT_ARQUIVO, // Agora em formato 'YYYYMMDD'
		nullIfEmpty(venda.SRC_LOAD),
		venda.DT_LOAD, // Agora formatado como TO_TIMESTAMP
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

	log.Printf("Venda salva com sucesso")
	return nil
}

func (r *VendaRepositoryDB) DeleteByIBMVenda(ctx context.Context, ibm string, dtTransacao string) error {
	query := `DELETE FROM STG_WS_VENDA WHERE CD_IBM_LOJA = :1 AND DT_TRANSACAO = :2`

	_, err := r.DB.ExecContext(ctx, query, ibm, dtTransacao)
	if err != nil {
		return fmt.Errorf("erro ao deletar venda por IBM e Transacao: %w", err)
	}

	log.Printf("Venda deletada com sucesso: CD_IBM_LOJA %s, DT_TRANSACAO %d", ibm, dtTransacao)
	return nil
}
