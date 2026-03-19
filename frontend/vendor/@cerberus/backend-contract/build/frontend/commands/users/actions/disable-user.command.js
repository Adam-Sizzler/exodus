"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DisableUserCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
const api_1 = require("../../../api");
var DisableUserCommand;
(function (DisableUserCommand) {
    DisableUserCommand.url = api_1.REST_API.USERS.ACTIONS.DISABLE;
    DisableUserCommand.TSQ_url = DisableUserCommand.url(':uuid');
    DisableUserCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.ACTIONS.DISABLE(':uuid'), 'post', 'Disable user');
    DisableUserCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    DisableUserCommand.ResponseSchema = zod_1.z.object({
        response: models_1.ExtendedUsersSchema,
    });
})(DisableUserCommand || (exports.DisableUserCommand = DisableUserCommand = {}));
