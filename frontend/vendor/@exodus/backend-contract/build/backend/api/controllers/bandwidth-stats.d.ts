export declare const BANDWIDTH_STATS_CONTROLLER: "bandwidth-stats";
export declare const BANDWIDTH_STATS_NODES_ROUTE: "nodes";
export declare const BANDWIDTH_STATS_USERS_ROUTE: "users";
export declare const BANDWIDTH_STATS_INTERNAL_SQUADS_ROUTE: "internal-squads";
export declare const BANDWIDTH_STATS_NODES_CONTROLLER: "bandwidth-stats/nodes";
export declare const BANDWIDTH_STATS_USERS_CONTROLLER: "bandwidth-stats/users";
export declare const BANDWIDTH_STATS_INTERNAL_SQUADS_CONTROLLER: "bandwidth-stats/internal-squads";
export declare const BANDWIDTH_STATS_ROUTES: {
    readonly NODES: {
        readonly GET: "";
        readonly GET_REALTIME: "realtime";
        readonly GET_USERS: (uuid: string) => string;
        readonly GET_USERS_BY_NODES: "users";
        readonly GET_USAGE: "usage";
    };
    readonly USERS: {
        readonly GET_BY_ID: (userId: string) => string;
    };
    readonly INTERNAL_SQUADS: {
        readonly GET_USAGE: (uuid: string) => string;
        readonly USER_USAGE: (squadUuid: string, userId: string) => string;
    };
};
//# sourceMappingURL=bandwidth-stats.d.ts.map