export const METADATA_CONTROLLER = 'metadata' as const;

export const METADATA_ROUTES = {
    NODE: {
        GET: (nodeId: string) => `node/${nodeId}`,
        UPSERT: (nodeId: string) => `node/${nodeId}`,
    },
    USER: {
        GET: (userId: string) => `user/${userId}`,
        UPSERT: (userId: string) => `user/${userId}`,
    },
} as const;
