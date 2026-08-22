"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SyncSnippetCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var SyncSnippetCommand;
(function (SyncSnippetCommand) {
    SyncSnippetCommand.url = api_1.REST_API.SNIPPETS.ACTIONS.SYNC;
    SyncSnippetCommand.TSQ_url = SyncSnippetCommand.url;
    SyncSnippetCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SNIPPETS_ROUTES.ACTIONS.SYNC, 'post', 'Sync snippet to affected config profiles', { scope: 'sync', kind: 'write' }, 'Trigger the sync of a snippet to all config profiles that reference it. Nodes which use affected config profiles will be restarted.');
    SyncSnippetCommand.RequestBodySchema = zod_1.z.object({
        name: zod_1.z
            .string()
            .min(2, 'Name must be at least 2 characters')
            .max(255, 'Name must be less than 255 characters')
            .regex(/^[A-Za-z0-9_\s-]+$/, 'Name can only contain letters, numbers, underscores, dashes and spaces'),
    });
})(SyncSnippetCommand || (exports.SyncSnippetCommand = SyncSnippetCommand = {}));
