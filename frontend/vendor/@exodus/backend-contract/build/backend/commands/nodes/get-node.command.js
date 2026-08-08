"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetNodeCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const node_response_1 = require("./node.response");
var GetNodeCommand;
(function (GetNodeCommand) {
    GetNodeCommand.url = api_1.REST_API.NODES.GET_BY_UUID;
    GetNodeCommand.TSQ_url = GetNodeCommand.url(':uuid');
    GetNodeCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODES_ROUTES.GET_BY_UUID(':uuid'), 'get', 'Get node by UUID', { scope: 'get', kind: 'read' });
    GetNodeCommand.RequestParamSchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
    });
    GetNodeCommand.ResponseSchema = node_response_1.NodeResponseSchema;
})(GetNodeCommand || (exports.GetNodeCommand = GetNodeCommand = {}));
