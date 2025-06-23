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
exports.QueryDefault = void 0;
const dbConnect_1 = require("@services/dbConnect");
class QueryDefault {
    schema() {
        const val = {};
        return val;
    }
    save(dados) {
        return __awaiter(this, void 0, void 0, function* () {
            const transaction = yield dbConnect_1.db_connect.transaction();
            let dadosSave;
            try {
                dadosSave = yield this.schema().create(dados, { transaction });
                yield transaction.commit();
            }
            catch (error) {
                yield transaction.rollback();
                console.error(`Erro ao salvar objeto: ${error.message}`);
            }
            return dadosSave;
        });
    }
    deletePorObjeto(where) {
        return __awaiter(this, void 0, void 0, function* () {
            // const transaction = await db_connect.transaction()
            // let dadosDelete
            try {
                const dadosDelete = yield this.schema().destroy({ where: where });
                return dadosDelete;
                // await transaction.commit()
            }
            catch (error) {
                //await transaction.rollback()
                //console.error(`Erro ao deletar objeto: ${(error as any).message}`)
            }
        });
    }
    listAll() {
        return __awaiter(this, void 0, void 0, function* () {
            return yield this.schema().findAll({
                raw: true
            });
        });
    }
    listCont() {
        return __awaiter(this, void 0, void 0, function* () {
            return yield this.schema().count();
        });
    }
    bulkCreate(dados) {
        return __awaiter(this, void 0, void 0, function* () {
            return yield this.schema().bulkCreate(dados);
        });
    }
}
exports.QueryDefault = QueryDefault;
