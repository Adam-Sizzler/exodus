"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetSubpageConfigsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var GetSubpageConfigsCommand;
(function (GetSubpageConfigsCommand) {
    GetSubpageConfigsCommand.url = api_1.REST_API.SUBSCRIPTION_PAGE_CONFIGS.GET_ALL;
    GetSubpageConfigsCommand.TSQ_url = GetSubpageConfigsCommand.url;
    GetSubpageConfigsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTION_PAGE_CONFIGS_ROUTES.GET_ALL, 'get', 'Get all subscription page configs', { scope: 'list', kind: 'read' });
    GetSubpageConfigsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            total: zod_1.z.number(),
            configs: zod_1.z.array(models_1.SubscriptionPageConfigSchema),
        }),
    });
})(GetSubpageConfigsCommand || (exports.GetSubpageConfigsCommand = GetSubpageConfigsCommand = {}));
