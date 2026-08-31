"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SetExternalSquadTagsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
var SetExternalSquadTagsCommand;
(function (SetExternalSquadTagsCommand) {
    SetExternalSquadTagsCommand.url = api_1.REST_API.EXTERNAL_SQUADS.TAGS.SET;
    SetExternalSquadTagsCommand.TSQ_url = SetExternalSquadTagsCommand.url;
    SetExternalSquadTagsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.EXTERNAL_SQUADS_ROUTES.TAGS.SET, 'patch', 'Set tags of External Squad', { scope: 'set-tags', kind: 'write' });
    SetExternalSquadTagsCommand.RequestBodySchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
        tags: models_1.TagsSchema,
    });
    SetExternalSquadTagsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            uuid: zod_1.z.uuid(),
            tags: zod_1.z.array(zod_1.z.string()),
        }),
    });
})(SetExternalSquadTagsCommand || (exports.SetExternalSquadTagsCommand = SetExternalSquadTagsCommand = {}));
