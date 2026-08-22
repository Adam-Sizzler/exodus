"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetSharedListsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
var GetSharedListsCommand;
(function (GetSharedListsCommand) {
    GetSharedListsCommand.url = api_1.REST_API.NODE_PLUGINS.SHARED_LISTS.GET_ALL;
    GetSharedListsCommand.TSQ_url = GetSharedListsCommand.url;
    GetSharedListsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODE_PLUGINS_ROUTES.SHARED_LISTS.GET_ALL, 'get', 'Get Shared Lists (Preview)', { scope: 'shared-lists-list', kind: 'read' }, 'Returns only the name, type and item count of every shared list. Use "Get Shared List by name" to fetch the items themselves.');
    GetSharedListsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            total: zod_1.z.number(),
            sharedLists: zod_1.z.array(models_1.SharedListPreviewSchema),
        }),
    });
})(GetSharedListsCommand || (exports.GetSharedListsCommand = GetSharedListsCommand = {}));
