export const USERS_CONTROLLER = 'users' as const;

export const USERS_ACTIONS_ROUTE = 'actions' as const;

export const USERS_ROUTES = {
    CREATE: '',
    UPDATE: '',
    GET: '',
    STREAM: 'stream',
    DELETE: (userId: string) => `${userId}`,
    GET_BY_ID: (userId: string) => `${userId}`,
    ACCESSIBLE_NODES: (userId: string) => `${userId}/accessible-nodes`,
    SUBSCRIPTION_REQUEST_HISTORY: (userId: string) => `${userId}/subscription-request-history`,
    ACTIONS: {
        ENABLE: (userId: string) => `${userId}/${USERS_ACTIONS_ROUTE}/enable`,
        DISABLE: (userId: string) => `${userId}/${USERS_ACTIONS_ROUTE}/disable`,
        RESET_TRAFFIC: (userId: string) => `${userId}/${USERS_ACTIONS_ROUTE}/reset-traffic`,
        REVOKE_SUBSCRIPTION: (userId: string) => `${userId}/${USERS_ACTIONS_ROUTE}/revoke`,
        EXTEND_EXPIRATION_DATE: (userId: string) => `${userId}/${USERS_ACTIONS_ROUTE}/extend`,
    },
    GET_BY: {
        SHORT_UUID: (shortUuid: string) => `by-short-uuid/${shortUuid}`,
        USERNAME: (username: string) => `by-username/${username}`,
    },

    BULK: {
        DELETE_BY_STATUS: 'bulk/delete-by-status',
        UPDATE: 'bulk/update',
        RESET_TRAFFIC: 'bulk/reset-traffic',
        REVOKE_SUBSCRIPTION: 'bulk/revoke-subscription',
        DELETE: 'bulk/delete',
        UPDATE_SQUADS: 'bulk/update-squads',
        EXTEND_EXPIRATION_DATE: 'bulk/extend-expiration-date',
        ALL: {
            UPDATE: 'bulk/all/update',
            RESET_TRAFFIC: 'bulk/all/reset-traffic',
            EXTEND_EXPIRATION_DATE: 'bulk/all/extend-expiration-date',
        },
    },

    TAGS: {
        GET: 'tags',
    },
    RESOLVE: 'resolve',
} as const;
