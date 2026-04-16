export declare const BANDWIDTH_STATS_CONTROLLER: "bandwidth-stats";
export declare const BANDWIDTH_STATS_NODES_ROUTE: "nodes";
export declare const BANDWIDTH_STATS_USERS_ROUTE: "users";
export declare const BANDWIDTH_STATS_NODES_CONTROLLER: "bandwidth-stats/nodes";
export declare const BANDWIDTH_STATS_USERS_CONTROLLER: "bandwidth-stats/users";
export declare const BANDWIDTH_STATS_ROUTES: {
    readonly NODES: {
        readonly GET: "";
        readonly GET_REALTIME: "realtime";
        readonly GET_USERS: (uuid: string) => string;
    };
    readonly USERS: {
        readonly GET_BY_UUID: (uuid: string) => string;
    };
    readonly LEGACY: {
        readonly NODES: {
            readonly GET_USERS: (uuid: string) => string;
        };
        readonly USERS: {
            readonly GET_BY_UUID: (uuid: string) => string;
        };
    };
};
//# sourceMappingURL=bandwidth-stats.d.ts.map