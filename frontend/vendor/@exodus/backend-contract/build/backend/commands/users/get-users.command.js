"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetUsersCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var GetUsersCommand;
(function (GetUsersCommand) {
    GetUsersCommand.url = api_1.REST_API.USERS.GET;
    GetUsersCommand.TSQ_url = GetUsersCommand.url;
    GetUsersCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.GET, 'get', 'Get all users using offset-based pagination', {
        scope: 'list',
        kind: 'read',
    }, 'Please note that the filters here are primarily intended for use by the frontend and rely on expensive operators such as LIKE under the hood. Misusing these filters may negatively impact the performance of your database.');
    GetUsersCommand.RequestQuerySchema = models_1.TanstackQueryRequestQuerySchema;
    GetUsersCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            users: zod_1.z.array(models_1.ExtendedUsersSchema),
            total: zod_1.z.number(),
        }),
    });
})(GetUsersCommand || (exports.GetUsersCommand = GetUsersCommand = {}));
