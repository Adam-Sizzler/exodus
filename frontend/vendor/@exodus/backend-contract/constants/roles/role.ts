export const ROLE = {
    ADMIN: 'ADMIN',
    API: 'API',
} as const;

export type TRole = typeof ROLE;
export type TRolesKeys = (typeof ROLE)[keyof typeof ROLE];
export type TRoleTypes = [keyof typeof ROLE][number];
