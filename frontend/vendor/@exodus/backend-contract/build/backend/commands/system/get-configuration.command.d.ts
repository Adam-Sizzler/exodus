import { z } from 'zod';
export declare namespace GetConfigurationCommand {
    const url: "/api/system/configuration";
    const TSQ_url: "/api/system/configuration";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            notifications: z.ZodObject<{
                webhook: z.ZodBoolean;
                bandwidthUsage: z.ZodNullable<z.ZodArray<z.ZodNumber>>;
                notConnectedAfter: z.ZodNullable<z.ZodArray<z.ZodNumber>>;
                expirationNotifications: z.ZodNullable<z.ZodArray<z.ZodNumber>>;
            }, z.core.$strip>;
            service: z.ZodObject<{
                cleanUsageHistory: z.ZodBoolean;
                disableUserUsageRecords: z.ZodBoolean;
                disableSrhRecords: z.ZodBoolean;
                exportToRedisStream: z.ZodBoolean;
            }, z.core.$strip>;
            misc: z.ZodObject<{
                shortUuidLength: z.ZodNumber;
                subPublicDomain: z.ZodString;
                userUsageIgnoreBelowBytes: z.ZodNumber;
            }, z.core.$strip>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-configuration.command.d.ts.map