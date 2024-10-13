package entities

type Compra struct {
	ID_STG_WS_COMPRA        int64   `json:"id_stg_ws_compra,omitempty"`
	DT_EMISSAO_NOTA         int64   `json:"dt_emissao_nota,omitempty"`
	DT_ENTRADA              int64   `json:"dt_entrada,omitempty"`
	CD_IBM_LOJA             string  `json:"cd_ibm_loja,omitempty"`
	RAZAO_SOCIAL_LOJA       string  `json:"razao_social_loja,omitempty"`
	CNPJ_FORNECEDOR         string  `json:"cnpj_fornecedor,omitempty"`
	NM_FORNECEDOR           string  `json:"nm_fornecedor,omitempty"`
	CNPJ_TRANSPORTADORA     string  `json:"cnpj_transportadora,omitempty"`
	NM_TRANSPORTADORA       string  `json:"nm_transportadora,omitempty"`
	CD_PRODUTO_FORNECEDOR   string  `json:"cd_produto_fornecedor,omitempty"`
	CD_EAN_PRODUTO          string  `json:"cd_ean_produto,omitempty"`
	DS_PRODUTO              string  `json:"ds_produto,omitempty"`
	CD_TP_PRODUTO           string  `json:"cd_tp_produto,omitempty"`
	DS_UN_MEDIDA            string  `json:"ds_un_medida,omitempty"`
	DS_UN_MEDIDA_CONVERTIDA string  `json:"ds_un_medida_convertida,omitempty"`
	NR_SERIE_NOTA           string  `json:"nr_serie_nota,omitempty"`
	NR_NOTA_FISCAL          string  `json:"nr_nota_fiscal,omitempty"`
	CD_ITEM_NOTA_FISCAL     string  `json:"cd_item_nota_fiscal,omitempty"`
	CD_CHAVE_NOTA_FISCAL    string  `json:"cd_chave_nota_fiscal,omitempty"`
	CD_NCM                  string  `json:"cd_ncm,omitempty"`
	TIPO_FRETE              string  `json:"tipo_frete,omitempty"`
	VL_ULTIMO_CUSTO         float64 `json:"vl_ultimo_custo,omitempty"`
	QT_PRODUTO              float64 `json:"qt_produto,omitempty"`
	QT_PRODUTO_CONVERTIDA   float64 `json:"qt_produto_convertida,omitempty"`
	VL_PRECO_COMPRA         float64 `json:"vl_preco_compra,omitempty"`
	VL_PIS                  float64 `json:"vl_pis,omitempty"`
	VL_ALIQUOTA_PIS         float64 `json:"vl_aliquota_pis,omitempty"`
	VL_COFINS               float64 `json:"vl_cofins,omitempty"`
	VL_ALIQUOTA_COFINS      float64 `json:"vl_aliquota_cofins,omitempty"`
	VL_ICMS                 float64 `json:"vl_icms,omitempty"`
	VL_ALIQUOTA_ICMS        float64 `json:"vl_aliquota_icms,omitempty"`
	VL_IPI                  float64 `json:"vl_ipi,omitempty"`
	VL_ALIQUOTA_IPI         float64 `json:"vl_aliquota_ipi,omitempty"`
	QT_PESO                 float64 `json:"qt_peso,omitempty"`
	VL_FRETE                float64 `json:"vl_frete,omitempty"`
	VL_TOTAL_ICMS           float64 `json:"vl_total_icms,omitempty"`
	VL_TOTAL_IPI            float64 `json:"vl_total_ipi,omitempty"`
	VL_TOTAL_COMPRA         float64 `json:"vl_total_compra,omitempty"`
	NM_SISTEMA              string  `json:"nm_sistema,omitempty"`
	SRC_LOAD                string  `json:"src_load,omitempty"`
	DT_LOAD                 string  `json:"dt_load,omitempty"`
	FL_CARGA_HISTORICA      int     `json:"fl_carga_historica,omitempty"`
	CD_IBM_LOJA_EAGLE       string  `json:"cd_ibm_loja_eagle,omitempty"`
	CD_EAN_PRODUTO_EAGLE    string  `json:"cd_ean_produto_eagle,omitempty"`
}

type IBM struct {
	Nro   string `json:"nro"`
	Razao string `json:"razao"`
}
