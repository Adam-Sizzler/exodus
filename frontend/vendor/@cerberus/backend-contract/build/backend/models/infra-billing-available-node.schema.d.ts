export declare const InfraBillingAvailableNodeSchema: import("zod").ZodObject<Pick<{
    uuid: import("zod").ZodString;
    name: import("zod").ZodString;
    address: import("zod").ZodString;
    port: import("zod").ZodNullable<import("zod").ZodNumber>;
    isConnected: import("zod").ZodBoolean;
    isDisabled: import("zod").ZodBoolean;
    isConnecting: import("zod").ZodBoolean;
    lastStatusChange: import("zod").ZodNullable<import("zod").ZodEffects<import("zod").ZodString, Date, string>>;
    lastStatusMessage: import("zod").ZodNullable<import("zod").ZodString>;
    xrayVersion: import("zod").ZodNullable<import("zod").ZodString>;
    nodeVersion: import("zod").ZodNullable<import("zod").ZodString>;
    xrayUptime: import("zod").ZodString;
    isTrafficTrackingActive: import("zod").ZodBoolean;
    trafficResetDay: import("zod").ZodNullable<import("zod").ZodNumber>;
    trafficLimitBytes: import("zod").ZodNullable<import("zod").ZodNumber>;
    trafficUsedBytes: import("zod").ZodNullable<import("zod").ZodNumber>;
    notifyPercent: import("zod").ZodNullable<import("zod").ZodNumber>;
    usersOnline: import("zod").ZodNullable<import("zod").ZodNumber>;
    viewPosition: import("zod").ZodNumber;
    countryCode: import("zod").ZodString;
    consumptionMultiplier: import("zod").ZodNumber;
    tags: import("zod").ZodArray<import("zod").ZodString, "many">;
    cpuCount: import("zod").ZodNullable<import("zod").ZodNumber>;
    cpuModel: import("zod").ZodNullable<import("zod").ZodString>;
    totalRam: import("zod").ZodNullable<import("zod").ZodString>;
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
            type: string;
            uuid: string;
            profileUuid: string;
            tag: string;
            network: string | null;
            security: string | null;
            port: number | null;
            rawInbound?: unknown;
        }, {
            type: string;
            uuid: string;
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
            type: string;
            uuid: string;
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
            type: string;
            uuid: string;
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
        createdAt: Date;
        updatedAt: Date;
        name: string;
        faviconLink: string | null;
        loginUrl: string | null;
    }, {
        uuid: string;
        createdAt: string;
        updatedAt: string;
        name: string;
        faviconLink: string | null;
        loginUrl: string | null;
    }>>;
}, "uuid" | "countryCode" | "name">, "strip", import("zod").ZodTypeAny, {
    uuid: string;
    countryCode: string;
    name: string;
}, {
    uuid: string;
    countryCode: string;
    name: string;
}>;
//# sourceMappingURL=infra-billing-available-node.schema.d.ts.map