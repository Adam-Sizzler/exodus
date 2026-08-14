"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetInternalSquadUsageCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var GetInternalSquadUsageCommand;
(function (GetInternalSquadUsageCommand) {
    GetInternalSquadUsageCommand.url = api_1.REST_API.BANDWIDTH_STATS.INTERNAL_SQUADS.GET_USAGE;
    GetInternalSquadUsageCommand.TSQ_url = GetInternalSquadUsageCommand.url(':uuid');
    GetInternalSquadUsageCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.BANDWIDTH_STATS_ROUTES.INTERNAL_SQUADS.GET_USAGE(':uuid'), 'get', 'Get internal squad users traffic usage for a period', { scope: 'internal-squad-usage', kind: 'read' }, 'Returns users whose total usage over the period on the given nodes is >= minTotalBytes, scoped to the nodes reachable via the internal squad inbounds. Underlying usage data is flushed to the database roughly every 2 minutes.');
    GetInternalSquadUsageCommand.RequestParamSchema = zod_1.z.object({
        uuid: zod_1.z.uuid().describe('Internal squad UUID'),
    });
    GetInternalSquadUsageCommand.RequestQuerySchema = zod_1.z.object({
        start: zod_1.z.iso.date().describe('Start date (YYYY-MM-DD)'),
        end: zod_1.z.iso.date().describe('End date (YYYY-MM-DD)'),
        minTotalBytes: zod_1.z.coerce
            .number()
            .min(0)
            .optional()
            .default(0)
            .describe('Only include users whose total usage over the period is >= this (bytes)'),
        limit: zod_1.z.coerce
            .number()
            .min(1)
            .max(1000)
            .optional()
            .default(250)
            .describe('Number of users to return, no more than 1000'),
        cursor: zod_1.z.coerce
            .number()
            .optional()
            .describe('Pass the nextCursor from the previous response. Omit on the first request.'),
    });
    GetInternalSquadUsageCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            squadUuid: zod_1.z.uuid(),
            users: zod_1.z.array(zod_1.z.object({
                id: zod_1.z.number(),
                totalBytes: zod_1.z.number().describe('Total used bytes over the period (raw bytes)'),
            })),
            nextCursor: zod_1.z
                .string()
                .nullable()
                .describe('Cursor to fetch the next page, or null if there are no more results'),
            hasMore: zod_1.z.boolean().describe('Whether there are more results to fetch'),
        }),
    });
})(GetInternalSquadUsageCommand || (exports.GetInternalSquadUsageCommand = GetInternalSquadUsageCommand = {}));
