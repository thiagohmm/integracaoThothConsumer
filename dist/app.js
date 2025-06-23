"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
/* eslint-disable @typescript-eslint/no-var-requires */
require('dotenv').config({ path: process.env.NODE_ENV === 'test' ? '.env.test' : '.env' });
const fs_1 = __importDefault(require("fs"));
const express = require('express');
const router = express.Router();
let app = express();
app = express().use(require('@middle/express'));
app.use('/', router);
fs_1.default.readdirSync(require('path').join(__dirname, 'routes')).forEach((route) => {
    app.use('', require(`./routes/${route}`));
});
module.exports = app;
