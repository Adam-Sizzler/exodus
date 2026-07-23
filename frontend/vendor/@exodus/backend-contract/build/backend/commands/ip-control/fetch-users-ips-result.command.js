"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.FetchUsersIpsResultCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var FetchUsersIpsResultCommand;
(function (FetchUsersIpsResultCommand) {
    FetchUsersIpsResultCommand.url = api_1.REST_API.IP_CONTROL.GET_FETCH_USERS_IPS_RESULT;
    FetchUsersIpsResultCommand.TSQ_url = FetchUsersIpsResultCommand.url(':jobId');
    FetchUsersIpsResultCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.IP_CONTROL_ROUTES.GET_FETCH_USERS_IPS_RESULT(':jobId'), 'get', 'Get Users IPs List Result by Job ID', { scope: 'fetch-users-ips-result', kind: 'read' });
    FetchUsersIpsResultCommand.RequestSchema = zod_1.z.object({
        jobId: zod_1.z.string(),
    });
    FetchUsersIpsResultCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            isCompleted: zod_1.z.boolean(),
            isFailed: zod_1.z.boolean(),
            result: zod_1.z
                .object({
                success: zod_1.z.boolean(),
                nodeUuid: zod_1.z.string().uuid(),
                users: zod_1.z.array(zod_1.z.object({
                    userId: zod_1.z.string(),
                    ips: zod_1.z.array(zod_1.z.object({
                        ip: zod_1.z.string(),
                        lastSeen: zod_1.z
                            .string()
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
})(FetchUsersIpsResultCommand || (exports.FetchUsersIpsResultCommand = FetchUsersIpsResultCommand = {}));
