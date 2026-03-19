"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SetInboundToManyHostsCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const api_1 = require("../../../api");
const models_1 = require("../../../models");
var SetInboundToManyHostsCommand;
(function (SetInboundToManyHostsCommand) {
    SetInboundToManyHostsCommand.url = api_1.REST_API.HOSTS.BULK.SET_INBOUND;
    SetInboundToManyHostsCommand.TSQ_url = SetInboundToManyHostsCommand.url;
    SetInboundToManyHostsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.HOSTS_ROUTES.BULK.SET_INBOUND, 'post', 'Set inbound to hosts by UUIDs');
    SetInboundToManyHostsCommand.RequestSchema = zod_1.z.object({
        uuids: zod_1.z.array(zod_1.z.string().uuid()),
        configProfileUuid: zod_1.z.string().uuid(),
        configProfileInboundUuid: zod_1.z.string().uuid(),
    });
    SetInboundToManyHostsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.array(models_1.HostsSchema),
    });
})(SetInboundToManyHostsCommand || (exports.SetInboundToManyHostsCommand = SetInboundToManyHostsCommand = {}));
