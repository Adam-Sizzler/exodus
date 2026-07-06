"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetUsersStreamCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var GetUsersStreamCommand;
(function (GetUsersStreamCommand) {
    GetUsersStreamCommand.url = api_1.REST_API.USERS.STREAM;
    GetUsersStreamCommand.TSQ_url = GetUsersStreamCommand.url;
    GetUsersStreamCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.STREAM, 'get', 'Get all users using cursor-based (keyset) pagination', { scope: 'stream', kind: 'read' });
    GetUsersStreamCommand.RequestQuerySchema = zod_1.z.object({
        cursor: zod_1.z
            .string()
            .regex(/^\d+$/, 'Cursor must be a positive integer string')
            .optional()
            .describe('Cursor for pagination — pass the nextCursor from the previous response. Omit on the first request.'),
        size: zod_1.z.coerce
            .number()
            .int()
            .min(1, 'Size (limit) must be greater than 0')
            .max(1000, 'Size (limit) must be less than 1000')
            .describe('Number of results to return, no more than 1000')
            .default(250),
    });
    GetUsersStreamCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            users: zod_1.z.array(models_1.ExtendedUsersSchema),
            nextCursor: zod_1.z
                .string()
                .nullable()
                .describe('Cursor to fetch the next page, or null if there are no more results'),
            hasMore: zod_1.z.boolean().describe('Whether there are more results to fetch'),
        }),
    });
})(GetUsersStreamCommand || (exports.GetUsersStreamCommand = GetUsersStreamCommand = {}));
