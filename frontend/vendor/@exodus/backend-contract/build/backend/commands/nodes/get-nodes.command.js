"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetNodesCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var GetNodesCommand;
(function (GetNodesCommand) {
    GetNodesCommand.url = api_1.REST_API.NODES.GET;
    GetNodesCommand.TSQ_url = GetNodesCommand.url;
    GetNodesCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODES_ROUTES.GET, 'get', 'Get nodes', {
        scope: 'list',
        kind: 'read',
    });
    GetNodesCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.array(models_1.NodesSchema),
    });

})(GetNodesCommand || (exports.GetNodesCommand = GetNodesCommand = {}));
