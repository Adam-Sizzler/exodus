"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BulkNodesActionsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var BulkNodesActionsCommand;
(function (BulkNodesActionsCommand) {
    BulkNodesActionsCommand.url = api_1.REST_API.NODES.BULK_ACTIONS.ACTIONS;
    BulkNodesActionsCommand.TSQ_url = BulkNodesActionsCommand.url;
    BulkNodesActionsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODES_ROUTES.BULK_ACTIONS.ACTIONS, 'post', 'Perform actions for many nodes', { scope: 'bulk-actions', kind: 'write' });
    BulkNodesActionsCommand.RequestSchema = zod_1.z.object({
        uuids: zod_1.z.array(zod_1.z.string().uuid()).min(1, 'Must be at least 1 Node UUID'),
        action: zod_1.z.nativeEnum(constants_1.NODES_BULK_ACTIONS),
    });
    BulkNodesActionsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            eventSent: zod_1.z.boolean(),
        }),
    });
})(BulkNodesActionsCommand || (exports.BulkNodesActionsCommand = BulkNodesActionsCommand = {}));
