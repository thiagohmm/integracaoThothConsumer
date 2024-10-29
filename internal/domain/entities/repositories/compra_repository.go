package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log" // Add the import statement for strconv
	"strconv"
	"strings"
	"time"

	"github.com/thiagohmm/integracaoThothConsumer/internal/domain/entities"
)

type CompraRepositoryDB struct {
	DB *sql.DB
}

func (r *CompraRepositoryDB) Save(ctx context.Context, compra entities.Compra) error {

	query := `INSERT INTO STG_WS_COMPRA (DT_EMISSAO_NOTA, DT_ENTRADA, CD_IBM_LOJA, RAZAO_SOCIAL_LOJA, CNPJ_FORNECEDOR, NM_FORNECEDOR, 
    CNPJ_TRANSPORTADORA, NM_TRANSPORTADORA, CD_PRODUTO_FORNECEDOR, CD_EAN_PRODUTO, DS_PRODUTO, CD_TP_PRODUTO, DS_UN_MEDIDA, 
    DS_UN_MEDIDA_CONVERTIDA, NR_SERIE_NOTA, NR_NOTA_FISCAL, CD_ITEM_NOTA_FISCAL, CD_CHAVE_NOTA_FISCAL, CD_NCM, 
    VL_ULTIMO_CUSTO, QT_PRODUTO, QT_PRODUTO_CONVERTIDA, VL_PRECO_COMPRA, VL_PIS, VL_ALIQUOTA_PIS, VL_COFINS, VL_ALIQUOTA_COFINS, 
    VL_ICMS, VL_ALIQUOTA_ICMS, VL_IPI, VL_ALIQUOTA_IPI, QT_PESO, VL_TOTAL_ICMS, VL_TOTAL_IPI, VL_TOTAL_COMPRA, 
    NM_SISTEMA, SRC_LOAD, DT_LOAD) 
    VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12, :13, :14, :15, :16, :17, :18, :19, :20, :21, :22, :23, :24, 
            :25, :26, :27, :28, :29, :30, :31, :32, :33, :34, :35, :36, :37, :38)`

	_, err := r.DB.ExecContext(ctx, query,
		compra.DT_EMISSAO_NOTA,
		compra.DT_ENTRADA,
		compra.CD_IBM_LOJA,
		compra.RAZAO_SOCIAL_LOJA,
		compra.CNPJ_FORNECEDOR,
		compra.NM_FORNECEDOR,
		compra.CNPJ_TRANSPORTADORA,
		compra.NM_TRANSPORTADORA,
		compra.CD_PRODUTO_FORNECEDOR,
		compra.CD_EAN_PRODUTO,
		compra.DS_PRODUTO,
		compra.CD_TP_PRODUTO,
		compra.DS_UN_MEDIDA,
		compra.DS_UN_MEDIDA_CONVERTIDA,
		compra.NR_SERIE_NOTA,
		compra.NR_NOTA_FISCAL,
		compra.CD_ITEM_NOTA_FISCAL,
		compra.CD_CHAVE_NOTA_FISCAL,
		compra.CD_NCM,
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
		compra.VL_TOTAL_ICMS,
		compra.VL_TOTAL_IPI,
		compra.VL_TOTAL_COMPRA,
		compra.NM_SISTEMA,
		compra.SRC_LOAD,
		compra.DT_LOAD,
	)
	if err != nil {
		log.Printf("Erro ao salvar a compra: %v. Valores: DT_EMISSAO_NOTA: %s, DT_ENTRADA: %s, CD_IBM_LOJA: %s, RAZAO_SOCIAL_LOJA: %s, CNPJ_FORNECEDOR: %s, NM_FORNECEDOR: %s, CNPJ_TRANSPORTADORA: %s, NM_TRANSPORTADORA: %s, CD_PRODUTO_FORNECEDOR: %s, CD_EAN_PRODUTO: %s, DS_PRODUTO: %s, CD_TP_PRODUTO: %s, DS_UN_MEDIDA: %s, DS_UN_MEDIDA_CONVERTIDA: %s, NR_SERIE_NOTA: %s, NR_NOTA_FISCAL: %s, CD_ITEM_NOTA_FISCAL: %s, CD_CHAVE_NOTA_FISCAL: %s, CD_NCM: %s, VL_ULTIMO_CUSTO: %f, QT_PRODUTO: %f, QT_PRODUTO_CONVERTIDA: %f, VL_PRECO_COMPRA: %f, VL_PIS: %f, VL_ALIQUOTA_PIS: %f, VL_COFINS: %f, VL_ALIQUOTA_COFINS: %f, VL_ICMS: %f, VL_ALIQUOTA_ICMS: %f, VL_IPI: %f, VL_ALIQUOTA_IPI: %f, QT_PESO: %f, VL_TOTAL_ICMS: %f, VL_TOTAL_IPI: %f, VL_TOTAL_COMPRA: %f, NM_SISTEMA: %s, SRC_LOAD: %s, DT_LOAD: %s",
			err, strconv.FormatInt(compra.DT_EMISSAO_NOTA, 10), strconv.FormatInt(compra.DT_ENTRADA, 10), compra.CD_IBM_LOJA, compra.RAZAO_SOCIAL_LOJA, compra.CNPJ_FORNECEDOR, compra.NM_FORNECEDOR, compra.CNPJ_TRANSPORTADORA, compra.NM_TRANSPORTADORA, compra.CD_PRODUTO_FORNECEDOR, compra.CD_EAN_PRODUTO, compra.DS_PRODUTO, compra.CD_TP_PRODUTO, compra.DS_UN_MEDIDA, compra.DS_UN_MEDIDA_CONVERTIDA, compra.NR_SERIE_NOTA, compra.NR_NOTA_FISCAL, compra.CD_ITEM_NOTA_FISCAL, compra.CD_CHAVE_NOTA_FISCAL, compra.CD_NCM, compra.VL_ULTIMO_CUSTO, compra.QT_PRODUTO, compra.QT_PRODUTO_CONVERTIDA, compra.VL_PRECO_COMPRA, compra.VL_PIS, compra.VL_ALIQUOTA_PIS, compra.VL_COFINS, compra.VL_ALIQUOTA_COFINS, compra.VL_ICMS, compra.VL_ALIQUOTA_ICMS, compra.VL_IPI, compra.VL_ALIQUOTA_IPI, compra.QT_PESO, compra.VL_TOTAL_ICMS, compra.VL_TOTAL_IPI, compra.VL_TOTAL_COMPRA, compra.NM_SISTEMA, compra.SRC_LOAD, compra.DT_LOAD)
		return err
	}

	log.Println("Compra salva com sucesso.")
	return nil
}

func parseDate(data string) (string, error) {
	if data == "" {
		return "", fmt.Errorf("data vazia")
	}
	parsedDate, err := time.Parse("20060102", data) // espera o formato YYYYMMDD
	if err != nil {
		return "", err
	}
	return parsedDate.Format("2006-01-02"), nil // Formato YYYY-MM-DD
}

func (r *CompraRepositoryDB) DeleteByIBMAndEntrada(ctx context.Context, ibm string, entrada string) error {
	// Garantir que os parâmetros estejam formatados corretamente
	ibm = strings.TrimSpace(ibm)
	entrada = strings.TrimSpace(entrada)

	query := `DELETE FROM STG_WS_COMPRA WHERE CD_IBM_LOJA = :1 AND DT_ENTRADA = :2`

	_, err := r.DB.ExecContext(ctx, query, ibm, entrada)
	if err != nil {
		return fmt.Errorf("erro ao deletar compra por IBM e Entrada: %w", err)
	}
	return nil
}
