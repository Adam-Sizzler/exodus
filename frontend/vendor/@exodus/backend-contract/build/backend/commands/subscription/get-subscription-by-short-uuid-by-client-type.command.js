"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetSubscriptionByShortUuidByClientTypeCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var GetSubscriptionByShortUuidByClientTypeCommand;
(function (GetSubscriptionByShortUuidByClientTypeCommand) {
    GetSubscriptionByShortUuidByClientTypeCommand.url = api_1.REST_API.SUBSCRIPTION.GET;
    GetSubscriptionByShortUuidByClientTypeCommand.TSQ_url = GetSubscriptionByShortUuidByClientTypeCommand.url(':shortUuid');
    GetSubscriptionByShortUuidByClientTypeCommand.RequestParamSchema = zod_1.z.object({
        shortUuid: zod_1.z.string(),
        clientType: zod_1.z.enum(constants_1.REQUEST_TEMPLATE_TYPE, {
            error: 'Invalid client type.'
        }),
    });

})(GetSubscriptionByShortUuidByClientTypeCommand || (exports.GetSubscriptionByShortUuidByClientTypeCommand = GetSubscriptionByShortUuidByClientTypeCommand = {}));
