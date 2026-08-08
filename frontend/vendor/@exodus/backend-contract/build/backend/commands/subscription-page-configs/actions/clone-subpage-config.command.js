"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CloneSubpageConfigCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
var CloneSubpageConfigCommand;
(function (CloneSubpageConfigCommand) {
    CloneSubpageConfigCommand.url = api_1.REST_API.SUBSCRIPTION_PAGE_CONFIGS.ACTIONS.CLONE;
    CloneSubpageConfigCommand.TSQ_url = CloneSubpageConfigCommand.url;
    CloneSubpageConfigCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTION_PAGE_CONFIGS_ROUTES.ACTIONS.CLONE, 'post', 'Clone subscription page config', { scope: 'clone', kind: 'write' });
    CloneSubpageConfigCommand.RequestBodySchema = zod_1.z.object({
        cloneFromUuid: zod_1.z.uuid(),
    });
    CloneSubpageConfigCommand.ResponseSchema = zod_1.z.object({
        response: models_1.SubscriptionPageConfigSchema.extend({
            config: zod_1.z.unknown(),
        }),
    });

})(CloneSubpageConfigCommand || (exports.CloneSubpageConfigCommand = CloneSubpageConfigCommand = {}));
