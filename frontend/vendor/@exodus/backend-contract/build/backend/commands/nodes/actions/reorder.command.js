"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ReorderNodeCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const api_1 = require("../../../api");
const models_1 = require("../../../models");
var ReorderNodeCommand;
(function (ReorderNodeCommand) {
    ReorderNodeCommand.url = api_1.REST_API.NODES.ACTIONS.REORDER;
    ReorderNodeCommand.TSQ_url = ReorderNodeCommand.url;
    ReorderNodeCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODES_ROUTES.ACTIONS.REORDER, 'post', 'Reorder nodes');
    ReorderNodeCommand.RequestSchema = zod_1.z.object({
        nodes: zod_1.z.array(models_1.NodesSchema.pick({
            viewPosition: true,
            uuid: true,
        })),
    });
    ReorderNodeCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.array(models_1.NodesSchema),
    });
})(ReorderNodeCommand || (exports.ReorderNodeCommand = ReorderNodeCommand = {}));
