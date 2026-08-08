"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CloneNodePluginCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
var CloneNodePluginCommand;
(function (CloneNodePluginCommand) {
    CloneNodePluginCommand.url = api_1.REST_API.NODE_PLUGINS.ACTIONS.CLONE;
    CloneNodePluginCommand.TSQ_url = CloneNodePluginCommand.url;
    CloneNodePluginCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODE_PLUGINS_ROUTES.ACTIONS.CLONE, 'post', 'Clone Node Plugin', { scope: 'clone', kind: 'write' });
    CloneNodePluginCommand.RequestBodySchema = zod_1.z.object({
        cloneFromUuid: zod_1.z.uuid(),
    });
    CloneNodePluginCommand.ResponseSchema = zod_1.z.object({
        response: models_1.NodePluginSchema.extend({
            pluginConfig: zod_1.z.unknown(),
        }),
    });

})(CloneNodePluginCommand || (exports.CloneNodePluginCommand = CloneNodePluginCommand = {}));
