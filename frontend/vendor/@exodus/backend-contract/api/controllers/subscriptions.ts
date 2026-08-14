export const SUBSCRIPTIONS_CONTROLLER = 'subscriptions' as const;

export const SUBSCRIPTIONS_ROUTES = {
    GET: '',
    GET_BY: {
        USERNAME: (username: string) => `by-username/${username}`,
        ID: (userId: string) => `by-id/${userId}`,
        SHORT_UUID: (shortUuid: string) => `by-short-uuid/${shortUuid}`,
        SHORT_UUID_RAW: (shortUuid: string) => `by-short-uuid/${shortUuid}/raw`,
    },
    SUBPAGE: {
        GET_CONFIG: (shortUuid: string) => `subpage-config/${shortUuid}`,
    },
    GET_CONNECTION_KEYS_BY_USER_ID: (userId: string) => `connection-keys/${userId}`,
} as const;
