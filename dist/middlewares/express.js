"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
require("dotenv/config");
const express_1 = __importDefault(require("express"));
const method_override_1 = __importDefault(require("method-override"));
const cors_1 = __importDefault(require("cors"));
const express_fileupload_1 = __importDefault(require("express-fileupload"));
module.exports = [
    express_1.default.urlencoded({ limit: '20MB', extended: true }),
    express_1.default.json({ limit: '20MB' }),
    express_1.default.raw({ limit: '20MB' }),
    express_1.default.text({ limit: '20MB' }),
    (0, express_fileupload_1.default)(),
    (0, method_override_1.default)('_method'),
    (0, cors_1.default)(),
];
