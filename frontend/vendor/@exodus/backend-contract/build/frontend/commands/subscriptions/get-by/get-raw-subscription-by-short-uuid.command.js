"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetRawSubscriptionByShortUuidCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
var GetRawSubscriptionByShortUuidCommand;
(function (GetRawSubscriptionByShortUuidCommand) {
    GetRawSubscriptionByShortUuidCommand.url = api_1.REST_API.SUBSCRIPTIONS.GET_BY.SHORT_UUID_RAW;
    GetRawSubscriptionByShortUuidCommand.TSQ_url = GetRawSubscriptionByShortUuidCommand.url(':shortUuid');
    GetRawSubscriptionByShortUuidCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTIONS_ROUTES.GET_BY.SHORT_UUID_RAW(':shortUuid'), 'get', 'Get Raw Subscription by Short UUID');
    GetRawSubscriptionByShortUuidCommand.RequestSchema = zod_1.z.object({
        shortUuid: zod_1.z.string(),
    });
    GetRawSubscriptionByShortUuidCommand.RequestQuerySchema = zod_1.z.object({
        withDisabledHosts: zod_1.z
            .string()
            .transform((str) => str === 'true')
            .optional()
            .default('false'),
    });
    GetRawSubscriptionByShortUuidCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            user: models_1.ExtendedUsersSchema,
            convertedUserInfo: zod_1.z.object({
                daysLeft: zod_1.z.number(),
                trafficLimit: zod_1.z.string(),
                trafficUsed: zod_1.z.string(),
                lifetimeTrafficUsed: zod_1.z.string(),
                isHwidLimited: zod_1.z.boolean(),
            }),
            headers: zod_1.z.record(zod_1.z.string(), zod_1.z.string().optional()),
            rawHosts: zod_1.z.array(zod_1.z.object({
                address: zod_1.z.optional(zod_1.z.nullable(zod_1.z.string())),
                alpn: zod_1.z.optional(zod_1.z.nullable(zod_1.z.string())),
                fingerprint: zod_1.z.optional(zod_1.z.nullable(zod_1.z.string())),
                host: zod_1.z.optional(zod_1.z.nullable(zod_1.z.string())),
                network: zod_1.z.optional(zod_1.z.nullable(zod_1.z.string())),
                password: zod_1.z.object({
                    ssPassword: zod_1.z.string(),
                    trojanPassword: zod_1.z.string(),
                    vlessPassword: zod_1.z.string(),
                }),
                path: zod_1.z.optional(zod_1.z.nullable(zod_1.z.string())),
                publicKey: zod_1.z.optional(zod_1.z.nullable(zod_1.z.string())),
                port: zod_1.z.optional(zod_1.z.nullable(zod_1.z.number())),
                protocol: zod_1.z.optional(zod_1.z.nullable(zod_1.z.string())),
                remark: zod_1.z.optional(zod_1.z.nullable(zod_1.z.string())),
                shortId: zod_1.z.optional(zod_1.z.nullable(zod_1.z.string())),
                sni: zod_1.z.optional(zod_1.z.nullable(zod_1.z.string())),
                spiderX: zod_1.z.optional(zod_1.z.nullable(zod_1.z.string())),
                tls: zod_1.z.optional(zod_1.z.nullable(zod_1.z.string())),
                rawSettings: zod_1.z.optional(zod_1.z.nullable(zod_1.z.object({
                    headerType: zod_1.z.optional(zod_1.z.nullable(zod_1.z.string())),
                    request: zod_1.z.optional(zod_1.z.nullable(zod_1.z.object({}))),
                }))),
                additionalParams: zod_1.z.optional(zod_1.z.nullable(zod_1.z.object({
                    mode: zod_1.z.optional(zod_1.z.nullable(zod_1.z.string())),
                    heartbeatPeriod: zod_1.z.optional(zod_1.z.nullable(zod_1.z.number())),
                }))),
                xHttpExtraParams: zod_1.z.optional(zod_1.z.nullable(zod_1.z.object({}))),
                muxParams: zod_1.z.optional(zod_1.z.nullable(zod_1.z.object({}))),
                sockoptParams: zod_1.z.optional(zod_1.z.nullable(zod_1.z.object({}))),
                serverDescription: zod_1.z.optional(zod_1.z.nullable(zod_1.z.string())),
                flow: zod_1.z.optional(zod_1.z.nullable(zod_1.z.string())),
                allowInsecure: zod_1.z.optional(zod_1.z.nullable(zod_1.z.boolean())),
                shuffleHost: zod_1.z.optional(zod_1.z.nullable(zod_1.z.boolean())),
                mihomoX25519: zod_1.z.optional(zod_1.z.nullable(zod_1.z.boolean())),
                mldsa65Verify: zod_1.z.optional(zod_1.z.nullable(zod_1.z.string())),
                encryption: zod_1.z.optional(zod_1.z.nullable(zod_1.z.string())),
                protocolOptions: zod_1.z.optional(zod_1.z.nullable(zod_1.z.object({
                    ss: zod_1.z.optional(zod_1.z.nullable(zod_1.z.object({
                        method: zod_1.z.optional(zod_1.z.nullable(zod_1.z.string())),
                    }))),
                }))),
                dbData: zod_1.z.optional(zod_1.z.object({
                    rawInbound: zod_1.z.nullable(zod_1.z.object({})),
                    inboundTag: zod_1.z.string(),
                    uuid: zod_1.z.string(),
                    configProfileUuid: zod_1.z.nullable(zod_1.z.string()),
                    configProfileInboundUuid: zod_1.z.nullable(zod_1.z.string()),
                    isDisabled: zod_1.z.boolean(),
                    viewPosition: zod_1.z.number(),
                    remark: zod_1.z.string(),
                    isHidden: zod_1.z.boolean(),
                    tag: zod_1.z.nullable(zod_1.z.string()),
                    vlessRouteId: zod_1.z.number().int().nullable(),
                })),
                xrayJsonTemplate: zod_1.z.optional(zod_1.z.nullable(zod_1.z.object({}))),
            })),
        }),
    });
})(GetRawSubscriptionByShortUuidCommand || (exports.GetRawSubscriptionByShortUuidCommand = GetRawSubscriptionByShortUuidCommand = {}));
