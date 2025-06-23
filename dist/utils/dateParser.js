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
exports.DateParser = void 0;
var DateParser;
(function (DateParser) {
    function Local(date, local) {
        return __awaiter(this, void 0, void 0, function* () {
            if (!date)
                date = new Date().getTime();
            if (!local)
                local = 'pt-BR';
            return new Date(date).toLocaleString(local, { timeZone: 'America/Sao_Paulo' });
        });
    }
    DateParser.Local = Local;
    function Timestamp(date) {
        return __awaiter(this, void 0, void 0, function* () {
            return Math.round(new Date(date).getTime() / 1000);
        });
    }
    DateParser.Timestamp = Timestamp;
    function Now() {
        return new Date().toUTCString();
    }
    DateParser.Now = Now;
    function shortDate(date, local) {
        const ended = new Date(date);
        let month = ended.toLocaleString(local, { timeZone: 'America/Sao_Paulo', month: 'long' });
        month = month.substring(0, 3);
        month = month.charAt(0).toUpperCase() + month.slice(1);
        return `${ended.getDate()}/${month}`;
    }
    DateParser.shortDate = shortDate;
    function stringDateToDBFormat(date, time) {
        return __awaiter(this, void 0, void 0, function* () {
            const dateValue = date + ' ' + time + ':00', timestampValue = yield Timestamp(dateValue), finalDateValue = yield Local(timestampValue * 1000, 'pt-br'), dbFragments = finalDateValue.split(' '), dbFragmentsDate = dbFragments[0].split('/'), finalDBValue = dbFragmentsDate[2] + '-' + dbFragmentsDate[1] + '-' + dbFragmentsDate[0] + ' ' + dbFragments[1];
            return finalDBValue;
        });
    }
    DateParser.stringDateToDBFormat = stringDateToDBFormat;
    function ptBrStringDateToDBFormat(date, time) {
        return __awaiter(this, void 0, void 0, function* () {
            const splittedDateValue = date.split('/'), dateValue = splittedDateValue[2] + '-' + splittedDateValue[1] + '-' + splittedDateValue[0] + ' ' + time;
            return dateValue;
        });
    }
    DateParser.ptBrStringDateToDBFormat = ptBrStringDateToDBFormat;
    function dBToPtBrStringDateFormat(date) {
        return __awaiter(this, void 0, void 0, function* () {
            const splittedDateTimeValue = date.split('-'), dateValue = splittedDateTimeValue[2] + '/' + splittedDateTimeValue[1] + '/' + splittedDateTimeValue[0];
            return dateValue;
        });
    }
    DateParser.dBToPtBrStringDateFormat = dBToPtBrStringDateFormat;
    // eslint-disable-next-line @typescript-eslint/ban-types
    function changeHours(date, nHours, increase) {
        return __awaiter(this, void 0, void 0, function* () {
            const dateVal = new Date(date);
            const newDate = increase ? dateVal.setTime(dateVal.getTime() + (nHours * 60 * 60 * 1000)) : dateVal.setTime(dateVal.getTime() - (nHours * 60 * 60 * 1000));
            return newDate;
        });
    }
    DateParser.changeHours = changeHours;
})(DateParser || (exports.DateParser = DateParser = {}));
