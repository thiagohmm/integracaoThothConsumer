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
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
var _a;
Object.defineProperty(exports, "__esModule", { value: true });
exports.validarToken = exports.JwtMiddleware = void 0;
const requestResultStatusEnum_1 = require("@enumerators/requestResultStatusEnum");
const jsonwebtoken_1 = __importDefault(require("jsonwebtoken"));
const versao = (_a = process.env.VERSION_API) === null || _a === void 0 ? void 0 : _a.toString();
var JwtMiddleware;
(function (JwtMiddleware) {
    function check(req, res, next) {
        return __awaiter(this, void 0, void 0, function* () {
            const rota = { method: req.originalMethod, url: req.originalUrl };
            const lista_rota_livre = [
                { method: 'GET', url: `/${versao}` },
            ];
            let permite = false;
            permite = lista_rota_livre.filter(item => item.method == rota.method && item.url == rota.url).length > 0;
            if (permite) {
                next();
            }
            else {
                const token = (req.headers['authorization']) ? req.headers['authorization'].replace(/Bearer /g, '') : null;
                const { status, error, mensagem } = yield (0, exports.validarToken)(token);
                if (error) {
                    mensagem;
                    return res.status(status).send({ mensagem: mensagem, error: error });
                }
                next();
            }
        });
    }
    JwtMiddleware.check = check;
})(JwtMiddleware || (exports.JwtMiddleware = JwtMiddleware = {}));
const validarToken = (token) => {
    const decoded = jsonwebtoken_1.default.decode(token[1]);
    const dateNow = Math.floor(Date.now() / 1000);
    if (decoded.exp <= dateNow) {
        return { status: requestResultStatusEnum_1.enumRequestResultStatus.expiredToken, error: 'Token expirado', message: `Acesse o endpoint ${process.env.ADW_ATENAGRPNOS_GENERATE_TOKEN} para gerar um novo, com seu usuário e senha. ` };
    }
    return { status: requestResultStatusEnum_1.enumRequestResultStatus.success, message: 'Sucesso na transação' };
};
exports.validarToken = validarToken;
