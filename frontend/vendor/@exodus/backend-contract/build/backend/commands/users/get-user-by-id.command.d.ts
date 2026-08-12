import { z } from 'zod';
export declare namespace GetUserByIdCommand {
    const url: (userId: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestParamSchema: z.ZodObject<{
        userId: z.ZodCoercedNumber<unknown>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
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
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestParam = z.infer<typeof RequestParamSchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-user-by-id.command.d.ts.map