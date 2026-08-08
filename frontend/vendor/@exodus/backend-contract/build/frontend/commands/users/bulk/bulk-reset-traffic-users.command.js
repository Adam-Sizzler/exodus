"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BulkResetTrafficUsersCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var BulkResetTrafficUsersCommand;
(function (BulkResetTrafficUsersCommand) {
    BulkResetTrafficUsersCommand.url = api_1.REST_API.USERS.BULK.RESET_TRAFFIC;
    BulkResetTrafficUsersCommand.TSQ_url = BulkResetTrafficUsersCommand.url;
    BulkResetTrafficUsersCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.BULK.RESET_TRAFFIC, 'post', 'Bulk reset traffic users by User IDs', { scope: 'bulk-reset-traffic', kind: 'write' });
    BulkResetTrafficUsersCommand.RequestBodySchema = zod_1.z.object({
        userIds: zod_1.z.array(zod_1.z.number()).min(1).max(500),
    });

})(BulkResetTrafficUsersCommand || (exports.BulkResetTrafficUsersCommand = BulkResetTrafficUsersCommand = {}));
