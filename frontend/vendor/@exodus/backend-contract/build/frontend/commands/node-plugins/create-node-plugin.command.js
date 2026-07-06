"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CreateNodePluginCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var CreateNodePluginCommand;
(function (CreateNodePluginCommand) {
    CreateNodePluginCommand.url = api_1.REST_API.NODE_PLUGINS.CREATE;
    CreateNodePluginCommand.TSQ_url = CreateNodePluginCommand.url;
    CreateNodePluginCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODE_PLUGINS_ROUTES.CREATE, 'post', 'Create Node Plugin', { scope: 'create', kind: 'write' });
    CreateNodePluginCommand.RequestSchema = zod_1.z.object({
        name: zod_1.z
            .string()
            .min(2, 'Name must be at least 2 characters')
            .max(30, 'Name must be less than 30 characters')
            .regex(/^[A-Za-z0-9_\s-]+$/, 'Name can only contain letters, numbers, underscores, dashes and spaces'),
    });
    CreateNodePluginCommand.ResponseSchema = zod_1.z.object({
        response: models_1.NodePluginSchema,
    });
})(CreateNodePluginCommand || (exports.CreateNodePluginCommand = CreateNodePluginCommand = {}));
