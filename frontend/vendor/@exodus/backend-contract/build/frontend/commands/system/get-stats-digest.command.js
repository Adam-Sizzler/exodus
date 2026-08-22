"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetStatsDigestCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var GetStatsDigestCommand;
(function (GetStatsDigestCommand) {
    GetStatsDigestCommand.url = api_1.REST_API.SYSTEM.STATS.DIGEST;
    GetStatsDigestCommand.TSQ_url = GetStatsDigestCommand.url;
    GetStatsDigestCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SYSTEM_ROUTES.STATS.DIGEST, 'get', 'Get Stats Digest', { scope: 'digest', kind: 'read' }, 'Aggregated statistics for a datetime range [start, end): created and expired users, total traffic, traffic spent by users created within the range and new HWID devices. Per-user traffic history is stored with daily granularity (UTC), so the "traffic by new users" metric snaps to whole days at the range edges.');
    GetStatsDigestCommand.RequestQuerySchema = zod_1.z.object({
        start: zod_1.z.iso
            .datetime({ offset: true })
            .describe('Start of the range, ISO 8601 datetime with timezone (e.g. 2026-07-15T00:00:00Z). Inclusive.'),
        end: zod_1.z.iso
            .datetime({ offset: true })
            .describe('End of the range, ISO 8601 datetime with timezone (e.g. 2026-07-16T00:00:00Z). Exclusive.'),
    });
    GetStatsDigestCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            users: zod_1.z.object({
                createdCount: zod_1.z.number(),
                expiredCount: zod_1.z.number(),
            }),
            traffic: zod_1.z.object({
                totalBytes: zod_1.z.string(),
                byUsersCreatedInRangeBytes: zod_1.z.string(),
            }),
            hwidDevices: zod_1.z.object({
                createdCount: zod_1.z.number(),
            }),
        }),
    });
})(GetStatsDigestCommand || (exports.GetStatsDigestCommand = GetStatsDigestCommand = {}));
