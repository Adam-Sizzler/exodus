"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SetNodePluginTagsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
var SetNodePluginTagsCommand;
(function (SetNodePluginTagsCommand) {
    SetNodePluginTagsCommand.url = api_1.REST_API.NODE_PLUGINS.TAGS.SET;
    SetNodePluginTagsCommand.TSQ_url = SetNodePluginTagsCommand.url;
    SetNodePluginTagsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODE_PLUGINS_ROUTES.TAGS.SET, 'patch', 'Set tags of Node Plugin', { scope: 'set-tags', kind: 'write' });
    SetNodePluginTagsCommand.RequestBodySchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
        tags: models_1.TagsSchema,
    });
    SetNodePluginTagsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            uuid: zod_1.z.uuid(),
            tags: zod_1.z.array(zod_1.z.string()),
        }),
    });
})(SetNodePluginTagsCommand || (exports.SetNodePluginTagsCommand = SetNodePluginTagsCommand = {}));
