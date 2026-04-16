"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BulkResetTrafficUsersCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const api_1 = require("../../../api");
var BulkResetTrafficUsersCommand;
(function (BulkResetTrafficUsersCommand) {
    BulkResetTrafficUsersCommand.url = api_1.REST_API.USERS.BULK.RESET_TRAFFIC;
    BulkResetTrafficUsersCommand.TSQ_url = BulkResetTrafficUsersCommand.url;
    BulkResetTrafficUsersCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.BULK.RESET_TRAFFIC, 'post', 'Bulk reset traffic users by UUIDs');
    BulkResetTrafficUsersCommand.RequestSchema = zod_1.z.object({
        uuids: zod_1.z
            .array(zod_1.z.string().uuid())
            .min(1, 'Must be at least 1 user UUID')
            .max(500, 'Maximum 500 user UUIDs'),
    });
    BulkResetTrafficUsersCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            affectedRows: zod_1.z.number(),
        }),
    });
})(BulkResetTrafficUsersCommand || (exports.BulkResetTrafficUsersCommand = BulkResetTrafficUsersCommand = {}));
