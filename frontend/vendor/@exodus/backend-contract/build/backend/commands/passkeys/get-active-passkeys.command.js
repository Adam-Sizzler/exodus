"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetAllPasskeysCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var GetAllPasskeysCommand;
(function (GetAllPasskeysCommand) {
    GetAllPasskeysCommand.url = api_1.REST_API.PASSKEYS.GET_ALL_PASSKEYS;
    GetAllPasskeysCommand.TSQ_url = GetAllPasskeysCommand.url;
    GetAllPasskeysCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.PASSKEYS_ROUTES.GET_ALL_PASSKEYS, 'get', 'Get all passkeys');
    GetAllPasskeysCommand.ResponseSchema = zod_1.z.object({
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
})(GetAllPasskeysCommand || (exports.GetAllPasskeysCommand = GetAllPasskeysCommand = {}));
