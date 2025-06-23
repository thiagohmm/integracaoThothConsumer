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

type VendaUseCase struct {
	Repo entities.VendaRepository
}

func NewVendaUseCase(repo entities.VendaRepository) *VendaUseCase {
	return &VendaUseCase{Repo: repo}
}

func stringOrDefaultWithDefault(b []byte, def string) string {
	if b == nil || len(b) == 0 {
		return def
	}
	return string(b)
}

// getIntFromMixed tenta extrair um inteiro do campo, seja ele numérico ou uma string numérica.
func getIntFromMixed(produto *fastjson.Value, key string) int {
	val := produto.Get(key)
	if val == nil {
		return 0
	}
	switch val.Type() {
	case fastjson.TypeNumber:
		return val.GetInt()
	case fastjson.TypeString:
		s := stringOrDefaultWithDefault(val.GetStringBytes(), "0")
		if i, err := strconv.Atoi(s); err == nil {
			return i
		}
	}
	return 0
}

// getStringFromMixed tenta extrair uma string do campo, seja ele numérico ou uma string.
func getStringFromMixed(produto *fastjson.Value, key, def string) string {
	val := produto.Get(key)
	if val == nil {
		return def
	}
	switch val.Type() {
	case fastjson.TypeString:
		return stringOrDefaultWithDefault(val.GetStringBytes(), def)
	case fastjson.TypeNumber:
		return strconv.Itoa(val.GetInt())
	default:
		return def
	}
}

