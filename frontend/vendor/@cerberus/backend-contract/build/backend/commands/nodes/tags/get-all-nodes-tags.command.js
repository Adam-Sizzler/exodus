"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetAllNodesTagsCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const api_1 = require("../../../api");
var GetAllNodesTagsCommand;
(function (GetAllNodesTagsCommand) {
    GetAllNodesTagsCommand.url = api_1.REST_API.NODES.TAGS.GET;
    GetAllNodesTagsCommand.TSQ_url = GetAllNodesTagsCommand.url;
    GetAllNodesTagsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODES_ROUTES.TAGS.GET, 'get', 'Get all existing nodes tags');
    GetAllNodesTagsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            tags: zod_1.z.array(zod_1.z.string()),
        }),
    });
})(GetAllNodesTagsCommand || (exports.GetAllNodesTagsCommand = GetAllNodesTagsCommand = {}));
