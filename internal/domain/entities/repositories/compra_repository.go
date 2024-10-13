package repositories

import (
	"context"
	"database/sql"
	"log"

	"github.com/thiagohmm/integracaoThothConsumer/internal/domain/entities"
)

type CompraRepositoryDB struct {
	DB *sql.DB
}

func (r *CompraRepositoryDB) Save(ctx context.Context, compra entities.Compra) error {
	query := `INSERT INTO STG_WS_COMPRA (DT_EMISSAO_NOTA, DT_ENTRADA, CD_IBM_LOJA, RAZAO_SOCIAL_LOJA, CNPJ_FORNECEDOR, NM_FORNECEDOR, 
	CNPJ_TRANSPORTADORA, NM_TRANSPORTADORA, CD_PRODUTO_FORNECEDOR, CD_EAN_PRODUTO, DS_PRODUTO, CD_TP_PRODUTO, DS_UN_MEDIDA, 
	DS_UN_MEDIDA_CONVERTIDA, NR_SERIE_NOTA, NR_NOTA_FISCAL, CD_ITEM_NOTA_FISCAL, CD_CHAVE_NOTA_FISCAL, CD_NCM, TIPO_FRETE, 
	VL_ULTIMO_CUSTO, QT_PRODUTO, QT_PRODUTO_CONVERTIDA, VL_PRECO_COMPRA, VL_PIS, VL_ALIQUOTA_PIS, VL_COFINS, VL_ALIQUOTA_COFINS, 
	VL_ICMS, VL_ALIQUOTA_ICMS, VL_IPI, VL_ALIQUOTA_IPI, QT_PESO, VL_FRETE, VL_TOTAL_ICMS, VL_TOTAL_IPI, VL_TOTAL_COMPRA, 
	NM_SISTEMA, SRC_LOAD, DT_LOAD, FL_CARGA_HISTORICA, CD_IBM_LOJA_EAGLE, CD_EAN_PRODUTO_EAGLE) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

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
		compra.TIPO_FRETE,
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
		compra.DT_LOAD,
		compra.FL_CARGA_HISTORICA,
		compra.CD_IBM_LOJA_EAGLE,
		compra.CD_EAN_PRODUTO_EAGLE,
	)

	if err != nil {
		log.Printf("Erro ao salvar a compra: %v", err)
		return err
	}

	log.Println("Compra salva com sucesso.")
	return nil

}

func (r *CompraRepositoryDB) DeleteByIBMAndEntrada(ctx context.Context, ibm string, dtaEntrada string) error {
	query := `DELETE FROM compras WHERE dtaentrada = $1 AND id IN (SELECT compra_id FROM ibms WHERE nro = $2)`
	_, err := r.DB.Exec(query, dtaEntrada, ibm)
	return err
}
