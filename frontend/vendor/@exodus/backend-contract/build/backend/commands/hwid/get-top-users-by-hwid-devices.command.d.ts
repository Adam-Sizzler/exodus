import { z } from 'zod';
export declare namespace GetTopUsersByHwidDevicesCommand {
    const url: "/api/hwid/devices/top-users";
    const TSQ_url: "/api/hwid/devices/top-users";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestQuerySchema: z.ZodObject<{
        start: z.ZodDefault<z.ZodCoercedNumber<unknown>>;
        size: z.ZodDefault<z.ZodCoercedNumber<unknown>>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            users: z.ZodArray<z.ZodObject<{
                id: z.ZodNumber;
                username: z.ZodString;
                devicesCount: z.ZodNumber;
            }, z.core.$strip>>;
            total: z.ZodNumber;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestQuery = z.infer<typeof RequestQuerySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-top-users-by-hwid-devices.command.d.ts.map