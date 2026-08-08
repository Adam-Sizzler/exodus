"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CreateNodeCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const node_response_1 = require("./node.response");
var CreateNodeCommand;
(function (CreateNodeCommand) {
    CreateNodeCommand.url = api_1.REST_API.NODES.CREATE;
    CreateNodeCommand.TSQ_url = CreateNodeCommand.url;
    CreateNodeCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODES_ROUTES.CREATE, 'post', 'Create a new node', { scope: 'create', kind: 'write' });
    CreateNodeCommand.RequestBodySchema = zod_1.z.object({
        name: zod_1.z.string().min(3).max(30),
        address: zod_1.z.string().min(2),
        port: zod_1.z.int().min(1).max(65535).optional(),
        proxyUrl: zod_1.z
            .string()
            .regex(/^socks5:\/\/(?:[^:@/\s]+(?::[^@/\s]*)?@)?[^:@/\s]+:\d{1,5}$/, 'Expected socks5://[user:pass@]host:port')
            .nullish(),
        isTrafficTrackingActive: zod_1.z.boolean().optional().default(false),
        trafficLimitBytes: zod_1.z.number().min(0).optional(),
        notifyPercent: zod_1.z.int().min(0).max(100).optional(),
        trafficResetDay: zod_1.z.int().min(1).max(31).optional(),
        countryCode: zod_1.z.string().max(2).toUpperCase().default('XX'),
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
        configProfile: zod_1.z.object({
            activeConfigProfileUuid: zod_1.z.uuid(),
            activeInbounds: zod_1.z.array(zod_1.z.uuid()),
        }),
        providerUuid: zod_1.z.uuid().nullish(),
        tags: zod_1.z.optional(zod_1.z
            .array(zod_1.z
            .string()
            .regex(/^[A-Z0-9_:]+$/, 'Tag can only contain uppercase letters, numbers, underscores and colons')
            .max(36, 'Each tag must be less than 36 characters'))
            .max(10, 'Maximum 10 tags')),
        activePluginUuid: zod_1.z.optional(zod_1.z.nullable(zod_1.z.uuid())),
        note: zod_1.z.optional(zod_1.z.string().max(255, 'Note must be less than 255 characters')),
    });
    CreateNodeCommand.ResponseSchema = node_response_1.NodeResponseSchema;

})(CreateNodeCommand || (exports.CreateNodeCommand = CreateNodeCommand = {}));
