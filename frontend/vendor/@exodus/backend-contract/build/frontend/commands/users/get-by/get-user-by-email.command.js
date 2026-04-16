"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetUserByEmailCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
const api_1 = require("../../../api");
var GetUserByEmailCommand;
(function (GetUserByEmailCommand) {
    GetUserByEmailCommand.url = api_1.REST_API.USERS.GET_BY.EMAIL;
    GetUserByEmailCommand.TSQ_url = GetUserByEmailCommand.url(':email');
    GetUserByEmailCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.GET_BY.EMAIL(':email'), 'get', 'Get users by email');
    GetUserByEmailCommand.RequestSchema = zod_1.z.object({
        email: zod_1.z.string().email(),
    });
    GetUserByEmailCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.array(models_1.ExtendedUsersSchema),
    });
})(GetUserByEmailCommand || (exports.GetUserByEmailCommand = GetUserByEmailCommand = {}));
