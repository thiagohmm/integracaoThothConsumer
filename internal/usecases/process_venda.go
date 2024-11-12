package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/thiagohmm/integracaoThothConsumer/internal/domain/entities"
	"github.com/valyala/fastjson"
	"go.opentelemetry.io/otel"
)

type VendaUseCase struct {
	Repo entities.VendaRepository
}

func NewVendaUseCase(repo entities.VendaRepository) *VendaUseCase {
	return &VendaUseCase{Repo: repo}
}

func (uc *VendaUseCase) ProcessarVenda(ctx context.Context, vendaData map[string]interface{}) error {
	tracer := otel.Tracer("ProcessarVenda")
	// Recuperar o UUID do contexto
	uuid, ok := ctx.Value("uuid").(string)
	if !ok {
		return fmt.Errorf("UUID não encontrado no contexto")
	}

	ctx, span := tracer.Start(ctx, uuid)
	defer span.End()

	vendaJSON, err := json.Marshal(vendaData)
	if err != nil {
		span.RecordError(err)
		log.Printf("Erro ao converter mapa em JSON: %v", err)
		return err
	}

	// Parsear o JSON usando fastjson
	var p fastjson.Parser
	v, err := p.ParseBytes(vendaJSON)
	if err != nil {
		log.Printf("Erro ao parsear JSON: %v", err)
		span.RecordError(err)
		return err
	}

	ibms := v.GetArray("vendas", "ibms")
	if ibms == nil {
		span.RecordError(err)
		return fmt.Errorf("IBMs não encontrados no objeto de venda")

	}

	var wg sync.WaitGroup
	goroutineLimit := make(chan struct{}, 50) // Limite de 10 goroutines simultâneas
	errorChan := make(chan error, len(ibms))  // Canal para capturar erros de goroutines

	// Processa cada IBM em goroutines
	for _, vendaIbms := range ibms {
		wg.Add(1)
		go func(vendaIbms *fastjson.Value) {
			defer wg.Done()

			goroutineLimit <- struct{}{} // Bloqueia se o limite de goroutines for atingido
			defer func() { <-goroutineLimit }()

			nroStr := sanitizeIBM(string(vendaIbms.GetStringBytes("nro")))
			dtStr := string(v.GetStringBytes("vendas", "dtavenda"))
			dtEntrada, err := strconv.ParseInt(dtStr, 10, 64)
			if err != nil {
				log.Printf("Erro ao converter DT_Venda para int64: %v", err)
				span.RecordError(err)
				errorChan <- err
				return
			}

			if err := uc.Repo.DeleteByIBMVenda(ctx, nroStr, dtStr); err != nil {
				log.Printf("Erro ao deletar IBM venda: %s, erro: %v", nroStr, err)
				span.RecordError(err)
				errorChan <- err
				return
			}

			novoIbm := entities.Venda{
				DT_TRANSACAO:      dtEntrada,
				CD_IBM_LOJA:       nroStr,
				RAZAO_SOCIAL_LOJA: stringOrDefault(vendaIbms.GetStringBytes("razao")),
				NM_SISTEMA:        stringOrDefault(vendaIbms.GetStringBytes("app")),
				DT_ARQUIVO:        stringOrDefault(v.GetStringBytes("vendas", "dtaenvio")),
				SRC_LOAD:          "API/Integração/Thoth",
				DT_LOAD:           string(time.Now().UTC().Format("2006-01-02T15:04:05.000Z")),
			}

			vendas := vendaIbms.GetArray("vendas")
			if len(vendas) == 0 {
				if err := uc.Repo.SalvarVenda(ctx, novoIbm); err != nil {
					log.Printf("Erro ao salvar IBM Venda: %v", err)
					span.RecordError(err)
					errorChan <- err
					return
				}
			}

			for _, venda := range vendas {
				novoIbm.HR_INICIO_TRANSACAO = stringOrDefault(venda.GetStringBytes("ini"))
				novoIbm.HR_FIM_TRANSACAO = stringOrDefault(venda.GetStringBytes("fim"))
				novoIbm.CD_TRANSACAO = stringOrDefault(venda.GetStringBytes("doc"))
				novoIbm.CPF_CNPJ_CLIENTE = stringOrDefault(venda.GetStringBytes("cpfcnpj"))
				novoIbm.NM_FORMA_PAGAMENTO = stringOrDefault(venda.GetStringBytes("formapagto"))
				novoIbm.NM_BANDEIRA = stringOrDefault(venda.GetStringBytes("bandeira"))
				novoIbm.CD_CCF = stringOrDefault(venda.GetStringBytes("ccf"))
				novoIbm.CD_MODELO_DOCTO = stringOrDefault(venda.GetStringBytes("moddoc"))

				produtos := venda.GetArray("produtos")
				for _, produto := range produtos {
					novoIbm.CD_EAN_PRODUTO = stringOrDefault(produto.GetStringBytes("ean"))
					novoIbm.QT_PRODUTO = stringOrDefault(produto.GetStringBytes("qtd"))
					novoIbm.VL_PRECO_UNITARIO = stringOrDefault(produto.GetStringBytes("preco"))
					novoIbm.VL_IMPOSTO = stringOrDefault(produto.GetStringBytes("imposto"))
					novoIbm.VL_FATURADO = stringOrDefault(produto.GetStringBytes("total"))
					novoIbm.VL_CUSTO_UNITARIO = stringOrDefault(produto.GetStringBytes("custo"))
					novoIbm.CD_DEPARTAMENTO = stringOrDefault(produto.GetStringBytes("dep"))
					tipo := produto.GetInt("tipo")
					itemTrans := produto.GetInt("trans")
					novoIbm.CD_TP_PRODUTO = strconv.FormatInt(int64(tipo), 10)
					novoIbm.DS_PRODUTO = stringOrDefault(produto.GetStringBytes("descricao"))
					novoIbm.CD_PROMOCAO = stringOrDefault(produto.GetStringBytes("codmix"))
					novoIbm.CD_EAN_EMBALAGEM = stringOrDefault(produto.GetStringBytes("eanpack"))
					novoIbm.CD_TP_TRANSACAO = strconv.FormatInt(int64(itemTrans), 10)
					novoIbm.VL_DESCONTO = stringOrDefault(produto.GetStringBytes("desconto"))
					novoIbm.CD_ITEM_TRANSACAO = strconv.FormatInt(int64(itemTrans), 10)

					// Salva o IBM atualizado
					if err := uc.Repo.SalvarVenda(ctx, novoIbm); err != nil {
						log.Printf("Erro ao salvar novo IBM: %v", err)
						span.RecordError(err)
						errorChan <- err
						return
					}
				}
			}
		}(vendaIbms)
	}

	wg.Wait()
	close(errorChan)

	// Verifica se houve algum erro nas goroutines
	for err := range errorChan {
		if err != nil {
			log.Printf("Erro durante o processamento: %v", err)
			span.RecordError(err)
			return err
		}
	}
	return nil
}
