"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.VerifyPasskeyRegistrationCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var VerifyPasskeyRegistrationCommand;
(function (VerifyPasskeyRegistrationCommand) {
    VerifyPasskeyRegistrationCommand.url = api_1.REST_API.PASSKEYS.VERIFY_REGISTRATION;
    VerifyPasskeyRegistrationCommand.TSQ_url = VerifyPasskeyRegistrationCommand.url;
    VerifyPasskeyRegistrationCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.PASSKEYS_ROUTES.VERIFY_REGISTRATION, 'post', 'Verify registration for passkey');
    VerifyPasskeyRegistrationCommand.RequestSchema = zod_1.z.object({
        response: zod_1.z.unknown(),
    });
    VerifyPasskeyRegistrationCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            verified: zod_1.z.boolean(),
        }),
    });
})(VerifyPasskeyRegistrationCommand || (exports.VerifyPasskeyRegistrationCommand = VerifyPasskeyRegistrationCommand = {}));
