"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.EnableUserCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
const api_1 = require("../../../api");
var EnableUserCommand;
(function (EnableUserCommand) {
    EnableUserCommand.url = api_1.REST_API.USERS.ACTIONS.ENABLE;
    EnableUserCommand.TSQ_url = EnableUserCommand.url(':uuid');
    EnableUserCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.ACTIONS.ENABLE(':uuid'), 'post', 'Enable user');
    EnableUserCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    EnableUserCommand.ResponseSchema = zod_1.z.object({
        response: models_1.ExtendedUsersSchema,
    });
})(EnableUserCommand || (exports.EnableUserCommand = EnableUserCommand = {}));
