"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.verifyToken = void 0;
const jsonwebtoken_1 = __importDefault(require("jsonwebtoken"));
const requestResultStatusEnum_1 = require("@enumerators/requestResultStatusEnum");
const verifyToken = (req, res) => {
    try {
        const bearer = req.headers.authorization;
        const token = bearer.split(' ');
        const decoded = jsonwebtoken_1.default.decode(token[1]);
        const dateNow = Math.floor(Date.now() / 1000);
        if (decoded.exp <= dateNow) {
            return res.status(requestResultStatusEnum_1.enumRequestResultStatus.expiredToken).json({
                error: 'Token expirado',
                message: `Acesse o endpoint ${process.env.ADW_ATENAGRPNOS_GENERATE_TOKEN} para gerar um novo, com seu usuário e senha. `
            });
        }
        res.status(requestResultStatusEnum_1.enumRequestResultStatus.success).json({ message: 'Sucesso na transação' });
    }
    catch (error) {
        return res.status(requestResultStatusEnum_1.enumRequestResultStatus.forcarNovoLogin).json({
            error: 'Falha ao decodificar o token ou não foi enviado',
            message: 'Verifique se possui um token válido'
        });
    }
};
exports.verifyToken = verifyToken;
