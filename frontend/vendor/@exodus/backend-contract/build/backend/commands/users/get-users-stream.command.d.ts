import { z } from 'zod';
export declare namespace GetUsersStreamCommand {
    const url: "/api/users/stream";
    const TSQ_url: "/api/users/stream";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestQuerySchema: z.ZodObject<{
        cursor: z.ZodOptional<z.ZodCoercedNumber<unknown>>;
        size: z.ZodDefault<z.ZodOptional<z.ZodCoercedNumber<unknown>>>;
        status: z.ZodOptional<z.ZodEnum<{
            readonly ACTIVE: "ACTIVE";
            readonly DISABLED: "DISABLED";
            readonly LIMITED: "LIMITED";
            readonly EXPIRED: "EXPIRED";
        }>>;
        trafficLimitStrategy: z.ZodOptional<z.ZodEnum<{
            readonly NO_RESET: "NO_RESET";
            readonly DAY: "DAY";
            readonly WEEK: "WEEK";
            readonly MONTH: "MONTH";
            readonly MONTH_ROLLING: "MONTH_ROLLING";
        }>>;
        telegramId: z.ZodOptional<z.ZodPipe<z.ZodPipe<z.ZodString, z.ZodTransform<number, string>>, z.ZodNumber>>;
        email: z.ZodOptional<z.ZodEmail>;
        tag: z.ZodOptional<z.ZodString>;
        externalSquadUuid: z.ZodOptional<z.ZodUUID>;
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
                naivePassword: z.ZodString;
                shadowtlsPassword: z.ZodString;
                hysteria2Password: z.ZodString;
                anytlsPassword: z.ZodString;
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
            nextCursor: z.ZodNullable<z.ZodString>;
            hasMore: z.ZodBoolean;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestQuery = z.infer<typeof RequestQuerySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-users-stream.command.d.ts.map