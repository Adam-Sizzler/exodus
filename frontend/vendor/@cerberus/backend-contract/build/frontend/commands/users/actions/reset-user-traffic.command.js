"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ResetUserTrafficCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
const api_1 = require("../../../api");
var ResetUserTrafficCommand;
(function (ResetUserTrafficCommand) {
    ResetUserTrafficCommand.url = api_1.REST_API.USERS.ACTIONS.RESET_TRAFFIC;
    ResetUserTrafficCommand.TSQ_url = ResetUserTrafficCommand.url(':uuid');
    ResetUserTrafficCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.ACTIONS.RESET_TRAFFIC(':uuid'), 'post', 'Reset user traffic');
    ResetUserTrafficCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    ResetUserTrafficCommand.ResponseSchema = zod_1.z.object({
        response: models_1.ExtendedUsersSchema,
    });
})(ResetUserTrafficCommand || (exports.ResetUserTrafficCommand = ResetUserTrafficCommand = {}));
