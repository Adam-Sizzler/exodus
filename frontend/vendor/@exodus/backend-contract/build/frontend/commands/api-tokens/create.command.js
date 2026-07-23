"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CreateApiTokenCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const api_tokens_schema_1 = require("../../models/api-tokens.schema");
var CreateApiTokenCommand;
(function (CreateApiTokenCommand) {
    CreateApiTokenCommand.url = api_1.REST_API.API_TOKENS.CREATE;
    CreateApiTokenCommand.TSQ_url = CreateApiTokenCommand.url;
    CreateApiTokenCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.API_TOKENS_ROUTES.CREATE, 'post', 'Create a new API token', { scope: 'create', kind: 'write' }, 'This endpoint is forbidden to use via "API-key". It can only be used with an admin JWT-token.');
    CreateApiTokenCommand.RequestSchema = zod_1.z.object({
        name: zod_1.z.string().min(2).max(30),
        expiresInDays: zod_1.z.number().min(1),
        scopes: zod_1.z.array(zod_1.z.string()).optional().default(['*']),
    });
    CreateApiTokenCommand.ResponseSchema = zod_1.z.object({
        response: api_tokens_schema_1.ApiTokensSchema.extend({
            token: zod_1.z.string(),
        }),
    });
})(CreateApiTokenCommand || (exports.CreateApiTokenCommand = CreateApiTokenCommand = {}));
