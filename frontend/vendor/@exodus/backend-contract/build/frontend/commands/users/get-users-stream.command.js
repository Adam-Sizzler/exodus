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
    GetUsersStreamCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.STREAM, 'get', 'Get all users using cursor-based (keyset) pagination with filtering options', { scope: 'stream', kind: 'read' });
    GetUsersStreamCommand.RequestQuerySchema = zod_1.z.object({
        cursor: zod_1.z.coerce
            .number()
            .optional()
            .describe('Cursor for pagination — pass the nextCursor from the previous response. Omit on the first request.'),
        size: zod_1.z.coerce
            .number()
            .min(1)
            .max(1000)
            .describe('Number of results to return, no more than 1000')
            .optional()
            .default(250),
        // Filtering
        status: zod_1.z.enum(constants_1.USERS_STATUS).optional().describe('Status to filter users by'),
        trafficLimitStrategy: zod_1.z
            .enum(constants_1.RESET_PERIODS)
            .optional()
            .describe('Traffic limit strategy to filter users by'),
        telegramId: zod_1.z
            .string()
            .transform(Number)
            .pipe(zod_1.z.number().nonnegative())
            .optional()
            .describe('Telegram ID to filter users by'),
        email: zod_1.z.email().optional().describe('Email to filter users by'),
        tag: zod_1.z.string().optional().describe('Tag to filter users by'),
        externalSquadUuid: zod_1.z.uuid().optional().describe('External squad UUID to filter users by'),
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
