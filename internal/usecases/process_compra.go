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

// ProcessarCompra realiza o processamento de compra de forma síncrona,
// sem o uso de goroutines para deleção e salvamento.
func (uc *CompraUseCase) ProcessarCompra(ctx context.Context, compraData map[string]interface{}) (bool, error) {
	tracer := otel.Tracer("ProcessarCompra")
	// Recuperar o UUID do contexto
	uuid, ok := ctx.Value("uuid").(string)
	if !ok {
		return false, fmt.Errorf("UUID não encontrado no contexto")
	}
	ctxDeleteTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	ctx, span := tracer.Start(ctx, uuid)
	defer span.End()

	var dataSave []entities.Compra

	compraJSON, err := json.Marshal(compraData)
	if err != nil {
		log.Printf("Erro ao converter mapa em JSON: %v", err)
		span.RecordError(err)
		return false, err
	}

	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		log.Printf("Falha ao carregar localização: %v", err)
		loc = time.Local
		span.RecordError(err)
	}
	dtNow := string(time.Now().In(loc).Format("2006-01-02T15:04:05.000Z"))

	// Parsear o JSON usando fastjson
	var p fastjson.Parser
	v, err := p.ParseBytes(compraJSON)
	if err != nil {
		log.Printf("Erro ao parsear JSON: %v", err)
		span.RecordError(err)
		return false, err
	}

	// Recupera o array de IBMs da compra
	ibms := v.GetArray("compras", "ibms")
	if ibms == nil {
		span.RecordError(fmt.Errorf("IBMs não encontrados no objeto de compra"))
		return false, err
	}

	for _, compraIbms := range ibms {
		nroStr := sanitizeIBM(string(compraIbms.GetStringBytes("nro")))
		// Obtém a data de entrada e remove caracteres desnecessários
		dtaentrada := string(v.GetStringBytes("compras", "dtaentrada"))
		dtStr := strings.ReplaceAll(dtaentrada, "-", "")

		log.Printf("Verificando se existe compra: %s, data: %s", nroStr, dtStr)
		existsCompra, err := uc.Repo.CheckIfExists(ctxDeleteTimeout, nroStr, dtStr)
		if err != nil {
			span.RecordError(err)
			log.Printf("Erro ao verificar existência de compra: %s, erro: %v", nroStr, err)
			return false, err
		}

		if existsCompra {

			log.Printf("Deletando IBM compra: %s, data: %s para UUID: %s", nroStr, dtStr, uuid)
			if err := uc.Repo.DeleteByIBMAndEntrada(ctxDeleteTimeout, nroStr, dtStr, uuid); err != nil {
				span.RecordError(err)
				log.Printf("Erro ao deletar IBM compra: %s, erro: %v", nroStr, err)
				return false, err
			}
		}
	}

	// Processar e salvar os IBMs de compra (de forma sequencial)
	log.Printf("Iniciando salvamento de IBMs da compra")
	for _, ibm := range ibms {
		nroStr := sanitizeIBM(string(ibm.GetStringBytes("nro")))
		dtStr := string(v.GetStringBytes("compras", "dtaentrada"))
		dtEntrada, err := strconv.ParseInt(dtStr, 10, 64)
		if err != nil {
			span.RecordError(err)
			log.Printf("Erro ao converter DT_ENTRADA para int64: %v", err)
			return false, err
		}
		novoIbm := entities.Compra{
			CD_IBM_LOJA:       nroStr,
			RAZAO_SOCIAL_LOJA: stringOrDefault(ibm.GetStringBytes("razao")),
			DT_ENTRADA:        dtEntrada,
			NM_SISTEMA:        stringOrDefault(ibm.GetStringBytes("app")),
			SRC_LOAD:          "API/Integração/Thoth",
			DT_LOAD:           dtNow,
		}

		notas := ibm.GetArray("notas")

		// Se houver notas, processar cada uma
		for _, nota := range notas {
			novoIbm.NR_NOTA_FISCAL = stringOrDefault(nota.GetStringBytes("nro"))
			novoIbm.NR_SERIE_NOTA = stringOrDefault(nota.GetStringBytes("serie"))
			emissaoStr := stringOrDefault(nota.GetStringBytes("emissao"))

			if emissaoStr == "" {
				emissaoStr = "0" // optionally use a default value or skip processing this nota
			}
			emissaoInt, err := strconv.ParseInt(emissaoStr, 10, 64)
			if err != nil {
				span.RecordError(err)
				log.Printf("Erro ao converter DT_EMISSAO_NOTA para int64: %v", err)
				continue
			}
			novoIbm.DT_EMISSAO_NOTA = emissaoInt
			novoIbm.CNPJ_FORNECEDOR = stringOrDefault(nota.GetStringBytes("fornecedor", "cnpj"))
			novoIbm.NM_FORNECEDOR = stringOrDefault(nota.GetStringBytes("fornecedor", "nome"))
			novoIbm.QT_PESO = float64(int64(parseFloatOrDefault(nota.Get("peso").GetStringBytes())))
			novoIbm.VL_TOTAL_IPI = parseFloat(nota.Get("total", "ipi"))
			novoIbm.VL_TOTAL_ICMS = parseFloat(nota.Get("total", "icms"))
			novoIbm.VL_TOTAL_COMPRA = parseFloat(nota.Get("total", "compra"))
			novoIbm.CNPJ_TRANSPORTADORA = stringOrDefault(nota.GetStringBytes("transportador", "cnpj"))
			novoIbm.NM_TRANSPORTADORA = stringOrDefault(nota.GetStringBytes("transportador", "nome"))
			novoIbm.CD_CHAVE_NOTA_FISCAL = stringOrDefault(nota.GetStringBytes("chavexml"))

			// Processar produtos da nota
			produtos := nota.GetArray("produtos")
			for _, produto := range produtos {
				novoIbm.CD_EAN_PRODUTO = stringOrDefault(produto.GetStringBytes("ean"))
				qtdStr := stringOrDefault(produto.GetStringBytes("qtd"))
				qtdFloat, err := strconv.ParseFloat(qtdStr, 64)
				if err != nil {
					span.RecordError(err)
					log.Printf("Erro ao converter quantidade: %v", err)
					continue
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

				dataSave = append(dataSave, novoIbm)
			}
		}
	}

	if err := uc.Repo.SalvarCompras(ctx, dataSave, uuid); err != nil {
		log.Printf("Erro ao salvar compras: %v", err)
		span.RecordError(err)
		return false, err
	}
	log.Printf("Compras salvas com sucesso: %d registros", len(dataSave))
	// Se tudo ocorrer bem, retornar true
	// e nil para o erro

	return true, nil
}
