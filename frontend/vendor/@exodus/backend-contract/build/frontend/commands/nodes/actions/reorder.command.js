"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ReorderNodesCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
var ReorderNodesCommand;
(function (ReorderNodesCommand) {
    ReorderNodesCommand.url = api_1.REST_API.NODES.ACTIONS.REORDER;
    ReorderNodesCommand.TSQ_url = ReorderNodesCommand.url;
    ReorderNodesCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODES_ROUTES.ACTIONS.REORDER, 'post', 'Reorder nodes', { scope: 'reorder', kind: 'write' });
    ReorderNodesCommand.RequestBodySchema = zod_1.z.object({
        nodes: zod_1.z.array(models_1.NodesSchema.pick({
            viewPosition: true,
            uuid: true,
        })),
    });
    ReorderNodesCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.array(models_1.NodesSchema),
    });
})(ReorderNodesCommand || (exports.ReorderNodesCommand = ReorderNodesCommand = {}));
