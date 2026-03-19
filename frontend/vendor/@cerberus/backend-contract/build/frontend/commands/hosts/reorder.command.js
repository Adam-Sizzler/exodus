"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ReorderHostCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const api_1 = require("../../api");
const models_1 = require("../../models");
var ReorderHostCommand;
(function (ReorderHostCommand) {
    ReorderHostCommand.url = api_1.REST_API.HOSTS.ACTIONS.REORDER;
    ReorderHostCommand.TSQ_url = ReorderHostCommand.url;
    ReorderHostCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.HOSTS_ROUTES.ACTIONS.REORDER, 'post', 'Reorder hosts');
    ReorderHostCommand.RequestSchema = zod_1.z.object({
        hosts: zod_1.z.array(models_1.HostsSchema.pick({
            viewPosition: true,
            uuid: true,
        })),
    });
    ReorderHostCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            isUpdated: zod_1.z.boolean(),
        }),
    });
})(ReorderHostCommand || (exports.ReorderHostCommand = ReorderHostCommand = {}));
