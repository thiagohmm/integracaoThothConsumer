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
	"go.opentelemetry.io/otel"
)

type EstoqueUseCase struct {
	Repo entities.EstoqueRepository
}

func NewEstoqueUseCase(repo entities.EstoqueRepository) *EstoqueUseCase {
	return &EstoqueUseCase{Repo: repo}
}

func (uc *EstoqueUseCase) ProcessarEstoque(ctx context.Context, estoqueData map[string]interface{}) error {
	tracer := otel.Tracer("ProcessarVenda")
	// Recuperar o UUID do contexto
	uuid, ok := ctx.Value("uuid").(string)
	if !ok {
		return fmt.Errorf("UUID não encontrado no contexto")
	}

	ctx, span := tracer.Start(ctx, uuid)
	defer span.End()

	estoqueJSON, err := json.Marshal(estoqueData)
	if err != nil {
		log.Printf("Erro ao converter mapa em JSON: %v", err)
		span.RecordError(err)
		return err
	}

	// Parsear o JSON usando fastjson
	var p fastjson.Parser
	v, err := p.ParseBytes(estoqueJSON)
	if err != nil {
		log.Printf("Erro ao parsear JSON: %v", err)
		span.RecordError(err)
		return err
	}

	// Limite de goroutines ativas
	const maxGoroutines = 50
	semaphore := make(chan struct{}, maxGoroutines)
	var wg sync.WaitGroup

	// Itera sobre as IBMs da compra
	ibms := v.GetArray("estoque", "ibms")
	if ibms == nil {
		return fmt.Errorf("IBMs não encontrados no objeto de compra")
	}

	for _, estoqueIbms := range ibms {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(estoqueIbms *fastjson.Value) {
			defer wg.Done()
			defer func() { <-semaphore }() // Libera o slot no semáforo

			nroStr := sanitizeIBM(string(estoqueIbms.GetStringBytes("nro")))
			dtaentrada := string(v.GetStringBytes("estoque", "dtaestoque"))
			dtStr := strings.ReplaceAll(dtaentrada, "-", "")

			log.Printf("Deletando IBM compra: %s, data: %s", nroStr, dtStr)
			if err := uc.Repo.DeleteByIBMEstoque(ctx, nroStr, dtStr); err != nil {
				log.Printf("Erro ao deletar IBM compra: %s, erro: %v", nroStr, err)
			}
		}(estoqueIbms)
	}

	// Espera todas as goroutines de deleção terminarem
	wg.Wait()

	// Processa e salva as IBMs em goroutines
	for _, ibm := range ibms {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(ibm *fastjson.Value) {
			defer wg.Done()
			defer func() { <-semaphore }()

			nroStr := sanitizeIBM(string(ibm.GetStringBytes("nro")))
			dtStr := string(v.GetStringBytes("estoque", "dtaestoque"))
			dtEstoque, err := strconv.ParseInt(dtStr, 10, 64)
			if err != nil {
				log.Printf("Erro ao converter DT_Estoque para int64: %v", err)
				span.RecordError(err)
				return
			}

			novoIbm := entities.Estoque{
				CD_IBM_LOJA:       nroStr,
				RAZAO_SOCIAL_LOJA: stringOrDefault(ibm.GetStringBytes("razao")),
				DT_ESTOQUE:        dtEstoque,
				NM_SISTEMA:        stringOrDefault(ibm.GetStringBytes("app")),
				SRC_LOAD:          "API/Integração/Thoth",
				DT_LOAD:           time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
			}

			produtos := ibm.GetArray("produtos")
			if len(produtos) == 0 {
				if err := uc.Repo.SalvarEstoque(ctx, novoIbm); err != nil {
					span.RecordError(err)
					log.Printf("Erro ao salvar IBM Estoque: %v", err)
				}
				return
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

				vlCustoMedio := parseFloat(produto.Get("vlrmedio"))
				if vlCustoMedio > float64(1<<53-1) {
					vlCustoMedioStr := fmt.Sprintf("%.2f", vlCustoMedio)
					vlCustoMedio, _ = strconv.ParseFloat(vlCustoMedioStr[:4], 64)
				}
				novoIbm.VL_CUSTO_MEDIO = vlCustoMedio

				if err := uc.Repo.SalvarEstoque(ctx, novoIbm); err != nil {
					span.RecordError(err)
					log.Printf("Erro ao salvar novo IBM Estoque: %v", err)
				}
			}
		}(ibm)
	}

	// Espera todas as goroutines de salvamento terminarem
	wg.Wait()

	log.Printf("Processamento de estoque concluído")
	return nil
}
