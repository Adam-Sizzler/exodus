"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetSubscriptionTemplatesTagsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var GetSubscriptionTemplatesTagsCommand;
(function (GetSubscriptionTemplatesTagsCommand) {
    GetSubscriptionTemplatesTagsCommand.url = api_1.REST_API.SUBSCRIPTION_TEMPLATE.TAGS.GET;
    GetSubscriptionTemplatesTagsCommand.TSQ_url = GetSubscriptionTemplatesTagsCommand.url;
    GetSubscriptionTemplatesTagsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTION_TEMPLATE_ROUTES.TAGS.GET, 'get', 'Get tags of Subscription Templates', { scope: 'list-tags', kind: 'read' });
    GetSubscriptionTemplatesTagsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            tags: zod_1.z.array(zod_1.z.string()),
        }),
    });
})(GetSubscriptionTemplatesTagsCommand || (exports.GetSubscriptionTemplatesTagsCommand = GetSubscriptionTemplatesTagsCommand = {}));
