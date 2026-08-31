"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DeleteSharedListCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
var DeleteSharedListCommand;
(function (DeleteSharedListCommand) {
    DeleteSharedListCommand.url = api_1.REST_API.NODE_PLUGINS.SHARED_LISTS.DELETE;
    DeleteSharedListCommand.TSQ_url = DeleteSharedListCommand.url;
    DeleteSharedListCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODE_PLUGINS_ROUTES.SHARED_LISTS.DELETE, 'delete', 'Delete Shared List by name', { scope: 'shared-lists-delete', kind: 'write' });
    DeleteSharedListCommand.RequestBodySchema = zod_1.z.object({
        name: models_1.SharedListNameSchema,
    });
})(DeleteSharedListCommand || (exports.DeleteSharedListCommand = DeleteSharedListCommand = {}));
