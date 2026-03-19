import { z } from 'zod';
export declare namespace GetExternalSquadByUuidCommand {
    const url: (uuid: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        uuid: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        uuid: string;
    }, {
        uuid: string;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            uuid: z.ZodString;
            viewPosition: z.ZodNumber;
            name: z.ZodString;
            info: z.ZodObject<{
                membersCount: z.ZodNumber;
            }, "strip", z.ZodTypeAny, {
                membersCount: number;
            }, {
                membersCount: number;
            }>;
            templates: z.ZodArray<z.ZodObject<{
                templateUuid: z.ZodString;
                templateType: z.ZodNativeEnum<{
                    readonly XRAY_JSON: "XRAY_JSON";
                    readonly XRAY_BASE64: "XRAY_BASE64";
                    readonly MIHOMO: "MIHOMO";
                    readonly STASH: "STASH";
                    readonly CLASH: "CLASH";
                    readonly SINGBOX: "SINGBOX";
                }>;
            }, "strip", z.ZodTypeAny, {
                templateUuid: string;
                templateType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64";
            }, {
                templateUuid: string;
                templateType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64";
            }>, "many">;
            subscriptionSettings: z.ZodNullable<z.ZodObject<{
                profileTitle: z.ZodOptional<z.ZodString>;
                supportLink: z.ZodOptional<z.ZodString>;
                profileUpdateInterval: z.ZodOptional<z.ZodNumber>;
                isProfileWebpageUrlEnabled: z.ZodOptional<z.ZodBoolean>;
                serveJsonAtBaseSubscription: z.ZodOptional<z.ZodBoolean>;
                isShowCustomRemarks: z.ZodOptional<z.ZodBoolean>;
                happAnnounce: z.ZodOptional<z.ZodNullable<z.ZodString>>;
                happRouting: z.ZodOptional<z.ZodNullable<z.ZodString>>;
                randomizeHosts: z.ZodOptional<z.ZodBoolean>;
            }, "strip", z.ZodTypeAny, {
                profileTitle?: string | undefined;
                supportLink?: string | undefined;
                profileUpdateInterval?: number | undefined;
                isProfileWebpageUrlEnabled?: boolean | undefined;
                serveJsonAtBaseSubscription?: boolean | undefined;
                isShowCustomRemarks?: boolean | undefined;
                happAnnounce?: string | null | undefined;
                happRouting?: string | null | undefined;
                randomizeHosts?: boolean | undefined;
            }, {
                profileTitle?: string | undefined;
                supportLink?: string | undefined;
                profileUpdateInterval?: number | undefined;
                isProfileWebpageUrlEnabled?: boolean | undefined;
                serveJsonAtBaseSubscription?: boolean | undefined;
                isShowCustomRemarks?: boolean | undefined;
                happAnnounce?: string | null | undefined;
                happRouting?: string | null | undefined;
                randomizeHosts?: boolean | undefined;
            }>>;
            hostOverrides: z.ZodNullable<z.ZodObject<{
                serverDescription: z.ZodOptional<z.ZodNullable<z.ZodString>>;
                vlessRouteId: z.ZodOptional<z.ZodNullable<z.ZodNumber>>;
            }, "strip", z.ZodTypeAny, {
                serverDescription?: string | null | undefined;
                vlessRouteId?: number | null | undefined;
            }, {
                serverDescription?: string | null | undefined;
                vlessRouteId?: number | null | undefined;
            }>>;
            responseHeaders: z.ZodNullable<z.ZodRecord<z.ZodString, z.ZodString>>;
            hwidSettings: z.ZodNullable<z.ZodObject<{
                enabled: z.ZodBoolean;
                fallbackDeviceLimit: z.ZodNumber;
                maxDevicesAnnounce: z.ZodNullable<z.ZodString>;
            }, "strip", z.ZodTypeAny, {
                enabled: boolean;
                fallbackDeviceLimit: number;
                maxDevicesAnnounce: string | null;
            }, {
                enabled: boolean;
                fallbackDeviceLimit: number;
                maxDevicesAnnounce: string | null;
            }>>;
            customRemarks: z.ZodNullable<z.ZodObject<{
                expiredUsers: z.ZodArray<z.ZodString, "many">;
                limitedUsers: z.ZodArray<z.ZodString, "many">;
                disabledUsers: z.ZodArray<z.ZodString, "many">;
                emptyHosts: z.ZodArray<z.ZodString, "many">;
                HWIDMaxDevicesExceeded: z.ZodArray<z.ZodString, "many">;
                HWIDNotSupported: z.ZodArray<z.ZodString, "many">;
            }, "strip", z.ZodTypeAny, {
                expiredUsers: string[];
                limitedUsers: string[];
                disabledUsers: string[];
                emptyHosts: string[];
                HWIDMaxDevicesExceeded: string[];
                HWIDNotSupported: string[];
            }, {
                expiredUsers: string[];
                limitedUsers: string[];
                disabledUsers: string[];
                emptyHosts: string[];
                HWIDMaxDevicesExceeded: string[];
                HWIDNotSupported: string[];
            }>>;
            subpageConfigUuid: z.ZodNullable<z.ZodString>;
            createdAt: z.ZodEffects<z.ZodString, Date, string>;
            updatedAt: z.ZodEffects<z.ZodString, Date, string>;
        }, "strip", z.ZodTypeAny, {
            uuid: string;
            createdAt: Date;
            updatedAt: Date;
            name: string;
            viewPosition: number;
            customRemarks: {
                expiredUsers: string[];
                limitedUsers: string[];
                disabledUsers: string[];
                emptyHosts: string[];
                HWIDMaxDevicesExceeded: string[];
                HWIDNotSupported: string[];
            } | null;
            hwidSettings: {
                enabled: boolean;
                fallbackDeviceLimit: number;
                maxDevicesAnnounce: string | null;
            } | null;
            info: {
                membersCount: number;
            };
            templates: {
                templateUuid: string;
                templateType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64";
            }[];
            subscriptionSettings: {
                profileTitle?: string | undefined;
                supportLink?: string | undefined;
                profileUpdateInterval?: number | undefined;
                isProfileWebpageUrlEnabled?: boolean | undefined;
                serveJsonAtBaseSubscription?: boolean | undefined;
                isShowCustomRemarks?: boolean | undefined;
                happAnnounce?: string | null | undefined;
                happRouting?: string | null | undefined;
                randomizeHosts?: boolean | undefined;
            } | null;
            hostOverrides: {
                serverDescription?: string | null | undefined;
                vlessRouteId?: number | null | undefined;
            } | null;
            responseHeaders: Record<string, string> | null;
            subpageConfigUuid: string | null;
        }, {
            uuid: string;
            createdAt: string;
            updatedAt: string;
            name: string;
            viewPosition: number;
            customRemarks: {
                expiredUsers: string[];
                limitedUsers: string[];
                disabledUsers: string[];
                emptyHosts: string[];
                HWIDMaxDevicesExceeded: string[];
                HWIDNotSupported: string[];
            } | null;
            hwidSettings: {
                enabled: boolean;
                fallbackDeviceLimit: number;
                maxDevicesAnnounce: string | null;
            } | null;
            info: {
                membersCount: number;
            };
            templates: {
                templateUuid: string;
                templateType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64";
            }[];
            subscriptionSettings: {
                profileTitle?: string | undefined;
                supportLink?: string | undefined;
                profileUpdateInterval?: number | undefined;
                isProfileWebpageUrlEnabled?: boolean | undefined;
                serveJsonAtBaseSubscription?: boolean | undefined;
                isShowCustomRemarks?: boolean | undefined;
                happAnnounce?: string | null | undefined;
                happRouting?: string | null | undefined;
                randomizeHosts?: boolean | undefined;
            } | null;
            hostOverrides: {
                serverDescription?: string | null | undefined;
                vlessRouteId?: number | null | undefined;
            } | null;
            responseHeaders: Record<string, string> | null;
            subpageConfigUuid: string | null;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            uuid: string;
            createdAt: Date;
            updatedAt: Date;
            name: string;
            viewPosition: number;
            customRemarks: {
                expiredUsers: string[];
                limitedUsers: string[];
                disabledUsers: string[];
                emptyHosts: string[];
                HWIDMaxDevicesExceeded: string[];
                HWIDNotSupported: string[];
            } | null;
            hwidSettings: {
                enabled: boolean;
                fallbackDeviceLimit: number;
                maxDevicesAnnounce: string | null;
            } | null;
            info: {
                membersCount: number;
            };
            templates: {
                templateUuid: string;
                templateType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64";
            }[];
            subscriptionSettings: {
                profileTitle?: string | undefined;
                supportLink?: string | undefined;
                profileUpdateInterval?: number | undefined;
                isProfileWebpageUrlEnabled?: boolean | undefined;
                serveJsonAtBaseSubscription?: boolean | undefined;
                isShowCustomRemarks?: boolean | undefined;
                happAnnounce?: string | null | undefined;
                happRouting?: string | null | undefined;
                randomizeHosts?: boolean | undefined;
            } | null;
            hostOverrides: {
                serverDescription?: string | null | undefined;
                vlessRouteId?: number | null | undefined;
            } | null;
            responseHeaders: Record<string, string> | null;
            subpageConfigUuid: string | null;
        };
    }, {
        response: {
            uuid: string;
            createdAt: string;
            updatedAt: string;
            name: string;
            viewPosition: number;
            customRemarks: {
                expiredUsers: string[];
                limitedUsers: string[];
                disabledUsers: string[];
                emptyHosts: string[];
                HWIDMaxDevicesExceeded: string[];
                HWIDNotSupported: string[];
            } | null;
            hwidSettings: {
                enabled: boolean;
                fallbackDeviceLimit: number;
                maxDevicesAnnounce: string | null;
            } | null;
            info: {
                membersCount: number;
            };
            templates: {
                templateUuid: string;
                templateType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64";
            }[];
            subscriptionSettings: {
                profileTitle?: string | undefined;
                supportLink?: string | undefined;
                profileUpdateInterval?: number | undefined;
                isProfileWebpageUrlEnabled?: boolean | undefined;
                serveJsonAtBaseSubscription?: boolean | undefined;
                isShowCustomRemarks?: boolean | undefined;
                happAnnounce?: string | null | undefined;
                happRouting?: string | null | undefined;
                randomizeHosts?: boolean | undefined;
            } | null;
            hostOverrides: {
                serverDescription?: string | null | undefined;
                vlessRouteId?: number | null | undefined;
            } | null;
            responseHeaders: Record<string, string> | null;
            subpageConfigUuid: string | null;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-external-squad-by-uuid.command.d.ts.map