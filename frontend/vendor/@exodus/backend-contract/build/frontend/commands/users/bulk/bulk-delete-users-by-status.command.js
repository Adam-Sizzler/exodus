"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BulkDeleteUsersByStatusCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
var BulkDeleteUsersByStatusCommand;
(function (BulkDeleteUsersByStatusCommand) {
    BulkDeleteUsersByStatusCommand.url = api_1.REST_API.USERS.BULK.DELETE_BY_STATUS;
    BulkDeleteUsersByStatusCommand.TSQ_url = BulkDeleteUsersByStatusCommand.url;
    BulkDeleteUsersByStatusCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.BULK.DELETE_BY_STATUS, 'post', 'Bulk delete users by status', { scope: 'bulk-delete-by-status', kind: 'write' });
    BulkDeleteUsersByStatusCommand.RequestBodySchema = zod_1.z.object({
        status: models_1.UsersSchema.shape.status,
    });
})(BulkDeleteUsersByStatusCommand || (exports.BulkDeleteUsersByStatusCommand = BulkDeleteUsersByStatusCommand = {}));
