import { z } from 'zod';
export declare namespace BulkUpdateUsersCommand {
    const url: "/api/users/bulk/update";
    const TSQ_url: "/api/users/bulk/update";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        userIds: z.ZodArray<z.ZodNumber>;
        fields: z.ZodObject<{
            status: z.ZodOptional<z.ZodEnum<{
                readonly ACTIVE: "ACTIVE";
                readonly DISABLED: "DISABLED";
                readonly LIMITED: "LIMITED";
                readonly EXPIRED: "EXPIRED";
            }>>;
            trafficLimitBytes: z.ZodOptional<z.ZodNumber>;
            trafficLimitStrategy: z.ZodOptional<z.ZodEnum<{
                readonly NO_RESET: "NO_RESET";
                readonly DAY: "DAY";
                readonly WEEK: "WEEK";
                readonly MONTH: "MONTH";
                readonly MONTH_ROLLING: "MONTH_ROLLING";
            }>>;
            expireAt: z.ZodOptional<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
            description: z.ZodOptional<z.ZodNullable<z.ZodString>>;
            telegramId: z.ZodOptional<z.ZodNullable<z.ZodNumber>>;
            email: z.ZodOptional<z.ZodNullable<z.ZodEmail>>;
            tag: z.ZodOptional<z.ZodNullable<z.ZodString>>;
            hwidDeviceLimit: z.ZodOptional<z.ZodNullable<z.ZodInt>>;
            externalSquadUuid: z.ZodOptional<z.ZodNullable<z.ZodUUID>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
}
//# sourceMappingURL=bulk-update-users.command.d.ts.map