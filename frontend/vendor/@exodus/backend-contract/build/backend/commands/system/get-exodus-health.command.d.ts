import { z } from 'zod';
export declare namespace GetExodusHealthCommand {
    const url: "/api/system/health";
    const TSQ_url: "/api/system/health";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            runtimeMetrics: z.ZodArray<z.ZodObject<{
                rss: z.ZodDefault<z.ZodOptional<z.ZodNumber>>;
                heapUsed: z.ZodDefault<z.ZodOptional<z.ZodNumber>>;
                heapTotal: z.ZodDefault<z.ZodOptional<z.ZodNumber>>;
                external: z.ZodDefault<z.ZodOptional<z.ZodNumber>>;
                arrayBuffers: z.ZodDefault<z.ZodOptional<z.ZodNumber>>;
                eventLoopDelayMs: z.ZodDefault<z.ZodOptional<z.ZodNumber>>;
                eventLoopP99Ms: z.ZodDefault<z.ZodOptional<z.ZodNumber>>;
                activeHandles: z.ZodDefault<z.ZodOptional<z.ZodNumber>>;
                uptime: z.ZodDefault<z.ZodOptional<z.ZodNumber>>;
                pid: z.ZodDefault<z.ZodOptional<z.ZodNumber>>;
                timestamp: z.ZodDefault<z.ZodOptional<z.ZodNumber>>;
                instanceId: z.ZodOptional<z.ZodUnion<readonly [z.ZodString, z.ZodNumber]>>;
                instanceType: z.ZodOptional<z.ZodString>;
                name: z.ZodOptional<z.ZodString>;
                startedAt: z.ZodOptional<z.ZodString>;
                uptimeSeconds: z.ZodOptional<z.ZodNumber>;
                runtime: z.ZodOptional<z.ZodObject<{
                    language: z.ZodOptional<z.ZodString>;
                    version: z.ZodOptional<z.ZodString>;
                    goos: z.ZodOptional<z.ZodString>;
                    goarch: z.ZodOptional<z.ZodString>;
                }, z.core.$strip>>;
                cpu: z.ZodOptional<z.ZodObject<{
                    cores: z.ZodOptional<z.ZodNumber>;
                    gomaxprocs: z.ZodOptional<z.ZodNumber>;
                    processPercent: z.ZodOptional<z.ZodNumber>;
                    processCpuSeconds: z.ZodOptional<z.ZodNumber>;
                    processCpuPercentMode: z.ZodOptional<z.ZodString>;
                }, z.core.$strip>>;
                memory: z.ZodOptional<z.ZodObject<{
                    rssBytes: z.ZodOptional<z.ZodNumber>;
                    allocBytes: z.ZodOptional<z.ZodNumber>;
                    totalAllocBytes: z.ZodOptional<z.ZodNumber>;
                    sysBytes: z.ZodOptional<z.ZodNumber>;
                    heapAllocBytes: z.ZodOptional<z.ZodNumber>;
                    heapSysBytes: z.ZodOptional<z.ZodNumber>;
                    heapIdleBytes: z.ZodOptional<z.ZodNumber>;
                    heapInuseBytes: z.ZodOptional<z.ZodNumber>;
                    heapReleasedBytes: z.ZodOptional<z.ZodNumber>;
                    stackInuseBytes: z.ZodOptional<z.ZodNumber>;
                    otherSysBytes: z.ZodOptional<z.ZodNumber>;
                    heapUsedPercent: z.ZodOptional<z.ZodNumber>;
                }, z.core.$strip>>;
                gc: z.ZodOptional<z.ZodObject<{
                    numGc: z.ZodOptional<z.ZodNumber>;
                    pauseTotalNs: z.ZodOptional<z.ZodNumber>;
                    lastPauseNs: z.ZodOptional<z.ZodNumber>;
                    lastGcUnixNano: z.ZodOptional<z.ZodNumber>;
                    gcCpuFraction: z.ZodOptional<z.ZodNumber>;
                    pauseP99Ms: z.ZodOptional<z.ZodNumber>;
                    pauseP99Source: z.ZodOptional<z.ZodString>;
                    gogc: z.ZodOptional<z.ZodNumber>;
                }, z.core.$strip>>;
                scheduler: z.ZodOptional<z.ZodObject<{
                    goroutines: z.ZodOptional<z.ZodNumber>;
                    cgoCalls: z.ZodOptional<z.ZodNumber>;
                    schedulerDelayMs: z.ZodOptional<z.ZodNumber>;
                    schedulerP99Ms: z.ZodOptional<z.ZodNumber>;
                    schedulerLatencySource: z.ZodOptional<z.ZodString>;
                }, z.core.$strip>>;
                process: z.ZodOptional<z.ZodObject<{
                    openFileDescriptors: z.ZodOptional<z.ZodNumber>;
                    threads: z.ZodOptional<z.ZodNumber>;
                }, z.core.$strip>>;
                collectedAt: z.ZodOptional<z.ZodString>;
                collectedAtUnix: z.ZodOptional<z.ZodNumber>;
            }, z.core.$strip>>;
            runtimeSummary: z.ZodOptional<z.ZodObject<{
                totalProcesses: z.ZodOptional<z.ZodNumber>;
                startedAt: z.ZodOptional<z.ZodString>;
                uptimeSeconds: z.ZodOptional<z.ZodNumber>;
                totalRssBytes: z.ZodOptional<z.ZodNumber>;
                heapAllocBytes: z.ZodOptional<z.ZodNumber>;
                averageCpuPercent: z.ZodOptional<z.ZodNumber>;
                averageGoroutines: z.ZodOptional<z.ZodNumber>;
                averageSchedulerDelayMs: z.ZodOptional<z.ZodNumber>;
                averageSchedulerP99Ms: z.ZodOptional<z.ZodNumber>;
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-exodus-health.command.d.ts.map