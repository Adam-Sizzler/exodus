"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetAllTagsCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const api_1 = require("../../../api");
var GetAllTagsCommand;
(function (GetAllTagsCommand) {
    GetAllTagsCommand.url = api_1.REST_API.USERS.TAGS.GET;
    GetAllTagsCommand.TSQ_url = GetAllTagsCommand.url;
    GetAllTagsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.TAGS.GET, 'get', 'Get all existing user tags');
    GetAllTagsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            tags: zod_1.z.array(zod_1.z.string()),
        }),
    });
})(GetAllTagsCommand || (exports.GetAllTagsCommand = GetAllTagsCommand = {}));
