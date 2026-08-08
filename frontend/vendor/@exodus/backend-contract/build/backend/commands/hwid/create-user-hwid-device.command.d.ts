import { z } from 'zod';
export declare namespace CreateUserHwidDeviceCommand {
    const url: "/api/hwid/devices";
    const TSQ_url: "/api/hwid/devices";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        hwid: z.ZodString;
        userId: z.ZodNumber;
        platform: z.ZodOptional<z.ZodString>;
        osVersion: z.ZodOptional<z.ZodString>;
        deviceModel: z.ZodOptional<z.ZodString>;
        userAgent: z.ZodOptional<z.ZodString>;
        requestIp: z.ZodOptional<z.ZodString>;
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
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=create-user-hwid-device.command.d.ts.map