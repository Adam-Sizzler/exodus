"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CreateSharedListCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
var CreateSharedListCommand;
(function (CreateSharedListCommand) {
    CreateSharedListCommand.url = api_1.REST_API.NODE_PLUGINS.SHARED_LISTS.CREATE;
    CreateSharedListCommand.TSQ_url = CreateSharedListCommand.url;
    CreateSharedListCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODE_PLUGINS_ROUTES.SHARED_LISTS.CREATE, 'post', 'Create Shared List', { scope: 'shared-lists-create', kind: 'write' });
    CreateSharedListCommand.RequestBodySchema = zod_1.z.object({
        name: models_1.SharedListNameSchema,
        config: zod_1.z.record(zod_1.z.string(), zod_1.z.unknown()),
    });
    CreateSharedListCommand.ResponseSchema = zod_1.z.object({
        response: models_1.SharedListsSchema,
    });
})(CreateSharedListCommand || (exports.CreateSharedListCommand = CreateSharedListCommand = {}));
