"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.CerberusWebhookEventSchema = exports.CerberusWebhookCrmEvents = exports.CerberusWebhookErrorsEvents = exports.CerberusWebhookServiceEvents = exports.CerberusWebhookNodeEvents = exports.CerberusWebhookUserHwidDevicesEvents = exports.CerberusWebhookUserEvents = void 0;
const zod_1 = __importDefault(require("zod"));
const constants_1 = require("../../constants");
const hwid_user_device_schema_1 = require("../hwid-user-device.schema");
const extended_users_schema_1 = require("../extended-users.schema");
const nodes_schema_1 = require("../nodes.schema");
exports.CerberusWebhookUserEvents = zod_1.default.object({
    scope: zod_1.default.literal(constants_1.EVENTS_SCOPES.USER),
    event: zod_1.default.enum((0, constants_1.toZodEnum)(constants_1.EVENTS.USER)),
    timestamp: zod_1.default
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    data: extended_users_schema_1.ExtendedUsersSchema,
    meta: zod_1.default
        .object({
        notConnectedAfterHours: zod_1.default.number().int().nullable().optional(),
    })
        .nullable(),
});
exports.CerberusWebhookUserHwidDevicesEvents = zod_1.default.object({
    scope: zod_1.default.literal(constants_1.EVENTS_SCOPES.USER_HWID_DEVICES),
    event: zod_1.default.enum((0, constants_1.toZodEnum)(constants_1.EVENTS.USER_HWID_DEVICES)),
    timestamp: zod_1.default
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    data: zod_1.default.object({
        user: extended_users_schema_1.ExtendedUsersSchema,
        hwidUserDevice: hwid_user_device_schema_1.HwidUserDeviceSchema,
    }),
});
exports.CerberusWebhookNodeEvents = zod_1.default.object({
    scope: zod_1.default.literal(constants_1.EVENTS_SCOPES.NODE),
    event: zod_1.default.enum((0, constants_1.toZodEnum)(constants_1.EVENTS.NODE)),
    timestamp: zod_1.default
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    data: nodes_schema_1.NodesSchema,
});
exports.CerberusWebhookServiceEvents = zod_1.default.object({
    scope: zod_1.default.literal(constants_1.EVENTS_SCOPES.SERVICE),
    event: zod_1.default.enum((0, constants_1.toZodEnum)(constants_1.EVENTS.SERVICE)),
    timestamp: zod_1.default
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    data: zod_1.default.object({
        loginAttempt: zod_1.default
            .object({
            username: zod_1.default.string(),
            ip: zod_1.default.string(),
            userAgent: zod_1.default.string(),
            description: zod_1.default.string().optional(),
            password: zod_1.default.string().optional(),
        })
            .optional(),
        panelVersion: zod_1.default.string().optional(),
    }),
});
exports.CerberusWebhookErrorsEvents = zod_1.default.object({
    scope: zod_1.default.literal(constants_1.EVENTS_SCOPES.ERRORS),
    event: zod_1.default.enum((0, constants_1.toZodEnum)(constants_1.EVENTS.ERRORS)),
    timestamp: zod_1.default
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    data: zod_1.default.object({
        description: zod_1.default.string(),
    }),
});
exports.CerberusWebhookCrmEvents = zod_1.default.object({
    scope: zod_1.default.literal(constants_1.EVENTS_SCOPES.CRM),
    event: zod_1.default.enum((0, constants_1.toZodEnum)(constants_1.EVENTS.CRM)),
    timestamp: zod_1.default
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    data: zod_1.default.object({
        providerName: zod_1.default.string(),
        nodeName: zod_1.default.string(),
        nextBillingAt: zod_1.default
            .string()
            .datetime()
            .transform((str) => new Date(str)),
        loginUrl: zod_1.default.string(),
    }),
});
exports.CerberusWebhookEventSchema = zod_1.default.discriminatedUnion('scope', [
    exports.CerberusWebhookUserEvents,
    exports.CerberusWebhookUserHwidDevicesEvents,
    exports.CerberusWebhookNodeEvents,
    exports.CerberusWebhookServiceEvents,
    exports.CerberusWebhookErrorsEvents,
    exports.CerberusWebhookCrmEvents,
]);
