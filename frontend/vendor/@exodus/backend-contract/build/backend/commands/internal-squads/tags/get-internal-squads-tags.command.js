"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetInternalSquadsTagsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var GetInternalSquadsTagsCommand;
(function (GetInternalSquadsTagsCommand) {
    GetInternalSquadsTagsCommand.url = api_1.REST_API.INTERNAL_SQUADS.TAGS.GET;
    GetInternalSquadsTagsCommand.TSQ_url = GetInternalSquadsTagsCommand.url;
    GetInternalSquadsTagsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INTERNAL_SQUADS_ROUTES.TAGS.GET, 'get', 'Get tags of Internal Squads', { scope: 'list-tags', kind: 'read' });
    GetInternalSquadsTagsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            tags: zod_1.z.array(zod_1.z.string()),
        }),
    });
})(GetInternalSquadsTagsCommand || (exports.GetInternalSquadsTagsCommand = GetInternalSquadsTagsCommand = {}));
