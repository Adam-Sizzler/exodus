import { z } from 'zod';
export declare namespace GetStatsNodesRealtimeUsageCommand {
    const url: "/api/bandwidth-stats/nodes/realtime";
    const TSQ_url: "/api/bandwidth-stats/nodes/realtime";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodArray<z.ZodObject<{
            nodeUuid: z.ZodString;
            nodeName: z.ZodString;
            countryCode: z.ZodString;
            downloadBytes: z.ZodNumber;
            uploadBytes: z.ZodNumber;
            totalBytes: z.ZodNumber;
            downloadSpeedBps: z.ZodNumber;
            uploadSpeedBps: z.ZodNumber;
            totalSpeedBps: z.ZodNumber;
        }, "strip", z.ZodTypeAny, {
            nodeUuid: string;
            nodeName: string;
            countryCode: string;
            downloadBytes: number;
            uploadBytes: number;
            totalBytes: number;
            downloadSpeedBps: number;
            uploadSpeedBps: number;
            totalSpeedBps: number;
        }, {
            nodeUuid: string;
            nodeName: string;
            countryCode: string;
            downloadBytes: number;
            uploadBytes: number;
            totalBytes: number;
            downloadSpeedBps: number;
            uploadSpeedBps: number;
            totalSpeedBps: number;
        }>, "many">;
    }, "strip", z.ZodTypeAny, {
        response: {
            nodeUuid: string;
            nodeName: string;
            countryCode: string;
            downloadBytes: number;
            uploadBytes: number;
            totalBytes: number;
            downloadSpeedBps: number;
            uploadSpeedBps: number;
            totalSpeedBps: number;
        }[];
    }, {
        response: {
            nodeUuid: string;
            nodeName: string;
            countryCode: string;
            downloadBytes: number;
            uploadBytes: number;
            totalBytes: number;
            downloadSpeedBps: number;
            uploadSpeedBps: number;
            totalSpeedBps: number;
        }[];
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-stats-nodes-realtime-usage.command.d.ts.map