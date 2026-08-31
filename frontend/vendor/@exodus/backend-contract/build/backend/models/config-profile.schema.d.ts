import { z } from 'zod';
export declare const ConfigProfileSchema: z.ZodObject<{
    uuid: z.ZodUUID;
    viewPosition: z.ZodInt;
    name: z.ZodString;
    tags: z.ZodDefault<z.ZodArray<z.ZodString>>;
    config: z.ZodUnknown;
    inbounds: z.ZodArray<z.ZodObject<{
        uuid: z.ZodUUID;
        profileUuid: z.ZodUUID;
        tag: z.ZodString;
        type: z.ZodString;
        network: z.ZodNullable<z.ZodString>;
        security: z.ZodNullable<z.ZodString>;
        port: z.ZodNullable<z.ZodNumber>;
        rawInbound: z.ZodNullable<z.ZodUnknown>;
    }, z.core.$strip>>;
    nodes: z.ZodArray<z.ZodObject<{
        uuid: z.ZodUUID;
        name: z.ZodString;
        countryCode: z.ZodString;
    }, z.core.$strip>>;
    createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
    updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
}, z.core.$strip>;
//# sourceMappingURL=config-profile.schema.d.ts.map