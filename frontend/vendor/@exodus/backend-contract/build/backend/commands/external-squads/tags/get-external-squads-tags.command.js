"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetExternalSquadsTagsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var GetExternalSquadsTagsCommand;
(function (GetExternalSquadsTagsCommand) {
    GetExternalSquadsTagsCommand.url = api_1.REST_API.EXTERNAL_SQUADS.TAGS.GET;
    GetExternalSquadsTagsCommand.TSQ_url = GetExternalSquadsTagsCommand.url;
    GetExternalSquadsTagsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.EXTERNAL_SQUADS_ROUTES.TAGS.GET, 'get', 'Get tags of External Squads', { scope: 'list-tags', kind: 'read' });
    GetExternalSquadsTagsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            tags: zod_1.z.array(zod_1.z.string()),
        }),
    });
})(GetExternalSquadsTagsCommand || (exports.GetExternalSquadsTagsCommand = GetExternalSquadsTagsCommand = {}));
