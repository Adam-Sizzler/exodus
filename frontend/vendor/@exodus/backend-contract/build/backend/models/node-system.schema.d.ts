import { z } from 'zod';
export declare const NetworkInterfaceSchema: z.ZodObject<{
    interface: z.ZodString;
    rxBytesPerSec: z.ZodNumber;
    txBytesPerSec: z.ZodNumber;
    rxTotal: z.ZodNumber;
    txTotal: z.ZodNumber;
}, "strip", z.ZodTypeAny, {
    interface: string;
    rxBytesPerSec: number;
    txBytesPerSec: number;
    rxTotal: number;
    txTotal: number;
}, {
    interface: string;
    rxBytesPerSec: number;
    txBytesPerSec: number;
    rxTotal: number;
    txTotal: number;
}>;
export declare const NodeSystemInfoSchema: z.ZodObject<{
    arch: z.ZodString;
    cpus: z.ZodNumber;
    cpuModel: z.ZodString;
    memoryTotal: z.ZodNumber;
    hostname: z.ZodString;
    platform: z.ZodString;
    release: z.ZodString;
    type: z.ZodString;
    version: z.ZodString;
    networkInterfaces: z.ZodArray<z.ZodString, "many">;
}, "strip", z.ZodTypeAny, {
    type: string;
    version: string;
    platform: string;
    arch: string;
    cpus: number;
    cpuModel: string;
    memoryTotal: number;
    hostname: string;
    release: string;
    networkInterfaces: string[];
}, {
    type: string;
    version: string;
    platform: string;
    arch: string;
    cpus: number;
    cpuModel: string;
    memoryTotal: number;
    hostname: string;
    release: string;
    networkInterfaces: string[];
}>;
export declare const NodeSystemStatsSchema: z.ZodObject<{
    memoryFree: z.ZodNumber;
    memoryUsed: z.ZodNumber;
    uptime: z.ZodNumber;
    loadAvg: z.ZodArray<z.ZodNumber, "many">;
    interface: z.ZodNullable<z.ZodObject<{
        interface: z.ZodString;
        rxBytesPerSec: z.ZodNumber;
        txBytesPerSec: z.ZodNumber;
        rxTotal: z.ZodNumber;
        txTotal: z.ZodNumber;
    }, "strip", z.ZodTypeAny, {
        interface: string;
        rxBytesPerSec: number;
        txBytesPerSec: number;
        rxTotal: number;
        txTotal: number;
    }, {
        interface: string;
        rxBytesPerSec: number;
        txBytesPerSec: number;
        rxTotal: number;
        txTotal: number;
    }>>;
}, "strip", z.ZodTypeAny, {
    interface: {
        interface: string;
        rxBytesPerSec: number;
        txBytesPerSec: number;
        rxTotal: number;
        txTotal: number;
    } | null;
    memoryFree: number;
    memoryUsed: number;
    uptime: number;
    loadAvg: number[];
}, {
    interface: {
        interface: string;
        rxBytesPerSec: number;
        txBytesPerSec: number;
        rxTotal: number;
        txTotal: number;
    } | null;
    memoryFree: number;
    memoryUsed: number;
    uptime: number;
    loadAvg: number[];
}>;
export type TNodeSystemStats = z.infer<typeof NodeSystemStatsSchema>;
export declare const NodeSystemSchema: z.ZodObject<{
    info: z.ZodObject<{
        arch: z.ZodString;
        cpus: z.ZodNumber;
        cpuModel: z.ZodString;
        memoryTotal: z.ZodNumber;
        hostname: z.ZodString;
        platform: z.ZodString;
        release: z.ZodString;
        type: z.ZodString;
        version: z.ZodString;
        networkInterfaces: z.ZodArray<z.ZodString, "many">;
    }, "strip", z.ZodTypeAny, {
        type: string;
        version: string;
        platform: string;
        arch: string;
        cpus: number;
        cpuModel: string;
        memoryTotal: number;
        hostname: string;
        release: string;
        networkInterfaces: string[];
    }, {
        type: string;
        version: string;
        platform: string;
        arch: string;
        cpus: number;
        cpuModel: string;
        memoryTotal: number;
        hostname: string;
        release: string;
        networkInterfaces: string[];
    }>;
    stats: z.ZodObject<{
        memoryFree: z.ZodNumber;
        memoryUsed: z.ZodNumber;
        uptime: z.ZodNumber;
        loadAvg: z.ZodArray<z.ZodNumber, "many">;
        interface: z.ZodNullable<z.ZodObject<{
            interface: z.ZodString;
            rxBytesPerSec: z.ZodNumber;
            txBytesPerSec: z.ZodNumber;
            rxTotal: z.ZodNumber;
            txTotal: z.ZodNumber;
        }, "strip", z.ZodTypeAny, {
            interface: string;
            rxBytesPerSec: number;
            txBytesPerSec: number;
            rxTotal: number;
            txTotal: number;
        }, {
            interface: string;
            rxBytesPerSec: number;
            txBytesPerSec: number;
            rxTotal: number;
            txTotal: number;
        }>>;
    }, "strip", z.ZodTypeAny, {
        interface: {
            interface: string;
            rxBytesPerSec: number;
            txBytesPerSec: number;
            rxTotal: number;
            txTotal: number;
        } | null;
        memoryFree: number;
        memoryUsed: number;
        uptime: number;
        loadAvg: number[];
    }, {
        interface: {
            interface: string;
            rxBytesPerSec: number;
            txBytesPerSec: number;
            rxTotal: number;
            txTotal: number;
        } | null;
        memoryFree: number;
        memoryUsed: number;
        uptime: number;
        loadAvg: number[];
    }>;
}, "strip", z.ZodTypeAny, {
    stats: {
        interface: {
            interface: string;
            rxBytesPerSec: number;
            txBytesPerSec: number;
            rxTotal: number;
            txTotal: number;
        } | null;
        memoryFree: number;
        memoryUsed: number;
        uptime: number;
        loadAvg: number[];
    };
    info: {
        type: string;
        version: string;
        platform: string;
        arch: string;
        cpus: number;
        cpuModel: string;
        memoryTotal: number;
        hostname: string;
        release: string;
        networkInterfaces: string[];
    };
}, {
    stats: {
        interface: {
            interface: string;
            rxBytesPerSec: number;
            txBytesPerSec: number;
            rxTotal: number;
            txTotal: number;
        } | null;
        memoryFree: number;
        memoryUsed: number;
        uptime: number;
        loadAvg: number[];
    };
    info: {
        type: string;
        version: string;
        platform: string;
        arch: string;
        cpus: number;
        cpuModel: string;
        memoryTotal: number;
        hostname: string;
        release: string;
        networkInterfaces: string[];
    };
}>;
export type TNetworkInterface = z.infer<typeof NetworkInterfaceSchema>;
export type TNodeSystemInfo = z.infer<typeof NodeSystemInfoSchema>;
export type TNodeSystem = z.infer<typeof NodeSystemSchema>;
//# sourceMappingURL=node-system.schema.d.ts.map