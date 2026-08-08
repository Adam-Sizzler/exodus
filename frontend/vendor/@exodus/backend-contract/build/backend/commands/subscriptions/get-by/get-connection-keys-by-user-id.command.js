"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetConnectionKeysByUserIdCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
var GetConnectionKeysByUserIdCommand;
(function (GetConnectionKeysByUserIdCommand) {
    GetConnectionKeysByUserIdCommand.url = api_1.REST_API.SUBSCRIPTIONS.GET_CONNECTION_KEYS_BY_USER_ID;
    GetConnectionKeysByUserIdCommand.TSQ_url = GetConnectionKeysByUserIdCommand.url(':userId');
    GetConnectionKeysByUserIdCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTIONS_ROUTES.GET_CONNECTION_KEYS_BY_USER_ID(':userId'), 'get', 'Get connection keys (base64 format) by user id', { scope: 'connection-keys', kind: 'read' });
    GetConnectionKeysByUserIdCommand.RequestParamSchema = zod_1.z.object({
        userId: models_1.numberParamSchema.describe('User ID'),
    });
    GetConnectionKeysByUserIdCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            enabledKeys: zod_1.z.array(zod_1.z.string()),
            hiddenKeys: zod_1.z.array(zod_1.z.string()),
            disabledKeys: zod_1.z.array(zod_1.z.string()),
        }),
    });

})(GetConnectionKeysByUserIdCommand || (exports.GetConnectionKeysByUserIdCommand = GetConnectionKeysByUserIdCommand = {}));
