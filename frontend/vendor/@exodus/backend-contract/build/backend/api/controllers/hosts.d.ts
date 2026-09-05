export declare const HOSTS_CONTROLLER: "hosts";
export declare const HOST_ACTIONS_ROUTE: "actions";
export declare const HOSTS_ROUTES: {
    readonly CREATE: "";
    readonly UPDATE: "";
    readonly GET: "";
    readonly GET_BY_UUID: (uuid: string) => string;
    readonly DELETE: (uuid: string) => string;
    readonly ACTIONS: {
        readonly REORDER: "actions/reorder";
        readonly CLONE: "actions/clone";
    };
    readonly BULK: {
        readonly ENABLE_HOSTS: "bulk/enable";
        readonly DISABLE_HOSTS: "bulk/disable";
        readonly DELETE_HOSTS: "bulk/delete";
        readonly UPDATE: "bulk/update";
    };
    readonly TAGS: {
        readonly GET: "tags";
    };
};
//# sourceMappingURL=hosts.d.ts.map