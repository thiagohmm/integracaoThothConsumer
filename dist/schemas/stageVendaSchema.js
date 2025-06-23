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
exports.StageVendaQuery = exports.StageVendaSchema = void 0;
const dbConnect_1 = require("../services/dbConnect");
const queryDefault_1 = require("../utils/queryDefault");
const sequelize_1 = require("sequelize");
const dateParser_1 = require("@utils/dateParser");
const tableName = 'STG_WS_VENDA';
exports.StageVendaSchema = dbConnect_1.db_connect.define(tableName, {
    ID_STG_WS_VENDA: {
        allowNull: false,
        autoIncrement: true,
        primaryKey: true,
        type: sequelize_1.DataTypes.NUMBER('38,0')
    },
    CD_CHAVE_TRANSACAO: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    DT_TRANSACAO: {
        allowNull: true,
        type: sequelize_1.DataTypes.NUMBER('8,0')
    },
    HR_INICIO_TRANSACAO: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(4),
    },
    HR_FIM_TRANSACAO: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(4),
    },
    CD_IBM_LOJA: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50),
    },
    RAZAO_SOCIAL_LOJA: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(255)
    },
    CD_DEPARTAMENTO: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    CPF_CNPJ_CLIENTE: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    CD_EAN_PRODUTO: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    DS_PRODUTO: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(255)
    },
    CD_TP_PRODUTO: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    CD_EAN_EMBALAGEM: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    CD_PROMOCAO: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    NM_FORMA_PAGAMENTO: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    NM_BANDEIRA: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    CD_MODELO_DOCTO: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    CD_CCF: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    CD_TRANSACAO: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    CD_ITEM_TRANSACAO: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(5)
    },
    CD_TP_TRANSACAO: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    QT_PRODUTO: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    VL_PRECO_UNITARIO: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    VL_IMPOSTO: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    VL_CUSTO_UNITARIO: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    VL_DESCONTO: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    VL_FATURADO: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    NM_SISTEMA: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    DT_ARQUIVO: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    SRC_LOAD: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(510)
    },
    DT_LOAD: {
        allowNull: true,
        defaultValue: dateParser_1.DateParser.Now(),
        type: sequelize_1.DataTypes.DATE
    },
    FL_CARGA_HISTORICA: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(1)
    },
    CD_IBM_LOJA_EAGLE: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    CD_EAN_PRODUTO_EAGLE: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    VL_CUSTO_EAGLE: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    NR_DDD_TELEFONE: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(3)
    },
    NR_TELEFONE: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(20)
    },
    VL_PIS: {
        allowNull: true,
        type: sequelize_1.DataTypes.NUMBER('18,3')
    },
    VL_COFINS: {
        allowNull: true,
        type: sequelize_1.DataTypes.NUMBER('18,3')
    },
    VL_ICMS: {
        allowNull: true,
        type: sequelize_1.DataTypes.NUMBER('18,3')
    },
}, {
    freezeTableName: true,
    timestamps: false,
    updatedAt: 'data_alteracao',
    createdAt: 'data_criacao'
});
class StageVendaQuery extends queryDefault_1.QueryDefault {
    schema() { return exports.StageVendaSchema; }
    listByIbmByDataTransacao(ibm, data_transacao) {
        return __awaiter(this, void 0, void 0, function* () {
            return yield this.schema().findAll({
                where: { CD_IBM_LOJA: ibm, DT_TRANSACAO: data_transacao },
                raw: true
            });
        });
    }
}
exports.StageVendaQuery = StageVendaQuery;
