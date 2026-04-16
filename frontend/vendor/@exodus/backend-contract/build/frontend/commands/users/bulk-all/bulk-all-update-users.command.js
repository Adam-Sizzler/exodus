"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BulkAllUpdateUsersCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const api_1 = require("../../../api");
const constants_2 = require("../../../constants");
const models_1 = require("../../../models");
var BulkAllUpdateUsersCommand;
(function (BulkAllUpdateUsersCommand) {
    BulkAllUpdateUsersCommand.url = api_1.REST_API.USERS.BULK.ALL.UPDATE;
    BulkAllUpdateUsersCommand.TSQ_url = BulkAllUpdateUsersCommand.url;
    BulkAllUpdateUsersCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.BULK.ALL.UPDATE, 'post', 'Bulk update all users');
    BulkAllUpdateUsersCommand.RequestSchema = zod_1.z.object({
        status: models_1.UsersSchema.shape.status.optional(),
        trafficLimitBytes: zod_1.z.optional(zod_1.z
            .number({
            invalid_type_error: 'Traffic limit must be a number',
        })
            .int('Traffic limit must be an integer')
            .min(0, 'Traffic limit must be 0 or greater')
            .describe('Traffic limit in bytes. 0 - unlimited')),
        trafficLimitStrategy: zod_1.z.optional(zod_1.z
            .nativeEnum(constants_2.RESET_PERIODS, {
            description: 'Available reset periods',
        })
            .describe('Traffic limit reset strategy')),
        expireAt: zod_1.z.optional(zod_1.z
            .string()
            .datetime({ local: true, offset: true, message: 'Invalid date format' })
            .transform((str) => new Date(str))
            .refine((date) => date > new Date(), {
            message: 'Expiration date cannot be in the past',
        })
            .describe('Expiration date: 2025-01-17T15:38:45.065Z')),
        description: zod_1.z.optional(zod_1.z.string().nullable()),
        telegramId: zod_1.z.optional(zod_1.z.number().int().nullable()),
        email: zod_1.z.optional(zod_1.z.string().email('Invalid email format').nullable()),
        tag: zod_1.z.optional(zod_1.z
            .string()
            .regex(/^[A-Z0-9_]+$/, 'Tag can only contain uppercase letters, numbers, underscores')
            .max(16, 'Tag must be less than 16 characters')
            .nullable()),
        hwidDeviceLimit: zod_1.z.optional(zod_1.z.number().int().min(0, 'Device limit must be non-negative').nullable()),
    });
    BulkAllUpdateUsersCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            eventSent: zod_1.z.boolean(),
        }),
    });
})(BulkAllUpdateUsersCommand || (exports.BulkAllUpdateUsersCommand = BulkAllUpdateUsersCommand = {}));
