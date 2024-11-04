package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/thiagohmm/integracaoThothConsumer/internal/domain/entities"
	"github.com/valyala/fastjson"
)

type EstoqueUseCase struct {
	Repo entities.EstoqueRepository
}

func NewEstoqueUseCase(repo entities.EstoqueRepository) *EstoqueUseCase {
	return &EstoqueUseCase{Repo: repo}
}

func (uc *EstoqueUseCase) ProcessarEstoque(ctx context.Context, estoqueData map[string]interface{}) error {

	estoqueJSON, err := json.Marshal(estoqueData)
	if err != nil {
		log.Printf("Erro ao converter mapa em JSON: %v", err)
		return err
	}

	// Parsear o JSON usando fastjson
	var p fastjson.Parser
	v, err := p.ParseBytes(estoqueJSON)
	if err != nil {
		log.Printf("Erro ao parsear JSON: %v", err)
		return err
	}

	deleteCounter := 0
	// Itera sobre as IBMs da compra
	ibms := v.GetArray("estoque", "ibms")
	if ibms == nil {
		return fmt.Errorf("IBMs não encontrados no objeto de compra")
	}

	for _, estoqueIbms := range ibms {
		nroStr := sanitizeIBM(string(estoqueIbms.GetStringBytes("nro")))

		// Obtém a data de entrada
		dtaentrada := string(v.GetStringBytes("estoque", "dtaestoque"))
		dtStr := strings.ReplaceAll(dtaentrada, "-", "")

		// Deleta o IBM com base no número e data
		log.Printf("Deletando IBM compra: %s, data: %s", nroStr, dtStr)
		if err := uc.Repo.DeleteByIBMEstoque(ctx, nroStr, dtStr); err != nil {
			log.Printf("Erro ao deletar IBM compra: %s, erro: %v", nroStr, err)
			continue
		}

		deleteCounter++
		if deleteCounter%100 == 0 {
			time.Sleep(500 * time.Millisecond)
		}
	}

	saveCounter := 0
	// Itera sobre as IBMs para salvar
	for _, ibm := range ibms {
		err := func() error {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Erro ao processar IBM Compra: %v", r)
				}
			}()

			nroStr := sanitizeIBM(string(ibm.GetStringBytes("nro")))

			// Obtém a data de emissão da nota
			dtStr := string(v.GetStringBytes("estoque", "dtaestoque"))
			dtEstoque, err := strconv.ParseInt(dtStr, 10, 64)
			if err != nil {
				log.Printf("Erro ao converter DT_Estoque para int64: %v", err)
				return err
			}
			novoIbm := entities.Estoque{
				CD_IBM_LOJA:       nroStr,
				RAZAO_SOCIAL_LOJA: stringOrDefault(ibm.GetStringBytes("razao")),
				DT_ESTOQUE:        dtEstoque,
				NM_SISTEMA:        stringOrDefault(ibm.GetStringBytes("app")),
				SRC_LOAD:          "API/Integração/Thoth",
				DT_LOAD:           string(time.Now().UTC().Format("2006-01-02T15:04:05.000Z")), // Formatar DT_LOAD corretamente
			}

			produtos := ibm.GetArray("produtos")
			if len(produtos) == 0 {
				fmt.Println("Salvando", novoIbm)
				if err := uc.Repo.SalvarEstoque(ctx, novoIbm); err != nil {
					return err
				}
			}

			for _, produto := range produtos {
				novoIbm.CD_EAN_PRODUTO = stringOrDefault(produto.GetStringBytes("ean"))
				novoIbm.CD_TP_PRODUTO = stringOrDefault(produto.GetStringBytes("tipo"))
				novoIbm.DT_ULTIMA_COMPRA = int64(parseFloatOrDefault(produto.Get("dtacompra").GetStringBytes()))
				novoIbm.DS_PRODUTO = stringOrDefault(produto.GetStringBytes("descricao"))
				novoIbm.VL_PRECO_UNITARIO = parseFloat(produto.Get("preco"))
				novoIbm.QT_INVENTARIO_ENTRADA = parseFloat(produto.Get("qtdentrada"))
				novoIbm.QT_INVENTARIO_SAIDA = parseFloat(produto.Get("qtdsaida"))
				novoIbm.QT_INICIAL_PRODUTO = parseFloat(produto.Get("qtdini"))
				novoIbm.QT_FINAL_PRODUTO = parseFloat(produto.Get("qtdfim"))
				novoIbm.VL_TOTAL_ESTOQUE = parseFloat(produto.Get("total"))

				// Tratamento de VL_CUSTO_MEDIO
				vlCustoMedio := parseFloat(produto.Get("vlrmedio"))
				if vlCustoMedio > float64(1<<53-1) { // Número máximo seguro para inteiros em JavaScript (Number.MAX_SAFE_INTEGER)
					vlCustoMedioStr := fmt.Sprintf("%.2f", vlCustoMedio)
					vlCustoMedio, _ = strconv.ParseFloat(vlCustoMedioStr[:4], 64)
				}
				novoIbm.VL_CUSTO_MEDIO = vlCustoMedio

				// Salvar IBM atualizado
				//fmt.Println("Salvando", novoIbm)
				if err := uc.Repo.SalvarEstoque(ctx, novoIbm); err != nil {
					log.Printf("Erro ao salvar novo IBM Estoque: %v", err)
				}
				saveCounter++
				if saveCounter%100 == 0 {
					time.Sleep(500 * time.Millisecond)
				}
			}

			return nil
		}()
		if err != nil {
			log.Printf("Erro ao processar IBM: %v", err)
		}
	}
	return nil

}
