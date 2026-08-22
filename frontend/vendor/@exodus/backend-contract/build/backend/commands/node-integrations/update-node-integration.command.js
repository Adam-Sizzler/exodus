"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UpdateNodeIntegrationCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var UpdateNodeIntegrationCommand;
(function (UpdateNodeIntegrationCommand) {
    UpdateNodeIntegrationCommand.url = api_1.REST_API.NODE_INTEGRATIONS.UPDATE;
    UpdateNodeIntegrationCommand.TSQ_url = UpdateNodeIntegrationCommand.url;
    UpdateNodeIntegrationCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODE_INTEGRATIONS_ROUTES.UPDATE, 'patch', 'Update Node Integration', { scope: 'update', kind: 'write' });
    UpdateNodeIntegrationCommand.RequestBodySchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
        name: zod_1.z
            .string()
            .min(2, 'Name must be at least 2 characters')
            .max(30, 'Name must be less than 30 characters')
            .optional(),
        description: zod_1.z.nullish(zod_1.z.string().max(255, 'Description must be less than 255 characters')),
        config: zod_1.z.record(zod_1.z.string(), zod_1.z.unknown()).optional(),
        restartNodes: zod_1.z.optional(zod_1.z.boolean()),
    });
    UpdateNodeIntegrationCommand.ResponseSchema = zod_1.z.object({
        response: models_1.NodeIntegrationSchema,
    });
})(UpdateNodeIntegrationCommand || (exports.UpdateNodeIntegrationCommand = UpdateNodeIntegrationCommand = {}));
