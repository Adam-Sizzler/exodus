"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BulkAllResetTrafficUsersCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var BulkAllResetTrafficUsersCommand;
(function (BulkAllResetTrafficUsersCommand) {
    BulkAllResetTrafficUsersCommand.url = api_1.REST_API.USERS.BULK.ALL.RESET_TRAFFIC;
    BulkAllResetTrafficUsersCommand.TSQ_url = BulkAllResetTrafficUsersCommand.url;
    BulkAllResetTrafficUsersCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.BULK.ALL.RESET_TRAFFIC, 'post', 'Reset user used traffic for all users', { scope: 'bulk-all-reset-traffic', kind: 'write' });
    BulkAllResetTrafficUsersCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            eventSent: zod_1.z.boolean(),
        }),
    });
})(BulkAllResetTrafficUsersCommand || (exports.BulkAllResetTrafficUsersCommand = BulkAllResetTrafficUsersCommand = {}));
