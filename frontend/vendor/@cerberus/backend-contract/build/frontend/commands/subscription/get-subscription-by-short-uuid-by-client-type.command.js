"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetSubscriptionByShortUuidByClientTypeCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const api_1 = require("../../api");
var GetSubscriptionByShortUuidByClientTypeCommand;
(function (GetSubscriptionByShortUuidByClientTypeCommand) {
    GetSubscriptionByShortUuidByClientTypeCommand.url = api_1.REST_API.SUBSCRIPTION.GET;
    GetSubscriptionByShortUuidByClientTypeCommand.TSQ_url = GetSubscriptionByShortUuidByClientTypeCommand.url(':shortUuid');
    GetSubscriptionByShortUuidByClientTypeCommand.RequestSchema = zod_1.z.object({
        shortUuid: zod_1.z.string(),
        clientType: zod_1.z.nativeEnum(constants_1.REQUEST_TEMPLATE_TYPE, {
            message: 'Invalid client type.',
        }),
    });
})(GetSubscriptionByShortUuidByClientTypeCommand || (exports.GetSubscriptionByShortUuidByClientTypeCommand = GetSubscriptionByShortUuidByClientTypeCommand = {}));
