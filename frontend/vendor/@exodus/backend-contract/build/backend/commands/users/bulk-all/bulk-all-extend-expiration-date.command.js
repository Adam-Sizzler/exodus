"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BulkAllExtendExpirationDateCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var BulkAllExtendExpirationDateCommand;
(function (BulkAllExtendExpirationDateCommand) {
    BulkAllExtendExpirationDateCommand.url = api_1.REST_API.USERS.BULK.ALL.EXTEND_EXPIRATION_DATE;
    BulkAllExtendExpirationDateCommand.TSQ_url = BulkAllExtendExpirationDateCommand.url;
    BulkAllExtendExpirationDateCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.BULK.ALL.EXTEND_EXPIRATION_DATE, 'post', 'Extend expiration date for all users by days', { scope: 'bulk-all-extend-expiration-date', kind: 'write' });
    BulkAllExtendExpirationDateCommand.RequestBodySchema = zod_1.z.object({
        extendDays: zod_1.z.int().min(1),
    });
})(BulkAllExtendExpirationDateCommand || (exports.BulkAllExtendExpirationDateCommand = BulkAllExtendExpirationDateCommand = {}));
