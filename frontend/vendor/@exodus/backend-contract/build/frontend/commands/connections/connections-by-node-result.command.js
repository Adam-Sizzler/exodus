"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ConnectionsByNodeResultCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var ConnectionsByNodeResultCommand;
(function (ConnectionsByNodeResultCommand) {
    ConnectionsByNodeResultCommand.url = api_1.REST_API.CONNECTIONS.CONNECTIONS_BY_NODE_RESULT;
    ConnectionsByNodeResultCommand.TSQ_url = ConnectionsByNodeResultCommand.url(':jobId');
    ConnectionsByNodeResultCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.CONNECTIONS_ROUTES.CONNECTIONS_BY_NODE_RESULT(':jobId'), 'get', 'Get Connections for Node by Job ID', { scope: 'by-node-result', kind: 'read' });
    ConnectionsByNodeResultCommand.RequestParamSchema = zod_1.z.object({
        jobId: zod_1.z.string(),
    });
    ConnectionsByNodeResultCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            isCompleted: zod_1.z.boolean(),
            isFailed: zod_1.z.boolean(),
            result: zod_1.z
                .object({
                success: zod_1.z.boolean(),
                nodeUuid: zod_1.z.uuid(),
                users: zod_1.z.array(zod_1.z.object({
                    userId: zod_1.z.number(),
                    ips: zod_1.z.array(zod_1.z.object({
                        ip: zod_1.z.string(),
                        lastSeen: zod_1.z.iso
                            .datetime({
                            local: true,
                            offset: true,
                        })
                            .transform((str) => new Date(str)),
                    })),
                })),
            })
                .nullable(),
        }),
    });
})(ConnectionsByNodeResultCommand || (exports.ConnectionsByNodeResultCommand = ConnectionsByNodeResultCommand = {}));
