export declare const EXTERNAL_SQUADS_CONTROLLER: "external-squads";
export declare const EXTERNAL_SQUADS_ROUTES: {
    readonly GET: "";
    readonly CREATE: "";
    readonly UPDATE: "";
    readonly GET_BY_UUID: (uuid: string) => string;
    readonly DELETE: (uuid: string) => string;
    readonly BULK_ACTIONS: {
        readonly ADD_USERS: (uuid: string) => string;
        readonly REMOVE_USERS: (uuid: string) => string;
    };
    readonly ACTIONS: {
        readonly REORDER: "actions/reorder";
    };
};
//# sourceMappingURL=external-squads.d.ts.map