"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UpdateNodePluginCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var UpdateNodePluginCommand;
(function (UpdateNodePluginCommand) {
    UpdateNodePluginCommand.url = api_1.REST_API.NODE_PLUGINS.UPDATE;
    UpdateNodePluginCommand.TSQ_url = UpdateNodePluginCommand.url;
    UpdateNodePluginCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODE_PLUGINS_ROUTES.UPDATE, 'patch', 'Update Node Plugin', { scope: 'update', kind: 'write' });
    UpdateNodePluginCommand.RequestBodySchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
        name: zod_1.z
            .string()
            .min(2, 'Name must be at least 2 characters')
            .max(30, 'Name must be less than 30 characters')
            .regex(/^[A-Za-z0-9_\s-]+$/, 'Name can only contain letters, numbers, underscores, dashes and spaces')
            .optional(),
        pluginConfig: zod_1.z.optional(zod_1.z.unknown()),
    });
    UpdateNodePluginCommand.ResponseSchema = zod_1.z.object({
        response: models_1.NodePluginSchema.extend({
            pluginConfig: zod_1.z.unknown(),
        }),
    });

})(UpdateNodePluginCommand || (exports.UpdateNodePluginCommand = UpdateNodePluginCommand = {}));
