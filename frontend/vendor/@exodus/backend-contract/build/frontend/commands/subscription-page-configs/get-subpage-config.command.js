"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetSubpageConfigCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var GetSubpageConfigCommand;
(function (GetSubpageConfigCommand) {
    GetSubpageConfigCommand.url = api_1.REST_API.SUBSCRIPTION_PAGE_CONFIGS.GET;
    GetSubpageConfigCommand.TSQ_url = GetSubpageConfigCommand.url(':uuid');
    GetSubpageConfigCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTION_PAGE_CONFIGS_ROUTES.GET(':uuid'), 'get', 'Get subscription page config by uuid', { scope: 'get', kind: 'read' });
    GetSubpageConfigCommand.RequestParamSchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
    });
    GetSubpageConfigCommand.ResponseSchema = zod_1.z.object({
        response: models_1.SubscriptionPageConfigSchema.extend({
            config: zod_1.z.unknown(),
        }),
    });
})(GetSubpageConfigCommand || (exports.GetSubpageConfigCommand = GetSubpageConfigCommand = {}));
