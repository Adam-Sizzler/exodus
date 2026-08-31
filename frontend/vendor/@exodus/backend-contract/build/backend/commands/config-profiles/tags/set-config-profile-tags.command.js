"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SetConfigProfileTagsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
var SetConfigProfileTagsCommand;
(function (SetConfigProfileTagsCommand) {
    SetConfigProfileTagsCommand.url = api_1.REST_API.CONFIG_PROFILES.TAGS.SET;
    SetConfigProfileTagsCommand.TSQ_url = SetConfigProfileTagsCommand.url;
    SetConfigProfileTagsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.CONFIG_PROFILES_ROUTES.TAGS.SET, 'patch', 'Set tags of Config Profile', { scope: 'set-tags', kind: 'write' });
    SetConfigProfileTagsCommand.RequestBodySchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
        tags: models_1.TagsSchema,
    });
    SetConfigProfileTagsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            uuid: zod_1.z.uuid(),
            tags: zod_1.z.array(zod_1.z.string()),
        }),
    });
})(SetConfigProfileTagsCommand || (exports.SetConfigProfileTagsCommand = SetConfigProfileTagsCommand = {}));
