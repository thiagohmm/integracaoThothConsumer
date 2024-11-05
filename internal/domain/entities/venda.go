package entities

type Venda struct {
	ID_STG_WS_VENDA      int64   `json:"id_stg_ws_venda,omitempty"`
	CD_CHAVE_TRANSACAO   string  `json:"cd_chave_transacao,omitempty"`
	DT_TRANSACAO         int64   `json:"dt_transacao,omitempty"`
	HR_INICIO_TRANSACAO  string  `json:"hr_inicio_transacao,omitempty"`
	HR_FIM_TRANSACAO     string  `json:"hr_fim_transacao,omitempty"`
	CD_IBM_LOJA          string  `json:"cd_ibm_loja,omitempty"`
	RAZAO_SOCIAL_LOJA    string  `json:"razao_social_loja,omitempty"`
	CD_DEPARTAMENTO      string  `json:"cd_departamento,omitempty"`
	CPF_CNPJ_CLIENTE     string  `json:"cpf_cnpj_cliente,omitempty"`
	CD_EAN_PRODUTO       string  `json:"cd_ean_produto,omitempty"`
	DS_PRODUTO           string  `json:"ds_produto,omitempty"`
	CD_TP_PRODUTO        string  `json:"cd_tp_produto,omitempty"`
	CD_EAN_EMBALAGEM     string  `json:"cd_ean_embalagem,omitempty"`
	CD_PROMOCAO          string  `json:"cd_promocao,omitempty"`
	NM_FORMA_PAGAMENTO   string  `json:"nm_forma_pagamento,omitempty"`
	NM_BANDEIRA          string  `json:"nm_bandeira,omitempty"`
	CD_MODELO_DOCTO      string  `json:"cd_modelo_docto,omitempty"`
	CD_CCF               string  `json:"cd_ccf,omitempty"`
	CD_TRANSACAO         string  `json:"cd_transacao,omitempty"`
	CD_ITEM_TRANSACAO    string  `json:"cd_item_transacao,omitempty"`
	CD_TP_TRANSACAO      string  `json:"cd_tp_transacao,omitempty"`
	QT_PRODUTO           string  `json:"qt_produto,omitempty"`
	VL_PRECO_UNITARIO    string  `json:"vl_preco_unitario,omitempty"`
	VL_IMPOSTO           string  `json:"vl_imposto,omitempty"`
	VL_CUSTO_UNITARIO    string  `json:"vl_custo_unitario,omitempty"`
	VL_DESCONTO          string  `json:"vl_desconto,omitempty"`
	VL_FATURADO          string  `json:"vl_faturado,omitempty"`
	NM_SISTEMA           string  `json:"nm_sistema,omitempty"`
	DT_ARQUIVO           string  `json:"dt_arquivo,omitempty"`
	SRC_LOAD             string  `json:"src_load,omitempty"`
	DT_LOAD              string  `json:"dt_load,omitempty"`
	FL_CARGA_HISTORICA   string  `json:"fl_carga_historica,omitempty"`
	CD_IBM_LOJA_EAGLE    string  `json:"cd_ibm_loja_eagle,omitempty"`
	CD_EAN_PRODUTO_EAGLE string  `json:"cd_ean_produto_eagle,omitempty"`
	VL_CUSTO_EAGLE       string  `json:"vl_custo_eagle,omitempty"`
	NR_DDD_TELEFONE      string  `json:"nr_ddd_telefone,omitempty"`
	NR_TELEFONE          string  `json:"nr_telefone,omitempty"`
	VL_PIS               float64 `json:"vl_pis,omitempty"`
	VL_COFINS            float64 `json:"vl_cofins,omitempty"`
	VL_ICMS              float64 `json:"vl_icms,omitempty"`
}
