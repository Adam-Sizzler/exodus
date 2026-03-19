"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BulkUpdateUsersSquadsCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const api_1 = require("../../../api");
var BulkUpdateUsersSquadsCommand;
(function (BulkUpdateUsersSquadsCommand) {
    BulkUpdateUsersSquadsCommand.url = api_1.REST_API.USERS.BULK.UPDATE_SQUADS;
    BulkUpdateUsersSquadsCommand.TSQ_url = BulkUpdateUsersSquadsCommand.url;
    BulkUpdateUsersSquadsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.BULK.UPDATE_SQUADS, 'post', 'Bulk update users internal squads by UUIDs');
    BulkUpdateUsersSquadsCommand.RequestSchema = zod_1.z.object({
        uuids: zod_1.z
            .array(zod_1.z.string().uuid())
            .min(1, 'Must be at least 1 user UUID')
            .max(500, 'Maximum 500 user UUIDs'),
        activeInternalSquads: zod_1.z.array(zod_1.z.string().uuid(), {
            invalid_type_error: 'Enabled internal squads must be an array of UUIDs',
        }),
    });
    BulkUpdateUsersSquadsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            affectedRows: zod_1.z.number(),
        }),
    });
})(BulkUpdateUsersSquadsCommand || (exports.BulkUpdateUsersSquadsCommand = BulkUpdateUsersSquadsCommand = {}));
