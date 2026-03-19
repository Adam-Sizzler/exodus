export declare const EVENTS: {
    readonly USER: {
        readonly CREATED: "user.created";
        readonly MODIFIED: "user.modified";
        readonly DELETED: "user.deleted";
        readonly REVOKED: "user.revoked";
        readonly DISABLED: "user.disabled";
        readonly ENABLED: "user.enabled";
        readonly LIMITED: "user.limited";
        readonly EXPIRED: "user.expired";
        readonly TRAFFIC_RESET: "user.traffic_reset";
        readonly EXPIRE_NOTIFY_EXPIRES_IN_72_HOURS: "user.expires_in_72_hours";
        readonly EXPIRE_NOTIFY_EXPIRES_IN_48_HOURS: "user.expires_in_48_hours";
        readonly EXPIRE_NOTIFY_EXPIRES_IN_24_HOURS: "user.expires_in_24_hours";
        readonly EXPIRE_NOTIFY_EXPIRED_24_HOURS_AGO: "user.expired_24_hours_ago";
        readonly FIRST_CONNECTED: "user.first_connected";
        readonly BANDWIDTH_USAGE_THRESHOLD_REACHED: "user.bandwidth_usage_threshold_reached";
        /**
         * Emitted when a user is not connected to the panel for a certain amount of time.
         * Active only when NOT_CONNECTED_USERS_NOTIFICATIONS_ENABLED is true in .env.
         * Note: For performance reasons, user entity's activeInternalSquads and lastConnectedNode fields are always empty.
         * Includes additional meta object with notConnectedAfterHours field.
         */
        readonly NOT_CONNECTED: "user.not_connected";
    };
    /**
     * HWID events are emitted only when HWID Device Limit feature is enabled.
     * Event payload contains user and HWID device entities.
     * Note: Telegram notifications are skipped for this event.
     * Note: For performance reasons, user entity's activeInternalSquads field is always empty.
     * Returns: user and HWID device entities.
     */
    readonly USER_HWID_DEVICES: {
        readonly ADDED: "user_hwid_devices.added";
        readonly DELETED: "user_hwid_devices.deleted";
    };
    readonly NODE: {
        readonly CREATED: "node.created";
        readonly MODIFIED: "node.modified";
        readonly DISABLED: "node.disabled";
        readonly ENABLED: "node.enabled";
        readonly DELETED: "node.deleted";
        readonly CONNECTION_LOST: "node.connection_lost";
        readonly CONNECTION_RESTORED: "node.connection_restored";
        readonly TRAFFIC_NOTIFY: "node.traffic_notify";
    };
    readonly SERVICE: {
        readonly PANEL_STARTED: "service.panel_started";
        readonly LOGIN_ATTEMPT_FAILED: "service.login_attempt_failed";
        readonly LOGIN_ATTEMPT_SUCCESS: "service.login_attempt_success";
    };
    readonly ERRORS: {
        readonly BANDWIDTH_USAGE_THRESHOLD_REACHED_MAX_NOTIFICATIONS: "errors.bandwidth_usage_threshold_reached_max_notifications";
    };
    readonly CRM: {
        readonly INFRA_BILLING_NODE_PAYMENT_IN_7_DAYS: "crm.infra_billing_node_payment_in_7_days";
        readonly INFRA_BILLING_NODE_PAYMENT_IN_48HRS: "crm.infra_billing_node_payment_in_48hrs";
        readonly INFRA_BILLING_NODE_PAYMENT_IN_24HRS: "crm.infra_billing_node_payment_in_24hrs";
        readonly INFRA_BILLING_NODE_PAYMENT_DUE_TODAY: "crm.infra_billing_node_payment_due_today";
        readonly INFRA_BILLING_NODE_PAYMENT_OVERDUE_24HRS: "crm.infra_billing_node_payment_overdue_24hrs";
        readonly INFRA_BILLING_NODE_PAYMENT_OVERDUE_48HRS: "crm.infra_billing_node_payment_overdue_48hrs";
        readonly INFRA_BILLING_NODE_PAYMENT_OVERDUE_7_DAYS: "crm.infra_billing_node_payment_overdue_7_days";
    };
    readonly CATCH_ALL_USER_EVENTS: "user.*";
    readonly CATCH_ALL_USER_HWID_DEVICES_EVENTS: "user_hwid_devices.*";
    readonly CATCH_ALL_NODE_EVENTS: "node.*";
    readonly CATCH_ALL_SERVICE_EVENTS: "service.*";
    readonly CATCH_ALL_ERRORS_EVENTS: "errors.*";
    readonly CATCH_ALL_CRM_EVENTS: "crm.*";
};
export type TNodeEvents = (typeof EVENTS.NODE)[keyof typeof EVENTS.NODE];
export type TUserEvents = (typeof EVENTS.USER)[keyof typeof EVENTS.USER];
export type TServiceEvents = (typeof EVENTS.SERVICE)[keyof typeof EVENTS.SERVICE];
export type TErrorsEvents = (typeof EVENTS.ERRORS)[keyof typeof EVENTS.ERRORS];
export type TCRMEvents = (typeof EVENTS.CRM)[keyof typeof EVENTS.CRM];
export type TUserHwidDevicesEvents = (typeof EVENTS.USER_HWID_DEVICES)[keyof typeof EVENTS.USER_HWID_DEVICES];
export type TAllEvents = TUserEvents | TNodeEvents | TServiceEvents | TErrorsEvents | TCRMEvents | TUserHwidDevicesEvents;
export type TAllEventChannels = 'telegram' | 'webhook';
export declare const EVENTS_SCOPES: {
    readonly USER: "user";
    readonly USER_HWID_DEVICES: "user_hwid_devices";
    readonly NODE: "node";
    readonly SERVICE: "service";
    readonly ERRORS: "errors";
    readonly CRM: "crm";
};
export type TEventsScope = (typeof EVENTS_SCOPES)[keyof typeof EVENTS_SCOPES];
type ObjectValues<T> = T[keyof T];
type NonEmptyArray<T> = [T, ...T[]];
export declare const toZodEnum: <T extends Record<string, string>>(obj: T) => NonEmptyArray<ObjectValues<T>>;
export {};
//# sourceMappingURL=events.d.ts.map