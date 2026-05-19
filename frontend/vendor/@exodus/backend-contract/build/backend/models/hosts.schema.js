"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.HostsSchema = void 0;
const zod_1 = require("zod");
const constants_1 = require("../constants");
const hosts_1 = require("../constants/hosts");
exports.HostsSchema = zod_1.z.object({
    uuid: zod_1.z.string().uuid(),
    viewPosition: zod_1.z.number().int(),
    remark: zod_1.z.string(),
    address: zod_1.z.string(),
    port: zod_1.z.number().int(),
    path: zod_1.z.string().nullable(),
    sni: zod_1.z.string().nullable(),
    host: zod_1.z.string().nullable(),
    alpn: zod_1.z.string().nullable(),
    fingerprint: zod_1.z.string().nullable(),
    isDisabled: zod_1.z.boolean().default(false),
    securityLayer: zod_1.z.nativeEnum(hosts_1.SECURITY_LAYERS).default(hosts_1.SECURITY_LAYERS.DEFAULT),
    muxParams: zod_1.z.nullable(zod_1.z.unknown()),
    singboxMuxParams: zod_1.z.nullable(zod_1.z.unknown()),
    clashMuxParams: zod_1.z.nullable(zod_1.z.unknown()),
    sockoptParams: zod_1.z.nullable(zod_1.z.unknown()),
    inbound: zod_1.z.object({
        configProfileUuid: zod_1.z.string().uuid().nullable(),
        configProfileInboundUuid: zod_1.z.string().uuid().nullable(),
    }),
    serverDescription: zod_1.z.string().max(30).nullable(),
    tag: zod_1.z.string().nullable(),
    isHidden: zod_1.z.boolean().default(false),
    overrideSniFromAddress: zod_1.z.boolean().default(false),
    keepSniBlank: zod_1.z.boolean().default(false),
    overrideProtocolCredential: zod_1.z.boolean().default(false),
    protocolCredential: zod_1.z.string().max(256).nullable(),
    allowInsecure: zod_1.z.boolean().default(false),
    shuffleHost: zod_1.z.boolean(),
    selectorNodesFirst: zod_1.z.boolean().default(false),
    mihomoX25519: zod_1.z.boolean(),
    nodes: zod_1.z.array(zod_1.z.string().uuid()),
    xrayJsonTemplateUuid: zod_1.z.string().uuid().nullable(),
    excludedInternalSquads: zod_1.z.array(zod_1.z.string().uuid()),
    excludeFromSubscriptionTypes: zod_1.z.array(zod_1.z.nativeEnum(constants_1.SUBSCRIPTION_TEMPLATE_TYPE)).optional(),
});
