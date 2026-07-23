"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetRemnawaveHealthCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var GetRemnawaveHealthCommand;
(function (GetRemnawaveHealthCommand) {
    GetRemnawaveHealthCommand.url = api_1.REST_API.SYSTEM.HEALTH;
    GetRemnawaveHealthCommand.TSQ_url = GetRemnawaveHealthCommand.url;
    GetRemnawaveHealthCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SYSTEM_ROUTES.HEALTH, 'get', 'Get Remnawave Health', { scope: 'remnawave-health', kind: 'read' });
    GetRemnawaveHealthCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            runtimeMetrics: zod_1.z.array(zod_1.z.object({
                rss: zod_1.z.number(),
                heapUsed: zod_1.z.number(),
                heapTotal: zod_1.z.number(),
                external: zod_1.z.number(),
                arrayBuffers: zod_1.z.number(),
                eventLoopDelayMs: zod_1.z.number(),
                eventLoopP99Ms: zod_1.z.number(),
                activeHandles: zod_1.z.number(),
                uptime: zod_1.z.number(),
                pid: zod_1.z.number(),
                timestamp: zod_1.z.number(),
                instanceId: zod_1.z.string(),
                instanceType: zod_1.z.string(),
            })),
        }),
    });
})(GetRemnawaveHealthCommand || (exports.GetRemnawaveHealthCommand = GetRemnawaveHealthCommand = {}));
