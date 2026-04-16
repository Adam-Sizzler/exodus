"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetUserByUuidCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
const api_1 = require("../../api");
var GetUserByUuidCommand;
(function (GetUserByUuidCommand) {
    GetUserByUuidCommand.url = api_1.REST_API.USERS.GET_BY_UUID;
    GetUserByUuidCommand.TSQ_url = GetUserByUuidCommand.url(':uuid');
    GetUserByUuidCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.GET_BY_UUID(':uuid'), 'get', 'Get user by UUID');
    GetUserByUuidCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    GetUserByUuidCommand.ResponseSchema = zod_1.z.object({
        response: models_1.ExtendedUsersSchema,
    });
})(GetUserByUuidCommand || (exports.GetUserByUuidCommand = GetUserByUuidCommand = {}));
