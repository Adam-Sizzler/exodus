"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetExodusHealthCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var GetExodusHealthCommand;
(function (GetExodusHealthCommand) {
    GetExodusHealthCommand.url = api_1.REST_API.SYSTEM.HEALTH;
    GetExodusHealthCommand.TSQ_url = GetExodusHealthCommand.url;
    GetExodusHealthCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SYSTEM_ROUTES.HEALTH, 'get', 'Get Exodus Health', { scope: 'exodus-health', kind: 'read' });
    GetExodusHealthCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            runtimeMetrics: zod_1.z.array(zod_1.z.object({
                rss: zod_1.z.number().optional().default(0),
                heapUsed: zod_1.z.number().optional().default(0),
                heapTotal: zod_1.z.number().optional().default(0),
                external: zod_1.z.number().optional().default(0),
                arrayBuffers: zod_1.z.number().optional().default(0),
                eventLoopDelayMs: zod_1.z.number().optional().default(0),
                eventLoopP99Ms: zod_1.z.number().optional().default(0),
                activeHandles: zod_1.z.number().optional().default(0),
                uptime: zod_1.z.number().optional().default(0),
                pid: zod_1.z.number().optional().default(0),
                timestamp: zod_1.z.number().optional().default(0),
                instanceId: zod_1.z.union([zod_1.z.string(), zod_1.z.number()]).optional(),
                instanceType: zod_1.z.string().optional(),
                name: zod_1.z.string().optional(),
                startedAt: zod_1.z.string().optional(),
                uptimeSeconds: zod_1.z.number().optional(),
                runtime: zod_1.z.object({
                    language: zod_1.z.string().optional(),
                    version: zod_1.z.string().optional(),
                    goos: zod_1.z.string().optional(),
                    goarch: zod_1.z.string().optional(),
                }).optional(),
                cpu: zod_1.z.object({
                    cores: zod_1.z.number().optional(),
                    gomaxprocs: zod_1.z.number().optional(),
                    processPercent: zod_1.z.number().optional(),
                    processCpuSeconds: zod_1.z.number().optional(),
                    processCpuPercentMode: zod_1.z.string().optional(),
                }).optional(),
                memory: zod_1.z.object({
                    rssBytes: zod_1.z.number().optional(),
                    allocBytes: zod_1.z.number().optional(),
                    totalAllocBytes: zod_1.z.number().optional(),
                    sysBytes: zod_1.z.number().optional(),
                    heapAllocBytes: zod_1.z.number().optional(),
                    heapSysBytes: zod_1.z.number().optional(),
                    heapIdleBytes: zod_1.z.number().optional(),
                    heapInuseBytes: zod_1.z.number().optional(),
                    heapReleasedBytes: zod_1.z.number().optional(),
                    stackInuseBytes: zod_1.z.number().optional(),
                    otherSysBytes: zod_1.z.number().optional(),
                    heapUsedPercent: zod_1.z.number().optional(),
                }).optional(),
                gc: zod_1.z.object({
                    numGc: zod_1.z.number().optional(),
                    pauseTotalNs: zod_1.z.number().optional(),
                    lastPauseNs: zod_1.z.number().optional(),
                    lastGcUnixNano: zod_1.z.number().optional(),
                    gcCpuFraction: zod_1.z.number().optional(),
                    pauseP99Ms: zod_1.z.number().optional(),
                    pauseP99Source: zod_1.z.string().optional(),
                    gogc: zod_1.z.number().optional(),
                }).optional(),
                scheduler: zod_1.z.object({
                    goroutines: zod_1.z.number().optional(),
                    cgoCalls: zod_1.z.number().optional(),
                    schedulerDelayMs: zod_1.z.number().optional(),
                    schedulerP99Ms: zod_1.z.number().optional(),
                    schedulerLatencySource: zod_1.z.string().optional(),
                }).optional(),
                process: zod_1.z.object({
                    openFileDescriptors: zod_1.z.number().optional(),
                    threads: zod_1.z.number().optional(),
                }).optional(),
                collectedAt: zod_1.z.string().optional(),
                collectedAtUnix: zod_1.z.number().optional(),
            })),
            runtimeSummary: zod_1.z.object({
                totalProcesses: zod_1.z.number().optional(),
                startedAt: zod_1.z.string().optional(),
                uptimeSeconds: zod_1.z.number().optional(),
                totalRssBytes: zod_1.z.number().optional(),
                heapAllocBytes: zod_1.z.number().optional(),
                averageCpuPercent: zod_1.z.number().optional(),
                averageGoroutines: zod_1.z.number().optional(),
                averageSchedulerDelayMs: zod_1.z.number().optional(),
                averageSchedulerP99Ms: zod_1.z.number().optional(),
            }).optional(),
        }),
    });

})(GetExodusHealthCommand || (exports.GetExodusHealthCommand = GetExodusHealthCommand = {}));

exports.GetExodusHealthCommand = exports.GetExodusHealthCommand;
