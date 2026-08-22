"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetUserByUsernameCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const user_response_1 = require("../user.response");
var GetUserByUsernameCommand;
(function (GetUserByUsernameCommand) {
    GetUserByUsernameCommand.url = api_1.REST_API.USERS.GET_BY.USERNAME;
    GetUserByUsernameCommand.TSQ_url = GetUserByUsernameCommand.url(':username');
    GetUserByUsernameCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.GET_BY.USERNAME(':username'), 'get', 'Get user by username', { scope: 'by-username', kind: 'read' });
    GetUserByUsernameCommand.RequestParamSchema = zod_1.z.object({
        username: zod_1.z.string(),
    });
    GetUserByUsernameCommand.ResponseSchema = user_response_1.UserResponseSchema;
})(GetUserByUsernameCommand || (exports.GetUserByUsernameCommand = GetUserByUsernameCommand = {}));
