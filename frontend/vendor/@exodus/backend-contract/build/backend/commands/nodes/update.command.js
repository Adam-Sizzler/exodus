"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UpdateNodeCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var UpdateNodeCommand;
(function (UpdateNodeCommand) {
    UpdateNodeCommand.url = api_1.REST_API.NODES.UPDATE;
    UpdateNodeCommand.TSQ_url = UpdateNodeCommand.url;
    UpdateNodeCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODES_ROUTES.UPDATE, 'patch', 'Update node', {
        scope: 'update',
        kind: 'write',
    });
    UpdateNodeCommand.RequestSchema = models_1.NodesSchema.pick({
        uuid: true,
    }).extend({
        name: zod_1.z.optional(zod_1.z.string().min(3, 'Min. 3 characters').max(30, 'Max. 30 characters')),
        address: zod_1.z.optional(zod_1.z.string().min(2, 'Min. 2 characters')),
        port: zod_1.z.optional(zod_1.z
            .number()
            .min(1, 'Port must be greater than 0')
            .max(65535, 'Port must be less than 65535')),
        proxyUrl: zod_1.z
            .string()
            .regex(/^socks5:\/\/(?:[^:@/\s]+(?::[^@/\s]*)?@)?[^:@/\s]+:\d{1,5}$/, 'Expected socks5://[user:pass@]host:port')
            .nullable()
            .optional(),
        apiSchema: zod_1.z.optional(zod_1.z.enum(['mtls', 'tls'])),
        apiPath: zod_1.z.optional(zod_1.z.string().max(255, 'API path must be less than 255 characters')),
        grpcAuthToken: zod_1.z.optional(zod_1.z.string().max(255, 'gRPC auth token must be less than 255 characters')),
        isTrafficTrackingActive: zod_1.z.optional(zod_1.z.boolean()),
        trafficLimitBytes: zod_1.z.optional(zod_1.z.number().min(0, 'Traffic limit must be greater than 0')),
        notifyPercent: zod_1.z.optional(zod_1.z
            .number()
            .min(0, 'Notify percent must be greater than 0')
            .max(100, 'Notify percent must be less than 100')),
        trafficResetDay: zod_1.z.optional(zod_1.z
            .number()
            .min(1, 'Traffic reset day must be greater than 0')
            .max(31, 'Traffic reset day must be less than 31')),
        countryCode: zod_1.z.optional(zod_1.z.string().max(2, 'Country code must be 2 characters').toUpperCase()),
        consumptionMultiplier: zod_1.z.optional(zod_1.z
            .number()
            .min(0.0, 'Consumption multiplier must be greater than 0.0')
            .max(100.0, 'Consumption multiplier must be less than 100.0')
            .transform((n) => Number(n.toFixed(1)))),
        nodeConsumptionMultiplier: zod_1.z.optional(zod_1.z
            .number()
            .min(0.0, 'Node consumption multiplier must be greater than 0.0')
            .max(100.0, 'Node consumption multiplier must be less than 100.0')
            .transform((n) => Number(n.toFixed(1)))),
        configProfile: zod_1.z
            .object({
            activeConfigProfileUuid: zod_1.z.string().uuid(),
            activeInbounds: zod_1.z.array(zod_1.z.string().uuid(), {
                invalid_type_error: 'Must be an array of UUIDs',
            }),
        })
            .optional(),
        providerUuid: zod_1.z.optional(zod_1.z.nullable(zod_1.z.string().uuid())),
        tags: zod_1.z.optional(zod_1.z
            .array(zod_1.z
            .string()
            .regex(/^[A-Z0-9_:]+$/, 'Tag can only contain uppercase letters, numbers, underscores and colons')
            .max(36, 'Each tag must be less than 36 characters'))
            .max(10, 'Maximum 10 tags')),
        activePluginUuid: zod_1.z.optional(zod_1.z.nullable(zod_1.z.string().uuid())),
        note: zod_1.z.optional(zod_1.z.string().max(255, 'Note must be less than 255 characters').nullable()),
    });
    UpdateNodeCommand.ResponseSchema = zod_1.z.object({
        response: models_1.NodesSchema,
    });
})(UpdateNodeCommand || (exports.UpdateNodeCommand = UpdateNodeCommand = {}));
