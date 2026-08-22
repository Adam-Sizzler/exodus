"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetPasskeysCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var GetPasskeysCommand;
(function (GetPasskeysCommand) {
    GetPasskeysCommand.url = api_1.REST_API.PASSKEYS.GET_ALL_PASSKEYS;
    GetPasskeysCommand.TSQ_url = GetPasskeysCommand.url;
    GetPasskeysCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.PASSKEYS_ROUTES.GET_ALL_PASSKEYS, 'get', 'Get passkeys', { scope: 'list', kind: 'read' });
    GetPasskeysCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            passkeys: zod_1.z.array(zod_1.z.object({
                id: zod_1.z.string(),
                name: zod_1.z.string(),
                createdAt: zod_1.z.iso
                    .datetime({ offset: true, local: true })
                    .transform((str) => new Date(str))
                    .describe('Created date. Format: 2025-01-17T15:38:45.065Z'),
                lastUsedAt: zod_1.z.iso
                    .datetime({ offset: true, local: true })
                    .transform((str) => new Date(str))
                    .describe('Last used date. Format: 2025-01-17T15:38:45.065Z'),
            })),
        }),
    });
})(GetPasskeysCommand || (exports.GetPasskeysCommand = GetPasskeysCommand = {}));
