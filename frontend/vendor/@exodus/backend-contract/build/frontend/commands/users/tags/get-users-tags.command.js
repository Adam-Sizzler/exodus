"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetUsersTagsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var GetUsersTagsCommand;
(function (GetUsersTagsCommand) {
    GetUsersTagsCommand.url = api_1.REST_API.USERS.TAGS.GET;
    GetUsersTagsCommand.TSQ_url = GetUsersTagsCommand.url;
    GetUsersTagsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.TAGS.GET, 'get', 'Get users tags', { scope: 'list-tags', kind: 'read' });
    GetUsersTagsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            tags: zod_1.z.array(zod_1.z.string()),
        }),
    });
})(GetUsersTagsCommand || (exports.GetUsersTagsCommand = GetUsersTagsCommand = {}));
