import { z } from 'zod';
export declare const UserTrafficSchema: z.ZodObject<{
    usedTrafficBytes: z.ZodNumber;
    lifetimeUsedTrafficBytes: z.ZodNumber;
    onlineAt: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
    firstConnectedAt: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
    lastConnectedNodeUuid: z.ZodNullable<z.ZodUUID>;
}, z.core.$strip>;
//# sourceMappingURL=user-traffic.schema.d.ts.map