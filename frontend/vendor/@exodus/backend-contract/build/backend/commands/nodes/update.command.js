"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UpdateNodeCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
const node_response_1 = require("./node.response");
var UpdateNodeCommand;
(function (UpdateNodeCommand) {
    UpdateNodeCommand.url = api_1.REST_API.NODES.UPDATE;
    UpdateNodeCommand.TSQ_url = UpdateNodeCommand.url;
    UpdateNodeCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODES_ROUTES.UPDATE, 'patch', 'Update node', {
        scope: 'update',
        kind: 'write',
    });
    UpdateNodeCommand.RequestBodySchema = models_1.NodesSchema.pick({
        uuid: true,
    }).extend({
        name: zod_1.z.optional(zod_1.z.string().min(3).max(30)),
        address: zod_1.z.optional(zod_1.z.string().min(2)),
        port: zod_1.z.optional(zod_1.z.number().min(1).max(65535)),
        proxyUrl: zod_1.z
            .string()
            .regex(/^socks5:\/\/(?:[^:@/\s]+(?::[^@/\s]*)?@)?[^:@/\s]+:\d{1,5}$/, 'Expected socks5://[user:pass@]host:port')
            .nullish(),
        isTrafficTrackingActive: zod_1.z.optional(zod_1.z.boolean()),
        trafficLimitBytes: zod_1.z.optional(zod_1.z.number().min(0)),
        notifyPercent: zod_1.z.optional(zod_1.z.number().min(0).max(100)),
        trafficResetDay: zod_1.z.optional(zod_1.z.number().min(1).max(31)),
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
        configProfile: zod_1.z
            .object({
            activeConfigProfileUuid: zod_1.z.uuid(),
            activeInbounds: zod_1.z.array(zod_1.z.uuid()),
        })
            .optional(),
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
        ips: zod_1.z.optional(models_1.NodeIpsSchema),
    });
    UpdateNodeCommand.ResponseSchema = node_response_1.NodeResponseSchema;
})(UpdateNodeCommand || (exports.UpdateNodeCommand = UpdateNodeCommand = {}));
