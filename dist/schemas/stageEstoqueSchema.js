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
exports.StageEstoqueQuery = exports.StageEstoqueSchema = void 0;
const dbConnect_1 = require("../services/dbConnect");
const queryDefault_1 = require("../utils/queryDefault");
const sequelize_1 = require("sequelize");
const tableName = 'STG_WS_ESTOQUE';
exports.StageEstoqueSchema = dbConnect_1.db_connect.define(tableName, {
    ID_STG_WS_ESTOQUE: {
        allowNull: false,
        autoIncrement: true,
        primaryKey: true,
        type: sequelize_1.DataTypes.INTEGER
    },
    DT_ESTOQUE: {
        allowNull: true,
        type: sequelize_1.DataTypes.INTEGER,
    },
    DT_ULTIMA_COMPRA: {
        allowNull: true,
        type: sequelize_1.DataTypes.INTEGER,
    },
    CD_IBM_LOJA: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50)
    },
    RAZAO_SOCIAL_LOJA: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(255)
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
    VL_PRECO_UNITARIO: {
        allowNull: true,
        type: sequelize_1.DataTypes.DOUBLE
    },
    QT_INVENTARIO_ENTRADA: {
        allowNull: true,
        type: sequelize_1.DataTypes.DOUBLE
    },
    QT_INVENTARIO_SAIDA: {
        allowNull: true,
        type: sequelize_1.DataTypes.DOUBLE
    },
    QT_INICIAL_PRODUTO: {
        allowNull: true,
        type: sequelize_1.DataTypes.DOUBLE
    },
    QT_FINAL_PRODUTO: {
        allowNull: true,
        type: sequelize_1.DataTypes.DOUBLE
    },
    VL_TOTAL_ESTOQUE: {
        allowNull: true,
        type: sequelize_1.DataTypes.DOUBLE
    },
    VL_CUSTO_MEDIO: {
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
        type: sequelize_1.DataTypes.DATE,
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
class StageEstoqueQuery extends queryDefault_1.QueryDefault {
    schema() { return exports.StageEstoqueSchema; }
    listByIbmByDataTransacao(ibm, data_estoque) {
        return __awaiter(this, void 0, void 0, function* () {
            return yield this.schema().findAll({
                where: { CD_IBM_LOJA: ibm, DT_ESTOQUE: data_estoque },
                raw: true
            });
        });
    }
}
exports.StageEstoqueQuery = StageEstoqueQuery;
