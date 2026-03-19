import { z } from 'zod';
export declare namespace GetTopUsersByHwidDevicesCommand {
    const url: "/api/hwid/devices/top-users";
    const TSQ_url: "/api/hwid/devices/top-users";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestQuerySchema: z.ZodObject<{
        start: z.ZodDefault<z.ZodNumber>;
        size: z.ZodDefault<z.ZodNumber>;
    }, "strip", z.ZodTypeAny, {
        start: number;
        size: number;
    }, {
        start?: number | undefined;
        size?: number | undefined;
    }>;
    type RequestQuery = z.infer<typeof RequestQuerySchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            users: z.ZodArray<z.ZodObject<{
                userUuid: z.ZodString;
                id: z.ZodNumber;
                username: z.ZodString;
                devicesCount: z.ZodNumber;
            }, "strip", z.ZodTypeAny, {
                username: string;
                id: number;
                userUuid: string;
                devicesCount: number;
            }, {
                username: string;
                id: number;
                userUuid: string;
                devicesCount: number;
            }>, "many">;
            total: z.ZodNumber;
        }, "strip", z.ZodTypeAny, {
            users: {
                username: string;
                id: number;
                userUuid: string;
                devicesCount: number;
            }[];
            total: number;
        }, {
            users: {
                username: string;
                id: number;
                userUuid: string;
                devicesCount: number;
            }[];
            total: number;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            users: {
                username: string;
                id: number;
                userUuid: string;
                devicesCount: number;
            }[];
            total: number;
        };
    }, {
        response: {
            users: {
                username: string;
                id: number;
                userUuid: string;
                devicesCount: number;
            }[];
            total: number;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-top-users-by-hwid-devices.command.d.ts.map