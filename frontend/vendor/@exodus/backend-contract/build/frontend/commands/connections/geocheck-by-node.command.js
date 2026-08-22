"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GeocheckByNodeCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var GeocheckByNodeCommand;
(function (GeocheckByNodeCommand) {
    GeocheckByNodeCommand.url = api_1.REST_API.CONNECTIONS.GEOCHECK_BY_NODE;
    GeocheckByNodeCommand.TSQ_url = GeocheckByNodeCommand.url(':nodeUuid');
    GeocheckByNodeCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.CONNECTIONS_ROUTES.GEOCHECK_BY_NODE(':nodeUuid'), 'post', 'Request Geocheck for Node', { scope: 'geocheck', kind: 'read' }, 'Queues a geocheck on the node and returns a job ID. Poll "Get Geocheck for Node by Job ID" for the result, the node may take up to a minute to answer.');
    GeocheckByNodeCommand.RequestParamSchema = zod_1.z.object({
        nodeUuid: zod_1.z.uuid().describe('Node UUID'),
    });
    GeocheckByNodeCommand.RequestBodySchema = zod_1.z
        .object({
        ip: zod_1.z.string().optional().describe('Check from this IP address'),
        interface: zod_1.z.string().optional().describe('Check from this network interface'),
    })
        .refine((data) => !(data.ip && data.interface), {
        message: 'Only one of "ip" or "interface" can be specified',
        path: ['ip'],
    });
    GeocheckByNodeCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            jobId: zod_1.z.string(),
        }),
    });
})(GeocheckByNodeCommand || (exports.GeocheckByNodeCommand = GeocheckByNodeCommand = {}));
