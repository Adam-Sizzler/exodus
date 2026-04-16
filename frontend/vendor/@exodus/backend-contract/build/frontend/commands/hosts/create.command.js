"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CreateHostCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const api_1 = require("../../api");
const models_1 = require("../../models");
var CreateHostCommand;
(function (CreateHostCommand) {
    CreateHostCommand.url = api_1.REST_API.HOSTS.CREATE;
    CreateHostCommand.TSQ_url = CreateHostCommand.url;
    CreateHostCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.HOSTS_ROUTES.CREATE, 'post', 'Create a new host');
    CreateHostCommand.RequestSchema = zod_1.z.object({
        inbound: zod_1.z.object({
            configProfileUuid: zod_1.z.string().uuid(),
            configProfileInboundUuid: zod_1.z.string().uuid(),
        }),
        remark: zod_1.z
            .string({
            invalid_type_error: 'Remark must be a string',
        })
            .min(1, {
            message: 'Remark must be at least 1 character',
        })
            .max(40, {
            message: 'Remark must be less than 40 characters',
        }),
        address: zod_1.z.string({
            invalid_type_error: 'Address must be a string',
        }),
        port: zod_1.z
            .number({
            invalid_type_error: 'Port must be an integer',
        })
            .int(),
        path: zod_1.z.string().optional(),
        sni: zod_1.z.string().optional(),
        host: zod_1.z.string().optional(),
        alpn: zod_1.z.optional(zod_1.z.nativeEnum(constants_1.ALPN).nullable()),
        fingerprint: zod_1.z.optional(zod_1.z.nativeEnum(constants_1.FINGERPRINTS).nullable()),
        isDisabled: zod_1.z.optional(zod_1.z.boolean().default(false)),
        securityLayer: zod_1.z.optional(zod_1.z.nativeEnum(constants_1.SECURITY_LAYERS).default(constants_1.SECURITY_LAYERS.DEFAULT)),
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
        isHidden: zod_1.z.optional(zod_1.z.boolean().default(false)),
        overrideSniFromAddress: zod_1.z.optional(zod_1.z.boolean().default(false)),
        keepSniBlank: zod_1.z.optional(zod_1.z.boolean().default(false)),
        allowInsecure: zod_1.z.optional(zod_1.z.boolean().default(false)),
        vlessRouteId: zod_1.z.optional(zod_1.z.number().int().min(0).max(65535).nullable()),
        shuffleHost: zod_1.z.optional(zod_1.z.boolean().default(false)),
        mihomoX25519: zod_1.z.optional(zod_1.z.boolean().default(false)),
        nodes: zod_1.z.optional(zod_1.z.array(zod_1.z.string().uuid())),
        xrayJsonTemplateUuid: zod_1.z.optional(zod_1.z.string().uuid().nullable()),
        excludedInternalSquads: zod_1.z
            .optional(zod_1.z.array(zod_1.z.string().uuid()))
            .describe('Optional. Internal squads from which the host will be excluded.'),
    });
    CreateHostCommand.ResponseSchema = zod_1.z.object({
        response: models_1.HostsSchema,
    });
})(CreateHostCommand || (exports.CreateHostCommand = CreateHostCommand = {}));
