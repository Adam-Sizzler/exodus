"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BulkDeleteUsersCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const api_1 = require("../../../api");
var BulkDeleteUsersCommand;
(function (BulkDeleteUsersCommand) {
    BulkDeleteUsersCommand.url = api_1.REST_API.USERS.BULK.DELETE;
    BulkDeleteUsersCommand.TSQ_url = BulkDeleteUsersCommand.url;
    BulkDeleteUsersCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.BULK.DELETE, 'post', 'Bulk delete users by UUIDs');
    BulkDeleteUsersCommand.RequestSchema = zod_1.z.object({
        uuids: zod_1.z
            .array(zod_1.z.string().uuid())
            .min(1, 'Must be at least 1 user UUID')
            .max(500, 'Maximum 500 user UUIDs'),
    });
    BulkDeleteUsersCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            affectedRows: zod_1.z.number(),
        }),
    });
})(BulkDeleteUsersCommand || (exports.BulkDeleteUsersCommand = BulkDeleteUsersCommand = {}));
