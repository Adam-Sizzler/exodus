"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetUserByIdCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
const user_response_1 = require("./user.response");
var GetUserByIdCommand;
(function (GetUserByIdCommand) {
    GetUserByIdCommand.url = api_1.REST_API.USERS.GET_BY_ID;
    GetUserByIdCommand.TSQ_url = GetUserByIdCommand.url(':userId');
    GetUserByIdCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.GET_BY_ID(':userId'), 'get', 'Get user by ID', { scope: 'by-id', kind: 'read' });
    GetUserByIdCommand.RequestParamSchema = zod_1.z.object({
        userId: models_1.numberParamSchema,
    });
    GetUserByIdCommand.ResponseSchema = user_response_1.UserResponseSchema;
})(GetUserByIdCommand || (exports.GetUserByIdCommand = GetUserByIdCommand = {}));
