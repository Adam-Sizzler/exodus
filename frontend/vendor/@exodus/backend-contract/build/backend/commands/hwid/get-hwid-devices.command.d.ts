import { z } from 'zod';
export declare namespace GetHwidDevicesCommand {
    const url: "/api/hwid/devices";
    const TSQ_url: "/api/hwid/devices";
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
            total: z.ZodNumber;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestQuery = z.infer<typeof RequestQuerySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-hwid-devices.command.d.ts.map