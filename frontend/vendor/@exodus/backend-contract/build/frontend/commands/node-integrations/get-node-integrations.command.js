"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetNodeIntegrationsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var GetNodeIntegrationsCommand;
(function (GetNodeIntegrationsCommand) {
    GetNodeIntegrationsCommand.url = api_1.REST_API.NODE_INTEGRATIONS.GET_ALL;
    GetNodeIntegrationsCommand.TSQ_url = GetNodeIntegrationsCommand.url;
    GetNodeIntegrationsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODE_INTEGRATIONS_ROUTES.GET_ALL, 'get', 'Get all Node Integrations', { scope: 'list', kind: 'read' });
    GetNodeIntegrationsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            total: zod_1.z.number(),
            nodeIntegrations: zod_1.z.array(models_1.NodeIntegrationSchema),
        }),
    });
})(GetNodeIntegrationsCommand || (exports.GetNodeIntegrationsCommand = GetNodeIntegrationsCommand = {}));
