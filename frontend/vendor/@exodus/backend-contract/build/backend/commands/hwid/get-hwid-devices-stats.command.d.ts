import { z } from 'zod';
export declare namespace GetHwidDevicesStatsCommand {
    const url: "/api/hwid/devices/stats";
    const TSQ_url: "/api/hwid/devices/stats";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            byPlatform: z.ZodArray<z.ZodObject<{
                platform: z.ZodString;
                count: z.ZodNumber;
                byApp: z.ZodArray<z.ZodObject<{
                    app: z.ZodString;
                    count: z.ZodNumber;
                }, z.core.$strip>>;
            }, z.core.$strip>>;
            stats: z.ZodObject<{
                totalUniqueDevices: z.ZodNumber;
                totalHwidDevices: z.ZodNumber;
                averageHwidDevicesPerUser: z.ZodNumber;
            }, z.core.$strip>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-hwid-devices-stats.command.d.ts.map