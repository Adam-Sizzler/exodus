"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.FindAllApiTokensCommand = void 0;
const zod_1 = require("zod");
const api_tokens_schema_1 = require("../../models/api-tokens.schema");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var FindAllApiTokensCommand;
(function (FindAllApiTokensCommand) {
    FindAllApiTokensCommand.url = api_1.REST_API.API_TOKENS.GET;
    FindAllApiTokensCommand.TSQ_url = FindAllApiTokensCommand.url;
    FindAllApiTokensCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.API_TOKENS_ROUTES.GET, 'get', 'Get all API tokens', 'This endpoint is forbidden to use via "API-key". It can only be used with admin JWT-token.');
    FindAllApiTokensCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            apiKeys: zod_1.z.array(api_tokens_schema_1.ApiTokensSchema),
            docs: zod_1.z.object({
                isDocsEnabled: zod_1.z.boolean(),
                scalarPath: zod_1.z.string().nullable(),
                swaggerPath: zod_1.z.string().nullable(),
            }),
        }),
    });
})(FindAllApiTokensCommand || (exports.FindAllApiTokensCommand = FindAllApiTokensCommand = {}));
