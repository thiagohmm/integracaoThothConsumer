package entities

import "context"

type CompraRepository interface {
	SaveCompra(ctx context.Context, compra Compra) error
	DeleteByIBMAndEntrada(ctx context.Context, ibm string, dtaEntrada string) error
}

type EstoqueRepository interface {
	SalvarEstoque(ctx context.Context, estoque Estoque) error
	DeleteByIBMEstoque(ctx context.Context, ibm string, dtEstoque string) error
}

type VendaRepository interface {
	SalvarVenda(ctx context.Context, venda Venda) error
	DeleteByIBMVenda(ctx context.Context, ibm string, dtVenda string) error
}
