import { z } from 'zod';

import { SUBSCRIPTION_TEMPLATE_TYPE } from '../constants';
import { ALPN, MIHOMO_IP_VERSION, SECURITY_LAYERS } from '../constants/hosts';

export const HostsSchema = z.object({
    uuid: z.uuid(),
    viewPosition: z.int(),
    remark: z.string(),
    address: z.string(),
    port: z.int(),
    path: z.string().nullable(),
    sni: z.string().nullable(),
    host: z.string().nullable(),
    alpn: z.enum(ALPN).nullable(),
    fingerprint: z.string().nullable(),
    isDisabled: z.boolean(),
    securityLayer: z.enum(SECURITY_LAYERS).default(SECURITY_LAYERS.DEFAULT),
    xhttpExtraParams: z.nullable(z.unknown()),
    muxParams: z.nullable(z.unknown()),
    singboxMuxParams: z.nullable(z.unknown()),
    clashMuxParams: z.string().nullable(),
    singboxCustomParams: z.nullable(z.unknown()),
    mihomoCustomParams: z.string().nullable(),
    sockoptParams: z.nullable(z.unknown()),
    finalMask: z.nullable(z.unknown()),

    inbound: z.object({
        configProfileUuid: z.uuid().nullable(),
        configProfileInboundUuid: z.uuid().nullable(),
    }),

    serverDescription: z.string().max(30).nullable(),
    tags: z.array(z.string()).default([]),
    isHidden: z.boolean().default(false),
    overrideSniFromAddress: z.boolean().default(false),
    keepSniBlank: z.boolean().default(false),
    overrideProtocolCredential: z.boolean().default(false),
    protocolCredential: z.string().nullable(),
    vlessRouteId: z.int().min(0).max(65535).nullable(),
    pinnedPeerCertSha256: z.string().nullable(),
    verifyPeerCertByName: z.string().nullable(),
    shuffleHost: z.boolean(),
    mihomoX25519: z.boolean(),
    mihomoIpVersion: z.enum(MIHOMO_IP_VERSION).nullable(),

    nodes: z.array(z.uuid()),
    xrayJsonTemplateUuid: z.uuid().nullable(),
    excludedInternalSquads: z.array(z.uuid()),
    excludeFromSubscriptionTypes: z.array(z.enum(SUBSCRIPTION_TEMPLATE_TYPE)),
});
