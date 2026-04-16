"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.OAuth2CallbackCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const api_1 = require("../../../api");
var OAuth2CallbackCommand;
(function (OAuth2CallbackCommand) {
    OAuth2CallbackCommand.url = api_1.REST_API.AUTH.OAUTH2.CALLBACK;
    OAuth2CallbackCommand.TSQ_url = OAuth2CallbackCommand.url;
    OAuth2CallbackCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.AUTH_ROUTES.OAUTH2.CALLBACK, 'post', 'Callback from OAuth2');
    OAuth2CallbackCommand.RequestSchema = zod_1.z.object({
        provider: zod_1.z.nativeEnum(constants_1.OAUTH2_PROVIDERS),
        code: zod_1.z.string(),
        state: zod_1.z.string(),
    });
    OAuth2CallbackCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            accessToken: zod_1.z.string(),
        }),
    });
})(OAuth2CallbackCommand || (exports.OAuth2CallbackCommand = OAuth2CallbackCommand = {}));
