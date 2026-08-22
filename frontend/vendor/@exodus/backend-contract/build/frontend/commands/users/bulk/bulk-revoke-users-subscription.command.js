"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BulkRevokeUsersSubscriptionCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var BulkRevokeUsersSubscriptionCommand;
(function (BulkRevokeUsersSubscriptionCommand) {
    BulkRevokeUsersSubscriptionCommand.url = api_1.REST_API.USERS.BULK.REVOKE_SUBSCRIPTION;
    BulkRevokeUsersSubscriptionCommand.TSQ_url = BulkRevokeUsersSubscriptionCommand.url;
    BulkRevokeUsersSubscriptionCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.BULK.REVOKE_SUBSCRIPTION, 'post', 'Revoke users subscription by User IDs', { scope: 'bulk-revoke-subscription', kind: 'write' });
    BulkRevokeUsersSubscriptionCommand.RequestBodySchema = zod_1.z.object({
        userIds: zod_1.z.array(zod_1.z.number()).min(1).max(500),
    });
})(BulkRevokeUsersSubscriptionCommand || (exports.BulkRevokeUsersSubscriptionCommand = BulkRevokeUsersSubscriptionCommand = {}));
