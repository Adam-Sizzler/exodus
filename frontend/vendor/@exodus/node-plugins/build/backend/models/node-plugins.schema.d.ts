import { z } from 'zod';
export declare const SharedListSchema: z.ZodObject<{
    name: z.ZodString;
    type: z.ZodEnum<["ipList"]>;
    items: z.ZodArray<z.ZodUnion<[z.ZodString, z.ZodString, z.ZodString]>, "many">;
}, "strip", z.ZodTypeAny, {
    type: "ipList";
    name: string;
    items: string[];
}, {
    type: "ipList";
    name: string;
    items: string[];
}>;
export declare const TorrentBlockerPluginSchema: z.ZodObject<{
    enabled: z.ZodBoolean;
    blockDuration: z.ZodNumber;
    ignoreLists: z.ZodObject<{
        ip: z.ZodOptional<z.ZodArray<z.ZodUnion<[z.ZodString, z.ZodString]>, "many">>;
        userId: z.ZodOptional<z.ZodArray<z.ZodNumber, "many">>;
    }, "strip", z.ZodTypeAny, {
        ip?: string[] | undefined;
        userId?: number[] | undefined;
    }, {
        ip?: string[] | undefined;
        userId?: number[] | undefined;
    }>;
    includeRuleTags: z.ZodOptional<z.ZodArray<z.ZodString, "many">>;
}, "strip", z.ZodTypeAny, {
    enabled: boolean;
    blockDuration: number;
    ignoreLists: {
        ip?: string[] | undefined;
        userId?: number[] | undefined;
    };
    includeRuleTags?: string[] | undefined;
}, {
    enabled: boolean;
    blockDuration: number;
    ignoreLists: {
        ip?: string[] | undefined;
        userId?: number[] | undefined;
    };
    includeRuleTags?: string[] | undefined;
}>;
export declare const ConnectionDropPluginSchema: z.ZodObject<{
    enabled: z.ZodBoolean;
    whitelistIps: z.ZodArray<z.ZodUnion<[z.ZodString, z.ZodString]>, "many">;
}, "strip", z.ZodTypeAny, {
    enabled: boolean;
    whitelistIps: string[];
}, {
    enabled: boolean;
    whitelistIps: string[];
}>;
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
export declare const NodePluginSchema: z.ZodObject<{
    sharedLists: z.ZodDefault<z.ZodOptional<z.ZodArray<z.ZodObject<{
        name: z.ZodString;
        type: z.ZodEnum<["ipList"]>;
        items: z.ZodArray<z.ZodUnion<[z.ZodString, z.ZodString, z.ZodString]>, "many">;
    }, "strip", z.ZodTypeAny, {
        type: "ipList";
        name: string;
        items: string[];
    }, {
        type: "ipList";
        name: string;
        items: string[];
    }>, "many">>>;
    torrentBlocker: z.ZodOptional<z.ZodObject<{
        enabled: z.ZodBoolean;
        blockDuration: z.ZodNumber;
        ignoreLists: z.ZodObject<{
            ip: z.ZodOptional<z.ZodArray<z.ZodUnion<[z.ZodString, z.ZodString]>, "many">>;
            userId: z.ZodOptional<z.ZodArray<z.ZodNumber, "many">>;
        }, "strip", z.ZodTypeAny, {
            ip?: string[] | undefined;
            userId?: number[] | undefined;
        }, {
            ip?: string[] | undefined;
            userId?: number[] | undefined;
        }>;
        includeRuleTags: z.ZodOptional<z.ZodArray<z.ZodString, "many">>;
    }, "strip", z.ZodTypeAny, {
        enabled: boolean;
        blockDuration: number;
        ignoreLists: {
            ip?: string[] | undefined;
            userId?: number[] | undefined;
        };
        includeRuleTags?: string[] | undefined;
    }, {
        enabled: boolean;
        blockDuration: number;
        ignoreLists: {
            ip?: string[] | undefined;
            userId?: number[] | undefined;
        };
        includeRuleTags?: string[] | undefined;
    }>>;
    ingressFilter: z.ZodOptional<z.ZodObject<{
        enabled: z.ZodBoolean;
        blockedIps: z.ZodArray<z.ZodUnion<[z.ZodString, z.ZodString, z.ZodString, z.ZodString]>, "many">;
    }, "strip", z.ZodTypeAny, {
        enabled: boolean;
        blockedIps: string[];
    }, {
        enabled: boolean;
        blockedIps: string[];
    }>>;
    egressFilter: z.ZodOptional<z.ZodObject<{
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
    }>>;
    connectionDrop: z.ZodOptional<z.ZodObject<{
        enabled: z.ZodBoolean;
        whitelistIps: z.ZodArray<z.ZodUnion<[z.ZodString, z.ZodString]>, "many">;
    }, "strip", z.ZodTypeAny, {
        enabled: boolean;
        whitelistIps: string[];
    }, {
        enabled: boolean;
        whitelistIps: string[];
    }>>;
}, "strip", z.ZodTypeAny, {
    sharedLists: {
        type: "ipList";
        name: string;
        items: string[];
    }[];
    torrentBlocker?: {
        enabled: boolean;
        blockDuration: number;
        ignoreLists: {
            ip?: string[] | undefined;
            userId?: number[] | undefined;
        };
        includeRuleTags?: string[] | undefined;
    } | undefined;
    ingressFilter?: {
        enabled: boolean;
        blockedIps: string[];
    } | undefined;
    egressFilter?: {
        enabled: boolean;
        blockedIps?: string[] | undefined;
        blockedPorts?: number[] | undefined;
    } | undefined;
    connectionDrop?: {
        enabled: boolean;
        whitelistIps: string[];
    } | undefined;
}, {
    sharedLists?: {
        type: "ipList";
        name: string;
        items: string[];
    }[] | undefined;
    torrentBlocker?: {
        enabled: boolean;
        blockDuration: number;
        ignoreLists: {
            ip?: string[] | undefined;
            userId?: number[] | undefined;
        };
        includeRuleTags?: string[] | undefined;
    } | undefined;
    ingressFilter?: {
        enabled: boolean;
        blockedIps: string[];
    } | undefined;
    egressFilter?: {
        enabled: boolean;
        blockedIps?: string[] | undefined;
        blockedPorts?: number[] | undefined;
    } | undefined;
    connectionDrop?: {
        enabled: boolean;
        whitelistIps: string[];
    } | undefined;
}>;
export type TNodePlugin = z.infer<typeof NodePluginSchema>;
//# sourceMappingURL=node-plugins.schema.d.ts.map