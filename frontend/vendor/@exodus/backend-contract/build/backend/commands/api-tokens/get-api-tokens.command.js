"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetApiTokensCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const api_tokens_schema_1 = require("../../models/api-tokens.schema");
var GetApiTokensCommand;
(function (GetApiTokensCommand) {
    GetApiTokensCommand.url = api_1.REST_API.API_TOKENS.GET;
    GetApiTokensCommand.TSQ_url = GetApiTokensCommand.url;
    GetApiTokensCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.API_TOKENS_ROUTES.GET, 'get', 'Get all API tokens', { scope: 'list', kind: 'read' }, 'This endpoint is forbidden to use via "API-key". It can only be used with admin JWT-token.');
    GetApiTokensCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            tokens: zod_1.z.array(api_tokens_schema_1.ApiTokensSchema),
        }),
    });

})(GetApiTokensCommand || (exports.GetApiTokensCommand = GetApiTokensCommand = {}));
