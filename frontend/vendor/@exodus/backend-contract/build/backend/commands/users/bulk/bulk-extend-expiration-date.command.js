"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BulkExtendExpirationDateCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const api_1 = require("../../../api");
var BulkExtendExpirationDateCommand;
(function (BulkExtendExpirationDateCommand) {
    BulkExtendExpirationDateCommand.url = api_1.REST_API.USERS.BULK.EXTEND_EXPIRATION_DATE;
    BulkExtendExpirationDateCommand.TSQ_url = BulkExtendExpirationDateCommand.url;
    BulkExtendExpirationDateCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.BULK.EXTEND_EXPIRATION_DATE, 'post', 'Extend expiration date for specified users by days');
    BulkExtendExpirationDateCommand.RequestSchema = zod_1.z.object({
        uuids: zod_1.z
            .array(zod_1.z.string().uuid())
            .min(1, 'Must be at least 1 user UUID')
            .max(500, 'Maximum 500 user UUIDs'),
        extendDays: zod_1.z
            .number()
            .int()
            .min(1, 'Extend days must be greater than 0')
            .max(9999, 'Maximum 9999 days'),
    });
    BulkExtendExpirationDateCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            affectedRows: zod_1.z.number(),
        }),
    });
})(BulkExtendExpirationDateCommand || (exports.BulkExtendExpirationDateCommand = BulkExtendExpirationDateCommand = {}));
