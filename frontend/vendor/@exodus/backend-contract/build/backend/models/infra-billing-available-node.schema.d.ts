export declare const InfraBillingAvailableNodeSchema: import("zod").ZodObject<Pick<{
    uuid: import("zod").ZodString;
    name: import("zod").ZodString;
    address: import("zod").ZodString;
    port: import("zod").ZodNullable<import("zod").ZodNumber>;
    proxyUrl: import("zod").ZodNullable<import("zod").ZodString>;
    isConnected: import("zod").ZodBoolean;
    isDisabled: import("zod").ZodBoolean;
    isConnecting: import("zod").ZodBoolean;
    lastStatusChange: import("zod").ZodNullable<import("zod").ZodEffects<import("zod").ZodString, Date, string>>;
    lastStatusMessage: import("zod").ZodNullable<import("zod").ZodString>;
    isTrafficTrackingActive: import("zod").ZodBoolean;
    trafficResetDay: import("zod").ZodNullable<import("zod").ZodNumber>;
    trafficLimitBytes: import("zod").ZodNullable<import("zod").ZodNumber>;
    trafficUsedBytes: import("zod").ZodNullable<import("zod").ZodNumber>;
    notifyPercent: import("zod").ZodNullable<import("zod").ZodNumber>;
    viewPosition: import("zod").ZodNumber;
    countryCode: import("zod").ZodString;
    consumptionMultiplier: import("zod").ZodNumber;
    nodeConsumptionMultiplier: import("zod").ZodNumber;
    tags: import("zod").ZodArray<import("zod").ZodString, "many">;
    createdAt: import("zod").ZodEffects<import("zod").ZodString, Date, string>;
    updatedAt: import("zod").ZodEffects<import("zod").ZodString, Date, string>;
    configProfile: import("zod").ZodObject<{
        activeConfigProfileUuid: import("zod").ZodNullable<import("zod").ZodString>;
        activeInbounds: import("zod").ZodArray<import("zod").ZodObject<{
            uuid: import("zod").ZodString;
            profileUuid: import("zod").ZodString;
            tag: import("zod").ZodString;
            type: import("zod").ZodString;
            network: import("zod").ZodNullable<import("zod").ZodString>;
            security: import("zod").ZodNullable<import("zod").ZodString>;
            port: import("zod").ZodNullable<import("zod").ZodNumber>;
            rawInbound: import("zod").ZodNullable<import("zod").ZodUnknown>;
        }, "strip", import("zod").ZodTypeAny, {
            uuid: string;
            type: string;
            profileUuid: string;
            tag: string;
            network: string | null;
            security: string | null;
            port: number | null;
            rawInbound?: unknown;
        }, {
            uuid: string;
            type: string;
            profileUuid: string;
            tag: string;
            network: string | null;
            security: string | null;
            port: number | null;
            rawInbound?: unknown;
        }>, "many">;
    }, "strip", import("zod").ZodTypeAny, {
        activeConfigProfileUuid: string | null;
        activeInbounds: {
            uuid: string;
            type: string;
            profileUuid: string;
            tag: string;
            network: string | null;
            security: string | null;
            port: number | null;
            rawInbound?: unknown;
        }[];
    }, {
        activeConfigProfileUuid: string | null;
        activeInbounds: {
            uuid: string;
            type: string;
            profileUuid: string;
            tag: string;
            network: string | null;
            security: string | null;
            port: number | null;
            rawInbound?: unknown;
        }[];
    }>;
    providerUuid: import("zod").ZodNullable<import("zod").ZodString>;
    provider: import("zod").ZodNullable<import("zod").ZodObject<{
        uuid: import("zod").ZodString;
        name: import("zod").ZodString;
        faviconLink: import("zod").ZodNullable<import("zod").ZodString>;
        loginUrl: import("zod").ZodNullable<import("zod").ZodString>;
        createdAt: import("zod").ZodEffects<import("zod").ZodString, Date, string>;
        updatedAt: import("zod").ZodEffects<import("zod").ZodString, Date, string>;
    }, "strip", import("zod").ZodTypeAny, {
        uuid: string;
        name: string;
        createdAt: Date;
        updatedAt: Date;
        faviconLink: string | null;
        loginUrl: string | null;
    }, {
        uuid: string;
        name: string;
        createdAt: string;
        updatedAt: string;
        faviconLink: string | null;
        loginUrl: string | null;
    }>>;
    activePluginUuid: import("zod").ZodNullable<import("zod").ZodString>;
    system: import("zod").ZodNullable<import("zod").ZodObject<{
        info: import("zod").ZodObject<{
            arch: import("zod").ZodString;
            cpus: import("zod").ZodNumber;
            cpuModel: import("zod").ZodString;
            memoryTotal: import("zod").ZodNumber;
            hostname: import("zod").ZodString;
            platform: import("zod").ZodString;
            release: import("zod").ZodString;
            type: import("zod").ZodString;
            version: import("zod").ZodString;
            networkInterfaces: import("zod").ZodArray<import("zod").ZodString, "many">;
        }, "strip", import("zod").ZodTypeAny, {
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
        stats: import("zod").ZodObject<{
            memoryFree: import("zod").ZodNumber;
            memoryUsed: import("zod").ZodNumber;
            uptime: import("zod").ZodNumber;
            loadAvg: import("zod").ZodArray<import("zod").ZodNumber, "many">;
            interface: import("zod").ZodNullable<import("zod").ZodObject<{
                interface: import("zod").ZodString;
                rxBytesPerSec: import("zod").ZodNumber;
                txBytesPerSec: import("zod").ZodNumber;
                rxTotal: import("zod").ZodNumber;
                txTotal: import("zod").ZodNumber;
            }, "strip", import("zod").ZodTypeAny, {
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
        }, "strip", import("zod").ZodTypeAny, {
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
    }, "strip", import("zod").ZodTypeAny, {
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
    }>>;
    versions: import("zod").ZodNullable<import("zod").ZodObject<{
        xray: import("zod").ZodString;
        node: import("zod").ZodString;
    }, "strip", import("zod").ZodTypeAny, {
        node: string;
        singbox: string;
    }, {
        node: string;
        singbox: string;
    }>>;
    singboxUptime: import("zod").ZodNumber;
    usersOnline: import("zod").ZodNumber;
    note: import("zod").ZodNullable<import("zod").ZodString>;
}, "uuid" | "name" | "countryCode">, "strip", import("zod").ZodTypeAny, {
    uuid: string;
    name: string;
    countryCode: string;
}, {
    uuid: string;
    name: string;
    countryCode: string;
}>;
//# sourceMappingURL=infra-billing-available-node.schema.d.ts.map