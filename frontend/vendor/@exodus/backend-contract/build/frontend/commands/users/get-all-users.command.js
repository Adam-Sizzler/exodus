"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetAllUsersCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var GetAllUsersCommand;
(function (GetAllUsersCommand) {
    GetAllUsersCommand.url = api_1.REST_API.USERS.GET;
    GetAllUsersCommand.TSQ_url = GetAllUsersCommand.url;
    GetAllUsersCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.GET, 'get', 'Get all users using offset-based pagination', {
        scope: 'list',
        kind: 'read',
    });
    GetAllUsersCommand.RequestQuerySchema = models_1.TanstackQueryRequestQuerySchema;
    GetAllUsersCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            users: zod_1.z.array(models_1.ExtendedUsersSchema),
            total: zod_1.z.number(),
        }),
    });
})(GetAllUsersCommand || (exports.GetAllUsersCommand = GetAllUsersCommand = {}));
