import { z } from 'zod';
export declare const SharedListSchema: z.ZodDiscriminatedUnion<"type", [z.ZodObject<{
    name: z.ZodString;
    type: z.ZodLiteral<"ipList">;
    items: z.ZodArray<z.ZodUnion<[z.ZodString, z.ZodString, z.ZodString]>, "many">;
}, "strip", z.ZodTypeAny, {
    type: "ipList";
    name: string;
    items: string[];
}, {
    type: "ipList";
    name: string;
    items: string[];
}>, z.ZodObject<{
    name: z.ZodString;
    type: z.ZodLiteral<"asList">;
    items: z.ZodArray<z.ZodNumber, "many">;
}, "strip", z.ZodTypeAny, {
    type: "asList";
    name: string;
    items: number[];
}, {
    type: "asList";
    name: string;
    items: number[];
}>]>;
export declare const IngressFilterPluginSchema: z.ZodObject<{
    enabled: z.ZodBoolean;
    blockedIps: z.ZodArray<z.ZodUnion<[z.ZodString, z.ZodString, z.ZodString, z.ZodString]>, "many">;
}, "strip", z.ZodTypeAny, {
    enabled: boolean;
    blockedIps: string[];
}, {
    enabled: boolean;
    blockedIps: string[];
}>;
export declare const EgressFilterPluginSchema: z.ZodObject<{
    enabled: z.ZodBoolean;
    blockedIps: z.ZodOptional<z.ZodArray<z.ZodUnion<[z.ZodString, z.ZodString, z.ZodString, z.ZodString]>, "many">>;
    blockedPorts: z.ZodOptional<z.ZodArray<z.ZodNumber, "many">>;
}, "strip", z.ZodTypeAny, {
    enabled: boolean;
    blockedIps?: string[] | undefined;
    blockedPorts?: number[] | undefined;
}, {
    enabled: boolean;
    blockedIps?: string[] | undefined;
    blockedPorts?: number[] | undefined;
}>;
export declare const HaproxyAuthPluginSchema: z.ZodObject<{
    enabled: z.ZodBoolean;
}, "strip", z.ZodTypeAny, {
    enabled: boolean;
}, {
    enabled: boolean;
}>;
export declare const NodePluginSchema: z.ZodObject<{
    sharedLists: z.ZodDefault<z.ZodOptional<z.ZodArray<typeof SharedListSchema, "many">>>;
    ingressFilter: z.ZodOptional<typeof IngressFilterPluginSchema>;
    egressFilter: z.ZodOptional<typeof EgressFilterPluginSchema>;
    haproxyAuth: z.ZodOptional<typeof HaproxyAuthPluginSchema>;
}, "strip", z.ZodTypeAny, {
    sharedLists: z.infer<typeof SharedListSchema>[];
    ingressFilter?: z.infer<typeof IngressFilterPluginSchema> | undefined;
    egressFilter?: z.infer<typeof EgressFilterPluginSchema> | undefined;
    haproxyAuth?: z.infer<typeof HaproxyAuthPluginSchema> | undefined;
}, {
    sharedLists?: z.infer<typeof SharedListSchema>[] | undefined;
    ingressFilter?: z.infer<typeof IngressFilterPluginSchema> | undefined;
    egressFilter?: z.infer<typeof EgressFilterPluginSchema> | undefined;
    haproxyAuth?: z.infer<typeof HaproxyAuthPluginSchema> | undefined;
}>;
export type TNodePlugin = z.infer<typeof NodePluginSchema>;
//# sourceMappingURL=node-plugins.schema.d.ts.map
