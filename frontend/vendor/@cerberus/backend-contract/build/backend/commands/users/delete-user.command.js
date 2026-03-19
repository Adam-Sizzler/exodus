"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DeleteUserCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const api_1 = require("../../api");
var DeleteUserCommand;
(function (DeleteUserCommand) {
    DeleteUserCommand.url = api_1.REST_API.USERS.DELETE;
    DeleteUserCommand.TSQ_url = DeleteUserCommand.url(':uuid');
    DeleteUserCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.DELETE(':uuid'), 'delete', 'Delete user');
    DeleteUserCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    DeleteUserCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            isDeleted: zod_1.z.boolean(),
        }),
    });
})(DeleteUserCommand || (exports.DeleteUserCommand = DeleteUserCommand = {}));
