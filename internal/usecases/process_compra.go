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
	// Ajustado para importar a interface correta
)

type CompraUseCase struct {
	Repo entities.CompraRepository // Troquei CompraRepositoryDB por CompraRepository (interface)
}

func NewCompraUseCase(repo entities.CompraRepository) *CompraUseCase {
	return &CompraUseCase{Repo: repo}
}

// Função auxiliar para converter data
func parseDate(data string) int {
	dateInt, _ := strconv.Atoi(strings.ReplaceAll(data, "-", ""))
	return dateInt
}

func parseFloat(value interface{}) float64 {
	if value == nil {
		return 0.0
	}
	floatValue, _ := strconv.ParseFloat(fmt.Sprintf("%v", value), 64)
	return floatValue
}

func (uc *CompraUseCase) ProcessarCompra(ctx context.Context, compraData map[string]interface{}) error {
	// Mapeia os dados de compra

	var persisteCompra map[string]interface{}
	var arquivoCompra []entities.Compra
	var novoIbm entities.Compra
	var novoIbmCompra map[string]interface{}

	// Verifica se "compra" é uma string JSON ou um objeto
	if compra, ok := compraData["compra"].(string); ok {
		err := json.Unmarshal([]byte(compra), &persisteCompra)
		if err != nil {
			log.Printf("Erro ao fazer parse da string JSON: %v", err)
			return err
		}
	} else if compra, ok := compraData["compra"].(map[string]interface{}); ok {
		persisteCompra = compra
	} else {
		log.Printf("A variável 'compra' não é uma string JSON válida nem um objeto.")
		return fmt.Errorf("formato de compra inválido")
	}

	deleteCounter := 0
	// Itera sobre as IBMs da compra
	ibms, ok := persisteCompra["compras"].(map[string]interface{})["ibms"].([]interface{})
	if !ok {
		return fmt.Errorf("IBMs não encontrados no objeto de compra")
	}

	for _, compraIbmsInterface := range ibms {
		compraIbms := compraIbmsInterface.(map[string]interface{})

		// Converte "nro" para string e aplica formatação
		nroStr := strings.TrimSpace(fmt.Sprintf("%v", compraIbms["nro"]))
		nroStr = strings.TrimLeft(nroStr, "0")
		if len(nroStr) > 10 {
			nroStr = nroStr[len(nroStr)-10:]
		} else {
			nroStr = fmt.Sprintf("%010s", nroStr)
		}

		// Obtém a data de entrada
		dtaentrada := persisteCompra["compras"].(map[string]interface{})["dtaentrada"]
		dtStr := strings.TrimSpace(fmt.Sprintf("%v", dtaentrada))
		dtStr = strings.ReplaceAll(dtStr, "-", "")

		// Deleta o IBM com base no número e data
		log.Printf("Deletando IBM compra: %s, data: %s", nroStr, dtStr)
		if err := uc.Repo.DeleteByIBMAndEntrada(ctx, nroStr, dtStr); err != nil {
			log.Printf("Erro ao deletar IBM compra: %v, erro: %v", compraIbms["nro"], err)
			continue
		}

		deleteCounter++
		if deleteCounter%100 == 0 {
			time.Sleep(500 * time.Millisecond) // Pausa para evitar sobrecarga
		}
	}

	// Verifica o tipo de compraData e faz o parse se for uma string
	compra, ok := compraData["compra"]
	if !ok {
		return fmt.Errorf("chave 'compra' não encontrada")
	}

	switch v := compra.(type) {
	case string:
		err := json.Unmarshal([]byte(v), &novoIbmCompra)
		if err != nil {
			log.Printf("Erro ao fazer parse da string JSON: %v", err)
			return err
		}
	case map[string]interface{}:
		novoIbmCompra = v
	default:
		return fmt.Errorf("A variável 'compra' não é uma string JSON válida nem um objeto")
	}

	saveCounter := 0

	// Itera sobre as IBMs da compra
	ibms, ok = novoIbmCompra["compras"].(map[string]interface{})["ibms"].([]interface{})
	if !ok {
		return fmt.Errorf("IBMs não encontrados no objeto de compra")
	}

	for _, ibmInterface := range ibms {
		ibm := ibmInterface.(map[string]interface{})

		// Tratamento de erro com defer para rollback em caso de falha
		err := func() error {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Erro ao processar IBM Compra: %v", r)
				}
			}()

			// Converter nro para string
			nroStr := fmt.Sprintf("%v", ibm["nro"])
			nroStr = strings.TrimLeft(strings.TrimSpace(nroStr), "0")
			if len(nroStr) > 10 {
				nroStr = nroStr[len(nroStr)-10:]
			} else {
				nroStr = fmt.Sprintf("%010s", nroStr)
			}

			// Preenchendo campos da nova IBM
			novoIbm.CD_IBM_LOJA = nroStr
			novoIbm.RAZAO_SOCIAL_LOJA = fmt.Sprintf("%v", ibm["razao"])
			dtaentradaStr, ok := novoIbmCompra["compras"].(map[string]interface{})["dtaentrada"].(string)
			if !ok {
				return fmt.Errorf("dtaentrada is not a string")
			}
			novoIbm.DT_ENTRADA = int64(parseDate(dtaentradaStr))
			novoIbm.NM_SISTEMA = fmt.Sprintf("%v", ibm["app"])
			novoIbm.SRC_LOAD = "API/Integração/Thoth"
			novoIbm.DT_LOAD = time.Now().Format(time.RFC3339)

			// Processando notas
			notas, ok := ibm["notas"].([]interface{})
			if !ok || len(notas) == 0 {
				arquivoCompra = append(arquivoCompra, novoIbm)
				if err := uc.Repo.Save(ctx, novoIbm); err != nil {
					return err
				}
				saveCounter++
				if saveCounter%100 == 0 {
					time.Sleep(500 * time.Millisecond)
				}
				return nil
			}

			for _, notaInterface := range notas {
				nota := notaInterface.(map[string]interface{})

				// Preenchendo dados da nota
				novoIbm.NR_NOTA_FISCAL = fmt.Sprintf("%v", nota["nro"])
				novoIbm.NR_SERIE_NOTA = fmt.Sprintf("%v", nota["serie"])
				emissaoStr, ok := nota["emissao"].(string)
				if !ok {
					return fmt.Errorf("emissao is not a string")
				}
				novoIbm.DT_EMISSAO_NOTA = int64(parseDate(emissaoStr))
				novoIbm.CNPJ_FORNECEDOR = fmt.Sprintf("%v", nota["fornecedor"].(map[string]interface{})["cnpj"])
				novoIbm.NM_FORNECEDOR = fmt.Sprintf("%v", nota["fornecedor"].(map[string]interface{})["nome"])
				novoIbm.QT_PESO = parseFloat(nota["total"].(map[string]interface{})["peso"])
				novoIbm.VL_TOTAL_IPI = parseFloat(nota["total"].(map[string]interface{})["vlripi"])
				novoIbm.VL_TOTAL_ICMS = parseFloat(nota["total"].(map[string]interface{})["vlricms"])
				novoIbm.VL_TOTAL_COMPRA = parseFloat(nota["total"].(map[string]interface{})["vlrnota"])
				novoIbm.CNPJ_TRANSPORTADORA = fmt.Sprintf("%v", nota["transportador"].(map[string]interface{})["cnpj"])
				novoIbm.NM_TRANSPORTADORA = fmt.Sprintf("%v", nota["transportador"].(map[string]interface{})["nome"])
				novoIbm.CD_CHAVE_NOTA_FISCAL = fmt.Sprintf("%v", nota["chavexml"])

				// Processando produtos
				produtos, ok := nota["produtos"].([]interface{})
				if !ok || len(produtos) == 0 {
					arquivoCompra = append(arquivoCompra, novoIbm)
					if err := uc.Repo.Save(ctx, novoIbm); err != nil {
						return err
					}
					saveCounter++
					if saveCounter%100 == 0 {
						time.Sleep(500 * time.Millisecond)
					}
					continue
				}

				for _, produtoInterface := range produtos {
					produto := produtoInterface.(map[string]interface{})

					// Preenchendo dados do produto
					novoIbm.CD_EAN_PRODUTO = fmt.Sprintf("%v", produto["ean"])
					novoIbm.QT_PRODUTO = parseFloat(produto["qtd"])
					novoIbm.VL_PRECO_COMPRA = parseFloat(produto["preco"])
					novoIbm.DS_PRODUTO = fmt.Sprintf("%v", produto["descricao"])
					novoIbm.CD_TP_PRODUTO = fmt.Sprintf("%v", produto["tipo"])
					novoIbm.VL_ALIQUOTA_IPI = parseFloat(produto["impostos"].(map[string]interface{})["ipi"].(map[string]interface{})["aliquota"])
					novoIbm.VL_IPI = parseFloat(produto["impostos"].(map[string]interface{})["ipi"].(map[string]interface{})["vlr"])
					novoIbm.VL_ALIQUOTA_ICMS = parseFloat(produto["impostos"].(map[string]interface{})["icms"].(map[string]interface{})["aliquota"])
					novoIbm.VL_ICMS = parseFloat(produto["impostos"].(map[string]interface{})["icms"].(map[string]interface{})["vlr"])
					novoIbm.VL_ALIQUOTA_PIS = parseFloat(produto["impostos"].(map[string]interface{})["pis"].(map[string]interface{})["aliquota"])
					novoIbm.VL_PIS = parseFloat(produto["impostos"].(map[string]interface{})["pis"].(map[string]interface{})["vlr"])
					novoIbm.VL_ALIQUOTA_COFINS = parseFloat(produto["impostos"].(map[string]interface{})["cofins"].(map[string]interface{})["aliquota"])
					novoIbm.VL_COFINS = parseFloat(produto["impostos"].(map[string]interface{})["cofins"].(map[string]interface{})["vlr"])
					novoIbm.CD_NCM = fmt.Sprintf("%v", produto["ncm"])
					novoIbm.CD_ITEM_NOTA_FISCAL = fmt.Sprintf("%v", produto["linha"])
					novoIbm.CD_PRODUTO_FORNECEDOR = fmt.Sprintf("%v", produto["codfornec"])
					novoIbm.QT_PRODUTO_CONVERTIDA = parseFloat(produto["qtdenf"])
					novoIbm.DS_UN_MEDIDA_CONVERTIDA = fmt.Sprintf("%v", produto["unconv"])
					novoIbm.DS_UN_MEDIDA = fmt.Sprintf("%v", produto["un"])
					novoIbm.VL_ULTIMO_CUSTO = parseFloat(produto["ultcusto"])

					arquivoCompra = append(arquivoCompra, novoIbm)
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
			log.Printf("Erro ao processar IBM Compra: %v", err)

		}

	}

	return nil

}
