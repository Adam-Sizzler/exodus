"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetUserByShortUuidCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const user_response_1 = require("../user.response");
var GetUserByShortUuidCommand;
(function (GetUserByShortUuidCommand) {
    GetUserByShortUuidCommand.url = api_1.REST_API.USERS.GET_BY.SHORT_UUID;
    GetUserByShortUuidCommand.TSQ_url = GetUserByShortUuidCommand.url(':shortUuid');
    GetUserByShortUuidCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.GET_BY.SHORT_UUID(':shortUuid'), 'get', 'Get user by Short UUID', { scope: 'by-short-uuid', kind: 'read' });
    GetUserByShortUuidCommand.RequestParamSchema = zod_1.z.object({
        shortUuid: zod_1.z.string(),
    });
    GetUserByShortUuidCommand.ResponseSchema = user_response_1.UserResponseSchema;

})(GetUserByShortUuidCommand || (exports.GetUserByShortUuidCommand = GetUserByShortUuidCommand = {}));
