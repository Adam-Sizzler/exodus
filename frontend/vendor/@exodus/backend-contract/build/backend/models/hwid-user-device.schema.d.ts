import { z } from 'zod';
export declare const HwidUserDeviceSchema: z.ZodObject<{
    hwid: z.ZodString;
    userId: z.ZodNumber;
    platform: z.ZodNullable<z.ZodString>;
    osVersion: z.ZodNullable<z.ZodString>;
    deviceModel: z.ZodNullable<z.ZodString>;
    userAgent: z.ZodNullable<z.ZodString>;
    requestIp: z.ZodNullable<z.ZodString>;
    createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
    updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
}, z.core.$strip>;
//# sourceMappingURL=hwid-user-device.schema.d.ts.map