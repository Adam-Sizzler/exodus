export declare const NODES_CONTROLLER: "nodes";
export declare const NODE_ACTIONS_ROUTE: "actions";
export declare const NODES_ROUTES: {
    readonly CREATE: "";
    readonly GET: "";
    readonly GET_BY_UUID: (uuid: string) => string;
    readonly UPDATE: "";
    readonly DELETE: (uuid: string) => string;
    readonly ACTIONS: {
        readonly ENABLE: (uuid: string) => string;
        readonly DISABLE: (uuid: string) => string;
        readonly RESTART: (uuid: string) => string;
        readonly RESET_TRAFFIC: (uuid: string) => string;
        readonly RESTART_ALL: "actions/restart-all";
        readonly REORDER: "actions/reorder";
    };
    readonly BULK_ACTIONS: {
        readonly PROFILE_MODIFICATION: "bulk-actions/profile-modification";
    };
    readonly TAGS: {
        readonly GET: "tags";
    };
};
//# sourceMappingURL=nodes.d.ts.map