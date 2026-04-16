"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DeleteNodeCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const api_1 = require("../../api");
var DeleteNodeCommand;
(function (DeleteNodeCommand) {
    DeleteNodeCommand.url = api_1.REST_API.NODES.DELETE;
    DeleteNodeCommand.TSQ_url = DeleteNodeCommand.url(':uuid');
    DeleteNodeCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODES_ROUTES.DELETE(':uuid'), 'delete', 'Delete a node');
    DeleteNodeCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    DeleteNodeCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            isDeleted: zod_1.z.boolean(),
        }),
    });
})(DeleteNodeCommand || (exports.DeleteNodeCommand = DeleteNodeCommand = {}));
