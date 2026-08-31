"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SetSubscriptionTemplateTagsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
var SetSubscriptionTemplateTagsCommand;
(function (SetSubscriptionTemplateTagsCommand) {
    SetSubscriptionTemplateTagsCommand.url = api_1.REST_API.SUBSCRIPTION_TEMPLATE.TAGS.SET;
    SetSubscriptionTemplateTagsCommand.TSQ_url = SetSubscriptionTemplateTagsCommand.url;
    SetSubscriptionTemplateTagsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTION_TEMPLATE_ROUTES.TAGS.SET, 'patch', 'Set tags of Subscription Template', { scope: 'set-tags', kind: 'write' });
    SetSubscriptionTemplateTagsCommand.RequestBodySchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
        tags: models_1.TagsSchema,
    });
    SetSubscriptionTemplateTagsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            uuid: zod_1.z.uuid(),
            tags: zod_1.z.array(zod_1.z.string()),
        }),
    });
})(SetSubscriptionTemplateTagsCommand || (exports.SetSubscriptionTemplateTagsCommand = SetSubscriptionTemplateTagsCommand = {}));
