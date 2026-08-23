import { z } from 'zod';
export declare const SharedListConfigSchema: z.ZodDiscriminatedUnion<[z.ZodObject<{
    type: z.ZodLiteral<"ipList">;
    items: z.ZodArray<z.ZodUnion<readonly [z.ZodCIDRv4, z.ZodString, z.ZodUnion<readonly [z.ZodIPv4, z.ZodString]>]>>;
}, z.core.$strip>, z.ZodObject<{
    type: z.ZodLiteral<"asList">;
    items: z.ZodArray<z.ZodInt>;
}, z.core.$strip>], "type">;
export declare const SharedListSchema: z.ZodDiscriminatedUnion<[z.ZodObject<{
    type: z.ZodLiteral<"ipList">;
    items: z.ZodArray<z.ZodUnion<readonly [z.ZodCIDRv4, z.ZodString, z.ZodUnion<readonly [z.ZodIPv4, z.ZodString]>]>>;
    name: z.ZodString;
}, z.core.$strip>, z.ZodObject<{
    type: z.ZodLiteral<"asList">;
    items: z.ZodArray<z.ZodInt>;
    name: z.ZodString;
}, z.core.$strip>], "type">;
export declare const IngressFilterPluginSchema: z.ZodObject<{
    enabled: z.ZodBoolean;
    blockedIps: z.ZodArray<z.ZodUnion<readonly [z.ZodUnion<readonly [z.ZodCIDRv4, z.ZodString]>, z.ZodUnion<readonly [z.ZodIPv4, z.ZodString]>, z.ZodString]>>;
}, z.core.$strip>;
export declare const EgressFilterPluginSchema: z.ZodObject<{
    enabled: z.ZodBoolean;
    blockedIps: z.ZodOptional<z.ZodArray<z.ZodUnion<readonly [z.ZodUnion<readonly [z.ZodCIDRv4, z.ZodString]>, z.ZodUnion<readonly [z.ZodIPv4, z.ZodString]>, z.ZodString]>>>;
    blockedPorts: z.ZodOptional<z.ZodArray<z.ZodInt>>;
}, z.core.$strip>;
export declare const preStartPluginSchema: z.ZodObject<{
    enabled: z.ZodDefault<z.ZodBoolean>;
    cleanupSockets: z.ZodOptional<z.ZodObject<{
        enabled: z.ZodBoolean;
        files: z.ZodArray<z.ZodString>;
    }, z.core.$strip>>;
}, z.core.$strip>;
export declare const HaproxyAuthPluginSchema: z.ZodObject<{
    enabled: z.ZodDefault<z.ZodBoolean>;
    inboundTags: z.ZodDefault<z.ZodOptional<z.ZodArray<z.ZodString>>>;
}, z.core.$strip>;
export declare const NodePluginSchema: z.ZodObject<{
    sharedLists: z.ZodDefault<z.ZodOptional<z.ZodArray<z.ZodDiscriminatedUnion<[z.ZodObject<{
        type: z.ZodLiteral<"ipList">;
        items: z.ZodArray<z.ZodUnion<readonly [z.ZodCIDRv4, z.ZodString, z.ZodUnion<readonly [z.ZodIPv4, z.ZodString]>]>>;
        name: z.ZodString;
    }, z.core.$strip>, z.ZodObject<{
        type: z.ZodLiteral<"asList">;
        items: z.ZodArray<z.ZodInt>;
        name: z.ZodString;
    }, z.core.$strip>], "type">>>>;
    ingressFilter: z.ZodOptional<z.ZodObject<{
        enabled: z.ZodBoolean;
        blockedIps: z.ZodArray<z.ZodUnion<readonly [z.ZodUnion<readonly [z.ZodCIDRv4, z.ZodString]>, z.ZodUnion<readonly [z.ZodIPv4, z.ZodString]>, z.ZodString]>>;
    }, z.core.$strip>>;
    egressFilter: z.ZodOptional<z.ZodObject<{
        enabled: z.ZodBoolean;
        blockedIps: z.ZodOptional<z.ZodArray<z.ZodUnion<readonly [z.ZodUnion<readonly [z.ZodCIDRv4, z.ZodString]>, z.ZodUnion<readonly [z.ZodIPv4, z.ZodString]>, z.ZodString]>>>;
        blockedPorts: z.ZodOptional<z.ZodArray<z.ZodInt>>;
    }, z.core.$strip>>;
    haproxyAuth: z.ZodOptional<z.ZodObject<{
        enabled: z.ZodDefault<z.ZodBoolean>;
        inboundTags: z.ZodDefault<z.ZodOptional<z.ZodArray<z.ZodString>>>;
    }, z.core.$strip>>;
    preStart: z.ZodOptional<z.ZodObject<{
        enabled: z.ZodDefault<z.ZodBoolean>;
        cleanupSockets: z.ZodOptional<z.ZodObject<{
            enabled: z.ZodBoolean;
            files: z.ZodArray<z.ZodString>;
        }, z.core.$strip>>;
    }, z.core.$strip>>;
}, z.core.$strip>;
export declare const NodePluginEditorSchema: z.ZodObject<{
    ingressFilter: z.ZodOptional<z.ZodObject<{
        enabled: z.ZodBoolean;
        blockedIps: z.ZodArray<z.ZodUnion<readonly [z.ZodUnion<readonly [z.ZodCIDRv4, z.ZodString]>, z.ZodUnion<readonly [z.ZodIPv4, z.ZodString]>, z.ZodString]>>;
    }, z.core.$strip>>;
    egressFilter: z.ZodOptional<z.ZodObject<{
        enabled: z.ZodBoolean;
        blockedIps: z.ZodOptional<z.ZodArray<z.ZodUnion<readonly [z.ZodUnion<readonly [z.ZodCIDRv4, z.ZodString]>, z.ZodUnion<readonly [z.ZodIPv4, z.ZodString]>, z.ZodString]>>>;
        blockedPorts: z.ZodOptional<z.ZodArray<z.ZodInt>>;
    }, z.core.$strip>>;
    haproxyAuth: z.ZodOptional<z.ZodObject<{
        enabled: z.ZodDefault<z.ZodBoolean>;
        inboundTags: z.ZodDefault<z.ZodOptional<z.ZodArray<z.ZodString>>>;
    }, z.core.$strip>>;
    preStart: z.ZodOptional<z.ZodObject<{
        enabled: z.ZodDefault<z.ZodBoolean>;
        cleanupSockets: z.ZodOptional<z.ZodObject<{
            enabled: z.ZodBoolean;
            files: z.ZodArray<z.ZodString>;
        }, z.core.$strip>>;
    }, z.core.$strip>>;
}, z.core.$strip>;
export type TSharedListConfig = z.infer<typeof SharedListConfigSchema>;
export type TNodePlugin = z.infer<typeof NodePluginSchema>;
export type TNodePluginEditor = z.infer<typeof NodePluginEditorSchema>;
//# sourceMappingURL=node-plugins.schema.d.ts.map