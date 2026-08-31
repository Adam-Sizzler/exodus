import { z } from 'zod';
export declare namespace CreateExternalSquadCommand {
    const url: "/api/external-squads/";
    const TSQ_url: "/api/external-squads/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        name: z.ZodString;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            uuid: z.ZodUUID;
            viewPosition: z.ZodInt;
            name: z.ZodString;
            tags: z.ZodDefault<z.ZodArray<z.ZodString>>;
            info: z.ZodObject<{
                membersCount: z.ZodNumber;
            }, z.core.$strip>;
            templates: z.ZodArray<z.ZodObject<{
                templateUuid: z.ZodUUID;
                templateType: z.ZodEnum<{
                    readonly XRAY_JSON: "XRAY_JSON";
                    readonly XRAY_BASE64: "XRAY_BASE64";
                    readonly MIHOMO: "MIHOMO";
                    readonly STASH: "STASH";
                    readonly CLASH: "CLASH";
                    readonly SINGBOX: "SINGBOX";
                }>;
            }, z.core.$strip>>;
            subscriptionSettings: z.ZodNullable<z.ZodObject<{
                serveJsonAtBaseSubscription: z.ZodOptional<z.ZodBoolean>;
                isShowCustomRemarks: z.ZodOptional<z.ZodBoolean>;
                randomizeHosts: z.ZodOptional<z.ZodBoolean>;
            }, z.core.$strip>>;
            hostOverrides: z.ZodNullable<z.ZodObject<{
                serverDescription: z.ZodOptional<z.ZodNullable<z.ZodString>>;
                vlessRouteId: z.ZodOptional<z.ZodNullable<z.ZodInt>>;
            }, z.core.$strip>>;
            responseHeadersAdd: z.ZodRecord<z.ZodString, z.ZodString>;
            responseHeadersRemove: z.ZodArray<z.ZodString>;
            hwidSettings: z.ZodNullable<z.ZodObject<{
                enabled: z.ZodBoolean;
                fallbackDeviceLimit: z.ZodNumber;
                maxDevicesAnnounce: z.ZodNullable<z.ZodString>;
            }, z.core.$strip>>;
            customRemarks: z.ZodNullable<z.ZodObject<{
                expiredUsers: z.ZodArray<z.ZodString>;
                limitedUsers: z.ZodArray<z.ZodString>;
                disabledUsers: z.ZodArray<z.ZodString>;
                emptyHosts: z.ZodArray<z.ZodString>;
                HWIDMaxDevicesExceeded: z.ZodArray<z.ZodString>;
                HWIDNotSupported: z.ZodArray<z.ZodString>;
            }, z.core.$strip>>;
            subpageConfigUuid: z.ZodNullable<z.ZodUUID>;
            createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=create-external-squad.command.d.ts.map