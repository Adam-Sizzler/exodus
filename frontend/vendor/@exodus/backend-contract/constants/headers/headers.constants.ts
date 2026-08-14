export const EXODUS_CLIENT_TYPE_HEADER = 'X-Exodus-Client-Type';

export const EXODUS_CLIENT_TYPE_BROWSER = 'browser';

export const EXODUS_REAL_IP_HEADER = 'x-exodus-real-ip';

export const EXODUS_BYPASS_HTTPS_RESTRCTIONS = {
    'x-forwarded-proto': 'https',
    'x-forwarded-for': '127.0.0.1',
} as const;
