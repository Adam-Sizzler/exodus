import { z } from 'zod';
export declare namespace GetUserHwidDevicesCommand {
    const url: (userId: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestParamSchema: z.ZodObject<{
        userId: z.ZodCoercedNumber<unknown>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            total: z.ZodNumber;
            devices: z.ZodArray<z.ZodObject<{
                hwid: z.ZodString;
                userId: z.ZodNumber;
                platform: z.ZodNullable<z.ZodString>;
                osVersion: z.ZodNullable<z.ZodString>;
                deviceModel: z.ZodNullable<z.ZodString>;
                userAgent: z.ZodNullable<z.ZodString>;
                requestIp: z.ZodNullable<z.ZodString>;
                createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
                updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestParam = z.infer<typeof RequestParamSchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-user-hwid-devices.command.d.ts.map