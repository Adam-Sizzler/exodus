"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DeleteNodeCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var DeleteNodeCommand;
(function (DeleteNodeCommand) {
    DeleteNodeCommand.url = api_1.REST_API.NODES.DELETE;
    DeleteNodeCommand.TSQ_url = DeleteNodeCommand.url(':uuid');
    DeleteNodeCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODES_ROUTES.DELETE(':uuid'), 'delete', 'Delete a node', { scope: 'delete', kind: 'write' });
    DeleteNodeCommand.RequestParamSchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
    });

})(DeleteNodeCommand || (exports.DeleteNodeCommand = DeleteNodeCommand = {}));
