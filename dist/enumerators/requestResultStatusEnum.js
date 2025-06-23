"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.enumRequestResultStatus = void 0;
var enumRequestResultStatus;
(function (enumRequestResultStatus) {
    enumRequestResultStatus[enumRequestResultStatus["success"] = 200] = "success";
    enumRequestResultStatus[enumRequestResultStatus["forbidden"] = 403] = "forbidden";
    enumRequestResultStatus[enumRequestResultStatus["restricao"] = 418] = "restricao";
    enumRequestResultStatus[enumRequestResultStatus["error"] = 413] = "error";
    enumRequestResultStatus[enumRequestResultStatus["forcarNovoLogin"] = 416] = "forcarNovoLogin";
    enumRequestResultStatus[enumRequestResultStatus["expiredToken"] = 401] = "expiredToken";
})(enumRequestResultStatus || (exports.enumRequestResultStatus = enumRequestResultStatus = {}));
