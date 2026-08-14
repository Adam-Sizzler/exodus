"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BulkExtendExpirationDateCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var BulkExtendExpirationDateCommand;
(function (BulkExtendExpirationDateCommand) {
    BulkExtendExpirationDateCommand.url = api_1.REST_API.USERS.BULK.EXTEND_EXPIRATION_DATE;
    BulkExtendExpirationDateCommand.TSQ_url = BulkExtendExpirationDateCommand.url;
    BulkExtendExpirationDateCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.BULK.EXTEND_EXPIRATION_DATE, 'post', 'Extend expiration date for specified users by days', { scope: 'bulk-extend-expiration-date', kind: 'write' });
    BulkExtendExpirationDateCommand.RequestBodySchema = zod_1.z.object({
        userIds: zod_1.z.array(zod_1.z.number()).min(1).max(500),
        extendDays: zod_1.z.int().min(1).max(9999),
    });
})(BulkExtendExpirationDateCommand || (exports.BulkExtendExpirationDateCommand = BulkExtendExpirationDateCommand = {}));
