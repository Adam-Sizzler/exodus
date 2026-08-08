"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetApiTokenScopesCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var GetApiTokenScopesCommand;
(function (GetApiTokenScopesCommand) {
    GetApiTokenScopesCommand.url = api_1.REST_API.API_TOKENS.GET_SCOPES;
    GetApiTokenScopesCommand.TSQ_url = GetApiTokenScopesCommand.url;
    GetApiTokenScopesCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.API_TOKENS_ROUTES.GET_SCOPES, 'get', 'Get available API token scopes', { scope: 'list-scopes', kind: 'read' }, 'Returns the catalog of scopes that can be granted to an API token, grouped by resource. Forbidden via "API-key", admin JWT only.');
    GetApiTokenScopesCommand.EndpointScopeSchema = zod_1.z.object({
        key: zod_1.z.string(),
        kind: zod_1.z.enum(['read', 'write']),
        method: zod_1.z.string(),
        path: zod_1.z.string(),
        description: zod_1.z.string(),
    });
    GetApiTokenScopesCommand.ResourceScopesSchema = zod_1.z.object({
        resource: zod_1.z.string(),
        resourceScopes: zod_1.z.array(zod_1.z.string()),
        endpoints: zod_1.z.array(GetApiTokenScopesCommand.EndpointScopeSchema),
    });
    GetApiTokenScopesCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            wildcard: zod_1.z.string(),
            resources: zod_1.z.array(GetApiTokenScopesCommand.ResourceScopesSchema),
        }),
    });

})(GetApiTokenScopesCommand || (exports.GetApiTokenScopesCommand = GetApiTokenScopesCommand = {}));
