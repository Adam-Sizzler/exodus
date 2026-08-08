import { z } from 'zod';
export declare namespace UpdateSubscriptionSettingsCommand {
    const url: "/api/subscription-settings/";
    const TSQ_url: "/api/subscription-settings/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        uuid: z.ZodUUID;
        serveJsonAtBaseSubscription: z.ZodOptional<z.ZodBoolean>;
        isShowCustomRemarks: z.ZodOptional<z.ZodBoolean>;
        customRemarks: z.ZodOptional<z.ZodObject<{
            expiredUsers: z.ZodArray<z.ZodString>;
            limitedUsers: z.ZodArray<z.ZodString>;
            disabledUsers: z.ZodArray<z.ZodString>;
            emptyHosts: z.ZodArray<z.ZodString>;
            HWIDMaxDevicesExceeded: z.ZodArray<z.ZodString>;
            HWIDNotSupported: z.ZodArray<z.ZodString>;
        }, z.core.$strip>>;
        customResponseHeaders: z.ZodOptional<z.ZodRecord<z.ZodString, z.ZodString>>;
        randomizeHosts: z.ZodOptional<z.ZodBoolean>;
        responseRules: z.ZodOptional<z.ZodObject<{
            version: z.ZodEnum<{
                readonly 1: "1";
            }>;
            settings: z.ZodOptional<z.ZodObject<{
                disableSubscriptionAccessByPath: z.ZodOptional<z.ZodBoolean>;
            }, z.core.$strip>>;
            rules: z.ZodArray<z.ZodObject<{
                name: z.ZodString;
                description: z.ZodOptional<z.ZodString>;
                enabled: z.ZodBoolean;
                operator: z.ZodEnum<{
                    readonly AND: "AND";
                    readonly OR: "OR";
                }>;
                conditions: z.ZodArray<z.ZodObject<{
                    headerName: z.ZodString;
                    operator: z.ZodEnum<{
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
                }, z.core.$strip>>;
                responseType: z.ZodEnum<{
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
                    }, z.core.$strip>>>;
                    applyHeadersToEnd: z.ZodOptional<z.ZodBoolean>;
                    subscriptionTemplate: z.ZodOptional<z.ZodString>;
                    ignoreHostXrayJsonTemplate: z.ZodOptional<z.ZodBoolean>;
                    ignoreServeJsonAtBaseSubscription: z.ZodOptional<z.ZodBoolean>;
                    additionalExtendedClientsRegex: z.ZodOptional<z.ZodArray<z.ZodString>>;
                    disableHwidCheck: z.ZodOptional<z.ZodBoolean>;
                    encryption: z.ZodOptional<z.ZodObject<{
                        method: z.ZodEnum<{
                            age1: "age1";
                            age1pq1: "age1pq1";
                        }>;
                        key: z.ZodString;
                    }, z.core.$strip>>;
                    excludeHostsByTags: z.ZodOptional<z.ZodArray<z.ZodString>>;
                }, z.core.$strip>>;
            }, z.core.$strip>>;
        }, z.core.$strip>>;
        hwidSettings: z.ZodOptional<z.ZodObject<{
            enabled: z.ZodBoolean;
            fallbackDeviceLimit: z.ZodNumber;
            maxDevicesAnnounce: z.ZodNullable<z.ZodString>;
        }, z.core.$strip>>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            uuid: z.ZodUUID;
            serveJsonAtBaseSubscription: z.ZodBoolean;
            isShowCustomRemarks: z.ZodBoolean;
            customRemarks: z.ZodObject<{
                expiredUsers: z.ZodArray<z.ZodString>;
                limitedUsers: z.ZodArray<z.ZodString>;
                disabledUsers: z.ZodArray<z.ZodString>;
                emptyHosts: z.ZodArray<z.ZodString>;
                HWIDMaxDevicesExceeded: z.ZodArray<z.ZodString>;
                HWIDNotSupported: z.ZodArray<z.ZodString>;
            }, z.core.$strip>;
            customResponseHeaders: z.ZodNullable<z.ZodRecord<z.ZodString, z.ZodString>>;
            randomizeHosts: z.ZodBoolean;
            responseRules: z.ZodNullable<z.ZodObject<{
                version: z.ZodEnum<{
                    readonly 1: "1";
                }>;
                settings: z.ZodOptional<z.ZodObject<{
                    disableSubscriptionAccessByPath: z.ZodOptional<z.ZodBoolean>;
                }, z.core.$strip>>;
                rules: z.ZodArray<z.ZodObject<{
                    name: z.ZodString;
                    description: z.ZodOptional<z.ZodString>;
                    enabled: z.ZodBoolean;
                    operator: z.ZodEnum<{
                        readonly AND: "AND";
                        readonly OR: "OR";
                    }>;
                    conditions: z.ZodArray<z.ZodObject<{
                        headerName: z.ZodString;
                        operator: z.ZodEnum<{
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
                    }, z.core.$strip>>;
                    responseType: z.ZodEnum<{
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
                        }, z.core.$strip>>>;
                        applyHeadersToEnd: z.ZodOptional<z.ZodBoolean>;
                        subscriptionTemplate: z.ZodOptional<z.ZodString>;
                        ignoreHostXrayJsonTemplate: z.ZodOptional<z.ZodBoolean>;
                        ignoreServeJsonAtBaseSubscription: z.ZodOptional<z.ZodBoolean>;
                        additionalExtendedClientsRegex: z.ZodOptional<z.ZodArray<z.ZodString>>;
                        disableHwidCheck: z.ZodOptional<z.ZodBoolean>;
                        encryption: z.ZodOptional<z.ZodObject<{
                            method: z.ZodEnum<{
                                age1: "age1";
                                age1pq1: "age1pq1";
                            }>;
                            key: z.ZodString;
                        }, z.core.$strip>>;
                        excludeHostsByTags: z.ZodOptional<z.ZodArray<z.ZodString>>;
                    }, z.core.$strip>>;
                }, z.core.$strip>>;
            }, z.core.$strip>>;
            hwidSettings: z.ZodNullable<z.ZodObject<{
                enabled: z.ZodBoolean;
                fallbackDeviceLimit: z.ZodNumber;
                maxDevicesAnnounce: z.ZodNullable<z.ZodString>;
            }, z.core.$strip>>;
            createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=update-subscription-settings.command.d.ts.map