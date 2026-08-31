"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SetInternalSquadTagsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
var SetInternalSquadTagsCommand;
(function (SetInternalSquadTagsCommand) {
    SetInternalSquadTagsCommand.url = api_1.REST_API.INTERNAL_SQUADS.TAGS.SET;
    SetInternalSquadTagsCommand.TSQ_url = SetInternalSquadTagsCommand.url;
    SetInternalSquadTagsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INTERNAL_SQUADS_ROUTES.TAGS.SET, 'patch', 'Set tags of Internal Squad', { scope: 'set-tags', kind: 'write' });
    SetInternalSquadTagsCommand.RequestBodySchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
        tags: models_1.TagsSchema,
    });
    SetInternalSquadTagsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            uuid: zod_1.z.uuid(),
            tags: zod_1.z.array(zod_1.z.string()),
        }),
    });
})(SetInternalSquadTagsCommand || (exports.SetInternalSquadTagsCommand = SetInternalSquadTagsCommand = {}));
