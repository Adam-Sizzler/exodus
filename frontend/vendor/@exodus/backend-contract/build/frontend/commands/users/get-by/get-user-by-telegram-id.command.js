"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetUserByTelegramIdCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
var GetUserByTelegramIdCommand;
(function (GetUserByTelegramIdCommand) {
    GetUserByTelegramIdCommand.url = api_1.REST_API.USERS.GET_BY.TELEGRAM_ID;
    GetUserByTelegramIdCommand.TSQ_url = GetUserByTelegramIdCommand.url(':telegramId');
    GetUserByTelegramIdCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.GET_BY.TELEGRAM_ID(':telegramId'), 'get', 'Get users by telegram ID', { scope: 'by-telegram-id', kind: 'read' });
    GetUserByTelegramIdCommand.RequestSchema = zod_1.z.object({
        telegramId: zod_1.z.string(),
    });
    GetUserByTelegramIdCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.array(models_1.ExtendedUsersSchema),
    });
})(GetUserByTelegramIdCommand || (exports.GetUserByTelegramIdCommand = GetUserByTelegramIdCommand = {}));
