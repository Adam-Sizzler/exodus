"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BulkNodesUpdateCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var BulkNodesUpdateCommand;
(function (BulkNodesUpdateCommand) {
    BulkNodesUpdateCommand.url = api_1.REST_API.NODES.BULK_ACTIONS.UPDATE;
    BulkNodesUpdateCommand.TSQ_url = BulkNodesUpdateCommand.url;
    BulkNodesUpdateCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODES_ROUTES.BULK_ACTIONS.UPDATE, 'post', 'Update many nodes', { scope: 'bulk-update', kind: 'write' });
    BulkNodesUpdateCommand.RequestBodySchema = zod_1.z.object({
        uuids: zod_1.z.array(zod_1.z.uuid()).min(1),
        fields: zod_1.z.object({
            countryCode: zod_1.z.optional(zod_1.z.string().max(2).toUpperCase()),
            consumptionMultiplier: zod_1.z.optional(zod_1.z
                .number()
                .min(0.0)
                .max(100.0)
                .transform((n) => Number(n.toFixed(1)))),
            nodeConsumptionMultiplier: zod_1.z.optional(zod_1.z
                .number()
                .min(0.0)
                .max(100.0)
                .transform((n) => Number(n.toFixed(1)))),
            providerUuid: zod_1.z.uuid().nullish(),
            tags: zod_1.z.optional(zod_1.z
                .array(zod_1.z
                .string()
                .regex(/^[A-Z0-9_:]+$/, 'Tag can only contain uppercase letters, numbers, underscores and colons')
                .max(36, 'Each tag must be less than 36 characters'))
                .max(10, 'Maximum 10 tags')),
            activePluginUuid: zod_1.z.uuid().nullish(),
            integrationUuids: zod_1.z.optional(zod_1.z.array(zod_1.z.uuid()).max(20, 'Maximum 20 integrations')),
            note: zod_1.z.optional(zod_1.z.string().max(255).nullable()),
        }),
    });
})(BulkNodesUpdateCommand || (exports.BulkNodesUpdateCommand = BulkNodesUpdateCommand = {}));
