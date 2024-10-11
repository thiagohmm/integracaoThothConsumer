package entities

import "context"

type CompraRepository interface {
	Save(ctx context.Context, compra Compra) error
	DeleteByIBMAndEntrada(ctx context.Context, ibm string, dtaEntrada string) error
}
