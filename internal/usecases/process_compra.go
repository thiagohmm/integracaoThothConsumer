package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
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

func parseFloat(value *fastjson.Value) float64 {
	if value == nil {
		return 0.0
	}
	floatValue, err := strconv.ParseFloat(string(value.GetStringBytes()), 64)
	if err != nil {
		return 0.0
	}
	// Formatar o valor com 3 casas decimais
	formattedValue, err := strconv.ParseFloat(fmt.Sprintf("%.3f", floatValue), 64)
	if err != nil {
		return 0.0
	}
	return formattedValue
}

func parseDate(data string) (string, error) {
	if data == "" {
		return "", fmt.Errorf("data vazia")
	}
	parsedDate, err := time.Parse("20060102", data)
	if err != nil {
		return "", err
	}
	return parsedDate.Format("2006-01-02"), nil
}

func sanitizeIBM(ibm string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, ibm)
}

func stringOrDefault(value []byte) string {
	if value == nil {
		return ""
	}
	return string(value)
}

func parseFloatOrDefault(value []byte) float64 {
	if value == nil {
		return 0.0
	}
	floatValue, err := strconv.ParseFloat(string(value), 64)
	if err != nil {
		return 0.0
	}
	return floatValue
}

func (uc *CompraUseCase) ProcessarCompra(ctx context.Context, compraData map[string]interface{}) error {
	compraJSON, err := json.Marshal(compraData)
	if err != nil {
		log.Printf("Erro ao converter mapa em JSON: %v", err)
		return err
	}

	// Parsear o JSON usando fastjson
	var p fastjson.Parser
	v, err := p.ParseBytes(compraJSON)
	if err != nil {
		log.Printf("Erro ao parsear JSON: %v", err)
		return err
	}

	// Limite de goroutines ativas
	const maxGoroutines = 10
	semaphore := make(chan struct{}, maxGoroutines)
	var wg sync.WaitGroup

	deleteCounter := 0
	// Itera sobre as IBMs da compra
	ibms := v.GetArray("compras", "ibms")
	if ibms == nil {
		return fmt.Errorf("IBMs não encontrados no objeto de compra")
	}

	for _, compraIbms := range ibms {

		wg.Add(1)
		semaphore <- struct{}{}

		go func(estoqueIbms *fastjson.Value) {
			defer wg.Done()
			defer func() { <-semaphore }() // Libera o slot no semáforo

			nroStr := sanitizeIBM(string(compraIbms.GetStringBytes("nro")))

			// Obtém a data de entrada
			dtaentrada := string(v.GetStringBytes("compras", "dtaentrada"))
			dtStr := strings.ReplaceAll(dtaentrada, "-", "")

			// Deleta o IBM com base no número e data
			log.Printf("Deletando IBM compra: %s, data: %s", nroStr, dtStr)
			if err := uc.Repo.DeleteByIBMAndEntrada(ctx, nroStr, dtStr); err != nil {
				log.Printf("Erro ao deletar IBM compra: %s, erro: %v", nroStr, err)

			}

			deleteCounter++
			if deleteCounter%100 == 0 {
				time.Sleep(500 * time.Millisecond)
			}
		}(compraIbms)
	}

	// Espera todas as goroutines de deleção terminarem
	wg.Wait()

	// Processa e salva as IBMs em goroutines
	saveCounter := 0
	// Itera sobre as IBMs para salvar
	for _, ibm := range ibms {
		wg.Add(1)
		semaphore <- struct{}{}

		func(ibm *fastjson.Value) {
			defer wg.Done()
			defer func() { <-semaphore }()

			err := func() error {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Erro ao processar IBM Compra: %v", r)
					}
				}()

				nroStr := sanitizeIBM(string(ibm.GetStringBytes("nro")))

				// Obtém a data de emissão da nota
				dtStr := string(v.GetStringBytes("compras", "dtaentrada"))
				dtEntrada, err := strconv.ParseInt(dtStr, 10, 64)
				if err != nil {
					log.Printf("Erro ao converter DT_ENTRADA para int64: %v", err)
					return err
				}
				novoIbm := entities.Compra{
					CD_IBM_LOJA:       nroStr,
					RAZAO_SOCIAL_LOJA: stringOrDefault(ibm.GetStringBytes("razao")),
					DT_ENTRADA:        dtEntrada,
					NM_SISTEMA:        stringOrDefault(ibm.GetStringBytes("app")),
					SRC_LOAD:          "API/Integração/Thoth",
					DT_LOAD:           string(time.Now().UTC().Format("2006-01-02T15:04:05.000Z")), // Formatar DT_LOAD corretamente
				}

				// novoIbm.DT_ENTRADA = dtEntrada

				notas := ibm.GetArray("notas")
				if len(notas) == 0 {
					fmt.Println("Salvando", novoIbm)
					if err := uc.Repo.SaveCompra(ctx, novoIbm); err != nil {
						return err
					}
				}

				for _, nota := range notas {
					novoIbm.NR_NOTA_FISCAL = stringOrDefault(nota.GetStringBytes("nro"))
					novoIbm.NR_SERIE_NOTA = stringOrDefault(nota.GetStringBytes("serie"))
					novoIbm.DT_EMISSAO_NOTA = int64(parseFloatOrDefault(nota.Get("emissao").GetStringBytes()))
					novoIbm.CNPJ_FORNECEDOR = stringOrDefault(nota.GetStringBytes("fornecedor", "cnpj"))
					novoIbm.NM_FORNECEDOR = stringOrDefault(nota.GetStringBytes("fornecedor", "nome"))
					novoIbm.QT_PESO = float64(int64(parseFloatOrDefault(nota.Get("peso").GetStringBytes())))
					novoIbm.VL_TOTAL_IPI = parseFloat(nota.Get("total", "ipi"))
					novoIbm.VL_TOTAL_ICMS = parseFloat(nota.Get("total", "icms"))
					novoIbm.VL_TOTAL_COMPRA = parseFloat(nota.Get("total", "compra"))
					novoIbm.CNPJ_TRANSPORTADORA = stringOrDefault(nota.GetStringBytes("transportador", "cnpj"))
					novoIbm.NM_TRANSPORTADORA = stringOrDefault(nota.GetStringBytes("transportador", "nome"))
					novoIbm.CD_CHAVE_NOTA_FISCAL = stringOrDefault(nota.GetStringBytes("chavexml"))

					// Processando produtos
					produtos := nota.GetArray("produtos")
					for _, produto := range produtos {
						novoIbm.CD_EAN_PRODUTO = stringOrDefault(produto.GetStringBytes("ean"))
						qtdStr := stringOrDefault(produto.GetStringBytes("qtd"))
						qtdFloat, err := strconv.ParseFloat(qtdStr, 64)
						if err != nil {
							log.Printf("Erro ao converter quantidade: %v", err)
							return err
						}
						novoIbm.QT_PRODUTO = qtdFloat
						novoIbm.VL_PRECO_COMPRA = parseFloat(produto.Get("preco"))
						novoIbm.DS_PRODUTO = stringOrDefault(produto.GetStringBytes("descricao"))
						novoIbm.CD_TP_PRODUTO = stringOrDefault(produto.GetStringBytes("tipo"))
						novoIbm.VL_ALIQUOTA_IPI = parseFloat(produto.Get("impostos", "ipi", "aliquota"))
						novoIbm.VL_IPI = parseFloat(produto.Get("impostos", "ipi", "vlr"))
						novoIbm.VL_ALIQUOTA_ICMS = parseFloat(produto.Get("impostos", "icms", "aliquota"))
						novoIbm.VL_ICMS = parseFloat(produto.Get("impostos", "icms", "vlr"))
						novoIbm.VL_ALIQUOTA_PIS = parseFloat(produto.Get("impostos", "pis", "aliquota"))
						novoIbm.VL_PIS = parseFloat(produto.Get("impostos", "pis", "vlr"))
						novoIbm.VL_ALIQUOTA_COFINS = parseFloat(produto.Get("impostos", "cofins", "aliquota"))
						novoIbm.VL_COFINS = parseFloat(produto.Get("impostos", "cofins", "vlr"))

						// Salvar IBM atualizado
						if err := uc.Repo.SaveCompra(ctx, novoIbm); err != nil {
							log.Printf("Erro ao salvar novo IBM: %v", err)
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
				log.Printf("Erro ao processar IBM: %v", err)
			}
		}(ibm)
		wg.Wait()
	}
	return nil
}
