"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetSubpageConfigsTagsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var GetSubpageConfigsTagsCommand;
(function (GetSubpageConfigsTagsCommand) {
    GetSubpageConfigsTagsCommand.url = api_1.REST_API.SUBSCRIPTION_PAGE_CONFIGS.TAGS.GET;
    GetSubpageConfigsTagsCommand.TSQ_url = GetSubpageConfigsTagsCommand.url;
    GetSubpageConfigsTagsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTION_PAGE_CONFIGS_ROUTES.TAGS.GET, 'get', 'Get tags of Subpage Configs', { scope: 'list-tags', kind: 'read' });
    GetSubpageConfigsTagsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            tags: zod_1.z.array(zod_1.z.string()),
        }),
    });
})(GetSubpageConfigsTagsCommand || (exports.GetSubpageConfigsTagsCommand = GetSubpageConfigsTagsCommand = {}));
