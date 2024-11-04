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

	var p fastjson.Parser
	v, err := p.ParseBytes(estoqueJSON)
	if err != nil {
		log.Printf("Erro ao parsear JSON: %v", err)
		return err
	}

	ibms := v.GetArray("estoque", "ibms")
	if ibms == nil {
		return fmt.Errorf("IBMs não encontrados no objeto de estoque")
	}

	dtaestoque := string(v.GetStringBytes("estoque", "dtaestoque"))
	dtStr := strings.ReplaceAll(dtaestoque, "-", "")
	deleteCounter, saveCounter := 0, 0

	for _, estoqueIbms := range ibms {
		nroStr := sanitizeIBM(string(estoqueIbms.GetStringBytes("nro")))

		if err := uc.Repo.DeleteByIBMEstoque(ctx, nroStr, dtStr); err != nil {
			log.Printf("Erro ao deletar IBM estoque: %s, erro: %v", nroStr, err)
			continue
		}

		deleteCounter++
		if deleteCounter%100 == 0 {
			time.Sleep(500 * time.Millisecond)
		}
	}

	for _, ibm := range ibms {
		err := uc.processarESalvarIbm(ctx, v, ibm, dtStr, &saveCounter)
		if err != nil {
			log.Printf("Erro ao processar IBM: %v", err)
		}
	}

	return nil
}

func (uc *EstoqueUseCase) processarESalvarIbm(
	ctx context.Context,
	v *fastjson.Value,
	ibm *fastjson.Value,
	dtStr string,
	saveCounter *int,
) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Erro ao processar IBM Estoque: %v", r)
		}
	}()

	nroStr := sanitizeIBM(string(ibm.GetStringBytes("nro")))
	dtEntrada, err := strconv.ParseInt(dtStr, 10, 64)
	if err != nil {
		log.Printf("Erro ao converter DT_ESTOQUE para int64: %v", err)
		return err
	}

	novoIbm := entities.Estoque{
		CD_IBM_LOJA:       nroStr,
		RAZAO_SOCIAL_LOJA: stringOrDefault(ibm.GetStringBytes("razao")),
		DT_ESTOQUE:        dtEntrada,
		NM_SISTEMA:        stringOrDefault(ibm.GetStringBytes("app")),
		SRC_LOAD:          "API/Integração/Thoth",
		DT_LOAD:           time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
	}

	produtos := ibm.GetArray("produtos")
	if len(produtos) == 0 {
		return uc.salvarEstoque(ctx, novoIbm, saveCounter)
	}

	for _, produto := range produtos {
		err := uc.processarProduto(&novoIbm, produto)
		if err != nil {
			log.Printf("Erro ao processar produto: %v", err)
			continue
		}

		if err := uc.salvarEstoque(ctx, novoIbm, saveCounter); err != nil {
			log.Printf("Erro ao salvar novo IBM Estoque: %v", err)
		}
	}

	return nil
}

func (uc *EstoqueUseCase) salvarEstoque(ctx context.Context, novoIbm entities.Estoque, saveCounter *int) error {
	if err := uc.Repo.SaveEstoque(ctx, novoIbm); err != nil {
		return err
	}

	*saveCounter++
	if *saveCounter%100 == 0 {
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

func (uc *EstoqueUseCase) processarProduto(novoIbm *entities.Estoque, produto *fastjson.Value) error {
	novoIbm.CD_EAN_PRODUTO = stringOrDefault(produto.GetStringBytes("ean"))
	novoIbm.CD_TP_PRODUTO = stringOrDefault(produto.GetStringBytes("tipo"))
	novoIbm.DT_ULTIMA_COMPRA = int64(parseFloatOrDefault(produto.Get("emissao").GetStringBytes()))
	novoIbm.DS_PRODUTO = stringOrDefault(produto.GetStringBytes("fornecedor", "cnpj"))
	novoIbm.VL_PRECO_UNITARIO = parseFloat(produto.Get("preco"))
	novoIbm.QT_INVENTARIO_ENTRADA = float64(int64(parseFloatOrDefault(produto.Get("qtdentrada").GetStringBytes())))
	novoIbm.QT_INVENTARIO_SAIDA = float64(int64(parseFloatOrDefault(produto.Get("qtdsaida").GetStringBytes())))
	novoIbm.QT_INICIAL_PRODUTO = float64(int64(parseFloatOrDefault(produto.Get("qtdini").GetStringBytes())))
	novoIbm.QT_FINAL_PRODUTO = float64(int64(parseFloatOrDefault(produto.Get("qtdfim").GetStringBytes())))
	novoIbm.VL_TOTAL_ESTOQUE = parseFloat(produto.Get("vlrfim"))

	vlCustoMedio := parseFloat(produto.Get("vlrmedio"))
	if vlCustoMedio > float64(1<<53-1) {
		vlCustoMedioStr := fmt.Sprintf("%.2f", vlCustoMedio)
		vlCustoMedio, _ = strconv.ParseFloat(vlCustoMedioStr[:4], 64)
	}
	novoIbm.VL_CUSTO_MEDIO = vlCustoMedio

	return nil
}
