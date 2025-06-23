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
Object.defineProperty(exports, "__esModule", { value: true });
exports.listenToQueue = listenToQueue;
const amqplib_1 = require("amqplib");
const dotenv_1 = __importDefault(require("dotenv"));
const stageCompraService_1 = require("@domains/stageCompra/services/stageCompraService");
const stageVendaService_1 = require("@domains/stageVenda/services/stageVendaService");
const stageCompraSchema_1 = require("@schemas/stageCompraSchema");
const stageVendaSchema_1 = require("@schemas/stageVendaSchema");
const stageEstoqueSchema_1 = require("@schemas/stageEstoqueSchema");
const stageEstoqueService_1 = require("@domains/stageEstoque/services/stageEstoqueService");
const stageLogErrorSchema_1 = require("@schemas/stageLogErrorSchema");
dotenv_1.default.config();
let connection = null;
let channel = null;
const queue = 'thothQueue';
let reconnecting = false;
function connectToRabbitMQ() {
    return __awaiter(this, void 0, void 0, function* () {
        const rabbitmqUrl = process.env.ENV_RABBITMQ;
        if (!rabbitmqUrl) {
            throw new Error('RabbitMQ URL is not defined.');
        }
        while (!connection || !channel) {
            try {
                console.log('Attempting to connect to RabbitMQ...');
                connection = yield (0, amqplib_1.connect)(rabbitmqUrl);
                console.log('Successfully connected to RabbitMQ');
                connection.on('close', handleReconnect);
                connection.on('error', (err) => {
                    console.error('Connection error:', err.message);
                    handleReconnect();
                });
                channel = yield connection.createChannel();
                console.log('Channel created');
                channel.on('close', () => {
                    console.warn('Channel closed. Triggering reconnection...');
                    handleReconnect();
                });
                channel.on('error', (err) => {
                    console.error('Channel error:', err.message);
                    handleReconnect();
                });
                yield assertQueue(queue);
                startConsuming();
            }
            catch (error) {
                console.error('Failed to connect to RabbitMQ:', error, 'Retrying in 5 seconds...');
                yield new Promise(resolve => setTimeout(resolve, 5000));
            }
        }
    });
}
const getRabbitMQChannel = () => __awaiter(void 0, void 0, void 0, function* () {
    if (!connection || !channel) {
        yield connectToRabbitMQ();
    }
    return channel;
});
const assertQueue = (queue) => __awaiter(void 0, void 0, void 0, function* () {
    const channel = yield getRabbitMQChannel();
    yield channel.assertQueue(queue, { durable: true });
    channel.prefetch(1);
});
function handleReconnect() {
    return __awaiter(this, void 0, void 0, function* () {
        if (reconnecting)
            return;
        reconnecting = true;
        console.warn('Reconnecting to RabbitMQ...');
        connection = null;
        channel = null;
        yield new Promise(resolve => setTimeout(resolve, 5000));
        reconnecting = false;
        yield connectToRabbitMQ();
    });
}
function startConsuming() {
    return __awaiter(this, void 0, void 0, function* () {
        const channel = yield getRabbitMQChannel();
        channel.consume(queue, (msg) => __awaiter(this, void 0, void 0, function* () {
            if (msg) {
                try {
                    const messageContent = JSON.parse(msg.content.toString());
                    if (messageContent.processa === 'compra') {
                        console.log('Processing compra');
                        const stageCompraQuery = new stageCompraSchema_1.StageCompraQuery();
                        yield (0, stageCompraService_1.PersisteDadosComprasThoth)(messageContent.dados, stageCompraQuery);
                    }
                    else if (messageContent.processa === 'venda') {
                        console.log('Processing venda');
                        const stageVendaQuery = new stageVendaSchema_1.StageVendaQuery();
                        yield (0, stageVendaService_1.PersisteDadosVendasThoth)(messageContent.dados, stageVendaQuery);
                    }
                    else if (messageContent.processa === 'estoque') {
                        console.log('Processing estoque');
                        const stageEstoqueQuery = new stageEstoqueSchema_1.StageEstoqueQuery();
                        yield (0, stageEstoqueService_1.PersisteDadosEstoqueThoth)(messageContent.dados, stageEstoqueQuery);
                    }
                    const currentChannel = yield getRabbitMQChannel();
                    currentChannel.ack(msg); // Always use the current channel
                }
                catch (err) {
                    const stageLogErrorQuery = new stageLogErrorSchema_1.StageLogErrorQuery();
                    const dadosErro = {
                        DATA: 0,
                        CD_IBM_LOJA: '',
                        MENSAGEM_ERRO: err.message,
                        JSON: JSON.stringify(msg),
                        INTERFACE: 'undefined',
                        DATA_ERRO: new Date(),
                    };
                    console.error('Error processing message:', dadosErro);
                    yield stageLogErrorQuery.save(dadosErro);
                    const currentChannel = yield getRabbitMQChannel();
                    currentChannel.nack(msg, false, false); // Reject message
                }
            }
        }), { noAck: false });
    });
}
function listenToQueue() {
    return __awaiter(this, void 0, void 0, function* () {
        try {
            yield connectToRabbitMQ();
            console.log(`Listening for messages in queue: ${queue}`);
        }
        catch (error) {
            console.error('Error setting up RabbitMQ consumer:', error);
        }
    });
}
// Start RabbitMQ connection
connectToRabbitMQ();
