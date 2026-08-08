"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetNodesTagsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var GetNodesTagsCommand;
(function (GetNodesTagsCommand) {
    GetNodesTagsCommand.url = api_1.REST_API.NODES.TAGS.GET;
    GetNodesTagsCommand.TSQ_url = GetNodesTagsCommand.url;
    GetNodesTagsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODES_ROUTES.TAGS.GET, 'get', 'Get nodes tags', { scope: 'list-tags', kind: 'read' });
    GetNodesTagsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            tags: zod_1.z.array(zod_1.z.string()),
        }),
    });

})(GetNodesTagsCommand || (exports.GetNodesTagsCommand = GetNodesTagsCommand = {}));
