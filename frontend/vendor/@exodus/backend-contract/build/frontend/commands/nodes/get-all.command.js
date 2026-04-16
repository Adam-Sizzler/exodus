"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetAllNodesCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const api_1 = require("../../api");
const models_1 = require("../../models");
var GetAllNodesCommand;
(function (GetAllNodesCommand) {
    GetAllNodesCommand.url = api_1.REST_API.NODES.GET;
    GetAllNodesCommand.TSQ_url = GetAllNodesCommand.url;
    GetAllNodesCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODES_ROUTES.GET, 'get', 'Get all nodes');
    GetAllNodesCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.array(models_1.NodesSchema),
    });
})(GetAllNodesCommand || (exports.GetAllNodesCommand = GetAllNodesCommand = {}));
