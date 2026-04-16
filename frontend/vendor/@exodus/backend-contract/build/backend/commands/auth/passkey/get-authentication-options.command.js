"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetPasskeyAuthenticationOptionsCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const api_1 = require("../../../api");
var GetPasskeyAuthenticationOptionsCommand;
(function (GetPasskeyAuthenticationOptionsCommand) {
    GetPasskeyAuthenticationOptionsCommand.url = api_1.REST_API.AUTH.PASSKEY.GET_AUTHENTICATION_OPTIONS;
    GetPasskeyAuthenticationOptionsCommand.TSQ_url = GetPasskeyAuthenticationOptionsCommand.url;
    GetPasskeyAuthenticationOptionsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.AUTH_ROUTES.PASSKEY.GET_AUTHENTICATION_OPTIONS, 'get', 'Get the authentication options for passkey');
    GetPasskeyAuthenticationOptionsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.unknown(),
    });
})(GetPasskeyAuthenticationOptionsCommand || (exports.GetPasskeyAuthenticationOptionsCommand = GetPasskeyAuthenticationOptionsCommand = {}));
