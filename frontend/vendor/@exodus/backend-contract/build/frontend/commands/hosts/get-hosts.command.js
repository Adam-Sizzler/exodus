"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetHostsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var GetHostsCommand;
(function (GetHostsCommand) {
    GetHostsCommand.url = api_1.REST_API.HOSTS.GET;
    GetHostsCommand.TSQ_url = GetHostsCommand.url;
    GetHostsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.HOSTS_ROUTES.GET, 'get', 'Get hosts', {
        scope: 'list',
        kind: 'read',
    });
    GetHostsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.array(models_1.HostsSchema),
    });

})(GetHostsCommand || (exports.GetHostsCommand = GetHostsCommand = {}));
