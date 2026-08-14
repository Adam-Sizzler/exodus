"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.EnableNodeCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const node_response_1 = require("../node.response");
var EnableNodeCommand;
(function (EnableNodeCommand) {
    EnableNodeCommand.url = api_1.REST_API.NODES.ACTIONS.ENABLE;
    EnableNodeCommand.TSQ_url = EnableNodeCommand.url(':uuid');
    EnableNodeCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODES_ROUTES.ACTIONS.ENABLE(':uuid'), 'post', 'Enable a node', { scope: 'enable', kind: 'write' });
    EnableNodeCommand.RequestParamSchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
    });
    EnableNodeCommand.ResponseSchema = node_response_1.NodeResponseSchema;
})(EnableNodeCommand || (exports.EnableNodeCommand = EnableNodeCommand = {}));
