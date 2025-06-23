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
	"go.opentelemetry.io/otel"
)

type EstoqueUseCase struct {
	Repo entities.EstoqueRepository
}

func NewEstoqueUseCase(repo entities.EstoqueRepository) *EstoqueUseCase {
	return &EstoqueUseCase{Repo: repo}
}

// ProcessarEstoque realiza o processamento de estoque de forma síncrona,
// primeiro deletando os registros existentes e, em seguida, salvando as novas informações.
func (uc *EstoqueUseCase) ProcessarEstoque(ctx context.Context, estoqueData map[string]interface{}) error {
	tracer := otel.Tracer("ProcessarEstoque")

	uuid, ok := ctx.Value("uuid").(string)
	if !ok {
		return fmt.Errorf("UUID não encontrado no contexto")
	}

	ctxDeleteTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	ctx, span := tracer.Start(ctx, uuid)
	defer span.End()

	var dataSave []entities.Estoque

	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		log.Printf("Falha ao carregar localização: %v", err)
		loc = time.Local
		span.RecordError(err)
	}
	dtNow := time.Now().In(loc).Format("2006-01-02T15:04:05.000Z")

	estoqueJSON, err := json.Marshal(estoqueData)
	if err != nil {
		log.Printf("Erro ao converter mapa em JSON: %v", err)
		span.RecordError(err)
		return err
	}

	var p fastjson.Parser
	v, err := p.ParseBytes(estoqueJSON)
	if err != nil {
		log.Printf("Erro ao parsear JSON: %v", err)
		span.RecordError(err)
		return err
	}

	// Verifica e formata a data de estoque
	dtaestoque := string(v.GetStringBytes("estoque", "dtaestoque"))
	if dtaestoque == "" {
		msg := "campo dtaestoque está vazio"
		log.Print(msg)
		span.RecordError(fmt.Errorf(msg))
		return fmt.Errorf(msg)
	}
	dtStr := strings.ReplaceAll(dtaestoque, "-", "")
	if len(dtStr) != 8 {
		msg := fmt.Sprintf("formato inválido em dtaestoque: %s", dtaestoque)
		log.Print(msg)
		span.RecordError(fmt.Errorf(msg))
		return fmt.Errorf(msg)
	}
	dtEstoqueInt, err := strconv.ParseInt(dtStr, 10, 64)
	if err != nil {
		log.Printf("Erro ao converter dtaestoque para int64: %v", err)
		span.RecordError(err)
		return err
	}

	ibms := v.GetArray("estoque", "ibms")
	if ibms == nil {
		err := fmt.Errorf("IBMs não encontrados no objeto de estoque")
		log.Print(err)
		span.RecordError(err)
		return err
	}

	// Deleção dos registros antigos
	for _, estoqueIbms := range ibms {
		nroStr := sanitizeIBM(string(estoqueIbms.GetStringBytes("nro")))

		log.Printf("Checkando se  estoque existe: %s, data: %s uuid %s", nroStr, dtStr, uuid)
		existEstoque, err := uc.Repo.CheckIfExists(ctxDeleteTimeout, nroStr, dtStr)
		if err != nil {
			log.Printf("Erro ao verificar existência de estoque: %s, erro: %v", nroStr, err)
			span.RecordError(err)
		}

		if existEstoque {
			log.Printf("Deletando IBM estoque: %s, data: %s uuid %s", nroStr, dtStr, uuid)

			if err := uc.Repo.DeleteByIBMEstoque(ctxDeleteTimeout, nroStr, dtStr, uuid); err != nil {
				log.Printf("Erro ao deletar IBM estoque: %s, erro: %v", nroStr, err)
				span.RecordError(err)
				return err
			}
		}
	}

	// Salvamento de novos dados
	for _, ibm := range ibms {
		nroStr := sanitizeIBM(string(ibm.GetStringBytes("nro")))
		novoIbm := entities.Estoque{
			CD_IBM_LOJA:       nroStr,
			RAZAO_SOCIAL_LOJA: stringOrDefault(ibm.GetStringBytes("razao")),
			DT_ESTOQUE:        dtEstoqueInt,
			NM_SISTEMA:        stringOrDefault(ibm.GetStringBytes("app")),
			SRC_LOAD:          "API/Integração/Thoth",
			DT_LOAD:           dtNow,
		}

		produtos := ibm.GetArray("produtos")
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

			dataSave = append(dataSave, novoIbm)
		}
	}

	if err := uc.Repo.SalvarEstoques(ctx, dataSave, uuid); err != nil {
		log.Printf("Erro ao salvar novos registros de estoque: %v", err)
		span.RecordError(err)
		return err
	}

	log.Printf("Processamento de estoque concluído com sucesso (%d registros)", len(dataSave))
	return nil
}
