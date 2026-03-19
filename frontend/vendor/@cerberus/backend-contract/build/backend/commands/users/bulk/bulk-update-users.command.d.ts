import { z } from 'zod';
export declare namespace BulkUpdateUsersCommand {
    const url: "/api/users/bulk/update";
    const TSQ_url: "/api/users/bulk/update";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        uuids: z.ZodArray<z.ZodString, "many">;
        fields: z.ZodObject<{
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
            externalSquadUuid: z.ZodOptional<z.ZodNullable<z.ZodString>>;
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
            externalSquadUuid?: string | null | undefined;
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
            externalSquadUuid?: string | null | undefined;
        }>;
    }, "strip", z.ZodTypeAny, {
        uuids: string[];
        fields: {
            status?: "DISABLED" | "LIMITED" | "EXPIRED" | "ACTIVE" | undefined;
            tag?: string | null | undefined;
            trafficLimitBytes?: number | undefined;
            description?: string | null | undefined;
            trafficLimitStrategy?: "MONTH" | "NO_RESET" | "DAY" | "WEEK" | undefined;
            expireAt?: Date | undefined;
            telegramId?: number | null | undefined;
            email?: string | null | undefined;
            hwidDeviceLimit?: number | null | undefined;
            externalSquadUuid?: string | null | undefined;
        };
    }, {
        uuids: string[];
        fields: {
            status?: "DISABLED" | "LIMITED" | "EXPIRED" | "ACTIVE" | undefined;
            tag?: string | null | undefined;
            trafficLimitBytes?: number | undefined;
            description?: string | null | undefined;
            trafficLimitStrategy?: "MONTH" | "NO_RESET" | "DAY" | "WEEK" | undefined;
            expireAt?: string | undefined;
            telegramId?: number | null | undefined;
            email?: string | null | undefined;
            hwidDeviceLimit?: number | null | undefined;
            externalSquadUuid?: string | null | undefined;
        };
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            affectedRows: z.ZodNumber;
        }, "strip", z.ZodTypeAny, {
            affectedRows: number;
        }, {
            affectedRows: number;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            affectedRows: number;
        };
    }, {
        response: {
            affectedRows: number;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=bulk-update-users.command.d.ts.map