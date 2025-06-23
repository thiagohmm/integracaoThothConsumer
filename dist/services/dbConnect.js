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
exports.db_connect = void 0;
/* eslint-disable @typescript-eslint/no-var-requires */
require('dotenv').config({
    path: process.env.NODE_ENV === 'test' ? '.env.test' : '.env'
});
const showSql = (command) => {
    if (command.indexOf('information_schema.tables') < 0 && command.indexOf('fn_tabelas_sequence') < 0) {
        while (command.indexOf('  ') >= 0) {
            command = command.split('\t').join(' ').split('\n').join(' ').split('  ').join(' ');
        }
        console.log(command, '\n');
    }
};
const cls = require('cls-hooked');
const namespace = cls.createNamespace('my-very-own-namespace');
const { Sequelize } = require('sequelize');
Sequelize.useCLS(namespace);
const libDir = `${process.env.libDir}` || '';
let wallet_location = null;
if (process.env.DB_DIALECT == 'oracledb' && libDir.length > 0) {
    const oracledb = require('oracledb');
    wallet_location = oracledb.initOracleClient({ libDir: libDir });
}
const storage = `${process.env.DB_STORAGE}`;
const dialect = `${process.env.DB_DIALECT}`;
const username = `${process.env.DB_USER}`;
const password = `${process.env.DB_PASSWD}`;
const host = `${process.env.DB_HOST}`;
const database = `${process.env.DB_DATABASE}`;
const port = parseInt(`${process.env.DB_PORT}`);
const schema = `${process.env.DB_SCHEMA}`;
const connectString = `${process.env.DB_CONNECTSTRING}`;
const node_env = `${process.env.NODE_ENV}`;
// eslint-disable-next-line @typescript-eslint/no-explicit-any
let config;
if (dialect == 'sqlite') {
    config = {
        storage: storage,
        dialect: dialect,
        logging: false
    };
}
else if (node_env == 'local') {
    config = {
        dialect: dialect,
        host: host,
        port: port,
        logging: (command) => {
            while (command.indexOf('  ') >= 0) {
                command = command.split('\n').join(' ').split('  ').join(' ');
            }
        },
        schema: schema,
        define: { 'createdAt': 'created_at', 'updatedAt': 'updated_at' }
    };
}
else {
    config = {
        quoteIdentifiers: true,
        logQueryParameters: true,
        dialect,
        logging: (command) => showSql(command),
        username,
        password,
        dialectOptions: {
            connectString,
            wallet_location,
            ssl: {
                require: true, // This will help you. But you will see nwe error
                rejectUnauthorized: false // This line will fix new error
            }
        },
        pool: {
            max: 15,
            min: 5,
            acquire: 60000,
            idle: 10000,
            evict: 20000
        }
    };
}
exports.db_connect = new Sequelize(database, username, password, config);
function authenticateDB() {
    return __awaiter(this, void 0, void 0, function* () {
        try {
            yield exports.db_connect.authenticate();
            console.log('Conexão estabelecida com sucesso.');
        }
        catch (error) {
            console.error('Não foi possível conectar ao banco de dados:', error);
        }
    });
}
authenticateDB();
