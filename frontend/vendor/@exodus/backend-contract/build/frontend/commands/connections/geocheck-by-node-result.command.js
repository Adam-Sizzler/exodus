"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GeocheckByNodeResultCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var GeocheckByNodeResultCommand;
(function (GeocheckByNodeResultCommand) {
    GeocheckByNodeResultCommand.url = api_1.REST_API.CONNECTIONS.GEOCHECK_BY_NODE_RESULT;
    GeocheckByNodeResultCommand.TSQ_url = GeocheckByNodeResultCommand.url(':jobId');
    GeocheckByNodeResultCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.CONNECTIONS_ROUTES.GEOCHECK_BY_NODE_RESULT(':jobId'), 'get', 'Get Geocheck for Node by Job ID', { scope: 'geocheck-result', kind: 'read' });
    GeocheckByNodeResultCommand.RequestParamSchema = zod_1.z.object({
        jobId: zod_1.z.string(),
    });
    GeocheckByNodeResultCommand.GeocheckImageSchema = zod_1.z.object({
        format: zod_1.z.literal('svg'),
        media_type: zod_1.z.literal('image/svg+xml'),
        encoding: zod_1.z.literal('base64'),
        data: zod_1.z.string().describe('Base64-encoded image, ready for a data: URL'),
    });
    GeocheckByNodeResultCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            isCompleted: zod_1.z.boolean(),
            isFailed: zod_1.z.boolean(),
            result: zod_1.z
                .object({
                success: zod_1.z.boolean(),
                nodeUuid: zod_1.z.uuid(),
                image: GeocheckByNodeResultCommand.GeocheckImageSchema.nullable(),
                rawReport: zod_1.z
                    .record(zod_1.z.string(), zod_1.z.unknown())
                    .nullable()
                    .describe('The full node report with the image object stripped out'),
                message: zod_1.z.string().nullable(),
            })
                .nullable(),
        }),
    });
})(GeocheckByNodeResultCommand || (exports.GeocheckByNodeResultCommand = GeocheckByNodeResultCommand = {}));
