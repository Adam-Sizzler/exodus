"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.ExodusWebhookEventSchema = exports.ExodusWebhookTorrentBlockerEvents = exports.ExodusWebhookCrmEvents = exports.ExodusWebhookErrorsEvents = exports.ExodusWebhookServiceEvents = exports.ExodusWebhookNodeEvents = exports.ExodusWebhookUserHwidDevicesEvents = exports.ExodusWebhookUserEvents = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const extended_users_schema_1 = require("../extended-users.schema");
const hwid_user_device_schema_1 = require("../hwid-user-device.schema");
const nodes_schema_1 = require("../nodes.schema");
exports.ExodusWebhookUserEvents = (zod_1.z || zod_1.default.z || zod_1.default).object({
    scope: (zod_1.z || zod_1.default.z || zod_1.default).literal(constants_1.EVENTS_SCOPES.USER),
    event: (zod_1.z || zod_1.default.z || zod_1.default).enum((0, constants_1.toZodEnum)(constants_1.EVENTS.USER)),
    timestamp: zod_1.default
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    data: extended_users_schema_1.ExtendedUsersSchema,
    meta: zod_1.default
        .object({
        notConnectedAfterHours: (zod_1.z || zod_1.default.z || zod_1.default).number().nullish(),
        expiration: (zod_1.z || zod_1.default.z || zod_1.default).number().nullish(),
    })
        .nullable(),
});
exports.ExodusWebhookUserHwidDevicesEvents = (zod_1.z || zod_1.default.z || zod_1.default).object({
    scope: (zod_1.z || zod_1.default.z || zod_1.default).literal(constants_1.EVENTS_SCOPES.USER_HWID_DEVICES),
    event: (zod_1.z || zod_1.default.z || zod_1.default).enum((0, constants_1.toZodEnum)(constants_1.EVENTS.USER_HWID_DEVICES)),
    timestamp: zod_1.default
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    data: (zod_1.z || zod_1.default.z || zod_1.default).object({
        user: extended_users_schema_1.ExtendedUsersSchema,
        hwidUserDevice: hwid_user_device_schema_1.HwidUserDeviceSchema,
    }),
});
exports.ExodusWebhookNodeEvents = (zod_1.z || zod_1.default.z || zod_1.default).object({
    scope: (zod_1.z || zod_1.default.z || zod_1.default).literal(constants_1.EVENTS_SCOPES.NODE),
    event: (zod_1.z || zod_1.default.z || zod_1.default).enum((0, constants_1.toZodEnum)(constants_1.EVENTS.NODE)),
    timestamp: zod_1.default
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    data: nodes_schema_1.NodesSchema,
});
exports.ExodusWebhookServiceEvents = (zod_1.z || zod_1.default.z || zod_1.default).object({
    scope: (zod_1.z || zod_1.default.z || zod_1.default).literal(constants_1.EVENTS_SCOPES.SERVICE),
    event: (zod_1.z || zod_1.default.z || zod_1.default).enum((0, constants_1.toZodEnum)(constants_1.EVENTS.SERVICE)),
    timestamp: zod_1.default
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    data: (zod_1.z || zod_1.default.z || zod_1.default).object({
        loginAttempt: zod_1.default
            .object({
            username: (zod_1.z || zod_1.default.z || zod_1.default).string(),
            ip: (zod_1.z || zod_1.default.z || zod_1.default).string(),
            userAgent: (zod_1.z || zod_1.default.z || zod_1.default).string(),
            description: (zod_1.z || zod_1.default.z || zod_1.default).string().optional(),
            password: (zod_1.z || zod_1.default.z || zod_1.default).string().optional(),
        })
            .optional(),
        panelVersion: (zod_1.z || zod_1.default.z || zod_1.default).string().optional(),
        subpageConfig: zod_1.default
            .object({
            action: (zod_1.z || zod_1.default.z || zod_1.default).enum((0, constants_1.toZodEnum)(constants_1.CRUD_ACTIONS)),
            uuid: (zod_1.z || zod_1.default.z || zod_1.default).uuid(),
        })
            .optional(),
        apiToken: zod_1.default
            .object({
            name: (zod_1.z || zod_1.default.z || zod_1.default).string(),
            uuid: (zod_1.z || zod_1.default.z || zod_1.default).uuid(),
            expireAt: zod_1.default
                .string()
                .datetime()
                .transform((str) => new Date(str)),
            scopes: (zod_1.z || zod_1.default.z || zod_1.default).array((zod_1.z || zod_1.default.z || zod_1.default).string()),
        })
            .optional(),
    }),
});
exports.ExodusWebhookErrorsEvents = (zod_1.z || zod_1.default.z || zod_1.default).object({
    scope: (zod_1.z || zod_1.default.z || zod_1.default).literal(constants_1.EVENTS_SCOPES.ERRORS),
    event: (zod_1.z || zod_1.default.z || zod_1.default).enum((0, constants_1.toZodEnum)(constants_1.EVENTS.ERRORS)),
    timestamp: zod_1.default
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    data: (zod_1.z || zod_1.default.z || zod_1.default).object({
        description: (zod_1.z || zod_1.default.z || zod_1.default).string(),
    }),
});
exports.ExodusWebhookCrmEvents = (zod_1.z || zod_1.default.z || zod_1.default).object({
    scope: (zod_1.z || zod_1.default.z || zod_1.default).literal(constants_1.EVENTS_SCOPES.CRM),
    event: (zod_1.z || zod_1.default.z || zod_1.default).enum((0, constants_1.toZodEnum)(constants_1.EVENTS.CRM)),
    timestamp: zod_1.default
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    data: (zod_1.z || zod_1.default.z || zod_1.default).object({
        providerName: (zod_1.z || zod_1.default.z || zod_1.default).string(),
        nodeName: (zod_1.z || zod_1.default.z || zod_1.default).string(),
        nextBillingAt: zod_1.default
            .string()
            .datetime()
            .transform((str) => new Date(str)),
        loginUrl: (zod_1.z || zod_1.default.z || zod_1.default).string(),
    }),
});
exports.ExodusWebhookTorrentBlockerEvents = (zod_1.z || zod_1.default.z || zod_1.default).object({
    scope: (zod_1.z || zod_1.default.z || zod_1.default).literal(constants_1.EVENTS_SCOPES.TORRENT_BLOCKER),
    event: (zod_1.z || zod_1.default.z || zod_1.default).enum((0, constants_1.toZodEnum)(constants_1.EVENTS.TORRENT_BLOCKER)),
    timestamp: zod_1.default
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    data: (zod_1.z || zod_1.default.z || zod_1.default).object({
        node: nodes_schema_1.NodesSchema,
        user: extended_users_schema_1.ExtendedUsersSchema,
        report: (zod_1.z || zod_1.default.z || zod_1.default).object({
            actionReport: (zod_1.z || zod_1.default.z || zod_1.default).object({
                blocked: (zod_1.z || zod_1.default.z || zod_1.default).boolean(),
                ip: (zod_1.z || zod_1.default.z || zod_1.default).string(),
                blockDuration: (zod_1.z || zod_1.default.z || zod_1.default).number(),
                willUnblockAt: zod_1.default
                    .string()
                    .datetime({ offset: true, local: true })
                    .transform((str) => new Date(str)),
                userId: (zod_1.z || zod_1.default.z || zod_1.default).string(),
                processedAt: zod_1.default
                    .string()
                    .datetime({ offset: true, local: true })
                    .transform((str) => new Date(str)),
            }),
            xrayReport: (zod_1.z || zod_1.default.z || zod_1.default).object({
                email: (zod_1.z || zod_1.default.z || zod_1.default).string().nullable(),
                level: (zod_1.z || zod_1.default.z || zod_1.default).number().nullable(),
                protocol: (zod_1.z || zod_1.default.z || zod_1.default).string().nullable(),
                network: (zod_1.z || zod_1.default.z || zod_1.default).string(),
                source: (zod_1.z || zod_1.default.z || zod_1.default).string().nullable(),
                destination: (zod_1.z || zod_1.default.z || zod_1.default).string(),
                routeTarget: (zod_1.z || zod_1.default.z || zod_1.default).string().nullable(),
                originalTarget: (zod_1.z || zod_1.default.z || zod_1.default).string().nullable(),
                inboundTag: (zod_1.z || zod_1.default.z || zod_1.default).string().nullable(),
                inboundName: (zod_1.z || zod_1.default.z || zod_1.default).string().nullable(),
                inboundLocal: (zod_1.z || zod_1.default.z || zod_1.default).string().nullable(),
                outboundTag: (zod_1.z || zod_1.default.z || zod_1.default).string().nullable(),
                ts: (zod_1.z || zod_1.default.z || zod_1.default).number(),
            }),
        }),
    }),
});
exports.ExodusWebhookEventSchema = (zod_1.z || zod_1.default.z || zod_1.default).discriminatedUnion('scope', [
    exports.ExodusWebhookUserEvents,
    exports.ExodusWebhookUserHwidDevicesEvents,
    exports.ExodusWebhookNodeEvents,
    exports.ExodusWebhookServiceEvents,
    exports.ExodusWebhookErrorsEvents,
    exports.ExodusWebhookCrmEvents,
    exports.ExodusWebhookTorrentBlockerEvents,
]);
