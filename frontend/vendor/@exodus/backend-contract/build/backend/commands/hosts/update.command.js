"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UpdateHostCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const api_1 = require("../../api");
const models_1 = require("../../models");
var UpdateHostCommand;
(function (UpdateHostCommand) {
    UpdateHostCommand.url = api_1.REST_API.HOSTS.UPDATE;
    UpdateHostCommand.TSQ_url = UpdateHostCommand.url;
    UpdateHostCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.HOSTS_ROUTES.UPDATE, 'patch', 'Update a host');
    UpdateHostCommand.RequestSchema = models_1.HostsSchema.pick({
        uuid: true,
    }).extend({
        inbound: zod_1.z
            .object({
            configProfileUuid: zod_1.z.string().uuid(),
            configProfileInboundUuid: zod_1.z.string().uuid(),
        })
            .optional(),
        remark: zod_1.z
            .string({
            invalid_type_error: 'Remark must be a string',
        })
            .max(40, {
            message: 'Remark must be less than 40 characters',
        })
            .optional(),
        address: zod_1.z
            .string({
            invalid_type_error: 'Address must be a string',
        })
            .optional(),
        port: zod_1.z
            .number({
            invalid_type_error: 'Port must be an integer',
        })
            .int()
            .optional(),
        path: zod_1.z.optional(zod_1.z.string()),
        sni: zod_1.z.optional(zod_1.z.string()),
        host: zod_1.z.optional(zod_1.z.string()),
        alpn: zod_1.z.optional(zod_1.z.nativeEnum(constants_1.ALPN).nullable()),
        fingerprint: zod_1.z.optional(zod_1.z.nativeEnum(constants_1.FINGERPRINTS).nullable()),
        isDisabled: zod_1.z.optional(zod_1.z.boolean()),
        securityLayer: zod_1.z.optional(zod_1.z.nativeEnum(constants_1.SECURITY_LAYERS)),
        xHttpExtraParams: zod_1.z.optional(zod_1.z.nullable(zod_1.z.unknown())),
        muxParams: zod_1.z.optional(zod_1.z.nullable(zod_1.z.unknown())),
        sockoptParams: zod_1.z.optional(zod_1.z.nullable(zod_1.z.unknown())),
        serverDescription: zod_1.z.optional(zod_1.z
            .string()
            .max(30, {
            message: 'Server description must be less than 30 characters',
        })
            .nullable()),
        tag: zod_1.z
            .optional(zod_1.z
            .string()
            .regex(/^[A-Z0-9_:]+$/, 'Tag can only contain uppercase letters, numbers, underscores and colons')
            .max(32, 'Tag must be less than 32 characters')
            .nullable())
            .describe('Optional. Host tag for categorization. Max 32 characters, uppercase letters, numbers, underscores and colons are allowed.'),
        isHidden: zod_1.z.optional(zod_1.z.boolean()),
        overrideSniFromAddress: zod_1.z.optional(zod_1.z.boolean()),
        keepSniBlank: zod_1.z.optional(zod_1.z.boolean()),
        vlessRouteId: zod_1.z.optional(zod_1.z.number().int().min(0).max(65535).nullable()),
        allowInsecure: zod_1.z.optional(zod_1.z.boolean()),
        shuffleHost: zod_1.z.optional(zod_1.z.boolean()),
        mihomoX25519: zod_1.z.optional(zod_1.z.boolean()),
        nodes: zod_1.z.optional(zod_1.z.array(zod_1.z.string().uuid())),
        xrayJsonTemplateUuid: zod_1.z.optional(zod_1.z.string().uuid().nullable()),
        excludedInternalSquads: zod_1.z
            .optional(zod_1.z.array(zod_1.z.string().uuid()))
            .describe('Optional. Internal squads from which the host will be excluded.'),
    });
    UpdateHostCommand.ResponseSchema = zod_1.z.object({
        response: models_1.HostsSchema,
    });
})(UpdateHostCommand || (exports.UpdateHostCommand = UpdateHostCommand = {}));
