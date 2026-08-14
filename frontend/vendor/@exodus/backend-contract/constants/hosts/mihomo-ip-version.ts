export const MIHOMO_IP_VERSION = {
    DUAL: 'dual',
    IPV4: 'ipv4',
    IPV6: 'ipv6',
    IPV4_PREFER: 'ipv4-prefer',
    IPV6_PREFER: 'ipv6-prefer',
} as const;

export type TMihomoIpVersion = (typeof MIHOMO_IP_VERSION)[keyof typeof MIHOMO_IP_VERSION];
export const MIHOMO_IP_VERSION_VALUES = Object.values(MIHOMO_IP_VERSION);
