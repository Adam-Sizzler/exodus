"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BulkAllResetTrafficUsersCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const api_1 = require("../../../api");
var BulkAllResetTrafficUsersCommand;
(function (BulkAllResetTrafficUsersCommand) {
    BulkAllResetTrafficUsersCommand.url = api_1.REST_API.USERS.BULK.ALL.RESET_TRAFFIC;
    BulkAllResetTrafficUsersCommand.TSQ_url = BulkAllResetTrafficUsersCommand.url;
    BulkAllResetTrafficUsersCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.BULK.ALL.RESET_TRAFFIC, 'post', 'Reset user used traffic for all users');
    BulkAllResetTrafficUsersCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            eventSent: zod_1.z.boolean(),
        }),
    });
})(BulkAllResetTrafficUsersCommand || (exports.BulkAllResetTrafficUsersCommand = BulkAllResetTrafficUsersCommand = {}));
