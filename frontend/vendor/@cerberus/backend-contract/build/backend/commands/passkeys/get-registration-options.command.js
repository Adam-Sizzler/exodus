"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetPasskeyRegistrationOptionsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var GetPasskeyRegistrationOptionsCommand;
(function (GetPasskeyRegistrationOptionsCommand) {
    GetPasskeyRegistrationOptionsCommand.url = api_1.REST_API.PASSKEYS.GET_REGISTRATION_OPTIONS;
    GetPasskeyRegistrationOptionsCommand.TSQ_url = GetPasskeyRegistrationOptionsCommand.url;
    GetPasskeyRegistrationOptionsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.PASSKEYS_ROUTES.GET_REGISTRATION_OPTIONS, 'get', 'Get registration options for passkey');
    GetPasskeyRegistrationOptionsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.unknown(),
    });
})(GetPasskeyRegistrationOptionsCommand || (exports.GetPasskeyRegistrationOptionsCommand = GetPasskeyRegistrationOptionsCommand = {}));
