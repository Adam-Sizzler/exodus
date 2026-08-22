"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UpdateHostCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
const host_response_1 = require("./host.response");
var UpdateHostCommand;
(function (UpdateHostCommand) {
    UpdateHostCommand.url = api_1.REST_API.HOSTS.UPDATE;
    UpdateHostCommand.TSQ_url = UpdateHostCommand.url;
    UpdateHostCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.HOSTS_ROUTES.UPDATE, 'patch', 'Update a host', { scope: 'update', kind: 'write' });
    UpdateHostCommand.RequestBodySchema = models_1.HostsSchema.pick({
        uuid: true,
    }).extend({
        inbound: zod_1.z
            .object({
            configProfileUuid: zod_1.z.uuid(),
            configProfileInboundUuid: zod_1.z.uuid(),
        })
            .optional(),
        remark: zod_1.z.string().min(1).max(100).optional(),
        address: zod_1.z.string().optional(),
        port: zod_1.z.int().optional(),
        path: zod_1.z.string().nullish(),
        sni: zod_1.z.string().nullish(),
        host: zod_1.z.string().nullish(),
        alpn: zod_1.z.enum(constants_1.ALPN).nullish(),
        fingerprint: zod_1.z.string().nullish(),
        isDisabled: zod_1.z.boolean().default(false),
        securityLayer: zod_1.z.optional(zod_1.z.enum(constants_1.SECURITY_LAYERS)),
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
        isHidden: zod_1.z.optional(zod_1.z.boolean()),
        overrideSniFromAddress: zod_1.z.optional(zod_1.z.boolean()),
        keepSniBlank: zod_1.z.optional(zod_1.z.boolean()),
        vlessRouteId: zod_1.z.optional(zod_1.z.int().min(0).max(65535).nullable()),
        pinnedPeerCertSha256: zod_1.z.string().nullish(),
        verifyPeerCertByName: zod_1.z.string().nullish(),
        shuffleHost: zod_1.z.optional(zod_1.z.boolean()),
        mihomoX25519: zod_1.z.optional(zod_1.z.boolean()),
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
    UpdateHostCommand.ResponseSchema = host_response_1.HostResponseSchema;
})(UpdateHostCommand || (exports.UpdateHostCommand = UpdateHostCommand = {}));
