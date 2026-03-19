"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.OAuth2AuthorizeCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const api_1 = require("../../../api");
var OAuth2AuthorizeCommand;
(function (OAuth2AuthorizeCommand) {
    OAuth2AuthorizeCommand.url = api_1.REST_API.AUTH.OAUTH2.AUTHORIZE;
    OAuth2AuthorizeCommand.TSQ_url = OAuth2AuthorizeCommand.url;
    OAuth2AuthorizeCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.AUTH_ROUTES.OAUTH2.AUTHORIZE, 'post', 'Initiate OAuth2 authorization');
    OAuth2AuthorizeCommand.RequestSchema = zod_1.z.object({
        provider: zod_1.z.nativeEnum(constants_1.OAUTH2_PROVIDERS),
    });
    OAuth2AuthorizeCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            authorizationUrl: zod_1.z.nullable(zod_1.z.string().url()),
        }),
    });
})(OAuth2AuthorizeCommand || (exports.OAuth2AuthorizeCommand = OAuth2AuthorizeCommand = {}));
