"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetInternalSquadUserUsageCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
var GetInternalSquadUserUsageCommand;
(function (GetInternalSquadUserUsageCommand) {
    GetInternalSquadUserUsageCommand.url = api_1.REST_API.BANDWIDTH_STATS.INTERNAL_SQUADS.USER_USAGE;
    GetInternalSquadUserUsageCommand.TSQ_url = GetInternalSquadUserUsageCommand.url(':squadUuid', ':userId');
    GetInternalSquadUserUsageCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.BANDWIDTH_STATS_ROUTES.INTERNAL_SQUADS.USER_USAGE(':squadUuid', ':userId'), 'get', 'Get a single user daily traffic usage on the internal squad nodes for a period', { scope: 'internal-squad-user-usage', kind: 'read' }, 'Returns users whose total usage over the period on the given nodes is >= minTotalBytes, scoped to the nodes reachable via the Internal Squad inbounds. Every day in the range is present (zero-filled). Underlying usage data is flushed to the database roughly every 2 minutes.');
    GetInternalSquadUserUsageCommand.RequestParamSchema = zod_1.z.object({
        squadUuid: zod_1.z.uuid().describe('Internal squad UUID'),
        userId: models_1.numberParamSchema,
    });
    GetInternalSquadUserUsageCommand.RequestQuerySchema = zod_1.z.object({
        start: zod_1.z.iso.date().describe('Start date (YYYY-MM-DD)'),
        end: zod_1.z.iso.date().describe('End date (YYYY-MM-DD)'),
    });
    GetInternalSquadUserUsageCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            days: zod_1.z.array(zod_1.z.object({
                date: zod_1.z.string().describe('Day (YYYY-MM-DD)'),
                nodes: zod_1.z.array(zod_1.z.object({
                    uuid: zod_1.z.uuid(),
                    totalBytes: zod_1.z
                        .number()
                        .describe('Used bytes on this node that day (raw bytes)'),
                })),
            })),
        }),
    });
})(GetInternalSquadUserUsageCommand || (exports.GetInternalSquadUserUsageCommand = GetInternalSquadUserUsageCommand = {}));
