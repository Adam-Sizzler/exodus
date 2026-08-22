"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetSharedListCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
var GetSharedListCommand;
(function (GetSharedListCommand) {
    GetSharedListCommand.url = api_1.REST_API.NODE_PLUGINS.SHARED_LISTS.GET;
    GetSharedListCommand.TSQ_url = GetSharedListCommand.url(':name');
    GetSharedListCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODE_PLUGINS_ROUTES.SHARED_LISTS.GET(':name'), 'get', 'Get Shared List by name', { scope: 'shared-lists-get', kind: 'read' });
    GetSharedListCommand.RequestParamSchema = zod_1.z.object({
        name: models_1.SharedListNameSchema,
    });
    GetSharedListCommand.ResponseSchema = zod_1.z.object({
        response: models_1.SharedListsSchema,
    });
})(GetSharedListCommand || (exports.GetSharedListCommand = GetSharedListCommand = {}));
