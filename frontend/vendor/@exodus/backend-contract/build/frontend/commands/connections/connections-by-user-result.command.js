"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ConnectionsByUserResultCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var ConnectionsByUserResultCommand;
(function (ConnectionsByUserResultCommand) {
    ConnectionsByUserResultCommand.url = api_1.REST_API.CONNECTIONS.CONNECTIONS_BY_USER_RESULT;
    ConnectionsByUserResultCommand.TSQ_url = ConnectionsByUserResultCommand.url(':jobId');
    ConnectionsByUserResultCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.CONNECTIONS_ROUTES.CONNECTIONS_BY_USER_RESULT(':jobId'), 'get', 'Get Connections for User by Job ID', { scope: 'by-user-result', kind: 'read' });
    ConnectionsByUserResultCommand.RequestParamSchema = zod_1.z.object({
        jobId: zod_1.z.string(),
    });
    ConnectionsByUserResultCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            isCompleted: zod_1.z.boolean(),
            isFailed: zod_1.z.boolean(),
            progress: zod_1.z.object({
                total: zod_1.z.number(),
                completed: zod_1.z.number(),
                percent: zod_1.z.number(),
            }),
            result: zod_1.z
                .object({
                success: zod_1.z.boolean(),
                userId: zod_1.z.number(),
                nodes: zod_1.z.array(zod_1.z.object({
                    nodeUuid: zod_1.z.uuid(),
                    nodeName: zod_1.z.string(),
                    countryCode: zod_1.z.string(),
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

})(ConnectionsByUserResultCommand || (exports.ConnectionsByUserResultCommand = ConnectionsByUserResultCommand = {}));
