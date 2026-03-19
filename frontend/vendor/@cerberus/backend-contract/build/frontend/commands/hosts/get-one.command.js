"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetOneHostCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const api_1 = require("../../api");
const models_1 = require("../../models");
var GetOneHostCommand;
(function (GetOneHostCommand) {
    GetOneHostCommand.url = api_1.REST_API.HOSTS.GET_BY_UUID;
    GetOneHostCommand.TSQ_url = GetOneHostCommand.url(':uuid');
    GetOneHostCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.HOSTS_ROUTES.GET_BY_UUID(':uuid'), 'get', 'Get a host by UUID');
    GetOneHostCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    GetOneHostCommand.ResponseSchema = zod_1.z.object({
        response: models_1.HostsSchema,
    });
})(GetOneHostCommand || (exports.GetOneHostCommand = GetOneHostCommand = {}));
