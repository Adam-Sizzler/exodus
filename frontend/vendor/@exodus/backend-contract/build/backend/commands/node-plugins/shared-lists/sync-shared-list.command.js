"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SyncSharedListCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
var SyncSharedListCommand;
(function (SyncSharedListCommand) {
    SyncSharedListCommand.url = api_1.REST_API.NODE_PLUGINS.SHARED_LISTS.ACTIONS.SYNC;
    SyncSharedListCommand.TSQ_url = SyncSharedListCommand.url;
    SyncSharedListCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODE_PLUGINS_ROUTES.SHARED_LISTS.ACTIONS.SYNC, 'post', 'Sync Shared List to nodes', { scope: 'shared-lists-sync', kind: 'write' }, 'Push every plugin referencing this shared list to the nodes it is active on.');
    SyncSharedListCommand.RequestBodySchema = zod_1.z.object({
        name: models_1.SharedListNameSchema,
    });
})(SyncSharedListCommand || (exports.SyncSharedListCommand = SyncSharedListCommand = {}));
