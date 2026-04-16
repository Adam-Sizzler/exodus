"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SetPortToManyHostsCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const api_1 = require("../../../api");
const models_1 = require("../../../models");
var SetPortToManyHostsCommand;
(function (SetPortToManyHostsCommand) {
    SetPortToManyHostsCommand.url = api_1.REST_API.HOSTS.BULK.SET_PORT;
    SetPortToManyHostsCommand.TSQ_url = SetPortToManyHostsCommand.url;
    SetPortToManyHostsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.HOSTS_ROUTES.BULK.SET_PORT, 'post', 'Set port to hosts by UUIDs');
    SetPortToManyHostsCommand.RequestSchema = zod_1.z.object({
        uuids: zod_1.z.array(zod_1.z.string().uuid()),
        port: zod_1.z.number().int().min(1).max(65535),
    });
    SetPortToManyHostsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.array(models_1.HostsSchema),
    });
})(SetPortToManyHostsCommand || (exports.SetPortToManyHostsCommand = SetPortToManyHostsCommand = {}));
