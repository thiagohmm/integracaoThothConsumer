"use strict";
/* eslint-disable indent */
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
exports.PersisteDadosComprasThoth = PersisteDadosComprasThoth;
const stageLogErrorSchema_1 = require("@schemas/stageLogErrorSchema");
const dateParser_1 = require("@utils/dateParser");
const sequelize_1 = require("sequelize");
function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}
function stageCompraMapearDados(compra, stageCompraQuery) {
    return __awaiter(this, void 0, void 0, function* () {
        var _a, _b;
        const arquivoCompra = [];
        const novoIbm = {};
        let novoIbmCompra;
        if (typeof compra === 'string') {
            try {
                novoIbmCompra = JSON.parse(compra);
            }
            catch (e) {
                console.error('Erro ao fazer parse da string JSON:', e);
            }
        }
        else if (typeof compra === 'object') {
            novoIbmCompra = compra; // Já é um objeto, não precisa fazer parse
        }
        else {
            console.error('A variável "compra" não é uma string JSON válida nem um objeto.');
        }
        let saveCounter = 0;
        for (const ibm of novoIbmCompra.compras.ibms || []) {
            try {
                // Converter ibm.nro para string
                let nroStr = ibm.nro.toString().replace(/\D/g, '').trim();
                // Verificar se o comprimento é maior que 10 e remover zeros à esquerda
                if (nroStr.length > 10) {
                    nroStr = nroStr.replace(/^0+/, '');
                }
                nroStr = nroStr.padStart(10, '0').slice(-10);
                novoIbm.CD_IBM_LOJA = nroStr;
                novoIbm.RAZAO_SOCIAL_LOJA = ibm.razao;
                novoIbm.DT_ENTRADA = parseInt(String(novoIbmCompra.compras.dtaentrada).replace(/\D/g, '').trim());
                novoIbm.NM_SISTEMA = ibm.app;
                novoIbm.SRC_LOAD = 'API/Integração/Thoth';
                novoIbm.DT_LOAD = dateParser_1.DateParser.Now();
                if (ibm.notas.length <= 0) {
                    arquivoCompra.push(novoIbm);
                    yield stageCompraQuery.save(novoIbm);
                    saveCounter++;
                    if (saveCounter % 100 === 0) {
                        yield sleep(500);
                    }
                    continue;
                }
                for (const nota of ibm.notas || []) {
                    novoIbm.NR_NOTA_FISCAL = nota.nro;
                    novoIbm.NR_SERIE_NOTA = nota.serie;
                    novoIbm.DT_EMISSAO_NOTA = parseInt(nota.emissao || '0');
                    novoIbm.CNPJ_FORNECEDOR = nota.fornecedor.cnpj;
                    novoIbm.NM_FORNECEDOR = nota.fornecedor.nome;
                    novoIbm.QT_PESO = nota.total.peso;
                    novoIbm.VL_TOTAL_IPI = nota.total.vlripi;
                    novoIbm.VL_TOTAL_ICMS = nota.total.vlricms;
                    novoIbm.VL_TOTAL_COMPRA = nota.total.vlrnota;
                    novoIbm.CNPJ_TRANSPORTADORA = nota.transportador.cnpj;
                    novoIbm.NM_TRANSPORTADORA = nota.transportador.nome;
                    novoIbm.CD_CHAVE_NOTA_FISCAL = nota.chavexml;
                    if (nota.produtos.length <= 0) {
                        arquivoCompra.push(novoIbm);
                        yield stageCompraQuery.save(novoIbm);
                        saveCounter++;
                        if (saveCounter % 100 === 0) {
                            yield sleep(500);
                        }
                        continue;
                    }
                    for (const produto of nota.produtos || []) {
                        novoIbm.CD_EAN_PRODUTO = produto.ean;
                        novoIbm.QT_PRODUTO = produto.qtd || 0.00;
                        novoIbm.VL_PRECO_COMPRA = produto.preco || 0.00;
                        novoIbm.DS_PRODUTO = produto.descricao;
                        novoIbm.CD_TP_PRODUTO = (produto.tipo).toString();
                        novoIbm.VL_ALIQUOTA_IPI = produto.impostos.ipi.aliquota;
                        novoIbm.VL_IPI = produto.impostos.ipi.vlr;
                        novoIbm.VL_ALIQUOTA_ICMS = produto.impostos.icms.aliquota;
                        novoIbm.VL_ICMS = produto.impostos.icms.vlr;
                        novoIbm.VL_ALIQUOTA_PIS = produto.impostos.pis.aliquota;
                        novoIbm.VL_PIS = produto.impostos.pis.vlr;
                        novoIbm.VL_ALIQUOTA_COFINS = produto.impostos.cofins.aliquota;
                        novoIbm.VL_COFINS = produto.impostos.cofins.vlr;
                        novoIbm.CD_NCM = produto.ncm;
                        novoIbm.CD_ITEM_NOTA_FISCAL = (produto.linha !== undefined && produto.linha !== '') ? String(produto.linha) : null;
                        novoIbm.CD_PRODUTO_FORNECEDOR = produto.codfornec;
                        novoIbm.QT_PRODUTO_CONVERTIDA = produto.qtdenf || 0;
                        novoIbm.DS_UN_MEDIDA_CONVERTIDA = produto.unconv;
                        novoIbm.DS_UN_MEDIDA = produto.un;
                        novoIbm.VL_ULTIMO_CUSTO = produto.ultcusto || 0;
                        arquivoCompra.push(yield stageCompraQuery.save(novoIbm));
                        saveCounter++;
                        if (saveCounter % 100 === 0) {
                            yield sleep(500);
                        }
                    }
                }
            }
            catch (error) {
                console.error(`Erro ao processar IBM Compra: ${ibm.nro}, erro: ${error.message}`);
                // Dependendo da sua necessidade, você pode adicionar lógica adicional aqui,
                // como registrar a IBM com erro em um arquivo de log ou banco de dados
                const stageLogErrorQuery = new stageLogErrorSchema_1.StageLogErrorQuery();
                const dadosErro = {
                    DATA: parseInt((_b = (_a = novoIbm.DT_ENTRADA) === null || _a === void 0 ? void 0 : _a.toString()) !== null && _b !== void 0 ? _b : '0'),
                    CD_IBM_LOJA: novoIbm.CD_IBM_LOJA,
                    MENSAGEM_ERRO: error.message,
                    JSON: JSON.stringify(novoIbmCompra),
                    INTERFACE: 'COMPRA',
                    DATA_ERRO: new Date()
                };
                yield stageLogErrorQuery.save(dadosErro);
            }
        }
        return arquivoCompra;
    });
}
function PersisteDadosComprasThoth(compra, stageCompraQuery) {
    return __awaiter(this, void 0, void 0, function* () {
        let persisteCompra;
        if (typeof compra === 'string') {
            try {
                persisteCompra = JSON.parse(compra);
            }
            catch (e) {
                console.error('Erro ao fazer parse da string JSON:', e);
            }
        }
        else if (typeof compra === 'object') {
            persisteCompra = compra; // Já é um objeto, não precisa fazer parse
        }
        else {
            console.error('A variável "compra" não é uma string JSON válida nem um objeto.');
        }
        let deleteCounter = 0;
        for (const compraIbms of persisteCompra.compras.ibms || []) {
            try {
                let nroStr = String(compraIbms.nro).replace(/\D/g, '').trim();
                if (nroStr.length > 10) {
                    nroStr = nroStr.replace(/^0+/, '');
                }
                nroStr = nroStr.padStart(10, '0').slice(-10);
                const dtaentrada = persisteCompra.compras.dtaentrada;
                const dtStr = String(dtaentrada !== undefined && dtaentrada !== null ? dtaentrada : '0').replace(/\D/g, '').trim();
                console.log(`Deletando IBM compra: ${nroStr}, data: ${dtStr}`);
                yield stageCompraQuery.deletePorObjeto({ CD_IBM_LOJA: nroStr, DT_ENTRADA: parseInt(dtStr || '0') });
                deleteCounter++;
                if (deleteCounter % 100 === 0) {
                    yield sleep(500);
                }
            }
            catch (error) {
                if (error instanceof sequelize_1.TimeoutError) {
                    console.error(`Timeout ao deletar IBM compra: ${compraIbms.nro}, erro: ${error.message}`);
                }
                else {
                    console.error(`Erro ao deletar IBM compra: ${compraIbms.nro}, erro: ${error.message}`);
                }
            }
        }
        return yield stageCompraMapearDados(compra, stageCompraQuery);
    });
}
