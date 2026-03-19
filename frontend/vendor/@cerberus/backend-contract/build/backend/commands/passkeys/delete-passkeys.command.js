"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DeletePasskeyCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var DeletePasskeyCommand;
(function (DeletePasskeyCommand) {
    DeletePasskeyCommand.url = api_1.REST_API.PASSKEYS.DELETE_PASSKEY;
    DeletePasskeyCommand.TSQ_url = DeletePasskeyCommand.url;
    DeletePasskeyCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.PASSKEYS_ROUTES.DELETE_PASSKEY, 'delete', 'Delete a passkey by ID');
    DeletePasskeyCommand.RequestSchema = zod_1.z.object({
        id: zod_1.z.string(),
    });
    DeletePasskeyCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            passkeys: zod_1.z.array(zod_1.z.object({
                id: zod_1.z.string(),
                name: zod_1.z.string(),
                createdAt: zod_1.z
                    .string({
                    invalid_type_error: 'Invalid date format',
                })
                    .datetime({ message: 'Invalid date format', offset: true, local: true })
                    .transform((str) => new Date(str))
                    .describe('Created date. Format: 2025-01-17T15:38:45.065Z'),
                lastUsedAt: zod_1.z
                    .string({
                    invalid_type_error: 'Invalid date format',
                })
                    .datetime({ message: 'Invalid date format', offset: true, local: true })
                    .transform((str) => new Date(str))
                    .describe('Last used date. Format: 2025-01-17T15:38:45.065Z'),
            })),
        }),
    });
})(DeletePasskeyCommand || (exports.DeletePasskeyCommand = DeletePasskeyCommand = {}));
