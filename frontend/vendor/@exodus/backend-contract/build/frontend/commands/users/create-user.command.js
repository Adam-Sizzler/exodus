"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CreateUserCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
const api_1 = require("../../api");
var CreateUserCommand;
(function (CreateUserCommand) {
    CreateUserCommand.url = api_1.REST_API.USERS.CREATE;
    CreateUserCommand.TSQ_url = CreateUserCommand.url;
    CreateUserCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.CREATE, 'post', 'Create a new user');
    CreateUserCommand.RequestSchema = zod_1.z.object({
        username: zod_1.z
            .string({
            required_error: 'Username is required',
            invalid_type_error: 'Username must be a string',
        })
            .regex(/^[a-zA-Z0-9_-]+$/, 'Username can only contain letters, numbers, underscores and dashes')
            .max(36, 'Username must be less than 36 characters')
            .min(3, 'Username must be at least 3 characters')
            .describe('Unique username for the user. Required. Must be 3-36 characters long and contain only letters, numbers, underscores and dashes.'),
        status: models_1.UsersSchema.shape.status
            .optional()
            .default(constants_1.USERS_STATUS.ACTIVE)
            .describe('Optional. User account status. Defaults to ACTIVE.'),
        shortUuid: zod_1.z
            .string({
            invalid_type_error: 'Short UUID must be a string',
        })
            .optional()
            .describe('Optional. Short UUID identifier for the user.'),
        trojanPassword: zod_1.z
            .string({
            invalid_type_error: 'Trojan password must be a string',
        })
            .min(8, 'Trojan password must be at least 8 characters')
            .max(32, 'Trojan password must be less than 32 characters')
            .optional()
            .describe('Optional. Password for Trojan protocol. Must be 8-32 characters.'),
        vlessUuid: zod_1.z
            .string({
            invalid_type_error: 'Vless UUID must be a string',
        })
            .uuid('Invalid Vless UUID format')
            .optional()
            .describe('Optional. UUID for VLESS protocol. Must be a valid UUID format.'),
        ssPassword: zod_1.z
            .string({
            invalid_type_error: 'SS password must be a string',
        })
            .min(8, 'SS password must be at least 8 characters')
            .max(32, 'SS password must be less than 32 characters')
            .optional()
            .describe('Optional. Password for Shadowsocks protocol. Must be 8-32 characters.'),
        trafficLimitBytes: zod_1.z
            .number({
            invalid_type_error: 'Traffic limit must be a number',
        })
            .int('Traffic limit must be an integer')
            .min(0, 'Traffic limit must be greater than 0')
            .optional()
            .describe('Optional. Traffic limit in bytes. Set to 0 for unlimited traffic.'),
        trafficLimitStrategy: zod_1.z.optional(models_1.UsersSchema.shape.trafficLimitStrategy
            .default(constants_1.RESET_PERIODS.NO_RESET)
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
        })),
        expireAt: zod_1.z
            .string({
            required_error: 'Expiration date is required',
            invalid_type_error: 'Invalid date format',
        })
            .datetime({ message: 'Invalid date format', offset: true, local: true })
            .transform((str) => new Date(str))
            .describe('Account expiration date. Required. Format: 2025-01-17T15:38:45.065Z'),
        createdAt: zod_1.z
            .string({
            invalid_type_error: 'Invalid date format',
        })
            .datetime({ message: 'Invalid date format', offset: true, local: true })
            .transform((str) => new Date(str))
            .optional()
            .describe('Optional. Account creation date. Format: 2025-01-17T15:38:45.065Z'),
        lastTrafficResetAt: zod_1.z
            .string({
            invalid_type_error: 'Invalid date format',
        })
            .datetime({ message: 'Invalid date format', offset: true, local: true })
            .transform((str) => new Date(str))
            .optional()
            .describe('Optional. Date of last traffic reset. Format: 2025-01-17T15:38:45.065Z'),
        description: zod_1.z
            .string()
            .optional()
            .describe('Optional. Additional notes or description for the user account.'),
        tag: zod_1.z
            .optional(zod_1.z
            .string()
            .regex(/^[A-Z0-9_]+$/, 'Tag can only contain uppercase letters, numbers, underscores')
            .max(16, 'Tag must be less than 16 characters')
            .nullable())
            .describe('Optional. User tag for categorization. Max 16 characters, uppercase letters, numbers and underscores only.'),
        telegramId: zod_1.z
            .optional(zod_1.z.number().int().nullable())
            .describe('Optional. Telegram user ID for notifications. Must be an integer.'),
        email: zod_1.z
            .optional(zod_1.z.string().email('Invalid email format').nullable())
            .describe('Optional. User email address. Must be a valid email format.'),
        hwidDeviceLimit: zod_1.z.optional(zod_1.z
            .number({ invalid_type_error: 'Device limit must be a number' })
            .int('Device limit must be an integer')
            .min(0, 'Device limit must be greater than 0')
            .describe('Optional. Maximum number of hardware devices allowed. Must be a positive integer.')),
        activeInternalSquads: zod_1.z
            .array(zod_1.z.string().uuid(), {
            invalid_type_error: 'Enabled internal squads must be an array',
        })
            .optional()
            .describe('Optional. Array of UUIDs representing enabled internal squads.'),
        uuid: zod_1.z.optional(zod_1.z
            .string()
            .uuid()
            .describe('Optional. Pass UUID to create user with specific UUID, otherwise it will be generated automatically.')),
        externalSquadUuid: zod_1.z
            .optional(zod_1.z.nullable(zod_1.z.string().uuid()))
            .describe('Optional. External squad UUID.'),
    });
    CreateUserCommand.ResponseSchema = zod_1.z.object({
        response: models_1.ExtendedUsersSchema,
    });
})(CreateUserCommand || (exports.CreateUserCommand = CreateUserCommand = {}));
