"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetAllHostTagsCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const api_1 = require("../../../api");
var GetAllHostTagsCommand;
(function (GetAllHostTagsCommand) {
    GetAllHostTagsCommand.url = api_1.REST_API.HOSTS.TAGS.GET;
    GetAllHostTagsCommand.TSQ_url = GetAllHostTagsCommand.url;
    GetAllHostTagsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.HOSTS_ROUTES.TAGS.GET, 'get', 'Get all existing host tags');
    GetAllHostTagsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            tags: zod_1.z.array(zod_1.z.string()),
        }),
    });
})(GetAllHostTagsCommand || (exports.GetAllHostTagsCommand = GetAllHostTagsCommand = {}));
