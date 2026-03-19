import { z } from 'zod';
export declare namespace BulkAllUpdateUsersCommand {
    const url: "/api/users/bulk/all/update";
    const TSQ_url: "/api/users/bulk/all/update";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        status: z.ZodOptional<z.ZodDefault<z.ZodNativeEnum<{
            readonly ACTIVE: "ACTIVE";
            readonly DISABLED: "DISABLED";
            readonly LIMITED: "LIMITED";
            readonly EXPIRED: "EXPIRED";
        }>>>;
        trafficLimitBytes: z.ZodOptional<z.ZodNumber>;
        trafficLimitStrategy: z.ZodOptional<z.ZodNativeEnum<{
            readonly NO_RESET: "NO_RESET";
            readonly DAY: "DAY";
            readonly WEEK: "WEEK";
            readonly MONTH: "MONTH";
        }>>;
        expireAt: z.ZodOptional<z.ZodEffects<z.ZodEffects<z.ZodString, Date, string>, Date, string>>;
        description: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        telegramId: z.ZodOptional<z.ZodNullable<z.ZodNumber>>;
        email: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        tag: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        hwidDeviceLimit: z.ZodOptional<z.ZodNullable<z.ZodNumber>>;
    }, "strip", z.ZodTypeAny, {
        status?: "DISABLED" | "LIMITED" | "EXPIRED" | "ACTIVE" | undefined;
        tag?: string | null | undefined;
        trafficLimitBytes?: number | undefined;
        description?: string | null | undefined;
        trafficLimitStrategy?: "MONTH" | "NO_RESET" | "DAY" | "WEEK" | undefined;
        expireAt?: Date | undefined;
        telegramId?: number | null | undefined;
        email?: string | null | undefined;
        hwidDeviceLimit?: number | null | undefined;
    }, {
        status?: "DISABLED" | "LIMITED" | "EXPIRED" | "ACTIVE" | undefined;
        tag?: string | null | undefined;
        trafficLimitBytes?: number | undefined;
        description?: string | null | undefined;
        trafficLimitStrategy?: "MONTH" | "NO_RESET" | "DAY" | "WEEK" | undefined;
        expireAt?: string | undefined;
        telegramId?: number | null | undefined;
        email?: string | null | undefined;
        hwidDeviceLimit?: number | null | undefined;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            eventSent: z.ZodBoolean;
        }, "strip", z.ZodTypeAny, {
            eventSent: boolean;
        }, {
            eventSent: boolean;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            eventSent: boolean;
        };
    }, {
        response: {
            eventSent: boolean;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=bulk-all-update-users.command.d.ts.map