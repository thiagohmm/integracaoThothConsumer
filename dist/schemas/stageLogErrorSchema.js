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
exports.StageLogErrorQuery = exports.StageLogErrorSchema = void 0;
const dbConnect_1 = require("../services/dbConnect");
const queryDefault_1 = require("../utils/queryDefault");
const sequelize_1 = require("sequelize");
const tableName = 'STG_WS_ERRO';
exports.StageLogErrorSchema = dbConnect_1.db_connect.define(tableName, {
    DATA: {
        allowNull: false,
        type: sequelize_1.DataTypes.DOUBLE
    },
    CD_IBM_LOJA: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(50),
    },
    MENSAGEM_ERRO: {
        allowNull: false,
        type: sequelize_1.DataTypes.STRING(255)
    },
    JSON: {
        allowNull: true,
        type: sequelize_1.DataTypes.BLOB
    },
    INTERFACE: {
        allowNull: true,
        type: sequelize_1.DataTypes.STRING(255)
    },
    DATA_ERRO: {
        allowNull: true,
        type: sequelize_1.DataTypes.DATE
    }
}, {
    freezeTableName: true,
    timestamps: false,
    primaryKey: false,
    hasPrimaryKey: false,
    indexes: [],
});
exports.StageLogErrorSchema.removeAttribute('id');
class StageLogErrorQuery extends queryDefault_1.QueryDefault {
    constructor() {
        super();
    }
    schema() { return exports.StageLogErrorSchema; }
    insertErrorLog(data) {
        return __awaiter(this, void 0, void 0, function* () {
            return yield this.schema().create(data);
        });
    }
}
exports.StageLogErrorQuery = StageLogErrorQuery;
