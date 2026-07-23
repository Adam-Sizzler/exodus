"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetStatusCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var GetStatusCommand;
(function (GetStatusCommand) {
    GetStatusCommand.url = api_1.REST_API.AUTH.GET_STATUS;
    GetStatusCommand.TSQ_url = GetStatusCommand.url;
    GetStatusCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.AUTH_ROUTES.GET_STATUS, 'get', 'Get the status of the authentication', { scope: 'get-status', kind: 'read' });
    GetStatusCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            isLoginAllowed: zod_1.z.boolean(),
            isRegisterAllowed: zod_1.z.boolean(),
            authentication: zod_1.z.nullable(zod_1.z.object({
                passkey: zod_1.z.object({
                    enabled: zod_1.z.boolean(),
                }),
                oauth2: zod_1.z.object({
                    providers: zod_1.z.record(zod_1.z.nativeEnum(constants_1.OAUTH2_PROVIDERS), zod_1.z.boolean()),
                }),
                password: zod_1.z.object({
                    enabled: zod_1.z.boolean(),
                }),
            })),
            branding: zod_1.z.object({
                title: zod_1.z.nullable(zod_1.z.string()),
                logoUrl: zod_1.z.nullable(zod_1.z.string()),
            }),
        }),
    });
})(GetStatusCommand || (exports.GetStatusCommand = GetStatusCommand = {}));
