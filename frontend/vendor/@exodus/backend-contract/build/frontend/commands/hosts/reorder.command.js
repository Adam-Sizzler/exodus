"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ReorderHostsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var ReorderHostsCommand;
(function (ReorderHostsCommand) {
    ReorderHostsCommand.url = api_1.REST_API.HOSTS.ACTIONS.REORDER;
    ReorderHostsCommand.TSQ_url = ReorderHostsCommand.url;
    ReorderHostsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.HOSTS_ROUTES.ACTIONS.REORDER, 'post', 'Reorder hosts', { scope: 'reorder', kind: 'write' });
    ReorderHostsCommand.RequestBodySchema = zod_1.z.object({
        hosts: zod_1.z.array(models_1.HostsSchema.pick({
            viewPosition: true,
            uuid: true,
        })),
    });
    ReorderHostsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            isUpdated: zod_1.z.boolean(),
        }),
    });
})(ReorderHostsCommand || (exports.ReorderHostsCommand = ReorderHostsCommand = {}));
