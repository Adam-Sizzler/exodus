import { z } from 'zod';

import { REST_API, SYSTEM_ROUTES } from '../../api';
import { getEndpointDetails } from '../../constants';

export namespace GetExodusHealthCommand {
    export const url = REST_API.SYSTEM.HEALTH;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        SYSTEM_ROUTES.HEALTH,
        'get',
        'Get Exodus Health',
        { scope: 'exodus-health', kind: 'read' },
    );

    export const ResponseSchema = z.object({
        response: z.object({
            runtimeMetrics: z.array(
                z.object({
                    rss: z.number().optional().default(0),
                    heapUsed: z.number().optional().default(0),
                    heapTotal: z.number().optional().default(0),
                    external: z.number().optional().default(0),
                    arrayBuffers: z.number().optional().default(0),
                    eventLoopDelayMs: z.number().optional().default(0),
                    eventLoopP99Ms: z.number().optional().default(0),
                    activeHandles: z.number().optional().default(0),
                    uptime: z.number().optional().default(0),
                    pid: z.number().optional().default(0),
                    timestamp: z.number().optional().default(0),
                    instanceId: z.union([z.string(), z.number()]).optional(),
                    instanceType: z.string().optional(),
                    name: z.string().optional(),
                    startedAt: z.string().optional(),
                    uptimeSeconds: z.number().optional(),
                    runtime: z
                        .object({
                            language: z.string().optional(),
                            version: z.string().optional(),
                            goos: z.string().optional(),
                            goarch: z.string().optional(),
                        })
                        .optional(),
                    cpu: z
                        .object({
                            cores: z.number().optional(),
                            gomaxprocs: z.number().optional(),
                            processPercent: z.number().optional(),
                            processCpuSeconds: z.number().optional(),
                            processCpuPercentMode: z.string().optional(),
                        })
                        .optional(),
                    memory: z
                        .object({
                            rssBytes: z.number().optional(),
                            allocBytes: z.number().optional(),
                            totalAllocBytes: z.number().optional(),
                            sysBytes: z.number().optional(),
                            heapAllocBytes: z.number().optional(),
                            heapSysBytes: z.number().optional(),
                            heapIdleBytes: z.number().optional(),
                            heapInuseBytes: z.number().optional(),
                            heapReleasedBytes: z.number().optional(),
                            stackInuseBytes: z.number().optional(),
                            otherSysBytes: z.number().optional(),
                            heapUsedPercent: z.number().optional(),
                        })
                        .optional(),
                    gc: z
                        .object({
                            numGc: z.number().optional(),
                            pauseTotalNs: z.number().optional(),
                            lastPauseNs: z.number().optional(),
                            lastGcUnixNano: z.number().optional(),
                            gcCpuFraction: z.number().optional(),
                            pauseP99Ms: z.number().optional(),
                            pauseP99Source: z.string().optional(),
                            gogc: z.number().optional(),
                        })
                        .optional(),
                    scheduler: z
                        .object({
                            goroutines: z.number().optional(),
                            cgoCalls: z.number().optional(),
                            schedulerDelayMs: z.number().optional(),
                            schedulerP99Ms: z.number().optional(),
                            schedulerLatencySource: z.string().optional(),
                        })
                        .optional(),
                    process: z
                        .object({
                            openFileDescriptors: z.number().optional(),
                            threads: z.number().optional(),
                        })
                        .optional(),
                    collectedAt: z.string().optional(),
                    collectedAtUnix: z.number().optional(),
                }),
            ),
            runtimeSummary: z
                .object({
                    totalProcesses: z.number().optional(),
                    startedAt: z.string().optional(),
                    uptimeSeconds: z.number().optional(),
                    totalRssBytes: z.number().optional(),
                    heapAllocBytes: z.number().optional(),
                    averageCpuPercent: z.number().optional(),
                    averageGoroutines: z.number().optional(),
                    averageSchedulerDelayMs: z.number().optional(),
                    averageSchedulerP99Ms: z.number().optional(),
                })
                .optional(),
        }),
    });

    export type Response = z.infer<typeof ResponseSchema>;
}
