"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BulkDeleteUsersCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var BulkDeleteUsersCommand;
(function (BulkDeleteUsersCommand) {
    BulkDeleteUsersCommand.url = api_1.REST_API.USERS.BULK.DELETE;
    BulkDeleteUsersCommand.TSQ_url = BulkDeleteUsersCommand.url;
    BulkDeleteUsersCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.BULK.DELETE, 'post', 'Bulk delete users by User IDs', { scope: 'bulk-delete', kind: 'write' });
    BulkDeleteUsersCommand.RequestBodySchema = zod_1.z.object({
        userIds: zod_1.z.array(zod_1.z.number()).min(1).max(500),
    });
})(BulkDeleteUsersCommand || (exports.BulkDeleteUsersCommand = BulkDeleteUsersCommand = {}));
