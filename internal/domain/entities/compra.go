package entities

type Compra struct {
	Compras Compras
}

type Compras struct {
	Ibms       []IBM
	DtaEntrada string
}

type IBM struct {
	Nro   string
	Razao string
	App   string
	Notas []Nota
}

type Nota struct {
	Nro           string
	Serie         string
	Emissao       string
	Fornecedor    Fornecedor
	Total         Total
	Frete         string
	Transportador Transportador
	ChaveXML      string
	Produtos      []Produto
}

type Fornecedor struct {
	Cnpj string
	Nome string
}

type Total struct {
	Peso     string
	VlrFrete string
	VlrIPI   string
	VlrICMS  string
	VlrNota  string
}

type Transportador struct {
	Cnpj string
	Nome string
}

type Produto struct {
	EAN           string
	Quantidade    float64
	Preco         float64
	Descricao     string
	Tipo          string
	Impostos      Impostos
	NCM           string
	Linha         int
	CodFornecedor string
	QuantidadeNF  float64
	UnidadeConv   string
	Unidade       string
	UltCusto      float64
}

type Impostos struct {
	IPI    AliquotaValor
	ICMS   AliquotaValor
	PIS    AliquotaValor
	COFINS AliquotaValor
}

type AliquotaValor struct {
	Aliquota float64
	Valor    float64
}
