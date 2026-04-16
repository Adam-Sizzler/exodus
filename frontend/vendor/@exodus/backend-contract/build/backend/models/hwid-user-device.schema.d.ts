import { z } from 'zod';
export declare const HwidUserDeviceSchema: z.ZodObject<{
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
}>;
//# sourceMappingURL=hwid-user-device.schema.d.ts.map