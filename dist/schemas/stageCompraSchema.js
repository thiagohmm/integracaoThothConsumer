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
exports.StageCompraQuery = exports.StageCompraSchema = void 0;
const dbConnect_1 = require("../services/dbConnect");
const queryDefault_1 = require("../utils/queryDefault");
const sequelize_1 = require("sequelize");
const tableName = 'STG_WS_COMPRA';
exports.StageCompraSchema = dbConnect_1.db_connect.define(tableName, {
    ID_STG_WS_COMPRA: {
        allowNull: false,
        autoIncrement: true,
        primaryKey: true,
        type: sequelize_1.DataTypes.DOUBLE
    },
    DT_EMISSAO_NOTA: {
        allowNull: true,
        type: sequelize_1.DataTypes.DOUBLE
    },
    DT_ENTRADA: {
        allowNull: true,
        type: sequelize_1.DataTypes.DOUBLE,
    },
    CD_IBM_LOJA: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50),
    },
    RAZAO_SOCIAL_LOJA: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(255)
    },
    CNPJ_FORNECEDOR: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    NM_FORNECEDOR: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(255)
    },
    CNPJ_TRANSPORTADORA: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    NM_TRANSPORTADORA: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(255)
    },
    CD_PRODUTO_FORNECEDOR: {
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
    DS_UN_MEDIDA: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    DS_UN_MEDIDA_CONVERTIDA: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    NR_SERIE_NOTA: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(10)
    },
    NR_NOTA_FISCAL: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    CD_ITEM_NOTA_FISCAL: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(10)
    },
    CD_CHAVE_NOTA_FISCAL: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(255),
    },
    CD_NCM: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    TIPO_FRETE: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    VL_ULTIMO_CUSTO: {
        allowNull: true,
        type: sequelize_1.DataTypes.DOUBLE
    },
    QT_PRODUTO: {
        allowNull: true,
        type: sequelize_1.DataTypes.DOUBLE
    },
    QT_PRODUTO_CONVERTIDA: {
        allowNull: true,
        type: sequelize_1.DataTypes.DOUBLE
    },
    VL_PRECO_COMPRA: {
        allowNull: true,
        type: sequelize_1.DataTypes.DOUBLE
    },
    VL_PIS: {
        allowNull: true,
        type: sequelize_1.DataTypes.DOUBLE
    },
    VL_ALIQUOTA_PIS: {
        allowNull: true,
        type: sequelize_1.DataTypes.DOUBLE
    },
    VL_COFINS: {
        allowNull: true,
        type: sequelize_1.DataTypes.DOUBLE
    },
    VL_ALIQUOTA_COFINS: {
        allowNull: true,
        type: sequelize_1.DataTypes.DOUBLE
    },
    VL_ICMS: {
        allowNull: true,
        type: sequelize_1.DataTypes.DOUBLE
    },
    VL_ALIQUOTA_ICMS: {
        allowNull: true,
        type: sequelize_1.DataTypes.DOUBLE
    },
    VL_IPI: {
        allowNull: true,
        type: sequelize_1.DataTypes.DOUBLE
    },
    VL_ALIQUOTA_IPI: {
        allowNull: true,
        type: sequelize_1.DataTypes.DOUBLE
    },
    QT_PESO: {
        allowNull: true,
        type: sequelize_1.DataTypes.DOUBLE
    },
    VL_FRETE: {
        allowNull: true,
        type: sequelize_1.DataTypes.DOUBLE
    },
    VL_TOTAL_ICMS: {
        allowNull: true,
        type: sequelize_1.DataTypes.DOUBLE
    },
    VL_TOTAL_IPI: {
        allowNull: true,
        type: sequelize_1.DataTypes.DOUBLE
    },
    VL_TOTAL_COMPRA: {
        allowNull: true,
        type: sequelize_1.DataTypes.DOUBLE
    },
    NM_SISTEMA: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    SRC_LOAD: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    DT_LOAD: {
        allowNull: true,
        type: sequelize_1.DataTypes.DATE
    },
    FL_CARGA_HISTORICA: {
        allowNull: true,
        type: sequelize_1.DataTypes.DOUBLE
    },
    CD_IBM_LOJA_EAGLE: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    CD_EAN_PRODUTO_EAGLE: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
}, {
    freezeTableName: true,
    timestamps: false,
    updatedAt: 'data_alteracao',
    createdAt: 'data_criacao'
});
class StageCompraQuery extends queryDefault_1.QueryDefault {
    schema() { return exports.StageCompraSchema; }
    listByIbmByDataTransacao(ibm, data_entrada) {
        return __awaiter(this, void 0, void 0, function* () {
            return yield this.schema().findAll({
                where: { CD_IBM_LOJA: ibm, DT_ENTRADA: data_entrada },
                raw: true
            });
        });
    }
}
exports.StageCompraQuery = StageCompraQuery;
