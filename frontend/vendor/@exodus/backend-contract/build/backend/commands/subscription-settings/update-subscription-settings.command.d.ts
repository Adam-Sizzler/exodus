import { z } from 'zod';
export declare namespace UpdateSubscriptionSettingsCommand {
    const url: "/api/subscription-settings/";
    const TSQ_url: "/api/subscription-settings/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        uuid: z.ZodString;
        profileTitle: z.ZodOptional<z.ZodString>;
        supportLink: z.ZodOptional<z.ZodString>;
        profileUpdateInterval: z.ZodOptional<z.ZodNumber>;
        isProfileWebpageUrlEnabled: z.ZodOptional<z.ZodBoolean>;
        serveJsonAtBaseSubscription: z.ZodOptional<z.ZodBoolean>;
        happAnnounce: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        happRouting: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        isShowCustomRemarks: z.ZodOptional<z.ZodBoolean>;
        customRemarks: z.ZodOptional<z.ZodObject<{
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
        customResponseHeaders: z.ZodOptional<z.ZodRecord<z.ZodString, z.ZodString>>;
        randomizeHosts: z.ZodOptional<z.ZodBoolean>;
        responseRules: z.ZodOptional<z.ZodObject<{
            version: z.ZodNativeEnum<{
                readonly 1: "1";
            }>;
            settings: z.ZodOptional<z.ZodObject<{
                disableSubscriptionAccessByPath: z.ZodOptional<z.ZodBoolean>;
            }, "strip", z.ZodTypeAny, {
                disableSubscriptionAccessByPath?: boolean | undefined;
            }, {
                disableSubscriptionAccessByPath?: boolean | undefined;
            }>>;
            rules: z.ZodArray<z.ZodObject<{
                name: z.ZodString;
                description: z.ZodOptional<z.ZodString>;
                enabled: z.ZodBoolean;
                operator: z.ZodNativeEnum<{
                    readonly AND: "AND";
                    readonly OR: "OR";
                }>;
                conditions: z.ZodArray<z.ZodObject<{
                    headerName: z.ZodString;
                    operator: z.ZodNativeEnum<{
                        readonly EQUALS: "EQUALS";
                        readonly NOT_EQUALS: "NOT_EQUALS";
                        readonly CONTAINS: "CONTAINS";
                        readonly NOT_CONTAINS: "NOT_CONTAINS";
                        readonly STARTS_WITH: "STARTS_WITH";
                        readonly NOT_STARTS_WITH: "NOT_STARTS_WITH";
                        readonly ENDS_WITH: "ENDS_WITH";
                        readonly NOT_ENDS_WITH: "NOT_ENDS_WITH";
                        readonly REGEX: "REGEX";
                        readonly NOT_REGEX: "NOT_REGEX";
                    }>;
                    value: z.ZodString;
                    caseSensitive: z.ZodBoolean;
                }, "strip", z.ZodTypeAny, {
                    value: string;
                    headerName: string;
                    operator: "EQUALS" | "NOT_EQUALS" | "CONTAINS" | "NOT_CONTAINS" | "STARTS_WITH" | "NOT_STARTS_WITH" | "ENDS_WITH" | "NOT_ENDS_WITH" | "REGEX" | "NOT_REGEX";
                    caseSensitive: boolean;
                }, {
                    value: string;
                    headerName: string;
                    operator: "EQUALS" | "NOT_EQUALS" | "CONTAINS" | "NOT_CONTAINS" | "STARTS_WITH" | "NOT_STARTS_WITH" | "ENDS_WITH" | "NOT_ENDS_WITH" | "REGEX" | "NOT_REGEX";
                    caseSensitive: boolean;
                }>, "many">;
                responseType: z.ZodNativeEnum<{
                    readonly BROWSER: "BROWSER";
                    readonly BLOCK: "BLOCK";
                    readonly STATUS_CODE_404: "STATUS_CODE_404";
                    readonly STATUS_CODE_451: "STATUS_CODE_451";
                    readonly SOCKET_DROP: "SOCKET_DROP";
                    readonly XRAY_JSON: "XRAY_JSON";
                    readonly XRAY_BASE64: "XRAY_BASE64";
                    readonly MIHOMO: "MIHOMO";
                    readonly STASH: "STASH";
                    readonly CLASH: "CLASH";
                    readonly SINGBOX: "SINGBOX";
                }>;
                responseModifications: z.ZodOptional<z.ZodObject<{
                    headers: z.ZodOptional<z.ZodArray<z.ZodObject<{
                        key: z.ZodString;
                        value: z.ZodString;
                    }, "strip", z.ZodTypeAny, {
                        value: string;
                        key: string;
                    }, {
                        value: string;
                        key: string;
                    }>, "many">>;
                    applyHeadersToEnd: z.ZodOptional<z.ZodOptional<z.ZodBoolean>>;
                    subscriptionTemplate: z.ZodOptional<z.ZodString>;
                    ignoreHostXrayJsonTemplate: z.ZodOptional<z.ZodBoolean>;
                    ignoreServeJsonAtBaseSubscription: z.ZodOptional<z.ZodBoolean>;
                    additionalExtendedClientsRegex: z.ZodOptional<z.ZodArray<z.ZodString, "many">>;
                    disableHwidCheck: z.ZodOptional<z.ZodBoolean>;
                    encryption: z.ZodOptional<z.ZodObject<{
                        method: z.ZodEnum<["age1", "age1pq1"]>;
                        key: z.ZodString;
                    }, "strip", z.ZodTypeAny, {
                        key: string;
                        method: "age1" | "age1pq1";
                    }, {
                        key: string;
                        method: "age1" | "age1pq1";
                    }>>;
                    excludeHostsByTags: z.ZodOptional<z.ZodArray<z.ZodString, "many">>;
                }, "strip", z.ZodTypeAny, {
                    headers?: {
                        value: string;
                        key: string;
                    }[] | undefined;
                    applyHeadersToEnd?: boolean | undefined;
                    subscriptionTemplate?: string | undefined;
                    ignoreHostXrayJsonTemplate?: boolean | undefined;
                    ignoreServeJsonAtBaseSubscription?: boolean | undefined;
                    additionalExtendedClientsRegex?: string[] | undefined;
                    disableHwidCheck?: boolean | undefined;
                    encryption?: {
                        key: string;
                        method: "age1" | "age1pq1";
                    } | undefined;
                    excludeHostsByTags?: string[] | undefined;
                }, {
                    headers?: {
                        value: string;
                        key: string;
                    }[] | undefined;
                    applyHeadersToEnd?: boolean | undefined;
                    subscriptionTemplate?: string | undefined;
                    ignoreHostXrayJsonTemplate?: boolean | undefined;
                    ignoreServeJsonAtBaseSubscription?: boolean | undefined;
                    additionalExtendedClientsRegex?: string[] | undefined;
                    disableHwidCheck?: boolean | undefined;
                    encryption?: {
                        key: string;
                        method: "age1" | "age1pq1";
                    } | undefined;
                    excludeHostsByTags?: string[] | undefined;
                }>>;
            }, "strip", z.ZodTypeAny, {
                name: string;
                enabled: boolean;
                operator: "AND" | "OR";
                conditions: {
                    value: string;
                    headerName: string;
                    operator: "EQUALS" | "NOT_EQUALS" | "CONTAINS" | "NOT_CONTAINS" | "STARTS_WITH" | "NOT_STARTS_WITH" | "ENDS_WITH" | "NOT_ENDS_WITH" | "REGEX" | "NOT_REGEX";
                    caseSensitive: boolean;
                }[];
                responseType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64" | "BROWSER" | "BLOCK" | "STATUS_CODE_404" | "STATUS_CODE_451" | "SOCKET_DROP";
                description?: string | undefined;
                responseModifications?: {
                    headers?: {
                        value: string;
                        key: string;
                    }[] | undefined;
                    applyHeadersToEnd?: boolean | undefined;
                    subscriptionTemplate?: string | undefined;
                    ignoreHostXrayJsonTemplate?: boolean | undefined;
                    ignoreServeJsonAtBaseSubscription?: boolean | undefined;
                    additionalExtendedClientsRegex?: string[] | undefined;
                    disableHwidCheck?: boolean | undefined;
                    encryption?: {
                        key: string;
                        method: "age1" | "age1pq1";
                    } | undefined;
                    excludeHostsByTags?: string[] | undefined;
                } | undefined;
            }, {
                name: string;
                enabled: boolean;
                operator: "AND" | "OR";
                conditions: {
                    value: string;
                    headerName: string;
                    operator: "EQUALS" | "NOT_EQUALS" | "CONTAINS" | "NOT_CONTAINS" | "STARTS_WITH" | "NOT_STARTS_WITH" | "ENDS_WITH" | "NOT_ENDS_WITH" | "REGEX" | "NOT_REGEX";
                    caseSensitive: boolean;
                }[];
                responseType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64" | "BROWSER" | "BLOCK" | "STATUS_CODE_404" | "STATUS_CODE_451" | "SOCKET_DROP";
                description?: string | undefined;
                responseModifications?: {
                    headers?: {
                        value: string;
                        key: string;
                    }[] | undefined;
                    applyHeadersToEnd?: boolean | undefined;
                    subscriptionTemplate?: string | undefined;
                    ignoreHostXrayJsonTemplate?: boolean | undefined;
                    ignoreServeJsonAtBaseSubscription?: boolean | undefined;
                    additionalExtendedClientsRegex?: string[] | undefined;
                    disableHwidCheck?: boolean | undefined;
                    encryption?: {
                        key: string;
                        method: "age1" | "age1pq1";
                    } | undefined;
                    excludeHostsByTags?: string[] | undefined;
                } | undefined;
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            version: "1";
            rules: {
                name: string;
                enabled: boolean;
                operator: "AND" | "OR";
                conditions: {
                    value: string;
                    headerName: string;
                    operator: "EQUALS" | "NOT_EQUALS" | "CONTAINS" | "NOT_CONTAINS" | "STARTS_WITH" | "NOT_STARTS_WITH" | "ENDS_WITH" | "NOT_ENDS_WITH" | "REGEX" | "NOT_REGEX";
                    caseSensitive: boolean;
                }[];
                responseType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64" | "BROWSER" | "BLOCK" | "STATUS_CODE_404" | "STATUS_CODE_451" | "SOCKET_DROP";
                description?: string | undefined;
                responseModifications?: {
                    headers?: {
                        value: string;
                        key: string;
                    }[] | undefined;
                    applyHeadersToEnd?: boolean | undefined;
                    subscriptionTemplate?: string | undefined;
                    ignoreHostXrayJsonTemplate?: boolean | undefined;
                    ignoreServeJsonAtBaseSubscription?: boolean | undefined;
                    additionalExtendedClientsRegex?: string[] | undefined;
                    disableHwidCheck?: boolean | undefined;
                    encryption?: {
                        key: string;
                        method: "age1" | "age1pq1";
                    } | undefined;
                    excludeHostsByTags?: string[] | undefined;
                } | undefined;
            }[];
            settings?: {
                disableSubscriptionAccessByPath?: boolean | undefined;
            } | undefined;
        }, {
            version: "1";
            rules: {
                name: string;
                enabled: boolean;
                operator: "AND" | "OR";
                conditions: {
                    value: string;
                    headerName: string;
                    operator: "EQUALS" | "NOT_EQUALS" | "CONTAINS" | "NOT_CONTAINS" | "STARTS_WITH" | "NOT_STARTS_WITH" | "ENDS_WITH" | "NOT_ENDS_WITH" | "REGEX" | "NOT_REGEX";
                    caseSensitive: boolean;
                }[];
                responseType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64" | "BROWSER" | "BLOCK" | "STATUS_CODE_404" | "STATUS_CODE_451" | "SOCKET_DROP";
                description?: string | undefined;
                responseModifications?: {
                    headers?: {
                        value: string;
                        key: string;
                    }[] | undefined;
                    applyHeadersToEnd?: boolean | undefined;
                    subscriptionTemplate?: string | undefined;
                    ignoreHostXrayJsonTemplate?: boolean | undefined;
                    ignoreServeJsonAtBaseSubscription?: boolean | undefined;
                    additionalExtendedClientsRegex?: string[] | undefined;
                    disableHwidCheck?: boolean | undefined;
                    encryption?: {
                        key: string;
                        method: "age1" | "age1pq1";
                    } | undefined;
                    excludeHostsByTags?: string[] | undefined;
                } | undefined;
            }[];
            settings?: {
                disableSubscriptionAccessByPath?: boolean | undefined;
            } | undefined;
        }>>;
        hwidSettings: z.ZodOptional<z.ZodObject<{
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
    }, "strip", z.ZodTypeAny, {
        uuid: string;
        profileTitle?: string | undefined;
        supportLink?: string | undefined;
        profileUpdateInterval?: number | undefined;
        isProfileWebpageUrlEnabled?: boolean | undefined;
        serveJsonAtBaseSubscription?: boolean | undefined;
        isShowCustomRemarks?: boolean | undefined;
        customRemarks?: {
            expiredUsers: string[];
            limitedUsers: string[];
            disabledUsers: string[];
            emptyHosts: string[];
            HWIDMaxDevicesExceeded: string[];
            HWIDNotSupported: string[];
        } | undefined;
        happAnnounce?: string | null | undefined;
        happRouting?: string | null | undefined;
        customResponseHeaders?: Record<string, string> | undefined;
        randomizeHosts?: boolean | undefined;
        responseRules?: {
            version: "1";
            rules: {
                name: string;
                enabled: boolean;
                operator: "AND" | "OR";
                conditions: {
                    value: string;
                    headerName: string;
                    operator: "EQUALS" | "NOT_EQUALS" | "CONTAINS" | "NOT_CONTAINS" | "STARTS_WITH" | "NOT_STARTS_WITH" | "ENDS_WITH" | "NOT_ENDS_WITH" | "REGEX" | "NOT_REGEX";
                    caseSensitive: boolean;
                }[];
                responseType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64" | "BROWSER" | "BLOCK" | "STATUS_CODE_404" | "STATUS_CODE_451" | "SOCKET_DROP";
                description?: string | undefined;
                responseModifications?: {
                    headers?: {
                        value: string;
                        key: string;
                    }[] | undefined;
                    applyHeadersToEnd?: boolean | undefined;
                    subscriptionTemplate?: string | undefined;
                    ignoreHostXrayJsonTemplate?: boolean | undefined;
                    ignoreServeJsonAtBaseSubscription?: boolean | undefined;
                    additionalExtendedClientsRegex?: string[] | undefined;
                    disableHwidCheck?: boolean | undefined;
                    encryption?: {
                        key: string;
                        method: "age1" | "age1pq1";
                    } | undefined;
                    excludeHostsByTags?: string[] | undefined;
                } | undefined;
            }[];
            settings?: {
                disableSubscriptionAccessByPath?: boolean | undefined;
            } | undefined;
        } | undefined;
        hwidSettings?: {
            enabled: boolean;
            fallbackDeviceLimit: number;
            maxDevicesAnnounce: string | null;
        } | undefined;
    }, {
        uuid: string;
        profileTitle?: string | undefined;
        supportLink?: string | undefined;
        profileUpdateInterval?: number | undefined;
        isProfileWebpageUrlEnabled?: boolean | undefined;
        serveJsonAtBaseSubscription?: boolean | undefined;
        isShowCustomRemarks?: boolean | undefined;
        customRemarks?: {
            expiredUsers: string[];
            limitedUsers: string[];
            disabledUsers: string[];
            emptyHosts: string[];
            HWIDMaxDevicesExceeded: string[];
            HWIDNotSupported: string[];
        } | undefined;
        happAnnounce?: string | null | undefined;
        happRouting?: string | null | undefined;
        customResponseHeaders?: Record<string, string> | undefined;
        randomizeHosts?: boolean | undefined;
        responseRules?: {
            version: "1";
            rules: {
                name: string;
                enabled: boolean;
                operator: "AND" | "OR";
                conditions: {
                    value: string;
                    headerName: string;
                    operator: "EQUALS" | "NOT_EQUALS" | "CONTAINS" | "NOT_CONTAINS" | "STARTS_WITH" | "NOT_STARTS_WITH" | "ENDS_WITH" | "NOT_ENDS_WITH" | "REGEX" | "NOT_REGEX";
                    caseSensitive: boolean;
                }[];
                responseType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64" | "BROWSER" | "BLOCK" | "STATUS_CODE_404" | "STATUS_CODE_451" | "SOCKET_DROP";
                description?: string | undefined;
                responseModifications?: {
                    headers?: {
                        value: string;
                        key: string;
                    }[] | undefined;
                    applyHeadersToEnd?: boolean | undefined;
                    subscriptionTemplate?: string | undefined;
                    ignoreHostXrayJsonTemplate?: boolean | undefined;
                    ignoreServeJsonAtBaseSubscription?: boolean | undefined;
                    additionalExtendedClientsRegex?: string[] | undefined;
                    disableHwidCheck?: boolean | undefined;
                    encryption?: {
                        key: string;
                        method: "age1" | "age1pq1";
                    } | undefined;
                    excludeHostsByTags?: string[] | undefined;
                } | undefined;
            }[];
            settings?: {
                disableSubscriptionAccessByPath?: boolean | undefined;
            } | undefined;
        } | undefined;
        hwidSettings?: {
            enabled: boolean;
            fallbackDeviceLimit: number;
            maxDevicesAnnounce: string | null;
        } | undefined;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            uuid: z.ZodString;
            profileTitle: z.ZodString;
            supportLink: z.ZodString;
            profileUpdateInterval: z.ZodNumber;
            isProfileWebpageUrlEnabled: z.ZodBoolean;
            serveJsonAtBaseSubscription: z.ZodBoolean;
            isShowCustomRemarks: z.ZodBoolean;
            customRemarks: z.ZodObject<{
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
            }>;
            happAnnounce: z.ZodNullable<z.ZodString>;
            happRouting: z.ZodNullable<z.ZodString>;
            customResponseHeaders: z.ZodNullable<z.ZodRecord<z.ZodString, z.ZodString>>;
            randomizeHosts: z.ZodBoolean;
            responseRules: z.ZodNullable<z.ZodObject<{
                version: z.ZodNativeEnum<{
                    readonly 1: "1";
                }>;
                settings: z.ZodOptional<z.ZodObject<{
                    disableSubscriptionAccessByPath: z.ZodOptional<z.ZodBoolean>;
                }, "strip", z.ZodTypeAny, {
                    disableSubscriptionAccessByPath?: boolean | undefined;
                }, {
                    disableSubscriptionAccessByPath?: boolean | undefined;
                }>>;
                rules: z.ZodArray<z.ZodObject<{
                    name: z.ZodString;
                    description: z.ZodOptional<z.ZodString>;
                    enabled: z.ZodBoolean;
                    operator: z.ZodNativeEnum<{
                        readonly AND: "AND";
                        readonly OR: "OR";
                    }>;
                    conditions: z.ZodArray<z.ZodObject<{
                        headerName: z.ZodString;
                        operator: z.ZodNativeEnum<{
                            readonly EQUALS: "EQUALS";
                            readonly NOT_EQUALS: "NOT_EQUALS";
                            readonly CONTAINS: "CONTAINS";
                            readonly NOT_CONTAINS: "NOT_CONTAINS";
                            readonly STARTS_WITH: "STARTS_WITH";
                            readonly NOT_STARTS_WITH: "NOT_STARTS_WITH";
                            readonly ENDS_WITH: "ENDS_WITH";
                            readonly NOT_ENDS_WITH: "NOT_ENDS_WITH";
                            readonly REGEX: "REGEX";
                            readonly NOT_REGEX: "NOT_REGEX";
                        }>;
                        value: z.ZodString;
                        caseSensitive: z.ZodBoolean;
                    }, "strip", z.ZodTypeAny, {
                        value: string;
                        headerName: string;
                        operator: "EQUALS" | "NOT_EQUALS" | "CONTAINS" | "NOT_CONTAINS" | "STARTS_WITH" | "NOT_STARTS_WITH" | "ENDS_WITH" | "NOT_ENDS_WITH" | "REGEX" | "NOT_REGEX";
                        caseSensitive: boolean;
                    }, {
                        value: string;
                        headerName: string;
                        operator: "EQUALS" | "NOT_EQUALS" | "CONTAINS" | "NOT_CONTAINS" | "STARTS_WITH" | "NOT_STARTS_WITH" | "ENDS_WITH" | "NOT_ENDS_WITH" | "REGEX" | "NOT_REGEX";
                        caseSensitive: boolean;
                    }>, "many">;
                    responseType: z.ZodNativeEnum<{
                        readonly BROWSER: "BROWSER";
                        readonly BLOCK: "BLOCK";
                        readonly STATUS_CODE_404: "STATUS_CODE_404";
                        readonly STATUS_CODE_451: "STATUS_CODE_451";
                        readonly SOCKET_DROP: "SOCKET_DROP";
                        readonly XRAY_JSON: "XRAY_JSON";
                        readonly XRAY_BASE64: "XRAY_BASE64";
                        readonly MIHOMO: "MIHOMO";
                        readonly STASH: "STASH";
                        readonly CLASH: "CLASH";
                        readonly SINGBOX: "SINGBOX";
                    }>;
                    responseModifications: z.ZodOptional<z.ZodObject<{
                        headers: z.ZodOptional<z.ZodArray<z.ZodObject<{
                            key: z.ZodString;
                            value: z.ZodString;
                        }, "strip", z.ZodTypeAny, {
                            value: string;
                            key: string;
                        }, {
                            value: string;
                            key: string;
                        }>, "many">>;
                        applyHeadersToEnd: z.ZodOptional<z.ZodOptional<z.ZodBoolean>>;
                        subscriptionTemplate: z.ZodOptional<z.ZodString>;
                        ignoreHostXrayJsonTemplate: z.ZodOptional<z.ZodBoolean>;
                        ignoreServeJsonAtBaseSubscription: z.ZodOptional<z.ZodBoolean>;
                        additionalExtendedClientsRegex: z.ZodOptional<z.ZodArray<z.ZodString, "many">>;
                        disableHwidCheck: z.ZodOptional<z.ZodBoolean>;
                        encryption: z.ZodOptional<z.ZodObject<{
                            method: z.ZodEnum<["age1", "age1pq1"]>;
                            key: z.ZodString;
                        }, "strip", z.ZodTypeAny, {
                            key: string;
                            method: "age1" | "age1pq1";
                        }, {
                            key: string;
                            method: "age1" | "age1pq1";
                        }>>;
                        excludeHostsByTags: z.ZodOptional<z.ZodArray<z.ZodString, "many">>;
                    }, "strip", z.ZodTypeAny, {
                        headers?: {
                            value: string;
                            key: string;
                        }[] | undefined;
                        applyHeadersToEnd?: boolean | undefined;
                        subscriptionTemplate?: string | undefined;
                        ignoreHostXrayJsonTemplate?: boolean | undefined;
                        ignoreServeJsonAtBaseSubscription?: boolean | undefined;
                        additionalExtendedClientsRegex?: string[] | undefined;
                        disableHwidCheck?: boolean | undefined;
                        encryption?: {
                            key: string;
                            method: "age1" | "age1pq1";
                        } | undefined;
                        excludeHostsByTags?: string[] | undefined;
                    }, {
                        headers?: {
                            value: string;
                            key: string;
                        }[] | undefined;
                        applyHeadersToEnd?: boolean | undefined;
                        subscriptionTemplate?: string | undefined;
                        ignoreHostXrayJsonTemplate?: boolean | undefined;
                        ignoreServeJsonAtBaseSubscription?: boolean | undefined;
                        additionalExtendedClientsRegex?: string[] | undefined;
                        disableHwidCheck?: boolean | undefined;
                        encryption?: {
                            key: string;
                            method: "age1" | "age1pq1";
                        } | undefined;
                        excludeHostsByTags?: string[] | undefined;
                    }>>;
                }, "strip", z.ZodTypeAny, {
                    name: string;
                    enabled: boolean;
                    operator: "AND" | "OR";
                    conditions: {
                        value: string;
                        headerName: string;
                        operator: "EQUALS" | "NOT_EQUALS" | "CONTAINS" | "NOT_CONTAINS" | "STARTS_WITH" | "NOT_STARTS_WITH" | "ENDS_WITH" | "NOT_ENDS_WITH" | "REGEX" | "NOT_REGEX";
                        caseSensitive: boolean;
                    }[];
                    responseType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64" | "BROWSER" | "BLOCK" | "STATUS_CODE_404" | "STATUS_CODE_451" | "SOCKET_DROP";
                    description?: string | undefined;
                    responseModifications?: {
                        headers?: {
                            value: string;
                            key: string;
                        }[] | undefined;
                        applyHeadersToEnd?: boolean | undefined;
                        subscriptionTemplate?: string | undefined;
                        ignoreHostXrayJsonTemplate?: boolean | undefined;
                        ignoreServeJsonAtBaseSubscription?: boolean | undefined;
                        additionalExtendedClientsRegex?: string[] | undefined;
                        disableHwidCheck?: boolean | undefined;
                        encryption?: {
                            key: string;
                            method: "age1" | "age1pq1";
                        } | undefined;
                        excludeHostsByTags?: string[] | undefined;
                    } | undefined;
                }, {
                    name: string;
                    enabled: boolean;
                    operator: "AND" | "OR";
                    conditions: {
                        value: string;
                        headerName: string;
                        operator: "EQUALS" | "NOT_EQUALS" | "CONTAINS" | "NOT_CONTAINS" | "STARTS_WITH" | "NOT_STARTS_WITH" | "ENDS_WITH" | "NOT_ENDS_WITH" | "REGEX" | "NOT_REGEX";
                        caseSensitive: boolean;
                    }[];
                    responseType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64" | "BROWSER" | "BLOCK" | "STATUS_CODE_404" | "STATUS_CODE_451" | "SOCKET_DROP";
                    description?: string | undefined;
                    responseModifications?: {
                        headers?: {
                            value: string;
                            key: string;
                        }[] | undefined;
                        applyHeadersToEnd?: boolean | undefined;
                        subscriptionTemplate?: string | undefined;
                        ignoreHostXrayJsonTemplate?: boolean | undefined;
                        ignoreServeJsonAtBaseSubscription?: boolean | undefined;
                        additionalExtendedClientsRegex?: string[] | undefined;
                        disableHwidCheck?: boolean | undefined;
                        encryption?: {
                            key: string;
                            method: "age1" | "age1pq1";
                        } | undefined;
                        excludeHostsByTags?: string[] | undefined;
                    } | undefined;
                }>, "many">;
            }, "strip", z.ZodTypeAny, {
                version: "1";
                rules: {
                    name: string;
                    enabled: boolean;
                    operator: "AND" | "OR";
                    conditions: {
                        value: string;
                        headerName: string;
                        operator: "EQUALS" | "NOT_EQUALS" | "CONTAINS" | "NOT_CONTAINS" | "STARTS_WITH" | "NOT_STARTS_WITH" | "ENDS_WITH" | "NOT_ENDS_WITH" | "REGEX" | "NOT_REGEX";
                        caseSensitive: boolean;
                    }[];
                    responseType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64" | "BROWSER" | "BLOCK" | "STATUS_CODE_404" | "STATUS_CODE_451" | "SOCKET_DROP";
                    description?: string | undefined;
                    responseModifications?: {
                        headers?: {
                            value: string;
                            key: string;
                        }[] | undefined;
                        applyHeadersToEnd?: boolean | undefined;
                        subscriptionTemplate?: string | undefined;
                        ignoreHostXrayJsonTemplate?: boolean | undefined;
                        ignoreServeJsonAtBaseSubscription?: boolean | undefined;
                        additionalExtendedClientsRegex?: string[] | undefined;
                        disableHwidCheck?: boolean | undefined;
                        encryption?: {
                            key: string;
                            method: "age1" | "age1pq1";
                        } | undefined;
                        excludeHostsByTags?: string[] | undefined;
                    } | undefined;
                }[];
                settings?: {
                    disableSubscriptionAccessByPath?: boolean | undefined;
                } | undefined;
            }, {
                version: "1";
                rules: {
                    name: string;
                    enabled: boolean;
                    operator: "AND" | "OR";
                    conditions: {
                        value: string;
                        headerName: string;
                        operator: "EQUALS" | "NOT_EQUALS" | "CONTAINS" | "NOT_CONTAINS" | "STARTS_WITH" | "NOT_STARTS_WITH" | "ENDS_WITH" | "NOT_ENDS_WITH" | "REGEX" | "NOT_REGEX";
                        caseSensitive: boolean;
                    }[];
                    responseType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64" | "BROWSER" | "BLOCK" | "STATUS_CODE_404" | "STATUS_CODE_451" | "SOCKET_DROP";
                    description?: string | undefined;
                    responseModifications?: {
                        headers?: {
                            value: string;
                            key: string;
                        }[] | undefined;
                        applyHeadersToEnd?: boolean | undefined;
                        subscriptionTemplate?: string | undefined;
                        ignoreHostXrayJsonTemplate?: boolean | undefined;
                        ignoreServeJsonAtBaseSubscription?: boolean | undefined;
                        additionalExtendedClientsRegex?: string[] | undefined;
                        disableHwidCheck?: boolean | undefined;
                        encryption?: {
                            key: string;
                            method: "age1" | "age1pq1";
                        } | undefined;
                        excludeHostsByTags?: string[] | undefined;
                    } | undefined;
                }[];
                settings?: {
                    disableSubscriptionAccessByPath?: boolean | undefined;
                } | undefined;
            }>>;
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
            createdAt: z.ZodEffects<z.ZodString, Date, string>;
            updatedAt: z.ZodEffects<z.ZodString, Date, string>;
        }, "strip", z.ZodTypeAny, {
            uuid: string;
            createdAt: Date;
            updatedAt: Date;
            profileTitle: string;
            supportLink: string;
            profileUpdateInterval: number;
            isProfileWebpageUrlEnabled: boolean;
            serveJsonAtBaseSubscription: boolean;
            isShowCustomRemarks: boolean;
            customRemarks: {
                expiredUsers: string[];
                limitedUsers: string[];
                disabledUsers: string[];
                emptyHosts: string[];
                HWIDMaxDevicesExceeded: string[];
                HWIDNotSupported: string[];
            };
            happAnnounce: string | null;
            happRouting: string | null;
            customResponseHeaders: Record<string, string> | null;
            randomizeHosts: boolean;
            responseRules: {
                version: "1";
                rules: {
                    name: string;
                    enabled: boolean;
                    operator: "AND" | "OR";
                    conditions: {
                        value: string;
                        headerName: string;
                        operator: "EQUALS" | "NOT_EQUALS" | "CONTAINS" | "NOT_CONTAINS" | "STARTS_WITH" | "NOT_STARTS_WITH" | "ENDS_WITH" | "NOT_ENDS_WITH" | "REGEX" | "NOT_REGEX";
                        caseSensitive: boolean;
                    }[];
                    responseType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64" | "BROWSER" | "BLOCK" | "STATUS_CODE_404" | "STATUS_CODE_451" | "SOCKET_DROP";
                    description?: string | undefined;
                    responseModifications?: {
                        headers?: {
                            value: string;
                            key: string;
                        }[] | undefined;
                        applyHeadersToEnd?: boolean | undefined;
                        subscriptionTemplate?: string | undefined;
                        ignoreHostXrayJsonTemplate?: boolean | undefined;
                        ignoreServeJsonAtBaseSubscription?: boolean | undefined;
                        additionalExtendedClientsRegex?: string[] | undefined;
                        disableHwidCheck?: boolean | undefined;
                        encryption?: {
                            key: string;
                            method: "age1" | "age1pq1";
                        } | undefined;
                        excludeHostsByTags?: string[] | undefined;
                    } | undefined;
                }[];
                settings?: {
                    disableSubscriptionAccessByPath?: boolean | undefined;
                } | undefined;
            } | null;
            hwidSettings: {
                enabled: boolean;
                fallbackDeviceLimit: number;
                maxDevicesAnnounce: string | null;
            } | null;
        }, {
            uuid: string;
            createdAt: string;
            updatedAt: string;
            profileTitle: string;
            supportLink: string;
            profileUpdateInterval: number;
            isProfileWebpageUrlEnabled: boolean;
            serveJsonAtBaseSubscription: boolean;
            isShowCustomRemarks: boolean;
            customRemarks: {
                expiredUsers: string[];
                limitedUsers: string[];
                disabledUsers: string[];
                emptyHosts: string[];
                HWIDMaxDevicesExceeded: string[];
                HWIDNotSupported: string[];
            };
            happAnnounce: string | null;
            happRouting: string | null;
            customResponseHeaders: Record<string, string> | null;
            randomizeHosts: boolean;
            responseRules: {
                version: "1";
                rules: {
                    name: string;
                    enabled: boolean;
                    operator: "AND" | "OR";
                    conditions: {
                        value: string;
                        headerName: string;
                        operator: "EQUALS" | "NOT_EQUALS" | "CONTAINS" | "NOT_CONTAINS" | "STARTS_WITH" | "NOT_STARTS_WITH" | "ENDS_WITH" | "NOT_ENDS_WITH" | "REGEX" | "NOT_REGEX";
                        caseSensitive: boolean;
                    }[];
                    responseType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64" | "BROWSER" | "BLOCK" | "STATUS_CODE_404" | "STATUS_CODE_451" | "SOCKET_DROP";
                    description?: string | undefined;
                    responseModifications?: {
                        headers?: {
                            value: string;
                            key: string;
                        }[] | undefined;
                        applyHeadersToEnd?: boolean | undefined;
                        subscriptionTemplate?: string | undefined;
                        ignoreHostXrayJsonTemplate?: boolean | undefined;
                        ignoreServeJsonAtBaseSubscription?: boolean | undefined;
                        additionalExtendedClientsRegex?: string[] | undefined;
                        disableHwidCheck?: boolean | undefined;
                        encryption?: {
                            key: string;
                            method: "age1" | "age1pq1";
                        } | undefined;
                        excludeHostsByTags?: string[] | undefined;
                    } | undefined;
                }[];
                settings?: {
                    disableSubscriptionAccessByPath?: boolean | undefined;
                } | undefined;
            } | null;
            hwidSettings: {
                enabled: boolean;
                fallbackDeviceLimit: number;
                maxDevicesAnnounce: string | null;
            } | null;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            uuid: string;
            createdAt: Date;
            updatedAt: Date;
            profileTitle: string;
            supportLink: string;
            profileUpdateInterval: number;
            isProfileWebpageUrlEnabled: boolean;
            serveJsonAtBaseSubscription: boolean;
            isShowCustomRemarks: boolean;
            customRemarks: {
                expiredUsers: string[];
                limitedUsers: string[];
                disabledUsers: string[];
                emptyHosts: string[];
                HWIDMaxDevicesExceeded: string[];
                HWIDNotSupported: string[];
            };
            happAnnounce: string | null;
            happRouting: string | null;
            customResponseHeaders: Record<string, string> | null;
            randomizeHosts: boolean;
            responseRules: {
                version: "1";
                rules: {
                    name: string;
                    enabled: boolean;
                    operator: "AND" | "OR";
                    conditions: {
                        value: string;
                        headerName: string;
                        operator: "EQUALS" | "NOT_EQUALS" | "CONTAINS" | "NOT_CONTAINS" | "STARTS_WITH" | "NOT_STARTS_WITH" | "ENDS_WITH" | "NOT_ENDS_WITH" | "REGEX" | "NOT_REGEX";
                        caseSensitive: boolean;
                    }[];
                    responseType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64" | "BROWSER" | "BLOCK" | "STATUS_CODE_404" | "STATUS_CODE_451" | "SOCKET_DROP";
                    description?: string | undefined;
                    responseModifications?: {
                        headers?: {
                            value: string;
                            key: string;
                        }[] | undefined;
                        applyHeadersToEnd?: boolean | undefined;
                        subscriptionTemplate?: string | undefined;
                        ignoreHostXrayJsonTemplate?: boolean | undefined;
                        ignoreServeJsonAtBaseSubscription?: boolean | undefined;
                        additionalExtendedClientsRegex?: string[] | undefined;
                        disableHwidCheck?: boolean | undefined;
                        encryption?: {
                            key: string;
                            method: "age1" | "age1pq1";
                        } | undefined;
                        excludeHostsByTags?: string[] | undefined;
                    } | undefined;
                }[];
                settings?: {
                    disableSubscriptionAccessByPath?: boolean | undefined;
                } | undefined;
            } | null;
            hwidSettings: {
                enabled: boolean;
                fallbackDeviceLimit: number;
                maxDevicesAnnounce: string | null;
            } | null;
        };
    }, {
        response: {
            uuid: string;
            createdAt: string;
            updatedAt: string;
            profileTitle: string;
            supportLink: string;
            profileUpdateInterval: number;
            isProfileWebpageUrlEnabled: boolean;
            serveJsonAtBaseSubscription: boolean;
            isShowCustomRemarks: boolean;
            customRemarks: {
                expiredUsers: string[];
                limitedUsers: string[];
                disabledUsers: string[];
                emptyHosts: string[];
                HWIDMaxDevicesExceeded: string[];
                HWIDNotSupported: string[];
            };
            happAnnounce: string | null;
            happRouting: string | null;
            customResponseHeaders: Record<string, string> | null;
            randomizeHosts: boolean;
            responseRules: {
                version: "1";
                rules: {
                    name: string;
                    enabled: boolean;
                    operator: "AND" | "OR";
                    conditions: {
                        value: string;
                        headerName: string;
                        operator: "EQUALS" | "NOT_EQUALS" | "CONTAINS" | "NOT_CONTAINS" | "STARTS_WITH" | "NOT_STARTS_WITH" | "ENDS_WITH" | "NOT_ENDS_WITH" | "REGEX" | "NOT_REGEX";
                        caseSensitive: boolean;
                    }[];
                    responseType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64" | "BROWSER" | "BLOCK" | "STATUS_CODE_404" | "STATUS_CODE_451" | "SOCKET_DROP";
                    description?: string | undefined;
                    responseModifications?: {
                        headers?: {
                            value: string;
                            key: string;
                        }[] | undefined;
                        applyHeadersToEnd?: boolean | undefined;
                        subscriptionTemplate?: string | undefined;
                        ignoreHostXrayJsonTemplate?: boolean | undefined;
                        ignoreServeJsonAtBaseSubscription?: boolean | undefined;
                        additionalExtendedClientsRegex?: string[] | undefined;
                        disableHwidCheck?: boolean | undefined;
                        encryption?: {
                            key: string;
                            method: "age1" | "age1pq1";
                        } | undefined;
                        excludeHostsByTags?: string[] | undefined;
                    } | undefined;
                }[];
                settings?: {
                    disableSubscriptionAccessByPath?: boolean | undefined;
                } | undefined;
            } | null;
            hwidSettings: {
                enabled: boolean;
                fallbackDeviceLimit: number;
                maxDevicesAnnounce: string | null;
            } | null;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=update-subscription-settings.command.d.ts.map