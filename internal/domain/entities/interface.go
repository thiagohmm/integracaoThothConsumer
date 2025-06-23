package entities

import (
	"context"
	"database/sql"
)

type CompraRepository interface {
	SaveCompraTx(ctx context.Context, tx *sql.Tx, compra Compra, uuid string) error
	SalvarCompras(ctx context.Context, compra []Compra, uuid string) error
	DeleteByIBMAndEntrada(ctx context.Context, ibm string, dtaEntrada string, uuid string) error
	CheckIfExists(ctx context.Context, ibm string, dtaEntrada string) (bool, error)
}

type EstoqueRepository interface {
	SalvarEstoqueComTx(ctx context.Context, tx *sql.Tx, estoque Estoque, uuid string) error
	SalvarEstoques(ctx context.Context, estoque []Estoque, uuid string) error
	DeleteByIBMEstoque(ctx context.Context, ibm string, dtEstoque string, uuid string) error
	CheckIfExists(ctx context.Context, ibm string, dtEstoque string) (bool, error)
}

type VendaRepository interface {
	SalvarVendaComTx(ctx context.Context, tx *sql.Tx, sqlvenda Venda, uuid string) error
	DeleteByIBMVenda(ctx context.Context, ibm string, dtVenda string, uuid string) error
	DeleteByIBMVendaZerado(ctx context.Context, ibm string, dtVenda string, uuid string) error
	SalvarVendas(ctx context.Context, venda []Venda, uuid string) error
	CheckIfExists(ctx context.Context, ibm string, dtVenda string) (bool, error)
	CheckIfExistsZerado(ctx context.Context, ibm string, dtVenda string) (bool, error)
}

type StatusRepository interface {
	UpdateStatusProcesso(ctx context.Context, uuid string, novoStatus string) error
}

type UowRepository interface {
	BeginTransaction(ctx context.Context) (context.Context, error)
	CommitTransaction(ctx context.Context) error
	RollbackTransaction(ctx context.Context) error
	CompraRepository() CompraRepository
	EstoqueRepository() EstoqueRepository
	VendaRepository() VendaRepository
}
