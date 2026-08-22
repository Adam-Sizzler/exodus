"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CreateNodeIntegrationCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var CreateNodeIntegrationCommand;
(function (CreateNodeIntegrationCommand) {
    CreateNodeIntegrationCommand.url = api_1.REST_API.NODE_INTEGRATIONS.CREATE;
    CreateNodeIntegrationCommand.TSQ_url = CreateNodeIntegrationCommand.url;
    CreateNodeIntegrationCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODE_INTEGRATIONS_ROUTES.CREATE, 'post', 'Create Node Integration', { scope: 'create', kind: 'write' });
    CreateNodeIntegrationCommand.RequestBodySchema = zod_1.z.object({
        name: zod_1.z
            .string()
            .min(2, 'Name must be at least 2 characters')
            .max(30, 'Name must be less than 30 characters'),
        description: zod_1.z.nullish(zod_1.z.string().max(255, 'Description must be less than 255 characters')),
        config: zod_1.z.record(zod_1.z.string(), zod_1.z.unknown()),
    });
    CreateNodeIntegrationCommand.ResponseSchema = zod_1.z.object({
        response: models_1.NodeIntegrationSchema,
    });
})(CreateNodeIntegrationCommand || (exports.CreateNodeIntegrationCommand = CreateNodeIntegrationCommand = {}));
