"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UpdateSharedListCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
var UpdateSharedListCommand;
(function (UpdateSharedListCommand) {
    UpdateSharedListCommand.url = api_1.REST_API.NODE_PLUGINS.SHARED_LISTS.UPDATE;
    UpdateSharedListCommand.TSQ_url = UpdateSharedListCommand.url;
    UpdateSharedListCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODE_PLUGINS_ROUTES.SHARED_LISTS.UPDATE, 'patch', 'Update Shared List', { scope: 'shared-lists-update', kind: 'write' });
    UpdateSharedListCommand.RequestBodySchema = zod_1.z.object({
        name: models_1.SharedListNameSchema,
        config: zod_1.z.record(zod_1.z.string(), zod_1.z.unknown()),
    });
    UpdateSharedListCommand.ResponseSchema = zod_1.z.object({
        response: models_1.SharedListsSchema,
    });
})(UpdateSharedListCommand || (exports.UpdateSharedListCommand = UpdateSharedListCommand = {}));
