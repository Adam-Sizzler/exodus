"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetAllHostsCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const api_1 = require("../../api");
const models_1 = require("../../models");
var GetAllHostsCommand;
(function (GetAllHostsCommand) {
    GetAllHostsCommand.url = api_1.REST_API.HOSTS.GET;
    GetAllHostsCommand.TSQ_url = GetAllHostsCommand.url;
    GetAllHostsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.HOSTS_ROUTES.GET, 'get', 'Get all hosts');
    GetAllHostsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.array(models_1.HostsSchema),
    });
})(GetAllHostsCommand || (exports.GetAllHostsCommand = GetAllHostsCommand = {}));
