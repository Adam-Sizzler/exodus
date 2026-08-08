"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetHttpStatsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var GetHttpStatsCommand;
(function (GetHttpStatsCommand) {
    GetHttpStatsCommand.url = api_1.REST_API.SYSTEM.STATS.HTTP;
    GetHttpStatsCommand.TSQ_url = GetHttpStatsCommand.url;
    GetHttpStatsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SYSTEM_ROUTES.STATS.HTTP, 'get', 'Get HTTP Stats', { scope: 'http', kind: 'read' });
    GetHttpStatsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            routes: zod_1.z.array(zod_1.z.object({
                method: zod_1.z.string(),
                route: zod_1.z.string(),
                count: zod_1.z.int32().nonnegative(),
            })),
            total: zod_1.z.int32().nonnegative(),
        }),
    });

})(GetHttpStatsCommand || (exports.GetHttpStatsCommand = GetHttpStatsCommand = {}));
