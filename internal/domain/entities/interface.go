package entities

import "context"

type CompraRepository interface {
	SaveCompra(ctx context.Context, compra Compra) error
	DeleteByIBMAndEntrada(ctx context.Context, ibm string, dtaEntrada string) error
}
