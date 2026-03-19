"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetUserByIdCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
const api_1 = require("../../../api");
var GetUserByIdCommand;
(function (GetUserByIdCommand) {
    GetUserByIdCommand.url = api_1.REST_API.USERS.GET_BY.ID;
    GetUserByIdCommand.TSQ_url = GetUserByIdCommand.url(':id');
    GetUserByIdCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.GET_BY.ID(':id'), 'get', 'Get user by ID');
    GetUserByIdCommand.RequestSchema = zod_1.z.object({
        id: zod_1.z.coerce.bigint(),
    });
    GetUserByIdCommand.ResponseSchema = zod_1.z.object({
        response: models_1.ExtendedUsersSchema,
    });
})(GetUserByIdCommand || (exports.GetUserByIdCommand = GetUserByIdCommand = {}));
