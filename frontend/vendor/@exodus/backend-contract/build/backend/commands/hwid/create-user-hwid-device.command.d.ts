import { z } from 'zod';
export declare namespace CreateUserHwidDeviceCommand {
    const url: "/api/hwid/devices";
    const TSQ_url: "/api/hwid/devices";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        hwid: z.ZodString;
        userUuid: z.ZodString;
        platform: z.ZodOptional<z.ZodString>;
        osVersion: z.ZodOptional<z.ZodString>;
        deviceModel: z.ZodOptional<z.ZodString>;
        userAgent: z.ZodOptional<z.ZodString>;
    }, "strip", z.ZodTypeAny, {
        hwid: string;
        userUuid: string;
        platform?: string | undefined;
        osVersion?: string | undefined;
        deviceModel?: string | undefined;
        userAgent?: string | undefined;
    }, {
        hwid: string;
        userUuid: string;
        platform?: string | undefined;
        osVersion?: string | undefined;
        deviceModel?: string | undefined;
        userAgent?: string | undefined;
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
//# sourceMappingURL=create-user-hwid-device.command.d.ts.map