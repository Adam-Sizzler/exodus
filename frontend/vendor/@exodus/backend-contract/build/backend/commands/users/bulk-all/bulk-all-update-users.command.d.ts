import { z } from 'zod';
export declare namespace BulkAllUpdateUsersCommand {
    const url: "/api/users/bulk/all/update";
    const TSQ_url: "/api/users/bulk/all/update";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
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
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
}
//# sourceMappingURL=bulk-all-update-users.command.d.ts.map