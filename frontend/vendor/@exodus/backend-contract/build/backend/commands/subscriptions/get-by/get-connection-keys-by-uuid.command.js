"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetConnectionKeysByUuidCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var GetConnectionKeysByUuidCommand;
(function (GetConnectionKeysByUuidCommand) {
    GetConnectionKeysByUuidCommand.url = api_1.REST_API.SUBSCRIPTIONS.GET_CONNECTION_KEYS_BY_UUID;
    GetConnectionKeysByUuidCommand.TSQ_url = GetConnectionKeysByUuidCommand.url(':uuid');
    GetConnectionKeysByUuidCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTIONS_ROUTES.GET_CONNECTION_KEYS_BY_UUID(':uuid'), 'get', 'Get connection keys (base64 format) by uuid', { scope: 'connection-keys', kind: 'read' });
    GetConnectionKeysByUuidCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string(),
    });
    GetConnectionKeysByUuidCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            enabledKeys: zod_1.z.array(zod_1.z.string()),
            hiddenKeys: zod_1.z.array(zod_1.z.string()),
            disabledKeys: zod_1.z.array(zod_1.z.string()),
        }),
    });
})(GetConnectionKeysByUuidCommand || (exports.GetConnectionKeysByUuidCommand = GetConnectionKeysByUuidCommand = {}));
