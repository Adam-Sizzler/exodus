"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DisableNodeCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const node_response_1 = require("../node.response");
var DisableNodeCommand;
(function (DisableNodeCommand) {
    DisableNodeCommand.url = api_1.REST_API.NODES.ACTIONS.DISABLE;
    DisableNodeCommand.TSQ_url = DisableNodeCommand.url(':uuid');
    DisableNodeCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODES_ROUTES.ACTIONS.DISABLE(':uuid'), 'post', 'Disable a node', { scope: 'disable', kind: 'write' });
    DisableNodeCommand.RequestParamSchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
    });
    DisableNodeCommand.ResponseSchema = node_response_1.NodeResponseSchema;
})(DisableNodeCommand || (exports.DisableNodeCommand = DisableNodeCommand = {}));
