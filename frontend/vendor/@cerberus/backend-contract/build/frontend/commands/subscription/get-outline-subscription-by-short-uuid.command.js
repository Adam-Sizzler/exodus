"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetOutlineSubscriptionByShortUuidCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
var GetOutlineSubscriptionByShortUuidCommand;
(function (GetOutlineSubscriptionByShortUuidCommand) {
    GetOutlineSubscriptionByShortUuidCommand.url = api_1.REST_API.SUBSCRIPTION.GET;
    GetOutlineSubscriptionByShortUuidCommand.TSQ_url = GetOutlineSubscriptionByShortUuidCommand.url(':shortUuid');
    GetOutlineSubscriptionByShortUuidCommand.RequestSchema = zod_1.z.object({
        shortUuid: zod_1.z.string(),
        type: zod_1.z.optional(zod_1.z.string()),
        encodedTag: zod_1.z.optional(zod_1.z.string()),
    });
})(GetOutlineSubscriptionByShortUuidCommand || (exports.GetOutlineSubscriptionByShortUuidCommand = GetOutlineSubscriptionByShortUuidCommand = {}));
