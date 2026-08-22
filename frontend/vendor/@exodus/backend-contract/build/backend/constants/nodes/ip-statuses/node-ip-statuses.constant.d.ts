export declare const NODE_IP_STATUSES: {
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
};
export type TNodeIpStatus = [keyof typeof NODE_IP_STATUSES][number];
export declare const NODE_IP_STATUSES_VALUES: ("INBOUND" | "OUTBOUND" | "MANAGEMENT" | "TRANSIT" | "MONITORING" | "RESERVE" | "BLOCKED" | "FLAGGED" | "DEPRECATED" | "UNKNOWN")[];
export declare const NODE_IP_STATUSES_KEYS: TNodeIpStatus[];
//# sourceMappingURL=node-ip-statuses.constant.d.ts.map