"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DeleteUserCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var DeleteUserCommand;
(function (DeleteUserCommand) {
    DeleteUserCommand.url = api_1.REST_API.USERS.DELETE;
    DeleteUserCommand.TSQ_url = DeleteUserCommand.url(':userId');
    DeleteUserCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.DELETE(':userId'), 'delete', 'Delete user', { scope: 'delete', kind: 'write' });
    DeleteUserCommand.RequestParamSchema = zod_1.z.object({
        userId: models_1.numberParamSchema,
    });

})(DeleteUserCommand || (exports.DeleteUserCommand = DeleteUserCommand = {}));
