export declare const INFRA_BILLING_CONTROLLER: "infra-billing";
export declare const INFRA_BILLING_ROUTES: {
    readonly GET_PROVIDERS: "providers";
    readonly CREATE_PROVIDER: "providers";
    readonly UPDATE_PROVIDER: "providers";
    readonly DELETE_PROVIDER: (uuid: string) => string;
    readonly GET_PROVIDER_BY_UUID: (uuid: string) => string;
    readonly GET_BILLING_NODES: "nodes";
    readonly CREATE_BILLING_NODE: "nodes";
    readonly UPDATE_BILLING_NODE: "nodes";
    readonly DELETE_BILLING_NODE: (uuid: string) => string;
    readonly GET_BILLING_HISTORY: "history";
    readonly CREATE_BILLING_HISTORY: "history";
    readonly DELETE_BILLING_HISTORY: (uuid: string) => string;
};
//# sourceMappingURL=infra-billing.d.ts.map