"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetConfigurationCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var GetConfigurationCommand;
(function (GetConfigurationCommand) {
    GetConfigurationCommand.url = api_1.REST_API.SYSTEM.CONFIGURATION;
    GetConfigurationCommand.TSQ_url = GetConfigurationCommand.url;
    GetConfigurationCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SYSTEM_ROUTES.CONFIGURATION, 'get', 'Get Exodus Configuration', { scope: 'configuration', kind: 'read' }, 'Returns some of the configuration values.');
    GetConfigurationCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            notifications: zod_1.z.object({
                webhook: zod_1.z.boolean().meta({
                    description: 'WEBHOOK_ENABLED',
                }),
                bandwidthUsage: zod_1.z.array(zod_1.z.number()).nullable().meta({
                    description: 'BANDWIDTH_USAGE_NOTIFICATIONS_THRESHOLD',
                }),
                notConnectedAfter: zod_1.z.array(zod_1.z.number()).nullable().meta({
                    description: 'NOT_CONNECTED_USERS_NOTIFICATIONS_AFTER_HOURS',
                }),
                expirationNotifications: zod_1.z.array(zod_1.z.number()).nullable().meta({
                    description: 'EXPIRATION_NOTIFICATIONS',
                }),
            }),
            service: zod_1.z.object({
                cleanUsageHistory: zod_1.z.boolean().meta({
                    description: 'SERVICE_CLEAN_USAGE_HISTORY',
                }),
                disableUserUsageRecords: zod_1.z.boolean().meta({
                    description: 'SERVICE_DISABLE_USER_USAGE_RECORDS',
                }),
                disableSrhRecords: zod_1.z.boolean().meta({
                    description: 'SERVICE_DISABLE_SRH_RECORDS',
                }),
                exportToRedisStream: zod_1.z.boolean().meta({
                    description: 'EXPORT_TO_STREAM_ENABLED',
                }),
            }),
            misc: zod_1.z.object({
                shortUuidLength: zod_1.z.number().meta({
                    description: 'SHORT_UUID_LENGTH',
                }),
                subPublicDomain: zod_1.z.string(),
                userUsageIgnoreBelowBytes: zod_1.z.number().meta({
                    description: 'USER_USAGE_IGNORE_BELOW_BYTES',
                }),
            }),
        }),
    });
})(GetConfigurationCommand || (exports.GetConfigurationCommand = GetConfigurationCommand = {}));
