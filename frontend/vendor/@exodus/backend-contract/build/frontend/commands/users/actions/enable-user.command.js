"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.EnableUserCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
const user_response_1 = require("../user.response");
var EnableUserCommand;
(function (EnableUserCommand) {
    EnableUserCommand.url = api_1.REST_API.USERS.ACTIONS.ENABLE;
    EnableUserCommand.TSQ_url = EnableUserCommand.url(':userId');
    EnableUserCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.ACTIONS.ENABLE(':userId'), 'post', 'Enable user', { scope: 'enable', kind: 'write' });
    EnableUserCommand.RequestParamSchema = zod_1.z.object({
        userId: models_1.numberParamSchema,
    });
    EnableUserCommand.ResponseSchema = user_response_1.UserResponseSchema;
})(EnableUserCommand || (exports.EnableUserCommand = EnableUserCommand = {}));
