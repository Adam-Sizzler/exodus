export declare const NODE_PLUGINS_CONTROLLER: "node-plugins";
export declare const NODE_PLUGINS_ROUTES: {
    readonly GET_ALL: "";
    readonly GET: (uuid: string) => string;
    readonly UPDATE: "";
    readonly DELETE: (uuid: string) => string;
    readonly CREATE: "";
    readonly ACTIONS: {
        readonly REORDER: "actions/reorder";
        readonly CLONE: "actions/clone";
        readonly SYNC: "actions/sync";
    };
    readonly EXECUTOR: "executor";
    readonly TORRENT_BLOCKER: {
        readonly GET_REPORTS: "torrent-blocker";
        readonly GET_REPORTS_STATS: "torrent-blocker/stats";
        readonly TRUNCATE_REPORTS: "torrent-blocker/truncate";
    };
    readonly SHARED_LISTS: {
        readonly GET_ALL: "shared-lists";
        readonly GET: "shared-lists/by-name";
        readonly CREATE: "shared-lists";
        readonly UPDATE: "shared-lists";
        readonly DELETE: "shared-lists";
        readonly ACTIONS: {
            readonly SYNC: "shared-lists/actions/sync";
        };
    };
    readonly TAGS: {
        readonly GET: "tags";
        readonly SET: "tags";
    };
};
//# sourceMappingURL=node-plugins.d.ts.map