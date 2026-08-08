"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.HostsSchema = void 0;
const zod_1 = require("zod");
const constants_1 = require("../constants");
const hosts_1 = require("../constants/hosts");
exports.HostsSchema = zod_1.z.object({
    uuid: zod_1.z.uuid(),
    viewPosition: zod_1.z.int(),
    remark: zod_1.z.string(),
    address: zod_1.z.string(),
    port: zod_1.z.int(),
    path: zod_1.z.string().nullable(),
    sni: zod_1.z.string().nullable(),
    host: zod_1.z.string().nullable(),
    alpn: zod_1.z.enum(hosts_1.ALPN).nullable(),
    fingerprint: zod_1.z.string().nullable(),
    isDisabled: zod_1.z.boolean(),
    securityLayer: zod_1.z.enum(hosts_1.SECURITY_LAYERS).default(hosts_1.SECURITY_LAYERS.DEFAULT),
    xhttpExtraParams: zod_1.z.nullable(zod_1.z.unknown()).optional(),
    muxParams: zod_1.z.nullable(zod_1.z.unknown()).optional(),
    sockoptParams: zod_1.z.nullable(zod_1.z.unknown()).optional(),
    finalMask: zod_1.z.nullable(zod_1.z.unknown()).optional(),
    inbound: zod_1.z.object({
        configProfileUuid: zod_1.z.uuid().nullable(),
        configProfileInboundUuid: zod_1.z.uuid().nullable(),
    }),
    serverDescription: zod_1.z.string().max(30).nullable(),
    tags: zod_1.z.array(zod_1.z.string()).default([]),
    isHidden: zod_1.z.boolean().default(false),
    overrideSniFromAddress: zod_1.z.boolean().default(false),
    keepSniBlank: zod_1.z.boolean().default(false),
    vlessRouteId: zod_1.z.int().min(0).max(65535).nullable(),
    pinnedPeerCertSha256: zod_1.z.string().nullable(),
    verifyPeerCertByName: zod_1.z.string().nullable(),
    shuffleHost: zod_1.z.boolean(),
    mihomoX25519: zod_1.z.boolean(),
    mihomoIpVersion: zod_1.z.enum(hosts_1.MIHOMO_IP_VERSION).nullable(),
    nodes: zod_1.z.array(zod_1.z.uuid()),
    xrayJsonTemplateUuid: zod_1.z.uuid().nullable(),
    excludedInternalSquads: zod_1.z.array(zod_1.z.uuid()),
    excludeFromSubscriptionTypes: zod_1.z.array(zod_1.z.enum(constants_1.SUBSCRIPTION_TEMPLATE_TYPE)),
});
