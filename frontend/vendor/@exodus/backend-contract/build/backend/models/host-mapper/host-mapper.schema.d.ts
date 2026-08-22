import { z } from 'zod';
export declare const XrayJsonHostMapperOperationsSchema: z.ZodDiscriminatedUnion<[z.ZodObject<{
    op: z.ZodLiteral<"copy">;
    from: z.ZodString;
    to: z.ZodString;
}, z.core.$strip>, z.ZodObject<{
    op: z.ZodLiteral<"set">;
    value: z.ZodUnion<readonly [z.ZodString, z.ZodNumber, z.ZodBoolean, z.ZodArray<z.ZodJSONSchema>, z.ZodRecord<z.ZodString, z.ZodJSONSchema>]>;
    to: z.ZodString;
}, z.core.$strip>, z.ZodObject<{
    op: z.ZodLiteral<"unset">;
    to: z.ZodString;
}, z.core.$strip>], "op">;
export declare const MihomoHostMapperOperationsSchema: z.ZodDiscriminatedUnion<[z.ZodObject<{
    op: z.ZodLiteral<"copy">;
    from: z.ZodString;
    to: z.ZodString;
}, z.core.$strip>, z.ZodObject<{
    op: z.ZodLiteral<"set">;
    value: z.ZodUnion<readonly [z.ZodString, z.ZodNumber, z.ZodBoolean, z.ZodArray<z.ZodJSONSchema>, z.ZodRecord<z.ZodString, z.ZodJSONSchema>]>;
    to: z.ZodString;
}, z.core.$strip>, z.ZodObject<{
    op: z.ZodLiteral<"unset">;
    to: z.ZodString;
}, z.core.$strip>], "op">;
export declare const SingBoxHostMapperOperationsSchema: z.ZodDiscriminatedUnion<[z.ZodObject<{
    op: z.ZodLiteral<"copy">;
    from: z.ZodString;
    to: z.ZodString;
}, z.core.$strip>, z.ZodObject<{
    op: z.ZodLiteral<"set">;
    value: z.ZodUnion<readonly [z.ZodString, z.ZodNumber, z.ZodBoolean, z.ZodArray<z.ZodJSONSchema>, z.ZodRecord<z.ZodString, z.ZodJSONSchema>]>;
    to: z.ZodString;
}, z.core.$strip>, z.ZodObject<{
    op: z.ZodLiteral<"unset">;
    to: z.ZodString;
}, z.core.$strip>], "op">;
export declare const Base64HostMapperOperationsSchema: z.ZodDiscriminatedUnion<[z.ZodObject<{
    op: z.ZodLiteral<"copy">;
    from: z.ZodString;
    to: z.ZodString;
}, z.core.$strip>, z.ZodObject<{
    op: z.ZodLiteral<"set">;
    value: z.ZodUnion<readonly [z.ZodString, z.ZodNumber, z.ZodBoolean, z.ZodArray<z.ZodJSONSchema>, z.ZodRecord<z.ZodString, z.ZodJSONSchema>]>;
    to: z.ZodString;
}, z.core.$strip>, z.ZodObject<{
    op: z.ZodLiteral<"unset">;
    to: z.ZodString;
}, z.core.$strip>], "op">;
export declare const HostMapperOperationsSchema: z.ZodDiscriminatedUnion<[z.ZodObject<{
    op: z.ZodLiteral<"copy">;
    from: z.ZodString;
    to: z.ZodString;
}, z.core.$strip>, z.ZodObject<{
    op: z.ZodLiteral<"set">;
    value: z.ZodUnion<readonly [z.ZodString, z.ZodNumber, z.ZodBoolean, z.ZodArray<z.ZodJSONSchema>, z.ZodRecord<z.ZodString, z.ZodJSONSchema>]>;
    to: z.ZodString;
}, z.core.$strip>, z.ZodObject<{
    op: z.ZodLiteral<"unset">;
    to: z.ZodString;
}, z.core.$strip>], "op">;
export declare const HostMapperSchema: z.ZodObject<{
    xrayJson: z.ZodOptional<z.ZodArray<z.ZodDiscriminatedUnion<[z.ZodObject<{
        op: z.ZodLiteral<"copy">;
        from: z.ZodString;
        to: z.ZodString;
    }, z.core.$strip>, z.ZodObject<{
        op: z.ZodLiteral<"set">;
        value: z.ZodUnion<readonly [z.ZodString, z.ZodNumber, z.ZodBoolean, z.ZodArray<z.ZodJSONSchema>, z.ZodRecord<z.ZodString, z.ZodJSONSchema>]>;
        to: z.ZodString;
    }, z.core.$strip>, z.ZodObject<{
        op: z.ZodLiteral<"unset">;
        to: z.ZodString;
    }, z.core.$strip>], "op">>>;
    mihomo: z.ZodOptional<z.ZodArray<z.ZodDiscriminatedUnion<[z.ZodObject<{
        op: z.ZodLiteral<"copy">;
        from: z.ZodString;
        to: z.ZodString;
    }, z.core.$strip>, z.ZodObject<{
        op: z.ZodLiteral<"set">;
        value: z.ZodUnion<readonly [z.ZodString, z.ZodNumber, z.ZodBoolean, z.ZodArray<z.ZodJSONSchema>, z.ZodRecord<z.ZodString, z.ZodJSONSchema>]>;
        to: z.ZodString;
    }, z.core.$strip>, z.ZodObject<{
        op: z.ZodLiteral<"unset">;
        to: z.ZodString;
    }, z.core.$strip>], "op">>>;
    base64: z.ZodOptional<z.ZodArray<z.ZodDiscriminatedUnion<[z.ZodObject<{
        op: z.ZodLiteral<"copy">;
        from: z.ZodString;
        to: z.ZodString;
    }, z.core.$strip>, z.ZodObject<{
        op: z.ZodLiteral<"set">;
        value: z.ZodUnion<readonly [z.ZodString, z.ZodNumber, z.ZodBoolean, z.ZodArray<z.ZodJSONSchema>, z.ZodRecord<z.ZodString, z.ZodJSONSchema>]>;
        to: z.ZodString;
    }, z.core.$strip>, z.ZodObject<{
        op: z.ZodLiteral<"unset">;
        to: z.ZodString;
    }, z.core.$strip>], "op">>>;
    singbox: z.ZodOptional<z.ZodArray<z.ZodDiscriminatedUnion<[z.ZodObject<{
        op: z.ZodLiteral<"copy">;
        from: z.ZodString;
        to: z.ZodString;
    }, z.core.$strip>, z.ZodObject<{
        op: z.ZodLiteral<"set">;
        value: z.ZodUnion<readonly [z.ZodString, z.ZodNumber, z.ZodBoolean, z.ZodArray<z.ZodJSONSchema>, z.ZodRecord<z.ZodString, z.ZodJSONSchema>]>;
        to: z.ZodString;
    }, z.core.$strip>, z.ZodObject<{
        op: z.ZodLiteral<"unset">;
        to: z.ZodString;
    }, z.core.$strip>], "op">>>;
}, z.core.$strip>;
export type THostMapper = z.infer<typeof HostMapperSchema>;
export type THostMapperOperation = z.infer<typeof HostMapperOperationsSchema>;
//# sourceMappingURL=host-mapper.schema.d.ts.map