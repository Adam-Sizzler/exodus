import { z } from 'zod';
export declare const NetworkInterfaceSchema: z.ZodObject<{
    interface: z.ZodString;
    rxBytesPerSec: z.ZodNumber;
    txBytesPerSec: z.ZodNumber;
    rxTotal: z.ZodNumber;
    txTotal: z.ZodNumber;
}, z.core.$strip>;
export declare const NodeSystemInfoSchema: z.ZodObject<{
    arch: z.ZodString;
    cpus: z.ZodInt;
    cpuModel: z.ZodString;
    memoryTotal: z.ZodNumber;
    hostname: z.ZodString;
    platform: z.ZodString;
    release: z.ZodString;
    type: z.ZodString;
    version: z.ZodString;
    networkInterfaces: z.ZodArray<z.ZodString>;
}, z.core.$strip>;
export declare const NodeSystemStatsSchema: z.ZodObject<{
    memoryFree: z.ZodNumber;
    memoryUsed: z.ZodNumber;
    uptime: z.ZodNumber;
    loadAvg: z.ZodArray<z.ZodNumber>;
    interface: z.ZodNullable<z.ZodObject<{
        interface: z.ZodString;
        rxBytesPerSec: z.ZodNumber;
        txBytesPerSec: z.ZodNumber;
        rxTotal: z.ZodNumber;
        txTotal: z.ZodNumber;
    }, z.core.$strip>>;
}, z.core.$strip>;
export type TNodeSystemStats = z.infer<typeof NodeSystemStatsSchema>;
export declare const NodeSystemSchema: z.ZodObject<{
    info: z.ZodObject<{
        arch: z.ZodString;
        cpus: z.ZodInt;
        cpuModel: z.ZodString;
        memoryTotal: z.ZodNumber;
        hostname: z.ZodString;
        platform: z.ZodString;
        release: z.ZodString;
        type: z.ZodString;
        version: z.ZodString;
        networkInterfaces: z.ZodArray<z.ZodString>;
    }, z.core.$strip>;
    stats: z.ZodObject<{
        memoryFree: z.ZodNumber;
        memoryUsed: z.ZodNumber;
        uptime: z.ZodNumber;
        loadAvg: z.ZodArray<z.ZodNumber>;
        interface: z.ZodNullable<z.ZodObject<{
            interface: z.ZodString;
            rxBytesPerSec: z.ZodNumber;
            txBytesPerSec: z.ZodNumber;
            rxTotal: z.ZodNumber;
            txTotal: z.ZodNumber;
        }, z.core.$strip>>;
    }, z.core.$strip>;
}, z.core.$strip>;
export type TNetworkInterface = z.infer<typeof NetworkInterfaceSchema>;
export type TNodeSystemInfo = z.infer<typeof NodeSystemInfoSchema>;
export type TNodeSystem = z.infer<typeof NodeSystemSchema>;
//# sourceMappingURL=node-system.schema.d.ts.map