"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetNodePluginsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var GetNodePluginsCommand;
(function (GetNodePluginsCommand) {
    GetNodePluginsCommand.url = api_1.REST_API.NODE_PLUGINS.GET_ALL;
    GetNodePluginsCommand.TSQ_url = GetNodePluginsCommand.url;
    GetNodePluginsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODE_PLUGINS_ROUTES.GET_ALL, 'get', 'Get all Node Plugins', { scope: 'list', kind: 'read' });
    GetNodePluginsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            total: zod_1.z.number(),
            nodePlugins: zod_1.z.array(models_1.NodePluginSchema),
        }),
    });

})(GetNodePluginsCommand || (exports.GetNodePluginsCommand = GetNodePluginsCommand = {}));
