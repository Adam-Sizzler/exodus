export declare const HWID_CONTROLLER: "hwid";
export declare const HWID_ROUTES: {
    readonly GET_ALL_HWID_DEVICES: "devices";
    readonly CREATE_USER_HWID_DEVICE: "devices";
    readonly GET_USER_HWID_DEVICES: (userUuid: string) => string;
    readonly DELETE_USER_HWID_DEVICE: "devices/delete";
    readonly DELETE_ALL_USER_HWID_DEVICES: "devices/delete-all";
    readonly STATS: "devices/stats";
    readonly TOP_USERS_BY_DEVICES: "devices/top-users";
};
//# sourceMappingURL=hwid.d.ts.map