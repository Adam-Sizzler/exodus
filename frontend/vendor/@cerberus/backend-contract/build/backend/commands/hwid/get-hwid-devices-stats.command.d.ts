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
            }, "strip", z.ZodTypeAny, {
                platform: string;
                count: number;
            }, {
                platform: string;
                count: number;
            }>, "many">;
            byApp: z.ZodArray<z.ZodObject<{
                app: z.ZodString;
                count: z.ZodNumber;
            }, "strip", z.ZodTypeAny, {
                count: number;
                app: string;
            }, {
                count: number;
                app: string;
            }>, "many">;
            stats: z.ZodObject<{
                totalUniqueDevices: z.ZodNumber;
                totalHwidDevices: z.ZodNumber;
                averageHwidDevicesPerUser: z.ZodNumber;
            }, "strip", z.ZodTypeAny, {
                totalUniqueDevices: number;
                totalHwidDevices: number;
                averageHwidDevicesPerUser: number;
            }, {
                totalUniqueDevices: number;
                totalHwidDevices: number;
                averageHwidDevicesPerUser: number;
            }>;
        }, "strip", z.ZodTypeAny, {
            stats: {
                totalUniqueDevices: number;
                totalHwidDevices: number;
                averageHwidDevicesPerUser: number;
            };
            byPlatform: {
                platform: string;
                count: number;
            }[];
            byApp: {
                count: number;
                app: string;
            }[];
        }, {
            stats: {
                totalUniqueDevices: number;
                totalHwidDevices: number;
                averageHwidDevicesPerUser: number;
            };
            byPlatform: {
                platform: string;
                count: number;
            }[];
            byApp: {
                count: number;
                app: string;
            }[];
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            stats: {
                totalUniqueDevices: number;
                totalHwidDevices: number;
                averageHwidDevicesPerUser: number;
            };
            byPlatform: {
                platform: string;
                count: number;
            }[];
            byApp: {
                count: number;
                app: string;
            }[];
        };
    }, {
        response: {
            stats: {
                totalUniqueDevices: number;
                totalHwidDevices: number;
                averageHwidDevicesPerUser: number;
            };
            byPlatform: {
                platform: string;
                count: number;
            }[];
            byApp: {
                count: number;
                app: string;
            }[];
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-hwid-devices-stats.command.d.ts.map