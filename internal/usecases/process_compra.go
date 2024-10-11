package usecases

import (
	"context"
	"fmt"

	"github.com/thiagohmm/integracaoThothConsumer/internal/domain/entities"
	// Ajustado para importar a interface correta
)

type CompraUseCase struct {
	Repo entities.CompraRepository // Troquei CompraRepositoryDB por CompraRepository (interface)
}

func NewCompraUseCase(repo entities.CompraRepository) *CompraUseCase {
	return &CompraUseCase{Repo: repo}
}

func (uc *CompraUseCase) ProcessarCompra(ctx context.Context, compraData map[string]interface{}) error {
	// Mapeia os dados de compra
	var compra entities.Compra

	// Verifica se a chave "compras" existe no mapa e é um map
	if compras, ok := compraData["compras"].(map[string]interface{}); ok {
		// Extrai o campo "dtaentrada" de "compras"
		if dtaentrada, ok := compras["dtaentrada"].(string); ok {
			compra.Compras.DtaEntrada = dtaentrada
		} else {
			return fmt.Errorf("campo 'dtaentrada' não encontrado ou com tipo incorreto")
		}

		err := uc.Repo.DeleteByIBMAndEntrada(ctx, compra.Compras.Ibms[0].Nro, compra.Compras.DtaEntrada)
		if err != nil {
			return err
		}

		// Extrai a lista de "ibms"
		if ibms, ok := compras["ibms"].([]interface{}); ok {
			for _, ibmData := range ibms {
				// Faz o cast de cada item de "ibms" para map
				if ibmMap, ok := ibmData.(map[string]interface{}); ok {
					var ibm entities.IBM

					// Mapeia o campo "nro"
					if nro, ok := ibmMap["nro"].(string); ok {
						ibm.Nro = nro
					} else {
						return fmt.Errorf("campo 'nro' não encontrado ou com tipo incorreto")
					}

					// Mapeia o campo "razao"
					if razao, ok := ibmMap["razao"].(string); ok {
						ibm.Razao = razao
					} else {
						return fmt.Errorf("campo 'razao' não encontrado ou com tipo incorreto")
					}

					// Adiciona o IBM mapeado à lista de Ibms em Compra
					compra.Compras.Ibms = append(compra.Compras.Ibms, ibm)
				}
			}
		} else {
			return fmt.Errorf("campo 'ibms' não encontrado ou com tipo incorreto")
		}
	} else {
		return fmt.Errorf("campo 'compras' não encontrado ou com tipo incorreto")
	}

	// Exemplo de uso do repositório
	if err := uc.Repo.Save(ctx, compra); err != nil {
		return err
	}

	return nil
}
