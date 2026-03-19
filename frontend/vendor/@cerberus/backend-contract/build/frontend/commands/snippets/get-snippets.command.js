"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetSnippetsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var GetSnippetsCommand;
(function (GetSnippetsCommand) {
    GetSnippetsCommand.url = api_1.REST_API.SNIPPETS.GET;
    GetSnippetsCommand.TSQ_url = GetSnippetsCommand.url;
    GetSnippetsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SNIPPETS_ROUTES.GET, 'get', 'Get snippets');
    GetSnippetsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            total: zod_1.z.number(),
            snippets: zod_1.z.array(models_1.SnippetsSchema),
        }),
    });
})(GetSnippetsCommand || (exports.GetSnippetsCommand = GetSnippetsCommand = {}));
