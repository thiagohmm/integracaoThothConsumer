package entities

type Estoque struct {
	ID_STG_WS_ESTOQUE     int64   `json:"id_stg_ws_estoque,omitempty"`
	DT_ESTOQUE            int64   `json:"dt_estoque,omitempty"`
	DT_ULTIMA_COMPRA      int64   `json:"dt_ultima_compra,omitempty"`
	CD_IBM_LOJA           string  `json:"cd_ibm_loja,omitempty"`
	RAZAO_SOCIAL_LOJA     string  `json:"razao_social_loja,omitempty"`
	CD_EAN_PRODUTO        string  `json:"cd_ean_produto,omitempty"`
	DS_PRODUTO            string  `json:"ds_produto,omitempty"`
	CD_TP_PRODUTO         string  `json:"cd_tp_produto,omitempty"`
	VL_PRECO_UNITARIO     float64 `json:"vl_preco_unitario,omitempty"`
	QT_INVENTARIO_ENTRADA float64 `json:"qt_inventario_entrada,omitempty"`
	QT_INVENTARIO_SAIDA   float64 `json:"qt_inventario_saida,omitempty"`
	QT_INICIAL_PRODUTO    float64 `json:"qt_inicial_produto,omitempty"`
	QT_FINAL_PRODUTO      float64 `json:"qt_final_produto,omitempty"`
	VL_TOTAL_ESTOQUE      float64 `json:"vl_total_estoque,omitempty"`
	VL_CUSTO_MEDIO        float64 `json:"vl_custo_medio,omitempty"`
	NM_SISTEMA            string  `json:"nm_sistema,omitempty"`
	SRC_LOAD              string  `json:"src_load,omitempty"`
	DT_LOAD               string  `json:"dt_load,omitempty"`
	FL_CARGA_HISTORICA    int64   `json:"fl_carga_historica,omitempty"`
	CD_IBM_LOJA_EAGLE     string  `json:"cd_ibm_loja_eagle,omitempty"`
	CD_EAN_PRODUTO_EAGLE  string  `json:"cd_ean_produto_eagle,omitempty"`
}
