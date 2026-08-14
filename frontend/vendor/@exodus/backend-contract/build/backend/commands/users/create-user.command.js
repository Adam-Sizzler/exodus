"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CreateUserCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const user_response_1 = require("./user.response");
var CreateUserCommand;
(function (CreateUserCommand) {
    CreateUserCommand.url = api_1.REST_API.USERS.CREATE;
    CreateUserCommand.TSQ_url = CreateUserCommand.url;
    CreateUserCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.CREATE, 'post', 'Create a new user', { scope: 'create', kind: 'write' });
    CreateUserCommand.RequestBodySchema = zod_1.z.object({
        username: zod_1.z
            .string()
            .regex(/^[a-zA-Z0-9_-]+$/, 'Username can only contain letters, numbers, underscores and dashes')
            .max(36, 'Username must be less than 36 characters')
            .min(3, 'Username must be at least 3 characters')
            .describe('Unique username for the user. Required. Must be 3-36 characters long and contain only letters, numbers, underscores and dashes.'),
        status: zod_1.z
            .enum(constants_1.USERS_STATUS)
            .default(constants_1.USERS_STATUS.ACTIVE)
            .optional()
            .describe('Optional. User account status. Defaults to ACTIVE.'),
        shortUuid: zod_1.z.string().optional().describe('Optional. Short UUID identifier for the user.'),
        trojanPassword: zod_1.z
            .string()
            .min(8, 'Trojan password must be at least 8 characters')
            .max(32, 'Trojan password must be less than 32 characters')
            .optional()
            .describe('Optional. Password for Trojan protocol. Must be 8-32 characters.'),
        vlessUuid: zod_1.z
            .uuid('Invalid Vless UUID format')
            .optional()
            .describe('Optional. UUID for VLESS protocol. Must be a valid UUID format.'),
        ssPassword: zod_1.z
            .string()
            .min(8, 'SS password must be at least 8 characters')
            .max(32, 'SS password must be less than 32 characters')
            .optional()
            .describe('Optional. Password for Shadowsocks protocol. Must be 8-32 characters.'),
        naivePassword: zod_1.z
            .string()
            .min(8, 'Naive password must be at least 8 characters')
            .max(32, 'Naive password must be less than 32 characters')
            .optional()
            .describe('Optional. Password for Naive protocol. Must be 8-32 characters.'),
        shadowtlsPassword: zod_1.z
            .string()
            .min(8, 'ShadowTLS password must be at least 8 characters')
            .max(32, 'ShadowTLS password must be less than 32 characters')
            .optional()
            .describe('Optional. Password for ShadowTLS protocol. Must be 8-32 characters.'),
        hysteria2Password: zod_1.z
            .string()
            .min(8, 'Hysteria2 password must be at least 8 characters')
            .max(32, 'Hysteria2 password must be less than 32 characters')
            .optional()
            .describe('Optional. Password for Hysteria2 protocol. Must be 8-32 characters.'),
        anytlsPassword: zod_1.z
            .string()
            .min(8, 'AnyTLS password must be at least 8 characters')
            .max(32, 'AnyTLS password must be less than 32 characters')
            .optional()
            .describe('Optional. Password for AnyTLS protocol. Must be 8-32 characters.'),
        trafficLimitBytes: zod_1.z
            .number()
            .min(0, 'Traffic limit must be greater than 0')
            .optional()
            .describe('Optional. Traffic limit in bytes. Set to 0 for unlimited traffic.'),
        trafficLimitStrategy: zod_1.z.optional(zod_1.z
            .enum(constants_1.RESET_PERIODS)
            .default(constants_1.RESET_PERIODS.NO_RESET)
            .describe('Available reset periods')),
        expireAt: zod_1.z.iso
            .datetime({ offset: true, local: true })
            .transform((str) => new Date(str))
            .describe('Account expiration date. Required. Format: 2025-01-17T15:38:45.065Z'),
        createdAt: zod_1.z.iso
            .datetime({ offset: true, local: true })
            .transform((str) => new Date(str))
            .optional()
            .describe('Optional. Account creation date. Format: 2025-01-17T15:38:45.065Z'),
        lastTrafficResetAt: zod_1.z.iso
            .datetime({ offset: true, local: true })
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
            .number()
            .nullish()
            .describe('Optional. Telegram user ID for notifications. Must be an integer.'),
        email: zod_1.z
            .email()
            .nullish()
            .describe('Optional. User email address. Must be a valid email format.'),
        hwidDeviceLimit: zod_1.z.optional(zod_1.z
            .int()
            .min(0)
            .describe('Optional. Maximum number of hardware devices allowed. Must be a positive integer.')),
        activeInternalSquads: zod_1.z
            .array(zod_1.z.uuid())
            .optional()
            .describe('Optional. Array of UUIDs representing enabled internal squads.'),
        externalSquadUuid: zod_1.z
            .optional(zod_1.z.nullable(zod_1.z.uuid()))
            .describe('Optional. External squad UUID.'),
    });
    CreateUserCommand.ResponseSchema = user_response_1.UserResponseSchema;
})(CreateUserCommand || (exports.CreateUserCommand = CreateUserCommand = {}));
