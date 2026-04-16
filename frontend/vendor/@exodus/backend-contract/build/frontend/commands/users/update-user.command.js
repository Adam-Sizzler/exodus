"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UpdateUserCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
const api_1 = require("../../api");
var UpdateUserCommand;
(function (UpdateUserCommand) {
    UpdateUserCommand.url = api_1.REST_API.USERS.UPDATE;
    UpdateUserCommand.TSQ_url = UpdateUserCommand.url;
    UpdateUserCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.UPDATE, 'patch', 'Update a user by UUID or username');
    UpdateUserCommand.RequestSchema = zod_1.z
        .object({
        username: zod_1.z.optional(zod_1.z.string().describe('Username of the user')),
        uuid: zod_1.z.optional(zod_1.z
            .string()
            .uuid()
            .describe('UUID of the user. UUID has higher priority than username, so if both are provided, username will be ignored.')),
        status: zod_1.z
            .enum([constants_1.USERS_STATUS.ACTIVE, constants_1.USERS_STATUS.DISABLED], {
            errorMap: () => ({
                message: "You can't change status to LIMITED or EXPIRED. These statuses handled by Exodus.",
            }),
        })
            .optional(),
        trafficLimitBytes: zod_1.z
            .number({
            invalid_type_error: 'Traffic limit must be a number',
        })
            .int('Traffic limit must be an integer')
            .min(0, 'Traffic limit must be greater than 0')
            .describe('Traffic limit in bytes. 0 - unlimited')
            .optional(),
        trafficLimitStrategy: models_1.UsersSchema.shape.trafficLimitStrategy
            .describe('Traffic limit reset strategy')
            .superRefine((val, ctx) => {
            if (val && !Object.values(constants_1.RESET_PERIODS).includes(val)) {
                ctx.addIssue({
                    code: zod_1.z.ZodIssueCode.invalid_enum_value,
                    message: 'Invalid traffic limit strategy',
                    path: ['trafficLimitStrategy'],
                    received: val,
                    options: Object.values(constants_1.RESET_PERIODS),
                });
            }
        })
            .optional(),
        expireAt: zod_1.z
            .string()
            .datetime({ local: true, offset: true, message: 'Invalid date format' })
            .transform((str) => new Date(str))
            .refine((date) => date > new Date(), {
            message: 'Expiration date cannot be in the past',
        })
            .describe('Expiration date: 2025-01-17T15:38:45.065Z')
            .optional(),
        description: zod_1.z.optional(zod_1.z.string().nullable()),
        tag: zod_1.z.optional(zod_1.z
            .string()
            .regex(/^[A-Z0-9_]+$/, 'Tag can only contain uppercase letters, numbers, underscores')
            .max(16, 'Tag must be less than 16 characters')
            .nullable()),
        telegramId: zod_1.z.optional(zod_1.z.number().int().nullable()),
        email: zod_1.z.optional(zod_1.z.string().email('Invalid email format').nullable()),
        hwidDeviceLimit: zod_1.z.optional(zod_1.z.number().int().min(0, 'Device limit must be non-negative').nullable()),
        activeInternalSquads: zod_1.z
            .array(zod_1.z.string().uuid(), {
            invalid_type_error: 'Enabled internal squads must be an array of UUIDs',
        })
            .optional(),
        externalSquadUuid: zod_1.z
            .optional(zod_1.z.nullable(zod_1.z.string().uuid()))
            .describe('Optional. External squad UUID.'),
    })
        .refine((data) => data.uuid || data.username, {
        message: 'Either uuid or username must be provided',
    });
    UpdateUserCommand.ResponseSchema = zod_1.z.object({
        response: models_1.ExtendedUsersSchema,
    });
})(UpdateUserCommand || (exports.UpdateUserCommand = UpdateUserCommand = {}));
