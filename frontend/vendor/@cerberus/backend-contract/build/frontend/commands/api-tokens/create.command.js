"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CreateApiTokenCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var CreateApiTokenCommand;
(function (CreateApiTokenCommand) {
    CreateApiTokenCommand.url = api_1.REST_API.API_TOKENS.CREATE;
    CreateApiTokenCommand.TSQ_url = CreateApiTokenCommand.url;
    CreateApiTokenCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.API_TOKENS_ROUTES.CREATE, 'post', 'Create a new API token', 'This endpoint is forbidden to use via "API-key". It can only be used with an admin JWT-token.');
    CreateApiTokenCommand.RequestSchema = zod_1.z.object({
        tokenName: zod_1.z.string(),
    });
    CreateApiTokenCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            token: zod_1.z.string(),
            uuid: zod_1.z.string(),
        }),
    });
})(CreateApiTokenCommand || (exports.CreateApiTokenCommand = CreateApiTokenCommand = {}));
