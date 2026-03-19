"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.FetchIpsResultCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var FetchIpsResultCommand;
(function (FetchIpsResultCommand) {
    FetchIpsResultCommand.url = api_1.REST_API.IP_CONTROL.GET_FETCH_IPS_RESULT;
    FetchIpsResultCommand.TSQ_url = FetchIpsResultCommand.url(':jobId');
    FetchIpsResultCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.IP_CONTROL_ROUTES.GET_FETCH_IPS_RESULT(':jobId'), 'get', 'Get IP List Result by Job ID');
    FetchIpsResultCommand.RequestSchema = zod_1.z.object({
        jobId: zod_1.z.string(),
    });
    FetchIpsResultCommand.ResponseSchema = zod_1.z.object({
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
                userUuid: zod_1.z.string().uuid(),
                userId: zod_1.z.string(),
                nodes: zod_1.z.array(zod_1.z.object({
                    nodeUuid: zod_1.z.string().uuid(),
                    nodeName: zod_1.z.string(),
                    countryCode: zod_1.z.string(),
                    ips: zod_1.z.array(zod_1.z.string()),
                })),
            })
                .nullable(),
        }),
    });
})(FetchIpsResultCommand || (exports.FetchIpsResultCommand = FetchIpsResultCommand = {}));
