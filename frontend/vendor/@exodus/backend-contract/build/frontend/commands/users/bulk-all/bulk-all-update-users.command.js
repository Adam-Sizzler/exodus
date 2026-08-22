"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BulkAllUpdateUsersCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const constants_2 = require("../../../constants");
var BulkAllUpdateUsersCommand;
(function (BulkAllUpdateUsersCommand) {
    BulkAllUpdateUsersCommand.url = api_1.REST_API.USERS.BULK.ALL.UPDATE;
    BulkAllUpdateUsersCommand.TSQ_url = BulkAllUpdateUsersCommand.url;
    BulkAllUpdateUsersCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.BULK.ALL.UPDATE, 'post', 'Bulk update all users', { scope: 'bulk-all-update-users', kind: 'write' });
    BulkAllUpdateUsersCommand.RequestBodySchema = zod_1.z.object({
        status: zod_1.z.enum(constants_1.USERS_STATUS).optional(),
        trafficLimitBytes: zod_1.z.optional(zod_1.z.number().min(0).describe('Traffic limit in bytes. 0 - unlimited')),
        trafficLimitStrategy: zod_1.z.optional(zod_1.z.enum(constants_2.RESET_PERIODS).describe('Available reset periods')),
        expireAt: zod_1.z.optional(zod_1.z.iso
            .datetime({ local: true, offset: true })
            .transform((str) => new Date(str))
            .refine((date) => date > new Date(), {
            error: 'Expiration date cannot be in the past',
        })
            .describe('Expiration date: 2025-01-17T15:38:45.065Z')),
        description: zod_1.z.string().nullish(),
        telegramId: zod_1.z.number().nullish(),
        email: zod_1.z.email().nullish(),
        tag: zod_1.z.optional(zod_1.z
            .string()
            .regex(/^[A-Z0-9_]+$/, 'Tag can only contain uppercase letters, numbers, underscores')
            .max(16, 'Tag must be less than 16 characters')
            .nullable()),
        hwidDeviceLimit: zod_1.z.int().min(0).nullish(),
    });
})(BulkAllUpdateUsersCommand || (exports.BulkAllUpdateUsersCommand = BulkAllUpdateUsersCommand = {}));
