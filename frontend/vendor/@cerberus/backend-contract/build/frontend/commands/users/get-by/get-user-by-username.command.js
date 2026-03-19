"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetUserByUsernameCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
const api_1 = require("../../../api");
var GetUserByUsernameCommand;
(function (GetUserByUsernameCommand) {
    GetUserByUsernameCommand.url = api_1.REST_API.USERS.GET_BY.USERNAME;
    GetUserByUsernameCommand.TSQ_url = GetUserByUsernameCommand.url(':username');
    GetUserByUsernameCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.GET_BY.USERNAME(':username'), 'get', 'Get user by username');
    GetUserByUsernameCommand.RequestSchema = zod_1.z.object({
        username: zod_1.z.string(),
    });
    GetUserByUsernameCommand.ResponseSchema = zod_1.z.object({
        response: models_1.ExtendedUsersSchema,
    });
})(GetUserByUsernameCommand || (exports.GetUserByUsernameCommand = GetUserByUsernameCommand = {}));
