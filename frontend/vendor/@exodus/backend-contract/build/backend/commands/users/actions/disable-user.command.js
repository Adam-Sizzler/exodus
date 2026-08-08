"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DisableUserCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
const user_response_1 = require("../user.response");
var DisableUserCommand;
(function (DisableUserCommand) {
    DisableUserCommand.url = api_1.REST_API.USERS.ACTIONS.DISABLE;
    DisableUserCommand.TSQ_url = DisableUserCommand.url(':userId');
    DisableUserCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.ACTIONS.DISABLE(':userId'), 'post', 'Disable user', { scope: 'disable', kind: 'write' });
    DisableUserCommand.RequestParamSchema = zod_1.z.object({
        userId: models_1.numberParamSchema,
    });
    DisableUserCommand.ResponseSchema = user_response_1.UserResponseSchema;

})(DisableUserCommand || (exports.DisableUserCommand = DisableUserCommand = {}));
