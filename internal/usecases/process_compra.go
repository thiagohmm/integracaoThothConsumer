package usecases

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/thiagohmm/integracaoThothConsumer/internal/domain/entities"
	"github.com/valyala/fastjson"
)

type CompraUseCase struct {
	Repo entities.CompraRepository
}

func NewCompraUseCase(repo entities.CompraRepository) *CompraUseCase {
	return &CompraUseCase{Repo: repo}
}

// Função auxiliar para converter data
func parseDate(data string) int {
	dateInt, _ := strconv.Atoi(strings.ReplaceAll(data, "-", ""))
	return dateInt
}

func parseFloat(value *fastjson.Value) float64 {
	if value == nil {
		return 0.0
	}
	floatValue, err := strconv.ParseFloat(string(value.GetStringBytes()), 64)
	if err != nil {
		return 0.0
	}
	return floatValue
}

func (uc *CompraUseCase) ProcessarCompra(ctx context.Context, compraData map[string]interface{}) error {
	var p fastjson.Parser
	// Verifica se "compra" é uma string JSON ou um objeto
	compraValue, ok := compraData["compra"].(string)
	if !ok {
		log.Printf("A variável 'compra' não é uma string JSON válida.")
		return fmt.Errorf("formato de compra inválido")
	}

	// Parseia a string JSON usando fastjson
	v, err := p.Parse(compraValue)
	if err != nil {
		log.Printf("Erro ao fazer parse da string JSON: %v", err)
		return err
	}

	deleteCounter := 0
	// Itera sobre as IBMs da compra
	ibms := v.GetArray("compras", "ibms")
	if ibms == nil {
		return fmt.Errorf("IBMs não encontrados no objeto de compra")
	}

	for _, compraIbms := range ibms {
		nroStr := strings.TrimSpace(string(compraIbms.GetStringBytes("nro")))
		nroStr = strings.TrimLeft(nroStr, "0")
		if len(nroStr) > 10 {
			nroStr = nroStr[len(nroStr)-10:]
		} else {
			nroStr = fmt.Sprintf("%010s", nroStr)
		}

		// Obtém a data de entrada
		dtaentrada := string(v.GetStringBytes("compras", "dtaentrada"))
		dtStr := strings.ReplaceAll(dtaentrada, "-", "")

		// Deleta o IBM com base no número e data
		log.Printf("Deletando IBM compra: %s, data: %s", nroStr, dtStr)
		if err := uc.Repo.DeleteByIBMAndEntrada(ctx, nroStr, dtStr); err != nil {
			log.Printf("Erro ao deletar IBM compra: %v, erro: %v", compraIbms.GetStringBytes("nro"), err)
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

			nroStr := string(ibm.GetStringBytes("nro"))
			nroStr = strings.TrimLeft(strings.TrimSpace(nroStr), "0")
			if len(nroStr) > 10 {
				nroStr = nroStr[len(nroStr)-10:]
			} else {
				nroStr = fmt.Sprintf("%010s", nroStr)
			}

			novoIbm := entities.Compra{
				CD_IBM_LOJA:       nroStr,
				RAZAO_SOCIAL_LOJA: string(ibm.GetStringBytes("razao")),
				DT_ENTRADA:        int64(parseDate(string(v.GetStringBytes("compras", "dtaentrada")))),
				NM_SISTEMA:        string(ibm.GetStringBytes("app")),
				SRC_LOAD:          "API/Integração/Thoth",
				DT_LOAD:           time.Now().Format(time.RFC3339),
			}

			notas := ibm.GetArray("notas")
			if len(notas) == 0 {
				if err := uc.Repo.Save(ctx, novoIbm); err != nil {
					return err
				}
				saveCounter++
				if saveCounter%100 == 0 {
					time.Sleep(500 * time.Millisecond)
				}
				return nil
			}

			for _, nota := range notas {
				novoIbm.NR_NOTA_FISCAL = string(nota.GetStringBytes("nro"))
				novoIbm.NR_SERIE_NOTA = string(nota.GetStringBytes("serie"))
				novoIbm.DT_EMISSAO_NOTA = int64(parseDate(string(nota.GetStringBytes("emissao"))))
				novoIbm.CNPJ_FORNECEDOR = string(nota.GetStringBytes("fornecedor", "cnpj"))
				novoIbm.NM_FORNECEDOR = string(nota.GetStringBytes("fornecedor", "nome"))
				novoIbm.QT_PESO = parseFloat(nota.Get("total", "peso"))
				novoIbm.VL_TOTAL_IPI = parseFloat(nota.Get("total", "vlripi"))
				novoIbm.VL_TOTAL_ICMS = parseFloat(nota.Get("total", "vlricms"))
				novoIbm.VL_TOTAL_COMPRA = parseFloat(nota.Get("total", "vlrnota"))
				novoIbm.CNPJ_TRANSPORTADORA = string(nota.GetStringBytes("transportador", "cnpj"))
				novoIbm.NM_TRANSPORTADORA = string(nota.GetStringBytes("transportador", "nome"))
				novoIbm.CD_CHAVE_NOTA_FISCAL = string(nota.GetStringBytes("chavexml"))

				// Processando produtos
				produtos := nota.GetArray("produtos")
				for _, produto := range produtos {
					novoIbm.CD_EAN_PRODUTO = string(produto.GetStringBytes("ean"))
					novoIbm.QT_PRODUTO = parseFloat(produto.Get("qtd"))
					novoIbm.VL_PRECO_COMPRA = parseFloat(produto.Get("preco"))
					novoIbm.DS_PRODUTO = string(produto.GetStringBytes("descricao"))
					novoIbm.CD_TP_PRODUTO = string(produto.GetStringBytes("tipo"))
					novoIbm.VL_ALIQUOTA_IPI = parseFloat(produto.Get("impostos", "ipi", "aliquota"))
					novoIbm.VL_IPI = parseFloat(produto.Get("impostos", "ipi", "vlr"))
					novoIbm.VL_ALIQUOTA_ICMS = parseFloat(produto.Get("impostos", "icms", "aliquota"))
					novoIbm.VL_ICMS = parseFloat(produto.Get("impostos", "icms", "vlr"))
					novoIbm.VL_ALIQUOTA_PIS = parseFloat(produto.Get("impostos", "pis", "aliquota"))
					novoIbm.VL_PIS = parseFloat(produto.Get("impostos", "pis", "vlr"))
					novoIbm.VL_ALIQUOTA_COFINS = parseFloat(produto.Get("impostos", "cofins", "aliquota"))
					novoIbm.VL_COFINS = parseFloat(produto.Get("impostos", "cofins", "vlr"))
					novoIbm.CD_NCM = string(produto.GetStringBytes("ncm"))
					novoIbm.CD_ITEM_NOTA_FISCAL = string(produto.GetStringBytes("linha"))
					novoIbm.CD_PRODUTO_FORNECEDOR = string(produto.GetStringBytes("codfornec"))
					novoIbm.QT_PRODUTO_CONVERTIDA = parseFloat(produto.Get("qtdenf"))
					novoIbm.DS_UN_MEDIDA_CONVERTIDA = string(produto.GetStringBytes("unconv"))
					novoIbm.DS_UN_MEDIDA = string(produto.GetStringBytes("un"))
					novoIbm.VL_ULTIMO_CUSTO = parseFloat(produto.Get("ultcusto"))

					if err := uc.Repo.Save(ctx, novoIbm); err != nil {
						return err
					}
					saveCounter++
					if saveCounter%100 == 0 {
						time.Sleep(500 * time.Millisecond)
					}
				}
			}
			return nil
		}()
		if err != nil {
			log.Printf("Erro ao salvar IBM compra: %v", err)
			return err
		}
	}

	return nil
}
