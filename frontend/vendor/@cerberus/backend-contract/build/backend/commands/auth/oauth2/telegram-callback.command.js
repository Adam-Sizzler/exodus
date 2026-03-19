"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.TelegramCallbackCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const api_1 = require("../../../api");
var TelegramCallbackCommand;
(function (TelegramCallbackCommand) {
    TelegramCallbackCommand.url = api_1.REST_API.AUTH.OAUTH2.TELEGRAM_CALLBACK;
    TelegramCallbackCommand.TSQ_url = TelegramCallbackCommand.url;
    TelegramCallbackCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.AUTH_ROUTES.OAUTH2.TELEGRAM_CALLBACK, 'post', 'Callback from Telegram OAuth2');
    TelegramCallbackCommand.RequestSchema = zod_1.z.object({
        id: zod_1.z.number(),
        first_name: zod_1.z.string(),
        last_name: zod_1.z.string().optional(),
        username: zod_1.z.string().optional(),
        photo_url: zod_1.z.string().optional(),
        auth_date: zod_1.z.number(),
        hash: zod_1.z.string(),
    });
    TelegramCallbackCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            accessToken: zod_1.z.string(),
        }),
    });
})(TelegramCallbackCommand || (exports.TelegramCallbackCommand = TelegramCallbackCommand = {}));
