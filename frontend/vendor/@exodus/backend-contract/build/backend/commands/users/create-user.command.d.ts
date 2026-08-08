import { z } from 'zod';
export declare namespace CreateUserCommand {
    const url: "/api/users/";
    const TSQ_url: "/api/users/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        username: z.ZodString;
        status: z.ZodOptional<z.ZodDefault<z.ZodEnum<{
            readonly ACTIVE: "ACTIVE";
            readonly DISABLED: "DISABLED";
            readonly LIMITED: "LIMITED";
            readonly EXPIRED: "EXPIRED";
        }>>>;
        shortUuid: z.ZodOptional<z.ZodString>;
        trojanPassword: z.ZodOptional<z.ZodString>;
        vlessUuid: z.ZodOptional<z.ZodUUID>;
        ssPassword: z.ZodOptional<z.ZodString>;
        trafficLimitBytes: z.ZodOptional<z.ZodNumber>;
        trafficLimitStrategy: z.ZodOptional<z.ZodDefault<z.ZodEnum<{
            readonly NO_RESET: "NO_RESET";
            readonly DAY: "DAY";
            readonly WEEK: "WEEK";
            readonly MONTH: "MONTH";
            readonly MONTH_ROLLING: "MONTH_ROLLING";
        }>>>;
        expireAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
        createdAt: z.ZodOptional<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
        lastTrafficResetAt: z.ZodOptional<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
        description: z.ZodOptional<z.ZodString>;
        tag: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        telegramId: z.ZodOptional<z.ZodNullable<z.ZodNumber>>;
        email: z.ZodOptional<z.ZodNullable<z.ZodEmail>>;
        hwidDeviceLimit: z.ZodOptional<z.ZodInt>;
        activeInternalSquads: z.ZodOptional<z.ZodArray<z.ZodUUID>>;
        externalSquadUuid: z.ZodOptional<z.ZodNullable<z.ZodUUID>>;
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
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=create-user.command.d.ts.map