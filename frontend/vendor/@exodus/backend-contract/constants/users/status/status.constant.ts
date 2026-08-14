export const USERS_STATUS = {
    ACTIVE: 'ACTIVE',
    DISABLED: 'DISABLED',
    LIMITED: 'LIMITED',
    EXPIRED: 'EXPIRED',
} as const;

export type TUsersStatus = [keyof typeof USERS_STATUS][number];
export const USERS_STATUS_VALUES = Object.values(USERS_STATUS);
