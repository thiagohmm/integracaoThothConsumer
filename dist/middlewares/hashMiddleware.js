"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.HashMiddleware = void 0;
require("dotenv/config");
var HashMiddleware;
(function (HashMiddleware) {
    function check(req, res, next) {
        if (req.headers.authorization === `Bearer ${process.env.HASH_FIXED}` ||
            req.originalUrl.includes('swagger')) {
            next();
        }
        else {
            res.redirect('/forbidden');
        }
    }
    HashMiddleware.check = check;
})(HashMiddleware || (exports.HashMiddleware = HashMiddleware = {}));
