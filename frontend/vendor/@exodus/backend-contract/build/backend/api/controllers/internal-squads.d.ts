export declare const INTERNAL_SQUADS_CONTROLLER: "internal-squads";
export declare const INTERNAL_SQUADS_ROUTES: {
    readonly GET: "";
    readonly CREATE: "";
    readonly UPDATE: "";
    readonly GET_BY_UUID: (uuid: string) => string;
    readonly DELETE: (uuid: string) => string;
    readonly ACCESSIBLE_NODES: (uuid: string) => string;
    readonly BULK_ACTIONS: {
        readonly ADD_USERS: (uuid: string) => string;
        readonly REMOVE_USERS: (uuid: string) => string;
    };
    readonly ACTIONS: {
        readonly REORDER: "actions/reorder";
    };
};
//# sourceMappingURL=internal-squads.d.ts.map