func (uc *VendaUseCase) ProcessarVenda(ctx context.Context, vendaData map[string]interface{}) error {
	tracer := otel.Tracer("ProcessarVenda")

	uuid, ok := ctx.Value("uuid").(string)
	if !ok {
		return fmt.Errorf("UUID não encontrado no contexto")
	}
	ctxDeleteTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)

	defer cancel()
	ctx, span := tracer.Start(ctx, uuid)
	defer span.End()

	var dataSave []entities.Venda

	vendaJSON, err := json.Marshal(vendaData)
	if err != nil {
		span.RecordError(fmt.Errorf("Erro ao converter mapa em JSON: %v", err))
		log.Printf("Erro ao converter mapa em JSON: %v", err)
		return err
	}

	var p fastjson.Parser
	v, err := p.ParseBytes(vendaJSON)
	if err != nil {
		log.Printf("Erro ao parsear JSON: %v", err)
		span.RecordError(fmt.Errorf("Erro ao parsear JSON: %v", err))
		return err
	}

	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		log.Printf("Falha ao carregar localização: %v", err)
		loc = time.Local
		span.RecordError(err)
	}
	dtNow := string(time.Now().In(loc).Format("2006-01-02T15:04:05.000Z"))

	ibms := v.GetArray("vendas", "ibms")
	if ibms == nil {
		log.Printf("IBMs não encontrados no objeto de venda")
		span.RecordError(fmt.Errorf("IBMs não encontrados"))
		fmt.Errorf("IBMs não encontrados")
		return err
	}

	for _, vendaIbms := range ibms {
		nroStr := sanitizeIBM(stringOrDefaultWithDefault(vendaIbms.GetStringBytes("nro"), ""))
		dtStr := getStringFromMixed(v.Get("vendas"), "dtavenda", "")
		dtEntrada, err := strconv.ParseInt(dtStr, 10, 64)
		if err != nil {
			log.Printf("Erro ao converter DT_Venda para int64: %v", err)
			span.RecordError(fmt.Errorf("erro ao converter DT_Venda para int64: %v", err))
			return err
		}

		existVendasZerada, err := uc.Repo.CheckIfExistsZerado(ctxDeleteTimeout, nroStr, dtStr)
		if err != nil {
			log.Printf("Erro ao verificar se a venda existe: %v", err)
			span.RecordError(fmt.Errorf("erro ao verificar se a venda existe: %v", err))
		}

		if existVendasZerada {
			if err := uc.Repo.DeleteByIBMVendaZerado(ctxDeleteTimeout, nroStr, dtStr, uuid); err != nil {
				log.Printf("Erro ao deletar IBM venda zerada: %s, erro: %v", nroStr, err)
				span.RecordError(fmt.Errorf("erro ao deletar IBM venda zerada: %s, erro: %v", nroStr, err))
				return err
			}
		}

		existVendas, err := uc.Repo.CheckIfExists(ctxDeleteTimeout, nroStr, dtStr)
		if err != nil {
			log.Printf("Erro ao verificar se a venda existe: %v", err)
			span.RecordError(fmt.Errorf("Erro ao verificar se a venda existe: %v", err))
		}

		if existVendas {
			if err := uc.Repo.DeleteByIBMVenda(ctxDeleteTimeout, nroStr, dtStr, uuid); err != nil {
				log.Printf("Erro ao deletar IBM venda: %s, erro: %v", nroStr, err)
				span.RecordError(fmt.Errorf("Erro ao deletar IBM venda: %s, erro: %v", nroStr, err))
				return err
			}
		}

		novoIbm := entities.Venda{
			DT_TRANSACAO:      dtEntrada,
			CD_IBM_LOJA:       nroStr,
			RAZAO_SOCIAL_LOJA: stringOrDefaultWithDefault(vendaIbms.GetStringBytes("razao"), ""),
			NM_SISTEMA:        stringOrDefaultWithDefault(vendaIbms.GetStringBytes("app"), ""),
			DT_ARQUIVO:        stringOrDefaultWithDefault(v.GetStringBytes("vendas", "dtaenvio"), ""),
			SRC_LOAD:          "API/Integração/Thoth",
			DT_LOAD:           dtNow,
		}

		vendas := vendaIbms.GetArray("vendas")
		for _, venda := range vendas {
			novoIbm.HR_INICIO_TRANSACAO = stringOrDefaultWithDefault(venda.GetStringBytes("ini"), "")
			novoIbm.HR_FIM_TRANSACAO = stringOrDefaultWithDefault(venda.GetStringBytes("fim"), "")
			novoIbm.CD_TRANSACAO = getStringFromMixed(venda, "doc", "")
			novoIbm.CPF_CNPJ_CLIENTE = stringOrDefaultWithDefault(venda.GetStringBytes("cpfcnpj"), "")
			pagto := strings.TrimSpace(stringOrDefaultWithDefault(venda.GetStringBytes("formapagto"), ""))
			if len(pagto) > 50 {
				pagto = pagto[:49]
			}
			novoIbm.NM_FORMA_PAGAMENTO = pagto
			novoIbm.NM_BANDEIRA = stringOrDefaultWithDefault(venda.GetStringBytes("bandeira"), "")
			novoIbm.CD_CCF = stringOrDefaultWithDefault(venda.GetStringBytes("ccf"), "")
			novoIbm.CD_MODELO_DOCTO = stringOrDefaultWithDefault(venda.GetStringBytes("moddoc"), "")

			produtos := venda.GetArray("produtos")
			for _, produto := range produtos {
				novoIbm.CD_EAN_PRODUTO = getStringFromMixed(produto, "ean", "")
				novoIbm.QT_PRODUTO = stringOrDefaultWithDefault(produto.GetStringBytes("qtd"), "")
				novoIbm.VL_PRECO_UNITARIO = stringOrDefaultWithDefault(produto.GetStringBytes("preco"), "")
				novoIbm.VL_IMPOSTO = stringOrDefaultWithDefault(produto.GetStringBytes("imposto"), "")
				novoIbm.VL_FATURADO = getStringFromMixed(produto, "total", "")
				novoIbm.VL_CUSTO_UNITARIO = stringOrDefaultWithDefault(produto.GetStringBytes("custo"), "")
				novoIbm.CD_DEPARTAMENTO = stringOrDefaultWithDefault(produto.GetStringBytes("dep"), "")
				tipoInt := getIntFromMixed(produto, "tipo")
				itemTransInt := getIntFromMixed(produto, "trans")
				tipo := strconv.Itoa(tipoInt)
				itemTrans := strconv.Itoa(itemTransInt)

				novoIbm.CD_TP_PRODUTO = tipo
				novoIbm.DS_PRODUTO = stringOrDefaultWithDefault(produto.GetStringBytes("descricao"), "")
				novoIbm.CD_PROMOCAO = stringOrDefaultWithDefault(produto.GetStringBytes("codmix"), "")
				novoIbm.CD_EAN_EMBALAGEM = stringOrDefaultWithDefault(produto.GetStringBytes("eanpack"), "")
				novoIbm.CD_TP_TRANSACAO = itemTrans
				novoIbm.VL_DESCONTO = stringOrDefaultWithDefault(produto.GetStringBytes("desconto"), "")
				novoIbm.CD_ITEM_TRANSACAO = itemTrans

				dataSave = append(dataSave, novoIbm)
			}
		}
	}

	if err := uc.Repo.SalvarVendas(ctx, dataSave, uuid); err != nil {
		log.Printf("Erro ao salvar vendas: %v", err)
		span.RecordError(err)
		return err
	}

	return nil
}
