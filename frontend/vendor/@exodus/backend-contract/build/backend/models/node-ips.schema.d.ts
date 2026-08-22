import { z } from 'zod';
export declare const NodeIpSchema: z.ZodObject<{
    ip: z.ZodUnion<readonly [z.ZodIPv4, z.ZodIPv6]>;
    status: z.ZodEnum<{
        readonly INBOUND: "INBOUND";
        readonly OUTBOUND: "OUTBOUND";
        readonly MANAGEMENT: "MANAGEMENT";
        readonly TRANSIT: "TRANSIT";
        readonly MONITORING: "MONITORING";
        readonly RESERVE: "RESERVE";
        readonly BLOCKED: "BLOCKED";
        readonly FLAGGED: "FLAGGED";
        readonly DEPRECATED: "DEPRECATED";
        readonly UNKNOWN: "UNKNOWN";
    }>;
}, z.core.$strip>;
export declare const NodeIpsSchema: z.ZodArray<z.ZodObject<{
    ip: z.ZodUnion<readonly [z.ZodIPv4, z.ZodIPv6]>;
    status: z.ZodEnum<{
        readonly INBOUND: "INBOUND";
        readonly OUTBOUND: "OUTBOUND";
        readonly MANAGEMENT: "MANAGEMENT";
        readonly TRANSIT: "TRANSIT";
        readonly MONITORING: "MONITORING";
        readonly RESERVE: "RESERVE";
        readonly BLOCKED: "BLOCKED";
        readonly FLAGGED: "FLAGGED";
        readonly DEPRECATED: "DEPRECATED";
        readonly UNKNOWN: "UNKNOWN";
    }>;
}, z.core.$strip>>;
export type TNodeIp = z.infer<typeof NodeIpSchema>;
export type TNodeIps = z.infer<typeof NodeIpsSchema>;
//# sourceMappingURL=node-ips.schema.d.ts.map