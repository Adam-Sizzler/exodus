"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SyncNodePluginCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var SyncNodePluginCommand;
(function (SyncNodePluginCommand) {
    SyncNodePluginCommand.url = api_1.REST_API.NODE_PLUGINS.ACTIONS.SYNC;
    SyncNodePluginCommand.TSQ_url = SyncNodePluginCommand.url;
    SyncNodePluginCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODE_PLUGINS_ROUTES.ACTIONS.SYNC, 'post', 'Sync Node Plugin to nodes', { scope: 'sync', kind: 'write' }, 'Push the current plugin config, including referenced shared lists, to every connected node this plugin is active on.');
    SyncNodePluginCommand.RequestBodySchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
    });
})(SyncNodePluginCommand || (exports.SyncNodePluginCommand = SyncNodePluginCommand = {}));
