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

type VendaUseCase struct {
	Repo entities.VendaRepository
}

func NewVendaUseCase(repo entities.VendaRepository) *VendaUseCase {
	return &VendaUseCase{Repo: repo}
}

func (uc *VendaUseCase) ProcessarVenda(ctx context.Context, vendaData map[string]interface{}) error {
	vendaJSON, err := json.Marshal(vendaData)
	if err != nil {
		log.Printf("Erro ao converter mapa em JSON: %v", err)
		return err
	}

	// Parsear o JSON usando fastjson
	var p fastjson.Parser
	v, err := p.ParseBytes(vendaJSON)
	if err != nil {
		log.Printf("Erro ao parsear JSON: %v", err)
		return err
	}

	deleteCounter := 0
	// Itera sobre as IBMs da compra
	ibms := v.GetArray("vendas", "ibms")
	if ibms == nil {
		return fmt.Errorf("IBMs não encontrados no objeto de venda")
	}

	for _, vendaIbms := range ibms {
		nroStr := sanitizeIBM(string(vendaIbms.GetStringBytes("nro")))

		// Obtém a data de entrada
		dtaentrada := string(v.GetStringBytes("vendas", "dtavenda"))
		dtStr := strings.ReplaceAll(dtaentrada, "-", "")

		// Deleta o IBM com base no número e data
		log.Printf("Deletando IBM Venda: %s, data: %s", nroStr, dtStr)
		if err := uc.Repo.DeleteByIBMVenda(ctx, nroStr, dtStr); err != nil {
			log.Printf("Erro ao deletar IBM venda: %s, erro: %v", nroStr, err)
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
			dtStr := string(v.GetStringBytes("vendas", "dtavenda"))
			dtEntrada, err := strconv.ParseInt(dtStr, 10, 64)
			if err != nil {
				log.Printf("Erro ao converter DT_Venda para int64: %v", err)
				return err
			}
			novoIbm := entities.Venda{
				DT_TRANSACAO:      dtEntrada,
				CD_IBM_LOJA:       nroStr,
				RAZAO_SOCIAL_LOJA: stringOrDefault(ibm.GetStringBytes("razao")),
				NM_SISTEMA:        stringOrDefault(ibm.GetStringBytes("app")),
				DT_ARQUIVO:        stringOrDefault(ibm.GetStringBytes("dtaenvio")),
				SRC_LOAD:          "API/Integração/Thoth",
				DT_LOAD:           string(time.Now().UTC().Format("2006-01-02T15:04:05.000Z")), // Formatar DT_LOAD corretamente
			}

			// novoIbm.DT_ENTRADA = dtEntrada

			vendas := ibm.GetArray("vendas")
			if len(vendas) == 0 {
				fmt.Println("Salvando", novoIbm)
				if err := uc.Repo.SalvarVenda(ctx, novoIbm); err != nil {
					return err
				}
			}

			for _, venda := range vendas {
				novoIbm.HR_INICIO_TRANSACAO = stringOrDefault(venda.GetStringBytes("inicio"))
				novoIbm.HR_FIM_TRANSACAO = stringOrDefault(venda.GetStringBytes("fim"))
				novoIbm.CD_TRANSACAO = stringOrDefault(venda.GetStringBytes("transacao"))
				novoIbm.CPF_CNPJ_CLIENTE = stringOrDefault(venda.GetStringBytes("cliente", "cnpj"))
				novoIbm.NM_FORMA_PAGAMENTO = stringOrDefault(venda.GetStringBytes("forma"))
				novoIbm.CD_CCF = stringOrDefault(venda.GetStringBytes("ccf"))
				novoIbm.CD_MODELO_DOCTO = stringOrDefault(venda.GetStringBytes("modelo"))

				// Processando produtos
				produtos := venda.GetArray("produtos")
				for _, produto := range produtos {
					novoIbm.CD_EAN_PRODUTO = stringOrDefault(produto.GetStringBytes("ean"))
					qtdStr := stringOrDefault(produto.GetStringBytes("qtd"))
					qtdFloat, err := strconv.ParseFloat(qtdStr, 64)
					if err != nil {
						log.Printf("Erro ao converter quantidade: %v", err)
						return err
					}
					novoIbm.QT_PRODUTO = strconv.FormatFloat(qtdFloat, 'f', -1, 64)
					novoIbm.VL_PRECO_UNITARIO = strconv.FormatFloat(parseFloat(produto.Get("preco")), 'f', -1, 64)
					novoIbm.VL_IMPOSTO = strconv.FormatFloat(parseFloat(produto.Get("imposto")), 'f', -1, 64)
					novoIbm.VL_FATURADO = strconv.FormatFloat(parseFloat(produto.Get("total")), 'f', -1, 64)
					novoIbm.VL_CUSTO_UNITARIO = strconv.FormatFloat(parseFloat(produto.Get("custo")), 'f', -1, 64)
					novoIbm.CD_DEPARTAMENTO = stringOrDefault(produto.GetStringBytes("dep"))
					novoIbm.CD_TP_PRODUTO = stringOrDefault(produto.GetStringBytes("tipo"))
					novoIbm.DS_PRODUTO = stringOrDefault(produto.GetStringBytes("descricao"))
					novoIbm.CD_PROMOCAO = stringOrDefault(produto.GetStringBytes("codmix"))
					novoIbm.CD_EAN_EMBALAGEM = stringOrDefault(produto.GetStringBytes("eanpack"))
					novoIbm.CD_TP_TRANSACAO = stringOrDefault(produto.GetStringBytes("trans"))
					novoIbm.VL_DESCONTO = stringOrDefault(produto.GetStringBytes("desconto"))
					novoIbm.CD_ITEM_TRANSACAO = stringOrDefault(produto.GetStringBytes("trans"))

					// Salvar IBM atualizado
					if err := uc.Repo.SalvarVenda(ctx, novoIbm); err != nil {
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
	}
	return nil
}
