"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ConnectionsByNodeCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var ConnectionsByNodeCommand;
(function (ConnectionsByNodeCommand) {
    ConnectionsByNodeCommand.url = api_1.REST_API.CONNECTIONS.CONNECTIONS_BY_NODE;
    ConnectionsByNodeCommand.TSQ_url = ConnectionsByNodeCommand.url(':nodeUuid');
    ConnectionsByNodeCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.CONNECTIONS_ROUTES.CONNECTIONS_BY_NODE(':nodeUuid'), 'post', 'Request Connections for Node', { scope: 'by-node', kind: 'read' });
    ConnectionsByNodeCommand.RequestParamSchema = zod_1.z.object({
        nodeUuid: zod_1.z.uuid().describe('Node UUID'),
    });
    ConnectionsByNodeCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            jobId: zod_1.z.string(),
        }),
    });
})(ConnectionsByNodeCommand || (exports.ConnectionsByNodeCommand = ConnectionsByNodeCommand = {}));
