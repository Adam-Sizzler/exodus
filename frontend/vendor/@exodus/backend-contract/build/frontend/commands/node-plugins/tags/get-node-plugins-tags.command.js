"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetNodePluginsTagsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var GetNodePluginsTagsCommand;
(function (GetNodePluginsTagsCommand) {
    GetNodePluginsTagsCommand.url = api_1.REST_API.NODE_PLUGINS.TAGS.GET;
    GetNodePluginsTagsCommand.TSQ_url = GetNodePluginsTagsCommand.url;
    GetNodePluginsTagsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODE_PLUGINS_ROUTES.TAGS.GET, 'get', 'Get tags of Node Plugins', { scope: 'list-tags', kind: 'read' });
    GetNodePluginsTagsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            tags: zod_1.z.array(zod_1.z.string()),
        }),
    });
})(GetNodePluginsTagsCommand || (exports.GetNodePluginsTagsCommand = GetNodePluginsTagsCommand = {}));
