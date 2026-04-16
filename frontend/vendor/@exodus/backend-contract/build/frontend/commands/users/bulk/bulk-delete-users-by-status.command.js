"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BulkDeleteUsersByStatusCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const api_1 = require("../../../api");
const models_1 = require("../../../models");
var BulkDeleteUsersByStatusCommand;
(function (BulkDeleteUsersByStatusCommand) {
    BulkDeleteUsersByStatusCommand.url = api_1.REST_API.USERS.BULK.DELETE_BY_STATUS;
    BulkDeleteUsersByStatusCommand.TSQ_url = BulkDeleteUsersByStatusCommand.url;
    BulkDeleteUsersByStatusCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.BULK.DELETE_BY_STATUS, 'post', 'Bulk delete users by status');
    BulkDeleteUsersByStatusCommand.RequestSchema = zod_1.z.object({
        status: models_1.UsersSchema.shape.status,
    });
    BulkDeleteUsersByStatusCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            affectedRows: zod_1.z.number(),
        }),
    });
})(BulkDeleteUsersByStatusCommand || (exports.BulkDeleteUsersByStatusCommand = BulkDeleteUsersByStatusCommand = {}));
