"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetSubpageConfigByShortUuidCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var GetSubpageConfigByShortUuidCommand;
(function (GetSubpageConfigByShortUuidCommand) {
    GetSubpageConfigByShortUuidCommand.url = api_1.REST_API.SUBSCRIPTIONS.SUBPAGE.GET_CONFIG;
    GetSubpageConfigByShortUuidCommand.TSQ_url = GetSubpageConfigByShortUuidCommand.url(':shortUuid');
    GetSubpageConfigByShortUuidCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTIONS_ROUTES.SUBPAGE.GET_CONFIG(':shortUuid'), 'get', 'Get Subpage Config by Short UUID');
    GetSubpageConfigByShortUuidCommand.RequestSchema = zod_1.z.object({
        shortUuid: zod_1.z.string(),
    });
    GetSubpageConfigByShortUuidCommand.RequestBodySchema = zod_1.z.object({
        requestHeaders: zod_1.z.record(zod_1.z.string(), zod_1.z.string()),
    });
    GetSubpageConfigByShortUuidCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            subpageConfigUuid: zod_1.z.string().uuid().nullable(),
            webpageAllowed: zod_1.z.boolean(),
        }),
    });
})(GetSubpageConfigByShortUuidCommand || (exports.GetSubpageConfigByShortUuidCommand = GetSubpageConfigByShortUuidCommand = {}));
