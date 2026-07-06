"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.FetchUsersIpsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var FetchUsersIpsCommand;
(function (FetchUsersIpsCommand) {
    FetchUsersIpsCommand.url = api_1.REST_API.IP_CONTROL.FETCH_USERS_IPS;
    FetchUsersIpsCommand.TSQ_url = FetchUsersIpsCommand.url(':nodeUuid');
    FetchUsersIpsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.IP_CONTROL_ROUTES.FETCH_USERS_IPS(':nodeUuid'), 'post', 'Request Users IPs List for Node', { scope: 'fetch-users-ips', kind: 'read' });
    FetchUsersIpsCommand.RequestSchema = zod_1.z.object({
        nodeUuid: zod_1.z.string().uuid(),
    });
    FetchUsersIpsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            jobId: zod_1.z.string(),
        }),
    });
})(FetchUsersIpsCommand || (exports.FetchUsersIpsCommand = FetchUsersIpsCommand = {}));
