import { z } from 'zod';
export declare namespace CreateNodeCommand {
    const url: "/api/nodes/";
    const TSQ_url: "/api/nodes/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        name: z.ZodString;
        address: z.ZodString;
        port: z.ZodOptional<z.ZodNumber>;
        isTrafficTrackingActive: z.ZodDefault<z.ZodOptional<z.ZodBoolean>>;
        trafficLimitBytes: z.ZodOptional<z.ZodNumber>;
        notifyPercent: z.ZodOptional<z.ZodNumber>;
        trafficResetDay: z.ZodOptional<z.ZodNumber>;
        countryCode: z.ZodDefault<z.ZodString>;
        consumptionMultiplier: z.ZodOptional<z.ZodEffects<z.ZodNumber, number, number>>;
        configProfile: z.ZodObject<{
            activeConfigProfileUuid: z.ZodString;
            activeInbounds: z.ZodArray<z.ZodString, "many">;
        }, "strip", z.ZodTypeAny, {
            activeConfigProfileUuid: string;
            activeInbounds: string[];
        }, {
            activeConfigProfileUuid: string;
            activeInbounds: string[];
        }>;
        providerUuid: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        tags: z.ZodOptional<z.ZodArray<z.ZodString, "many">>;
    }, "strip", z.ZodTypeAny, {
        countryCode: string;
        name: string;
        address: string;
        isTrafficTrackingActive: boolean;
        configProfile: {
            activeConfigProfileUuid: string;
            activeInbounds: string[];
        };
        tags?: string[] | undefined;
        port?: number | undefined;
        trafficLimitBytes?: number | undefined;
        trafficResetDay?: number | undefined;
        notifyPercent?: number | undefined;
        consumptionMultiplier?: number | undefined;
        providerUuid?: string | null | undefined;
    }, {
        name: string;
        address: string;
        configProfile: {
            activeConfigProfileUuid: string;
            activeInbounds: string[];
        };
        tags?: string[] | undefined;
        countryCode?: string | undefined;
        port?: number | undefined;
        trafficLimitBytes?: number | undefined;
        isTrafficTrackingActive?: boolean | undefined;
        trafficResetDay?: number | undefined;
        notifyPercent?: number | undefined;
        consumptionMultiplier?: number | undefined;
        providerUuid?: string | null | undefined;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            uuid: z.ZodString;
            name: z.ZodString;
            address: z.ZodString;
            port: z.ZodNullable<z.ZodNumber>;
            isConnected: z.ZodBoolean;
            isDisabled: z.ZodBoolean;
            isConnecting: z.ZodBoolean;
            lastStatusChange: z.ZodNullable<z.ZodEffects<z.ZodString, Date, string>>;
            lastStatusMessage: z.ZodNullable<z.ZodString>;
            xrayVersion: z.ZodNullable<z.ZodString>;
            nodeVersion: z.ZodNullable<z.ZodString>;
            xrayUptime: z.ZodString;
            isTrafficTrackingActive: z.ZodBoolean;
            trafficResetDay: z.ZodNullable<z.ZodNumber>;
            trafficLimitBytes: z.ZodNullable<z.ZodNumber>;
            trafficUsedBytes: z.ZodNullable<z.ZodNumber>;
            notifyPercent: z.ZodNullable<z.ZodNumber>;
            usersOnline: z.ZodNullable<z.ZodNumber>;
            viewPosition: z.ZodNumber;
            countryCode: z.ZodString;
            consumptionMultiplier: z.ZodNumber;
            tags: z.ZodArray<z.ZodString, "many">;
            cpuCount: z.ZodNullable<z.ZodNumber>;
            cpuModel: z.ZodNullable<z.ZodString>;
            totalRam: z.ZodNullable<z.ZodString>;
            createdAt: z.ZodEffects<z.ZodString, Date, string>;
            updatedAt: z.ZodEffects<z.ZodString, Date, string>;
            configProfile: z.ZodObject<{
                activeConfigProfileUuid: z.ZodNullable<z.ZodString>;
                activeInbounds: z.ZodArray<z.ZodObject<{
                    uuid: z.ZodString;
                    profileUuid: z.ZodString;
                    tag: z.ZodString;
                    type: z.ZodString;
                    network: z.ZodNullable<z.ZodString>;
                    security: z.ZodNullable<z.ZodString>;
                    port: z.ZodNullable<z.ZodNumber>;
                    rawInbound: z.ZodNullable<z.ZodUnknown>;
                }, "strip", z.ZodTypeAny, {
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
            }, "strip", z.ZodTypeAny, {
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
            providerUuid: z.ZodNullable<z.ZodString>;
            provider: z.ZodNullable<z.ZodObject<{
                uuid: z.ZodString;
                name: z.ZodString;
                faviconLink: z.ZodNullable<z.ZodString>;
                loginUrl: z.ZodNullable<z.ZodString>;
                createdAt: z.ZodEffects<z.ZodString, Date, string>;
                updatedAt: z.ZodEffects<z.ZodString, Date, string>;
            }, "strip", z.ZodTypeAny, {
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
        }, "strip", z.ZodTypeAny, {
            tags: string[];
            uuid: string;
            createdAt: Date;
            updatedAt: Date;
            provider: {
                uuid: string;
                createdAt: Date;
                updatedAt: Date;
                name: string;
                faviconLink: string | null;
                loginUrl: string | null;
            } | null;
            countryCode: string;
            name: string;
            port: number | null;
            viewPosition: number;
            trafficLimitBytes: number | null;
            address: string;
            isDisabled: boolean;
            isConnected: boolean;
            isConnecting: boolean;
            lastStatusChange: Date | null;
            lastStatusMessage: string | null;
            xrayVersion: string | null;
            nodeVersion: string | null;
            xrayUptime: string;
            isTrafficTrackingActive: boolean;
            trafficResetDay: number | null;
            trafficUsedBytes: number | null;
            notifyPercent: number | null;
            usersOnline: number | null;
            consumptionMultiplier: number;
            cpuCount: number | null;
            cpuModel: string | null;
            totalRam: string | null;
            configProfile: {
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
            };
            providerUuid: string | null;
        }, {
            tags: string[];
            uuid: string;
            createdAt: string;
            updatedAt: string;
            provider: {
                uuid: string;
                createdAt: string;
                updatedAt: string;
                name: string;
                faviconLink: string | null;
                loginUrl: string | null;
            } | null;
            countryCode: string;
            name: string;
            port: number | null;
            viewPosition: number;
            trafficLimitBytes: number | null;
            address: string;
            isDisabled: boolean;
            isConnected: boolean;
            isConnecting: boolean;
            lastStatusChange: string | null;
            lastStatusMessage: string | null;
            xrayVersion: string | null;
            nodeVersion: string | null;
            xrayUptime: string;
            isTrafficTrackingActive: boolean;
            trafficResetDay: number | null;
            trafficUsedBytes: number | null;
            notifyPercent: number | null;
            usersOnline: number | null;
            consumptionMultiplier: number;
            cpuCount: number | null;
            cpuModel: string | null;
            totalRam: string | null;
            configProfile: {
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
            };
            providerUuid: string | null;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            tags: string[];
            uuid: string;
            createdAt: Date;
            updatedAt: Date;
            provider: {
                uuid: string;
                createdAt: Date;
                updatedAt: Date;
                name: string;
                faviconLink: string | null;
                loginUrl: string | null;
            } | null;
            countryCode: string;
            name: string;
            port: number | null;
            viewPosition: number;
            trafficLimitBytes: number | null;
            address: string;
            isDisabled: boolean;
            isConnected: boolean;
            isConnecting: boolean;
            lastStatusChange: Date | null;
            lastStatusMessage: string | null;
            xrayVersion: string | null;
            nodeVersion: string | null;
            xrayUptime: string;
            isTrafficTrackingActive: boolean;
            trafficResetDay: number | null;
            trafficUsedBytes: number | null;
            notifyPercent: number | null;
            usersOnline: number | null;
            consumptionMultiplier: number;
            cpuCount: number | null;
            cpuModel: string | null;
            totalRam: string | null;
            configProfile: {
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
            };
            providerUuid: string | null;
        };
    }, {
        response: {
            tags: string[];
            uuid: string;
            createdAt: string;
            updatedAt: string;
            provider: {
                uuid: string;
                createdAt: string;
                updatedAt: string;
                name: string;
                faviconLink: string | null;
                loginUrl: string | null;
            } | null;
            countryCode: string;
            name: string;
            port: number | null;
            viewPosition: number;
            trafficLimitBytes: number | null;
            address: string;
            isDisabled: boolean;
            isConnected: boolean;
            isConnecting: boolean;
            lastStatusChange: string | null;
            lastStatusMessage: string | null;
            xrayVersion: string | null;
            nodeVersion: string | null;
            xrayUptime: string;
            isTrafficTrackingActive: boolean;
            trafficResetDay: number | null;
            trafficUsedBytes: number | null;
            notifyPercent: number | null;
            usersOnline: number | null;
            consumptionMultiplier: number;
            cpuCount: number | null;
            cpuModel: string | null;
            totalRam: string | null;
            configProfile: {
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
            };
            providerUuid: string | null;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=create.command.d.ts.map