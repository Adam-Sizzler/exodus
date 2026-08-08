export declare const USERS_CONTROLLER: "users";
export declare const USERS_ACTIONS_ROUTE: "actions";
export declare const USERS_ROUTES: {
    readonly CREATE: "";
    readonly UPDATE: "";
    readonly GET: "";
    readonly STREAM: "stream";
    readonly DELETE: (userId: string) => string;
    readonly GET_BY_ID: (userId: string) => string;
    readonly ACCESSIBLE_NODES: (userId: string) => string;
    readonly SUBSCRIPTION_REQUEST_HISTORY: (userId: string) => string;
    readonly ACTIONS: {
        readonly ENABLE: (userId: string) => string;
        readonly DISABLE: (userId: string) => string;
        readonly RESET_TRAFFIC: (userId: string) => string;
        readonly REVOKE_SUBSCRIPTION: (userId: string) => string;
        readonly EXTEND_EXPIRATION_DATE: (userId: string) => string;
    };
    readonly GET_BY: {
        readonly SHORT_UUID: (shortUuid: string) => string;
        readonly USERNAME: (username: string) => string;
    };
    readonly BULK: {
        readonly DELETE_BY_STATUS: "bulk/delete-by-status";
        readonly UPDATE: "bulk/update";
        readonly RESET_TRAFFIC: "bulk/reset-traffic";
        readonly REVOKE_SUBSCRIPTION: "bulk/revoke-subscription";
        readonly DELETE: "bulk/delete";
        readonly UPDATE_SQUADS: "bulk/update-squads";
        readonly EXTEND_EXPIRATION_DATE: "bulk/extend-expiration-date";
        readonly ALL: {
            readonly UPDATE: "bulk/all/update";
            readonly RESET_TRAFFIC: "bulk/all/reset-traffic";
            readonly EXTEND_EXPIRATION_DATE: "bulk/all/extend-expiration-date";
        };
    };
    readonly TAGS: {
        readonly GET: "tags";
    };
    readonly RESOLVE: "resolve";
};
//# sourceMappingURL=users.d.ts.map