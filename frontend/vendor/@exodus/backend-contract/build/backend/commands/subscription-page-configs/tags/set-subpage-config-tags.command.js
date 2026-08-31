"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SetSubpageConfigTagsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
var SetSubpageConfigTagsCommand;
(function (SetSubpageConfigTagsCommand) {
    SetSubpageConfigTagsCommand.url = api_1.REST_API.SUBSCRIPTION_PAGE_CONFIGS.TAGS.SET;
    SetSubpageConfigTagsCommand.TSQ_url = SetSubpageConfigTagsCommand.url;
    SetSubpageConfigTagsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTION_PAGE_CONFIGS_ROUTES.TAGS.SET, 'patch', 'Set tags of Subpage Config', { scope: 'set-tags', kind: 'write' });
    SetSubpageConfigTagsCommand.RequestBodySchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
        tags: models_1.TagsSchema,
    });
    SetSubpageConfigTagsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            uuid: zod_1.z.uuid(),
            tags: zod_1.z.array(zod_1.z.string()),
        }),
    });
})(SetSubpageConfigTagsCommand || (exports.SetSubpageConfigTagsCommand = SetSubpageConfigTagsCommand = {}));
