export declare const USERS_CONTROLLER: "users";
export declare const USERS_ACTIONS_ROUTE: "actions";
export declare const USERS_ROUTES: {
    readonly CREATE: "";
    readonly UPDATE: "";
    readonly GET: "";
    readonly DELETE: (uuid: string) => string;
    readonly GET_BY_UUID: (uuid: string) => string;
    readonly ACCESSIBLE_NODES: (uuid: string) => string;
    readonly SUBSCRIPTION_REQUEST_HISTORY: (uuid: string) => string;
    readonly ACTIONS: {
        readonly ENABLE: (uuid: string) => string;
        readonly DISABLE: (uuid: string) => string;
        readonly RESET_TRAFFIC: (uuid: string) => string;
        readonly REVOKE_SUBSCRIPTION: (uuid: string) => string;
    };
    readonly GET_BY: {
        readonly ID: (id: string) => string;
        readonly SHORT_UUID: (shortUuid: string) => string;
        readonly USERNAME: (username: string) => string;
        readonly SUBSCRIPTION_UUID: (subscriptionUuid: string) => string;
        readonly TELEGRAM_ID: (telegramId: string) => string;
        readonly EMAIL: (email: string) => string;
        readonly TAG: (tag: string) => string;
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
};
//# sourceMappingURL=users.d.ts.map