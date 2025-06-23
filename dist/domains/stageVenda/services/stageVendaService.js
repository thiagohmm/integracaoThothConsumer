"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.PersisteDadosVendasThoth = PersisteDadosVendasThoth;
const stageLogErrorSchema_1 = require("@schemas/stageLogErrorSchema");
const dateParser_1 = require("@utils/dateParser");
const sequelize_1 = require("sequelize");
function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}
function stageVendaMapearDados(venda, stageVendaQuery) {
    return __awaiter(this, void 0, void 0, function* () {
        var _a, _b, _c, _d;
        const arquivoVenda = [];
        const novoIbm = {};
        let novoIbmVenda;
        if (typeof venda === 'string') {
            try {
                novoIbmVenda = JSON.parse(venda);
            }
            catch (e) {
                console.error('Erro ao fazer parse da string JSON:', e);
            }
        }
        else if (typeof venda === 'object') {
            novoIbmVenda = venda; // Já é um objeto, não precisa fazer parse
        }
        else {
            console.error('A variável "venda" não é uma string JSON válida nem um objeto.');
        }
        let saveCounter = 0;
        for (const ibm of novoIbmVenda.vendas.ibms || []) {
            try {
                const dtaentrada = parseInt(novoIbmVenda.vendas.dtavenda || '0');
                const dtStr = String(dtaentrada !== undefined && dtaentrada !== null ? dtaentrada : '0').replace(/\D/g, '').trim();
                novoIbm.DT_TRANSACAO = parseInt(dtStr || '0');
                // Converter ibm.nro para string
                let nroStr = ibm.nro.toString().replace(/\D/g, '').trim();
                // Verificar se o comprimento é maior que 10 e remover zeros à esquerda
                if (nroStr.length > 10) {
                    nroStr = nroStr.replace(/^0+/, '');
                }
                nroStr = nroStr.padStart(10, '0').slice(-10);
                // Atribuir o valor resultante a novoIbm.CD_IBM_LOJA
                novoIbm.CD_IBM_LOJA = nroStr;
                novoIbm.RAZAO_SOCIAL_LOJA = ibm.razao;
                novoIbm.NM_SISTEMA = ibm.app;
                novoIbm.DT_ARQUIVO = novoIbmVenda.vendas.dtaenvio;
                novoIbm.SRC_LOAD = 'API/Integração/Thoth';
                novoIbm.DT_LOAD = dateParser_1.DateParser.Now();
                if (ibm.vendas.length < 0) {
                    arquivoVenda.push(novoIbm);
                    yield stageVendaQuery.save(novoIbm);
                    saveCounter++;
                    if (saveCounter % 100 === 0) {
                        yield sleep(500);
                    }
                    continue;
                }
                for (const venda of ibm.vendas || []) {
                    novoIbm.HR_INICIO_TRANSACAO = venda.ini;
                    novoIbm.HR_FIM_TRANSACAO = venda.fim;
                    novoIbm.CD_TRANSACAO = venda.doc;
                    novoIbm.CPF_CNPJ_CLIENTE = venda.cpfcnpj;
                    novoIbm.NM_FORMA_PAGAMENTO = venda.formapagto;
                    novoIbm.NM_BANDEIRA = venda.bandeira;
                    novoIbm.CD_CCF = venda.ccf;
                    novoIbm.CD_MODELO_DOCTO = venda.moddoc;
                    if (venda.produtos.length < 0) {
                        arquivoVenda.push(novoIbm);
                        yield stageVendaQuery.save(novoIbm);
                        saveCounter++;
                        if (saveCounter % 100 === 0) {
                            yield sleep(500);
                        }
                        continue;
                    }
                    for (const produto of venda.produtos || []) {
                        novoIbm.CD_EAN_PRODUTO = produto.ean;
                        novoIbm.QT_PRODUTO = produto.qtd;
                        novoIbm.VL_PRECO_UNITARIO = produto.preco;
                        novoIbm.VL_IMPOSTO = produto.imposto;
                        novoIbm.VL_FATURADO = produto.total;
                        novoIbm.VL_CUSTO_UNITARIO = produto.custo;
                        novoIbm.CD_DEPARTAMENTO = produto.dep;
                        novoIbm.CD_TP_PRODUTO = produto.tipo;
                        novoIbm.DS_PRODUTO = produto.descricao;
                        novoIbm.CD_PROMOCAO = produto.codmix;
                        novoIbm.CD_EAN_EMBALAGEM = produto.eanpack;
                        novoIbm.CD_TP_TRANSACAO = produto.trans;
                        novoIbm.VL_DESCONTO = produto.desconto;
                        novoIbm.CD_ITEM_TRANSACAO = produto.trans;
                        try {
                            arquivoVenda.push(yield stageVendaQuery.save(novoIbm));
                            saveCounter++;
                            if (saveCounter % 100 === 0) {
                                yield sleep(500);
                            }
                        }
                        catch (error) {
                            console.error(`Erro ao salvar produto da venda: ${produto.ean}, erro: ${error.message}`);
                            const stageLogErrorQuery = new stageLogErrorSchema_1.StageLogErrorQuery();
                            const dadosErro = {
                                DATA: parseInt((_b = (_a = novoIbm.DT_TRANSACAO) === null || _a === void 0 ? void 0 : _a.toString()) !== null && _b !== void 0 ? _b : '0'),
                                CD_IBM_LOJA: novoIbm.CD_IBM_LOJA,
                                MENSAGEM_ERRO: error.message,
                                JSON: JSON.stringify(novoIbmVenda),
                                INTERFACE: 'VENDA',
                                DATA_ERRO: new Date()
                            };
                            yield stageLogErrorQuery.save(dadosErro);
                        }
                    }
                }
            }
            catch (error) {
                console.error(`Erro ao processar IBM Venda: ${ibm.nro}, erro: ${error.message}`);
                const stageLogErrorQuery = new stageLogErrorSchema_1.StageLogErrorQuery();
                const dadosErro = {
                    DATA: parseInt((_d = (_c = novoIbm.DT_TRANSACAO) === null || _c === void 0 ? void 0 : _c.toString()) !== null && _d !== void 0 ? _d : '0'),
                    CD_IBM_LOJA: novoIbm.CD_IBM_LOJA,
                    MENSAGEM_ERRO: error.message,
                    JSON: JSON.stringify(novoIbmVenda),
                    INTERFACE: 'VENDA',
                    DATA_ERRO: new Date()
                };
                yield stageLogErrorQuery.save(dadosErro);
            }
        }
        return arquivoVenda;
    });
}
function PersisteDadosVendasThoth(venda, stageVendaQuery) {
    return __awaiter(this, void 0, void 0, function* () {
        let persisteVenda;
        if (typeof venda === 'string') {
            try {
                persisteVenda = JSON.parse(venda);
            }
            catch (e) {
                console.error('Erro ao fazer parse da string JSON:', e);
            }
        }
        else if (typeof venda === 'object') {
            persisteVenda = venda; // Já é um objeto, não precisa fazer parse
        }
        else {
            console.error('A variável "venda" não é uma string JSON válida nem um objeto.');
        }
        let deleteCounter = 0;
        for (const vendasIbms of persisteVenda.vendas.ibms || []) {
            try {
                let nroStr = String(vendasIbms.nro).replace(/\D/g, '').trim();
                // Verificar se o comprimento é maior que 10 e remover zeros à esquerda
                if (nroStr.length > 10) {
                    nroStr = nroStr.replace(/^0+/, '');
                }
                nroStr = nroStr.padStart(10, '0').slice(-10);
                const dtaentrada = persisteVenda.vendas.dtavenda;
                const dtStr = String(dtaentrada !== undefined && dtaentrada !== null ? dtaentrada : '0').replace(/\D/g, '').trim();
                console.log(`Deletando IBM Venda: ${nroStr}, data: ${dtStr}`);
                yield stageVendaQuery.deletePorObjeto({ CD_IBM_LOJA: nroStr, DT_TRANSACAO: parseInt(dtStr || '0') });
                deleteCounter++;
                if (deleteCounter % 100 === 0) {
                    yield sleep(500);
                }
            }
            catch (error) {
                if (error instanceof sequelize_1.TimeoutError) {
                    console.error(`Timeout ao deletar IBM Venda: ${vendasIbms.nro}, erro: ${error.message}`);
                }
                else {
                    console.error(`Erro ao deletar IBM Venda: ${vendasIbms.nro}, erro: ${error.message}`);
                }
            }
        }
        return yield stageVendaMapearDados(venda, stageVendaQuery);
    });
}
