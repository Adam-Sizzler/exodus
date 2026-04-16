"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetSubscriptionByShortUuidCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
var GetSubscriptionByShortUuidCommand;
(function (GetSubscriptionByShortUuidCommand) {
    GetSubscriptionByShortUuidCommand.url = api_1.REST_API.SUBSCRIPTION.GET;
    GetSubscriptionByShortUuidCommand.TSQ_url = GetSubscriptionByShortUuidCommand.url(':shortUuid');
    GetSubscriptionByShortUuidCommand.RequestSchema = zod_1.z.object({
        shortUuid: zod_1.z.string(),
    });
})(GetSubscriptionByShortUuidCommand || (exports.GetSubscriptionByShortUuidCommand = GetSubscriptionByShortUuidCommand = {}));
