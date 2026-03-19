"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BulkRevokeUsersSubscriptionCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const api_1 = require("../../../api");
var BulkRevokeUsersSubscriptionCommand;
(function (BulkRevokeUsersSubscriptionCommand) {
    BulkRevokeUsersSubscriptionCommand.url = api_1.REST_API.USERS.BULK.REVOKE_SUBSCRIPTION;
    BulkRevokeUsersSubscriptionCommand.TSQ_url = BulkRevokeUsersSubscriptionCommand.url;
    BulkRevokeUsersSubscriptionCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.BULK.REVOKE_SUBSCRIPTION, 'post', 'Revoke users subscription by User UUIDs');
    BulkRevokeUsersSubscriptionCommand.RequestSchema = zod_1.z.object({
        uuids: zod_1.z
            .array(zod_1.z.string().uuid())
            .min(1, 'Must be at least 1 user UUID')
            .max(500, 'Maximum 500 user UUIDs'),
    });
    BulkRevokeUsersSubscriptionCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            affectedRows: zod_1.z.number(),
        }),
    });
})(BulkRevokeUsersSubscriptionCommand || (exports.BulkRevokeUsersSubscriptionCommand = BulkRevokeUsersSubscriptionCommand = {}));
