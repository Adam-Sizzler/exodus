import { z } from 'zod';
export declare namespace GetUsersCommand {
    const url: "/api/users/";
    const TSQ_url: "/api/users/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestQuerySchema: z.ZodObject<{
        start: z.ZodDefault<z.ZodCoercedNumber<unknown>>;
        size: z.ZodDefault<z.ZodCoercedNumber<unknown>>;
        filters: z.ZodOptional<z.ZodPreprocess<z.ZodArray<z.ZodObject<{
            id: z.ZodString;
            value: z.ZodUnknown;
        }, z.core.$strip>>>>;
        filterModes: z.ZodOptional<z.ZodPreprocess<z.ZodRecord<z.ZodString, z.ZodString>>>;
        globalFilterMode: z.ZodOptional<z.ZodString>;
        sorting: z.ZodOptional<z.ZodPreprocess<z.ZodArray<z.ZodObject<{
            id: z.ZodString;
            desc: z.ZodBoolean;
        }, z.core.$strip>>>>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            users: z.ZodArray<z.ZodObject<{
                id: z.ZodNumber;
                shortUuid: z.ZodString;
                username: z.ZodString;
                status: z.ZodEnum<{
                    readonly ACTIVE: "ACTIVE";
                    readonly DISABLED: "DISABLED";
                    readonly LIMITED: "LIMITED";
                    readonly EXPIRED: "EXPIRED";
                }>;
                trafficLimitBytes: z.ZodNumber;
                trafficLimitStrategy: z.ZodEnum<{
                    readonly NO_RESET: "NO_RESET";
                    readonly DAY: "DAY";
                    readonly WEEK: "WEEK";
                    readonly MONTH: "MONTH";
                    readonly MONTH_ROLLING: "MONTH_ROLLING";
                }>;
                expireAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
                telegramId: z.ZodNullable<z.ZodNumber>;
                email: z.ZodNullable<z.ZodEmail>;
                description: z.ZodNullable<z.ZodString>;
                tag: z.ZodNullable<z.ZodString>;
                hwidDeviceLimit: z.ZodNullable<z.ZodInt>;
                externalSquadUuid: z.ZodNullable<z.ZodUUID>;
                trojanPassword: z.ZodString;
                vlessUuid: z.ZodUUID;
                ssPassword: z.ZodString;
                lastTriggeredThreshold: z.ZodInt;
                subRevokedAt: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
                lastTrafficResetAt: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
                createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
                updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
                subscriptionUrl: z.ZodString;
                activeInternalSquads: z.ZodArray<z.ZodObject<{
                    uuid: z.ZodUUID;
                    name: z.ZodString;
                }, z.core.$strip>>;
                userTraffic: z.ZodObject<{
                    usedTrafficBytes: z.ZodNumber;
                    lifetimeUsedTrafficBytes: z.ZodNumber;
                    onlineAt: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
                    firstConnectedAt: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
                    lastConnectedNodeUuid: z.ZodNullable<z.ZodUUID>;
                }, z.core.$strip>;
            }, z.core.$strip>>;
            total: z.ZodNumber;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestQuery = z.infer<typeof RequestQuerySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-users.command.d.ts.map