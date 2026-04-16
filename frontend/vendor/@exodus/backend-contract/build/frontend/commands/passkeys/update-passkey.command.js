"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UpdatePasskeyCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var UpdatePasskeyCommand;
(function (UpdatePasskeyCommand) {
    UpdatePasskeyCommand.url = api_1.REST_API.PASSKEYS.UPDATE_PASSKEY;
    UpdatePasskeyCommand.TSQ_url = UpdatePasskeyCommand.url;
    UpdatePasskeyCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.PASSKEYS_ROUTES.UPDATE_PASSKEY, 'patch', 'Update passkey');
    UpdatePasskeyCommand.RequestSchema = zod_1.z.object({
        id: zod_1.z.string(),
        name: zod_1.z
            .string()
            .min(2, 'Name must be at least 2 characters')
            .max(30, 'Name must be less than 30 characters')
            .regex(/^[A-Za-z0-9_\s-]+$/, 'Name can only contain letters, numbers, underscores, dashes and spaces'),
    });
    UpdatePasskeyCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            passkeys: zod_1.z.array(zod_1.z.object({
                id: zod_1.z.string(),
                name: zod_1.z.string(),
                createdAt: zod_1.z
                    .string()
                    .datetime({ message: 'Invalid date format', offset: true, local: true })
                    .transform((str) => new Date(str))
                    .describe('Created date. Format: 2025-01-17T15:38:45.065Z'),
                lastUsedAt: zod_1.z
                    .string()
                    .datetime({ message: 'Invalid date format', offset: true, local: true })
                    .transform((str) => new Date(str))
                    .describe('Last used date. Format: 2025-01-17T15:38:45.065Z'),
            })),
        }),
    });
})(UpdatePasskeyCommand || (exports.UpdatePasskeyCommand = UpdatePasskeyCommand = {}));
