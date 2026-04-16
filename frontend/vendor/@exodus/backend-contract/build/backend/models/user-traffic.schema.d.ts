import { z } from 'zod';
export declare const UserTrafficSchema: z.ZodObject<{
    usedTrafficBytes: z.ZodNumber;
    lifetimeUsedTrafficBytes: z.ZodNumber;
    onlineAt: z.ZodNullable<z.ZodEffects<z.ZodString, Date, string>>;
    firstConnectedAt: z.ZodNullable<z.ZodEffects<z.ZodString, Date, string>>;
    lastConnectedNodeUuid: z.ZodNullable<z.ZodString>;
}, "strip", z.ZodTypeAny, {
    usedTrafficBytes: number;
    lifetimeUsedTrafficBytes: number;
    onlineAt: Date | null;
    firstConnectedAt: Date | null;
    lastConnectedNodeUuid: string | null;
}, {
    usedTrafficBytes: number;
    lifetimeUsedTrafficBytes: number;
    onlineAt: string | null;
    firstConnectedAt: string | null;
    lastConnectedNodeUuid: string | null;
}>;
//# sourceMappingURL=user-traffic.schema.d.ts.map