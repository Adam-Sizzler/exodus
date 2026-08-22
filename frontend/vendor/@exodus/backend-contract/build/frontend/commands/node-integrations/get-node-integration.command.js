"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetNodeIntegrationCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var GetNodeIntegrationCommand;
(function (GetNodeIntegrationCommand) {
    GetNodeIntegrationCommand.url = api_1.REST_API.NODE_INTEGRATIONS.GET;
    GetNodeIntegrationCommand.TSQ_url = GetNodeIntegrationCommand.url(':uuid');
    GetNodeIntegrationCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODE_INTEGRATIONS_ROUTES.GET(':uuid'), 'get', 'Get Node Integration by uuid', { scope: 'get', kind: 'read' });
    GetNodeIntegrationCommand.RequestParamSchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
    });
    GetNodeIntegrationCommand.ResponseSchema = zod_1.z.object({
        response: models_1.NodeIntegrationSchema,
    });
})(GetNodeIntegrationCommand || (exports.GetNodeIntegrationCommand = GetNodeIntegrationCommand = {}));
