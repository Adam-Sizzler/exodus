import { z } from 'zod';
export declare namespace GetAllHwidDevicesCommand {
    const url: "/api/hwid/devices";
    const TSQ_url: "/api/hwid/devices";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestQuerySchema: z.ZodObject<{
        start: z.ZodDefault<z.ZodNumber>;
        size: z.ZodDefault<z.ZodNumber>;
        filters: z.ZodOptional<z.ZodEffects<z.ZodArray<z.ZodObject<{
            id: z.ZodString;
            value: z.ZodUnknown;
        }, "strip", z.ZodTypeAny, {
            id: string;
            value?: unknown;
        }, {
            id: string;
            value?: unknown;
        }>, "many">, {
            id: string;
            value?: unknown;
        }[], unknown>>;
        filterModes: z.ZodOptional<z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, unknown>>;
        globalFilterMode: z.ZodOptional<z.ZodString>;
        sorting: z.ZodOptional<z.ZodEffects<z.ZodArray<z.ZodObject<{
            id: z.ZodString;
            desc: z.ZodBoolean;
        }, "strip", z.ZodTypeAny, {
            id: string;
            desc: boolean;
        }, {
            id: string;
            desc: boolean;
        }>, "many">, {
            id: string;
            desc: boolean;
        }[], unknown>>;
    }, "strip", z.ZodTypeAny, {
        start: number;
        size: number;
        filters?: {
            id: string;
            value?: unknown;
        }[] | undefined;
        filterModes?: Record<string, string> | undefined;
        globalFilterMode?: string | undefined;
        sorting?: {
            id: string;
            desc: boolean;
        }[] | undefined;
    }, {
        start?: number | undefined;
        size?: number | undefined;
        filters?: unknown;
        filterModes?: unknown;
        globalFilterMode?: string | undefined;
        sorting?: unknown;
    }>;
    type RequestQuery = z.infer<typeof RequestQuerySchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
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
            total: z.ZodNumber;
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
//# sourceMappingURL=get-all-hwid-devices.command.d.ts.map