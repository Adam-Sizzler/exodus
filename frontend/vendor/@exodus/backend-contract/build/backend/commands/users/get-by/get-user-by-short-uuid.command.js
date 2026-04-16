"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetUserByShortUuidCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
const api_1 = require("../../../api");
var GetUserByShortUuidCommand;
(function (GetUserByShortUuidCommand) {
    GetUserByShortUuidCommand.url = api_1.REST_API.USERS.GET_BY.SHORT_UUID;
    GetUserByShortUuidCommand.TSQ_url = GetUserByShortUuidCommand.url(':shortUuid');
    GetUserByShortUuidCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.GET_BY.SHORT_UUID(':shortUuid'), 'get', 'Get user by Short UUID');
    GetUserByShortUuidCommand.RequestSchema = zod_1.z.object({
        shortUuid: zod_1.z.string(),
    });
    GetUserByShortUuidCommand.ResponseSchema = zod_1.z.object({
        response: models_1.ExtendedUsersSchema,
    });
})(GetUserByShortUuidCommand || (exports.GetUserByShortUuidCommand = GetUserByShortUuidCommand = {}));
