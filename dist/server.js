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
/* eslint-disable semi */
/* eslint-disable @typescript-eslint/no-var-requires */
require('dotenv').config({
    path: process.env.NODE_ENV === 'test' ? '.env.test' : '.env'
});
const colors = require('colors');
console.log('');
console.log(colors.yellow('Starting...'));
require("dotenv/config");
const queue_1 = require("@utils/queue");
//const app = require('./app')
//const appWs = require('./server-ws')
// const server = app.listen(process.env.PORT, () => {
//   console.log(`SGBD: ${process.env.DB_DIALECT}`)
//   console.log(`Host API: ${process.env.PORT}`)
//   console.log(colors.green('*** Integração Thoth API BD: No ar! ***'))
// })
function startListening() {
    return __awaiter(this, void 0, void 0, function* () {
        try {
            console.log('Listening to queue');
            yield (0, queue_1.listenToQueue)();
        }
        catch (error) {
            console.error('Error in listenToQueue', error);
            // Aguarde um pouco antes de tentar novamente para evitar sobrecarregar o sistema
            yield new Promise(resolve => setTimeout(resolve, 1000));
        }
    });
}
startListening();
//server.timeout = 100000 * 6
//const wss = appWs(server)
//CronManager.runAllCrons(wss)
