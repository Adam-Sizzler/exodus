"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ReorderNodePluginCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
var ReorderNodePluginCommand;
(function (ReorderNodePluginCommand) {
    ReorderNodePluginCommand.url = api_1.REST_API.NODE_PLUGINS.ACTIONS.REORDER;
    ReorderNodePluginCommand.TSQ_url = ReorderNodePluginCommand.url;
    ReorderNodePluginCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODE_PLUGINS_ROUTES.ACTIONS.REORDER, 'post', 'Reorder Node Plugins', { scope: 'reorder', kind: 'write' });
    ReorderNodePluginCommand.RequestBodySchema = zod_1.z.object({
        items: zod_1.z.array(models_1.NodePluginSchema.pick({
            viewPosition: true,
            uuid: true,
        })),
    });
    ReorderNodePluginCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            total: zod_1.z.number(),
            nodePlugins: zod_1.z.array(models_1.NodePluginSchema),
        }),
    });

})(ReorderNodePluginCommand || (exports.ReorderNodePluginCommand = ReorderNodePluginCommand = {}));

