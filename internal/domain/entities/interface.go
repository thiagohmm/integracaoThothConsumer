package entities

type CompraRepository interface {
	Save(compra Compra) error
	DeleteByIBMAndEntrada(ibm string, dtaEntrada string) error
}
