"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.VerifyPasskeyAuthenticationCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const api_1 = require("../../../api");
var VerifyPasskeyAuthenticationCommand;
(function (VerifyPasskeyAuthenticationCommand) {
    VerifyPasskeyAuthenticationCommand.url = api_1.REST_API.AUTH.PASSKEY.VERIFY_AUTHENTICATION;
    VerifyPasskeyAuthenticationCommand.TSQ_url = VerifyPasskeyAuthenticationCommand.url;
    VerifyPasskeyAuthenticationCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.AUTH_ROUTES.PASSKEY.VERIFY_AUTHENTICATION, 'post', 'Verify the authentication for passkey');
    VerifyPasskeyAuthenticationCommand.RequestSchema = zod_1.z.object({
        response: zod_1.z.unknown(),
    });
    VerifyPasskeyAuthenticationCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            accessToken: zod_1.z.string(),
        }),
    });
})(VerifyPasskeyAuthenticationCommand || (exports.VerifyPasskeyAuthenticationCommand = VerifyPasskeyAuthenticationCommand = {}));
