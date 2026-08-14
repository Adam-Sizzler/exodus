"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DeleteNodePluginCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var DeleteNodePluginCommand;
(function (DeleteNodePluginCommand) {
    DeleteNodePluginCommand.url = api_1.REST_API.NODE_PLUGINS.DELETE;
    DeleteNodePluginCommand.TSQ_url = DeleteNodePluginCommand.url(':uuid');
    DeleteNodePluginCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODE_PLUGINS_ROUTES.DELETE(':uuid'), 'delete', 'Delete Node Plugin', { scope: 'delete', kind: 'write' });
    DeleteNodePluginCommand.RequestParamSchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
    });
})(DeleteNodePluginCommand || (exports.DeleteNodePluginCommand = DeleteNodePluginCommand = {}));
