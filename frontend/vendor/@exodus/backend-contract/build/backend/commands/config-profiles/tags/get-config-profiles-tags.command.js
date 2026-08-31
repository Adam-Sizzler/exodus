"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetConfigProfilesTagsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var GetConfigProfilesTagsCommand;
(function (GetConfigProfilesTagsCommand) {
    GetConfigProfilesTagsCommand.url = api_1.REST_API.CONFIG_PROFILES.TAGS.GET;
    GetConfigProfilesTagsCommand.TSQ_url = GetConfigProfilesTagsCommand.url;
    GetConfigProfilesTagsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.CONFIG_PROFILES_ROUTES.TAGS.GET, 'get', 'Get tags of Config Profiles', { scope: 'list-tags', kind: 'read' });
    GetConfigProfilesTagsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            tags: zod_1.z.array(zod_1.z.string()),
        }),
    });
})(GetConfigProfilesTagsCommand || (exports.GetConfigProfilesTagsCommand = GetConfigProfilesTagsCommand = {}));
