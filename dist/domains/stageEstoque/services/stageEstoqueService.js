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
exports.PersisteDadosEstoqueThoth = PersisteDadosEstoqueThoth;
const dateParser_1 = require("@utils/dateParser");
const stageLogErrorSchema_1 = require("@schemas/stageLogErrorSchema");
const sequelize_1 = require("sequelize");
function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}
function stageEstoqueMapearDados(estoque, stageEstoqueQuery) {
    return __awaiter(this, void 0, void 0, function* () {
        var _a, _b;
        const arquivoEstoque = [];
        const novoIbm = {};
        let novoIbmEstoque;
        if (typeof estoque === 'string') {
            try {
                novoIbmEstoque = JSON.parse(estoque);
            }
            catch (e) {
                console.error('Erro ao fazer parse da string JSON:', e);
            }
        }
        else if (typeof estoque === 'object') {
            novoIbmEstoque = estoque; // Já é um objeto, não precisa fazer parse
        }
        else {
            console.error('A variável "estoque" não é uma string JSON válida nem um objeto.');
        }
        let saveCounter = 0;
        for (const ibmEstoque of novoIbmEstoque.estoque.ibms || []) {
            try {
                // Converter ibm.nro para string
                let nroStr = ibmEstoque.nro.toString().replace(/\D/g, '').trim();
                // Verificar se o comprimento é maior que 10 e remover zeros à esquerda
                if (nroStr.length > 10) {
                    nroStr = nroStr.replace(/^0+/, '');
                }
                nroStr = nroStr.padStart(10, '0').slice(-10);
                novoIbm.CD_IBM_LOJA = nroStr;
                novoIbm.RAZAO_SOCIAL_LOJA = ibmEstoque.razao;
                novoIbm.DT_ESTOQUE = parseInt(String(novoIbmEstoque.estoque.dtaestoque || '0').replace(/\D/g, '').trim());
                novoIbm.NM_SISTEMA = ibmEstoque.app;
                novoIbm.SRC_LOAD = 'API/Integração/Thoth';
                novoIbm.DT_LOAD = dateParser_1.DateParser.Now();
                if (ibmEstoque.produtos.length === 0) {
                    arquivoEstoque.push(novoIbm);
                    yield stageEstoqueQuery.save(novoIbm);
                    saveCounter++;
                    if (saveCounter % 100 === 0) {
                        yield sleep(500);
                    }
                    continue;
                }
                for (const produto of ibmEstoque.produtos || []) {
                    novoIbm.CD_EAN_PRODUTO = produto.ean;
                    novoIbm.CD_TP_PRODUTO = produto.tipo.toString();
                    novoIbm.DT_ULTIMA_COMPRA = Number(produto.dtacompra);
                    novoIbm.DS_PRODUTO = produto.descricao;
                    novoIbm.VL_PRECO_UNITARIO = Number(produto.preco || 0);
                    novoIbm.QT_INVENTARIO_ENTRADA = Number(produto.qtdentrada || 0);
                    novoIbm.QT_INVENTARIO_SAIDA = Number(produto.qtdsaida || 0);
                    novoIbm.QT_INICIAL_PRODUTO = Number(produto.qtdini || 0);
                    novoIbm.QT_FINAL_PRODUTO = Number(produto.qtdfim || 0);
                    novoIbm.VL_TOTAL_ESTOQUE = Number(produto.vlrfim || 0);
                    novoIbm.VL_CUSTO_MEDIO = Number(produto.vlrmedio || 0);
                    if (novoIbm.VL_CUSTO_MEDIO > Number.MAX_SAFE_INTEGER) {
                        novoIbm.VL_CUSTO_MEDIO = parseFloat(Number(produto.vlrmedio.toString().slice(0, 4)).toFixed(2));
                    }
                    arquivoEstoque.push(yield stageEstoqueQuery.save(novoIbm));
                    saveCounter++;
                    if (saveCounter % 100 === 0) {
                        yield sleep(500);
                    }
                }
            }
            catch (error) {
                console.error(`Erro ao processar IBM Estoque: ${ibmEstoque.nro}, erro: ${error.message}`);
                const stageLogErrorQuery = new stageLogErrorSchema_1.StageLogErrorQuery();
                const dadosErro = {
                    DATA: parseInt((_b = (_a = novoIbm.DT_ESTOQUE) === null || _a === void 0 ? void 0 : _a.toString()) !== null && _b !== void 0 ? _b : '0'),
                    CD_IBM_LOJA: novoIbm.CD_IBM_LOJA,
                    MENSAGEM_ERRO: error.message,
                    JSON: JSON.stringify(novoIbmEstoque),
                    INTERFACE: 'ESTOQUE',
                    DATA_ERRO: new Date()
                };
                yield stageLogErrorQuery.save(dadosErro);
            }
        }
        return arquivoEstoque;
    });
}
function PersisteDadosEstoqueThoth(estoque, stageEstoqueQuery) {
    return __awaiter(this, void 0, void 0, function* () {
        let persisteEstoque;
        if (typeof estoque === 'string') {
            try {
                persisteEstoque = JSON.parse(estoque);
            }
            catch (e) {
                console.error('Erro ao fazer parse da string JSON:', e);
            }
        }
        else if (typeof estoque === 'object') {
            persisteEstoque = estoque; // Já é um objeto, não precisa fazer parse
        }
        else {
            console.error('A variável "estoque" não é uma string JSON válida nem um objeto.');
        }
        let deleteCounter = 0;
        for (const estoqueIbms of persisteEstoque.estoque.ibms || []) {
            try {
                const nroStr = String(estoqueIbms.nro).replace(/\D/g, '').trim();
                const dtaentrada = persisteEstoque.estoque.dtaestoque;
                const dtStr = String(dtaentrada !== undefined && dtaentrada !== null ? dtaentrada : '0').replace(/\D/g, '').trim();
                yield stageEstoqueQuery.deletePorObjeto({ CD_IBM_LOJA: nroStr, DT_ESTOQUE: parseInt(dtStr || '0') });
                deleteCounter++;
                if (deleteCounter % 100 === 0) {
                    yield sleep(500);
                }
            }
            catch (error) {
                if (error instanceof sequelize_1.TimeoutError) {
                    console.error(`Timeout ao deletar IBM estoque: ${estoqueIbms.nro}, erro: ${error.message}`);
                }
                else {
                    console.error(`Erro ao deletar IBM estoque: ${estoqueIbms.nro}, erro: ${error.message}`);
                }
            }
        }
        return yield stageEstoqueMapearDados(estoque, stageEstoqueQuery);
    });
}
