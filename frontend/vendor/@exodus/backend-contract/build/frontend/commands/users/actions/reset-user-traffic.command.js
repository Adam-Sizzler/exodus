"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ResetUserTrafficCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
const user_response_1 = require("../user.response");
var ResetUserTrafficCommand;
(function (ResetUserTrafficCommand) {
    ResetUserTrafficCommand.url = api_1.REST_API.USERS.ACTIONS.RESET_TRAFFIC;
    ResetUserTrafficCommand.TSQ_url = ResetUserTrafficCommand.url(':userId');
    ResetUserTrafficCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.ACTIONS.RESET_TRAFFIC(':userId'), 'post', 'Reset user traffic', { scope: 'reset-traffic', kind: 'write' });
    ResetUserTrafficCommand.RequestParamSchema = zod_1.z.object({
        userId: models_1.numberParamSchema,
    });
    ResetUserTrafficCommand.ResponseSchema = user_response_1.UserResponseSchema;
})(ResetUserTrafficCommand || (exports.ResetUserTrafficCommand = ResetUserTrafficCommand = {}));
