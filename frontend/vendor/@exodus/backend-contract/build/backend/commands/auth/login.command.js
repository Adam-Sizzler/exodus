"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LoginCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var LoginCommand;
(function (LoginCommand) {
    LoginCommand.url = api_1.REST_API.AUTH.LOGIN;
    LoginCommand.TSQ_url = LoginCommand.url;
    LoginCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.AUTH_ROUTES.LOGIN, 'post', 'Login as superadmin', { scope: 'login', kind: 'write' });
    LoginCommand.RequestBodySchema = zod_1.z.object({
        username: zod_1.z.string().describe('Username of the user'),
        password: zod_1.z.string().describe('Password of the user'),
    });
    LoginCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            accessToken: zod_1.z.string(),
        }),
    });

})(LoginCommand || (exports.LoginCommand = LoginCommand = {}));
