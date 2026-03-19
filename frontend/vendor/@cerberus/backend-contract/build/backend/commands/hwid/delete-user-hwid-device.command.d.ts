import { z } from 'zod';
export declare namespace DeleteUserHwidDeviceCommand {
    const url: "/api/hwid/devices/delete";
    const TSQ_url: "/api/hwid/devices/delete";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        userUuid: z.ZodString;
        hwid: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        hwid: string;
        userUuid: string;
    }, {
        hwid: string;
        userUuid: string;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            total: z.ZodNumber;
            devices: z.ZodArray<z.ZodObject<{
                hwid: z.ZodString;
                userUuid: z.ZodString;
                platform: z.ZodNullable<z.ZodString>;
                osVersion: z.ZodNullable<z.ZodString>;
                deviceModel: z.ZodNullable<z.ZodString>;
                userAgent: z.ZodNullable<z.ZodString>;
                createdAt: z.ZodEffects<z.ZodString, Date, string>;
                updatedAt: z.ZodEffects<z.ZodString, Date, string>;
            }, "strip", z.ZodTypeAny, {
                hwid: string;
                createdAt: Date;
                updatedAt: Date;
                userUuid: string;
                platform: string | null;
                osVersion: string | null;
                deviceModel: string | null;
                userAgent: string | null;
            }, {
                hwid: string;
                createdAt: string;
                updatedAt: string;
                userUuid: string;
                platform: string | null;
                osVersion: string | null;
                deviceModel: string | null;
                userAgent: string | null;
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            devices: {
                hwid: string;
                createdAt: Date;
                updatedAt: Date;
                userUuid: string;
                platform: string | null;
                osVersion: string | null;
                deviceModel: string | null;
                userAgent: string | null;
            }[];
            total: number;
        }, {
            devices: {
                hwid: string;
                createdAt: string;
                updatedAt: string;
                userUuid: string;
                platform: string | null;
                osVersion: string | null;
                deviceModel: string | null;
                userAgent: string | null;
            }[];
            total: number;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            devices: {
                hwid: string;
                createdAt: Date;
                updatedAt: Date;
                userUuid: string;
                platform: string | null;
                osVersion: string | null;
                deviceModel: string | null;
                userAgent: string | null;
            }[];
            total: number;
        };
    }, {
        response: {
            devices: {
                hwid: string;
                createdAt: string;
                updatedAt: string;
                userUuid: string;
                platform: string | null;
                osVersion: string | null;
                deviceModel: string | null;
                userAgent: string | null;
            }[];
            total: number;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=delete-user-hwid-device.command.d.ts.map