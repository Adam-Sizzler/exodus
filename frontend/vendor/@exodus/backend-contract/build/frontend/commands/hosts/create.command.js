"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CreateHostCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
const host_response_1 = require("./host.response");
var CreateHostCommand;
(function (CreateHostCommand) {
    CreateHostCommand.url = api_1.REST_API.HOSTS.CREATE;
    CreateHostCommand.TSQ_url = CreateHostCommand.url;
    CreateHostCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.HOSTS_ROUTES.CREATE, 'post', 'Create a new host', { scope: 'create', kind: 'write' });
    CreateHostCommand.RequestBodySchema = zod_1.z.object({
        inbound: zod_1.z.object({
            configProfileUuid: zod_1.z.uuid(),
            configProfileInboundUuid: zod_1.z.uuid(),
        }),
        remark: zod_1.z.string().min(1).max(100),
        address: zod_1.z.string(),
        port: zod_1.z.int(),
        path: zod_1.z.string().nullish(),
        sni: zod_1.z.string().nullish(),
        host: zod_1.z.string().nullish(),
        alpn: zod_1.z.enum(constants_1.ALPN).nullish(),
        fingerprint: zod_1.z.string().nullish(),
        isDisabled: zod_1.z.optional(zod_1.z.boolean().default(false)),
        securityLayer: zod_1.z.optional(zod_1.z.enum(constants_1.SECURITY_LAYERS).default(constants_1.SECURITY_LAYERS.DEFAULT)),
        xhttpExtraParams: zod_1.z.unknown().nullish(),
        muxParams: zod_1.z.unknown().nullish(),
        sockoptParams: zod_1.z.unknown().nullish(),
        finalMask: zod_1.z.unknown().nullish(),
        serverDescription: zod_1.z.string().max(30).nullish(),
        tags: zod_1.z.optional(zod_1.z
            .array(zod_1.z
            .string()
            .regex(/^[A-Z0-9_:]+$/, 'Tag can only contain uppercase letters, numbers, underscores and colons')
            .max(36, 'Each tag must be less than 36 characters'))
            .max(10, 'Maximum 10 tags')),
        isHidden: zod_1.z.optional(zod_1.z.boolean().default(false)),
        overrideSniFromAddress: zod_1.z.optional(zod_1.z.boolean().default(false)),
        keepSniBlank: zod_1.z.optional(zod_1.z.boolean().default(false)),
        pinnedPeerCertSha256: zod_1.z.string().nullish(),
        verifyPeerCertByName: zod_1.z.string().nullish(),
        vlessRouteId: zod_1.z.int().min(0).max(65535).nullish(),
        shuffleHost: zod_1.z.optional(zod_1.z.boolean().default(false)),
        mihomoX25519: zod_1.z.optional(zod_1.z.boolean().default(false)),
        mihomoIpVersion: zod_1.z.enum(constants_1.MIHOMO_IP_VERSION).nullish(),
        nodes: zod_1.z.optional(zod_1.z.array(zod_1.z.uuid())),
        xrayJsonTemplateUuid: zod_1.z.uuid().nullish(),
        excludedInternalSquads: zod_1.z
            .optional(zod_1.z.array(zod_1.z.uuid()))
            .describe('Optional. Internal squads from which the host will be excluded.'),
        excludeFromSubscriptionTypes: zod_1.z
            .optional(zod_1.z.array(zod_1.z.enum(constants_1.SUBSCRIPTION_TEMPLATE_TYPE)))
            .describe('Optional. Subscription types from which the host will be excluded from.'),
        mapper: models_1.HostMapperSchema.optional(),
    });
    CreateHostCommand.ResponseSchema = host_response_1.HostResponseSchema;
})(CreateHostCommand || (exports.CreateHostCommand = CreateHostCommand = {}));
