import { z } from 'zod';
export declare const ConfigProfileInboundsSchema: z.ZodObject<{
    uuid: z.ZodString;
    profileUuid: z.ZodString;
    tag: z.ZodString;
    type: z.ZodString;
    network: z.ZodNullable<z.ZodString>;
    security: z.ZodNullable<z.ZodString>;
    port: z.ZodNullable<z.ZodNumber>;
    rawInbound: z.ZodNullable<z.ZodUnknown>;
}, "strip", z.ZodTypeAny, {
    type: string;
    uuid: string;
    profileUuid: string;
    tag: string;
    network: string | null;
    security: string | null;
    port: number | null;
    rawInbound?: unknown;
}, {
    type: string;
    uuid: string;
    profileUuid: string;
    tag: string;
    network: string | null;
    security: string | null;
    port: number | null;
    rawInbound?: unknown;
}>;
//# sourceMappingURL=config-profile-inbounds.schema.d.ts.map