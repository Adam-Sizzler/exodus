"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.FetchIpsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var FetchIpsCommand;
(function (FetchIpsCommand) {
    FetchIpsCommand.url = api_1.REST_API.IP_CONTROL.FETCH_IPS;
    FetchIpsCommand.TSQ_url = FetchIpsCommand.url(':uuid');
    FetchIpsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.IP_CONTROL_ROUTES.FETCH_IPS(':uuid'), 'post', 'Request IP List for User', { scope: 'fetch-ips', kind: 'read' });
    FetchIpsCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    FetchIpsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            jobId: zod_1.z.string(),
        }),
    });
})(FetchIpsCommand || (exports.FetchIpsCommand = FetchIpsCommand = {}));
