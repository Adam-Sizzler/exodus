"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UpdateUserCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const user_response_1 = require("./user.response");
var UpdateUserCommand;
(function (UpdateUserCommand) {
    UpdateUserCommand.url = api_1.REST_API.USERS.UPDATE;
    UpdateUserCommand.TSQ_url = UpdateUserCommand.url;
    UpdateUserCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.UPDATE, 'patch', 'Update a user', { scope: 'update', kind: 'write' }, 'Update a user by ID or username. Exactly one of the fields must be provided.');
    UpdateUserCommand.RequestBodySchema = zod_1.z
        .object({
        username: zod_1.z.optional(zod_1.z.string().describe('Username of the user')),
        id: zod_1.z.optional(zod_1.z.number().describe('ID of the user')),
        status: zod_1.z.enum([constants_1.USERS_STATUS.ACTIVE, constants_1.USERS_STATUS.DISABLED]).optional(),
        trafficLimitBytes: zod_1.z
            .number()
            .min(0)
            .describe('Traffic limit in bytes. 0 - unlimited')
            .optional(),
        trafficLimitStrategy: zod_1.z
            .enum(constants_1.RESET_PERIODS)
            .describe('Traffic limit reset strategy')
            .optional(),
        expireAt: zod_1.z.iso
            .datetime({ local: true, offset: true })
            .transform((str) => new Date(str))
            .refine((date) => date > new Date(), {
            error: 'Expiration date cannot be in the past',
        })
            .describe('Expiration date: 2025-01-17T15:38:45.065Z')
            .optional(),
        description: zod_1.z.optional(zod_1.z.string().nullable()),
        tag: zod_1.z.optional(zod_1.z
            .string()
            .regex(/^[A-Z0-9_]+$/, 'Tag can only contain uppercase letters, numbers, underscores')
            .max(16, 'Tag must be less than 16 characters')
            .nullable()),
        telegramId: zod_1.z.number().nullish(),
        email: zod_1.z.email().nullish(),
        hwidDeviceLimit: zod_1.z.int().min(0).nullish(),
        activeInternalSquads: zod_1.z.array(zod_1.z.uuid()).optional(),
        externalSquadUuid: zod_1.z
            .optional(zod_1.z.nullable(zod_1.z.uuid()))
            .describe('Optional. External squad UUID.'),
    })
        .refine((d) => d.username ?? d.id, {
        error: 'At least one of username, id must be provided',
    });
    UpdateUserCommand.ResponseSchema = user_response_1.UserResponseSchema;

})(UpdateUserCommand || (exports.UpdateUserCommand = UpdateUserCommand = {}));
