"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.RegisterCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const api_1 = require("../../api");
var RegisterCommand;
(function (RegisterCommand) {
    RegisterCommand.url = api_1.REST_API.AUTH.REGISTER;
    RegisterCommand.TSQ_url = RegisterCommand.url;
    RegisterCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.AUTH_ROUTES.REGISTER, 'post', 'Register as superadmin');
    RegisterCommand.RequestSchema = zod_1.z.object({
        username: zod_1.z.string(),
        password: zod_1.z
            .string()
            .min(24, 'Password must contain at least 24 characters')
            .regex(/^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9]).{24,}$/, 'Password must contain uppercase and lowercase letters and numbers, and be at least 24 characters long.'),
    });
    RegisterCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            accessToken: zod_1.z.string(),
        }),
    });
})(RegisterCommand || (exports.RegisterCommand = RegisterCommand = {}));
