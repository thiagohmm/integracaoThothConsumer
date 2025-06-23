package entities

type Status struct {
	TipoOperacao string `db:"TIPO_OPERACAO" json:"tipo_operacao"` // VARCHAR2(50)
	Ibm          string `db:"IBM" json:"ibm"`                     // VARCHAR2(100)
	DataEnvio    string `db:"DATAENVIO" json:"data_envio"`        // VARCHAR2(30)
	Processo     string `db:"PROCESSO" json:"processo"`           // VARCHAR2(100)
	Status       string `db:"STATUS" json:"status"`               // VARCHAR2(50)
}
