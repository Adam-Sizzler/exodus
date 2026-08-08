import { z } from 'zod';
export declare const ConfigProfileInboundsSchema: z.ZodObject<{
    uuid: z.ZodUUID;
    profileUuid: z.ZodUUID;
    tag: z.ZodString;
    type: z.ZodString;
    network: z.ZodNullable<z.ZodString>;
    security: z.ZodNullable<z.ZodString>;
    port: z.ZodNullable<z.ZodNumber>;
    rawInbound: z.ZodNullable<z.ZodUnknown>;
}, z.core.$strip>;
//# sourceMappingURL=config-profile-inbounds.schema.d.ts.map