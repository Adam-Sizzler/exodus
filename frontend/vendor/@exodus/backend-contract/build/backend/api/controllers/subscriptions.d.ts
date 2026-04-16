export declare const SUBSCRIPTIONS_CONTROLLER: "subscriptions";
export declare const SUBSCRIPTIONS_ROUTES: {
    readonly GET: "";
    readonly GET_BY: {
        readonly USERNAME: (username: string) => string;
        readonly UUID: (uuid: string) => string;
        readonly SHORT_UUID: (shortUuid: string) => string;
        readonly SHORT_UUID_RAW: (shortUuid: string) => string;
    };
    readonly SUBPAGE: {
        readonly GET_CONFIG: (shortUuid: string) => string;
    };
};
//# sourceMappingURL=subscriptions.d.ts.map