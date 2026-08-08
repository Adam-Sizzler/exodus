"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetHostCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const host_response_1 = require("./host.response");
var GetHostCommand;
(function (GetHostCommand) {
    GetHostCommand.url = api_1.REST_API.HOSTS.GET_BY_UUID;
    GetHostCommand.TSQ_url = GetHostCommand.url(':uuid');
    GetHostCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.HOSTS_ROUTES.GET_BY_UUID(':uuid'), 'get', 'Get a host by UUID', { scope: 'get', kind: 'read' });
    GetHostCommand.RequestParamSchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
    });
    GetHostCommand.ResponseSchema = host_response_1.HostResponseSchema;

})(GetHostCommand || (exports.GetHostCommand = GetHostCommand = {}));
