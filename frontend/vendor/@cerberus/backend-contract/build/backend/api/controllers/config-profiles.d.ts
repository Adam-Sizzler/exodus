export declare const CONFIG_PROFILES_CONTROLLER: "config-profiles";
export declare const CONFIG_PROFILES_ROUTES: {
    readonly GET: "";
    readonly CREATE: "";
    readonly UPDATE: "";
    readonly GET_BY_UUID: (uuid: string) => string;
    readonly DELETE: (uuid: string) => string;
    readonly GET_INBOUNDS_BY_PROFILE_UUID: (uuid: string) => string;
    readonly GET_COMPUTED_CONFIG_BY_PROFILE_UUID: (uuid: string) => string;
    readonly GET_ALL_INBOUNDS: "inbounds";
    readonly ACTIONS: {
        readonly REORDER: "actions/reorder";
    };
};
//# sourceMappingURL=config-profiles.d.ts.map